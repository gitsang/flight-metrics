package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gitsang/flight-metrics/pkg/configer"
	"github.com/gitsang/flight-metrics/pkg/syntax"
	"github.com/krisukox/google-flights-api/flights"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type FlightConfig struct {
	Cron             string `json:"cron" yaml:"cron"`
	OutboundDate     string `json:"outboundDate" yaml:"outboundDate"`
	ReturnDate       string `json:"returnDate" yaml:"returnDate"`
	DepartureAirport string `json:"departureAirport" yaml:"departureAirport"`
	ArrivalAirport   string `json:"arrivalAirport" yaml:"arrivalAirport"`
}

type Config struct {
	Currency string         `json:"currency" yaml:"currency"`
	Lang     string         `json:"lang" yaml:"lang"`
	Flights  []FlightConfig `json:"flights"`
}

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "This is an example",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var rootFlags = struct {
	Config string
}{}

var cfger *configer.Configer

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootFlags.Config, "config", "c",
		"./configs/config.yml", "config file")

	cfger = configer.New(
		configer.WithTemplate((*Config)(nil)),
		configer.WithEnvBind(
			configer.WithEnvPrefix("FLIGHTM"),
			configer.WithEnvDelim("_"),
		),
		configer.WithFlagBind(
			configer.WithCommand(rootCmd),
			configer.WithFlagPrefix("flightm"),
			configer.WithFlagDelim("."),
		),
	)
}

func run() {
	var c Config
	err := cfger.Load(&c, rootFlags.Config)
	if err != nil {
		panic(err)
	}

	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	// croner
	cronParser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour |
			cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	croner := cron.New(cron.WithParser(cronParser))
	croner.Start()
	defer croner.Stop()

	// session
	session, err := flights.New()
	if err != nil {
		logger.Error("Error creating session", slog.Any("err", err))
		return
	}

	// currency
	currency, err := currency.ParseISO(c.Currency)
	if err != nil {
		logger.Error("Error parsing currency", slog.Any("err", err))
		return
	}

	// lang
	lang, err := language.Parse(c.Lang)
	if err != nil {
		logger.Error("Error parsing language", slog.Any("err", err))
		return
	}

	go func() {
		for _, flightConfig := range c.Flights {
			flightConfig := flightConfig
			logger := logger.With(slog.Any("flightConfig", flightConfig))
			schedule, err := cronParser.Parse(flightConfig.Cron)
			if err != nil {
				logger.Error("Error parsing cron", slog.Any("err", err))
				continue
			}

			entryID := croner.Schedule(schedule, cron.FuncJob(func() {
				var (
					ctx       = context.Background()
					entryTime = time.Now()
					logger    = logger.With(
						slog.String("executedAt", time.Now().Format(time.RFC3339)),
					)
				)
				defer func() {
					logger = logger.With(
						slog.String("duration", time.Since(entryTime).String()),
					)
					logger.Info("Search finished")
				}()

				outboundDate, err := time.Parse(time.DateOnly, flightConfig.OutboundDate)
				if err != nil {
					logger = logger.With(slog.Any("err", err))
					return
				}

				returnDate, err := time.Parse(time.DateOnly, flightConfig.ReturnDate)
				if err != nil {
					logger = logger.With(slog.Any("err", err))
					return
				}

				args := flights.Args{
					Date:        outboundDate,
					ReturnDate:  returnDate,
					SrcCities:   nil,
					SrcAirports: []string{flightConfig.DepartureAirport},
					DstCities:   nil,
					DstAirports: []string{flightConfig.ArrivalAirport},
					Options: flights.Options{
						Travelers: flights.Travelers{Adults: 1},
						Currency:  currency,
						Stops:     flights.AnyStops,
						Class:     flights.Economy,
						TripType:  flights.RoundTrip,
						Lang:      lang,
					},
				}

				offers, _, err := session.GetOffers(ctx, args)
				if err != nil {
					logger.Error("Error getting offers", slog.Any("err", err))
					return
				}

				url, err := session.SerializeURL(ctx, args)
				if err != nil {
					logger.Error("Error serializing url", slog.Any("err", err))
					return
				}
				logger.Info("flight search", slog.Any("url", url))

				for _, offer := range offers {
					outboundFlights := []Flight{}
					for _, flight := range offer.Flight {
						outboundFlights = append(outboundFlights,
							Flight{
								DepTime:      flight.DepTime,
								ArrTime:      flight.ArrTime,
								DepAirport:   flight.DepAirportCode,
								ArrAirport:   flight.ArrAirportCode,
								Duration:     flight.Duration.String(),
								AirlineName:  flight.AirlineName,
								FlightNumber: flight.FlightNumber,
								Airplane:     flight.Airplane,
							},
						)
					}
					returnFlights := []Flight{}
					for _, flight := range offer.ReturnFlight {
						returnFlights = append(returnFlights,
							Flight{
								DepTime:      flight.DepTime,
								ArrTime:      flight.ArrTime,
								DepAirport:   flight.DepAirportCode,
								ArrAirport:   flight.ArrAirportCode,
								Duration:     flight.Duration.String(),
								AirlineName:  flight.AirlineName,
								FlightNumber: flight.FlightNumber,
								Airplane:     flight.Airplane,
							},
						)
					}

					logger.Info("Found offer",
						slog.Group("outbound",
							slog.Any("flights", outboundFlights),
							slog.String("duration", offer.FlightDuration.String()),
						),
						slog.Group("return",
							slog.Any("flights", returnFlights),
							slog.String("duration", offer.ReturnFlightDuration.String()),
						),
						slog.String("price", fmt.Sprintf("%0.0f %s", offer.Price, c.Currency)),
					)

					outboundDurDay := DurDay(outboundFlights[0].DepTime.Truncate(24*time.Hour),
						outboundFlights[len(outboundFlights)-1].ArrTime)
					returnDurDay := DurDay(returnFlights[0].DepTime.Truncate(24*time.Hour),
						returnFlights[len(returnFlights)-1].ArrTime)
					RoundTripGauge.WithLabelValues(
						flightConfig.DepartureAirport,
						flightConfig.ArrivalAirport,
						flightConfig.OutboundDate,
						flightConfig.ReturnDate,
						// outbound
						outboundFlights[0].DepTime.Format(time.TimeOnly),
						outboundFlights[len(outboundFlights)-1].ArrTime.Format(time.TimeOnly)+
							syntax.If(outboundDurDay > 0, fmt.Sprintf("(+%d)", int64(outboundDurDay)), ""),
						offer.FlightDuration.String(),
						strings.Join(AllAirlineName(outboundFlights), " & "),
						strings.Join(AllMetadata(outboundFlights), "\n"),
						// return
						returnFlights[0].DepTime.Format(time.TimeOnly),
						returnFlights[len(returnFlights)-1].ArrTime.Format(time.TimeOnly)+
							syntax.If(returnDurDay > 0, fmt.Sprintf("(+%d)", int64(returnDurDay)), ""),
						offer.ReturnFlightDuration.String(),
						strings.Join(AllAirlineName(returnFlights), " & "),
						strings.Join(AllMetadata(returnFlights), "\n"),
					).Set(offer.Price)
				}
			}))

			logger.Info("Flight search scheduled", slog.Any("entryID", entryID))
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	logger.Info("Listening on :2112")
	http.ListenAndServe(":2112", nil)
}

func main() {
	rootCmd.Execute()
}

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	g "github.com/serpapi/google-search-results-golang"
)

var (
	RoundTripGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "serpapi",
		Subsystem:   "google_flights",
		Name:        "round_trip",
		Help:        "",
		ConstLabels: map[string]string{},
	}, []string{
		"departure_id", "arrival_id", "outbound_date", "return_date",
		"outbound_routes", "outbound_departure_time", "outbound_arrival_time",
		"return_routes", "return_departure_time", "return_arrival_time",
	})
)

func mask(in string) string {
	if len(in) < 3 {
		return in
	}
	return in[0:3] + "****" + in[len(in)-3:]
}

func Search(departureId, arrivalId, outboundDate, returnDate, apiKey string) error {
	// Set the parameters
	parameter := map[string]string{
		"engine":        "google_flights",
		"departure_id":  departureId,
		"arrival_id":    arrivalId,
		"outbound_date": outboundDate,
		"return_date":   returnDate,
		"currency":      "CNY",
		"hl":            "zh-cn",
	}

	search := g.NewGoogleSearch(parameter, apiKey)
	outboundResultRaw, err := search.GetJSON()
	if err != nil {
		return err
	}
	outboundJsonBytes, _ := json.Marshal(outboundResultRaw)
	var outboundResult Result
	err = json.Unmarshal(outboundJsonBytes, &outboundResult)
	if err != nil {
		return err
	}

	for _, outboundFlights := range outboundResult.BestFlights {
		parameter["departure_token"] = outboundFlights.DepartureToken
		search := g.NewGoogleSearch(parameter, apiKey)
		returnResultRaw, err := search.GetJSON()
		if err != nil {
			return err
		}
		returnJsonBytes, _ := json.MarshalIndent(returnResultRaw, "", "  ")
		var returnResult Result
		err = json.Unmarshal(returnJsonBytes, &returnResult)
		if err != nil {
			return err
		}

		for _, returnFlights := range returnResult.BestFlights {
			outboundRoutes := []string{}
			outboundDepartureTime := outboundFlights.Flights[0].DepartureAirport.Time
			outboundArrivalTime := outboundFlights.Flights[len(outboundFlights.Flights)-1].ArrivalAirport.Time
			for _, outboundFlight := range outboundFlights.Flights {
				outboundRoutes = append(outboundRoutes,
					fmt.Sprintf("%s %s(%s) --%s %s(%s)--> %s %s(%s)",
						outboundFlight.DepartureAirport.Time, outboundFlight.DepartureAirport.Name, outboundFlight.DepartureAirport.Id,
						outboundFlight.Airline, outboundFlight.FlightNumber,
						(time.Duration(outboundFlight.Duration)*time.Minute).String(),
						outboundFlight.ArrivalAirport.Time, outboundFlight.ArrivalAirport.Name, outboundFlight.ArrivalAirport.Id,
					),
				)
			}

			returnRoutes := []string{}
			returnDepartureTime := returnFlights.Flights[0].DepartureAirport.Time
			returnArrivalTime := returnFlights.Flights[len(returnFlights.Flights)-1].ArrivalAirport.Time
			for _, returnFlight := range returnFlights.Flights {
				returnRoutes = append(returnRoutes,
					fmt.Sprintf("%s %s(%s) --%s %s(%s)--> %s %s(%s)",
						returnFlight.DepartureAirport.Time, returnFlight.DepartureAirport.Name, returnFlight.DepartureAirport.Id,
						returnFlight.Airline, returnFlight.FlightNumber,
						(time.Duration(returnFlight.Duration)*time.Minute).String(),
						returnFlight.ArrivalAirport.Time, returnFlight.ArrivalAirport.Name, returnFlight.ArrivalAirport.Id,
					),
				)
			}

			RoundTripGauge.WithLabelValues(
				departureId, arrivalId, outboundDate, returnDate,
				strings.Join(outboundRoutes, " =>> "), outboundDepartureTime, outboundArrivalTime,
				strings.Join(returnRoutes, " =>> "), returnDepartureTime, returnArrivalTime,
			).Set(float64(returnFlights.Price))
		}
	}

	return nil
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", slog.Any("err", err))
	}
	searches := strings.Split(os.Getenv("FLIGHTM_SEARCHES"), ";")

	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	croner := cron.New()
	croner.Start()
	defer croner.Stop()

	go func() {
		for _, search := range searches {
			search = strings.TrimSpace(search)
			if search == "" {
				continue
			}
			parameters := strings.Split(search, ",")
			if len(parameters) != 6 {
				slog.Error("Error parsing search string", slog.String("search", search))
				continue
			}
			schedule, err := cronParser.Parse(parameters[0])
			if err != nil {
				slog.Error("Error parsing cron string", slog.Any("err", err), slog.String("search", search))
				continue
			}
			entryID := croner.Schedule(schedule, cron.FuncJob(func() {
				err = Search(parameters[1], parameters[2], parameters[3], parameters[4], parameters[5])
				if err != nil {
					slog.Error("Search failed", slog.String("search", search))
				} else {
					slog.Info("Search success", slog.String("search", search))
				}
			}))
			slog.Info("Scheduled", slog.Any("id", entryID),
				slog.String("cron", parameters[0]),
				slog.String("departure", parameters[1]),
				slog.String("arrival", parameters[2]),
				slog.String("outbound", parameters[3]),
				slog.String("return", parameters[4]),
				slog.String("api-key", mask(parameters[5])),
			)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

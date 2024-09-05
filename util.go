package main

import (
	"slices"
	"strings"
	"time"
)

func Mask(in string) string {
	if len(in) < 3 {
		return in
	}
	return in[0:3] + "****" + in[len(in)-3:]
}

func DurDay(start, end time.Time) int {
	start = start.Truncate(24 * time.Hour)
	end = end.Truncate(24 * time.Hour)
	diff := end.Sub(start)
	return int(diff.Hours() / 24)
}

func AllAirlineName(flights []Flight) []string {
	airlineMaps := make(map[string]struct{})
	for _, flight := range flights {
		airlineMaps[flight.AirlineName] = struct{}{}
	}
	airlineNames := make([]string, 0, len(airlineMaps))
	for k := range airlineMaps {
		airlineNames = append(airlineNames, k)
	}
	slices.Sort(airlineNames)
	return airlineNames
}

func AllAirlineCode(flights []Flight) []string {
	airlineCodeMaps := make(map[string]struct{})
	for _, flight := range flights {
		airlineCode := strings.Split(flight.FlightNumber, " ")[0]
		airlineCodeMaps[airlineCode] = struct{}{}
	}
	airlineCodes := make([]string, 0, len(airlineCodeMaps))
	for k := range airlineCodeMaps {
		airlineCodes = append(airlineCodes, k)
	}
	slices.Sort(airlineCodes)
	return airlineCodes
}

func AllFlightNumber(flights []Flight) []string {
	flightNumbers := make([]string, 0, len(flights))
	for _, flight := range flights {
		flightNumbers = append(flightNumbers, flight.FlightNumber)
	}
	return flightNumbers
}

func AllMetadata(flights []Flight) []string {
	metadatas := make([]string, 0, len(flights))
	for _, flight := range flights {
		metadatas = append(metadatas, flight.String())
	}
	return metadatas
}

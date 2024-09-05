package main

import (
	"encoding/json"
	"time"
)

type Flight struct {
	DepTime      time.Time
	ArrTime      time.Time
	DepAirport   string
	ArrAirport   string
	Duration     string
	AirlineName  string
	FlightNumber string
	Airplane     string
}

func (f Flight) String() string {
	jsonBytes, _ := json.Marshal(f)
	return string(jsonBytes)
}

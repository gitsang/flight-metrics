package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RoundTripGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "flight",
		Subsystem:   "google",
		Name:        "round_trip",
		Help:        "",
		ConstLabels: map[string]string{},
	}, []string{
		"departure_airport", "arrival_airport", "outbound_date", "return_date",
		"outbound_departure_time", "outbound_arrival_time", "outbound_duration", "outbound_airline_name", "outbound_flight_numbers",
		"return_departure_time", "return_arrival_time", "return_duration", "return_airline_name", "return_flight_numbers",
	})
)

package main

type SearchParameters struct {
	ArrivalId    string `json:"arrival_id,omitempty"`
	Currency     string `json:"currency,omitempty"`
	DepartureId  string `json:"departure_id,omitempty"`
	Engine       string `json:"engine,omitempty"`
	Gl           string `json:"gl,omitempty"`
	Hl           string `json:"hl,omitempty"`
	OutboundDate string `json:"outbound_date,omitempty"`
	ReturnDate   string `json:"return_date,omitempty"`
}

type SearchMetadata struct {
	CreatedAt        string  `json:"created_at,omitempty"`
	GoogleFlightsUrl string  `json:"google_flights_url,omitempty"`
	Id               string  `json:"id,omitempty"`
	JsonEndpoint     string  `json:"json_endpoint,omitempty"`
	PrettifyHtmlFile string  `json:"prettify_html_file,omitempty"`
	ProcessedAt      string  `json:"processed_at,omitempty"`
	RawHtmlFile      string  `json:"raw_html_file,omitempty"`
	Status           string  `json:"status,omitempty"`
	TotalTimeTaken   float64 `json:"total_time_taken,omitempty"`
}

type PriceInsights struct {
	LowestPrice       int     `json:"lowest_price,omitempty"`
	PriceHistory      [][]int `json:"price_history,omitempty"`
	PriceLevel        string  `json:"price_level,omitempty"`
	TypicalPriceRange []int   `json:"typical_price_range,omitempty"`
}

type Airport struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Time string `json:"time,omitempty"`
}

type FlightExtensions []string

type Flight struct {
	Airline          string           `json:"airline,omitempty"`
	Airplane         string           `json:"airplane,omitempty"`
	ArrivalAirport   Airport          `json:"arrival_airport,omitempty"`
	DepartureAirport Airport          `json:"departure_airport,omitempty"`
	Duration         int              `json:"duration,omitempty"`
	Extensions       FlightExtensions `json:"extensions,omitempty"`
	FlightNumber     string           `json:"flight_number,omitempty"`
}

type CarbonEmissions struct {
	DifferencePercent   int `json:"difference_percent,omitempty"`
	ThisFlight          int `json:"this_flight,omitempty"`
	TypicalForThisRoute int `json:"typical_for_this_route,omitempty"`
}

type Layover struct {
	Duration  int    `json:"duration,omitempty"`
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Overnight bool   `json:"overnight,omitempty"`
}

type AirlineFlight struct {
	AirlineLogo     string          `json:"airline_logo,omitempty"`
	CarbonEmissions CarbonEmissions `json:"carbon_emissions,omitempty"`
	DepartureToken  string          `json:"departure_token,omitempty"`
	Flights         []Flight        `json:"flights,omitempty"`
	Layovers        []Layover       `json:"layovers,omitempty"`
	Price           int             `json:"price,omitempty"`
	TotalDuration   int             `json:"total_duration,omitempty"`
	Type            string          `json:"type,omitempty"`
}

type Result struct {
	SearchParameters SearchParameters `json:"search_parameters,omitempty"`
	SearchMetadata   SearchMetadata   `json:"search_metadata,omitempty"`
	BestFlights      []AirlineFlight  `json:"best_flights,omitempty"`
	OtherFlights     []AirlineFlight  `json:"other_flights,omitempty"`
	PriceInsights    PriceInsights    `json:"price_insights,omitempty"`
}

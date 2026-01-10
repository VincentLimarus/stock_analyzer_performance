package dto

type TradeData struct {
	SerialCode    string       `json:"serialCode"`
	Status        string       `json:"status"`
	Code          string       `json:"code"`
	Lot           int          `json:"lot"`
	AveragePrice  AveragePrice `json:"averagePrice"`
	GainLoss      GainLoss     `json:"gainLoss"`
	DividenAmount *float64     `json:"dividenAmount"`
}

type GainLoss struct {
	Gross float64 `json:"gross"`
	Fee   float64 `json:"fee"`
	Net   float64 `json:"net"`
}

type AveragePrice struct {
	PercentageChange float64 `json:"percentageChange"`
	Buy              float64 `json:"buy"`
	Sell             float64 `json:"sell"`
}

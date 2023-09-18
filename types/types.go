package types

type FindPriceResponse struct {
	Price  float64 `json:"price"`
	Ticker string  `json:"ticker"`
}

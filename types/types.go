package types

type FindPriceResponse struct {
	Price float64 `json:"price"`
	Coin  string  `json:"coin"`
}

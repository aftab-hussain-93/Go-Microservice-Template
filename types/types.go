package types

import (
	"time"

	err "github.com/aftab-hussain-93/crypto-price-finder-microservice/error"
)

type ErrResponse struct {
	Error *err.Err `json:"error"`
	// the operation being performed, such as updateUser, deleteUser etc. Typically processed by the handler
	Operation string    `json:"operation"`
	RequestID int       `json:"requestId"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

type FindPriceResponse struct {
	Price float64 `json:"price"`
	Coin  string  `json:"coin"`
}

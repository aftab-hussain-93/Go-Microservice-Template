package types

import "time"

type ErrResponse struct {
	Error      string    `json:"error"`
	Code       int       `json:"code"`
	Timestamp  time.Time `json:"timestamp"`
	RequestID  string    `json:"requestId"`
	Service    string    `json:"service"`
	Resolution string    `json:"resolution"`
}

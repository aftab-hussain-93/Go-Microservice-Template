package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/aftab-hussain-93/crypto-price-finder-microservice/types"
)

var (
	ErrInvalid         = errors.New("invalid")
	ErrUnAuthenticated = errors.New("unauthenticated") // http unauth
	ErrForbidden       = errors.New("forbidden")       // forbidden
	ErrNotFound        = errors.New("notfound")
	ErrTimeout         = errors.New("timeout")
	ErrConflict        = errors.New("conflict")
	ErrDBFailure       = errors.New("dberror")  // internal server error
	ErrExternalFailure = errors.New("external") // internal server error
	ErrUnknown         = errors.New("unknown")  // internal server error
)

var ErrorToHTTPStatus = map[error]int{
	ErrInvalid:         http.StatusBadRequest,
	ErrUnAuthenticated: http.StatusUnauthorized,
	ErrForbidden:       http.StatusForbidden,
	ErrNotFound:        http.StatusNotFound,
	ErrTimeout:         http.StatusRequestTimeout,
	ErrConflict:        http.StatusConflict,
	ErrDBFailure:       http.StatusInternalServerError,
	ErrExternalFailure: http.StatusInternalServerError,
	ErrUnknown:         http.StatusInternalServerError,
}

func ConvertToHTTPError(ctx context.Context, err error) (status int, e types.ErrResponse) {
	e.Error = err.Error()
	e.Timestamp = time.Now()
	e.Service = "PriceFinderService"
	e.RequestID = ctx.Value(RequestID).(string)
	e.Code = http.StatusInternalServerError // default error
	status = http.StatusInternalServerError // default error

	return
}

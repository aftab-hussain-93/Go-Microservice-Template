package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	err "github.com/aftab-hussain-93/crypto-price-finder-microservice/error"
	"github.com/aftab-hussain-93/crypto-price-finder-microservice/types"
)

func (s *JSONAPIServer) errorHandler(ctx context.Context, e error) (status int, res types.ErrResponse) {
	if op := ctx.Value(Operation); op != nil {
		res.Operation = op.(string)
	}
	if reqid := ctx.Value(RequestID); reqid != nil {
		res.RequestID = reqid.(int)
	}
	if srv := ctx.Value(Service); srv != nil {
		res.Service = srv.(string)
	}
	res.Timestamp = time.Now()
	var myErr *err.Err
	if errors.As(e, &myErr) {
		// handled error
		res.Error = myErr
		status = myErr.GetHTTPStatusCode()
	} else {
		// default error
		res.Error = &err.Err{
			Code:       err.ErrInternal,
			Message:    "Unexpected error",
			Resolution: "Please try again later",
			Err:        e,
		}
	}
	if status == 0 {
		status = http.StatusInternalServerError
	}
	return
}

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aftab-hussain-93/crypto-price-finder-microservice/types"
)

type JSONAPIServer struct {
	svc  PriceFinder
	addr string
}

func NewJSONAPIServer(add string, svc PriceFinder) *JSONAPIServer {
	return &JSONAPIServer{
		addr: add,
		svc:  svc,
	}
}

func (s *JSONAPIServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", makeErrorHandler(s.handleFindPrice))

	err := http.ListenAndServe(s.addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *JSONAPIServer) handleFindPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	key := r.URL.Query().Get("ticker")

	price, err := s.svc.FindPrice(ctx, key)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, &types.FindPriceResponse{
		Price:  price,
		Ticker: key,
	})
}

type ReqID string

func makeErrorHandler(fn func(context.Context, http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ReqID("requestID"), 10001)
		err := fn(ctx, w, r)
		if err != nil {
			err := writeJSON(w, 400, map[string]any{"error": err.Error()})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`Internal server error`))
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

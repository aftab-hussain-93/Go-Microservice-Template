package main

import (
	"context"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"

	"github.com/aftab-hussain-93/crypto-price-finder-microservice/types"
)

type JSONAPIServer struct {
	svc PriceFinder
	cfg JSONAPIServerConf
}

type JSONAPIServerConf struct {
	Debug   bool
	Address string
}

func (s *JSONAPIServer) getAddr() string {
	return s.cfg.Address
}

func NewJSONAPIServer(cfg JSONAPIServerConf, svc PriceFinder) *JSONAPIServer {
	return &JSONAPIServer{
		svc: svc,
		cfg: cfg,
	}
}

func (s *JSONAPIServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger", s.handleAPIDoc)

	mux.HandleFunc("/api/v1/prices/", s.makeErrorHandler(s.handleFindPrice))

	err := http.ListenAndServe(s.getAddr(), mux)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *JSONAPIServer) handleAPIDoc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`Hello World!`))
}

func (s *JSONAPIServer) handleFindPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	escapedPath := strings.TrimSpace(strings.TrimPrefix(r.URL.EscapedPath(), "/api/v1/prices/"))
	path := strings.SplitN(escapedPath, "/", 1)
	coin := path[0]
	price, err := s.svc.FindPrice(ctx, coin)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, &types.FindPriceResponse{
		Price: price,
		Coin:  coin,
	})
}

func (s *JSONAPIServer) makeErrorHandler(fn func(context.Context, http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestID, rand.Intn(999999))
		slog.InfoContext(ctx, "Request received", slog.Any("path", r.URL.Path))
		err := fn(ctx, w, r)
		if err != nil {
			status, jsonErr := s.errorHandler(ctx, err)
			err := writeJSON(w, status, jsonErr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`Internal server error`))
			}
		}
	}
}

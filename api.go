package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"strings"

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
	mux.HandleFunc("/", makeErrorHandler(s._routeIndex))

	err := http.ListenAndServe(s.addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *JSONAPIServer) handleAPIDoc(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`Hello World!`))
	return nil
}

func (s *JSONAPIServer) handleFindPrice(ctx context.Context, coin string, w http.ResponseWriter, r *http.Request) error {
	price, err := s.svc.FindPrice(ctx, coin)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, &types.FindPriceResponse{
		Price: price,
		Coin:  coin,
	})
}

func (s *JSONAPIServer) _routeIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	escapedPath := strings.TrimSpace(strings.TrimPrefix(r.URL.EscapedPath(), "/"))
	if escapedPath == "" {
		// doc router
		return s.handleAPIDoc(ctx, w, r)
	}
	path := strings.Split(escapedPath, "/")
	if len(path) == 1 {
		// coin data router
		return s.handleFindPrice(ctx, path[0], w, r)
	}
	return s._handleNotFound(ctx, w, r)
}

func (s *JSONAPIServer) _handleNotFound(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`404 not found`))
	return nil
}

func makeErrorHandler(fn func(context.Context, http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestID, 10001)
		slog.InfoContext(ctx, "Request received", slog.Any("path", r.URL.Path))
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

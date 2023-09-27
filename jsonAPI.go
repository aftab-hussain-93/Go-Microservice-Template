package main

import (
	"context"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aftab-hussain-93/crypto-price-finder-microservice/types"
)

type JSONAPIServer struct {
	svc PriceFinder
	cfg *JSONAPIServerConf
}

type JSONAPIServerConf struct {
	Debug   bool
	Address string
}

func (s *JSONAPIServer) getAddr() string {
	return s.cfg.Address
}

func NewJSONAPIServer(cfg *JSONAPIServerConf, svc PriceFinder) *JSONAPIServer {
	return &JSONAPIServer{
		svc: svc,
		cfg: cfg,
	}
}

func (s *JSONAPIServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger", s.handleAPIDoc)

	mux.HandleFunc("/api/v1/prices/", s.makeErrorHandler(s.handleFindPrice))

	srv := http.Server{
		Addr:         s.getAddr(),
		ErrorLog:     log.Default(), // Integrate with slog and write to specific out
		Handler:      mux,
		ReadTimeout:  5 * time.Second, // default values
		WriteTimeout: 5 * time.Second, // default values
	}
	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			// errored out
			slog.Error("listenAndServe error - ", slog.String("err", err.Error()))
		}
		close(done)
	}()

	select {
	case q := <-quit:
		// terminate called
		slog.Error("interrupt received, shutting down server", q)
	case <-done:
		// server errored out
	}

	ctx, shutdown := context.WithTimeout(context.Background(), 7*time.Second)
	defer shutdown()

	err := srv.Shutdown(ctx)
	if err != nil {
		slog.Error("errored out while shutting down server")
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

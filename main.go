package main

import (
	"log/slog"
	"os"
	"time"
)

func main() {
	pc := &priceFinder{
		m: &mockPriceFinder{},
	}
	// init global slog logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
	svc := NewLoggingWrapper(pc)
	cfg := &JSONAPIServerConf{
		Debug:          true,
		Address:        ":3000",
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    4 * time.Second,
		HandlerTimeout: 9 * time.Second,
	}

	server := NewJSONAPIServer(cfg, svc)

	server.Run()
}

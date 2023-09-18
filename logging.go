package main

import (
	"context"
	"log/slog"
	"time"
)

type logger struct {
	next PriceFinder
}

func NewLogger(next PriceFinder) PriceFinder {
	return &logger{
		next: next,
	}
}

func (l *logger) FindPrice(ctx context.Context, key string) (price float64, err error) {
	defer func(begin time.Time) {
		slog.Info("FindPrice", slog.Duration("latency", time.Since(begin)), slog.Float64("price", price), slog.Any("err", err))
	}(time.Now())

	return l.next.FindPrice(ctx, key)
}

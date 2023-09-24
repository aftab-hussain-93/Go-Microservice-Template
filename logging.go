package main

import (
	"context"
	"log/slog"
	"time"
)

type lgr struct {
	next PriceFinder
}

func NewLoggingWrapper(next PriceFinder) PriceFinder {
	return &lgr{
		next: next,
	}
}

func (l *lgr) FindPrice(ctx context.Context, key string) (price float64, err error) {
	defer func(begin time.Time, ctx context.Context) {
		slog.InfoContext(ctx,
			"FindPrice",
			slog.Duration("latency",
				time.Since(begin)),
			slog.String("coin", key),
			slog.Float64("price", price),
			slog.Any("err", err))

	}(time.Now(), ctx)

	return l.next.FindPrice(ctx, key)
}

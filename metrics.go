package main

import (
	"context"
	"fmt"
)

type metrics struct {
	next PriceFinder
}

func (m *metrics) FindPrice(ctx context.Context, key string) (float64, error) {
	fmt.Println("logging to prometheus")
	return m.next.FindPrice(ctx, key)
}

func NewMetric(next PriceFinder) PriceFinder {
	return &metrics{
		next: next,
	}
}

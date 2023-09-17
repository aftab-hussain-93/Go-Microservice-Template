package main

import (
	"context"
	"fmt"
)

type PriceChecker interface {
	CheckPrice(context.Context, string) (float64, error)
}

type priceChecker struct {
	m MockPriceChecker
}

var prices = map[string]float64{
	"ETH": 2_000,
	"BTC": 20_000,
}

func (s *priceChecker) CheckPrice(ctx context.Context, key string) (float64, error) {
	return s.m.CheckPrice(ctx, key)
}

type MockPriceChecker struct {
}

func (m *MockPriceChecker) CheckPrice(ctx context.Context, key string) (float64, error) {
	if price, ok := prices[key]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("price not found for key (%s)", key)
}

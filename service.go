package main

import (
	"context"
	"fmt"
)

type ContextKey string

const (
	UsrCtxKey ContextKey = "user" // Logged in user
	RequestID ContextKey = "requestID"
	Operation ContextKey = "operation"
	Service   ContextKey = "service"
)

type PriceFinder interface {
	FindPrice(context.Context, string) (float64, error)
}

type priceFinder struct {
	m *mockPriceFinder
}

var prices = map[string]float64{
	"ETH": 2_000,
	"BTC": 20_000,
}

func (s *priceFinder) FindPrice(ctx context.Context, key string) (float64, error) {
	return s.m.FindPrice(ctx, key)
}

type mockPriceFinder struct {
}

func (m *mockPriceFinder) FindPrice(ctx context.Context, key string) (float64, error) {
	if price, ok := prices[key]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("price not found for key (%s)", key)
}

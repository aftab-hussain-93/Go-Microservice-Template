package main

import (
	"github.com/redis/go-redis/v9"
)

func New(addrs string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       addrs,
		ClientName: "pricefinder",
	})
}

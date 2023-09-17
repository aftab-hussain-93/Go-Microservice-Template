package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	pc := &priceFinder{}
	lg := NewLogger(pc)

	cl := NewMetric(lg)

	price, err := cl.FindPrice(context.Background(), "GHT")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(price)
}

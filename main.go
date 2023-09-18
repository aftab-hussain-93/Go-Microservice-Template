package main

func main() {
	pc := &priceFinder{
		m: &mockPriceFinder{},
	}
	svc := NewMetric(NewLogger(pc))

	server := NewJSONAPIServer(":3000", svc)

	server.Run()

	// price, err := svc.FindPrice(context.Background(), "GHT")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(price)
}

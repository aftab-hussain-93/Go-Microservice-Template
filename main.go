package main

func main() {
	pc := &priceFinder{
		m: &mockPriceFinder{},
	}
	svc := NewMetric(NewLoggingWrapper(pc))
	cfg := JSONAPIServerConf{
		Debug:   true,
		Address: ":3000",
	}

	server := NewJSONAPIServer(cfg, svc)

	server.Run()
}

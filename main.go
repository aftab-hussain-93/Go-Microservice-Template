package main

import "time"

func main() {
	pc := &priceFinder{
		m: &mockPriceFinder{},
	}
	svc := NewLoggingWrapper(pc)
	cfg := &JSONAPIServerConf{
		Debug:        true,
		Address:      ":3000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	server := NewJSONAPIServer(cfg, svc)

	server.Run()
}

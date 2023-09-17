package main

type server struct {
	svc PriceFinder
}

func NewTransport() *server {
	return &server{}
}

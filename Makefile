build:
	go build -o ./bin/pricefinder .
run: build
	./bin/pricefinder
lint:
	golangci-lint run ./...
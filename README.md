# crypto-price-finder-microservice

A production ready microservice that can be used to fetch the price of a particular crypto currency. Comes baked with unit tests, EFK, prometheus support, e2e tests and a robust CI-CD pipeline using Github Actions.

# setup

- husky needs to be installed (https://github.com/go-courier/husky)
- golangci-lint needs to be installed (https://github.com/golangci/golangci-lint)

# todo

- Support input and output coin
- Add timeout middleware
- Use luna client for http requests
- Use HTTP request coalescing in Luna client
- Use kafka for logging
- Implement Open API spec
- Implement GRPC
- Add persistance for rate limiting based on incoming IP address, logged in user data etc.

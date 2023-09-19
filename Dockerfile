FROM golang:1.21

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
# COPY go.mod ./
RUN go mod download && go mod verify

COPY . ./
RUN make build

EXPOSE 3000
CMD ["./bin/pricefinder"]
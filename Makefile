.PHONY: build test lint run tidy
     
build:
  go build -o bin/pricemon ./cmd/pricemon

test:
  go test -race -count=1 ./...

lint:
  golangci-lint run

run:
  go run ./cmd/pricemon

tidy:
  go mod tidy

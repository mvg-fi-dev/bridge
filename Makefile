.PHONY: run test fmt

run:
	PORT?=8080
	go run ./cmd/bridge-api

test:
	go test ./...

fmt:
	gofmt -w .

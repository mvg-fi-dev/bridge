.PHONY: run test fmt

run:
	PORT?=8080
	go run ./cmd/bridge-api

worker:
	go run ./cmd/bridge-worker

test:
	go test ./...

fmt:
	gofmt -w .

.PHONY: build

lint:
	golangci-lint run

format:
	go fix ./...
	go fmt ./...

test:
	go test ./... -v

build: format lint test
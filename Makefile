.PHONY: build

lint:
	golangci-lint run

format:
	go fix ./...
	go fmt ./...

test:
	go test $(shell go list ./... | grep -v /integrations/) -v

integration:
	go test ./integrations/... -v -count=1 -timeout=120s

build: format lint test
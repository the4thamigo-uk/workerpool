GOPATH=$(shell go env GOPATH)

.PHONY: deps
deps:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.31.0

.PHONY: all
all: test build

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run ./...

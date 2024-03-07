.DEFAULT_GOAL := build

build: deps vet
	@go build github.com/siwei-luo/brot
.PHONY: build

deps:
	@go get -u
	@go mod tidy
.PHONY: deps

fmt:
	@go fmt ./...
.PHONY: fmt

vet: fmt
	@go vet ./...
.PHONY: vet

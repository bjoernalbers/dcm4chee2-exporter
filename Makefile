VERSION := $(shell git describe --tags | tr -d v)

.PHONY: all
all:
	go test ./...
	go build -ldflags '-X main.version=$(VERSION)'

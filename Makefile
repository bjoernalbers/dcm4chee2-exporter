.PHONY: all
all:
	go test ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

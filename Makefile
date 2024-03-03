.PHONY: test build

VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS := -X main.version=$(VERSION)

test:
	go test ./... -v

build: test
	go build -o ./bin/envLoader -ldflags "${LDFLAGS}" *.go

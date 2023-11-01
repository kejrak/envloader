.PHONY: test build

test:
	go test ./... -v

build: test
	go build -o envLoader -ldflags "-X main.version=0.1.0" *.go

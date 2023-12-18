.PHONY: test build

test:
	go test ./... -v

build: test
	go build -o ./bin/envLoader -ldflags "-X main.version=0.4.0" *.go

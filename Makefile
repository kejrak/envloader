.PHONY: test build

test:
	go test ./... -v

build: test
	go build -o ./bin/envLoader -ldflags "-X main.version=0.5.0" *.go

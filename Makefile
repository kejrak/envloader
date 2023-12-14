.PHONY: test build

test:
	go test ./... -v

build: test
	go build -o ./mock/envLoader -ldflags "-X main.version=0.3.0" *.go

build: test lint compile

test:
	go test

lint:
	golangci-lint run

compile:
	go build -ldflags "-s -w" ./cmd/ods

.PHONY: clean

clean:
	go clean -r -cache -testcache



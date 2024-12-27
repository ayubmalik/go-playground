build: test compile

test:
	go test

compile:
	go build -ldflags "-s -w" ./cmd/ods

.PHONY: clean

clean:
	go clean -r -cache -testcache



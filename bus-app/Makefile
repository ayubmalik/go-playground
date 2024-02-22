VERSION=$(shell git describe --tags --long)
LDFLAGS="-s -w -X main.version=${VERSION}"
BINARY=bus-app

build:
	GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o $(BINARY) ./

test: clean
	go test ./...

clean:
	@rm -rf dist
	@go clean -testcache


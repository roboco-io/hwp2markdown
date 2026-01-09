.PHONY: build test lint clean release install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"
BINARY := hwp2markdown

build:
	go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/hwp2markdown

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/ dist/

install: build
	cp bin/$(BINARY) $(GOPATH)/bin/

release: clean
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-x64.exe ./cmd/hwp2markdown
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-macos-x64 ./cmd/hwp2markdown
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-macos-arm64 ./cmd/hwp2markdown
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-x64 ./cmd/hwp2markdown

run: build
	./bin/$(BINARY) $(ARGS)

fmt:
	go fmt ./...

tidy:
	go mod tidy

deps:
	go mod download

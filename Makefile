VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
PREFIX  ?= $(HOME)/.local
LDFLAGS := -X main.version=$(VERSION)

.PHONY: build test install

build:
	go build -ldflags "$(LDFLAGS)" -o hometab .

test:
	go vet ./...
	go test ./...

install: build
	install -Dm755 hometab $(PREFIX)/bin/hometab
	./hometab gen-man $(PREFIX)/share/man/man1

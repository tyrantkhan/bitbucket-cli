.PHONY: build install test lint clean setup

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.0.0-dev")
BUILD_DATE := $(shell date +%Y-%m-%d)
LDFLAGS := -ldflags "-X github.com/tyrantkhan/bb/cmd.Version=$(VERSION) -X github.com/tyrantkhan/bb/cmd.BuildDate=$(BUILD_DATE)"

setup:
	mise install
	lefthook install

build:
	go build $(LDFLAGS) -o bb .

install:
	go install $(LDFLAGS) .

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -f bb

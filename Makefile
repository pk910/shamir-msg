# shamir-msg
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
VERSION := $(shell git rev-parse --short HEAD)

GOLDFLAGS += -X 'github.com/pk910/shamir-msg/buildinfo.BuildVersion="$(VERSION)"'
GOLDFLAGS += -X 'github.com/pk910/shamir-msg/buildinfo.Buildtime="$(BUILDTIME)"'
GOLDFLAGS += -X 'github.com/pk910/shamir-msg/buildinfo.BuildRelease="$(RELEASE)"'

.PHONY: all test clean

all: test build

test:
	go test ./...

build:
	@echo version: $(VERSION)
	go build -v -o bin/ -ldflags="-s -w $(GOLDFLAGS)" .

clean:
	rm -f bin/*

GO ?= go
OUTPUT_DIR ?= build
VERSION ?= $(shell git describe --tags --always)
VERSION := $(VERSION)
BUILD = $(GO) build -v -ldflags="-s -w -X main.version=$(VERSION)" -trimpath
BINARY := dlnapass dlnafind

.PHONY: all clean windows_x64 linux_x64 linux_arm64

default: all

all: windows_x64 linux_x64 linux_arm64

windows_x64:
	set -e; for bin in $(BINARY); do \
		GOOS=windows GOARCH=amd64 $(BUILD) -o ./$(OUTPUT_DIR)/$${bin}_$(VERSION)_$@.exe ./cmd/$$bin; \
	done

linux_x64:
	set -e; for bin in $(BINARY); do \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD) -o ./$(OUTPUT_DIR)/$${bin}_$(VERSION)_$@ ./cmd/$$bin; \
	done

linux_arm64:
	set -e; for bin in $(BINARY); do \
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(BUILD) -o ./$(OUTPUT_DIR)/$${bin}_$(VERSION)_$@ ./cmd/$$bin; \
	done

clean:
	rm -f $(OUTPUT_DIR)/*

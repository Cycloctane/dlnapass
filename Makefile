OUTPUT_DIR=build
BUILD=go build -v -ldflags="-s -w" -trimpath

.PHONY: all windows_x64 linux_x64 linux_static_x64 linux_static_arm64

default: all

all: windows_x64 linux_x64 linux_static_x64 linux_static_arm64

windows_x64:
	GOOS=windows GOARCH=amd64 $(BUILD) -o ./$(OUTPUT_DIR)/dlnapass_$@.exe ./cmd/dlnapass
	GOOS=windows GOARCH=amd64 $(BUILD) -o ./$(OUTPUT_DIR)/dlnafind_$@.exe ./cmd/dlnafind

linux_x64:
	GOOS=linux GOARCH=amd64 $(BUILD) -o ./$(OUTPUT_DIR)/dlnapass_$@ ./cmd/dlnapass
	GOOS=linux GOARCH=amd64 $(BUILD) -o ./$(OUTPUT_DIR)/dlnafind_$@ ./cmd/dlnafind

linux_static_x64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD) -o ./$(OUTPUT_DIR)/dlnapass_$@ ./cmd/dlnapass
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD) -o ./$(OUTPUT_DIR)/dlnafind_$@ ./cmd/dlnafind

linux_static_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(BUILD) -o ./$(OUTPUT_DIR)/dlnapass_$@ ./cmd/dlnapass
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(BUILD) -o ./$(OUTPUT_DIR)/dlnafind_$@ ./cmd/dlnafind

clean:
	rm -f $(OUTPUT_DIR)/*

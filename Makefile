OUTPUT_DIR=build
BINARY=dlnafind dlnapass

.PHONY: all $(BINARY) clean

all: $(BINARY)

$(BINARY):
	GOOS=windows GOARCH=amd64 go build -v -o ./$(OUTPUT_DIR)/$@_windows_x64.exe -ldflags="-s -w" ./cmd/$@
	GOOS=linux GOARCH=amd64 go build -v -o ./$(OUTPUT_DIR)/$@_linux_x64 -ldflags="-s -w" ./cmd/$@

clean:
	rm -f $(OUTPUT_DIR)/*

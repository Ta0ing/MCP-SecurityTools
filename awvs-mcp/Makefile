# AWVS MCP Makefile

.PHONY: build run-stdio run-sse clean

BUILD_DIR=./bin
BINARY_NAME=awvs-mcp

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

run-stdio: build
	$(BUILD_DIR)/$(BINARY_NAME) stdio

run-sse: build
	$(BUILD_DIR)/$(BINARY_NAME) sse --port 8080

clean:
	rm -rf $(BUILD_DIR)

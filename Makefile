# Build variables
BINARY_NAME=gollama
GO=go
GOFLAGS=-v
BUILD_DIR=bin
CMD_DIR=./cmd/cli

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)


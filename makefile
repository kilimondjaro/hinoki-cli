# Variables
APP_NAME := hinoki-planner
BIN_DIR := ./bin
OUTPUT_ARM := $(BIN_DIR)/$(APP_NAME)-darwin-arm64
OUTPUT_AMD := $(BIN_DIR)/$(APP_NAME)-darwin-amd64
UNIVERSAL_BIN := $(BIN_DIR)/$(APP_NAME)-darwin-universal
ARCHIVE := $(UNIVERSAL_BIN).tar.gz
SHA_FILE := $(ARCHIVE).sha256

# Ensure bin directory exists
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Build for macOS ARM
build-arm: $(BIN_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_ARM)

# Build for macOS AMD
build-amd: $(BIN_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_AMD)

# Combine ARM and AMD binaries into a universal binary
universal: build-arm build-amd
	lipo -create -output $(UNIVERSAL_BIN) $(OUTPUT_ARM) $(OUTPUT_AMD)

# Archive the universal binary
archive: universal
	tar -czvf $(ARCHIVE) -C $(BIN_DIR) $(APP_NAME)-darwin-universal

# Generate SHA256 checksum
checksum: archive
	shasum -a 256 $(ARCHIVE) > $(SHA_FILE)

# Clean up the build directory
clean:
	rm -rf $(BIN_DIR)

# Default target
all: build-arm build-amd universal archive checksum
.PHONY: run start build clean help

# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := windows
    BINARY_EXT := .exe
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        DETECTED_OS := linux
    endif
    ifeq ($(UNAME_S),Darwin)
        DETECTED_OS := darwin
    endif
    BINARY_EXT :=
endif

# Build configuration
APP_NAME := localsend
BUILD_DIR := build
BINARY_NAME := $(APP_NAME)$(BINARY_EXT)
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

# Go build flags for maximum performance
LDFLAGS := -ldflags="-s -w"
BUILD_FLAGS := -trimpath

help:
	@echo "Available targets:"
	@echo "  make run     - Run application natively (go run)"
	@echo "  make build   - Build optimized binary to $(BUILD_DIR)/"
	@echo "  make start   - Run the built binary"
	@echo "  make clean   - Remove build directory"
	@echo ""
	@echo "Detected OS: $(DETECTED_OS)"

run:
	@echo "Running application natively..."
	@go run main.go

build:
	@echo "Building for $(DETECTED_OS)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_PATH) main.go
	@echo "Build complete: $(BINARY_PATH)"

start:
	@if [ ! -f "$(BINARY_PATH)" ]; then \
		echo "Binary not found. Building first..."; \
		$(MAKE) build; \
	fi
	@echo "Starting $(BINARY_NAME)..."
	@./$(BINARY_PATH)

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

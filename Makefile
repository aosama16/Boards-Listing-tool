APP_NAME := boards_merger
BUILD_DIR := ./build

.PHONY: all build clean test

ifeq ($(OS),Windows_NT)
    RM = rmdir /Q /S
    MKDIR = mkdir $(subst /,\,$1)
else
    RM = rm -rf
    MKDIR = mkdir -p $1
endif

all: build

build:
	$(MKDIR) $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go

clean:
	$(RM) $(BUILD_DIR)

test:
	go test -v ./...

cov: build
	go test -coverprofile=$(BUILD_DIR)/coverage.out ./...
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
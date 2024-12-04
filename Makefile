CLI_APP_NAME=cli_boards_merger
WEB_APP_NAME=web_boards_merger
BUILD_DIR=./build

.PHONY: all build clean test

ifeq ($(OS),Windows_NT)
    RM = rmdir /Q /S
    MKDIR = mkdir $(subst /,\,$1)
else
    RM = rm -rf
    MKDIR = mkdir -p $1
endif

all: build

build: build-cli build-web

build-cli:
	-$(MKDIR) $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(CLI_APP_NAME) ./cmd/cli/main.go

build-web:
	-$(MKDIR) $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(WEB_APP_NAME) ./cmd/web/main.go

clean:
	-$(RM) $(BUILD_DIR)

test:
	go test -v ./...

cov: build
	go test -coverprofile=$(BUILD_DIR)/coverage.out ./...
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
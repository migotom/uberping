# Go
GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-s -w"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mybinary

# Project
UPING_PROJECT=cmd/uping/*

BIN_DIR=bin
BUILD_DIR=build

CONFIG_EXAMPLE=api-config.example.toml
HOSTS_EXAMPLE=hosts-file.example.txt

# Build
all: build
build: build-prepare build-linux-amd64 build-linux-386 build-darwin-amd64 clean

# Clean
clean:
	rm $(BUILD_DIR)/*

build-prepare:
	[ -d $(BUILD_DIR) ] || mkdir -p $(BUILD_DIR)
	[ -d $(BIN_DIR) ] || mkdir -p $(BIN_DIR)
	cp $(CONFIG_EXAMPLE) $(BUILD_DIR)/
	cp $(HOSTS_EXAMPLE) $(BUILD_DIR)/

# Cross compile
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/uping $(UPING_PROJECT)
	zip -j -9 $(BIN_DIR)/uping.linux.amd64.zip $(BUILD_DIR)/*
	rm $(BUILD_DIR)/uping

build-linux-386:
	GOOS=linux GOARCH=386 $(GOBUILD) -o $(BUILD_DIR)/uping $(UPING_PROJECT)
	zip -j -9 $(BIN_DIR)/uping.linux.386.zip $(BUILD_DIR)/*
	rm $(BUILD_DIR)/uping

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/uping $(UPING_PROJECT)
	zip -j -9 $(BIN_DIR)/uping.darwin.amd64.zip $(BUILD_DIR)/*
	rm $(BUILD_DIR)/uping
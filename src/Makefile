.PHONY: all build-RegisterBusinessFunction clean

SRC_DIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
PROJECT_ROOT := $(realpath $(SRC_DIR)/..)
BUILD_DIR := $(PROJECT_ROOT)/infrastructure/.aws-sam/build

all: build

build: build-RegisterBusinessFunction

build-RegisterBusinessFunction:
	@echo "Building RegisterBusinessFunction..."
	mkdir -p $(BUILD_DIR)/RegisterBusinessFunction
	cd $(SRC_DIR)/cmd/auth/register && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/RegisterBusinessFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/RegisterBusinessFunction/bootstrap

clean:
	rm -rf $(BUILD_DIR)
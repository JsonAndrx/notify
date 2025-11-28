.PHONY: all build-RegisterBusinessFunction build-RegenerateAPIKeyFunction build-AccountInfoFunction build-PlanUsageFunction build-SendNotificationFunction clean

SRC_DIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
PROJECT_ROOT := $(realpath $(SRC_DIR)/..)
BUILD_DIR := $(PROJECT_ROOT)/infrastructure/.aws-sam/build

all: build

build: build-RegisterBusinessFunction build-RegenerateAPIKeyFunction build-AccountInfoFunction build-PlanUsageFunction build-SendNotificationFunction

build-RegisterBusinessFunction:
	@echo "Building RegisterBusinessFunction..."
	mkdir -p $(BUILD_DIR)/RegisterBusinessFunction
	cd $(SRC_DIR)/cmd/auth/register && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/RegisterBusinessFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/RegisterBusinessFunction/bootstrap

build-RegenerateAPIKeyFunction:
	@echo "Building RegenerateAPIKeyFunction..."
	mkdir -p $(BUILD_DIR)/RegenerateAPIKeyFunction
	cd $(SRC_DIR)/cmd/account/regenerate-key && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/RegenerateAPIKeyFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/RegenerateAPIKeyFunction/bootstrap

build-AccountInfoFunction:
	@echo "Building AccountInfoFunction..."
	mkdir -p $(BUILD_DIR)/AccountInfoFunction
	cd $(SRC_DIR)/cmd/account/info && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/AccountInfoFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/AccountInfoFunction/bootstrap

build-PlanUsageFunction:
	@echo "Building PlanUsageFunction..."
	mkdir -p $(BUILD_DIR)/PlanUsageFunction
	cd $(SRC_DIR)/cmd/plan/usage && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/PlanUsageFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/PlanUsageFunction/bootstrap

build-SendNotificationFunction:
	@echo "Building SendNotificationFunction..."
	mkdir -p $(BUILD_DIR)/SendNotificationFunction
	cd $(SRC_DIR)/cmd/notifications/send && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/SendNotificationFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/SendNotificationFunction/bootstrap

clean:
	rm -rf $(BUILD_DIR)
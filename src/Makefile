.PHONY: all build-RegisterBusinessFunction build-RegenerateAPIKeyFunction build-AccountInfoFunction build-PlanUsageFunction build-SendWhatsAppFunction build-SendSMSFunction build-SendEmailFunction clean

SRC_DIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
PROJECT_ROOT := $(realpath $(SRC_DIR)/..)
BUILD_DIR := $(PROJECT_ROOT)/infrastructure/aws-sam/build

all: build

build: build-RegisterBusinessFunction build-RegenerateAPIKeyFunction build-AccountInfoFunction build-PlanUsageFunction build-SendWhatsAppFunction build-SendSMSFunction build-SendEmailFunction

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

build-SendWhatsAppFunction:
	@echo "Building SendWhatsAppFunction..."
	mkdir -p $(BUILD_DIR)/SendWhatsAppFunction
	cd $(SRC_DIR)/cmd/notifications/whatsapp && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/SendWhatsAppFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/SendWhatsAppFunction/bootstrap

build-SendSMSFunction:
	@echo "Building SendSMSFunction..."
	mkdir -p $(BUILD_DIR)/SendSMSFunction
	cd $(SRC_DIR)/cmd/notifications/sms && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/SendSMSFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/SendSMSFunction/bootstrap

build-SendEmailFunction:
	@echo "Building SendEmailFunction..."
	mkdir -p $(BUILD_DIR)/SendEmailFunction
	cd $(SRC_DIR)/cmd/notifications/email && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/SendEmailFunction/bootstrap main.go
	chmod +x $(BUILD_DIR)/SendEmailFunction/bootstrap

clean:
	rm -rf $(BUILD_DIR)
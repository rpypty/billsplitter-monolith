ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
BIN_NAME := billsplitter
BIN_PATH := $(ROOT_DIR)/$(BIN_NAME)

vps-back-deploy:
	cd $(ROOT_DIR) && go build -o $(BIN_PATH) ./cmd/main.go
	nohup $(BIN_PATH) > /var/log/billsplitter.log 2>&1 &

vps-back-stop:
	pkill -9 -f "$(BIN_PATH)"
	

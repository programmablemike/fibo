#!/usr/bin/env make
SHELL := bash
CC := go build
CC_ARGS := --mod=vendor
CURRENT_DIR := $(shell pwd)
BIN_DIR := $(CURRENT_DIR)/bin

GOOS ?= darwin
GOARCH ?= arm64

list-build-targets:
	go tool dist list

$(BIN_DIR):
	mkdir -p $@

build: $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(CC) $(CC_ARGS) -o $(BIN_DIR)/fibo_$(GOOS)_$(GOARCH)

update-deps:
	go get -u
	go mod vendor
	git add go.mod go.sum vendor/*

clean:
	rm -rf $(BIN_DIR)
#!/usr/bin/env make
SHELL := bash
CC := go build
CC_ARGS := --mod=vendor
CURRENT_DIR := $(shell pwd)

build:
	mkdir -p $(CURRENT_DIR)/bin
	GOOS="darwin" GOARCH="arm64" $(CC) $(CC_ARGS) -o bin/fibo_darwin_arm64
	GOOS="darwin" GOARCH="amd64" $(CC) $(CC_ARGS) -o bin/fibo_darwin_amd64
	GOOS="linux"  GOARCH="arm64" $(CC) $(CC_ARGS) -o bin/fibo_linux_arm64
	GOOS="linux"  GOARCH="386"   $(CC) $(CC_ARGS) -o bin/fibo_linux_386
	GOOS="linux"  GOARCH="amd64" $(CC) $(CC_ARGS) -o bin/fibo_linux_amd64

clean:
	rm -rf $(CURRENT_DIR)/bin
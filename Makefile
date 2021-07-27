#!/usr/bin/env make
.PHONY: build docker list-build-targets update-deps k6
SHELL := bash
CC_CMD := go build
CC_OPTS := --mod=vendor
TEST_CMD := go test
TEST_OPTS := -v

GOOS ?= darwin
GOARCH ?= arm64

list-build-targets:
	go tool dist list

fibo_$(GOOS)_$(GOARCH):
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(CC_CMD) $(CC_OPTS) -o $@

build: fibo_$(GOOS)_$(GOARCH)

test:
	$(TEST_CMD) $(TEST_OPTS) ./...

docker:
	docker build -f ./docker/Dockerfile -t fibo:dev .

update-deps:
	go get -u
	go mod vendor
	git add go.mod go.sum vendor/*

k6:
	k6 run ./k6/worker.js

clean:
	rm fibo_*
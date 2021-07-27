#!/usr/bin/env make
.PHONY: list-build-targets build test docker docker-build docker-logs-app docker-logs-db docker-start docker-stop update-deps k6 clean
SHELL := bash
CC_CMD := go build
CC_OPTS := --mod=vendor
TEST_CMD := go test
TEST_OPTS := -v
GO_SRCS := $(shell find . -type f -name "*.go")

GOOS ?= darwin
GOARCH ?= arm64

list-build-targets:
	go tool dist list

fibo_$(GOOS)_$(GOARCH): $(GO_SRCS)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(CC_CMD) $(CC_OPTS) -o $@

build: fibo_$(GOOS)_$(GOARCH)

test:
	$(TEST_CMD) $(TEST_OPTS) ./... --bench=.

docker-build:
	docker-compose build

docker-logs-app:
	docker-compose logs fibo

docker-logs-db:
	docker-compose logs postgres

docker-start:
	docker-compose up --detach

docker-stop:
	docker-compose down -v

update-deps:
	go get -u
	go mod vendor
	git add go.mod go.sum vendor/*

k6:
	k6 run ./k6/calculate-worker.js

clean:
	rm fibo_*
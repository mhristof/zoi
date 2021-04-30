MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
ifeq ($(word 1,$(subst ., ,$(MAKE_VERSION))),4)
.SHELLFLAGS := -eu -o pipefail -c
endif
.DEFAULT_GOAL := all
.ONESHELL:

UNAME := $(shell uname | tr A-Z a-z)
GO_SRC := $(shell find ./ -name '*.go')

.PHONY: all
all: ./bin/zoi.darwin ./bin/zoi.linux

./bin/zoi.darwin: $(GO_SRC)
	GOOS=darwin go build -o $@ main.go

./bin/zoi.linux: $(GO_SRC)
	GOOS=linux go build -o $@ main.go

.PHONY: fast-test
fast-test:  ## Run fast tests
	go test ./... -tags fast

.PHONY: test
test:	## Run all tests
	go test -v ./...

.PHONY: simple
simple: ./bin/zoi.$(UNAME)
	./bin/zoi.$(UNAME) ./tests/simple.py | python | grep '^https'

.PHONY: help
help:           ## Show this help.
	@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ## /:/g' | column -t -s:

.PHONY: clean
clean:
	rm -rf bin/*

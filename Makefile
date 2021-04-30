MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
ifeq ($(word 1,$(subst ., ,$(MAKE_VERSION))),4)
.SHELLFLAGS := -eu -o pipefail -c
endif
.DEFAULT_GOAL := all
.ONESHELL:

GIT_REF := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git name-rev --tags --name-only $(GIT_REF))
UNAME := $(shell uname | tr A-Z a-z)
GO_SRC := $(shell find ./ -name '*.go')

.PHONY: all
all: ./bin/zoi.darwin ./bin/zoi.linux

./bin/zoi.%: $(GO_SRC)
	GOOS=$* go build -o $@ -ldflags "-X github.com/mhristof/zoi/cmd.version=$(GIT_TAG)+$(GIT_REF)" main.go

.PHONY: fast-test
fast-test:  ## Run fast tests
	go test ./... -tags fast

.PHONY: test
test:	## Run all tests
	go test -v ./...

.PHONY: simple
simple: ./bin/zoi.$(UNAME)
	./bin/zoi.$(UNAME) ./tests/simple.py | python | grep '^https'

.PHONY: install
install: ./bin/zoi.$(UNAME)
	cp ./bin/zoi.$(UNAME) ~/bin/zoi

.PHONY: help
help:           ## Show this help.
	@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.* ## /:/g' | column -t -s:

.PHONY: clean
clean:
	rm -rf bin/*

SHELL := /bin/bash

-include Makefile.env

.PHONY: all check format vet lint build test generate tidy

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check               to do static check"
	@echo "  build               to create bin directory and build"

check: vet

format:
	go fmt ./...

vet:
	go vet ./...

build: tidy format check
	go build ./...

tidy:
	go mod tidy
	go mod verify

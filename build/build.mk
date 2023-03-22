SHELL := /usr/bin/bash

SRC=./tarantella-server
BIN=$(notdir $(SRC))
BIN_PATH=./.bin/$(EXEC)

$(BIN_PATH):
	@echo "Building $@..."
	@mkdir -p $(@D)
	GOBIN=$(abspath $(@D)) go install $(SRC)

.PHONY: build
build: $(BIN_PATH) # Build application

.PHONY: clean
clean: # Clean all artefacts
	rm -rf .bin
	rm -rf .build
	go clean

.PHONY: dependencies
dependencies: # Install all required dependencies
	go get .../.

.PHONY: run
run: # Run application in development mode
	test -f .env || touch .env
	$(shell cat .env) go run $(SRC)

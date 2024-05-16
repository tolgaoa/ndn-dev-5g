# Makefile

SRC_DIR := ./cmd

# Set the name of the Go module
MODULE_NAME := http1server

build:
	go build -o ./bin/http1server $(SRC_DIR)/main.go

# Run "make image VERSION=x.x.x"
image:
	docker build . -t tolgaomeratalay/http1server:$(VERSION)
push:
	docker push tolgaomeratalay/http1server:$(VERSION)

run:
	go run $(SRC_DIR)/main.go

deps:
	go mod download

clean:
	rm -f ./bin/http1server

.PHONY: build run deps clean

# Makefile

SRC_DIR := ./cmd

# Set the name of the Go module
MODULE_NAME := http1client

build:
	go build -o ./bin/http1client $(SRC_DIR)/main.go

# Run "make image VERSION=x.x.x"
image:
	docker build . -t tolgaomeratalay/http1client:$(VERSION)
push:
	docker push tolgaomeratalay/http1client:$(VERSION)

run:
	go run $(SRC_DIR)/main.go

deps:
	go mod download

clean:
	rm -f ./bin/http1client

.PHONY: build run deps clean

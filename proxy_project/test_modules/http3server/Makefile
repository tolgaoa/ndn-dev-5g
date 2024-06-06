# Makefile

SRC_DIR := ./cmd

# Set the name of the Go module
MODULE_NAME := http3server

build:
	go build -o ./bin/http3server $(SRC_DIR)/main.go

# Run "make image VERSION=x.x.x"
image:
	docker build . -t tolgaomeratalay/http3server:$(VERSION)
push:
	docker push tolgaomeratalay/http3server:$(VERSION)

run:
	go run $(SRC_DIR)/main.go

deps:
	go mod download

clean:
	rm -f ./bin/http3server

.PHONY: build run deps clean

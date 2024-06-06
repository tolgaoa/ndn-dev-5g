# Makefile

SRC_DIR := ./cmd

# Set the name of the Go module
MODULE_NAME := proxy

build:
	go build -o ./bin/proxy $(SRC_DIR)/main.go

# Run "make image VERSION=x.x.x"
image:
	docker build . -t tolgaomeratalay/biproxy:$(VERSION)
push:
	docker push tolgaomeratalay/biproxy:$(VERSION)

run:
	go run $(SRC_DIR)/main.go

deps:
	go mod download

clean:
	rm -f ./bin/proxy

.PHONY: build run deps clean

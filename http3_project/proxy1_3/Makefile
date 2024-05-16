# Makefile

SRC_DIR := ./cmd

# Set the name of the Go module
MODULE_NAME := proxy1_3

build:
	go build -o ./bin/proxy1_3 $(SRC_DIR)/main.go

# Run "make image VERSION=x.x.x"
image:
	docker build . -t tolgaomeratalay/httpproxy1_3:$(VERSION)
push:
	docker push tolgaomeratalay/httpproxy1_3:$(VERSION)

run:
	go run $(SRC_DIR)/main.go

deps:
	go mod download

clean:
	rm -f ./bin/proxy1_3

.PHONY: build run deps clean

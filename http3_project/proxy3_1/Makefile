# Makefile

SRC_DIR := ./cmd

# Set the name of the Go module
MODULE_NAME := proxy3_1

build:
	go build -o ./bin/proxy3_1 $(SRC_DIR)/main.go

# Run "make image VERSION=x.x.x"
image:
	docker build . -t tolgaomeratalay/proxy3_1:$(VERSION)
push:
	docker push tolgaomeratalay/proxy3_1:$(VERSION)

run:
	go run $(SRC_DIR)/main.go

deps:
	go mod download

clean:
	rm -f ./bin/proxy3_1

.PHONY: build run deps clean

# Makefile for storage package

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=storage

# Build targets
.PHONY: all build clean test coverage lint tidy

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

test:
	$(GOTEST) -v ./...

test-short:
	$(GOTEST) -v -short ./...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

lint:
	golangci-lint run

tidy:
	$(GOMOD) tidy

deps:
	$(GOGET) -u ./...

# Docker targets
.PHONY: docker-build docker-run

docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run --rm -it $(BINARY_NAME)

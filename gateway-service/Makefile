.PHONY: all deps build test run clean

# Project Settings
PROJECT_NAME := gateway-service

# Dependencies
deps:
	go mod tidy


# Build
build:
	go build -o $(PROJECT_NAME) ./cmd

# Testing (add your test command)
test: 
	go test ./...

# Run the Service
run: build
	./$(PROJECT_NAME)

# Cleanup
clean:
	rm -f $(PROJECT_NAME) 
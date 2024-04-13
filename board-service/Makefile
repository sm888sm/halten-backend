.PHONY: all deps proto build test run clean

# Project Settings
PROJECT_NAME := board-service

# Dependencies
deps:
	go mod tidy

# Protobuf Compilation
proto:
	protoc --proto_path=./api/proto --go-grpc_out=./api/pb --go-grpc_opt=paths=source_relative \
       --go_out=./api/pb --go_opt paths=source_relative \
       ./api/proto/*.proto

# Build
build: proto deps
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
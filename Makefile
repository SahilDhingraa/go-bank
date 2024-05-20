# Description: Makefile for gobank project

build:
	@echo "Building the project..."
	@go build -o bin/gobank

run: build
	@echo "Running the project..."
	@./bin/gobank

test:
	@echo "Running tests..."
	@go test -v ./...
	
clean:
	@echo "Cleaning the project..."
	@rm -rf bin
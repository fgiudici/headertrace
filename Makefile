.PHONY: help generate build test clean

help:
	@echo "Available targets:"
	@echo "  make generate - Generate code from OpenAPI spec using oapi-codegen"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove generated files and build artifacts"

generate:
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen -config api/config.yaml api/openapi.yaml

build:
	@echo "Building application..."
	go build -o bin/headertrace .

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning up..."
	rm -f bin/headertrace

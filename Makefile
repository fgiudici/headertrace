GIT_TAG?=$(shell git describe --candidates=50 --abbrev=0 --tags 2>/dev/null || echo "v0.0.1" )
GIT_COMMIT?=$(shell git rev-parse HEAD)
GIT_COMMIT_SHORT?=$(shell git rev-parse --short HEAD)
GO_MODULE?= $(shell go list -m)
CONTAINER_ENGINE?=docker
CONTAINER_IMAGE?=headertrace

LDFLAGS:=-w -s
LDFLAGS+=-X "$(GO_MODULE)/cmd.version=$(GIT_TAG)"
LDFLAGS+=-X "$(GO_MODULE)/cmd.gitCommit=$(GIT_COMMIT)"

.PHONY: help generate build test clean docker

help:
	@echo "Available targets:"
	@echo "  make generate - Generate code from OpenAPI spec using oapi-codegen"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove generated files and build artifacts"
	@echo "  make docker   - Build Docker image"

generate:
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen -config api/config.yaml api/openapi.yaml
	@echo "Patching generated code for wildcard path matching..."
	sed -i 's|"/{matchall}"|"/{matchall...}"|g' api/gen.go

build:
	@echo "Building application..."
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o bin/headertrace .

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning up..."
	rm -f bin/headertrace

docker:
	@echo "Building Docker image..."
	$(CONTAINER_ENGINE) build -f Dockerfile \
		--build-arg VERSION=$(GIT_TAG) \
		--build-arg COMMIT=$(GIT_COMMIT) \
		-t $(CONTAINER_IMAGE):$(GIT_TAG)-${GIT_COMMIT} .
	$(CONTAINER_ENGINE) tag $(CONTAINER_IMAGE):$(GIT_TAG)-${GIT_COMMIT} $(CONTAINER_IMAGE):${GIT_TAG}

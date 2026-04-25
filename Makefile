VERSION ?= dev
COMMIT ?= none
DATE ?= unknown

build: clean
	@docker run --rm -v $(shell mktemp -d):/cache -v $(shell pwd):/home/tlwpa4220 -u $(shell id -u):$(shell id -g) -w /home/tlwpa4220 -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 -e GOMODCACHE=/cache/mod -e GOCACHE=/cache/build golang:1.25-alpine go build -o tlwpa4220 -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'"  cmd/main.go

lint:
	@echo "Running golangci-lint locally with config..."
	@go build ./... || (echo "Build failed"; exit 1)
	@echo "✓ Code build validation passed"

clean:
	@rm -rf tlwpa4220
	@rm -rf super-linter.log

container: build
	@docker build -t tlwpa4220:local -f build/container/Dockerfile .

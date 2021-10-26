APP_VERSION ?= $(shell git describe --abbrev=5 --dirty --tags --always)

BINDIR := $(PWD)/bin
OUTPUT_DIR := $(PWD)/_output

GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64

PATH := $(BINDIR):$(PATH)
SHELL := env PATH='$(PATH)' /bin/sh

all: build

# Run tests
test: fmt vet
	@# Disable --race until https://github.com/kubernetes-sigs/controller-runtime/issues/1171 is fixed.
	ginkgo --randomizeAllSpecs --randomizeSuites --failOnPending --flakeAttempts=2 \
			--cover --coverprofile cover.out --trace --progress  $(TEST_ARGS)\
			./pkg/... ./cmd/...

# Build sail binary
build: fmt vet
	go build -o $(OUTPUT_DIR)/sail ./cmd/sail

# Cross compiler
build-all: fmt vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(OUTPUT_DIR)/sail_$(APP_VERSION)_linux_amd64 ./cmd/sail
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o $(OUTPUT_DIR)/sail_$(APP_VERSION)_linux_arm64 ./cmd/sail
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o $(OUTPUT_DIR)/sail_$(APP_VERSION)_darwin_amd64 ./cmd/sail
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -o $(OUTPUT_DIR)/sail_$(APP_VERSION)_darwin_arm64 ./cmd/sail

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

lint:
	$(BINDIR)/golangci-lint run --timeout 2m0s ./pkg/... ./cmd/...

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go install github.com/onsi/ginkgo/ginkgo@v1.16.4

	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(BINDIR) latest
SHELL = /bin/bash
GO-VER = go1.24

default: build

# #### GO Binary Management ####
deps-go-binary:
	echo "Expect: $(GO-VER)" && \
		echo "Actual: $$(go version)" && \
	 	go version | grep $(GO-VER) > /dev/null

# #### CLEAN ####

clean: deps-go-binary
	rm -rf build/*
	go clean --modcache

# #### DEPS ####

deps-modules: deps-go-binary
	go mod download

deps-counterfeiter: deps-modules
	go install github.com/maxbrunsfeld/counterfeiter/v6@latest

deps-ginkgo: deps-go-binary
	go install github.com/onsi/ginkgo/v2/ginkgo@latest

deps-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

deps: deps-modules deps-counterfeiter deps-ginkgo deps-golangci-lint

# #### BUILD ####

SRC = $(shell find . -name "*.go" | grep -v "_test\." )
VERSION := $(or $(VERSION), dev)
LDFLAGS="-X github.com/cf-platform-eng/tileinspect/version.Version=$(VERSION)"

build/tileinspect: $(SRC) deps
	go build -o build/tileinspect -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go

build: build/tileinspect

build-all: build-linux build-darwin build-windows

build-linux: build/tileinspect-linux-amd64 build/tileinspect-linux-arm64
build/tileinspect-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o build/tileinspect-linux-amd64 -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go
build/tileinspect-linux-arm64:
	GOARCH=arm64 GOOS=linux go build -o build/tileinspect-linux-arm64 -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go


build-darwin: build/tileinspect-darwin-amd64 build/tileinspect-darwin-arm64
build/tileinspect-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o build/tileinspect-darwin-amd64 -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go
build/tileinspect-darwin-arm64:
	GOARCH=arm64 GOOS=darwin go build -o build/tileinspect-darwin-arm64 -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go

build-windows: build/tileinspect-windows-amd64.exe build/tileinspect-windows-arm64.exe
build/tileinspect-windows-amd64.exe:
	GOARCH=amd64 GOOS=windows go build -o build/tileinspect-windows-amd64.exe -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go
build/tileinspect-windows-arm64.exe:
	GOARCH=arm64 GOOS=windows go build -o build/tileinspect-windows-arm64.exe -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go

build-image: build/tileinspect-linux-amd64
	docker build --tag cfplatformeng/tileinspect:${VERSION} --file Dockerfile .

# #### TESTS ####
test: deps lint
	ginkgo -skip-package features -r .

test-features: deps lint
	ginkgo -tags feature -r features

lint: deps-golangci-lint
	golangci-lint run

.PHONY: set-pipeline
set-pipeline: ci/pipeline.yaml
	fly -t ppe-isv set-pipeline -p tileinspect -c ci/pipeline.yaml

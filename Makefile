SHELL = /bin/bash
GO-VER = go1.17

default: build

# #### GO Binary Management ####
deps-go-binary:
	echo "Expect: $(GO-VER)" && \
		echo "Actual: $$(go version)" && \
	 	go version | grep $(GO-VER) > /dev/null


HAS_GO_IMPORTS := $(shell command -v goimports;)

deps-goimports: deps-go-binary
ifndef HAS_GO_IMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif

# #### CLEAN ####
clean: deps-go-binary
	rm -rf build/*
	go clean --modcache


# #### DEPS ####

deps-modules: deps-goimports deps-go-binary
	go mod download

deps-counterfeiter: deps-modules
	command -v counterfeiter >/dev/null 2>&1 || go get -u github.com/maxbrunsfeld/counterfeiter/v6

deps-ginkgo: deps-go-binary
	command -v ginkgo >/dev/null 2>&1 || go get -u github.com/onsi/ginkgo/ginkgo github.com/onsi/gomega

deps: deps-modules deps-counterfeiter deps-ginkgo


# #### BUILD ####
SRC = $(shell find . -name "*.go" | grep -v "_test\." )

VERSION := $(or $(VERSION), "dev")

LDFLAGS="-X github.com/cf-platform-eng/tileinspect/version.Version=$(VERSION)"

build/tileinspect: $(SRC) deps
	go build -o build/tileinspect -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go

build: build/tileinspect

build-all: build-linux build-darwin

build-linux: build/tileinspect-linux

build/tileinspect-linux:
	GOARCH=amd64 GOOS=linux go build -o build/tileinspect-linux -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go

build-darwin: build/tileinspect-darwin

build/tileinspect-darwin:
	GOARCH=amd64 GOOS=darwin go build -o build/tileinspect-darwin -ldflags ${LDFLAGS} ./cmd/tileinspect/main.go

build-image: build/tileinspect-linux
	docker build --tag cfplatformeng/tileinspect:${VERSION} --file Dockerfile .

test: deps lint
	ginkgo -skipPackage features -r .

test-features: deps lint
	ginkgo -tags feature -r features

lint: deps-goimports
	git ls-files | grep '.go$$' | xargs goimports -l -w

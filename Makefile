SHELL = /bin/bash
GO-VER = go1.12

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

deps: deps-goimports deps-go-binary
	go mod download

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
	docker build --tag cfplatformeng/needs:${VERSION} --file Dockerfile .

test: deps lint
	ginkgo -r .

lint: deps-goimports
	git ls-files | grep '.go$$' | xargs goimports -l -w

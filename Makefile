PROJECT_NAME := "dougkirkley/trex"
PKG := "github.com/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

GO=go

.PHONY: all dep build clean test lint

all: build

format: ## format the files
	@gofmt -w ${GO_FILES}
	
lint: ## Lint the files 
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests  
	@go test -short ${PKG_LIST}

clean: ## Remove previous build 
	rm -Rf $(DIST_ROOT)
	$(GO) clean -i ./...

build-osx:
	@echo Build OSX amd64
	env GOOS=darwin GOARCH=amd64 $(GO) build -o trex-darwin-amd64

build-linux:
	@echo Build linux amd64
	env GOOS=linux GOARCH=amd64 $(GO) build -o trex-linux-amd64

build: build-osx build-linux

help: ## Display this help screen 
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' 

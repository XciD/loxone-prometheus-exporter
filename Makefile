BUILD_DIR 		:= build

.DEFAULT_GOAL := build

.PHONY: init
init:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/modocache/gover
	go mod download

.PHONY: lint
lint:
	golangci-lint run --config golangci.yml

.PHONY: build
build:
	env GOOS=linux go build -o $(BUILD_DIR)/exporter .

.PHONY: format
format:
	go fmt $(go list)
	goimports -e -w -d $(shell find . -type f -name '*.go' -print)
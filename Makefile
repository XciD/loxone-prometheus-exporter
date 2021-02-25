BUILD_DIR 		    := build
.DEFAULT_GOAL       := build
BUILD_DIR 		    := build
NAME 				:= loxone-prometheus-exporter
REGISTRY 			:= xcid
TAG					:= latest

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

.PHONY: fmt
fmt:
	goimports -e -w -d $(shell find . -type f -name '*.go' -print)
	gofmt -e -w -d $(shell find . -type f -name '*.go' -print)

.PHONY: cover
cover:
	$(GOPATH)/bin/gover . coverage.txt

.PHONY: docker-build
docker-build:
	docker build -t $(NAME) .

.PHONY: docker-push
docker-push:
	docker tag $(NAME) $(REGISTRY)/$(NAME):$(TAG)
	docker push $(REGISTRY)/$(NAME):$(TAG)

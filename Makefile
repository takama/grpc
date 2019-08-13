PROJECT = github.com/takama/grpc
APP = grpc
BIN = grpc
SERVICE_NAME ?= $(shell echo "$(APP)" | tr - _)

# Use the v0.0.0 tag for testing, it shouldn't clobber any release builds
RELEASE ?= v0.0.0
GOOS ?= linux
GOARCH ?= amd64
CA_DIR ?= certs

# Configs for GKE
GKE_PROJECT_ID ?= drs-017
GKE_PROJECT_ZONE ?= europe-west4-b
GKE_CLUSTER_NAME ?= domingo-01

KUBE_CONTEXT ?= gke_$(GKE_PROJECT_ID)_$(GKE_PROJECT_ZONE)_$(GKE_CLUSTER_NAME)

REGISTRY ?= gcr.io/$(GKE_PROJECT_ID)

# Common configuration
GRPC_SERVER_PORT ?= 8000
GRPC_EXTERNAL_PORT ?= 8000
GRPC_INFO_PORT ?= 8080
GRPC_INFO_EXTERNAL_PORT ?= 8080
GRPC_LOGGER_LEVEL ?= 0
GRPC_CONFIG_PATH ?= /etc/$(SERVICE_NAME)/default.conf

# Namespace: dev, prod, username ...
NAMESPACE ?= test
VALUES ?= values-$(NAMESPACE)

CONTAINER_IMAGE ?= $(REGISTRY)/$(APP)
CONTAINER_NAME ?= $(APP)

REPO_INFO=$(shell git config --get remote.origin.url)
REPO_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
RELEASE_DATE=$(shell date +%FT%T%Z)

ifndef REPO_COMMIT
REPO_COMMIT = git-$(shell git rev-parse --short HEAD)
endif

BUILD =? $(RELEASE)
DEPLOY_PARAMS ?= --wait
ifneq ("$(NAMESPACE)","prod")
BUILD = $(RELEASE)-$(REPO_COMMIT)-$(NAMESPACE)
DEPLOY_PARAMS = --wait --recreate-pods --force
endif

LDFLAGS = "-s -w \
	-X $(PROJECT)/pkg/version.RELEASE=$(RELEASE) \
	-X $(PROJECT)/pkg/version.DATE=$(RELEASE_DATE) \
	-X $(PROJECT)/pkg/version.REPO=$(REPO_INFO) \
	-X $(PROJECT)/pkg/version.COMMIT=$(REPO_COMMIT) \
	-X $(PROJECT)/pkg/version.BRANCH=$(REPO_BRANCH)"

GO_PACKAGES=$(shell go list $(PROJECT)/pkg/...)

BUILDTAGS=

all: build

check-all: fmt imports test lint

project:
	@echo "+ $@"
ifneq ("$(GKE_PROJECT_ID)", "$(shell gcloud config get-value project)")
	@gcloud config set project $(GKE_PROJECT_ID)
endif
ifneq ("$(GKE_PROJECT_ZONE)", "$(shell gcloud config get-value compute/zone)")
	@gcloud config set compute/zone $(GKE_PROJECT_ZONE)
endif
ifneq ("$(GKE_CLUSTER_NAME)", "$(shell gcloud config get-value container/cluster)")
	@gcloud config set container/cluster $(GKE_CLUSTER_NAME)
endif

cluster:
	@echo "+ $@"
ifneq ("$(KUBE_CONTEXT)", "$(shell kubectl config get-clusters | grep $(KUBE_CONTEXT))")
	@gcloud container clusters get-credentials $(GKE_CLUSTER_NAME) --zone $(GKE_PROJECT_ZONE) --project $(GKE_PROJECT_ID)
endif
ifneq ("$(KUBE_CONTEXT)", "$(shell kubectl config current-context)")
	@kubectl config use-context $(KUBE_CONTEXT)
endif

vendor: bootstrap
	@echo "+ $@"
	@go mod tidy

contracts:
	@echo "+ $@"
	@$(MAKE) -C contracts generate

compile: contracts vendor test lint
	@echo "+ $@"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo \
		-ldflags $(LDFLAGS) -o bin/$(GOOS)-$(GOARCH)/$(BIN) $(PROJECT)/cmd

certs:
ifeq ("$(wildcard $(CA_DIR)/ca-certificates.crt)","")
	@echo "+ $@"
	@docker run --name $(CONTAINER_NAME)-certs -d alpine:latest \
	sh -c "apk --update upgrade && apk add ca-certificates && update-ca-certificates"
	@docker wait $(CONTAINER_NAME)-certs
	@mkdir -p $(CA_DIR)
	@docker cp $(CONTAINER_NAME)-certs:/etc/ssl/certs/ca-certificates.crt $(CA_DIR)
	@docker rm -f $(CONTAINER_NAME)-certs
endif

build: compile certs
	@echo "+ $@"
	@docker build --pull -t $(CONTAINER_IMAGE):$(BUILD) .

push: build project
	@echo "+ $@"
	@docker push $(CONTAINER_IMAGE):$(BUILD)

run: clean build
	@echo "+ $@"
	@docker run --name $(CONTAINER_NAME) \
		-p $(GRPC_EXTERNAL_PORT):$(GRPC_SERVER_PORT) \
		-p $(GRPC_INFO_EXTERNAL_PORT):$(GRPC_INFO_PORT) \
		-e "GRPC_SERVER_PORT=$(GRPC_SERVER_PORT)" \
		-e "GRPC_INFO_PORT=$(GRPC_INFO_PORT)" \
		-e "GRPC_LOGGER_LEVEL=$(GRPC_LOGGER_LEVEL)" \
		-e "GRPC_CONFIG_PATH=$(GRPC_CONFIG_PATH)" \
		-v $(shell pwd)/config/default.conf:$(GRPC_CONFIG_PATH):ro \
		-d $(CONTAINER_IMAGE):$(BUILD)
	@sleep 1
	@docker logs $(CONTAINER_NAME)

logs:
	@echo "+ $@"
	@docker logs -f $(CONTAINER_NAME)

deploy: push cluster
	@echo "+ $@"
	@helm upgrade $(CONTAINER_NAME)-$(NAMESPACE) -f .helm/$(VALUES).yaml .helm --kube-context $(KUBE_CONTEXT) \
		--namespace $(NAMESPACE) --version=$(RELEASE) --set image.tag=$(BUILD) -i $(DEPLOY_PARAMS)

charts:
	@echo "+ $@"
	@helm template .helm -n $(APP)-$(NAMESPACE) --namespace $(NAMESPACE) -f .helm/$(VALUES).yaml

test:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"go test -race -cover {{.Dir}}"{{end}}' $(GO_PACKAGES) | xargs -L 1 sh -c

cover:
	@echo "+ $@"
	@echo "mode: set" > coverage.txt
	@go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}} && \
		cat {{.Dir}}/.coverprofile | sed 1d >> coverage.txt"{{end}}' $(GO_PACKAGES) | xargs -L 1 sh -c

fmt:
	@echo "+ $@"
	@go list -f '"gofmt -w -s -l {{.Dir}}"' $(GO_PACKAGES) | xargs -L 1 sh -c

imports:
	@echo "+ $@"
	@go list -f '"goimports -w {{.Dir}}"' ${GO_PACKAGES} | xargs -L 1 sh -c

lint: bootstrap
	@echo "+ $@"
	@golangci-lint run --enable-all --skip-dirs vendor ./...

HAS_RUNNED := $(shell docker ps | grep $(CONTAINER_NAME))
HAS_EXITED := $(shell docker ps -a | grep $(CONTAINER_NAME))

stop:
ifdef HAS_RUNNED
	@echo "+ $@"
	@docker stop $(CONTAINER_NAME)
endif

start: stop
	@echo "+ $@"
	@docker start $(CONTAINER_NAME)

rm:
ifdef HAS_EXITED
	@echo "+ $@"
	@docker rm $(CONTAINER_NAME)
endif

version:
	@./bumper.sh

clean: stop rm
	@rm -f bin/$(GOOS)-$(GOARCH)/$(BIN)

HAS_LINT := $(shell command -v golangci-lint;)
HAS_IMPORTS := $(shell command -v goimports;)
HAS_GCLOUD := $(shell command -v gcloud;)
HAS_DOCKER_GCR := $(shell command -v docker-credential-gcr)

bootstrap:
ifndef HAS_LINT
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif
ifndef HAS_IMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif
ifndef HAS_GCLOUD
	@echo "gcloud cli utility should be installed"
	@echo "Pre-compiled binaries for your platform:
	@echo https://console.cloud.google.com/storage/browser/cloud-sdk-release?authuser=0"
	@exit 1
endif
ifdef HAS_GCLOUD
ifndef HAS_DOCKER_GCR
	@gcloud components install docker-credential-gcr -q
	@docker-credential-gcr configure-docker
endif
endif

.PHONY: all \
	project \
	cluster \
	vendor \
	contracts \
	compile \
	build \
	certs \
	push \
	run \
	logs \
	deploy \
	charts \
	test \
	cover \
	fmt \
	lint \
	stop \
	start \
	rm \
	version \
	clean \
	bootstrap

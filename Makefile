.PHONY: build test clean 


NAME   := device-ethernetip-go
ORG    := $(org)
REPOSITORY := $(repo)
PLATFORM:=

ifeq ($(platform),amd64)
PLATFORM:=
endif

ifeq ($(platform),arm64)
PLATFORM:=-arm64
endif

REPO   := ${ORG}/${NAME}${PLATFORM}
TAG    := $(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
IMG    := ${REPO}:${TAG}
LATEST := ${REPO}:latest
DOCKER_EMAIL := kemal.durkaya@konzek.com


GO=go
CGO=CGO_ENABLED=1 GO111MODULE=on $(GO)

MICROSERVICES=cmd/${NAME}

.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)

GOFLAGS=-ldflags "-X ${NAME}.Version=$(VERSION)"

GIT_SHA=$(shell git rev-parse HEAD)

tidy:
	go mod tidy


build: $(MICROSERVICES)
	$(CGO) install -tags=safe


cmd/device-ethernetip-go:
	$(CGO) build $(GOFLAGS) -o $@ ./cmd


docker-login: ## login to DockerHub with credentials found in env
	docker login ${REPOSITORY} -u ${DOCKER_EMAIL}

docker: 
	docker build \
	    --build-arg http_proxy \
	    --build-arg https_proxy \
		-f Dockerfile \
		-t ${REPO}:${TAG} \
		.

docker-buildx:
	docker buildx build --platform=linux/$(platform) -f Dockerfile -t ${IMG} . -o type=docker
	

image-tag:

	docker image tag ${IMG} ${REPOSITORY}/${IMG}

repo-push:

	docker image push  ${REPOSITORY}/${IMG}


test:
	$(CGO) test ./... -coverprofile=coverage.out
	$(CGO) vet ./...
	gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")`" = "" ]
	./bin/test-attribution-txt.sh
	
clean:
	rm -f $(MICROSERVICES)


vendor:
	$(GO) mod vendor

BINARY_NAME=echoserver

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_HASH   := $(shell git rev-parse --verify HEAD)
GIT_TAG    := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null || echo "")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

ifdef GIT_TAG
	VERSION := $(GIT_TAG)-$(GIT_HASH)
else
	VERSION := $(GIT_BRANCH)-$(GIT_HASH)
endif

buildmacos:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go

buildlinux:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go

buildwindows:
	GOARCH=amd64 GOOS=window go build -o ${BINARY_NAME}-windows main.go

run:
	./${BINARY_NAME}

build_and_run:
	buildwindows run

azure:
	echo "Building ${VERSION}"
	az acr build -t common/echoserver:${VERSION} -r ${Registry} .
clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all
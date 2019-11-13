# Go parameters
GO_CMD=go
GO_BUILD=${GO_CMD} build
GO_CLEAN=${GO_CMD} clean
GO_TEST=${GO_CMD} test

# Environments
ARCH=amd64
ENV_LINUX=linux
ENV_OSX=darwin
ENV_WINDOWS=windows

# Binary names
BINARY_NAME=feedback-service
BINARY_LINUX=${BINARY_NAME}-linux
BINARY_OSX=${BINARY_NAME}-max
BINARY_WINDOWS=${BINARY_NAME}.exe

all: clean test build
build:
	env GOOS=${ENV_LINUX} GOARCH=${ARCH} ${GO_BUILD} -o ${BINARY_LINUX}
	env GOOS=${BINARY_OSX} GOARCH=${ARCH} ${GO_BUILD} -o ${ENV_OSX}
	env GOOS=${BINARY_WINDOWS} GOARCH=${ARCH} ${GO_BUILD} -o ${ENV_WINDOWS}
test:
	${GO_TEST} -v ./...
clean:
	${GO_CLEAN}
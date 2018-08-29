GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME="eau"

build:
	${GOBUILD} -o ${GOPATH}/bin/${BINARY_NAME} -v


.DEFAULT_GOAL=help
BINARY_NAME=course
OUTPUT_DIR=_output

.PHONY: deps test build install run clean all help


deps:
	go mod download

test: deps
	go test ./...

build: deps
	go build -o ${OUTPUT_DIR}/bin/${BINARY_NAME} .

cross: deps
	GOARCH=amd64 GOOS=darwin go build -o ${OUTPUT_DIR}/bin/${BINARY_NAME}-darwin-amd64 .
	GOARCH=arm64 GOOS=darwin go build -o ${OUTPUT_DIR}/bin/${BINARY_NAME}-darwin-arm64 .
	GOARCH=amd64 GOOS=linux  go build -o ${OUTPUT_DIR}/bin/${BINARY_NAME}-linux-amd64 .
	GOARCH=arm64 GOOS=linux  go build -o ${OUTPUT_DIR}/bin/${BINARY_NAME}-linux-arm64 .

run: build
#	vault agent -config vault-agent.hcl
	${OUTPUT_DIR}/bin/${BINARY_NAME} --config course.yaml

install: build
	echo "install: ToDo"

all: build install

clean:
	go clean
	rm -rf ${OUTPUT_DIR}

help:
	@echo "Usage:"
	@echo "    make [help]:  display this info"
	@echo "    make deps:    only dowload dependencies"
	@echo "    make test:    run go test"
	@echo "    make build:   run go build"
	@echo "    make run:     execute the binary"
	@echo "    make install: install to local system"
	@echo "    make all:     build and install"
	@echo "    make clean:   run go clean and remove files"


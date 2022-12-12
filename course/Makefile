.DEFAULT_GOAL=help
BINARY_NAME=course
OUTPUT_DIR=dist

.PHONY: deps test build cross image clean help


deps:
	go mod download

test-unit: deps
	go test ./...

${OUTPUT_DIR}:
	mkdir -p ${OUTPUT_DIR}

build: deps ${OUTPUT_DIR}
	go build -o ${OUTPUT_DIR}/${BINARY_NAME} .

build-cross: deps ${OUTPUT_DIR}
	GOARCH=amd64 GOOS=darwin go build -o ${OUTPUT_DIR}/${BINARY_NAME}-darwin-amd64 .
	GOARCH=arm64 GOOS=darwin go build -o ${OUTPUT_DIR}/${BINARY_NAME}-darwin-arm64 .
	GOARCH=amd64 GOOS=linux  go build -o ${OUTPUT_DIR}/${BINARY_NAME}-linux-amd64 .
	GOARCH=arm64 GOOS=linux  go build -o ${OUTPUT_DIR}/${BINARY_NAME}-linux-arm64 .

image: ${OUTPUT_DIR}
	GOOS=linux  go build -o ${OUTPUT_DIR}/${BINARY_NAME}-linux .
	docker build -t youngshawn/course:latest -f Dockerfile .

image-cross: cross ${OUTPUT_DIR}
	echo "ToDo"

local-run: image
	echo "ToDo"
#	docker-compose up

test-integration: local-run
	echo "ToDo"

release-binary:
	echo "ToDo"

release-image:
	echo "ToDo"

release: release-binary release-image
	echo "release: binary, image"

clean:
	go clean
	rm -rf ${OUTPUT_DIR}

help:
	@echo "Usage:"
	@echo "    make [help]:  display this info"
	@echo "    make deps:    only dowload dependencies"
	@echo "    make test:    go test"
	@echo "    make build:   go build"
	@echo "    make cross:   go cross build"
	@echo "    make image:   docker build image"
	@echo "    make clean:   run go clean and remove files"

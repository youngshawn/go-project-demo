.DEFAULT_GOAL := help
BINARY_NAME := course
OUTPUT_DIR := dist
GIT_TAG ?=

.PHONY: checkout deps test build run image docker-run cross-build cross-image integration-test release-binary release-image release clean help


${OUTPUT_DIR}:
	mkdir -p ${OUTPUT_DIR}

checkout:
ifneq (${GIT_TAG},)
	git checkout ${GIT_TAG}
endif
	$(eval GIT_TAG := $(shell git describe --tags 2>/dev/null))
ifeq (${GIT_TAG},)
	$(eval GIT_TAG := latest)
endif

deps: checkout
	go mod download

test: checkout deps
	go test ./...

build: checkout deps ${OUTPUT_DIR}
	go build -o ${OUTPUT_DIR}/${BINARY_NAME}_${GIT_TAG} .

run: checkout build
	./course_${GIT_TAG} --config course.yaml

image: checkout deps ${OUTPUT_DIR}
	GOOS=linux  go build -o ${OUTPUT_DIR}/${BINARY_NAME}_${GIT_TAG}_linux .
	docker build -t youngshawn/course:${GIT_TAG} -f Dockerfile --build-arg VERSION=${GIT_TAG} .

docker-run: checkout image
	IMG_TAG=${GIT_TAG} docker-compose up -d

cross-build: checkout deps ${OUTPUT_DIR}
	GOARCH=amd64 GOOS=darwin go build -o ${OUTPUT_DIR}/${BINARY_NAME}_${GIT_TAG}_darwin_amd64 .
	GOARCH=arm64 GOOS=darwin go build -o ${OUTPUT_DIR}/${BINARY_NAME}_${GIT_TAG}_darwin_arm64 .
	GOARCH=amd64 GOOS=linux  go build -o ${OUTPUT_DIR}/${BINARY_NAME}_${GIT_TAG}_linux_amd64 .
	GOARCH=arm64 GOOS=linux  go build -o ${OUTPUT_DIR}/${BINARY_NAME}_${GIT_TAG}_linux_arm64 .

cross-image: checkout cross-build ${OUTPUT_DIR}
	echo "ToDo"

integration-test: docker-run
	echo "ToDo"

release-binary: cross-build
	echo "ToDo"

release-image: cross-image
	echo "ToDo"

release: release-binary release-image
	echo "ToDo"

clean:
	go clean
	rm -rf ${OUTPUT_DIR}

define MAKE_HELP_INFO
Usage: make <Target>

Targets:
    [help]              display this help info
    deps                download go dependency modules
    test                run go test
    build               run go build for local platform
	run                 run locally in host
	image               docker build image with Dockerfile
	docker-run          run locally in docker
	cross-build         run go build for multiple platforms
	cross-image         ToDo: build multiple platform images
	integration-test    ToDo: run integration test
	release-binary      ToDo: release binaries
	release-image       ToDo: release images
	release             ToDo: release binaries and images
	clean               run go clean and remove output files

endef
export MAKE_HELP_INFO
help:
	@echo "$$MAKE_HELP_INFO"

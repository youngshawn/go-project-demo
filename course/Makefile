.DEFAULT_GOAL := help
MODULE_NAME := course
DIST_DIR := dist
BINARY_NAME := ${MODULE_NAME}
IMAGE_NAME := ${MODULE_NAME}
GIT_TAG ?=

.PHONY: checkout deps test build run image docker-run cross-build cross-image integration-test release-binary release-image release clean help


${DIST_DIR}:
	mkdir -p ${DIST_DIR}

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

build: checkout deps ${DIST_DIR}
	go build -o ${DIST_DIR}/${BINARY_NAME}_${GIT_TAG} .

run: checkout build
	${DIST_DIR}/${BINARY_NAME}_${GIT_TAG} --config course.yaml

image: checkout deps ${DIST_DIR}
	GOOS=linux  go build -o ${DIST_DIR}/${BINARY_NAME}_${GIT_TAG}_linux .
	docker build -t youngshawn/${IMAGE_NAME}:${GIT_TAG} -f Dockerfile --build-arg GIT_TAG=${GIT_TAG}  .

docker-run: checkout image
	GIT_TAG=${GIT_TAG} docker-compose -f docker-compose.yaml up -d

cross-build: checkout deps ${DIST_DIR}
	GOARCH=amd64 GOOS=darwin go build -o ${DIST_DIR}/${BINARY_NAME}_${GIT_TAG}_darwin_amd64 .
	GOARCH=arm64 GOOS=darwin go build -o ${DIST_DIR}/${BINARY_NAME}_${GIT_TAG}_darwin_arm64 .
	GOARCH=amd64 GOOS=linux  go build -o ${DIST_DIR}/${BINARY_NAME}_${GIT_TAG}_linux_amd64 .
	GOARCH=arm64 GOOS=linux  go build -o ${DIST_DIR}/${BINARY_NAME}_${GIT_TAG}_linux_arm64 .

cross-image: checkout cross-build ${DIST_DIR}
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
	rm -rf ${DIST_DIR}

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

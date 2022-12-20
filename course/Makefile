.DEFAULT_GOAL := help
MODULE_NAME := course
DIST_DIR := dist
BIN_NAME := ${MODULE_NAME}
IMG_NAME := ${MODULE_NAME}
CHECKOUT ?=
VERSION = $(shell git describe --tags --always)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_TIME = $(shell date -R)

define LDFLAGS
"-X 'github.com/youngshawn/go-project-demo/course/config.Version=${VERSION}' \
-X 'github.com/youngshawn/go-project-demo/course/config.GitCommit=${GIT_COMMIT}' \
-X 'github.com/youngshawn/go-project-demo/course/config.BuildTime=${BUILD_TIME}'"
endef


.PHONY: show-vars checkout deps test build run image docker-run cross-build cross-image test-integration release-binary release-image release clean help

show-vars:
	@echo "VERSION = ${VERSION}"
	@echo "GIT_COMMIT = ${GIT_COMMIT}"
	@echo "BUILD_TIME = ${BUILD_TIME}"

${DIST_DIR}:
	mkdir -p ${DIST_DIR}

checkout:
ifneq (${CHECKOUT},)
	-git checkout ${CHECKOUT}
endif

deps: checkout
	go mod download

test: checkout deps
	go test ./...

build: checkout deps ${DIST_DIR}
	go build -ldflags ${LDFLAGS} -o ${DIST_DIR}/${BIN_NAME}_${VERSION} .

run: checkout build
	${DIST_DIR}/${BIN_NAME}_${VERSION} --config course.yaml

image: checkout deps ${DIST_DIR}
	GOOS=linux  go build -ldflags ${LDFLAGS} -o ${DIST_DIR}/${BIN_NAME}_${VERSION}_linux .
	docker build -t youngshawn/${IMG_NAME}:${VERSION} -f Dockerfile --build-arg BIN_VERSION=${VERSION}  .

docker-run: checkout image
	IMG_TAG=${VERSION} docker-compose -f docker-compose.yaml up -d

cross-build: checkout deps ${DIST_DIR}
	GOARCH=amd64 GOOS=darwin  go build -ldflags ${LDFLAGS} -o ${DIST_DIR}/${BIN_NAME}_${VERSION}_darwin_amd64 .
	GOARCH=arm64 GOOS=darwin  go build -ldflags ${LDFLAGS} -o ${DIST_DIR}/${BIN_NAME}_${VERSION}_darwin_arm64 .
	GOARCH=amd64 GOOS=linux   go build -ldflags ${LDFLAGS} -o ${DIST_DIR}/${BIN_NAME}_${VERSION}_linux_amd64 .
	GOARCH=arm64 GOOS=linux   go build -ldflags ${LDFLAGS} -o ${DIST_DIR}/${BIN_NAME}_${VERSION}_linux_arm64 .

cross-image: checkout cross-build ${DIST_DIR}
	@echo "ToDo"

test-integration: docker-run
	@echo "ToDo"

release-binary: cross-build
	@echo "ToDo"

release-image: cross-image
	@echo "ToDo"

release: release-binary release-image
	@echo "ToDo"

clean:
	go clean
	rm -rf ${DIST_DIR}
	-docker image rm youngshawn/${IMG_NAME}:${VERSION}
	-IMG_TAG=${VERSION} docker-compose -f docker-compose.yaml down

define MAKE_HELP_INFO
Usage: make <Target> [CHECKOUT=<branch|tag|commit>]

Targets:
	[help]              display this help info
	show-vars           show variables
	checkout            run git checkout CHECKOUT 
	deps                download go dependency modules
	test                run go test
	build               run go build for local platform
	run                 run locally in host
	image               docker build image with Dockerfile
	docker-run          run locally in docker
	cross-build         run go build for multiple platforms
	cross-image         ToDo: build multiple platform images
	test-integration    ToDo: run integration test
	release-binary      ToDo: release binaries
	release-image       ToDo: release images
	release             ToDo: release binaries and images
	clean               run go clean and remove output files
endef
export MAKE_HELP_INFO
help:
	@echo "$${MAKE_HELP_INFO}"

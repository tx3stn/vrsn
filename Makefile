BINARY_NAME=vrsn
DIR ?= ./...
PWD ?= $(shell pwd)
VERSION ?= $(shell head -n 1 VERSION)

define circleci-docker
	docker run --rm -v ${PWD}/.circleci:/repo circleci/circleci-cli:alpine 
endef

.PHONY: build
build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X github.com/tx3stn/vrsn/cmd.Version=${VERSION}" -o ${BINARY_NAME}

.PHONY: build-image
build-image:
	@docker build --tag ${BINARY_NAME}:local .

.PHONY: build-image-demo-gif
build-image-demo-gif:
	@docker build --tag ${BINARY_NAME}-vhs:demo -f ./.docker/demo-gif.Dockerfile .

.PHONY: demo-gif
demo-gif: build build-image-demo-gif
	@docker run --rm -v ${PWD}:/vhs ${BINARY_NAME}-vhs:demo /vhs/scripts/demo.tape

.PHONY: fmt
fmt:
	@go fmt ${DIR}

.PHONY: install
install: build
	@sudo cp ./${BINARY_NAME} /usr/bin/${BINARY_NAME}

.PHONY: lint
lint:
	@golangci-lint run -v ${DIR}

.PHONY: push-tag
push-tag:
	@git tag -a ${VERSION} -m "Release ${VERSION}"
	@git push origin ${VERSION}

.PHONY: test
test:
	@CGO_ENABLED=1 go test ${DIR} -race -cover

.PHONY: validate-ci
validate-ci:
	@$(circleci-docker) config validate /repo/config.yml

.PHONY: validate-orb
validate-orb:
	@$(circleci-docker) orb validate /repo/orb.yml

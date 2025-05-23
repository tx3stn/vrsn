BINARY_NAME=vrsn
DIR ?= ./...
PWD ?= $(shell pwd)
VERSION ?= $(shell head -n 1 VERSION)

define ajv-docker
	docker run --rm -v "${PWD}":/repo weibeld/ajv-cli:5.0.0 ajv --spec draft7
endef

define circleci-docker
	docker run --rm -v ${PWD}/.circleci:/repo circleci/circleci-cli:alpine
endef

.PHONY: build
build:
	@CGO_ENABLED=0 go build -ldflags "-X github.com/tx3stn/vrsn/cmd.Version=${VERSION}" -o ${BINARY_NAME}

.PHONY: build-image
build-image:
	@docker build --tag ${BINARY_NAME}:local .

.PHONY: generate-gifs
generate-gifs:
	@docker build --tag ${BINARY_NAME}-vhs:demo -f ./.docker/demo-gif.Dockerfile .
	@docker run --rm -v ${PWD}:/vhs ${BINARY_NAME}-vhs:demo /vhs/.scripts/demo.tape

.PHONY: install
install: build
	@sudo cp ./${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

.PHONY: lint
lint:
	@golangci-lint run -v ${DIR}

.PHONY: schema-validate
schema-validate:
	@$(ajv-docker) compile -s /repo/.schema/vrsn.json

.PHONY: test
test:
	@CGO_ENABLED=1 go test ${DIR} -race -cover

.PHONY: validate-ci
validate-ci:
	@$(circleci-docker) config validate /repo/config.yml

.PHONY: validate-orb
validate-orb:
	@$(circleci-docker) orb validate /repo/orb.yml

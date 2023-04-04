.PHONY: setup test

default: help

#❓ help: @ Displays all commands and tooling
help:
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#'  | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

#🔍 check: @ Runs all code verifications
check: lint lint.ci test
check.report: lint lint.ci test.report

#🔍 lint.ci: @ Strictly runs a code formatter
lint.ci:
	@go fmt ./...

lint:
	@golangci-lint run --fix

#📦 build: @ Builds and compiles dependencies
build: SHELL:=/bin/bash
build: setup
	@go build -v -o .

#📦 setup: @ Installs and compiles dependencies
setup: SHELL:=/bin/bash
setup: setup.server
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/joho/godotenv/cmd/godotenv@latest
	@go install github.com/jstemmer/go-junit-report/v2@latest

install.xunit-viewer:
	npm i -g xunit-viewer

#📦 setup.server: @ Installs and compiles the server
setup.server: SHELL:=/bin/bash
setup.server:
	@go get .

#start: @ ‍💻 Starts a server.
start: SHELL:=/bin/bash
start:
	@go run .

#🧪 test.cleanup: @ Removes all artifacts possibly left behind from previous testing
test.cleanup: SHELL:=/bin/bash
test.cleanup:
	@rm -f test-report.xml 2> /dev/null || true

#🧪 test: @ Runs all test suites
test: SHELL:=/bin/bash
test: test.cleanup
	@godotenv -f ./.env go test -count=1 -v ./...

#🧪 test.report: @ Runs all test suites and creates a test report
test.report: SHELL:=/bin/bash
test.report: test.cleanup
	@godotenv -f ./.env go test -count=1 -v ./... | go-junit-report -iocopy -out test-report.xml
	@./render-report.sh test-report.xml test-report.html

#🧪 test.ci: @ Runs all test suites
test.ci: SHELL:=/bin/bash
test.ci:
	@go test -count=1 -v ./... | go-junit-report -set-exit-code -iocopy -out test-report.xml

#🐳 docker.build: @ Builds a new local docker image
docker.build: SHELL:=/bin/bash
docker.build: TAG:=latest
docker.build:
	@echo "🐳👁️  Build nft-imx docker image"
	@source .env && docker build -t nft-imx:$(TAG) -f Dockerfile ..

#🐳 docker.rm: @ Removes a running or terminated local docker image
docker.rm: SHELL:=/bin/bash
docker.rm: TAG:=latest
docker.rm:
	@echo "🐳‍👁️  Removing docker container"
	@docker rm --force nft-imx

#🐳 docker.run: @ Runs a local docker image
docker.run: docker.rm
docker.run: SHELL:=/bin/bash
docker.run: REQUESTS_SERVICE_PORT:=4000
docker.run: TAG:=latest
docker.run: ENTRYPOINT:=./nft-imx
docker.run:
	@echo "🐳‍💻  Running a local docker image"
	@docker run --entrypoint "${ENTRYPOINT}" --name=nft-imx -p $(REQUESTS_SERVICE_PORT):4000 nft-imx:$(TAG)

#🐳 services.start: @ Starts docker services, requires a local authz image
services.start: SHELL:=/bin/bash
services.start: DOCKER_COMPOSE_FILE:=docker-compose.services.yml
services.start:
	@echo "🐳‍👁️  Starting docker-compose services"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) -p "nft-services" up -d

#🐳 services.stop: @ Stop docker services
services.stop: SHELL:=/bin/bash
services.stop: DOCKER_COMPOSE_FILE:=docker-compose.services.yml
services.stop:
	@echo "🐳‍👁️  Removing docker-compose services"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) -p "nft-services" stop && docker-compose -f $(DOCKER_COMPOSE_FILE) -p "nft-services" down --remove-orphans

#🐳 nft.build: @ Builds the nft-imx docker images
nft.build: nft.stop docker.build

#🐳 nft.start: @ Starts docker nft-imx, requires make nft.build
nft.start: SHELL:=/bin/bash
nft.start: DOCKER_COMPOSE_FILE:=docker-compose.nft.yml
nft.start:
	@echo "🐳‍👁️  Starting docker-compose services"
	@source .env && docker-compose -f $(DOCKER_COMPOSE_FILE) -p "nft-imx" up -d

#🐳 nft.stop: @ Stops docker nft-imx
nft.stop: SHELL:=/bin/bash
nft.stop: DOCKER_COMPOSE_FILE:=docker-compose.nft.yml
nft.stop:
	@echo "🐳‍👁️  Removing docker-compose services"
	@source .env && docker-compose -f $(DOCKER_COMPOSE_FILE) -p "nft-imx" stop && docker-compose -f $(DOCKER_COMPOSE_FILE) -p "nft-imx" down -v --remove-orphans

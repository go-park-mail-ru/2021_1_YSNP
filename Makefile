PROJECT_DIR := ${CURDIR}

MAIN_BINARY=main
AUTH_BINARY=auth
CHAT_BINARY=chat


DOCKER_DIR := ${CURDIR}/docker

BUILD_DOCKER=build
MAIN_DOCKER=main
AUTH_DOCKER=auth
CHAT_DOCKER=chat

## build: build go files
build:
	rm -rf ${MAIN_BINARY} && rm -rf ${AUTH_BINARY} && rm -rf ${CHAT_BINARY}
	go build -o ${MAIN_BINARY} cmd/main/main.go
	go build -o ${AUTH_BINARY} cmd/auth/main.go
	go build -o ${CHAT_BINARY} cmd/chat/main.go

## build-docker: build docker files
build-docker:
	docker build -t ${BUILD_DOCKER} -f ${DOCKER_DIR}/build.Dockerfile .
	docker build -t ${MAIN_DOCKER} -f ${DOCKER_DIR}/main.Dockerfile .
	docker build -t ${AUTH_DOCKER} -f ${DOCKER_DIR}/auth.Dockerfile .
	docker build -t ${CHAT_DOCKER} -f ${DOCKER_DIR}/chat.Dockerfile .

## run: run docker containers
run:
	docker-compose up --build --no-deps

## run-background: run docker containers background
run-background:
	docker-compose up -d --build --no-deps

## stop: stop docker containers
stop:
	docker-compose down

## rm: remove docker containers
rm:
	docker rm -vf $$(docker ps -a -q) || true

## build-and-run: build docker and run
build-and-run:
	make build-docker
	make run

## build-and-run-background: build docker and run background
build-and-run-background:
	make build-docker
	make run-background

## lint: run go liners
lint:
	golangci-lint run

## test-func: run go test func
test-func:
	go test -coverprofile cover ./... -coverpkg ./...
	go tool cover -func cover
	rm -rf cover

## test-html: run go test html
test-html:
	go test -coverprofile cover ./... -coverpkg ./...
	go tool cover -html cover
	rm -rf cover

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
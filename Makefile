# Makefile

.DEFAULT_GOAL := all

# -----------------------------------------------------------------
#        ENV VARIABLE
# -----------------------------------------------------------------

PKGS=$(shell go list ./... | grep -v /vendor/)
DOCKER_IP=$(shell if [ -z "$(DOCKER_MACHINE_NAME)" ]; then echo 'localhost'; else docker-machine ip $(DOCKER_MACHINE_NAME); fi)

# -----------------------------------------------------------------
#        Version
# -----------------------------------------------------------------

VERSION=0.0.1
BUILDDATE=$(shell date -u '+%s')
BUILDHASH=$(shell git rev-parse --short HEAD)
VERSION_FLAG=-ldflags "-X main.Version=$(VERSION) -X main.GitHash=$(BUILDHASH) -X main.BuildStmp=$(BUILDDATE)"

# -----------------------------------------------------------------
#        Main targets
# -----------------------------------------------------------------

all: clean build ## Clean and build the project

help: ## Print this message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

clean: ## Clean the project
	@go clean
	@rm -Rf .tmp *.log build/

build: format ## Build all libraries and binaries
	@go build -v $(VERSION_FLAG) -o ../../bin/todolist todolist.go

format: ## Format all packages
	@go fmt $(PKGS)

teardownTest: ## Tear down mongodb and redis for integration tests
	@$(shell docker kill todolist-mongo-test 2&>/dev/null 1&>/dev/null)
	@$(shell docker kill todolist-redis-test 2&>/dev/null 1&>/dev/null)
	@$(shell docker rm todolist-mongo-test 2&>/dev/null 1&>/dev/null)
	@$(shell docker rm todolist-redis-test 2&>/dev/null 1&>/dev/null)

setupTest: teardownTest ## Start mongodb for integration tests
	@docker run -d --name todolist-mongo-test -p "27017:27017" mongo:latest
	@docker run -d --name todolist-redis-test -p "6379:6379" redis:latest

test: setupTest ## Start tests with a mongodb and a redis docker images
	@export URL_DB=$(DOCKER_IP); go test -v $(PKGS);

start: ## Start the program
	@todolist -p 8020 -l debug -d 2

stop: ## Stop the program
	@killall todolist

# -----------------------------------------------------------------
#        Docker targets
# -----------------------------------------------------------------

dockerBuild: ## Build a docker image of the program
	docker build -t xavmarc/todolist:latest .

dockerClean: ## Remove the docker image of the program
	docker rmi -f xavmarc/todolist:latest

dockerUpRedis: ## Start the program with its redis
	@export URL_DB=$(DOCKER_IP); docker-compose -f docker-compose-redis.yml up -d

dockerUpMongo: ## Start the program with its mongodb
	@export URL_DB=$(DOCKER_IP); docker-compose -f docker-compose-mongo.yml up -d

dockerDownRedis: ## Stop the program and the redis and remove the containers
	@export URL_DB=$(DOCKER_IP); docker-compose -f docker-compose-redis.yml down

dockerDownMongo: ## Stop the program and the redis and remove the containers
	@export URL_DB=$(DOCKER_IP); docker-compose -f docker-compose-mongo.yml down

dockerBuildUpRedis: dockerDownRedis dockerBuild dockerUpRedis ## Stop, build and launch the docker images of the program

dockerBuildUpMongo: dockerDownMongo dockerBuild dockerUpMongo ## Stop, build and launch the docker images of the program

dockerWatch: ## Watch the status of the docker container
	@watch -n1 'docker ps | grep bookstore'

dockerLogsRedis: ## Print the logs of the container
	@export URL_DB=$(DOCKER_IP); docker-compose -f docker-compose-redis.yml logs -f

dockerLogsMongo: ## Print the logs of the container
	@export URL_DB=$(DOCKER_IP); docker-compose -f docker-compose-mongo.yml logs -f

.PHONY: all test clean teardownTest setupTest
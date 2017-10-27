
DOCKER := $(shell command -v docker 2> /dev/null)
DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null)

dchecks:
ifndef DOCKER
	$(error "Command is not available please install Docker")
endif
ifndef DOCKER_COMPOSE
	$(error "Command is not available please install docker-compose")
endif

setup: checks ## Docker: setup, build/rebuild app container
	cp --no-clobber sample.env .env
	docker-compose build

run: checks ## Docker: run app in container
	docker-compose up

run-cmd: checks ## Docker: run arbitary command in container, eg. `make drun-cmd go test`
	docker-compose run coverage $(filter-out $@,$(MAKECMDGOALS))

test: checks ## Docker: run unit tests in container
	docker-compose run coverage go test

%:
	@true

.PHONY: help

help:
	@echo 'Usage: make <command>'
	@echo
	@echo 'where <command> is one of the following:'
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

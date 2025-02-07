.PHONY: build up down restart logs clean help lint lint-install lint-run go-clean

DOCKER_COMPOSE = docker-compose
APP_NAME = orderbot
GOLANGCI_LINT_VERSION = v1.54.2

help:
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

go-clean:
	go clean -cache -modcache -i -r
	go mod tidy
	go mod verify

build: go-clean
	$(DOCKER_COMPOSE) build --no-cache

up:
	$(DOCKER_COMPOSE) up -d

down: 
	$(DOCKER_COMPOSE) down

restart: down up

logs:
	$(DOCKER_COMPOSE) logs -f app

logs-all:
	$(DOCKER_COMPOSE) logs -f

ps: 
	$(DOCKER_COMPOSE) ps

clean: down
	$(DOCKER_COMPOSE) down -v
	rm -rf logs/*
	$(MAKE) go-clean

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

lint-run:
	golangci-lint run --config=.golangci.yml

lint: lint-install lint-run

init:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Создан файл .env. Пожалуйста, настройте переменные окружения."; \
	else \
		echo "Файл .env уже существует."; \
	fi 
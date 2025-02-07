.PHONY: build up down restart logs clean help

DOCKER_COMPOSE = docker-compose
APP_NAME = orderbot

help:
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build:
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

init:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Создан файл .env. Пожалуйста, настройте переменные окружения."; \
	else \
		echo "Файл .env уже существует."; \
	fi 
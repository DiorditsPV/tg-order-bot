.PHONY: build up down restart logs clean help

# Переменные
DOCKER_COMPOSE = docker-compose
APP_NAME = orderbot

help: ## Показать это сообщение
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Собрать все образы
	$(DOCKER_COMPOSE) build --no-cache

up: ## Запустить все контейнеры
	$(DOCKER_COMPOSE) up -d

down: ## Остановить все контейнеры
	$(DOCKER_COMPOSE) down

restart: down up ## Перезапустить все контейнеры

logs: ## Показать логи приложения
	$(DOCKER_COMPOSE) logs -f app

logs-all: ## Показать логи всех сервисов
	$(DOCKER_COMPOSE) logs -f

ps: ## Показать статус контейнеров
	$(DOCKER_COMPOSE) ps

clean: down ## Очистить все данные (включая volumes)
	$(DOCKER_COMPOSE) down -v
	rm -rf logs/*

init: ## Инициализировать проект (создать .env если не существует)
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Создан файл .env. Пожалуйста, настройте переменные окружения."; \
	else \
		echo "Файл .env уже существует."; \
	fi 
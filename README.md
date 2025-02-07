# OrderBot

Telegram бот для приема заказов в кафе/ресторане с возможностью управления меню, корзиной и оформлением заказов.

## Возможности

- Просмотр меню по категориям
- Управление корзиной (добавление/удаление товаров)
- Расчет стоимости заказа
- Оформление заказа
- Логирование действий пользователей

## Технологии

- Go 1.21+
- PostgreSQL 15
- Redis 7
- Docker & Docker Compose

## Структура проекта

```
.
├── cmd/
│   └── main.go           # Точка входа приложения
├── internal/
│   ├── domain/          # Бизнес-модели и конфигурация
│   ├── handler/         # Обработчики Telegram API
│   ├── repository/      # Работа с хранилищем данных
│   ├── service/         # Бизнес-логика
│   └── tools/           # Вспомогательные инструменты
├── menu/                # Ресурсы меню
├── docker-compose.yml   # Конфигурация Docker
├── Dockerfile          # Сборка приложения
└── Makefile           # Команды управления
```

## Установка и запуск

1. Клонировать репозиторий:

```bash
git clone <repository-url>
cd orderbot
```

2. Создать и настроить .env файл:

```bash
make init
```

3. Запустить приложение:

```bash
make up
```

## Команды Makefile

- `make init` - Инициализация проекта
- `make up` - Запуск всех сервисов
- `make down` - Остановка всех сервисов
- `make logs` - Просмотр логов приложения
- `make logs-all` - Просмотр логов всех сервисов
- `make restart` - Перезапуск всех сервисов
- `make clean` - Очистка всех данных

## Конфигурация

Основные настройки в `.env` файле:

```env
# Bot settings
BOT_TOKEN=your_bot_token_here

# PostgreSQL settings
POSTGRES_USER=orderbot
POSTGRES_PASSWORD=your_password
POSTGRES_DB=orderbot_db

# Logging settings
LOG_TO_CHANNEL=true
LOG_CHANNEL_ID=your_channel_id
```

## Логирование

- Все действия пользователей логируются
- Логи сохраняются в файл и отправляются в Telegram канал
- Поддерживается батчевая отправка логов
- Форматирование в Markdown

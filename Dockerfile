FROM golang:1.21-alpine AS builder

WORKDIR /app

# Установка зависимостей для сборки
RUN apk add --no-cache git

# Копирование и загрузка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bot cmd/main.go

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копирование бинарного файла из предыдущего этапа
COPY --from=builder /app/bot .
COPY --from=builder /app/menu ./menu

# Создание пользователя без прав root
RUN adduser -D appuser
USER appuser

CMD ["./bot"] 
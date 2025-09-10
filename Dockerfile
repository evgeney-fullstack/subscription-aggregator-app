FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копируем собранное приложение
COPY --from=builder /app/main .
# Копируем конфигурационные файлы
COPY --from=builder /app/config.env .
COPY --from=builder /app/cert.pem .
COPY --from=builder /app/key.pem .
# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Устанавливаем зависимости для PostgreSQL клиента (если нужно)
RUN apk add --no-cache postgresql-client

EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
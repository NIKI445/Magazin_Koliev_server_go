FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum сначала для кеширования зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download
RUN go mod verify

# Копируем весь код
COPY . .

# Собираем приложение (исправлен путь)
RUN go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 3000
CMD ["./main"]
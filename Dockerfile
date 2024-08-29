# Указываем базовый образ
FROM golang:1.20-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем бинарный файл
RUN go build -o main .

# Новый этап для создания минимального образа
FROM alpine:latest

# Устанавливаем сертификаты для HTTPS-запросов
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/main .

# Копируем папку static в контейнер
COPY --from=builder /app/static ./static

# Открываем порт для приложения
EXPOSE 8080

# Команда запуска контейнера
CMD ["./main"]
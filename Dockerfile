# Используем официальный образ Golang для сборки приложения
FROM golang:1.19-bullseye as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем модульные файлы и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Сборка Go-приложения
RUN go build -o main .

# Финальный образ для выполнения
FROM debian:bullseye-slim

# Устанавливаем зависимости для Chromium
RUN apt-get update && apt-get install -y \
    chromium \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Копируем бинарный файл из стадии сборки
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

# Указываем переменные среды для Chromium
ENV ROD_BROWSER_PATH=/usr/bin/chromium

# Запуск приложения
CMD ["./main"]
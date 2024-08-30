# Указываем базовый образ
FROM golang:1.21-alpine

# Устанавливаем необходимые зависимости
RUN apk add --no-cache git chromium nss freetype ttf-freefont ca-certificates \
    harfbuzz bash udev

# Создаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем оставшиеся файлы проекта
COPY . .

# Создаем директорию для кеша
RUN mkdir -p /app/cache

# Сборка Go приложения
RUN go build -o main .

# Указываем порт, который будет открыт
EXPOSE 8080

# Команда запуска приложения
CMD ["./main"]

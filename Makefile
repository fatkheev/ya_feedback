# Название Docker Compose файла
COMPOSE_FILE=docker-compose.yml

# Название сервиса (можно указать несколько, разделенных пробелом)
SERVICES=feedback-service

# Функции Makefile

.PHONY: all build up down restart logs clean

# Сборка Docker образа
build:
	docker-compose -f $(COMPOSE_FILE) build

# Запуск контейнеров в фоне
up:
	docker-compose -f $(COMPOSE_FILE) up -d

# Остановка контейнеров
down:
	docker-compose -f $(COMPOSE_FILE) down

# Перезапуск контейнеров
restart:
	docker-compose -f $(COMPOSE_FILE) restart $(SERVICES)

# Просмотр логов
logs:
	docker-compose -f $(COMPOSE_FILE) logs -f $(SERVICES)

# Удаление всех контейнеров, образов и сетей, связанных с проектом
clean: down
	docker-compose -f $(COMPOSE_FILE) rm -f
	docker system prune -f
# Makefile

# Переменные
PROJECT_NAME=ya_feedback
SERVICE_NAME=feedback-service

# Команды
up:
	docker-compose up -d

down:
	docker-compose down

clean: down
	docker-compose rm -f

rebuild: clean
	docker-compose build

logs:
	docker-compose logs -f

.PHONY: up down clean rebuild logs
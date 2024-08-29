# Makefile

# Переменные
PROJECT_NAME=ya_feedback
SERVICE_NAME=feedback-service

# Команды
up:
	docker-compose up -d $(SERVICE_NAME)

down:
	docker-compose down

clean:
	docker-compose stop $(SERVICE_NAME)
	docker-compose rm -f $(SERVICE_NAME)
	docker rmi $(shell docker-compose images -q $(SERVICE_NAME))

rebuild: clean
	docker-compose build $(SERVICE_NAME)

logs:
	docker-compose logs -f $(SERVICE_NAME)

.PHONY: up down clean rebuild logs
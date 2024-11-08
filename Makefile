# Переменные
PROJECT_NAME=ya_feedback
SERVICE_NAME=feedback-service
PORT=8080

# Запуск контейнера
up:
	docker-compose up -d $(SERVICE_NAME)

# Остановка контейнера
down:
	docker-compose down

# Очистка (удаление контейнеров и образов)
clean:
	@echo "Stopping and removing containers..."
	docker-compose stop $(SERVICE_NAME)
	docker-compose rm -f $(SERVICE_NAME)
	@echo "Removing Docker images..."
	docker images -q $(SERVICE_NAME) | xargs -r docker rmi
	sudo rm -rf cache
	@echo "Clean completed."

# Пересборка контейнера (сначала очистка, затем сборка)
rebuild: clean
	docker-compose build $(SERVICE_NAME)

# остановка процеса (сначала очистка, затем остановка)
stop: clean
	@echo "Stopping process using port $(PORT)..."
	@sudo lsof -t -i:$(PORT) | xargs -r sudo kill -9

# Просмотр логов контейнера
logs:
	docker-compose logs -f $(SERVICE_NAME)

.PHONY: up down clean rebuild logs

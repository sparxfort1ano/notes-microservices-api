export PROJECT_ROOT=${shell pwd}
export AUTH_ROOT=${PROJECT_ROOT}/auth/
export NOTES_ROOT=${PROJECT_ROOT}/notes/

.DEFAULT_GOAL := help

DC_AUTH = docker compose -p auth-env -f docker-compose.auth.yml --env-file .env.auth
DC_NOTES = docker compose -p notes-env -f docker-compose.notes.yml --env-file .env.notes
DC_NGINX = docker compose -p gateway-env -f docker-compose.nginx.yml

auth-env-up: ## Auth-env: Launch the auth microservice environment
	@$(DC_AUTH) up -d auth-postgres

auth-env-down: ## Auth-env: Stop the microservice environment
	@$(DC_AUTH) down auth-postgres

auth-env-cleanup: ## Auth-env: Clear the microservice environment
	@read -p "Очистить все volume файлы окружения auth? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		$(DC_AUTH) down auth-postgres -v && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

auth-port-forwarder: ## Auth-env: Start the socat container for port forwarding
	@$(DC_AUTH) up -d auth-port-forward

auth-port-close: ## Auth-env: Stop the socat container
	@$(DC_AUTH) down auth-port-forward

auth-pgadmin-up: ## Auth-env: Start the PgAdmin container
	@$(DC_AUTH) up -d auth-pgadmin

auth-pgadmin-down: ## Auth-env: Stop the PgAdmin container
	@$(DC_AUTH) down auth-pgadmin


auth-run: ## Auth-Go: Execute the Go application locally (for local development and testing)
	@set -a && . ./.env.auth && set +a && \
	export POSTGRES_HOST=localhost && \
	export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs/auth && \
	cd ${AUTH_ROOT} && \
	go mod tidy && \
	go run ./cmd/

auth-deploy: ## Auth-Go: Start the Go application in the Docker Compose service (for deploying)
	@$(DC_AUTH) up -d --build auth

auth-undeploy: ## Auth-Go: Stop the Go application in the Docker Compose service
	@$(DC_AUTH) down auth


auth-api-test: ## Auth-Test: Execute script to test Auth API
	@${AUTH_ROOT}/test-scripts/api.sh


auth-go-get: ## Auth-Util: Execute go get command
	@if [ -z "$(path)" ]; then \
		echo "Отсутствует необходимый параметр path. Пример: make auth-go-get path=github.com/gin-gonic/gin"; \
		exit 1; \
	fi; \
	cd ${AUTH_ROOT} && \
	go get $(path)



notes-env-up: ## Notes-env: Launch the notes microservice environment
	@$(DC_NOTES) up -d notes-mongodb notes-redis

notes-env-down: ## Notes-env: Stop the microservice environment
	@$(DC_NOTES) down notes-mongodb notes-redis

notes-env-cleanup: ## Notes-env: Clear the microservice environment
	@read -p "Очистить все volume файлы окружения notes? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		$(DC_NOTES) down notes-mongodb notes-redis -v && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

notes-port-forwarder: ## Notes-env: Start the socat container for port forwarding
	@$(DC_NOTES) up -d notes-port-forward

notes-port-close: ## Notes-env: Stop the socat container
	@$(DC_NOTES) down notes-port-forward


notes-run: ## Notes-Go: Execute the Go application locally (for local development and testing)
	@set -a && . ./.env.notes && set +a && \
	export REDIS_HOST=localhost && \
	export MONGO_INITDB_HOST=localhost && \
	export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs/notes && \
	cd ${NOTES_ROOT} && \
	go mod tidy && \
	go run ./cmd/

notes-deploy: ## Notes-Go: Start the Go application in the Docker Compose service (for deploying)
	@$(DC_NOTES) up -d --build notes

notes-undeploy: ## Notes-Go: Stop the Go application in the Docker Compose service
	@$(DC_NOTES) down notes


notes-api-test: ## Notes-Test: Execute script to test Notes API
	@bash ${NOTES_ROOT}/test-scripts/api.sh && \
	bash ${NOTES_ROOT}/test-scripts/api2.sh

notes-cache-test: ## Notes-Test: Execute script to test Notes cache
	@bash ${NOTES_ROOT}/test-scripts/cache.sh


notes-go-get: ## Notes-Util: Execute go get command
	@if [ -z "$(path)" ]; then \
		echo "Отсутствует необходимый параметр path. Пример: make notes-go-get path=github.com/gin-gonic/gin"; \
		exit 1; \
	fi; \
	cd ${NOTES_ROOT} && \
	go get $(path)



gateway-up: ## Gateway: Launch the Nginx API gateway
	@$(DC_NGINX) up -d

gateway-down: ## Gateway: Stop the Nginx API gateway
	@$(DC_NGINX) down

gateway-test: ## Gateway: Execute global API testing through Nginx
	@bash ${PROJECT_ROOT}/nginx/global_test.sh



ps: ## Env: View running Docker Compose services
	@echo "=== Auth Services ===" && \
	${DC_AUTH} ps && \
	echo "=== Notes Services ===" && \
	${DC_NOTES} ps && \
	echo "=== Other Services ===" && \
	${DC_NGINX} ps

network-create: ## Env: Single network for multiple Docker-compose files
	@docker network create microservice-net

help: ## Show help for commands
	@echo "=== Help ==="
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
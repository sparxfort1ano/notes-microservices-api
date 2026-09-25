export PROJECT_ROOT=${shell pwd}
export AUTH_ROOT=${PROJECT_ROOT}/auth/
export NOTES_ROOT=${PROJECT_ROOT}/notes/

.DEFAULT_GOAL := help

auth-env-up: ## Auth-env: Launch the auth microservice environment
	@docker compose --env-file .env.auth up -d auth-postgres

auth-env-down: ## Auth-env: Stop the microservice environment
	@docker compose --env-file .env.auth down auth-postgres

auth-env-cleanup: ## Auth-env: Clear the microservice environment
	@read -p "Очистить все volume файлы окружения auth? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down auth-postgres -v && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

auth-port-forwarder: ## Auth-env: Start the socat container for port forwarding
	@docker compose --env-file .env.auth up -d auth-port-forward

auth-port-close: ## Auth-env: Stop the socat container
	@docker compose --env-file .env.auth down auth-port-forward

auth-pgadmin-up: ## Auth-env: Start the PgAdmin container
	@docker compose --env-file .env.auth up -d auth-pgadmin

auth-pgadmin-down: ## Auth-env: Stop the PgAdmin container
	@docker compose --env-file .env.auth down auth-pgadmin


auth-run: ## Auth-Go: Execute the Go application locally (for local development and testing)
	@set -a && . ./.env.auth && set +a && \
	export POSTGRES_HOST=localhost && \
	export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs/auth && \
	cd ${AUTH_ROOT} && \
	go mod tidy && \
	go run ./cmd/

auth-deploy: ## Auth-Go: Start the Go application in the Docker Compose service (for deploying)
	@docker compose --env-file .env.auth up -d --build auth

auth-undeploy: ## Auth-Go: Stop the Go application in the Docker Compose service
	@docker compose --env-file .env.auth down auth


auth-api-test: ## Auth-Test: Execute script to test Auth API
	@${AUTH_ROOT}/test-scripts/api.sh


auth-go-get: ## Auth-Util: Execute go get command
	@if [ -z "$(path)" ]; then \
		echo "Отсутствует необходимый параметр path. Пример: make auth-go-get path=github.com/gin-gonic/gin"; \
		exit 1; \
	fi; \
	cd ${AUTH_ROOT} && \
	go get $(path)



ps: ## Env: View running Docker Compose services
	@docker compose --env-file .env.auth ps



help: ## Show help for commands
	@echo "=== Help ==="
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
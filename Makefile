include .env
export

export PROJECT_ROOT=${shell pwd}

.DEFAULT_GOAL := help

ps: ## Env: View running Docker Compose services
	@docker compose ps



auth-env-up: ## Auth-env: Launch the auth microservice environment
	@docker compose up -d auth-postgres

auth-env-down: ## Auth-env: Stop the microservice environment
	@docker compose down auth-postgres

auth-env-cleanup: ## Auth-env: Clear the microservice environment
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down -v && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

auth-port-forwarder: ## Auth-env: Start the socat container for port forwarding
	@docker compose up -d auth-port-forward

auth-port-close: ## Auth-env: Stop the socat container
	@docker compose down auth-port-forward



help: ## Show help for commands
	@echo "=== Help ==="
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
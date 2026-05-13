# ============================================
# Development commands (docker-compose.yml)
# ============================================

dev:
	docker-compose up -d

lint:
	docker-compose exec api golangci-lint run --fix
	
# ============================================
# CLI Commands (via Docker)
# ============================================

migrate:
	@echo "🔨 Building CLI..."
	@docker-compose exec api go build -o bin/cli ./cmd/cli
	@echo "🔄 Running migrate command..."
	@docker-compose exec api ./bin/cli migrate

execute-workflow:
	@echo "🔨 Building CLI..."
	@docker-compose exec api go build -o bin/cli ./cmd/cli
	@echo "🔄 Running execute-workflow command..."
	@docker-compose exec api ./bin/cli execute-workflow
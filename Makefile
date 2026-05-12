# ============================================
# Development commands (docker-compose.yml)
# ============================================

dev:
	docker-compose up -d
	
# ============================================
# CLI Commands (via Docker)
# ============================================

migrate:
	@echo "🔨 Building CLI..."
	@docker-compose exec api go build -o bin/cli ./cmd/cli
	@echo "🔄 Running migrate command..."
	@docker-compose exec api ./bin/cli migrate
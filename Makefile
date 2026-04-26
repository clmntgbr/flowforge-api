.PHONY: dev dev-logs dev-down dev-restart dev-rebuild prod prod-logs prod-down prod-restart build clean shell tests coverage help

# ============================================
# Development commands (docker-compose.yml)
# ============================================

dev:
	docker-compose up -d

dev-logs:
	docker-compose logs -f

dev-down:
	docker-compose down

dev-restart:
	docker-compose restart

dev-rebuild:
	docker-compose down
	docker-compose up --build

# ============================================
# Production commands (docker-compose.prod.yml)
# ============================================

prod:
	docker-compose -f docker-compose.prod.yml up --build

prod-d:
	docker-compose -f docker-compose.prod.yml up -d --build

prod-logs:
	docker-compose -f docker-compose.prod.yml logs -f

prod-down:
	docker-compose -f docker-compose.prod.yml down

prod-restart:
	docker-compose -f docker-compose.prod.yml restart

prod-rebuild:
	docker-compose -f docker-compose.prod.yml down
	docker-compose -f docker-compose.prod.yml up --build

# ============================================
# Build specific images
# ============================================

build-dev:
	docker build --target development -t uplink-api:dev .

build-prod:
	docker build --target production -t uplink-api:prod .

# ============================================
# Utility commands
# ============================================

shell:
	docker-compose exec api sh

shell-prod:
	docker-compose -f docker-compose.prod.yml exec api sh

test:
	docker-compose exec api go test ./... -v

clean:
	docker-compose down -v
	docker-compose -f docker-compose.prod.yml down -v
	rm -rf tmp/
	docker system prune -f

clean-all:
	docker-compose down -v --rmi all
	docker-compose -f docker-compose.prod.yml down -v --rmi all
	rm -rf tmp/
	docker system prune -af --volumes

lint:
	docker-compose exec api golangci-lint run --fix

coverage:
	docker-compose exec api go test -coverprofile=coverage.out ./tests/handler_test -coverpkg=./handler

tests:
	docker-compose exec api go test ./tests/handler_test -v

# ============================================
# CLI Commands (via Docker)
# ============================================

start-workflow:
	@echo "🔨 Building CLI..."
	@docker-compose exec api go build -o bin/cli ./cmd/cli
	@echo "🔄 Running workflow:start command..."
	@docker-compose exec api ./bin/cli workflow:start

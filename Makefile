export PGSQL_URL=postgresql://postgres:admin123@localhost:5432/distributed_system?sslmode=disable

# =============================================================================
# Database Migration Commands
# =============================================================================
migrate-create:
	@migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	@migrate -database ${PGSQL_URL} -path migrations up

migrate-down:
	@migrate -database ${PGSQL_URL} -path migrations down

migrate-reset:
	@migrate -database ${PGSQL_URL} -path migrations reset

migrate-version:
	@migrate -database ${PGSQL_URL} -path migrations version

# =============================================================================
# Seed Commands
# =============================================================================
seed-admin:
	@echo "Running admin seed..."
	@cd cmd/seeder && go run main.go

# =============================================================================
# Build Commands
# =============================================================================
build-controller:
	@echo "Building controller service..."
	@go build -o bin/controller ./cmd/controller

build-agent:
	@echo "Building agent service..."
	@go build -o bin/agent ./cmd/agents

build-worker:
	@echo "Building worker service..."
	@go build -o bin/worker ./cmd/worker

build-seeder:
	@echo "Building seeder..."
	@go build -o bin/seeder ./cmd/seeder

build-all: build-controller build-agent build-worker build-seeder
	@echo "All services built successfully!"

# =============================================================================
# Run Commands
# =============================================================================
run-controller:
	@echo "Starting controller service on port 8080..."
	@echo "Config: config/config.yaml"
	@go run ./cmd/controller/main.go

run-agent:
	@echo "Starting agent service on port 8081..."
	@echo "Config: config/agent-config.yaml"
	@go run ./cmd/agents/main.go

run-worker:
	@echo "Starting worker service on port 8082..."
	@echo "Config: config/worker-config.yaml"
	@go run ./cmd/worker/main.go

run-seeder:
	@echo "Running admin seeder..."
	@echo "Config: config/config.yaml"
	@go run ./cmd/seeder/main.go

# =============================================================================
# Docker Commands
# =============================================================================
docker-build:
	@echo "Building Docker images..."
	@docker-compose build

docker-up:
	@echo "Starting all services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping all services..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

docker-logs-controller:
	@docker-compose logs -f controller

docker-logs-agent:
	@docker-compose logs -f agent

docker-logs-worker:
	@docker-compose logs -f worker

docker-ps:
	@docker-compose ps

docker-restart:
	@echo "Restarting all services..."
	@docker-compose restart

docker-clean:
	@echo "Cleaning up Docker resources..."
	@docker-compose down -v
	@docker system prune -f

docker-seed:
	@echo "Running admin seeder in Docker..."
	@docker-compose run --rm admin-seeder

# =============================================================================
# Development Commands
# =============================================================================
dev: docker-up
	@echo "Starting development environment..."
	@sleep 5
	@make docker-logs

dev-controller:
	@go run ./cmd/controller/main.go

dev-agent:
	@go run ./cmd/agents/main.go

# =============================================================================
# Test Commands
# =============================================================================
test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# =============================================================================
# Clean Commands
# =============================================================================
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f config.json
	@echo "Clean completed!"

clean-all: clean docker-clean
	@echo "All cleaned up!"

# =============================================================================
# Documentation Commands
# =============================================================================
docs-serve:
	@echo "Starting API Documentation Server..."
	@echo "Open browser to: http://localhost:8080"
	@echo ""
	@echo "Available views:"
	@echo "  - http://localhost:8080/index.html (Swagger UI)"
	@echo "  - http://localhost:8080/redoc.html (Redoc - Beautiful)"
	@echo "  - http://localhost:8080/swagger.yaml (Raw YAML)"
	@echo ""
	@echo "Press Ctrl+C to stop"
	@echo ""
	@python3 -m http.server 8080 --directory docs

docs-view:
	@echo "Opening Swagger UI documentation in browser..."
	@start docs/index.html || open docs/index.html

docs-redoc:
	@echo "Opening Redoc documentation in browser..."
	@start docs/redoc.html || open docs/redoc.html

docs-online:
	@echo "Opening Swagger Editor Online..."
	@echo "1. Go to: https://editor.swagger.io/"
	@echo "2. Import file: docs/swagger.yaml"
	@start https://editor.swagger.io/ || open https://editor.swagger.io/

# =============================================================================
# Setup Commands
# =============================================================================
setup: migrate-up seed-admin
	@echo "Setup completed!"
	@echo "Default admin credentials:"
	@echo "  Email: admin@distributed-system.com"
	@echo "  Password: Admin123!@#"

# =============================================================================
# Help
# =============================================================================
help:
	@echo "Distributed Configuration System - Makefile Commands"
	@echo ""
	@echo "Database Migration:"
	@echo "  make migrate-create name=migration_name    Create new migration"
	@echo "  make migrate-up                           Run all up migrations"
	@echo "  make migrate-down                         Rollback last migration"
	@echo "  make migrate-reset                        Reset all migrations"
	@echo "  make migrate-version                      Show current migration version"
	@echo ""
	@echo "Seed:"
	@echo "  make seed-admin                           Seed default admin user"
	@echo ""
	@echo "Build:"
	@echo "  make build-controller                     Build controller service"
	@echo "  make build-agent                          Build agent service"
	@echo "  make build-seeder                         Build seeder"
	@echo "  make build-all                            Build all services"
	@echo ""
	@echo "Run:"
	@echo "  make run-controller                       Run controller service"
	@echo "  make run-agent                            Run agent service"
	@echo "  make run-seeder                           Run seeder"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build                         Build Docker images"
	@echo "  make docker-up                            Start all services"
	@echo "  make docker-down                          Stop all services"
	@echo "  make docker-logs                          View all logs"
	@echo "  make docker-logs-controller               View controller logs"
	@echo "  make docker-logs-agent                    View agent logs"
	@echo "  make docker-ps                            Show running containers"
	@echo "  make docker-restart                      Restart all services"
	@echo "  make docker-seed                          Run admin seeder in Docker"
	@echo "  make docker-clean                         Remove containers and volumes"
	@echo ""
	@echo "Development:"
	@echo "  make dev                                  Start dev environment with Docker"
	@echo "  make dev-controller                       Run controller in dev mode"
	@echo "  make dev-agent                            Run agent in dev mode"
	@echo ""
	@echo "Test:"
	@echo "  make test                                 Run tests"
	@echo "  make test-coverage                        Run tests with coverage"
	@echo ""
	@echo "Documentation:"
	@echo "  make docs-serve                           Serve API docs locally (port 8080)"
	@echo "  make docs-view                            Open Swagger UI in browser"
	@echo "  make docs-redoc                           Open Redoc in browser (beautiful UI)"
	@echo "  make docs-online                          Open Swagger Editor online"
	@echo ""
	@echo "Clean:"
	@echo "  make clean                                Clean build artifacts"
	@echo "  make clean-all                            Clean everything including Docker"
	@echo ""
	@echo "Setup:"
	@echo "  make setup                                Run migrations and seed"
	@echo ""
	@echo "Help:"
	@echo "  make help                                 Show this help message"

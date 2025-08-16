.PHONY: help start-dev start-prod build-dev build-prod stop clean logs

# Default target
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make start-dev   - Start development environment (DB on port 5440)"
	@echo "  make build-dev   - Build and start development environment"
	@echo "  make dev-local   - Start dev with LocalStack for S3 simulation"
	@echo ""
	@echo "Production:"
	@echo "  make start-prod  - Start production environment (DB on port 5432)"
	@echo "  make build-prod  - Build and start production environment"
	@echo ""
	@echo "LocalStack (S3 Simulation):"
	@echo "  make localstack        - Start LocalStack and create S3 bucket"
	@echo "  make localstack-stop   - Stop LocalStack"
	@echo "  make localstack-status - Check LocalStack status"
	@echo ""
	@echo "Utilities:"
	@echo "  make stop        - Stop all containers"
	@echo "  make clean       - Stop containers and remove volumes"
	@echo "  make logs        - View logs from all services"
	@echo "  make logs-api    - View API logs"
	@echo "  make logs-db     - View database logs"
	@echo "  make test        - Run API tests"

# Start commands (without rebuild)
start-dev:
	DB_PORT_EXTERNAL=5440 docker-compose --profile dev up -d

start-prod:
	docker-compose --profile prod up -d

# Build commands (with rebuild)
build-dev:
	DB_PORT_EXTERNAL=5440 docker-compose --profile dev up --build

build-prod:
	docker-compose --profile prod up --build

# Convenience aliases
start: start-dev
build: build-dev

# Stop and cleanup commands
stop:
	docker-compose --profile dev down
	docker-compose --profile prod down

clean:
	docker-compose --profile dev down -v
	docker-compose --profile prod down -v

# Logs commands
logs:
	docker-compose logs -f

logs-api:
	docker-compose logs -f api api-dev

logs-db:
	docker-compose logs -f postgres

# Database access
db-dev:
	docker-compose exec postgres psql -U ${DB_USER} -d ${DB_NAME}

db-prod:
	docker-compose exec postgres psql -U ${DB_USER} -d ${DB_NAME}

# Test API
test:
	./test_api.sh

# LocalStack commands
localstack:
	@chmod +x localstack.sh
	./localstack.sh setup

localstack-stop:
	@chmod +x localstack.sh
	./localstack.sh stop

localstack-status:
	@chmod +x localstack.sh
	./localstack.sh status

# Development with LocalStack
dev-local:
	@echo "Starting development with LocalStack..."
	@make localstack
	LOCALSTACK_ENDPOINT=http://localhost:4566 AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test AWS_S3_BUCKET=dona-tutti-s3 DB_PORT_EXTERNAL=5440 docker-compose --profile dev up --build
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Running the Application
```bash
# Development with Docker (recommended)
docker-compose up --build

# Local development
go run main.go

# Environment setup for local development
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=microservice_user
export DB_PASSWORD=microservice_password
export DB_NAME=microservice_db
export DB_SSLMODE=disable
export API_PORT=9999
```

### Testing
```bash
# Run API tests
./test_api.sh

# Manual testing with curl examples
# See curl_examples.md for detailed examples
```

### Database Operations
```bash
# Access database in Docker
docker-compose exec postgres psql -U microservice_user -d microservice_db

# Run migrations (handled automatically on startup)
# Migrations are in migrations/ directory and run via migrations.Up() in main.go
```

### Docker Operations
```bash
# Build and start services
docker-compose up --build

# Stop services
docker-compose down

# View logs
docker-compose logs -f api
docker-compose logs -f postgres

# Restart API only
docker-compose restart api
```

## Architecture Overview

This is a Go REST API for the "Dona Tutti" donation platform, built with Echo framework and GORM ORM.

### Technology Stack
- **Framework**: Echo v4 (HTTP web framework)
- **ORM**: GORM with PostgreSQL driver
- **Database**: PostgreSQL 15
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

### Domain Structure
The application follows Domain-Driven Design with clean architecture:

```
/{domain}/
├── model.go        # Database models with GORM tags + domain entities
├── repository.go   # Data access layer
├── service.go      # Business logic layer
├── handlers.go     # HTTP handlers (Echo)
└── {domain}.go     # Domain-specific types and validation
```

**Current domains:**
- `campaign` - Main donation campaigns
- `campaigncategory` - Categories for campaigns
- `donation` - Individual donations
- `donor` - Donor profiles
- `organizer` - Campaign organizers
- `user` - User management with JWT auth

### Key Architecture Patterns

1. **Separation of Concerns**: Database models (with GORM tags) are separate from domain entities
2. **Repository Pattern**: Data access abstracted through interfaces
3. **Service Layer**: Business logic isolated from HTTP concerns
4. **Clean Mapping**: `FromEntity()` and `ToEntity()` methods convert between database and domain models

### Database Schema
- Uses UUID primary keys with `uuid_generate_v4()` extension
- Automatic timestamps with GORM `autoCreateTime`/`autoUpdateTime`
- Database constraints defined in GORM struct tags
- Migrations in `migrations/` directory, executed automatically

### API Structure
- Base path: `/api`
- Swagger docs: `/swagger/index.html`
- JWT authentication on protected endpoints
- JSON request/response format
- RESTful endpoint conventions

### Environment Configuration
Environment variables are defined in `.env` file and have defaults in `database/connection.go`:
- Database connection settings (DB_HOST, DB_PORT, etc.)
- API_PORT (default: 9999)
- JWT configuration for authentication

### RBAC System
- **Roles**: admin, donor, guest with predefined permissions
- **JWT Integration**: Role information included in token claims  
- **Middleware**: Flexible RBAC middleware supporting role, permission, and ownership-based authorization
- **Database**: Normalized schema with roles, permissions, and role_permissions tables
- **Usage**: See `RBAC_USAGE.md` and `examples/rbac_usage.go` for implementation patterns

### Development Notes
- Server runs on port 9999 by default
- Database migrations run automatically on startup (includes RBAC migration)
- GORM logging is configured for SQL query visibility
- Connection pooling configured for production use
- CORS enabled for frontend integration
- Run `./test_rbac.sh` to test RBAC implementation
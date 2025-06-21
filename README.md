# Microservice Go - Campaign Management API

A Go microservice for managing campaigns, articles, categories, and organizers using Go-kit architecture with PostgreSQL database.

## Features

- **Articles Management**: CRUD operations for articles
- **Campaign Management**: Create, read, and list campaigns
- **Category Management**: Manage campaign categories
- **Organizer Management**: Manage campaign organizers
- **PostgreSQL Integration**: Full database persistence
- **Docker Support**: Containerized application with Docker Compose
- **Go-kit Architecture**: Clean, modular architecture

## Quick Start with Docker

### Prerequisites

- Docker
- Docker Compose
- go

### 1. Clone and Setup

```bash
git clone <repository-url>
cd microservice_go
cp .env.example .env
```

### 2. Start Services

```bash
# Start PostgreSQL and API
docker-compose up -d

# View logs
docker-compose logs -f api
```

### 3. Test the API

The API will be available at `http://localhost:9999`

```bash
# List all campaigns
curl http://localhost:9999/campaigns

# List all categories
curl http://localhost:9999/categories

# List all organizers
curl http://localhost:9999/organizers

# List all articles
curl http://localhost:9999/articles
```

## API Endpoints

### Campaigns
- `GET /campaigns` - List all campaigns
- `GET /campaigns/:id` - Get specific campaign
- `POST /campaigns` - Create new campaign

### Categories
- `GET /categories` - List all categories
- `GET /categories/:id` - Get specific category

### Organizers
- `GET /organizers` - List all organizers
- `GET /organizers/:id` - Get specific organizer

### Articles
- `GET /articles` - List all articles
- `GET /articles/:id` - Get specific article
- `PUT /articles/:id` - Update article

## Development Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 15+

### 1. Install Dependencies

```bash
go mod download
```

### 2. Setup Database

```bash
# Start PostgreSQL (or use Docker)
docker run --name postgres -e POSTGRES_PASSWORD=microservice_password -e POSTGRES_USER=microservice_user -e POSTGRES_DB=microservice_db -p 5432:5432 -d postgres:15-alpine

# Run migrations
psql -h localhost -U microservice_user -d microservice_db -f migrations/001_initial_schema.sql
```

### 3. Set Environment Variables

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=microservice_user
export DB_PASSWORD=microservice_password
export DB_NAME=microservice_db
export DB_SSLMODE=disable
export API_PORT=9999
```

### 4. Run the Application

```bash
go run main.go
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `microservice_user` |
| `DB_PASSWORD` | Database password | `microservice_password` |
| `DB_NAME` | Database name | `microservice_db` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `API_PORT` | API port | `9999` |

## Docker Commands

```bash
# Build and start services
docker-compose up --build

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Restart API only
docker-compose restart api

# Access database
docker-compose exec postgres psql -U microservice_user -d microservice_db
```

## Database Schema

The application uses PostgreSQL with the following main tables:

- `articles` - Article content
- `campaigns` - Campaign information
- `campaign_categories` - Campaign categories
- `organizers` - Campaign organizers

## Testing

Use the provided `test.http` file with your HTTP client (VS Code REST Client, Postman, etc.) to test all endpoints.

## Architecture

The application follows Go-kit patterns with:

- **Domain Models**: Pure business entities
- **Repositories**: Data access layer
- **Services**: Business logic layer
- **Endpoints**: Go-kit endpoint layer
- **Transport**: HTTP transport layer

Each domain (campaign, article, etc.) is organized in separate packages following this structure.

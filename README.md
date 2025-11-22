# Dona Tutti - Sistema de Donaciones Transparentes y Auditables

Este es el backend relacionado al prototipo del sistema de donaciones transparentes y auditables. Proporciona una API REST completa para la gestión de campañas de donación, organizadores, donantes y transacciones.

## Características

- **Gestión de Campañas de Donación**: Crear, leer y administrar campañas con transparencia total
- **Gestión de Donantes**: Registro y seguimiento de donantes
- **Gestión de Donaciones**: Procesamiento y registro de transacciones
- **Gestión de Organizadores**: Administrar organizadores de campañas
- **Gestión de Categorías**: Categorización de campañas
- **Sistema de Autenticación JWT**: Seguridad basada en tokens
- **Sistema RBAC**: Control de acceso basado en roles
- **Auditoría Completa**: Registro detallado de todas las operaciones

## Inicio Rápido del sistema:

### Requisitos Previos instalados.

- Docker
- Docker Compose
- go

### 1. Clonar y Configurar

```bash
git clone <repository-url>
cd microservice_go
cp .env.example .env
```

### 2. Iniciar Servicios

```bash
# Iniciar PostgreSQL y API
docker-compose up -d

# Ver logs
docker-compose logs -f api
```

### 3. Probar la API

La API estará disponible en `http://localhost:9999`

```bash
# Listar todas las campañas
curl http://localhost:9999/campaigns

# Listar todas las categorías
curl http://localhost:9999/categories

# Listar todos los organizadores
curl http://localhost:9999/organizers

```

## Endpoints de la API

### Campañas
- `GET /campaigns` - Listar todas las campañas
- `GET /campaigns/:id` - Obtener campaña específica
- `POST /campaigns` - Crear nueva campaña

### Categorías
- `GET /categories` - Listar todas las categorías
- `GET /categories/:id` - Obtener categoría específica

### Organizadores
- `GET /organizers` - Listar todos los organizadores
- `GET /organizers/:id` - Obtener organizador específico


## Variables de Entorno

| Variable | Descripción | Por Defecto |
|----------|-------------|-------------|
| `DB_HOST` | Host de la base de datos | `localhost` |
| `DB_PORT` | Puerto de la base de datos | `5432` |
| `DB_USER` | Usuario de la base de datos | `microservice_user` |
| `DB_PASSWORD` | Contraseña de la base de datos | `microservice_password` |
| `DB_NAME` | Nombre de la base de datos | `microservice_db` |
| `DB_SSLMODE` | Modo SSL | `disable` |
| `API_PORT` | Puerto de la API | `9999` |

## Comandos Docker

```bash
# Construir e iniciar servicios
docker-compose up --build

# Detener servicios
docker-compose down

# Ver logs
docker-compose logs -f

# Reiniciar solo la API
docker-compose restart api

# Acceder a la base de datos
docker-compose exec postgres psql -U microservice_user -d microservice_db
```


## Pruebas

Utiliza el archivo `test.http` proporcionado con tu cliente HTTP (VS Code REST Client, Postman, etc.) para probar todos los endpoints.

## Tecnologías y Arquitectura

### Stack Tecnológico

- **Lenguaje**: Go 1.21+
- **Framework Web**: Echo v4
- **ORM**: GORM con driver PostgreSQL
- **Base de Datos**: PostgreSQL 15
- **Autenticación**: JWT (JSON Web Tokens)
- **Documentación API**: Swagger/OpenAPI
- **Contenedorización**: Docker & Docker Compose

### Arquitectura

La aplicación sigue los principios de **Clean Architecture** y **Domain-Driven Design (DDD)** con una estructura modular por dominios:

#### Estructura de Capas

- **Modelos de Dominio**: Entidades de negocio puras separadas de modelos de base de datos
- **Repositorios**: Capa de acceso a datos con interfaces para abstracción
- **Servicios**: Capa de lógica de negocio aislada de preocupaciones HTTP
- **Handlers**: Controladores HTTP con Echo framework
- **Middleware**: Autenticación JWT y RBAC

#### Organización por Dominios

Cada dominio (campaign, donation, donor, organizer, user) está organizado en paquetes separados con:

```
/{dominio}/
├── model.go        # Modelos de base de datos (GORM)
├── {dominio}.go    # Entidades de dominio
├── repository.go   # Capa de acceso a datos
├── service.go      # Lógica de negocio
└── handlers.go     # Controladores HTTP
```

### Patrones de Diseño

- **Repository Pattern**: Abstracción del acceso a datos
- **Service Layer Pattern**: Separación de lógica de negocio
- **DTO Pattern**: Objetos de transferencia de datos para API
- **Middleware Pattern**: Autenticación y autorización centralizadas

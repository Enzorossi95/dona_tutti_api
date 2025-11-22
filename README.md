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

**Opción A: Usando Makefile (Recomendado)**

```bash
# Desarrollo (DB en puerto 5440, con hot-reload)
make start-dev

# Producción (DB en puerto 5432)
make start-prod

# Ver todos los comandos disponibles
make help
```

**Opción B: Usando Docker Compose directamente**

```bash
# Desarrollo
DB_PORT_EXTERNAL=5440 docker-compose --profile dev up -d

# Producción
docker-compose --profile prod up -d

# Ver logs del API
docker-compose logs -f api        # Producción
docker-compose logs -f api-dev    # Desarrollo
```

> ⚠️ **Importante**: No uses solo `docker-compose up -d` sin especificar un perfil, ya que esto solo levantará PostgreSQL sin el servicio API.

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


## Perfiles de Entorno

El proyecto utiliza **Docker Compose profiles** para gestionar diferentes entornos:

### Desarrollo (`dev`)
- Base de datos expuesta en puerto **5440** (configurable con `DB_PORT_EXTERNAL`)
- Hot-reload habilitado con volúmenes montados
- Servicio: `api-dev`
- Configuración LocalStack para S3 (opcional)

### Producción (`prod`)
- Base de datos expuesta en puerto **5432** (estándar PostgreSQL)
- Sin hot-reload, imagen optimizada
- Servicio: `api`
- Configuración para AWS S3 real

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
| `DB_PORT_EXTERNAL` | Puerto externo de PostgreSQL | `5432` (prod), `5440` (dev) |

## Comandos Docker

### Usando Makefile (Recomendado)

```bash
# Iniciar servicios
make start-dev          # Desarrollo
make start-prod         # Producción

# Construir e iniciar servicios
make build-dev          # Desarrollo con rebuild
make build-prod         # Producción con rebuild

# Detener servicios
make stop               # Detiene ambos perfiles

# Ver logs
make logs               # Todos los servicios
make logs-api           # Solo API
make logs-db            # Solo base de datos

# Limpiar (detener y eliminar volúmenes)
make clean

# Ejecutar tests
make test

# Ver ayuda completa
make help
```

### Usando Docker Compose directamente

```bash
# Construir e iniciar servicios
docker-compose --profile dev up --build    # Desarrollo
docker-compose --profile prod up --build   # Producción

# Detener servicios
docker-compose --profile dev down          # Desarrollo
docker-compose --profile prod down         # Producción

# Ver logs
docker-compose logs -f api                 # API Producción
docker-compose logs -f api-dev             # API Desarrollo
docker-compose logs -f postgres            # Base de datos

# Reiniciar solo la API
docker-compose --profile prod restart api      # Producción
docker-compose --profile dev restart api-dev   # Desarrollo

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

# Distributed Configuration System

A sophisticated distributed configuration management system built with **Clean Architecture** and **Microservices pattern**, designed for centralized configuration management with real-time synchronization across distributed agents.

---

## 🏗️ System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    CONTROLLER (Port 8080)                        │
│  - Centralized configuration management                          │
│  - Admin authentication with JWT                                 │
│  - Agent registration & management                               │
│  - PostgreSQL + Redis for data & caching                        │
└────────────┬────────────────────────────────────────────────────┘
             │
             │ Polling Version (every 10s via Redis)
             │
┌────────────▼───────────────────────────────────────────────────┐
│                    AGENT (Port 8081)                            │
│  - Background service (no public HTTP endpoints)                │
│  - Version checking via Redis                                    │
│  - Fetch & push configurations to Workers                        │
│  - Local cache (config.json)                                    │
└────────────┬────────────────────────────────────────────────────┘
             │ POST /config + Internal Key
             │
┌────────────▼───────────────────────────────────────────────────┐
│                    WORKER (Port 8082)                           │
│  - Receive config from Agent                                     │
│  - Execute HTTP GET tasks to configured URLs                    │
│  - In-memory configuration storage                              │
└─────────────────────────────────────────────────────────────────┘
```

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                            │
│  HTTP Handlers, Middleware, Request/Response                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Use Case Layer                            │
│  Business Logic, Domain Rules, Orchestration                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                           │
│  Data Access, Database Operations, Caching                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                             │
│  Entities, Business Models, Interfaces                      │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔄 System Flow

### 1. Configuration Update Flow

```
Admin → Controller → PostgreSQL → Redis Cache
                    ↓
                 Update Version
                    ↓
Agent (poll every 10s) → Detect Version Change
                    ↓
              Fetch Full Config
                    ↓
              Push to Worker
                    ↓
              Worker Update Memory
```

### 2. Authentication Flow

**Admin Authentication:**
```
Admin → POST /login → Controller (verify credentials)
                      ↓
                   Generate JWT
                      ↓
                 Return Token
```

**Agent Registration:**
```
Agent → GET /agent/admin → Generate Registration Token
                        ↓
Agent → POST /agent/register (with token)
                      ↓
                 Receive Agent UUID
```

---

## 🛠️ Tech Stack

### Core Technologies

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Language** | Go | 1.24.0 | Primary programming language |
| **Web Framework** | Gin | v1.9.1 | HTTP router & middleware |
| **Database** | PostgreSQL | 15-alpine | Primary data storage |
| **Cache** | Redis | 7-alpine | Version caching & performance |
| **ORM** | GORM | v1.25.5 | Database operations |
| **Authentication** | JWT | v5.2.0 | Admin authentication |
| **Password Hashing** | bcrypt | Latest | Secure password storage |
| **Config Management** | Viper | v1.18.2 | Configuration file loading |
| **Container** | Docker | Latest | Application containerization |
| **Container Orchestration** | Docker Compose | v3.8 | Multi-container orchestration |

### Go Libraries & Dependencies

```go
// Web Framework
github.com/gin-gonic/gin v1.9.1

// Database & ORM
gorm.io/gorm v1.25.5
gorm.io/driver/postgres v1.5.4

// Authentication & Security
github.com/golang-jwt/jwt/v5 v5.2.0
golang.org/x/crypto v0.17.0

// Configuration
github.com/spf13/viper v1.18.2

// Redis
github.com/redis/go-redis/v9 v9.4.0

// Database Migration
github.com/golang-migrate/migrate/v4 v4.16.2

// UUID Generation
github.com/google/uuid v1.5.0
```

### Development Tools

| Tool | Purpose |
|------|---------|
| **Go Modules** | Dependency management |
| **Docker** | Containerization |
| **Docker Compose** | Multi-container management |
| **Make** | Build automation |
| **Git** | Version control |

---

## 🚀 Quick Start Guide

### Prerequisites

- Go 1.24.0 or higher
- Docker & Docker Compose (optional)
- PostgreSQL 15+ (if not using Docker)
- Redis 7+ (if not using Docker)

---

### Option 1: Docker Setup (Recommended) ⭐

**Fastest way to run the entire system:**

```bash
# 1. Clone the repository
git clone <repository-url>
cd distributed_system

# 2. Start all services with Docker Compose
docker-compose up -d

# 3. Wait for services to be healthy (~30 seconds)
docker-compose ps

# 4. Run database seeder (create default admin)
docker-compose run --rm admin-seeder

# 5. Check logs
docker-compose logs -f
```

**Services will be available at:**
- Controller API: http://localhost:8080
- Worker API: http://localhost:8082
- PostgreSQL: localhost:5432
- Redis: localhost:6379

**Default Admin Credentials:**
- Email: `admin@distributed-system.com`
- Password: `Admin123!@#`

---

### Option 2: Local Development Setup

**For development with hot-reload:**

```bash
# 1. Install dependencies
go mod download

# 2. Start infrastructure (PostgreSQL + Redis)
docker-compose up -d postgres redis

# 3. Run database migrations
make migrate-up

# 4. Seed admin user
make seed-admin

# 5. Start services (in separate terminals)

# Terminal 1 - Controller
make run-controller
# or: go run ./cmd/controller/main.go

# Terminal 2 - Agent
make run-agent
# or: go run ./cmd/agents/main.go

# Terminal 3 - Worker
make run-worker
# or: go run ./cmd/worker/main.go
```

---

## 📡 API Endpoints

### Controller Service (Port 8080)

#### Authentication
```bash
# Login Admin
POST /login
Content-Type: application/json

{
  "email": "admin@distributed-system.com",
  "password": "Admin123!@#"
}

Response: JWT Token
```

#### Configuration Management
```bash
# Create Configuration (Admin only)
POST /config/admin
Authorization: Bearer {JWT_TOKEN}
{
  "config_url": "https://api.example.com/task",
  "pooling_interval": 30
}

# Get Current Configuration (Admin)
GET /config/admin
Authorization: Bearer {JWT_TOKEN}

# Update Configuration (Admin)
PUT /config/admin
Authorization: Bearer {JWT_TOKEN}
{
  "config_url": "https://api.example.com/new-task",
  "pooling_interval": 60
}
```

#### Agent Management
```bash
# Generate Registration Token (Admin)
GET /agent/admin
Authorization: Bearer {JWT_TOKEN}

# Register Agent
POST /agent/register
Authorization: Bearer {REGISTRATION_TOKEN}

# Get Configuration Version (Agent)
GET /config/version
Authorization: Bearer {AGENT_TOKEN}

# Get Full Configuration (Agent)
GET /config/agent
Authorization: Bearer {AGENT_TOKEN}
```

### Worker Service (Port 8082)

```bash
# Execute Task (Public)
GET /hit

# Get Current Configuration (Public)
GET /config

# Health Check (Public)
GET /health

# Receive Config Update (Agent only)
POST /config
X-Internal-Key: {INTERNAL_KEY}
{
  "config_url": "https://api.example.com/task",
  "pooling_interval": 30,
  "version": 1,
  "uuid": "..."
}
```

---

## 🗄️ Database Schema

### Tables

**admin**
| Column | Type | Description |
|--------|------|-------------|
| uuid | TEXT (PK) | Unique identifier |
| email | TEXT | Admin email (unique) |
| password | TEXT | Bcrypt hashed password |
| created_at | TIMESTAMP | Creation timestamp |

**config**
| Column | Type | Description |
|--------|------|-------------|
| uuid | TEXT (PK) | Unique identifier |
| version | INT | Auto-increment version |
| config_url | TEXT | Target URL for task execution |
| pooling_interval | INT | Polling interval in seconds (min: 30) |
| created_at | TIMESTAMP | Creation timestamp |

**agents**
| Column | Type | Description |
|--------|------|-------------|
| id | TEXT (PK) | Agent UUID |
| created_at | TIMESTAMP | Registration timestamp |

---

## 🔒 Security Features

1. **Multi-Layer Authentication**
   - JWT token for Admin users
   - Bearer token for Agent registration
   - HMAC-SHA256 for internal service communication

2. **Password Security**
   - Bcrypt hashing with default cost factor
   - No plaintext password storage

3. **Internal Communication**
   - X-Internal-Key header for Agent → Worker communication
   - Separate internal keys per environment

4. **API Security**
   - Middleware-based authentication
   - Role-based access control (Admin vs Agent)
   - Token-based registration for new Agents

---

## 📁 Project Structure

```
distributed_system/
├── cmd/                        # Application entry points
│   ├── controller/            # Controller service (API server)
│   ├── agents/                # Agent service (background worker)
│   ├── worker/                # Worker service (task executor)
│   └── seeder/                # Database seeder (admin user)
├── internal/                   # Private application code
│   ├── domain/                # Domain entities & interfaces
│   │   ├── admin/             # Admin domain
│   │   ├── config/            # Config domain
│   │   ├── agents/            # Agent domain
│   │   └── worker/            # Worker domain
│   ├── usecase/               # Business logic layer
│   ├── repository/            # Data access layer
│   ├── delivery/              # Delivery layer (HTTP)
│   │   └── http/
│   │       ├── handler/       # HTTP handlers
│   │       └── middleware/    # Middleware
│   ├── infrastructure/        # Infrastructure
│   │   ├── database/          # Database connection
│   │   ├── cache/             # Redis cache
│   │   └── redis/             # Redis client
│   └── config/                # Configuration loading
├── pkg/                        # Public packages
│   ├── crypto/                # Cryptography utilities
│   ├── errors/                # Error handling
│   ├── response/              # HTTP response formatter
│   └── utils/                 # Utility functions
├── config/                     # Configuration files
│   ├── config.yaml            # Controller config
│   ├── agent-config.yaml      # Agent config
│   └── worker-config.yaml     # Worker config
├── migrations/                 # Database migrations
├── docker/                     # Docker files
│   ├── Dockerfile.controller
│   ├── Dockerfile.agent
│   ├── Dockerfile.worker
│   └── Dockerfile.seeder
├── docs/                       # API documentation
│   ├── swagger.yaml           # OpenAPI 3.0 spec
│   ├── index.html             # Swagger UI
│   └── redoc.html             # Redoc documentation
├── seeds/                      # Database seeders
│   └── admin_seed.go          # Admin seeder logic
├── docker-compose.yml          # Multi-container orchestration
├── Makefile                    # Build automation
├── go.mod                      # Go module definition
└── go.sum                      # Dependency checksums
```

---

## 🧪 Testing the System

### 1. Test Admin Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@distributed-system.com","password":"Admin123!@#"}'
```

### 2. Create Configuration
```bash
# Replace {JWT_TOKEN} with actual token from login
curl -X POST http://localhost:8080/config/admin \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"config_url":"https://jsonplaceholder.typicode.com/posts/1","pooling_interval":30}'
```

### 3. Test Worker
```bash
# Check worker health
curl http://localhost:8082/health

# Execute task (will hit configured URL)
curl http://localhost:8082/hit

# View current config
curl http://localhost:8082/config
```

---

## 📚 Documentation

- **API Documentation**: Open `docs/index.html` or run `make docs-view`
- **Swagger UI**: Interactive API testing at `docs/index.html`
- **Redoc**: Beautiful API docs at `docs/redoc.html`
- **OpenAPI Spec**: `docs/swagger.yaml`

---

## 🎯 Key Features

✅ **Clean Architecture**: Separation of concerns with layered architecture
✅ **Microservices**: Independent, scalable services
✅ **Real-time Sync**: Redis-based version checking for instant updates
✅ **Secure**: Multi-layer authentication with JWT & HMAC
✅ **Containerized**: Full Docker support for easy deployment
✅ **Type-Safe**: Strong typing with Go
✅ **ORM**: GORM for database operations
✅ **Caching**: Redis for performance optimization
✅ **Migration**: Database versioning with golang-migrate
✅ **Auto-seeding**: Default admin user creation

---

## 🛠️ Available Commands

```bash
# Database
make migrate-up              # Run database migrations
make seed-admin              # Create default admin user
make setup                   # Migrate + Seed

# Build
make build-all               # Build all services
make build-controller        # Build controller
make build-agent             # Build agent
make build-worker            # Build worker

# Run (Local)
make run-controller          # Run controller service
make run-agent               # Run agent service
make run-worker              # Run worker service

# Docker
make docker-up               # Start all services
make docker-down             # Stop all services
make docker-logs             # View logs
make docker-seed             # Run seeder in Docker

# Documentation
make docs-view               # Open Swagger UI
make docs-redoc              # Open Redoc
make docs-serve              # Serve docs locally

# Help
make help                    # Show all commands
```

---

## 🌐 Architecture Patterns Used

| Pattern | Implementation | Purpose |
|---------|---------------|---------|
| **Clean Architecture** | Layered structure (Delivery/UseCase/Repository/Domain) | Separation of concerns, testability |
| **Repository Pattern** | Abstract data access with interfaces | Decouple business logic from data layer |
| **Dependency Injection** | Constructor injection throughout | Loose coupling, easy testing |
| **Middleware Pattern** | Gin middleware for auth/validation | Cross-cutting concerns |
| **Service Layer Pattern** | UseCase layer for business logic | Encapsulate business rules |
| **Background Worker** | Agent as daemon process | Async task processing |
| **Version Caching** | Redis for config version tracking | Performance optimization |

---

## 💡 Design Decisions

### Why Redis for Versioning?
- **Performance**: ~1-5ms vs ~50-100ms HTTP overhead
- **Scalability**: Handle 1000+ agents efficiently
- **Separation**: Version tracking decoupled from API
- **Future-proof**: Can upgrade to Pub/Sub for real-time

### Why PostgreSQL as Primary Database?
- **ACID Compliance**: Strong data consistency
- **Relationship Support**: Foreign keys, constraints
- **Mature**: Battle-tested, reliable
- **SQL**: Powerful querying capabilities

### Why Clean Architecture?
- **Testability**: Easy to mock dependencies
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Swap implementations without changing business logic
- **Scalability**: Easy to add features

---

## 📊 System Scalability

| Component | Horizontal Scaling | Vertical Scaling |
|-----------|-------------------|------------------|
| Controller | ✅ Yes (load balancer) | ✅ Yes |
| Agent | ✅ Yes (multiple instances) | ✅ Yes |
| Worker | ✅ Yes (multiple instances) | ✅ Yes |
| PostgreSQL | ✅ Yes (read replicas) | ✅ Yes |
| Redis | ✅ Yes (cluster mode) | ✅ Yes |

---

## 🔧 Configuration

All configuration is YAML-based and loaded via Viper:

- **Controller**: `config/config.yaml`
- **Agent**: `config/agent-config.yaml`
- **Worker**: `config/worker-config.yaml`

Environment variables can override config values.

---

## 📝 Notes for Recruiters

- **Code Quality**: Follows Go best practices and Effective Go guidelines
- **Error Handling**: Comprehensive error handling with wrapped errors
- **Logging**: Structured logging throughout the system
- **API Design**: RESTful principles with OpenAPI 3.0 specification
- **Security**: Defense in depth with multiple authentication layers
- **Testing**: Unit test ready with dependency injection
- **Documentation**: Comprehensive inline comments and API docs
- **Containerization**: Production-ready Docker setup
- **Orchestration**: Docker Compose for local development

---

## 🚀 Deployment Ready

The system is containerized and ready for deployment to:

- **Docker Swarm**
- **Kubernetes** (with K8s manifests)
- **Cloud Platforms**: AWS, GCP, Azure
- **PaaS**: Heroku, DigitalOcean App Platform

---

## 📧 Contact

For questions or technical discussions about this project, please reach out through the repository issues or contact channels.

---

**Built with ❤️ using Go, Clean Architecture, and Microservices**

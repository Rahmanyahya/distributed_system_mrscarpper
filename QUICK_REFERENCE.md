# Distributed System - Quick Reference Card

## 🎯 Project Snapshot

| Aspect | Details |
|--------|---------|
| **Project Type** | Distributed Configuration Management System |
| **Architecture** | Clean Architecture + Microservices |
| **Language** | Go 1.24.0 |
| **Services** | 3 (Controller, Agent, Worker) |
| **Database** | PostgreSQL 15 + Redis 7 |
| **Deployment** | Docker + Docker Compose |

---

## ⚡ 30-Second Setup

```bash
# Docker (Fastest)
git clone <repo> && cd distributed_system
docker-compose up -d
docker-compose run --rm admin-seeder
# Done! Services at localhost:8080, localhost:8082
```

---

## 🏗️ Architecture at a Glance

```
Admin ──► Controller (8080) ──► PostgreSQL + Redis
                                    │
                                    ▼
Agent (Background) ◄──── Polling (10s)
                                    │
                                    ▼
                                Worker (8082)
```

**Key Points:**
- **Controller**: Central API, config management, auth
- **Agent**: Background worker, syncs config to workers
- **Worker**: Executes HTTP tasks, stores config in memory

---

## 📦 Tech Stack Summary

| Layer | Technology | Version |
|-------|-----------|---------|
| **Language** | Go | 1.24.0 |
| **Web Framework** | Gin | v1.9.1 |
| **Database** | PostgreSQL | 15-alpine |
| **Cache** | Redis | 7-alpine |
| **ORM** | GORM | v1.25.5 |
| **Auth** | JWT | v5.2.0 |
| **Crypto** | bcrypt | golang.org/x/crypto |
| **Config** | Viper | v1.18.2 |
| **Container** | Docker | Latest |
| **Orchestration** | Docker Compose | v3.8 |

---

## 🔑 Key Endpoints

### Controller (8080)
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| POST | `/login` | Public | Admin login |
| POST | `/config/admin` | JWT | Create config |
| GET | `/config/admin` | JWT | Get config |
| PUT | `/config/admin` | JWT | Update config |
| GET | `/agent/admin` | JWT | Generate agent token |
| POST | `/agent/register` | Token | Register agent |
| GET | `/config/version` | Agent | Get version |
| GET | `/config/agent` | Agent | Get full config |

### Worker (8082)
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/health` | Health check |
| GET | `/hit` | Execute task |
| GET | `/config` | Get current config |
| POST | `/config` | Receive config (Agent only) |

---

## 🧪 Quick Test

```bash
# 1. Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@distributed-system.com","password":"Admin123!@#"}'

# 2. Create Config (replace JWT_TOKEN)
curl -X POST http://localhost:8080/config/admin \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"config_url":"https://jsonplaceholder.typicode.com/posts/1","pooling_interval":30}'

# 3. Test Worker
curl http://localhost:8082/hit
```

---

## 🗄️ Database Schema

**3 Tables:**
1. **admin** (uuid, email, password, created_at)
2. **config** (uuid, version, config_url, pooling_interval, created_at)
3. **agents** (id, created_at)

---

## 🔐 Security Layers

1. **Admin** → JWT Token authentication
2. **Agent** → Bearer Token registration
3. **Worker** → X-Internal-Key header
4. **Passwords** → Bcrypt hashing

---

## 📁 Folder Structure

```
cmd/          # Entry points (controller, agent, worker, seeder)
internal/
  domain/     # Entities & interfaces
  usecase/    # Business logic
  repository/ # Data access
  delivery/   # HTTP handlers & middleware
config/       # YAML configs
migrations/   # DB migrations
docker/       # Dockerfiles
docs/         # API documentation (Swagger/Redoc)
```

---

## 🚀 Make Commands

```bash
make docker-up      # Start all services
make docker-seed    # Create admin user
make docs-view      # Open API docs
make help           # Show all commands
```

---

## ✨ Highlights

- ✅ Clean Architecture (testable, maintainable)
- ✅ Microservices (scalable, independent)
- ✅ Real-time sync (Redis-based versioning)
- ✅ Secure (JWT, bcrypt, HMAC)
- ✅ Containerized (Docker ready)
- ✅ Documented (OpenAPI 3.0)
- ✅ Type-safe (Go + GORM)

---

## 📚 Documentation

- Full README: `README.md`
- Recruiter Guide: `RECRUITER_README.md`
- API Docs: `docs/index.html` (Swagger UI)
- OpenAPI Spec: `docs/swagger.yaml`

---

**Built for scale, designed for maintainability.**

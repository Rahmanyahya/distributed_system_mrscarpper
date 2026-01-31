# 🚀 Setup & Running Guide

## ✅ What Was Fixed

All entry points now support flexible config paths:
- ✅ **Controller** (`cmd/controller/main.go`) - Fixed config path
- ✅ **Agent** (`cmd/agents/main.go`) - Fixed config path
- ✅ **Worker** (`cmd/worker/main.go`) - Fixed config path
- ✅ **Seeder** (`cmd/seeder/main.go`) - Fixed config path

**All services now:**
- Use `CONFIG_PATH` environment variable (optional)
- Default to `config` directory when run from project root
- Work with `make run-*` commands

---

## 🎯 Quick Start (2 Ways)

### Option 1: Docker (Recommended - Easiest) ⭐

```bash
# 1. Start all services
docker-compose up -d

# 2. Wait for services to be healthy (~30 seconds)
docker-compose ps

# 3. Run database seeder
docker-compose run --rm admin-seeder

# 4. Check logs
docker-compose logs -f
```

**Done!** All services running at:
- Controller: http://localhost:8080
- Worker: http://localhost:8082
- PostgreSQL: localhost:5432
- Redis: localhost:6379

---

### Option 2: Local Development (For Coding)

**Step 1: Start Infrastructure**
```bash
# Start PostgreSQL and Redis only
docker-compose up -d postgres redis
```

**Step 2: Setup Database**
```bash
# Run migrations
make migrate-up

# Or manually:
migrate -database postgresql://postgres:admin123@localhost:5432/distributed_system?sslmode=disable -path migrations up

# Create admin user
make seed-admin
```

**Step 3: Start Services (3 Terminals)**

```bash
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

## 📋 Available Commands

### Infrastructure
```bash
# Start infrastructure only
docker-compose up -d postgres redis

# Start all services (Docker)
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f
```

### Database
```bash
make migrate-up      # Run migrations
make seed-admin      # Create default admin
make setup           # Migrate + Seed
```

### Run Services (Local)
```bash
make run-controller  # Port 8080
make run-agent       # Background service
make run-worker      # Port 8082
make run-seeder      # Create admin user
```

### Docker
```bash
make docker-up       # Start all services
make docker-down     # Stop all services
make docker-seed     # Run seeder in Docker
make docker-logs     # View logs
```

### Build
```bash
make build-all       # Build all services
make build-controller
make build-agent
make build-worker
```

---

## 🧪 Test the System

### 1. Test Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@distributed-system.com","password":"Admin123!@#"}'
```

**Expected Response:**
```json
{
  "data": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 2. Create Configuration
```bash
# Replace JWT_TOKEN with token from login
curl -X POST http://localhost:8080/config/admin \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"config_url":"https://jsonplaceholder.typicode.com/posts/1","pooling_interval":30}'
```

### 3. Test Worker
```bash
# Health check
curl http://localhost:8082/health

# Execute task (after config is pushed by Agent)
curl http://localhost:8082/hit

# View current config
curl http://localhost:8082/config
```

---

## 🔧 Configuration Files

All configs are in `config/` directory:

- **config.yaml** - Controller config (database, redis, server)
- **agent-config.yaml** - Agent config (controller URL, worker URL, internal keys)
- **worker-config.yaml** - Worker config (port, internal key)

**Environment Override:**
```bash
# Override default config path
CONFIG_PATH=/path/to/config go run ./cmd/controller/main.go
```

---

## ⚠️ Troubleshooting

### Port Already in Use
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -i :8080
kill -9 <PID>
```

### Database Connection Error
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Restart
docker-compose restart postgres
```

### Redis Connection Error
```bash
# Check if Redis is running
docker-compose ps redis

# Check logs
docker-compose logs redis

# Restart
docker-compose restart redis
```

### Config Not Found
```bash
# Make sure you're in project root
pwd
# Should be: /path/to/distributed_system

# Check config files exist
ls -la config/
# Should see: config.yaml, agent-config.yaml, worker-config.yaml
```

### Migration Failed
```bash
# Drop and recreate database
docker-compose down -v
docker-compose up -d postgres

# Wait a few seconds, then run migrations
make migrate-up
make seed-admin
```

---

## 📊 Service Architecture

```
┌─────────────────────────────────────────┐
│         Controller (8080)               │
│  - Admin API                            │
│  - Config Management                    │
│  - Agent Registration                   │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│      Agent (Background)                 │
│  - Polls Controller every 10s           │
│  - Pushes config to Worker              │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│         Worker (8082)                   │
│  - Executes HTTP tasks                  │
│  - Stores config in memory              │
└─────────────────────────────────────────┘
```

---

## 🎯 Default Credentials

**Admin User:**
- Email: `admin@distributed-system.com`
- Password: `Admin123!@#`

**Internal Keys:**
- Agent Internal Key: `bintang_kejora_agent`
- Worker Internal Key: `worker_internal_key_2024`

---

## 📚 More Documentation

- **Full README**: `README.md`
- **Recruiter Guide**: `RECRUITER_README.md`
- **Quick Reference**: `QUICK_REFERENCE.md`
- **API Docs**: Run `make docs-view`

---

## 🆘 Still Having Issues?

1. **Check Docker**: `docker-compose ps`
2. **Check Logs**: `docker-compose logs -f`
3. **Check Ports**: Nothing else using 8080, 8082, 5432, 6379
4. **Clean Restart**: `docker-compose down -v && docker-compose up -d`

---

**Happy Coding! 🚀**

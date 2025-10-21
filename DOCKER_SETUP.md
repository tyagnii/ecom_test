# Docker Development Environment

## ✅ **Complete Docker Setup**

I have successfully created a comprehensive Docker development environment with all the services you requested:

### **🐳 Services Implemented:**

#### 1. **PostgreSQL Database**
- **Image**: postgres:15-alpine
- **Port**: 5432
- **Database**: ecom_test
- **Credentials**: postgres/postgres
- **Features**: Health checks, data persistence, initialization scripts

#### 2. **Migration Service**
- **Purpose**: Runs database migrations automatically
- **Dependencies**: Waits for PostgreSQL to be healthy
- **Behavior**: Runs once and exits after successful migration

#### 3. **API Service**
- **Purpose**: Runs the banner click tracking API
- **Port**: 8080
- **Dependencies**: Waits for migration service to complete
- **Features**: Health checks, auto-restart, optimized build

#### 4. **Adminer (Bonus)**
- **Purpose**: Web-based database administration
- **Port**: 8081
- **URL**: http://localhost:8081

### **📁 Files Created:**

#### **Docker Configuration:**
- `Dockerfile` - Multi-stage build for the application
- `dev/docker-compose.yml` - Main Docker Compose configuration
- `dev/docker-compose.override.yml` - Development overrides
- `.dockerignore` - Docker build exclusions

#### **Development Scripts:**
- `dev/start.sh` - Start development environment
- `dev/stop.sh` - Stop development environment  
- `dev/test-api.sh` - API testing script
- `dev/README.md` - Comprehensive development guide

### **🚀 Quick Start:**

```bash
# 1. Navigate to the project directory
cd /home/tas/Yandex.Disk/git/ecom_test

# 2. Start the development environment
./dev/start.sh

# 3. Test the API
./dev/test-api.sh

# 4. Stop when done
./dev/stop.sh
```

### **🔧 Service Architecture:**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │    Migration    │    │   API Service   │
│   (Port 5432)   │◄───┤   (One-time)    │◄───┤   (Port 8080)   │
│                 │    │                 │    │                 │
│ - Database      │    │ - Runs migrations│    │ - REST API      │
│ - Health checks │    │ - Waits for DB  │    │ - Health checks │
│ - Data persist  │    │ - Exits after   │    │ - Auto-restart  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **📋 Available Endpoints:**

Once running, you can access:

- **API Server**: http://localhost:8080
  - `GET /health` - Health check
  - `GET /api/v1/counter/<bannerID>` - Record click
  - `POST /api/v1/stats/<bannerID>` - Get statistics

- **Database Admin**: http://localhost:8081
  - Web interface for database management
  - Server: postgres, Username: postgres, Password: postgres

- **PostgreSQL**: localhost:5432
  - Direct database access
  - Database: ecom_test

### **🧪 Testing Commands:**

```bash
# Test health endpoint
curl http://localhost:8080/health

# Create a test banner (via database)
docker-compose exec postgres psql -U postgres -d ecom_test -c "INSERT INTO banners (name) VALUES ('Test Banner');"

# Test counter endpoint
curl http://localhost:8080/api/v1/counter/1

# Test stats endpoint
curl -X POST http://localhost:8080/api/v1/stats/1 \
  -H "Content-Type: application/json" \
  -d '{"banner_id": 1, "ts_from": "2025-01-01T00:00:00Z", "ts_to": "2025-01-31T23:59:59Z"}'
```

### **🔍 Development Features:**

#### **Health Monitoring:**
- PostgreSQL health checks
- API health endpoint
- Service dependency management
- Automatic restart policies

#### **Data Persistence:**
- PostgreSQL data volume
- Migration state tracking
- Development data seeding

#### **Development Tools:**
- Hot reloading support (via override file)
- Debug logging enabled
- Database admin interface
- Comprehensive testing scripts

### **📊 Service Dependencies:**

```
postgres (healthy) → migrate (success) → api (running)
```

1. **PostgreSQL** starts and becomes healthy
2. **Migration** runs after PostgreSQL is ready
3. **API** starts after migration completes successfully

### **🛠️ Development Commands:**

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f api
docker-compose logs -f postgres

# Restart specific service
docker-compose restart api

# Connect to database
docker-compose exec postgres psql -U postgres -d ecom_test

# Run migrations manually
docker-compose exec api ./main migrate

# Clean restart (removes all data)
docker-compose down -v && docker-compose up -d
```

### **🔒 Security Features:**

- Non-root user in containers
- Minimal Alpine Linux base images
- Health checks for all services
- Network isolation
- Volume persistence for data

### **📈 Production Considerations:**

The setup includes production-ready features:
- Multi-stage Docker builds
- Health checks
- Proper logging
- Resource optimization
- Security best practices

### **🎉 Ready to Use!**

The Docker development environment is now complete and ready for development. Simply run `./dev/start.sh` to get started with:

- ✅ PostgreSQL database with migrations
- ✅ API server with all endpoints
- ✅ Database administration interface
- ✅ Comprehensive testing tools
- ✅ Development scripts and documentation

Your banner click tracking system is now fully containerized and ready for development! 🚀

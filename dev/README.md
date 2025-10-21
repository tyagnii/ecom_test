# Development Environment

This directory contains Docker Compose configuration and scripts for the E-commerce Test development environment.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose installed
- Git (for cloning the repository)

### Start Development Environment
```bash
# Make scripts executable
chmod +x dev/*.sh

# Start all services
./dev/start.sh
```

### Stop Development Environment
```bash
./dev/stop.sh
```

## 📦 Services

### 1. PostgreSQL Database
- **Port**: 5432
- **Database**: ecom_test
- **Username**: postgres
- **Password**: postgres
- **Health Check**: Built-in PostgreSQL health check

### 2. Migration Service
- **Purpose**: Runs database migrations automatically
- **Dependencies**: Waits for PostgreSQL to be healthy
- **Behavior**: Runs once and exits

### 3. API Service
- **Port**: 8080
- **Health Check**: HTTP endpoint `/health`
- **Dependencies**: Waits for migration service to complete
- **Restart Policy**: Unless stopped

### 4. Adminer (Database Admin)
- **Port**: 8081
- **Purpose**: Web-based database administration
- **URL**: http://localhost:8081

## 🔧 Development Commands

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f postgres
```

### Database Operations
```bash
# Connect to database
docker-compose exec postgres psql -U postgres -d ecom_test

# Run migrations manually
docker-compose exec api ./main migrate

# Create test data
docker-compose exec postgres psql -U postgres -d ecom_test -c "INSERT INTO banners (name) VALUES ('Test Banner');"
```

### API Testing
```bash
# Run automated tests
./dev/test-api.sh

# Manual testing
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/counter/1
```

### Service Management
```bash
# Restart specific service
docker-compose restart api

# Rebuild and restart
docker-compose up --build -d api

# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

## 🧪 Testing the API

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. Create Test Banner
```bash
# Connect to database and create a banner
docker-compose exec postgres psql -U postgres -d ecom_test -c "INSERT INTO banners (name) VALUES ('Test Banner');"
```

### 3. Test Counter Endpoint
```bash
# Record a click
curl http://localhost:8080/api/v1/counter/1
```

### 4. Test Stats Endpoint
```bash
curl -X POST http://localhost:8080/api/v1/stats/1 \
  -H "Content-Type: application/json" \
  -d '{
    "banner_id": 1,
    "ts_from": "2025-01-01T00:00:00Z",
    "ts_to": "2025-01-31T23:59:59Z"
  }'
```

## 📁 File Structure

```
dev/
├── docker-compose.yml    # Docker Compose configuration
├── start.sh             # Start development environment
├── stop.sh              # Stop development environment
├── test-api.sh          # API testing script
└── README.md            # This file
```

## 🔍 Troubleshooting

### Services Not Starting
```bash
# Check Docker status
docker info

# Check service logs
docker-compose logs

# Restart services
docker-compose down && docker-compose up -d
```

### Database Connection Issues
```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Test database connection
docker-compose exec postgres pg_isready -U postgres -d ecom_test
```

### API Not Responding
```bash
# Check API logs
docker-compose logs api

# Test API health
curl http://localhost:8080/health

# Restart API service
docker-compose restart api
```

### Port Conflicts
If you have port conflicts, modify the ports in `docker-compose.yml`:
```yaml
ports:
  - "8080:8080"  # Change 8080 to another port
  - "5432:5432"  # Change 5432 to another port
```

## 🚀 Production Deployment

For production deployment, consider:
- Using environment variables for sensitive data
- Setting up proper logging
- Configuring reverse proxy (nginx)
- Setting up monitoring and alerting
- Using managed database services

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [API Documentation](../api_examples.md)
- [Migration Guide](../MIGRATIONS.md)

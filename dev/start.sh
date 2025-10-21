#!/bin/bash

# Development environment startup script
set -e

echo "🚀 Starting E-commerce Test Development Environment"
echo "=================================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

echo "📦 Building and starting services..."
docker-compose up --build -d

echo "⏳ Waiting for services to be ready..."
sleep 10

echo "🔍 Checking service health..."

# Check PostgreSQL
echo "  - PostgreSQL:"
if docker-compose exec postgres pg_isready -U postgres -d ecom_test > /dev/null 2>&1; then
    echo "    ✅ PostgreSQL is ready"
else
    echo "    ❌ PostgreSQL is not ready"
fi

# Check API
echo "  - API Service:"
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "    ✅ API is ready"
else
    echo "    ❌ API is not ready"
fi

echo ""
echo "🎉 Development environment is ready!"
echo ""
echo "📋 Available services:"
echo "  - API Server:     http://localhost:8080"
echo "  - Database Admin: http://localhost:8081"
echo "  - PostgreSQL:     localhost:5432"
echo ""
echo "🔧 Useful commands:"
echo "  - View logs:       docker-compose logs -f"
echo "  - Stop services:  docker-compose down"
echo "  - Restart API:     docker-compose restart api"
echo "  - Database shell: docker-compose exec postgres psql -U postgres -d ecom_test"
echo ""
echo "🧪 Test the API:"
echo "  curl http://localhost:8080/health"
echo "  curl http://localhost:8080/api/v1/counter/1"

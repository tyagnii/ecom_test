#!/bin/bash

# Development environment stop script
set -e

echo "🛑 Stopping E-commerce Test Development Environment"
echo "=================================================="

# Stop and remove containers
echo "📦 Stopping services..."
docker-compose down

echo "🧹 Cleaning up..."
# Remove unused images (optional)
# docker image prune -f

echo "✅ Development environment stopped successfully!"
echo ""
echo "💡 To start again, run: ./dev/start.sh"
echo "💡 To remove all data, run: docker-compose down -v"

#!/bin/bash

# API testing script for development environment
set -e

echo "🧪 Testing E-commerce Test API"
echo "============================="

API_URL="http://localhost:8080"

# Check if API is running
echo "🔍 Checking API health..."
if ! curl -f "$API_URL/health" > /dev/null 2>&1; then
    echo "❌ API is not running. Please start the development environment first:"
    echo "   ./dev/start.sh"
    exit 1
fi

echo "✅ API is running"

# Test health endpoint
echo ""
echo "1. Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s "$API_URL/health")
echo "   Response: $HEALTH_RESPONSE"

# Test counter endpoint (this will fail if no banners exist, which is expected)
echo ""
echo "2. Testing counter endpoint..."
echo "   Note: This will fail if no banners exist in the database"
COUNTER_RESPONSE=$(curl -s -w "%{http_code}" "$API_URL/api/v1/counter/1" || echo "404")
echo "   Response: $COUNTER_RESPONSE"

# Test stats endpoint
echo ""
echo "3. Testing stats endpoint..."
STATS_PAYLOAD='{
  "banner_id": 1,
  "ts_from": "2025-01-01T00:00:00Z",
  "ts_to": "2025-01-31T23:59:59Z"
}'
STATS_RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d "$STATS_PAYLOAD" "$API_URL/api/v1/stats/1" || echo "404")
echo "   Response: $STATS_RESPONSE"

echo ""
echo "🎉 API testing completed!"
echo ""
echo "💡 To create test data, connect to the database:"
echo "   docker-compose exec postgres psql -U postgres -d ecom_test"
echo "   INSERT INTO banners (name) VALUES ('Test Banner');"
echo ""
echo "💡 Then test the counter endpoint:"
echo "   curl $API_URL/api/v1/counter/1"

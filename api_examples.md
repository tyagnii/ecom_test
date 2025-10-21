# API Documentation

This document describes the REST API endpoints for the banner click tracking system.

## Base URL
```
http://localhost:8080
```

## Endpoints

### 1. Counter Endpoint
**Record a click for a banner**

```
GET /api/v1/counter/<bannerID>
```

**Description:** Increments the click counter for the specified banner by 1.

**Parameters:**
- `bannerID` (path parameter): The ID of the banner to record a click for

**Response:**
```json
{
  "banner_id": 1,
  "click_count": 42,
  "timestamp": "2025-01-27T10:30:00Z",
  "message": "Click recorded successfully"
}
```

**Status Codes:**
- `200 OK`: Click recorded successfully
- `400 Bad Request`: Invalid banner ID
- `404 Not Found`: Banner not found
- `500 Internal Server Error`: Server error

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/counter/1
```

### 2. Stats Endpoint
**Get banner statistics for a time period**

```
POST /api/v1/stats/<bannerID>
```

**Description:** Returns statistics for a banner within a specified time period.

**Parameters:**
- `bannerID` (path parameter): The ID of the banner to get statistics for

**Request Body:**
```json
{
  "banner_id": 1,
  "ts_from": "2025-01-01T00:00:00Z",
  "ts_to": "2025-01-31T23:59:59Z"
}
```

**Response:**
```json
{
  "banner_id": 1,
  "total_clicks": 150,
  "first_click": "2025-01-01T08:30:00Z",
  "last_click": "2025-01-31T18:45:00Z",
  "period_start": "2025-01-01T00:00:00Z",
  "period_end": "2025-01-31T23:59:59Z",
  "clicks_in_period": 45
}
```

**Status Codes:**
- `200 OK`: Statistics retrieved successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Banner not found
- `500 Internal Server Error`: Server error

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/stats/1 \
  -H "Content-Type: application/json" \
  -d '{
    "banner_id": 1,
    "ts_from": "2025-01-01T00:00:00Z",
    "ts_to": "2025-01-31T23:59:59Z"
  }'
```

### 3. Health Check
**Check API server health**

```
GET /health
```

**Description:** Returns the health status of the API server.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-27T10:30:00Z",
  "version": "1.0.0"
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/health
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error Type",
  "message": "Detailed error message"
}
```

**Common Error Types:**
- `Invalid banner ID`: Banner ID is not a valid number
- `Banner not found`: Banner with the specified ID doesn't exist
- `Invalid request body`: Request JSON is malformed
- `Invalid time range`: Time range parameters are invalid
- `Internal server error`: Unexpected server error

## Usage Examples

### Recording Clicks

```bash
# Record a click for banner ID 1
curl -X GET http://localhost:8080/api/v1/counter/1

# Record a click for banner ID 2
curl -X GET http://localhost:8080/api/v1/counter/2
```

### Getting Statistics

```bash
# Get stats for banner 1 for the last 30 days
curl -X POST http://localhost:8080/api/v1/stats/1 \
  -H "Content-Type: application/json" \
  -d '{
    "banner_id": 1,
    "ts_from": "2025-01-01T00:00:00Z",
    "ts_to": "2025-01-31T23:59:59Z"
  }'

# Get stats for banner 2 for today
curl -X POST http://localhost:8080/api/v1/stats/2 \
  -H "Content-Type: application/json" \
  -d '{
    "banner_id": 2,
    "ts_from": "2025-01-27T00:00:00Z",
    "ts_to": "2025-01-27T23:59:59Z"
  }'
```

### JavaScript Examples

```javascript
// Record a click
async function recordClick(bannerId) {
  const response = await fetch(`http://localhost:8080/api/v1/counter/${bannerId}`);
  const data = await response.json();
  return data;
}

// Get statistics
async function getStats(bannerId, fromDate, toDate) {
  const response = await fetch(`http://localhost:8080/api/v1/stats/${bannerId}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      banner_id: bannerId,
      ts_from: fromDate,
      ts_to: toDate
    })
  });
  const data = await response.json();
  return data;
}

// Usage
recordClick(1).then(data => console.log('Click recorded:', data));
getStats(1, '2025-01-01T00:00:00Z', '2025-01-31T23:59:59Z')
  .then(data => console.log('Stats:', data));
```

### Python Examples

```python
import requests
import json
from datetime import datetime

# Record a click
def record_click(banner_id):
    url = f"http://localhost:8080/api/v1/counter/{banner_id}"
    response = requests.get(url)
    return response.json()

# Get statistics
def get_stats(banner_id, from_date, to_date):
    url = f"http://localhost:8080/api/v1/stats/{banner_id}"
    data = {
        "banner_id": banner_id,
        "ts_from": from_date,
        "ts_to": to_date
    }
    response = requests.post(url, json=data)
    return response.json()

# Usage
click_data = record_click(1)
print("Click recorded:", click_data)

stats_data = get_stats(1, "2025-01-01T00:00:00Z", "2025-01-31T23:59:59Z")
print("Stats:", stats_data)
```

## Starting the API Server

### Using CLI
```bash
# Start with default port (8080)
go run main.go api

# Start with custom port
go run main.go api --port 3000
```

### Environment Variables
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=ecom_test
export DB_SSLMODE=disable
```

## Performance Considerations

- The API server includes connection pooling for database operations
- Each request is logged with timing information
- Graceful shutdown is supported (SIGTERM, SIGINT)
- Timeouts are configured for read (15s), write (15s), and idle (60s) operations

## Security Notes

- Input validation is performed on all parameters
- SQL injection protection through parameterized queries
- Error messages don't expose internal system details
- CORS headers can be added if needed for web applications

## Monitoring

- Health check endpoint for monitoring systems
- Request logging with timing information
- Error logging for debugging
- Graceful shutdown handling

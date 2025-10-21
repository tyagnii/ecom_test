# Database Migrations

This project includes a PostgreSQL migration system for managing database schema changes.

## Migration Files

The migration files are located in `db/migrations/`:

- `001_create_banners_table.sql` - Creates the banners table
- `002_create_clicks_table.sql` - Creates the clicks table with foreign key to banners

## Database Schema

### Banners Table
```sql
CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Clicks Table
```sql
CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    bannerid INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bannerid) REFERENCES banners(id) ON DELETE CASCADE ON UPDATE CASCADE
);
```

## Running Migrations

### Prerequisites

1. Install dependencies:
```bash
go mod tidy
```

2. Set up your PostgreSQL database and configure connection:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=ecom_test
export DB_SSLMODE=disable
```

### Commands

#### Run all pending migrations:
```bash
go run main.go migrate
```

#### Check migration status:
```bash
go run main.go migrate status
```

#### Run migrations with custom database settings:
```bash
go run main.go migrate --host=localhost --port=5432 --user=postgres --password=your_password --dbname=ecom_test
```

## Migration System Features

- **Version Control**: Each migration has a version number (001, 002, etc.)
- **Transaction Safety**: Each migration runs in a transaction
- **Tracking**: Applied migrations are tracked in `schema_migrations` table
- **Idempotent**: Safe to run multiple times
- **Ordered Execution**: Migrations run in version order

## Creating New Migrations

1. Create a new SQL file in `db/migrations/` with the format: `XXX_description.sql`
2. Use the next sequential number (e.g., `003_add_indexes.sql`)
3. Write your SQL DDL statements
4. Run migrations with `go run main.go migrate`

## Example Usage in Code

```go
package main

import (
    "database/sql"
    "log"
    
    "github.com/tyagnii/ecom_test/db"
    "github.com/tyagnii/ecom_test/db/migrations"
)

func main() {
    // Get database configuration
    config := db.GetConfigFromEnv()
    
    // Connect to database
    database, err := db.Connect(config)
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()
    
    // Run migrations
    migrator := migrations.NewMigrator(database)
    if err := migrator.RunMigrations("db/migrations"); err != nil {
        log.Fatal(err)
    }
    
    // Your application code here...
}
```

## Troubleshooting

- Ensure PostgreSQL is running and accessible
- Check that the database exists
- Verify connection parameters
- Check PostgreSQL logs for detailed error messages

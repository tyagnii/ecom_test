package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Version string
	Name    string
	Up      string
	Down    string
}

// Migrator handles database migrations
type Migrator struct {
	db *sql.DB
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// CreateMigrationsTable creates the migrations tracking table
func (m *Migrator) CreateMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	
	_, err := m.db.Exec(query)
	return err
}

// GetAppliedMigrations returns a list of applied migration versions
func (m *Migrator) GetAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)
	
	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	
	return applied, nil
}

// RunMigrations executes all pending migrations
func (m *Migrator) RunMigrations(migrationsDir string) error {
	// Create migrations table if it doesn't exist
	if err := m.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	// Get applied migrations
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	
	// Read migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	var migrations []Migration
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") && !strings.HasPrefix(file.Name(), ".") {
			version := strings.Split(file.Name(), "_")[0]
			if _, exists := applied[version]; !exists {
				content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
				if err != nil {
					return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
				}
				
				migrations = append(migrations, Migration{
					Version: version,
					Name:    file.Name(),
					Up:      string(content),
				})
			}
		}
	}
	
	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	
	// Execute pending migrations
	for _, migration := range migrations {
		log.Printf("Running migration: %s", migration.Name)
		
		// Start transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", migration.Name, err)
		}
		
		// Execute migration
		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration.Name, err)
		}
		
		// Record migration as applied
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}
		
		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Name, err)
		}
		
		log.Printf("Successfully applied migration: %s", migration.Name)
	}
	
	if len(migrations) == 0 {
		log.Println("No pending migrations found")
	} else {
		log.Printf("Applied %d migrations", len(migrations))
	}
	
	return nil
}

// GetMigrationStatus returns the status of all migrations
func (m *Migrator) GetMigrationStatus(migrationsDir string) error {
	// Get applied migrations
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	
	// Read migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	fmt.Println("Migration Status:")
	fmt.Println("================")
	
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") && !strings.HasPrefix(file.Name(), ".") {
			version := strings.Split(file.Name(), "_")[0]
			status := "PENDING"
			if applied[version] {
				status = "APPLIED"
			}
			fmt.Printf("%s: %s (%s)\n", version, file.Name(), status)
		}
	}
	
	return nil
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/tyagnii/ecom_test/db/migrations"
)

var (
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
	dbSSLMode  string
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long: `Run database migrations for the application.
This command will execute all pending SQL migrations in the migrations directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

// migrateStatusCmd represents the migrate status command
var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  `Show the status of all database migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		showMigrationStatus()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateStatusCmd)

	// Database connection flags
	migrateCmd.PersistentFlags().StringVar(&dbHost, "host", "localhost", "Database host")
	migrateCmd.PersistentFlags().IntVar(&dbPort, "port", 5432, "Database port")
	migrateCmd.PersistentFlags().StringVar(&dbUser, "user", "postgres", "Database user")
	migrateCmd.PersistentFlags().StringVar(&dbPassword, "password", "", "Database password")
	migrateCmd.PersistentFlags().StringVar(&dbName, "dbname", "ecom_test", "Database name")
	migrateCmd.PersistentFlags().StringVar(&dbSSLMode, "sslmode", "disable", "SSL mode")
}

func getDatabaseURL() string {
	if dbPassword == "" {
		dbPassword = os.Getenv("DB_PASSWORD")
	}
	
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
}

func connectToDatabase() (*sql.DB, error) {
	dbURL := getDatabaseURL()
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}

func runMigrations() {
	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	
	// Get migrations directory
	migrationsDir := filepath.Join("db", "migrations")
	
	// Create migrator and run migrations
	migrator := migrations.NewMigrator(db)
	
	log.Println("Starting database migrations...")
	if err := migrator.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	
	log.Println("Migrations completed successfully!")
}

func showMigrationStatus() {
	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	
	// Get migrations directory
	migrationsDir := filepath.Join("db", "migrations")
	
	// Create migrator and show status
	migrator := migrations.NewMigrator(db)
	
	if err := migrator.GetMigrationStatus(migrationsDir); err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/tyagnii/ecom_test/api"
)

var (
	apiPort int
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Long:  `Start the HTTP API server for banner click tracking.`,
	Run: func(cmd *cobra.Command, args []string) {
		startAPIServer()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().IntVarP(&apiPort, "port", "p", 8080, "Port to run the API server on")
}

func startAPIServer() {
	// Connect to database
	database, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create API server
	server := api.NewServer(database)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down API server...")
		if err := server.Stop(); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
		os.Exit(0)
	}()

	// Start server
	log.Printf("Starting API server on port %d", apiPort)
	if err := server.Start(apiPort); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}

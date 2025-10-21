package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tyagnii/ecom_test/app"
	"github.com/tyagnii/ecom_test/cache"
	"github.com/tyagnii/ecom_test/db"
)

// Server represents the API server
type Server struct {
	handler *APIHandler
	server  *http.Server
}

// NewServer creates a new API server
func NewServer(database *sql.DB) *Server {
	// Create repository and service
	repo := db.NewRepository(database)
	service := app.NewService(repo)
	
	// Create cache and cached repository
	cacheInstance := cache.NewInMemoryCache(cache.DefaultCleanupInterval)
	cachedRepo := cache.NewCachedRepository(repo, cacheInstance)
	
	// Create API handler with cached repository
	handler := NewAPIHandler(service, cachedRepo)
	
	return &Server{
		handler: handler,
	}
}

// Start starts the API server
func (s *Server) Start(port int) error {
	// Setup routes
	mux := s.handler.SetupRoutes()
	
	// Create HTTP server
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Printf("Starting API server on port %d", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET  /api/v1/counter/<bannerID>  - Record a click for a banner")
	log.Printf("  POST /api/v1/stats/<bannerID>    - Get banner statistics")
	log.Printf("  GET  /health                     - Health check")
	
	return s.server.ListenAndServe()
}

// Stop stops the API server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// GetHandler returns the API handler (for testing)
func (s *Server) GetHandler() *APIHandler {
	return s.handler
}

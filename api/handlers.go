package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tyagnii/ecom_test/app"
	"github.com/tyagnii/ecom_test/db"
)

// APIHandler provides HTTP API handlers
type APIHandler struct {
	service *app.Service
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(service *app.Service) *APIHandler {
	return &APIHandler{service: service}
}

// CounterRequest represents a counter request
type CounterRequest struct {
	BannerID int `json:"banner_id"`
}

// CounterResponse represents a counter response
type CounterResponse struct {
	BannerID   int       `json:"banner_id"`
	ClickCount int       `json:"click_count"`
	Timestamp  time.Time `json:"timestamp"`
	Message    string    `json:"message"`
}

// StatsRequest represents a stats request
type StatsRequest struct {
	BannerID int       `json:"banner_id"`
	TsFrom   time.Time `json:"ts_from"`
	TsTo     time.Time `json:"ts_to"`
}

// StatsResponse represents a stats response
type StatsResponse struct {
	BannerID    int       `json:"banner_id"`
	TotalClicks int       `json:"total_clicks"`
	FirstClick  time.Time `json:"first_click,omitempty"`
	LastClick   time.Time `json:"last_click,omitempty"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	ClicksInPeriod int    `json:"clicks_in_period"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// CounterHandler handles GET /api/v1/counter/<bannerID>
func (h *APIHandler) CounterHandler(w http.ResponseWriter, r *http.Request) {
	// Extract banner ID from URL path
	bannerIDStr := r.URL.Path[len("/api/v1/counter/"):]
	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid banner ID", "Banner ID must be a number")
		return
	}

	// Validate banner ID
	if bannerID <= 0 {
		h.sendError(w, http.StatusBadRequest, "Invalid banner ID", "Banner ID must be positive")
		return
	}

	// Check if banner exists
	bannerService := app.NewBannerService(h.service)
	_, err = bannerService.GetBanner(bannerID)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Banner not found", fmt.Sprintf("Banner with ID %d not found", bannerID))
		return
	}

	// Record the click
	clickService := app.NewClickService(h.service)
	click, err := clickService.RecordClick(bannerID, time.Now())
	if err != nil {
		log.Printf("Failed to record click for banner %d: %v", bannerID, err)
		h.sendError(w, http.StatusInternalServerError, "Failed to record click", "Internal server error")
		return
	}

	// Get updated click count for this banner
	stats, err := clickService.GetClickStats(bannerID)
	if err != nil {
		log.Printf("Failed to get click stats for banner %d: %v", bannerID, err)
		// Don't fail the request, just use the click we recorded
		stats = &db.ClickStats{
			BannerID:    bannerID,
			TotalClicks: 1,
		}
	}

	// Prepare response
	response := CounterResponse{
		BannerID:   bannerID,
		ClickCount: stats.TotalClicks,
		Timestamp:  click.Timestamp,
		Message:    "Click recorded successfully",
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// StatsHandler handles POST /api/v1/stats/<bannerID>
func (h *APIHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract banner ID from URL path
	bannerIDStr := r.URL.Path[len("/api/v1/stats/"):]
	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid banner ID", "Banner ID must be a number")
		return
	}

	// Validate banner ID
	if bannerID <= 0 {
		h.sendError(w, http.StatusBadRequest, "Invalid banner ID", "Banner ID must be positive")
		return
	}

	// Parse request body
	var req StatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body", "Failed to parse JSON")
		return
	}

	// Validate request
	if req.TsFrom.IsZero() || req.TsTo.IsZero() {
		h.sendError(w, http.StatusBadRequest, "Invalid time range", "ts_from and ts_to are required")
		return
	}

	if req.TsFrom.After(req.TsTo) {
		h.sendError(w, http.StatusBadRequest, "Invalid time range", "ts_from must be before ts_to")
		return
	}

	// Check if banner exists
	bannerService := app.NewBannerService(h.service)
	_, err = bannerService.GetBanner(bannerID)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Banner not found", fmt.Sprintf("Banner with ID %d not found", bannerID))
		return
	}

	// Get click service
	clickService := app.NewClickService(h.service)

	// Get overall stats for the banner
	overallStats, err := clickService.GetClickStats(bannerID)
	if err != nil {
		log.Printf("Failed to get overall stats for banner %d: %v", bannerID, err)
		h.sendError(w, http.StatusInternalServerError, "Failed to get stats", "Internal server error")
		return
	}

	// Get clicks in the specified time period
	clicksInPeriod, err := clickService.GetClicksForBannerInDateRange(bannerID, req.TsFrom, req.TsTo)
	if err != nil {
		log.Printf("Failed to get clicks in period for banner %d: %v", bannerID, err)
		h.sendError(w, http.StatusInternalServerError, "Failed to get period stats", "Internal server error")
		return
	}

	// Prepare response
	response := StatsResponse{
		BannerID:      bannerID,
		TotalClicks:   overallStats.TotalClicks,
		PeriodStart:   req.TsFrom,
		PeriodEnd:     req.TsTo,
		ClicksInPeriod: len(clicksInPeriod),
	}

	// Add first and last click times if available
	if overallStats.TotalClicks > 0 {
		response.FirstClick = overallStats.FirstClick
		response.LastClick = overallStats.LastClick
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// HealthHandler handles health check
func (h *APIHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendError sends an error response
func (h *APIHandler) sendError(w http.ResponseWriter, statusCode int, error, message string) {
	response := ErrorResponse{
		Error:   error,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SetupRoutes sets up the API routes
func (h *APIHandler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/v1/counter/", h.CounterHandler)
	mux.HandleFunc("/api/v1/stats/", h.StatsHandler)
	mux.HandleFunc("/health", h.HealthHandler)

	// Add middleware for logging
	return h.addLoggingMiddleware(mux)
}

// addLoggingMiddleware adds logging middleware
func (h *APIHandler) addLoggingMiddleware(handler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
	return mux
}

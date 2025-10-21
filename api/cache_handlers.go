package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tyagnii/ecom_test/cache"
)

// CacheStatsResponse represents cache statistics response
type CacheStatsResponse struct {
	Stats cache.CacheStats `json:"stats"`
}

// CacheManagementHandler provides cache management endpoints
type CacheManagementHandler struct {
	cachedRepo *cache.CachedRepository
}

// NewCacheManagementHandler creates a new cache management handler
func NewCacheManagementHandler(cachedRepo *cache.CachedRepository) *CacheManagementHandler {
	return &CacheManagementHandler{
		cachedRepo: cachedRepo,
	}
}

// GetCacheStatsHandler handles GET /api/v1/cache/stats
func (h *CacheManagementHandler) GetCacheStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := h.cachedRepo.GetCacheStats()
	
	response := CacheStatsResponse{
		Stats: stats,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ClearCacheHandler handles POST /api/v1/cache/clear
func (h *CacheManagementHandler) ClearCacheHandler(w http.ResponseWriter, r *http.Request) {
	h.cachedRepo.ClearCache()
	
	response := map[string]interface{}{
		"message": "Cache cleared successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// WarmCacheHandler handles POST /api/v1/cache/warm
func (h *CacheManagementHandler) WarmCacheHandler(w http.ResponseWriter, r *http.Request) {
	err := h.cachedRepo.WarmCache()
	if err != nil {
		response := map[string]interface{}{
			"error":   "Failed to warm cache",
			"message": err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"message": "Cache warmed successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// InvalidateBannerCacheHandler handles POST /api/v1/cache/banner/{id}/invalidate
func (h *CacheManagementHandler) InvalidateBannerCacheHandler(w http.ResponseWriter, r *http.Request) {
	// Extract banner ID from URL path
	bannerIDStr := r.URL.Path[len("/api/v1/cache/banner/"):]
	bannerIDStr = bannerIDStr[:len(bannerIDStr)-len("/invalidate")]
	
	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		response := map[string]interface{}{
			"error":   "Invalid banner ID",
			"message": "Banner ID must be a number",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	h.cachedRepo.InvalidateBannerCache(bannerID)

	response := map[string]interface{}{
		"message":    "Banner cache invalidated successfully",
		"banner_id":  bannerID,
		"status":     "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SetupCacheRoutes sets up cache management routes
func (h *CacheManagementHandler) SetupCacheRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/cache/stats", h.GetCacheStatsHandler)
	mux.HandleFunc("/api/v1/cache/clear", h.ClearCacheHandler)
	mux.HandleFunc("/api/v1/cache/warm", h.WarmCacheHandler)
	mux.HandleFunc("/api/v1/cache/banner/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[len("/api/v1/cache/banner/"):] == "" {
			http.NotFound(w, r)
			return
		}
		h.InvalidateBannerCacheHandler(w, r)
	})
}

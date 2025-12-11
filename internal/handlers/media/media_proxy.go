package media

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"MrRSS/internal/cache"
	"MrRSS/internal/handlers/core"
	"MrRSS/internal/utils"
)

// validateMediaURL validates that the URL is HTTP/HTTPS and properly formatted
func validateMediaURL(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("URL must use HTTP or HTTPS")
	}

	return nil
}

// HandleMediaProxy serves cached media or downloads and caches it
func HandleMediaProxy(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if media cache is enabled
	mediaCacheEnabled, _ := h.DB.GetSetting("media_cache_enabled")
	if mediaCacheEnabled != "true" {
		http.Error(w, "Media cache is disabled", http.StatusForbidden)
		return
	}

	// Get URL from query parameter
	mediaURL := r.URL.Query().Get("url")
	if mediaURL == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	// Validate mediaURL (must be HTTP/HTTPS and valid format)
	if err := validateMediaURL(mediaURL); err != nil {
		http.Error(w, "Invalid url parameter: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get optional referer from query parameter
	referer := r.URL.Query().Get("referer")

	// Get media cache directory
	cacheDir, err := utils.GetMediaCacheDir()
	if err != nil {
		log.Printf("Failed to get media cache directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Initialize media cache
	mediaCache, err := cache.NewMediaCache(cacheDir)
	if err != nil {
		log.Printf("Failed to initialize media cache: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get media (from cache or download)
	data, contentType, err := mediaCache.Get(mediaURL, referer)
	if err != nil {
		log.Printf("Failed to get media %s: %v", mediaURL, err)
		http.Error(w, "Failed to fetch media", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year

	// Write response
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write media response: %v", err)
	}
}

// HandleMediaCacheCleanup performs manual cleanup of media cache
func HandleMediaCacheCleanup(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get media cache directory
	cacheDir, err := utils.GetMediaCacheDir()
	if err != nil {
		log.Printf("Failed to get media cache directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Initialize media cache
	mediaCache, err := cache.NewMediaCache(cacheDir)
	if err != nil {
		log.Printf("Failed to initialize media cache: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get settings
	maxAgeDaysStr, _ := h.DB.GetSetting("media_cache_max_age_days")
	maxSizeMBStr, _ := h.DB.GetSetting("media_cache_max_size_mb")

	maxAgeDays, err := strconv.Atoi(maxAgeDaysStr)
	if err != nil || maxAgeDays <= 0 {
		maxAgeDays = 7 // Default
	}

	maxSizeMB, err := strconv.Atoi(maxSizeMBStr)
	if err != nil || maxSizeMB <= 0 {
		maxSizeMB = 100 // Default
	}

	// Cleanup by age
	ageCount, err := mediaCache.CleanupOldFiles(maxAgeDays)
	if err != nil {
		log.Printf("Failed to cleanup old media files: %v", err)
	}

	// Cleanup by size
	sizeCount, err := mediaCache.CleanupBySize(maxSizeMB)
	if err != nil {
		log.Printf("Failed to cleanup media files by size: %v", err)
	}

	totalCleaned := ageCount + sizeCount
	log.Printf("Media cache cleanup: removed %d files", totalCleaned)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success":       true,
		"files_cleaned": totalCleaned,
	}
	json.NewEncoder(w).Encode(response)
}

// HandleMediaCacheInfo returns information about the media cache
func HandleMediaCacheInfo(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get media cache directory
	cacheDir, err := utils.GetMediaCacheDir()
	if err != nil {
		log.Printf("Failed to get media cache directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Initialize media cache
	mediaCache, err := cache.NewMediaCache(cacheDir)
	if err != nil {
		log.Printf("Failed to initialize media cache: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get cache size
	cacheSize, err := mediaCache.GetCacheSize()
	if err != nil {
		log.Printf("Failed to get cache size: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert to MB
	cacheSizeMB := float64(cacheSize) / (1024 * 1024)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"cache_size_mb": cacheSizeMB,
	}
	json.NewEncoder(w).Encode(response)
}

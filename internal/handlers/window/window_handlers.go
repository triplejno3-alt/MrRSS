package window

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"MrRSS/internal/handlers/core"
)

// WindowState represents the saved window state
type WindowState struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	Maximized bool `json:"maximized"`
}

// HandleGetWindowState returns the saved window state from database
func HandleGetWindowState(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	x, err := h.DB.GetSetting("window_x")
	if err != nil {
		// Return empty state if not found
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"x":         "",
			"y":         "",
			"width":     "",
			"height":    "",
			"maximized": "",
		})
		return
	}
	y, _ := h.DB.GetSetting("window_y")
	width, _ := h.DB.GetSetting("window_width")
	height, _ := h.DB.GetSetting("window_height")
	maximized, _ := h.DB.GetSetting("window_maximized")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"x":         x,
		"y":         y,
		"width":     width,
		"height":    height,
		"maximized": maximized,
	})
}

// HandleSaveWindowState saves the current window state to database
func HandleSaveWindowState(h *core.Handler, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var state WindowState
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to strings for database storage and check for errors
	if err := h.DB.SetSetting("window_x", fmt.Sprintf("%d", state.X)); err != nil {
		log.Printf("Failed to save window_x: %v", err)
		http.Error(w, "Failed to save window state", http.StatusInternalServerError)
		return
	}
	if err := h.DB.SetSetting("window_y", fmt.Sprintf("%d", state.Y)); err != nil {
		log.Printf("Failed to save window_y: %v", err)
		http.Error(w, "Failed to save window state", http.StatusInternalServerError)
		return
	}
	if err := h.DB.SetSetting("window_width", fmt.Sprintf("%d", state.Width)); err != nil {
		log.Printf("Failed to save window_width: %v", err)
		http.Error(w, "Failed to save window state", http.StatusInternalServerError)
		return
	}
	if err := h.DB.SetSetting("window_height", fmt.Sprintf("%d", state.Height)); err != nil {
		log.Printf("Failed to save window_height: %v", err)
		http.Error(w, "Failed to save window state", http.StatusInternalServerError)
		return
	}
	if err := h.DB.SetSetting("window_maximized", fmt.Sprintf("%t", state.Maximized)); err != nil {
		log.Printf("Failed to save window_maximized: %v", err)
		http.Error(w, "Failed to save window state", http.StatusInternalServerError)
		return
	}

	log.Printf("Window state saved: x=%d, y=%d, width=%d, height=%d, maximized=%t",
		state.X, state.Y, state.Width, state.Height, state.Maximized)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

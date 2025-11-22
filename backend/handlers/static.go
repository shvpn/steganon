package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// ServeStatic serves frontend static files
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Serve frontend files
	frontendDir := "../frontend"

	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	fullPath := filepath.Join(frontendDir, path)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, fullPath)
}

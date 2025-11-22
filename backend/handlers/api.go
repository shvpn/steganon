package handlers

import (
	"encoding/json"
	"image"
	"image/png"
	"log"
	"net/http"
	"steganography/utils"
)

// EncodeRequest represents the encoding request structure
type EncodeRequest struct {
	Message string `json:"message"`
}

// DecodeResponse represents the decoding response structure
type DecodeResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleEncode processes image encoding requests
func HandleEncode(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10 MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		sendJSONError(w, "Failed to parse form")
		return
	}

	// Get the image file
	file, _, err := r.FormFile("image")
	if err != nil {
		sendJSONError(w, "No image file provided")
		return
	}
	defer file.Close()

	// Get the message
	message := r.FormValue("message")
	if message == "" {
		sendJSONError(w, "No message provided")
		return
	}

	// Get optional password
	password := r.FormValue("password")

	// Encrypt message if password provided
	encryptedMessage, err := utils.EncryptMessage(message, password)
	if err != nil {
		sendJSONError(w, "Failed to encrypt message")
		return
	}

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		sendJSONError(w, "Failed to decode image")
		return
	}

	// Encode encrypted message into image
	encodedImg := utils.EncodeMessageInImage(img, string(encryptedMessage))

	// Send the encoded image back
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=encoded.png")

	if err := png.Encode(w, encodedImg); err != nil {
		log.Printf("Failed to encode PNG: %v", err)
	}
}

// HandleDecode processes image decoding requests
func HandleDecode(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10 MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.Header().Set("Content-Type", "application/json")
		sendJSONError(w, "Failed to parse form")
		return
	}

	// Get the image file
	file, _, err := r.FormFile("image")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		sendJSONError(w, "No image file provided")
		return
	}
	defer file.Close()

	// Get optional password
	password := r.FormValue("password")

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		sendJSONError(w, "Failed to decode image")
		return
	}

	// Decode message from image
	encryptedMessage := utils.DecodeMessageFromImage(img)

	// Decrypt message if password provided
	message, err := utils.DecryptMessage([]byte(encryptedMessage), password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		sendJSONError(w, err.Error())
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DecodeResponse{Message: message})
}

// enableCORS sets CORS headers for cross-origin requests
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// sendJSONError sends a JSON error response
func sendJSONError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

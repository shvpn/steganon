package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type EncodeRequest struct {
	Message string `json:"message"`
}

type DecodeResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Encrypt message using AES with password
func encryptMessage(message string, password string) ([]byte, error) {
	if password == "" {
		// No password, return plain message
		return []byte(message), nil
	}

	// Hash password with SHA-256 to get 32-byte key
	hash := sha256.Sum256([]byte(password))
	key := hash[:]

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and append nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(message), nil)
	return ciphertext, nil
}

// Decrypt message using AES with password
func decryptMessage(encrypted []byte, password string) (string, error) {
	if password == "" {
		// No password, return as plain text
		return string(encrypted), nil
	}

	// Hash password with SHA-256 to get 32-byte key
	hash := sha256.Sum256([]byte(password))
	key := hash[:]

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("invalid password or corrupted data")
	}

	return string(plaintext), nil
}

// Encode hides a message in an image using LSB steganography
func encodeMessage(img image.Image, message string) *image.RGBA {
	bounds := img.Bounds()
	encoded := image.NewRGBA(bounds)

	// Copy original image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			encoded.Set(x, y, img.At(x, y))
		}
	}

	// Add length delimiter
	messageBytes := []byte(message)
	messageLen := len(messageBytes)
	dataToHide := append([]byte(fmt.Sprintf("%08d", messageLen)), messageBytes...)

	bitIndex := 0
	totalBits := len(dataToHide) * 8

	for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalBits; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalBits; x++ {
			r, g, b, a := encoded.At(x, y).RGBA()

			// Convert to 8-bit values
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Encode in RGB channels (3 bits per pixel)
			if bitIndex < totalBits {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := (dataToHide[byteIndex] >> (7 - bitOffset)) & 1
				r8 = (r8 & 0xFE) | bit
				bitIndex++
			}

			if bitIndex < totalBits {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := (dataToHide[byteIndex] >> (7 - bitOffset)) & 1
				g8 = (g8 & 0xFE) | bit
				bitIndex++
			}

			if bitIndex < totalBits {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := (dataToHide[byteIndex] >> (7 - bitOffset)) & 1
				b8 = (b8 & 0xFE) | bit
				bitIndex++
			}

			encoded.Set(x, y, image.RGBA{r8, g8, b8, a8})
		}
	}

	return encoded
}

// Decode extracts a hidden message from an image
func decodeMessage(img image.Image) string {
	bounds := img.Bounds()

	// First, extract the length (8 bytes = 64 bits)
	lengthBytes := make([]byte, 8)
	bitIndex := 0

	for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < 64; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && bitIndex < 64; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Extract LSB from each channel
			if bitIndex < 64 {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := uint8((r >> 8) & 1)
				lengthBytes[byteIndex] |= bit << (7 - bitOffset)
				bitIndex++
			}

			if bitIndex < 64 {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := uint8((g >> 8) & 1)
				lengthBytes[byteIndex] |= bit << (7 - bitOffset)
				bitIndex++
			}

			if bitIndex < 64 {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := uint8((b >> 8) & 1)
				lengthBytes[byteIndex] |= bit << (7 - bitOffset)
				bitIndex++
			}
		}
	}

	// Parse the length
	var messageLen int
	fmt.Sscanf(string(lengthBytes), "%08d", &messageLen)

	if messageLen <= 0 || messageLen > 1000000 {
		return ""
	}

	// Extract the message
	messageBytes := make([]byte, messageLen)
	totalBits := messageLen * 8
	bitIndex = 0

	for y := bounds.Min.Y; y < bounds.Max.Y && bitIndex < totalBits; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && bitIndex < totalBits; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Skip the first 64 bits (length header)
			currentBit := bitIndex + 64

			// Calculate pixel position for current bit
			pixelBit := currentBit % 3
			pixelIndex := currentBit / 3

			if pixelIndex != (x-bounds.Min.X)+(y-bounds.Min.Y)*(bounds.Max.X-bounds.Min.X) {
				continue
			}

			if pixelBit == 0 && bitIndex < totalBits {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := uint8((r >> 8) & 1)
				messageBytes[byteIndex] |= bit << (7 - bitOffset)
				bitIndex++
			}

			if pixelBit <= 1 && bitIndex < totalBits {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := uint8((g >> 8) & 1)
				messageBytes[byteIndex] |= bit << (7 - bitOffset)
				bitIndex++
			}

			if pixelBit <= 2 && bitIndex < totalBits {
				byteIndex := bitIndex / 8
				bitOffset := bitIndex % 8
				bit := uint8((b >> 8) & 1)
				messageBytes[byteIndex] |= bit << (7 - bitOffset)
				bitIndex++
			}
		}
	}

	return string(messageBytes)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleEncode(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to parse form"})
		return
	}

	// Get the image file
	file, _, err := r.FormFile("image")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No image file provided"})
		return
	}
	defer file.Close()

	// Get the message
	message := r.FormValue("message")
	if message == "" {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No message provided"})
		return
	}

	// Get optional password
	password := r.FormValue("password")

	// Encrypt message if password provided
	encryptedMessage, err := encryptMessage(message, password)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to encrypt message"})
		return
	}

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to decode image"})
		return
	}

	// Encode encrypted message into image
	encodedImg := encodeMessage(img, string(encryptedMessage))

	// Send the encoded image back
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=encoded.png")

	err = png.Encode(w, encodedImg)
	if err != nil {
		log.Printf("Failed to encode PNG: %v", err)
	}
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to parse form"})
		return
	}

	// Get the image file
	file, _, err := r.FormFile("image")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No image file provided"})
		return
	}
	defer file.Close()

	// Get optional password
	password := r.FormValue("password")

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to decode image"})
		return
	}

	// Decode message from image
	encryptedMessage := decodeMessage(img)

	// Decrypt message if password provided
	message, err := decryptMessage([]byte(encryptedMessage), password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DecodeResponse{Message: message})
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
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

func main() {
	// API endpoints
	http.HandleFunc("/api/encode", handleEncode)
	http.HandleFunc("/api/decode", handleDecode)

	// Serve frontend
	http.HandleFunc("/", serveStatic)

	port := ":8080"
	fmt.Printf("🚀 Steganography server running on http://localhost%s\n", port)
	fmt.Println("📁 Frontend available at: http://localhost:8080")
	fmt.Println("🔧 API endpoints:")
	fmt.Println("   - POST /api/encode (encode message into image)")
	fmt.Println("   - POST /api/decode (decode message from image)")

	log.Fatal(http.ListenAndServe(port, nil))
}

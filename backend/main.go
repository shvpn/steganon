package main

import (
	"fmt"
	"log"
	"net/http"
	"steganography/handlers"
)

func main() {
	// API endpoints
	http.HandleFunc("/api/encode", handlers.HandleEncode)
	http.HandleFunc("/api/decode", handlers.HandleDecode)

	// Serve frontend
	http.HandleFunc("/", handlers.ServeStatic)

	port := ":8080"
	printServerInfo(port)

	log.Fatal(http.ListenAndServe(port, nil))
}

// printServerInfo displays server startup information
func printServerInfo(port string) {
	fmt.Printf("🚀 Steganography server running on http://localhost%s\n", port)
	fmt.Println("📁 Frontend available at: http://localhost:8080")
	fmt.Println("🔧 API endpoints:")
	fmt.Println("   - POST /api/encode (encode message into image)")
	fmt.Println("   - POST /api/decode (decode message from image)")
	fmt.Println("\n✨ Features:")
	fmt.Println("   - LSB Steganography")
	fmt.Println("   - AES-256-GCM Encryption")
	fmt.Println("   - SHA-256 Password Hashing")
}

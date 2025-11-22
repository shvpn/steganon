# Backend Structure

## Overview

The backend is organized into clean, modular packages for maintainability and scalability.

## Package Organization

### `main.go`
- Application entry point
- HTTP server configuration
- Route registration
- Server startup information

### `handlers/`
API and request handlers

#### `api.go`
- **HandleEncode**: Processes image encoding requests
- **HandleDecode**: Processes image decoding requests
- **enableCORS**: CORS middleware
- **sendJSONError**: Error response helper

#### `static.go`
- **ServeStatic**: Serves frontend files

### `utils/`
Utility functions and core logic

#### `crypto.go`
Cryptographic operations:
- **EncryptMessage**: AES-256-GCM encryption with SHA-256 hashed password
- **DecryptMessage**: AES-256-GCM decryption with password verification

#### `steganography.go`
LSB steganography implementation:
- **EncodeMessageInImage**: Hides message in image pixels
- **DecodeMessageFromImage**: Extracts message from image pixels
- Helper functions for bit manipulation and data extraction

## Code Quality

### Clean Code Principles
- ✅ Single Responsibility: Each function has one clear purpose
- ✅ Descriptive Names: Functions and variables clearly indicate their purpose
- ✅ Small Functions: Complex logic broken into smaller, testable units
- ✅ DRY (Don't Repeat Yourself): Common logic extracted to helper functions
- ✅ Error Handling: Proper error propagation and handling
- ✅ Comments: Clear documentation for complex operations

### Design Patterns
- **Separation of Concerns**: Handlers, utilities, and main logic separated
- **Modular Architecture**: Easy to extend and maintain
- **Dependency Injection**: Functions accept dependencies as parameters

## Building & Running

```bash
# Navigate to backend
cd backend

# Build the application
go build -o stegano

# Run the application
./stegano
# or
go run main.go
```

## Testing

```bash
# Run tests (when implemented)
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./utils
go test ./handlers
```

## Adding New Features

### Adding a New API Endpoint

1. Add handler function in `handlers/api.go`:
```go
func HandleNewFeature(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

2. Register route in `main.go`:
```go
http.HandleFunc("/api/new-feature", handlers.HandleNewFeature)
```

### Adding New Utility Functions

1. Create or update file in `utils/` package
2. Export function with capital first letter
3. Import in handlers: `import "steganography/utils"`

## Performance Considerations

- Image processing uses efficient byte manipulation
- Minimal memory allocations
- Concurrent request handling via Go's goroutines
- No global state for thread safety

## Security

- Password hashing with SHA-256
- AES-256-GCM for authenticated encryption
- Input validation on all endpoints
- No data persistence (stateless)
- CORS configuration for production use

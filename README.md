# Steganography Tool 🔐

A modern web-based steganography tool that allows you to hide secret messages inside images using LSB (Least Significant Bit) encoding. Built with HTML, CSS, JavaScript frontend and Go backend.

## Features

- **Encode Messages**: Hide secret text messages inside PNG images
- **Decode Messages**: Extract hidden messages from encoded images
- **LSB Steganography**: Uses Least Significant Bit technique for invisible message embedding
- **User-Friendly Interface**: Modern, responsive dark-themed UI
- **Real-time Preview**: See your images before encoding/decoding
- **Download Encoded Images**: Save your steganographic images locally
- **Copy to Clipboard**: Easily copy decoded messages

## Technology Stack

### Frontend
- HTML5
- CSS3 (Modern gradient design with dark theme)
- Vanilla JavaScript (No frameworks required)

### Backend
- Go (Golang)
- Standard library (no external dependencies)
- PNG image processing

## Project Structure

```
Stegano/
├── backend/
│   ├── main.go          # Go server with encode/decode endpoints
│   └── go.mod           # Go module file
└── frontend/
    ├── index.html       # Main HTML file
    ├── style.css        # Styling
    └── script.js        # Frontend logic
```

## Installation & Setup

### Prerequisites

- Go 1.21 or higher installed ([Download Go](https://golang.org/dl/))

### Steps

1. **Navigate to the backend directory**:
   ```bash
   cd backend
   ```

2. **Run the Go server**:
   ```bash
   go run main.go
   ```

3. **Access the application**:
   Open your browser and visit: `http://localhost:8080`

## Usage

### Encoding a Message

1. Click on the **Encode** tab
2. Select an image file (PNG recommended for lossless encoding)
3. Enter your secret message in the text area
4. Click **Encode Message**
5. Download the encoded image once processing is complete

### Decoding a Message

1. Click on the **Decode** tab
2. Select an encoded image (previously created with this tool)
3. Click **Decode Message**
4. The hidden message will be displayed
5. Click **Copy to Clipboard** to copy the message

## How It Works

### LSB Steganography with Optional Encryption

The tool combines **Least Significant Bit (LSB)** steganography with **AES encryption**:

1. **Encoding**: 
   - (Optional) Encrypts the message using AES-GCM with SHA-256 hashed password
   - Converts the (encrypted) message to binary
   - Replaces the least significant bits of RGB pixel values with message bits
   - Creates a new image that looks identical to the original
   - Stores message length as a header for accurate extraction

2. **Decoding**:
   - Reads the LSB of each pixel's RGB values
   - Extracts the message length from the header
   - Reconstructs the original (encrypted) message from the bits
   - (Optional) Decrypts the message using the provided password

### Encryption Details

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Derivation**: SHA-256 hash of password
- **Authentication**: GCM provides both encryption and authentication
- **Nonce**: Randomly generated for each encryption

## API Endpoints

### POST `/api/encode`
- **Description**: Encode a message into an image
- **Parameters**: 
  - `image` (file): The carrier image
  - `message` (string): The secret message to hide
  - `password` (string, optional): Password for encryption
- **Response**: PNG image file with encoded message

### POST `/api/decode`
- **Description**: Extract a hidden message from an image
- **Parameters**: 
  - `image` (file): The encoded image
  - `password` (string, optional): Password for decryption (if encrypted)
- **Response**: JSON object with the decoded message

## Important Notes

- Use PNG images for best results (JPG compression may destroy hidden data)
- The image must be large enough to hold your message
- Encoded images look identical to originals to the human eye
- Only images encoded with this tool can be decoded by it

## Security Considerations

⚠️ **This tool is for educational purposes**:
- LSB steganography is detectable with steganalysis tools
- Password protection adds encryption layer (AES-256-GCM)
- Password is hashed with SHA-256 before use as encryption key
- Without password, messages are hidden but not encrypted
- Always use strong passwords for sensitive data
- The same password must be used for both encoding and decoding

## License

This project is open source and available for educational purposes.

## Contributing

Feel free to fork, modify, and submit pull requests!

## Credits

Built with ❤️ for cryptography enthusiasts and students.

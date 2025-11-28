# III. Implementation Details

This section details the core algorithms and code structure used in the project.

## Technology Stack
-   **Language**: Go (Golang) 1.21+
-   **Standard Libraries**: `image`, `image/png`, `crypto/*`, `net/http`
-   **Frontend**: HTML5, CSS3, JavaScript (ES6+)

## 1. Cryptographic Implementation (`utils/crypto.go`)

The project uses the Go standard library's `crypto` packages.

### Key Derivation
Since AES requires a fixed-size key (32 bytes for AES-256), user passwords are hashed.
```go
// Hash password with SHA-256 to get 32-byte key
hash := sha256.Sum256([]byte(password))
key := hash[:]
```

### Encryption (AES-GCM)
We use Galois/Counter Mode (GCM) because it provides authenticated encryption.
1.  **Nonce Generation**: A unique 12-byte nonce is generated for every encryption using `crypto/rand`.
2.  **Sealing**: The `gcm.Seal` function encrypts the data and appends an authentication tag.
3.  **Packing**: The final payload structure is `[Nonce (12 bytes)] + [Ciphertext] + [Auth Tag]`.

### Decryption
1.  **Unpacking**: The nonce is extracted from the first 12 bytes of the data.
2.  **Opening**: `gcm.Open` decrypts the ciphertext and verifies the authentication tag. If the password is wrong or data corrupted, it returns an error.

## 2. Steganography Implementation (`utils/steganography.go`)

The core logic relies on manipulating the **Least Significant Bit (LSB)** of image pixels.

### Data Preparation
Before embedding, the message is prefixed with a length header to ensure the decoder knows exactly how much data to read.
-   **Header**: 8 bytes (string representation of length, e.g., "00000123").
-   **Payload**: The actual message bytes (or encrypted bytes).

### Encoding Algorithm (`EncodeMessageInImage`)
1.  **Image Copy**: A mutable copy (`image.RGBA`) of the source image is created.
2.  **Iteration**: The code iterates through every pixel (x, y).
3.  **Bit Insertion**:
    -   The algorithm processes the Red, Green, and Blue channels sequentially.
    -   For each channel, the LSB is cleared (`val & 0xFE`) and replaced with a bit from the data (`val | bit`).
    -   This allows storing 3 bits per pixel.
4.  **Termination**: The loop stops once all data bits are embedded.

```go
// Example bit insertion logic
bit := (data[byteIndex] >> (7 - bitOffset)) & 1
channel = (channel & 0xFE) | bit
```

### Decoding Algorithm (`DecodeMessageFromImage`)
1.  **Header Extraction**: The first 64 bits (8 bytes * 8 bits) are read from the image pixels to reconstruct the length header.
2.  **Length Parsing**: The header string is parsed to an integer (e.g., "00000010" -> 10 bytes).
3.  **Payload Extraction**: The loop continues reading LSBs until the specified number of bytes is retrieved.

## 3. API Handlers (`handlers/api.go`)

The handlers act as controllers, bridging the HTTP layer with the logic layer.

-   **Multipart Parsing**: Handles file uploads up to 10MB.
-   **Error Handling**: Returns JSON error responses (`{"error": "message"}`) for invalid inputs or processing failures.
-   **CORS**: Middleware is implemented to allow cross-origin requests during development or specific domains in production.

## 4. Frontend Implementation

-   **File API**: Uses `FileReader` to show image previews before upload.
-   **Fetch API**: Sends asynchronous `POST` requests to the backend.
-   **Blob Handling**: When receiving the encoded image, the frontend converts the binary response into a Blob URL for downloading.

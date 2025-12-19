# Steganography Tool By Ork Senghout Group 2 (Individual) 
## Note: deploy and update script are for automate deployment.However, this might not work on some server. So we can deploy it manually throught this README

A modern web-based steganography tool that allows you to hide secret messages inside images using LSB (Least Significant Bit) encoding. Built with HTML, CSS, JavaScript frontend and Go backend.

## Features

### Core Features
- **Encode Messages**: Hide secret text messages inside PNG images
- **Decode Messages**: Extract hidden messages from encoded images
- **LSB Steganography**: Uses Least Significant Bit technique for invisible message embedding
- **Password Protection**: Optional AES-256-GCM encryption with SHA-256 hashed passwords
- **Real-time Preview**: See your images before encoding/decoding
- **Download Encoded Images**: Save your steganographic images locally
- **Copy to Clipboard**: Easily copy decoded messages

### Performance & Responsiveness
- **Fast Processing**: Efficient Go backend handles images quickly
- **Responsive Design**: Mobile-first design works on all devices
- **Minimal Dependencies**: No external libraries required, reducing overhead
- **Concurrent Handling**: Go's goroutines enable multiple simultaneous requests
- **Optimized Frontend**: Vanilla JavaScript for fast load times

### Security Features
- **AES-256-GCM Encryption**: Military-grade encryption for message protection
- **SHA-256 Password Hashing**: Secure key derivation from passwords
- **CORS Protection**: Configurable cross-origin resource sharing
- **No Data Storage**: Messages are never stored on the server
- **Client-Side Validation**: Input validation before API calls
- **HTTPS Ready**: Designed for secure SSL/TLS deployment

## Technology Stack

### Frontend
- HTML5
- CSS
- JavaScript

### Backend
- Go (Golang)
- Standard library
- PNG image processing

## Project Structure

```
Stegano/
├── backend/
│   ├── main.go                    # Application entry point
│   ├── go.mod                     # Go module file
│   ├── handlers/
│   │   ├── api.go                 # API handlers (encode/decode)
│   │   └── static.go              # Static file serving
│   └── utils/
│       ├── crypto.go              # AES encryption/decryption
│       └── steganography.go       # LSB steganography logic
├── frontend/
│   ├── index.html                 # Main HTML file
│   ├── style.css                  # Styling
│   └── script.js                  # Frontend logic
├── deploy.sh                      # Automated deployment script
├── update.sh                      # Quick update script
├── DEPLOYMENT.md                  # Deployment guide
└── README.md                      # Project documentation

```

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

## Installation & Setup -My server is Ubuntu2.24


## Production Deployment

### Prerequisites
- Ubuntu/Debian server (tested on Ubuntu 22.04)
- Domain name pointed to your server
- Root or sudo access

### Step 1: Server Setup

```bash
# Update system packages
sudo apt update

# Install Go
sudo apt install go-lang -y

# Verify Go installation
go version
```

### Step 2: Deploy Application

```bash
# Create application directory
sudo mkdir -p /var/www/stegano
cd /var/www/stegano

# Clone or copy your project files
# Upload your Stegano folder to /var/www/stegano

# Build the Go application
cd /var/www/stegano/backend
go build -o stegano main.go
```

### Step 3: Create Systemd Service

Create a service file to run your application as a daemon:

```bash
sudo nano /etc/systemd/system/stegano.service
```

Add the following content:

```ini
[Unit]
Description=Steganography Tool Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/stegano/backend
ExecStart=/var/www/stegano/backend/stegano
Restart=always
RestartSec=5
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable stegano
sudo systemctl start stegano
sudo systemctl status stegano
```

### Step 4: Install Nginx

```bash
# Install Nginx
sudo apt install nginx -y

# Create Nginx configuration
sudo nano /etc/nginx/sites-available/stegano
```

Add the following configuration:

```nginx
server {
    listen 80;
    server_name stegano.shvpn.live;

    # Increase max upload size for images
    client_max_body_size 20M;

    # Serve frontend files
    location / {
        root /var/www/stegano/frontend;
        try_files $uri $uri/ /index.html;
    }

    # Proxy API requests to Go backend
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable the site:

```bash
# Create symbolic link
sudo ln -s /etc/nginx/sites-available/stegano /etc/nginx/sites-enabled/

# Test Nginx configuration
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx
```

### Step 5: Configure DNS

Point your domain to your server:

1. Go to your DNS provider (e.g., Cloudflare, GoDaddy, Namecheap)
2. Add an **A Record**:
   - **Type**: A
   - **Name**: stegano (or @ for root domain)
   - **Value**: {YourIP}
   - **TTL**: Automatic or 300

Wait for DNS propagation (usually 5-15 minutes).

### Step 6: Install SSL Certificate with Certbot

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain SSL certificate
sudo certbot --nginx -d stegano.shvpn.live

# Follow the prompts:
# - Enter your email address
# - Agree to terms of service
# - Choose whether to redirect HTTP to HTTPS (recommended: Yes)
```

Certbot will automatically:
- Obtain the SSL certificate
- Configure Nginx to use HTTPS
- Set up auto-renewal

Test auto-renewal:

```bash
sudo certbot renew --dry-run
```

### Step 7: Configure Firewall (You can Just Diable firewall)

```bash
# Allow Nginx through firewall
sudo ufw allow 'Nginx Full'
sudo ufw allow OpenSSH
sudo ufw enable
sudo ufw status
```

### Step 8: Update Go Backend for Production(Optinal- Security consideration)

Update the CORS settings in `backend/main.go` for production:

```go
func enableCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "https://stegano.shvpn.live")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
```

Rebuild and restart:

```bash
cd /var/www/stegano/backend
go build -o stegano main.go
sudo systemctl restart stegano
```


### Verification

Visit **https://stegano.shvpn.live** in your browser. You should see:
- ✅ Secure HTTPS connection (padlock icon)
- ✅ Your steganography tool interface
- ✅ Ability to encode/decode images

### Monitoring & Maintenance

```bash
# View application logs
sudo journalctl -u stegano -f

# Check Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# Restart services if needed
sudo systemctl restart stegano
sudo systemctl restart nginx

# Check SSL certificate expiration
sudo certbot certificates
```

### Performance Optimization

Add caching headers in Nginx:

```nginx
# Add inside server block
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    root /var/www/stegano/frontend;
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

Enable Gzip compression:

```nginx
# Add inside server block
gzip on;
gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
gzip_min_length 1000;
```

### Troubleshooting

**Service won't start:**
```bash
sudo journalctl -u stegano -n 50
```

**Nginx errors:**
```bash
sudo nginx -t
sudo tail -f /var/log/nginx/error.log
```

**SSL certificate issues:**
```bash
sudo certbot certificates
sudo certbot renew --force-renewal
```

**Port already in use:**
```bash
sudo lsof -i :8080
sudo kill -9 <PID>
```

## License

This project is open source and available for educational purposes.

## Contributing

Feel free to fork, modify, and submit pull requests!

## Credits

Lect.Mr. Meas Sothearath
AI

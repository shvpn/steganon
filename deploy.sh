#!/bin/bash

# Deployment script for Steganography Tool
# Usage: ./deploy.sh

set -e

echo "🚀 Starting deployment of Steganography Tool..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_DIR="/var/www/stegano"
DOMAIN="stegano.shvpn.live"
SERVER_IP="54.169.190.81"

echo -e "${YELLOW}Step 1: Installing dependencies...${NC}"
sudo apt update
sudo apt install -y nginx certbot python3-certbot-nginx

echo -e "${GREEN}✓ Dependencies installed${NC}"

echo -e "${YELLOW}Step 2: Creating application directory...${NC}"
sudo mkdir -p $APP_DIR
sudo chown -R $USER:$USER $APP_DIR

echo -e "${GREEN}✓ Directory created${NC}"

echo -e "${YELLOW}Step 3: Building Go application...${NC}"
cd $APP_DIR/backend
go build -o stegano main.go

echo -e "${GREEN}✓ Application built${NC}"

echo -e "${YELLOW}Step 4: Creating systemd service...${NC}"
sudo tee /etc/systemd/system/stegano.service > /dev/null <<EOF
[Unit]
Description=Steganography Tool Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=$APP_DIR/backend
ExecStart=$APP_DIR/backend/stegano
Restart=always
RestartSec=5
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable stegano
sudo systemctl start stegano

echo -e "${GREEN}✓ Service created and started${NC}"

echo -e "${YELLOW}Step 5: Configuring Nginx...${NC}"
sudo tee /etc/nginx/sites-available/stegano > /dev/null <<EOF
server {
    listen 80;
    server_name $DOMAIN;

    # Increase max upload size for images
    client_max_body_size 20M;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # Serve frontend files
    location / {
        root $APP_DIR/frontend;
        try_files \$uri \$uri/ /index.html;
        
        # Caching for static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # Proxy API requests to Go backend
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # Timeout settings for large file uploads
        proxy_connect_timeout 300;
        proxy_send_timeout 300;
        proxy_read_timeout 300;
    }

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
    gzip_min_length 1000;
}
EOF

sudo ln -sf /etc/nginx/sites-available/stegano /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

echo -e "${GREEN}✓ Nginx configured${NC}"

echo -e "${YELLOW}Step 6: Configuring firewall...${NC}"
sudo ufw allow 'Nginx Full'
sudo ufw allow OpenSSH
echo "y" | sudo ufw enable

echo -e "${GREEN}✓ Firewall configured${NC}"

echo -e "${YELLOW}Step 7: Obtaining SSL certificate...${NC}"
echo -e "${YELLOW}Please enter your email for SSL certificate notifications:${NC}"
read -p "Email: " EMAIL

sudo certbot --nginx -d $DOMAIN --non-interactive --agree-tos --email $EMAIL --redirect

echo -e "${GREEN}✓ SSL certificate obtained${NC}"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}🎉 Deployment complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "Your application is now available at:"
echo -e "${GREEN}https://$DOMAIN${NC}"
echo ""
echo "Useful commands:"
echo "  - View logs: sudo journalctl -u stegano -f"
echo "  - Restart app: sudo systemctl restart stegano"
echo "  - Check status: sudo systemctl status stegano"
echo "  - Nginx logs: sudo tail -f /var/log/nginx/error.log"
echo ""

#!/bin/bash

# Quick update script for Steganography Tool
# Usage: ./update.sh

set -e

APP_DIR="/var/www/stegano"

echo "🔄 Updating Steganography Tool..."

# Build new version
cd $APP_DIR/backend
go build -o stegano main.go

# Restart service
sudo systemctl restart stegano

# Check status
if sudo systemctl is-active --quiet stegano; then
    echo "✅ Update successful! Service is running."
    sudo systemctl status stegano --no-pager
else
    echo "❌ Update failed! Service is not running."
    sudo journalctl -u stegano -n 20
    exit 1
fi

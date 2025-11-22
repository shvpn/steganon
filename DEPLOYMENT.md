# Steganography Tool - Deployment Guide

## This deployment is not really great.Howerver, it is best for low budget.
**In this project i will use my own domain that point from cloudeflare while using nginx to forward proxy on port 8080 and using cert bot for ssl certificate**

### Option 1: Automated Deployment

1. **Upload files to your server:**
```bash
scp -r Stegano/ user@<ServerIP>:/tmp/
```

2. **SSH into your server:**
```bash
ssh user@<ServerIP>
```

3. **Move files and run deployment script:**
```bash
sudo mv /tmp/Stegano /var/www/stegano
cd /var/www/stegano
chmod +x deploy.sh
sudo ./deploy.sh
```

The script will automatically:
- Install all dependencies
- Build the Go application
- Create and start the systemd service
- Configure Nginx
- Set up SSL with Certbot
- Configure firewall

### Option 2: Manual Deployment

Follow the detailed steps in the main README.md file under "Production Deployment" section.

## DNS Configuration

Before deployment, ensure your DNS is configured:

1. Log into your DNS provider
2. Create an A record:
   - **Type**: A
   - **Name**: stegano
   - **Value**: ServerIP
   - **TTL**: 300 (or Auto)

3. Wait for propagation (5-15 minutes)
4. Verify: `nslookup stegano.shvpn.live`

## Post-Deployment Checklist

- [ ] Application accessible at https://stegano.shvpn.live
- [ ] SSL certificate is active (padlock icon)
- [ ] Image encoding works
- [ ] Image decoding works
- [ ] Password protection works
- [ ] File download works
- [ ] Mobile responsive design works

## Monitoring

### Check Application Status
```bash
sudo systemctl status stegano
```

### View Application Logs
```bash
sudo journalctl -u stegano -f
```

### View Nginx Logs
```bash
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Check SSL Certificate
```bash
sudo certbot certificates
```

## Updating the Application

After making changes to your code:

```bash
cd /var/www/stegano
chmod +x update.sh
./update.sh
```

Or manually:
```bash
cd /var/www/stegano/backend
go build -o stegano main.go
sudo systemctl restart stegano
```

## Performance Tuning

### Increase Upload Size (if needed)
Edit `/etc/nginx/sites-available/stegano`:
```nginx
client_max_body_size 50M;  # Increase from 20M
```

### Enable HTTP/2
Nginx with SSL automatically uses HTTP/2 for better performance.

### Monitor Resource Usage
```bash
htop
df -h
free -m
```

## Troubleshooting

### Service Won't Start
```bash
sudo journalctl -u stegano -n 50 --no-pager
```

### Port 8080 Already in Use
```bash
sudo lsof -i :8080
sudo kill -9 <PID>
```

### SSL Certificate Issues
```bash
sudo certbot renew --force-renewal
sudo systemctl restart nginx
```

### Application Not Accessible
1. Check firewall: `sudo ufw status`
2. Check Nginx: `sudo nginx -t`
3. Check service: `sudo systemctl status stegano`
4. Check DNS: `nslookup stegano.shvpn.live`

## Security Best Practices

1. **Keep system updated:**
```bash
sudo apt update && sudo apt upgrade -y
```

2. **Monitor failed login attempts:**
```bash
sudo tail -f /var/log/auth.log
```

3. **Enable automatic security updates:**
```bash
sudo apt install unattended-upgrades -y
sudo dpkg-reconfigure -plow unattended-upgrades
```

4. **Backup regularly:**
```bash
# Backup application
tar -czf stegano-backup-$(date +%Y%m%d).tar.gz /var/www/stegano

# Backup Nginx config
sudo cp /etc/nginx/sites-available/stegano /etc/nginx/sites-available/stegano.backup
```

## Support

For issues or questions:
- Check logs first
- Review the main README.md
- Verify DNS and SSL configuration
- Ensure all services are running

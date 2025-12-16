# Nginx Configuration

This directory contains the Nginx reverse proxy configuration for Suipic.

## Structure

```
nginx/
├── nginx.conf          # Main Nginx configuration
├── conf.d/
│   └── default.conf   # Site configuration with routing
├── ssl/
│   ├── README.md      # SSL certificate guide
│   ├── cert.pem       # SSL certificate (not in repo)
│   └── key.pem        # SSL private key (not in repo)
└── README.md          # This file
```

## Configuration Files

### nginx.conf

Main Nginx configuration with:
- Worker processes and connections
- Logging configuration
- GZIP compression
- Client body size limit (100MB for photo uploads)
- MIME types
- Performance optimizations

### conf.d/default.conf

Site-specific configuration with:
- Upstream backend and frontend definitions
- HTTP server block (port 80)
- HTTPS server block (port 443, commented out)
- Reverse proxy settings
- Headers for proper forwarding

## Routing

Nginx routes requests as follows:

```
http://localhost/api/* → backend:3000
http://localhost/*     → frontend:3001
```

### API Routes

All requests starting with `/api` are proxied to the backend service:

```nginx
location /api {
    proxy_pass http://backend;
    # ... proxy headers
}
```

### Frontend Routes

All other requests are proxied to the frontend service:

```nginx
location / {
    proxy_pass http://frontend;
    # ... proxy headers
}
```

## Enabling HTTPS

To enable HTTPS support:

1. **Generate or obtain SSL certificates** (see `ssl/README.md`)

2. **Place certificates in the ssl/ directory:**
   ```bash
   docker/nginx/ssl/cert.pem
   docker/nginx/ssl/key.pem
   ```

3. **Edit conf.d/default.conf:**
   - Uncomment the HTTPS server block (lines starting with #)
   - Uncomment the HTTP to HTTPS redirect
   - Update `server_name` to your domain

4. **Restart Nginx:**
   ```bash
   docker-compose restart nginx
   ```

### Example HTTPS Configuration

After uncommenting and configuring:

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # ... rest of configuration
}
```

## Custom Configuration

### Changing Upload Size Limit

Edit `nginx.conf`:

```nginx
http {
    client_max_body_size 500M;  # Change from 100M to 500M
}
```

### Adding Rate Limiting

Add to `conf.d/default.conf`:

```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

server {
    location /api {
        limit_req zone=api burst=20;
        proxy_pass http://backend;
        # ...
    }
}
```

### Adding Custom Headers

Add to server block in `conf.d/default.conf`:

```nginx
server {
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

### Caching Static Assets

Add to `conf.d/default.conf`:

```nginx
location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf)$ {
    proxy_pass http://frontend;
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## Proxy Headers

The configuration sets these headers for proper request handling:

- `X-Real-IP`: Client's real IP address
- `X-Forwarded-For`: Chain of proxy IPs
- `X-Forwarded-Proto`: Original protocol (http/https)
- `Host`: Original host header
- `Upgrade` & `Connection`: For WebSocket support

## Performance Tuning

### Worker Processes

Edit `nginx.conf`:

```nginx
worker_processes auto;  # Uses number of CPU cores
```

### Connections

```nginx
events {
    worker_connections 2048;  # Increase from 1024
}
```

### Timeouts

Add to `conf.d/default.conf`:

```nginx
proxy_connect_timeout 75s;
proxy_send_timeout 300s;
proxy_read_timeout 300s;
```

## Security Headers

Add security headers in `conf.d/default.conf`:

```nginx
server {
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
}
```

## WebSocket Support

The configuration includes WebSocket support via these headers:

```nginx
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection 'upgrade';
proxy_cache_bypass $http_upgrade;
```

## Logging

### Access Logs

Located at: `/var/log/nginx/access.log` (inside container)

View with:
```bash
docker-compose logs nginx
docker-compose exec nginx tail -f /var/log/nginx/access.log
```

### Error Logs

Located at: `/var/log/nginx/error.log` (inside container)

View with:
```bash
docker-compose exec nginx tail -f /var/log/nginx/error.log
```

### Custom Log Format

Edit `nginx.conf` to customize:

```nginx
log_format custom '$remote_addr - $remote_user [$time_local] '
                  '"$request" $status $body_bytes_sent '
                  '"$http_referer" "$http_user_agent" '
                  'rt=$request_time';

access_log /var/log/nginx/access.log custom;
```

## Testing Configuration

Test configuration without restarting:

```bash
docker-compose exec nginx nginx -t
```

Reload configuration:

```bash
docker-compose exec nginx nginx -s reload
```

## Troubleshooting

### 502 Bad Gateway

Check if backend/frontend services are running:
```bash
docker-compose ps
docker-compose logs backend
docker-compose logs frontend
```

### 413 Request Entity Too Large

Increase `client_max_body_size` in `nginx.conf`

### SSL Certificate Errors

Verify certificate files exist and are valid:
```bash
ls -la docker/nginx/ssl/
openssl x509 -in docker/nginx/ssl/cert.pem -text -noout
```

### Connection Timeouts

Increase timeout values in `conf.d/default.conf`:
```nginx
proxy_connect_timeout 90s;
proxy_read_timeout 600s;
```

## Advanced Topics

### Load Balancing

For multiple backend instances:

```nginx
upstream backend {
    least_conn;
    server backend1:3000;
    server backend2:3000;
    server backend3:3000;
}
```

### Health Checks

```nginx
upstream backend {
    server backend:3000 max_fails=3 fail_timeout=30s;
}
```

### IP Whitelisting

Restrict access to specific IPs:

```nginx
location /admin {
    allow 192.168.1.0/24;
    deny all;
    proxy_pass http://backend;
}
```

## Further Reading

- [Nginx Documentation](https://nginx.org/en/docs/)
- [Nginx Reverse Proxy Guide](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/)
- [SSL Configuration Guide](https://ssl-config.mozilla.org/)

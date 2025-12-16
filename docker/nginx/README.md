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
- Worker processes and connections (auto-tuned)
- Event-driven architecture with epoll
- Logging configuration (main and detailed formats)
- GZIP compression (30+ MIME types, level 6)
- Client body size limit (100MB for photo uploads)
- Buffer optimizations
- **Rate limiting zones** (3 zones for different traffic types)
- Security headers (X-Frame-Options, X-Content-Type-Options, etc.)
- Performance tuning (keepalive, timeouts, etc.)

### conf.d/default.conf

Site-specific configuration with:
- **Upstream definitions** with health checks and keepalive
- HTTP server block (port 80)
- **HTTPS server block** (port 443) - pre-configured
- **Rate limiting** applied to routes
- Reverse proxy settings with optimized buffering
- Headers for proper forwarding
- **Static asset caching** (1 year for images, fonts, etc.)
- Health check endpoint (`/health`)

## Features

### Rate Limiting

The configuration implements three rate limiting zones:

| Zone | Limit | Burst | Applied To | Purpose |
|------|-------|-------|------------|---------|
| `login_limit` | 5 req/min | 2 | `/api/auth/login` | Prevent brute force attacks |
| `api_limit` | 10 req/s | 20 | `/api/*` | Protect API endpoints |
| `general_limit` | 30 req/s | 50 | `/*` | Protect frontend |

Additionally, connection limiting restricts each IP to 10 concurrent connections.

**To adjust rate limits**, edit `nginx.conf`:
```nginx
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
```

Then update the corresponding `limit_req` directives in `conf.d/default.conf`:
```nginx
limit_req zone=api_limit burst=20 nodelay;
```

### Gzip Compression

Enabled for optimal bandwidth usage:
- Compression level: 6 (balanced performance/ratio)
- Minimum file size: 256 bytes
- 30+ MIME types including:
  - Text: HTML, CSS, JS, JSON, XML, plain text
  - Fonts: TTF, OTF, EOT, WOFF
  - Images: SVG, ICO
  - Data: JSON, XML, CSV

### SSL/TLS Support

Pre-configured HTTPS server with:
- **TLS 1.2 and 1.3** support
- Modern cipher suites (ECDHE, ChaCha20, AES-GCM)
- **OCSP stapling** for certificate validation
- **HSTS** with preload support
- Session caching (50MB, 1 day timeout)
- HTTP/2 enabled

**Certificates are automatically loaded** from `/etc/nginx/ssl/` when present:
- `cert.pem` - SSL certificate (public key)
- `key.pem` - Private key
- `chain.pem` - Certificate chain (optional)

### Security Headers

Applied to all responses:
- `X-Frame-Options: SAMEORIGIN` - Prevent clickjacking
- `X-Content-Type-Options: nosniff` - Prevent MIME sniffing
- `X-XSS-Protection: 1; mode=block` - XSS protection
- `Referrer-Policy: no-referrer-when-downgrade` - Privacy
- `Strict-Transport-Security` (HTTPS only) - Force HTTPS

### Performance Optimizations

- **HTTP/2** enabled on HTTPS
- **Keepalive connections** to upstream services (pool of 32)
- **Static asset caching** with 1-year expiration
- **Connection pooling** to backends
- **Optimized buffers** for efficient memory usage
- **TCP optimizations** (tcp_nopush, tcp_nodelay)

## Routing

Nginx routes requests as follows:

```
http://localhost/health       → nginx health check
http://localhost/api/* → backend:3000
http://localhost/*     → frontend:3001
```

### API Routes

All requests starting with `/api` are proxied to the backend service with rate limiting:

```nginx
location /api {
    limit_req zone=api_limit burst=20 nodelay;
    limit_conn addr 10;
    proxy_pass http://backend;
    # ... proxy headers
}
```

### Login Route

Special rate limiting for login endpoint:

```nginx
location ~ ^/api/auth/login$ {
    limit_req zone=login_limit burst=2 nodelay;
    limit_conn addr 5;
    proxy_pass http://backend;
    # ... proxy headers
}
```

### Frontend Routes

All other requests are proxied to the frontend service:

```nginx
location / {
    limit_req zone=general_limit burst=50 nodelay;
    limit_conn addr 10;
    proxy_pass http://frontend;
    # ... proxy headers with WebSocket support
}
```

### Static Assets

Optimized caching for static files:

```nginx
location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
    # ... proxy to frontend
}
```

## Upstream Definitions

Two upstream backends with health checks:

```nginx
upstream backend {
    server backend:3000 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

upstream frontend {
    server frontend:3001 max_fails=3 fail_timeout=30s;
    keepalive 32;
}
```

## Enabling HTTPS

HTTPS is pre-configured and will activate automatically when certificates are present.

### 1. Obtain SSL Certificates

**Development (self-signed):**
```bash
cd docker/nginx/ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout key.pem \
  -out cert.pem \
  -subj "/C=US/ST=State/L=City/O=Suipic/CN=localhost"
```

**Production (Let's Encrypt):**
```bash
# Install certbot
sudo apt-get install certbot

# Generate certificate (nginx must be stopped)
docker-compose stop nginx
sudo certbot certonly --standalone -d yourdomain.com

# Copy certificates
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem docker/nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem docker/nginx/ssl/key.pem
sudo chown $USER:$USER docker/nginx/ssl/*.pem
```

See `ssl/README.md` for more details.

### 2. Enable HTTP to HTTPS Redirect (Optional)

Edit `conf.d/default.conf` and uncomment these lines in the HTTP server block:

```nginx
# location / {
#     return 301 https://$host$request_uri;
# }
```

### 3. Update Domain Name (Optional)

Replace `server_name _;` with your domain:

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    # ...
}
```

### 4. Restart Nginx

```bash
docker-compose restart nginx
```

### 5. Verify HTTPS

```bash
curl -k https://localhost/health
docker-compose exec nginx nginx -T | grep ssl_certificate
```

## Certificate Mounting

SSL certificates are mounted via Docker volumes in `docker-compose.yml`:

```yaml
nginx:
  volumes:
    - ./docker/nginx/ssl:/etc/nginx/ssl:ro
```

The `:ro` flag mounts certificates as read-only for security.

**Important:** Certificate files are excluded from git via `.gitignore`.

## Custom Configuration

### Changing Upload Size Limit

Edit `nginx.conf`:

```nginx
http {
    client_max_body_size 500M;  # Change from 100M to 500M
}
```

### Adjusting Rate Limits

Edit rate zones in `nginx.conf`:

```nginx
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=20r/s;  # Increased from 10r/s
```

Then update bursts in `conf.d/default.conf`:

```nginx
limit_req zone=api_limit burst=50 nodelay;  # Increased from 20
```

### Adding Custom Headers

Add to server blocks in `conf.d/default.conf`:

```nginx
server {
    add_header Custom-Header "value" always;
}
```

## Proxy Headers

The configuration sets these headers for proper request handling:

- `X-Real-IP`: Client's real IP address
- `X-Forwarded-For`: Chain of proxy IPs
- `X-Forwarded-Proto`: Original protocol (http/https)
- `Host`: Original host header
- `Upgrade` & `Connection`: For WebSocket support (frontend)

## Health Check

Nginx includes a `/health` endpoint that always returns 200 OK:

```bash
curl http://localhost/health
# Output: healthy
```

This is useful for load balancer health checks and monitoring.

## Testing Configuration

Test configuration without restarting:

```bash
docker-compose exec nginx nginx -t
```

Reload configuration (no downtime):

```bash
docker-compose exec nginx nginx -s reload
```

View full configuration:

```bash
docker-compose exec nginx nginx -T
```

## Logging

### Access Logs

View access logs:
```bash
docker-compose logs nginx
docker-compose exec nginx tail -f /var/log/nginx/access.log
```

### Error Logs

View error logs:
```bash
docker-compose exec nginx tail -f /var/log/nginx/error.log
```

### Log Formats

Two log formats are configured:

1. **main**: Standard access log format
2. **detailed**: Includes request time, upstream response time, caching status

To use detailed format, edit `nginx.conf`:
```nginx
access_log /var/log/nginx/access.log detailed;
```

## Performance Tuning

### Worker Processes

Default is `auto` (matches CPU cores). For specific tuning:

```nginx
worker_processes 4;  # Set to number of CPU cores
```

### Worker Connections

Default is 1024. For high-traffic:

```nginx
events {
    worker_connections 2048;
}
```

### Keepalive Connections

Increase upstream keepalive pool:

```nginx
upstream backend {
    server backend:3000;
    keepalive 64;  # Increased from 32
}
```

## Troubleshooting

### Rate Limiting Too Strict

Users getting 503 errors? Increase burst values:

```nginx
limit_req zone=api_limit burst=50 nodelay;  # Increased from 20
```

### 502 Bad Gateway

Check if backend/frontend services are running:
```bash
docker-compose ps
docker-compose logs backend frontend
curl -v http://backend:3000/api/health
```

### 413 Request Entity Too Large

Increase `client_max_body_size` in `nginx.conf`:
```nginx
client_max_body_size 500M;
```

### SSL Certificate Errors

Verify certificate files:
```bash
docker-compose exec nginx ls -la /etc/nginx/ssl/
docker-compose exec nginx openssl x509 -in /etc/nginx/ssl/cert.pem -text -noout
```

### Connection Timeouts

Increase timeout values in `conf.d/default.conf`:
```nginx
proxy_connect_timeout 90s;
proxy_read_timeout 600s;
```

### High Memory Usage

Reduce buffer sizes in `nginx.conf`:
```nginx
client_body_buffer_size 64k;  # Reduced from 128k
```

## Advanced Topics

### Load Balancing

For multiple backend instances:

```nginx
upstream backend {
    least_conn;
    server backend1:3000 max_fails=3 fail_timeout=30s;
    server backend2:3000 max_fails=3 fail_timeout=30s;
    server backend3:3000 max_fails=3 fail_timeout=30s;
    keepalive 32;
}
```

### IP Whitelisting

Restrict access to specific IPs:

```nginx
location /admin {
    allow 192.168.1.0/24;
    allow 10.0.0.0/8;
    deny all;
    proxy_pass http://backend;
}
```

### Request Buffering

Disable buffering for large uploads:

```nginx
location /api/photos {
    proxy_request_buffering off;
    proxy_buffering off;
    proxy_pass http://backend;
}
```

### Custom Error Pages

```nginx
error_page 502 503 504 /50x.html;
location = /50x.html {
    root /usr/share/nginx/html;
}
```

## Further Reading

- [Nginx Documentation](https://nginx.org/en/docs/)
- [Nginx Reverse Proxy Guide](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/)
- [SSL Configuration Guide](https://ssl-config.mozilla.org/)
- [Nginx Rate Limiting](https://www.nginx.com/blog/rate-limiting-nginx/)
- [Nginx HTTP/2 Guide](https://www.nginx.com/blog/http2-module-nginx/)

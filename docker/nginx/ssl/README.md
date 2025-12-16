# SSL Certificate Configuration

This directory should contain your SSL certificates for HTTPS configuration.

## Self-Signed Certificates (Development)

To generate self-signed certificates for local development:

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout key.pem \
  -out cert.pem \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
```

## Production Certificates

For production, use certificates from a trusted Certificate Authority (CA) like:
- Let's Encrypt (free, automated)
- DigiCert
- Comodo
- GlobalSign

### Let's Encrypt with Certbot

```bash
# Install certbot
sudo apt-get install certbot

# Generate certificate
sudo certbot certonly --standalone -d yourdomain.com

# Copy certificates to this directory
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ./cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ./key.pem
```

## Required Files

- `cert.pem` - SSL certificate (public key)
- `key.pem` - Private key

## Security Notes

- Never commit actual certificate files to version control
- Keep private keys secure and restrict file permissions (chmod 600)
- Use strong encryption algorithms (TLS 1.2+)
- Regularly renew certificates before expiration
- Consider using certificate management tools for automation

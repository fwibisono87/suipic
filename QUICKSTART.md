# Suipic Quick Start Guide

Get Suipic running in 5 minutes!

## Prerequisites

- **Docker Desktop** installed and running
- **4GB RAM** available
- **10GB disk space** available

## Installation Steps

### 1. Clone the Repository

```bash
git clone <repository-url>
cd suipic
```

### 2. Run Setup Script

**Linux/Mac:**
```bash
chmod +x setup.sh
./setup.sh --docker
```

**Windows (PowerShell):**
```powershell
.\setup.ps1 -Docker
```

### 3. Configure Environment (Optional)

Edit `.env` to change default passwords and settings:

```bash
nano .env  # or use any text editor
```

**Minimum recommended changes:**
- `POSTGRES_PASSWORD` - Change from default
- `JWT_SECRET` - Generate with: `openssl rand -base64 32`
- `ADMIN_PASSWORD` - Change from default
- `MINIO_ROOT_PASSWORD` - Change from default

### 4. Start Services

**Linux/Mac:**
```bash
./docker-start.sh
```

**Windows:**
```batch
docker-start.bat
```

Or manually:
```bash
docker-compose up -d
```

### 5. Verify Installation

Wait 30-60 seconds for all services to start, then check:

```bash
# View logs
docker-compose logs -f

# Check service status
docker-compose ps
```

All services should show "healthy" or "running" status.

### 6. Access the Application

Open your browser and navigate to:

- **Application**: http://localhost
- **MinIO Console**: http://localhost:9001

**Default Login Credentials:**
- Email: `admin@suipic.local`
- Password: `admin123`

**⚠️ Important**: Change the admin password immediately after first login!

## What Gets Set Up Automatically

The setup process automatically:

1. ✅ Creates PostgreSQL database
2. ✅ Runs database migrations
3. ✅ Creates admin user (if configured in .env)
4. ✅ Creates MinIO bucket for photo storage
5. ✅ Initializes Elasticsearch
6. ✅ Configures Nginx reverse proxy
7. ✅ Sets up all service connections

## First Steps After Login

1. **Change admin password**
   - Click on user menu
   - Go to Settings
   - Update password

2. **Create a photographer account**
   - Go to Admin panel
   - Create new photographer user

3. **Create your first album**
   - Click "New Album"
   - Give it a name and description

4. **Upload photos**
   - Open the album
   - Click "Upload Photos"
   - Select and upload your photos

5. **Try the search**
   - Use the search bar
   - Filter by date, stars, state
   - Search in comments

## Useful Commands

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f postgres
```

### Restart Services
```bash
# All services
docker-compose restart

# Specific service
docker-compose restart backend
```

### Stop Services
```bash
# Stop but keep data
docker-compose stop

# Stop and remove containers (keeps data)
docker-compose down

# Stop and remove everything including data
docker-compose down -v
```

### Check Service Health
```bash
# Service status
docker-compose ps

# Backend health
curl http://localhost/api/health

# Elasticsearch health
curl http://localhost:9200/_cluster/health
```

## Troubleshooting

### Services Won't Start

```bash
# Check logs for errors
docker-compose logs

# Check if ports are in use
# On Linux/Mac:
sudo lsof -i :80
sudo lsof -i :5432

# On Windows:
netstat -ano | findstr :80
netstat -ano | findstr :5432
```

### Can't Login

1. Check if admin user was created:
```bash
docker-compose logs backend | grep -i admin
```

2. Verify credentials in `.env`:
```bash
grep ADMIN .env
```

3. Try resetting: stop services, remove volumes, restart
```bash
docker-compose down -v
docker-compose up -d
```

### Photos Won't Upload

1. Check MinIO is running:
```bash
docker-compose logs minio
curl http://localhost:9000/minio/health/live
```

2. Check backend logs:
```bash
docker-compose logs backend | grep -i minio
```

### Search Not Working

1. Check Elasticsearch:
```bash
curl http://localhost:9200/_cluster/health
docker-compose logs elasticsearch
```

2. Re-index photos:
   - Go to an album
   - Click "Re-index" or use API: `POST /api/albums/:id/index`

## Need More Help?

- **Detailed Setup**: See [SETUP.md](SETUP.md)
- **Docker Guide**: See [DOCKER.md](DOCKER.md)
- **API Docs**: See [backend/README.md](backend/README.md)
- **Development**: See [AGENTS.md](AGENTS.md)

## Next Steps

Once you're up and running:

1. Review [SETUP.md](SETUP.md) for production deployment
2. Configure SSL/HTTPS for production use
3. Set up regular database backups
4. Review security checklist
5. Customize branding and settings

## Getting Updates

```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
docker-compose build --no-cache
docker-compose up -d
```

## Uninstalling

To completely remove Suipic:

```bash
# Stop and remove everything
docker-compose down -v --rmi all

# Delete the repository
cd ..
rm -rf suipic
```

---

**Happy photo managing! 📸**

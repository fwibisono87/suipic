#!/bin/bash

# Suipic Docker Quick Start Script
set -e

echo "======================================"
echo "   Suipic Docker Setup"
echo "======================================"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed."
    echo "Please install Docker from https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Error: Docker Compose is not installed."
    echo "Please install Docker Compose from https://docs.docker.com/compose/install/"
    exit 1
fi

# Check if .env exists, if not copy from .env.docker
if [ ! -f .env ]; then
    echo "Creating .env file from .env.docker..."
    cp .env.docker .env
    echo "✓ .env file created"
    echo ""
    echo "⚠️  WARNING: Please edit .env and change default passwords before deploying to production!"
    echo ""
else
    echo "✓ .env file already exists"
fi

# Create necessary directories
echo "Creating required directories..."
mkdir -p docker/nginx/conf.d docker/nginx/ssl docker/postgres/init-scripts
echo "✓ Directories created"
echo ""

# Ask user what to start
echo "What would you like to start?"
echo "1) Full stack (all services including backend and frontend)"
echo "2) Infrastructure only (PostgreSQL, Elasticsearch, MinIO)"
echo "3) Production mode (with production optimizations)"
echo ""
read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        echo ""
        echo "Starting full stack..."
        docker-compose up -d
        ;;
    2)
        echo ""
        echo "Starting infrastructure services only..."
        docker-compose -f docker-compose.infra.yml up -d
        ;;
    3)
        echo ""
        echo "Starting in production mode..."
        docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
        ;;
    *)
        echo "Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "Waiting for services to be ready..."
sleep 10

# Check service status
echo ""
echo "======================================"
echo "   Service Status"
echo "======================================"
docker-compose ps

echo ""
echo "======================================"
echo "   Access Information"
echo "======================================"
echo "Frontend:          http://localhost"
echo "Backend API:       http://localhost/api"
echo "Backend (direct):  http://localhost:3000"
echo "MinIO Console:     http://localhost:9001"
echo "Elasticsearch:     http://localhost:9200"
echo "PostgreSQL:        localhost:5432"
echo ""
echo "Default credentials:"
echo "  Admin User:  admin@suipic.local / admin123"
echo "  MinIO:       minioadmin / minioadmin"
echo "  PostgreSQL:  suipic / suipic_password"
echo ""
echo "======================================"
echo ""
echo "Useful commands:"
echo "  View logs:       docker-compose logs -f"
echo "  Stop services:   docker-compose down"
echo "  Restart:         docker-compose restart"
echo "  See all commands: make help"
echo ""
echo "Setup complete! 🚀"

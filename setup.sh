#!/bin/bash
set -e

echo "=========================================="
echo "Suipic Initial Setup Script"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_info() {
    echo "ℹ $1"
}

# Check if running with Docker
if [ "$1" == "--docker" ]; then
    DOCKER_MODE=true
    print_info "Running in Docker mode"
else
    DOCKER_MODE=false
    print_info "Running in local development mode"
fi

echo ""
echo "Step 1: Environment Configuration"
echo "----------------------------------"

# Check for .env file
if [ "$DOCKER_MODE" = true ]; then
    if [ ! -f ".env" ]; then
        print_warning ".env file not found, creating from .env.docker template"
        cp .env.docker .env
        print_success "Created .env file"
        print_warning "Please review and update .env with your configuration!"
    else
        print_success ".env file exists"
    fi
else
    # Local development mode
    if [ ! -f "backend/.env" ]; then
        print_warning "backend/.env file not found, creating from .env.example"
        cp backend/.env.example backend/.env
        print_success "Created backend/.env file"
        print_warning "Please review and update backend/.env with your configuration!"
    else
        print_success "backend/.env file exists"
    fi
fi

echo ""
echo "Step 2: Prerequisites Check"
echo "----------------------------"

if [ "$DOCKER_MODE" = true ]; then
    # Check Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    print_success "Docker is installed"

    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    print_success "Docker Compose is installed"
else
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21+ first."
        exit 1
    fi
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go $GO_VERSION is installed"

    # Check PostgreSQL
    if ! command -v psql &> /dev/null; then
        print_warning "PostgreSQL client not found. Make sure PostgreSQL server is running."
    else
        print_success "PostgreSQL client is installed"
    fi

    # Check if required services are running
    print_info "Checking if required services are accessible..."
    
    # Source environment variables
    if [ -f "backend/.env" ]; then
        export $(grep -v '^#' backend/.env | xargs)
    fi
fi

echo ""
echo "Step 3: Backend Setup"
echo "---------------------"

if [ "$DOCKER_MODE" = false ]; then
    cd backend
    
    print_info "Downloading Go dependencies..."
    go mod download
    print_success "Dependencies downloaded"
    
    print_info "Running database migrations..."
    if go run cmd/migrate/main.go -action=up 2>/dev/null; then
        print_success "Database migrations completed"
    else
        print_warning "Migration failed. Ensure PostgreSQL is running and configured correctly."
    fi
    
    cd ..
fi

echo ""
echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""

if [ "$DOCKER_MODE" = true ]; then
    echo "Next steps:"
    echo "1. Review and update .env file with your configuration"
    echo "2. Start the application with: ./docker-start.sh"
    echo "   Or manually: docker-compose up -d"
    echo ""
    echo "The application will:"
    echo "  - Create MinIO bucket automatically"
    echo "  - Run database migrations automatically"
    echo "  - Create admin user (if configured in .env)"
else
    echo "Next steps:"
    echo "1. Ensure PostgreSQL, Elasticsearch, and MinIO are running"
    echo "2. Review and update backend/.env with your configuration"
    echo "3. Start the backend: cd backend && go run main.go"
    echo ""
    echo "For infrastructure only (using Docker):"
    echo "  docker-compose -f docker-compose.infra.yml up -d"
    echo ""
    echo "The application will:"
    echo "  - Create MinIO bucket automatically on first start"
    echo "  - Create admin user on first start (if configured)"
fi

echo ""
echo "Default credentials:"
echo "  Admin: admin@suipic.local / admin123"
echo "  MinIO: minioadmin / minioadmin"
echo ""
print_warning "Remember to change default passwords in production!"
echo ""

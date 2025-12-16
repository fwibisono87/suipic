# Suipic Initial Setup Script (PowerShell)
param(
    [switch]$Docker
)

function Write-Success {
    param([string]$Message)
    Write-Host "✓ " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Error {
    param([string]$Message)
    Write-Host "✗ " -ForegroundColor Red -NoNewline
    Write-Host $Message
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠ " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ " -NoNewline
    Write-Host $Message
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Suipic Initial Setup Script" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

if ($Docker) {
    Write-Info "Running in Docker mode"
} else {
    Write-Info "Running in local development mode"
}

Write-Host ""
Write-Host "Step 1: Environment Configuration"
Write-Host "----------------------------------"

# Check for .env file
if ($Docker) {
    if (-not (Test-Path ".env")) {
        Write-Warning ".env file not found, creating from .env.docker template"
        Copy-Item ".env.docker" ".env"
        Write-Success "Created .env file"
        Write-Warning "Please review and update .env with your configuration!"
    } else {
        Write-Success ".env file exists"
    }
} else {
    # Local development mode
    if (-not (Test-Path "backend\.env")) {
        Write-Warning "backend\.env file not found, creating from .env.example"
        Copy-Item "backend\.env.example" "backend\.env"
        Write-Success "Created backend\.env file"
        Write-Warning "Please review and update backend\.env with your configuration!"
    } else {
        Write-Success "backend\.env file exists"
    }
}

Write-Host ""
Write-Host "Step 2: Prerequisites Check"
Write-Host "----------------------------"

if ($Docker) {
    # Check Docker
    try {
        $null = docker --version
        Write-Success "Docker is installed"
    } catch {
        Write-Error "Docker is not installed. Please install Docker Desktop first."
        exit 1
    }

    # Check Docker Compose
    try {
        $null = docker-compose --version
        Write-Success "Docker Compose is installed"
    } catch {
        try {
            $null = docker compose version
            Write-Success "Docker Compose (plugin) is installed"
        } catch {
            Write-Error "Docker Compose is not installed. Please install Docker Compose first."
            exit 1
        }
    }
} else {
    # Check Go
    try {
        $goVersion = go version
        Write-Success "Go is installed: $goVersion"
    } catch {
        Write-Error "Go is not installed. Please install Go 1.21+ first."
        exit 1
    }

    # Check PostgreSQL client
    try {
        $null = psql --version
        Write-Success "PostgreSQL client is installed"
    } catch {
        Write-Warning "PostgreSQL client not found. Make sure PostgreSQL server is running."
    }
}

Write-Host ""
Write-Host "Step 3: Backend Setup"
Write-Host "---------------------"

if (-not $Docker) {
    Push-Location backend
    
    Write-Info "Downloading Go dependencies..."
    go mod download
    Write-Success "Dependencies downloaded"
    
    Write-Info "Running database migrations..."
    try {
        go run cmd/migrate/main.go -action=up 2>$null
        Write-Success "Database migrations completed"
    } catch {
        Write-Warning "Migration failed. Ensure PostgreSQL is running and configured correctly."
    }
    
    Pop-Location
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Setup Complete!" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

if ($Docker) {
    Write-Host "Next steps:"
    Write-Host "1. Review and update .env file with your configuration"
    Write-Host "2. Start the application with: .\docker-start.bat"
    Write-Host "   Or manually: docker-compose up -d"
    Write-Host ""
    Write-Host "The application will:"
    Write-Host "  - Create MinIO bucket automatically"
    Write-Host "  - Run database migrations automatically"
    Write-Host "  - Create admin user (if configured in .env)"
} else {
    Write-Host "Next steps:"
    Write-Host "1. Ensure PostgreSQL, Elasticsearch, and MinIO are running"
    Write-Host "2. Review and update backend\.env with your configuration"
    Write-Host "3. Start the backend: cd backend; go run main.go"
    Write-Host ""
    Write-Host "For infrastructure only (using Docker):"
    Write-Host "  docker-compose -f docker-compose.infra.yml up -d"
    Write-Host ""
    Write-Host "The application will:"
    Write-Host "  - Create MinIO bucket automatically on first start"
    Write-Host "  - Create admin user on first start (if configured)"
}

Write-Host ""
Write-Host "Default credentials:"
Write-Host "  Admin: admin@suipic.local / admin123"
Write-Host "  MinIO: minioadmin / minioadmin"
Write-Host ""
Write-Warning "Remember to change default passwords in production!"
Write-Host ""

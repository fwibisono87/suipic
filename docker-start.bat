@echo off
setlocal enabledelayedexpansion

echo ======================================
echo    Suipic Docker Setup
echo ======================================
echo.

REM Check if Docker is installed
docker --version >nul 2>&1
if errorlevel 1 (
    echo Error: Docker is not installed.
    echo Please install Docker Desktop from https://docs.docker.com/get-docker/
    exit /b 1
)

REM Check if Docker Compose is installed
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo Error: Docker Compose is not installed.
    echo Please install Docker Desktop which includes Docker Compose
    exit /b 1
)

REM Check if .env exists, if not copy from .env.docker
if not exist .env (
    echo Creating .env file from .env.docker...
    copy .env.docker .env >nul
    echo [92m✓[0m .env file created
    echo.
    echo [93m⚠️  WARNING: Please edit .env and change default passwords before deploying to production![0m
    echo.
) else (
    echo [92m✓[0m .env file already exists
)

REM Create necessary directories
echo Creating required directories...
if not exist docker\nginx\conf.d mkdir docker\nginx\conf.d
if not exist docker\nginx\ssl mkdir docker\nginx\ssl
if not exist docker\postgres\init-scripts mkdir docker\postgres\init-scripts
echo [92m✓[0m Directories created
echo.

REM Ask user what to start
echo What would you like to start?
echo 1) Full stack (all services including backend and frontend)
echo 2) Infrastructure only (PostgreSQL, Elasticsearch, MinIO)
echo 3) Production mode (with production optimizations)
echo.
set /p choice="Enter your choice (1-3): "

if "%choice%"=="1" (
    echo.
    echo Starting full stack...
    docker-compose up -d
) else if "%choice%"=="2" (
    echo.
    echo Starting infrastructure services only...
    docker-compose -f docker-compose.infra.yml up -d
) else if "%choice%"=="3" (
    echo.
    echo Starting in production mode...
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
) else (
    echo Invalid choice. Exiting.
    exit /b 1
)

echo.
echo Waiting for services to be ready...
timeout /t 10 /nobreak >nul

REM Check service status
echo.
echo ======================================
echo    Service Status
echo ======================================
docker-compose ps

echo.
echo ======================================
echo    Access Information
echo ======================================
echo Frontend:          http://localhost
echo Backend API:       http://localhost/api
echo Backend (direct):  http://localhost:3000
echo MinIO Console:     http://localhost:9001
echo Elasticsearch:     http://localhost:9200
echo PostgreSQL:        localhost:5432
echo.
echo Default credentials:
echo   Admin User:  admin@suipic.local / admin123
echo   MinIO:       minioadmin / minioadmin
echo   PostgreSQL:  suipic / suipic_password
echo.
echo ======================================
echo.
echo Useful commands:
echo   View logs:       docker-compose logs -f
echo   Stop services:   docker-compose down
echo   Restart:         docker-compose restart
echo.
echo Setup complete! 🚀

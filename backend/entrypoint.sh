#!/bin/sh
set -e

echo "=========================================="
echo "Starting Suipic Backend Initialization"
echo "=========================================="

# Wait for PostgreSQL to be ready
echo "[1/5] Waiting for PostgreSQL..."
until nc -z -v -w30 $DB_HOST $DB_PORT 2>/dev/null
do
  echo "  Waiting for PostgreSQL connection at $DB_HOST:$DB_PORT..."
  sleep 2
done
echo "  ✓ PostgreSQL is ready!"

# Wait for Elasticsearch to be ready
echo "[2/5] Waiting for Elasticsearch..."
until curl -f http://elasticsearch:9200/_cluster/health > /dev/null 2>&1
do
  echo "  Waiting for Elasticsearch connection..."
  sleep 2
done
echo "  ✓ Elasticsearch is ready!"

# Wait for MinIO to be ready
echo "[3/5] Waiting for MinIO..."
until curl -f http://minio:9000/minio/health/live > /dev/null 2>&1
do
  echo "  Waiting for MinIO connection..."
  sleep 2
done
echo "  ✓ MinIO is ready!"

# Run database migrations
echo "[4/5] Running database migrations..."
if [ -f "./cmd/migrate/main.go" ]; then
  go run cmd/migrate/main.go -action=up
  echo "  ✓ Migrations completed!"
else
  echo "  ⚠ Migration script not found, skipping..."
fi

# Note: Admin user and MinIO bucket initialization happens automatically
# when the application starts via services initialization
echo "[5/5] Initializing services..."
echo "  - Admin user will be created if credentials are provided"
echo "  - MinIO bucket will be created automatically"
echo "  - ElasticSearch index will be initialized on first use"

echo "=========================================="
echo "Initialization Complete!"
echo "=========================================="

# Start the application
echo "Starting application..."
exec "$@"

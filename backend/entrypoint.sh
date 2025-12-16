#!/bin/sh
set -e

echo "Starting Suipic backend..."

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL..."
until nc -z -v -w30 $DB_HOST $DB_PORT
do
  echo "Waiting for PostgreSQL connection..."
  sleep 2
done
echo "PostgreSQL is up!"

# Wait for Elasticsearch to be ready
echo "Waiting for Elasticsearch..."
until curl -f http://elasticsearch:9200/_cluster/health > /dev/null 2>&1
do
  echo "Waiting for Elasticsearch connection..."
  sleep 2
done
echo "Elasticsearch is up!"

# Run database migrations
echo "Running database migrations..."
if [ -f "./cmd/migrate/main.go" ]; then
  go run cmd/migrate/main.go -action=up || echo "Migration failed or already applied"
fi

# Start the application
echo "Starting application..."
exec "$@"

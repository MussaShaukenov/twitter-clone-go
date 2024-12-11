#!/bin/sh

echo "Running migrations..."

# Use the Docker Compose service name 'db' as the hostname
goose -dir "$GOOSE_MIGRATION_DIR" postgres "$DATABASE_URL" up

echo "Starting the application..."
exec "$@"

#!/bin/sh

set -e

echo "Running DB migrations..."
/usr/bin/migrate -path /app/db/migrate -database "$DB_SOURCE" -verbose up

echo "Starting the app..."
exec "$@"
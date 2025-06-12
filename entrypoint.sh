#!/bin/sh

set -e

echo "Running database migrations..."
migrate -path db/migrate -database "postgresql://root:secret@postgres:5432/bilitioo?sslmode=disable" -verbose up

echo "Starting the application..."
exec go run main.go

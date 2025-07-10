#!/bin/bash

echo "Starting Page Hoppers development environment..."

# Start the database
echo "Starting PostgreSQL database..."
docker-compose up -d postgres

# Wait for database to be ready
echo "Waiting for database to be ready..."
until docker-compose exec -T postgres pg_isready -U pagehoppers_user -d pagehoppers_db; do
  echo "Database is not ready yet. Waiting..."
  sleep 2
done

echo "Database is ready!"

# Start the backend
echo "Starting backend server..."
cd page-hoppers-backend
go run main.go 
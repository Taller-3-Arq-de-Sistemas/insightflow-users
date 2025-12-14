#!/bin/sh
set -e

# Ensure data directory exists for SQLite database
mkdir -p /app/data

echo "Running migrations..."
./migrate up

echo "Seeding database..."
./seed

echo "Starting server..."
./server

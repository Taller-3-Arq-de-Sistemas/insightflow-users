#!/bin/sh
set -e

echo "Running migrations..."
./migrate up

echo "Seeding database..."
./seed

echo "Starting server..."
./server

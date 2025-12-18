#!/bin/bash
set -e

echo "========================================"
echo "Building Candidate Organizer for Heroku"
echo "========================================"

# Build Frontend
echo ""
echo "1. Building Next.js frontend..."
cd frontend
npm install --include=dev
npm run build
cd ..

# Build Backend
echo ""
echo "2. Building Go backend..."
cd backend

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the Go binary
echo "   Compiling Go server..."
go build -o bin/server ./cmd/server

# Make it executable
chmod +x bin/server

cd ..

# Make start-services.sh executable
chmod +x start-services.sh

echo ""
echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo "Frontend: frontend/.next/"
echo "Backend:  backend/bin/server"
echo "========================================"

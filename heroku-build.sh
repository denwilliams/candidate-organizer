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

# Copy static assets to standalone directory (required for standalone mode)
echo "   Copying static assets to standalone..."
cp -r .next/static .next/standalone/.next/static
if [ -d "public" ]; then
    cp -r public .next/standalone/public
fi

cd ..

# Prepare Backend for Go buildpack
echo ""
echo "2. Preparing backend files for Go buildpack..."
echo "   Copying backend source to root..."
cp -r backend/cmd backend/internal backend/migrations .

# Make start-services.sh executable
chmod +x start-services.sh

echo ""
echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo "Frontend: frontend/.next/"
echo "Backend:  backend files copied to root"
echo "         (Go buildpack will compile binary)"
echo "========================================"

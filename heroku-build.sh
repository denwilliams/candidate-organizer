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

echo ""
echo "   Verifying build output..."
echo "   Contents of .next/:"
ls -la .next/ | head -20
echo ""
echo "   Contents of .next/standalone/:"
ls -la .next/standalone/ 2>/dev/null || echo "   (standalone directory not found)"

# Copy static assets to standalone directory (required for standalone mode)
echo ""
echo "   Copying static assets to standalone..."

# Check if this is a monorepo structure (frontend/ subdirectory in standalone)
if [ -d ".next/standalone/frontend" ]; then
    echo "   Detected monorepo structure"
    mkdir -p .next/standalone/frontend/.next
    cp -r .next/static .next/standalone/frontend/.next/static
    echo "   ✓ Copied .next/static to standalone/frontend/.next/static"

    if [ -d "public" ]; then
        cp -r public .next/standalone/frontend/public
        echo "   ✓ Copied public to standalone/frontend/public"
    else
        echo "   (no public directory to copy)"
    fi
else
    echo "   Detected standard structure"
    mkdir -p .next/standalone/.next
    cp -r .next/static .next/standalone/.next/static
    echo "   ✓ Copied .next/static to standalone/.next/static"

    if [ -d "public" ]; then
        cp -r public .next/standalone/public
        echo "   ✓ Copied public to standalone/public"
    else
        echo "   (no public directory to copy)"
    fi
fi

echo ""
echo "   Final standalone structure:"
ls -la .next/standalone/ | head -15

echo ""
echo "   Checking for server.js location..."
if [ -f ".next/standalone/frontend/server.js" ]; then
    echo "   ✓ Found: .next/standalone/frontend/server.js (monorepo structure)"
    ls -la .next/standalone/frontend/ | head -10
elif [ -f ".next/standalone/server.js" ]; then
    echo "   ✓ Found: .next/standalone/server.js (standard structure)"
else
    echo "   ⚠ WARNING: server.js not found in expected locations!"
fi

cd ..

# Prepare Backend for Go buildpack
echo ""
echo "2. Preparing backend files for Go buildpack..."
echo "   Copying backend source to root..."
cp -r backend/cmd backend/internal backend/migrations .
echo "   ✓ Copied cmd, internal, migrations to root"

echo ""
echo "   Verifying backend files at root:"
ls -la cmd internal migrations 2>/dev/null | head -10

# Make start-services.sh executable
chmod +x start-services.sh

echo ""
echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo "Frontend: frontend/.next/standalone/"
echo "  - Static assets: ✓"
echo "  - Server: frontend/.next/standalone/server.js"
echo ""
echo "Backend: Files copied to root"
echo "  - Go buildpack will compile to: bin/server"
echo "========================================"

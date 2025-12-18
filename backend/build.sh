#!/bin/bash
set -e

echo "Building frontend..."
cd ../frontend
npm install
npm run build
cd ../backend

echo "Copying static files..."
mkdir -p static
cp -r ../frontend/out/* static/

echo "Build complete!"

#!/bin/bash

# Process manager script for running Next.js frontend and Go backend
# If either process crashes, kill both and exit so Heroku restarts the dyno

set -e

echo "Starting Candidate Organizer services..."

# Trap EXIT to ensure cleanup
cleanup() {
    echo "Shutting down services..."
    kill $(jobs -p) 2>/dev/null || true
    exit
}
trap cleanup EXIT INT TERM

# Start Go backend on port 8080
echo "Starting Go backend on port 8080..."
cd backend
PORT=8080 ./bin/server &
BACKEND_PID=$!
cd ..

# Give backend a moment to start
sleep 2

# Start Next.js frontend on $PORT (from Heroku)
echo "Starting Next.js frontend on port $PORT..."
cd frontend/.next/standalone
node server.js &
FRONTEND_PID=$!
cd ../../..

echo "Services started:"
echo "  - Backend PID: $BACKEND_PID"
echo "  - Frontend PID: $FRONTEND_PID"

# Monitor both processes
while true; do
    # Check if backend is still running
    if ! kill -0 $BACKEND_PID 2>/dev/null; then
        echo "ERROR: Backend process ($BACKEND_PID) has died!"
        echo "Killing all processes to trigger dyno restart..."
        kill $FRONTEND_PID 2>/dev/null || true
        exit 1
    fi

    # Check if frontend is still running
    if ! kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "ERROR: Frontend process ($FRONTEND_PID) has died!"
        echo "Killing all processes to trigger dyno restart..."
        kill $BACKEND_PID 2>/dev/null || true
        exit 1
    fi

    # Check every 5 seconds
    sleep 5
done

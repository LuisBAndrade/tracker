#!/bin/bash

# --------------------------
# Kill any old backend on port 3000
# --------------------------
OLD_PID=$(lsof -ti:3000)
if [ -n "$OLD_PID" ]; then
    echo "Killing old backend process on port 3000 (PID $OLD_PID)"
    kill -9 $OLD_PID
fi

# --------------------------
# Start Go backend
# --------------------------
cd server
export $(cat .env | xargs)   # load environment variables
go build -o app ./cmd/api    # build Go binary
./app &                      # run backend in background
BACKEND_PID=$!
echo "Backend started (PID $BACKEND_PID)"

# --------------------------
# Serve frontend (production)
# --------------------------
cd ../frontend

# Install local serve if not installed
if [ ! -d "node_modules/serve" ]; then
    npm install serve
fi

npx serve -s dist -l 5173 &
FRONTEND_PID=$!
echo "Frontend started (PID $FRONTEND_PID)"

# --------------------------
# Wait for both processes
# --------------------------
wait $BACKEND_PID $FRONTEND_PID



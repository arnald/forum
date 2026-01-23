#!/bin/sh
set -e

echo "🚀 Starting Forum Application..."

# Ensure database directory exists
mkdir -p db/data

# Function to handle shutdown gracefully
cleanup() {
    echo ""
    echo "🛑 Shutting down services..."
    kill $SERVER_PID $CLIENT_PID 2>/dev/null || true
    wait $SERVER_PID $CLIENT_PID 2>/dev/null || true
    echo "✅ Services stopped"
    exit 0
}

# Trap signals for graceful shutdown
trap cleanup SIGTERM SIGINT SIGQUIT

# Start backend server in background
echo "🔧 Starting backend server on port ${SERVER_PORT:-8080}..."
./server &
SERVER_PID=$!

# Wait a moment for server to initialize
sleep 2

# Check if server started successfully
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "❌ Backend server failed to start"
    exit 1
fi

echo "✅ Backend server started (PID: $SERVER_PID)"

# Start frontend client in background
echo "🌐 Starting frontend client on port ${CLIENT_PORT:-3001}..."
./client &
CLIENT_PID=$!

# Wait a moment for client to initialize
sleep 2

# Check if client started successfully
if ! kill -0 $CLIENT_PID 2>/dev/null; then
    echo "❌ Frontend client failed to start"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi

echo "✅ Frontend client started (PID: $CLIENT_PID)"
echo ""
echo "🎉 Forum application is running!"
echo "   Frontend: http://localhost:${CLIENT_PORT:-3001}"
echo "   Backend:  http://localhost:${SERVER_PORT:-8080}/api/v1"
echo ""
echo "Press Ctrl+C to stop..."

# Wait for both processes
wait $SERVER_PID $CLIENT_PID

#!/bin/bash
# Start LoFP development servers
set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Starting Legends of Future Past..."

# Load environment variables from .env if present
if [ -f "$ROOT_DIR/.env" ]; then
  set -a
  source "$ROOT_DIR/.env"
  set +a
fi

# Kill existing processes on our ports
kill $(lsof -ti:4993) 2>/dev/null || true
kill $(lsof -ti:4992) 2>/dev/null || true
sleep 1

# Start backend
cd "$ROOT_DIR/engine"
go run cmd/lofp/main.go &
BACKEND_PID=$!

# Start frontend
cd "$ROOT_DIR/frontend"
npx vite --port 4992 &
FRONTEND_PID=$!

echo "Backend PID: $BACKEND_PID (port 4993)"
echo "Frontend PID: $FRONTEND_PID (port 4992)"
echo ""
echo "Open http://localhost:4992 to play!"
echo "Press Ctrl+C to stop."

trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0" INT TERM
wait

#!/bin/bash

# Start the web dashboard server
# This starts the server with web UI for viewing jobs and managing companies

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

echo "🌐 Starting Intern Job Tracker Web Dashboard..."
echo "📂 Project Dir: $DIR"
echo "----------------------------------------"
echo "📊 Dashboard: http://localhost:8080"
echo "📡 API: http://localhost:8080/api/jobs"
echo "----------------------------------------"
echo "Press Ctrl+C to stop"
echo ""

export PATH="/opt/homebrew/bin:$PATH"
go run ./cmd/server -recipient="+17246807862"

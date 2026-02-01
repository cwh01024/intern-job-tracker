#!/bin/bash

# Run-once job check script
# This script runs the job tracker once, scrapes all companies, and sends notifications

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

echo "� Starting Intern Job Tracker..."
echo "📱 Recipient: +17246807862"
echo "⏰ Time: $(date)"
echo "----------------------------------------"

# Run one-time check (not server mode)
export PATH="/opt/homebrew/bin:$PATH"
go run ./cmd/server -run-once -recipient="+17246807862"

echo "----------------------------------------"
echo "✅ Check complete!"

#!/bin/bash

# Get the directory where this script is located
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

echo "🚀 Starting Intern Job Tracker..."
echo "📱 Recipient: +17246807862"
echo "📂 Project Dir: $DIR"
echo "----------------------------------------"

# Run the server
# Using +1 for US country code based on 724 area code
go run ./cmd/server -recipient="+17246807862"

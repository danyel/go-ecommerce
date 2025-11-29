#!/usr/bin/env bash
set -e

echo "📦 Running DevContainer setup..."

# Install Go dependencies
go mod tidy

# Install air
echo "📦 Installing Air..."
go install github.com/air-verse/air@latest

# Install goose
echo "📦 Installing Goose..."
go install github.com/pressly/goose/v3/cmd/goose@latest

# Load NVM
echo "📦 Installing Node.js via NVM..."
export NVM_DIR="/usr/local/nvm"
. $NVM_DIR/nvm.sh

nvm install --lts
npm install -g tailwindcss

echo "🎉 Dev container setup completed!"

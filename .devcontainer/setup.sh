#!/usr/bin/env bash
set -e

echo "📦 Running DevContainer setup..."

# Force kill any stubborn image-level GOROOT defaults for interactive bash sessions
if ! grep -q "export GOROOT=" /root/.bashrc; then
    echo "export GOROOT=/usr/local/go" >> /root/.bashrc
    echo "export PATH=/usr/local/go/bin:\$PATH" >> /root/.bashrc
fi

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
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

nvm install --lts
npm install -g tailwindcss

echo "🎉 Dev container setup completed!"

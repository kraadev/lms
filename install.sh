#!/usr/bin/env bash
# LMS One-Command VPS Installer (Linux)
set -e

echo "=============================================="
echo "🚀 LMS Platform Installer"
echo "=============================================="

if ! command -v docker &> /dev/null; then
    echo "[!] Docker not detected. Please install Docker first: https://docs.docker.com/engine/install/"
    exit 1
fi

if [ ! -f ".env" ]; then
    echo "[*] Creating .env from .env.example..."
    cp .env.example .env
fi

echo "[*] Starting LMS multi-container stack via Docker Compose..."
docker compose up -d --build

echo "=============================================="
echo "✨ LMS Platform is now live!"
echo "   Frontend: http://localhost:3000"
echo "   Backend:  http://localhost:8080"
echo "=============================================="

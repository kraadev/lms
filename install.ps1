# LMS One-Command Windows / PowerShell Installer
Write-Host "==============================================" -ForegroundColor Cyan
Write-Host "🚀 LMS Platform Installer" -ForegroundColor Cyan
Write-Host "==============================================" -ForegroundColor Cyan

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "[!] Docker not detected. Please install Docker Desktop: https://www.docker.com/products/docker-desktop" -ForegroundColor Red
    exit 1
}

if (-not (Test-Path ".env")) {
    Write-Host "[*] Creating .env from .env.example..." -ForegroundColor Yellow
    Copy-Item ".env.example" ".env"
}

Write-Host "[*] Launching containers via docker compose..." -ForegroundColor Green
docker compose up -d --build

Write-Host "==============================================" -ForegroundColor Cyan
Write-Host "✨ LMS Platform is now live!" -ForegroundColor Green
Write-Host "   Frontend: http://localhost:3000"
Write-Host "   Backend:  http://localhost:8080"
Write-Host "==============================================" -ForegroundColor Cyan

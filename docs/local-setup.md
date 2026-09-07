# Local Development & Setup Guide

Panduan teknis langkah demi langkah untuk menyiapkan, menjalankan, dan melakukan *debugging* aplikasi LMS di lingkungan pengembangan lokal.

---

## 📋 1. Kebutuhan Lingkungan

Pastikan perangkat Anda telah terinstal:
- **Go**: 1.21+ ([go.dev](https://go.dev/))
- **Node.js**: LTS 18.x atau 20.x ([nodejs.org](https://nodejs.org/))
- **npm**: 9.x+ (bawaan Node.js) atau `pnpm`
- **Git**: 2.30+

---

## 🗄️ 2. Konfigurasi Backend (Go)

1. Masuk ke direktori backend:
   ```bash
   cd backend
   ```
2. Salin template konfigurasi:
   ```bash
   cp .env.example .env
   ```
3. Unduh seluruh dependensi modul Go:
   ```bash
   go mod download
   ```
4. Jalankan seeder database (membuat schema dan pengguna bawaan):
   ```bash
   go run ./cmd/seed
   ```
5. Jalankan server backend:
   ```bash
   go run ./cmd/api
   ```
   > Server REST API siap di `http://localhost:8080` dan WebSocket endpoint di `ws://localhost:8080/ws`.

---

## 🖥️ 3. Konfigurasi Frontend (Nuxt 4)

1. Buka jendela terminal baru dan masuk ke direktori frontend:
   ```bash
   cd frontend
   ```
2. Salin template konfigurasi:
   ```bash
   cp .env.example .env
   ```
3. Pasang dependensi npm:
   ```bash
   npm install
   ```
4. Jalankan server development Nuxt 4:
   ```bash
   npm run dev
   ```
   > Buka peramban di `http://localhost:3000` untuk mengakses antarmuka LMS.

---

## 🛠️ 4. Pengujian Otomatis

- **Backend Tests**:
  ```bash
  cd backend
  go test -v ./internal/tests/...
  ```
- **Frontend Type Check**:
  ```bash
  cd frontend
  npx vue-tsc --noEmit
  ```

---

## ❓ 5. Troubleshooting Masalah Umum

- **Port 8080 atau 3000 sedang digunakan**: Ganti port pada berkas `.env` masing-masing modul.
- **Database Lock (SQLite)**: Pastikan tidak ada proses Go ganda yang sedang mengakses berkas `lms.db` secara bersamaan.

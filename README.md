<div align="center">

# 🎓 LMS — Learning Management System

Modern, scalable, and secure full-stack Learning Management System designed for seamless virtual learning and classroom collaboration.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Nuxt 4](https://img.shields.io/badge/Nuxt-4-00DC82?logo=nuxt.js&logoColor=white)](https://nuxt.com/)
[![Vue 3](https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![Tailwind CSS](https://img.shields.io/badge/TailwindCSS-3.4-38B2AC?logo=tailwind-css&logoColor=white)](https://tailwindcss.com/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

</div>

---

## 📖 Overview

**LMS (Learning Management System)** adalah platform pembelajaran digital berbasis web yang dirancang untuk menghubungkan pengajar (*Teachers*), peserta didik (*Students*), dan pengelola sistem (*Administrators*) dalam satu ekosistem terpadu.

Aplikasi ini menyederhanakan siklus pendidikan daring mulai dari pengelolaan silabus materi, pengumpulan dan penilaian tugas, kuis berbasis server yang aman, ruang diskusi real-time, hingga tatap muka interaktif berbasis video konferensi WebRTC berlatensi rendah.

---

## ✨ Key Features

- 👥 **Role-Based Access Control (RBAC)**  
  Tiga peran pengguna dengan batasan akses ketat (*server-side enforcement*): **Admin** (manajemen sistem/pengguna), **Teacher** (manajemen kelas, materi, penilaian), dan **Student** (pembelajaran, kuis, pengumpulan tugas).

- 📹 **Real-Time Video Classrooms & Meetings**  
  Tatap muka virtual langsung di dalam aplikasi ditenagai oleh **LiveKit WebRTC** dengan fallback penemuan peer ganda (REST + WebSocket), kontrol audio/video, dan pembagian hak moderator untuk pengajar.

- 💬 **Real-Time Classroom Chat**  
  Ruang obrolan interaktif per kelas berbasis **Gorilla WebSocket** dengan *in-memory connection hub*, pembersihan koneksi terputus otomatis, dan riwayat pesan persisten.

- 📚 **Classroom & Material Management**  
  Pembuatan ruang kelas dinamis, pembagian modul ajar (dokumen, teks, tautan), dan pengelolaan daftar anggota kelas.

- 📝 **Assignments & Automated Late Tracking**  
  Pengumpulan tugas digital dengan pendeteksian otomatis status keterlambatan (*late submission*), antarmuka penilaian, serta formulir masukan (*feedback*) dari pengajar.

- 🧠 **Zero-Leakage Quiz System**  
  Mekanisme pengerjaan kuis aman dengan proteksi server-side: kunci jawaban dan metadata penilaian tidak terekspos ke sisi klien hingga kuis selesai dievaluasi secara otomatis oleh sistem.

- 🔔 **In-App Notifications**  
  Pemberitahuan instan mengenai jadwal kuis, tugas baru yang ditugaskan, batas pengumpulan, nilai yang telah dirilis, dan jadwal kelas virtual.

- 🛡️ **Enterprise-Grade Security**  
  Penyimpanan sandi terenkripsi (bcrypt), otentikasi sesi ganda via JWT (HttpOnly Cookie & Bearer Header), sanitasi nama file penyimpanan UUID, dan pencegahan *path-traversal*.

---

## 🛠️ Tech Stack

### Frontend
- **Framework:** [Nuxt 4](https://nuxt.com/) (Vue 3 + Composition API + Nitro Engine)
- **Language:** [TypeScript](https://www.typescriptlang.org/)
- **Styling:** [Tailwind CSS](https://tailwindcss.com/)
- **State Management:** [Pinia](https://pinia.vuejs.org/)
- **WebRTC Client:** [LiveKit Client SDK](https://docs.livekit.io/)
- **Icons:** [Lucide Vue Next](https://lucide.dev/)

### Backend
- **Language & Runtime:** [Go (Golang)](https://go.dev/)
- **HTTP Routing:** [Chi Router (`go-chi/chi`)](https://github.com/go-chi/chi)
- **Realtime Protocols:** Gorilla WebSocket & LiveKit Go Protocol SDK
- **Database Engine:** PostgreSQL (Production default) dengan SQLite fallback untuk pengujian lokal
- **Authentication:** JWT (`golang-jwt/jwt/v5`) & `golang.org/x/crypto/bcrypt`

---

## 📋 Prerequisites

Sebelum memulai instalasi dan menjalankan proyek secara lokal, pastikan sistem Anda telah memenuhi prasyarat berikut:

- **Git**: Versi 2.30 atau lebih baru ([Unduh Git](https://git-scm.com/))
- **Node.js**: Versi LTS 18.x atau 20.x atau lebih baru ([Unduh Node.js](https://nodejs.org/))
- **npm**: Versi 9.x atau lebih baru (otomatis terpasang bersama Node.js) atau alternatif package manager seperti `pnpm` / `yarn`
- **Go**: Versi 1.21 atau lebih baru ([Unduh Go](https://go.dev/dl/))
- **Database Engine** *(Opsional)*: PostgreSQL 14+ (jika tidak dikonfigurasi, sistem otomatis menggunakan SQLite lokal sebagai fallback tanpa dependensi eksternal)

---

## 💻 Getting Started & Installation

Ikuti langkah-langkah berikut untuk mengkloning repositori, memasang dependensi, dan menjalankan lingkungan *development* lokal:

### 1. Kloning Repositori

```bash
# Clone repositori dari GitHub
git clone https://github.com/kraadev/lms.git

# Masuk ke direktori proyek
cd lms
```

### 2. Konfigurasi & Menjalankan Backend (Go)

Buka terminal pertama untuk backend:

```bash
# Masuk ke direktori backend
cd backend

# Salin berkas environment konfigurasi
cp .env.example .env

# Unduh seluruh dependensi Go modules
go mod download

# (Opsional) Jalankan migrasi dan database seeder akun bawaan
go run ./cmd/seed

# Jalankan server backend dalam mode development
go run ./cmd/api
```
> Backend API akan aktif di `http://localhost:8080` dan WebSocket di `ws://localhost:8080/ws`.

---

### 3. Konfigurasi & Menjalankan Frontend (Nuxt 4)

Buka terminal kedua untuk frontend:

```bash
# Masuk ke direktori frontend
cd frontend

# Salin berkas environment konfigurasi
cp .env.example .env

# Pasang dependensi modul npm
npm install

# Jalankan server development Nuxt 4 (Hot Module Replacement aktif)
npm run dev
```
> Antarmuka web frontend akan dapat diakses melalui peramban di `http://localhost:3000`.

---

### 4. Akun Pengujian Lokal (Seed Data)

Setelah seeder dijalankan, Anda dapat langsung masuk (*login*) menggunakan kredensial bawaan berikut:

| Peran (Role) | Email | Password | Hak Akses |
| :--- | :--- | :--- | :--- |
| **Admin** | `admin@lms.local` | `admin123` | Manajemen pengguna, kelas, dan konfigurasi sistem |
| **Teacher** | `teacher1@lms.local` | `password123` | Manajemen materi ajar, pembuatan kuis, dan penilaian tugas |
| **Student** | `student1@lms.local` | `password123` | Akses materi, pengerjaan kuis, dan pengumpulan tugas |

---

## 🤝 Contributing

Kontribusi selalu terbuka dan sangat kami hargai! Untuk berkontribusi:
1. *Fork* repositori ini
2. Buat *feature branch* (`git checkout -b feature/AmazingFeature`)
3. *Commit* perubahan Anda (`git commit -m 'feat: Add AmazingFeature'`)
4. *Push* ke branch Anda (`git push origin feature/AmazingFeature`)
5. Buka **Pull Request**

---

## 📄 License

Didistribusikan di bawah Lisensi MIT. Lihat file `LICENSE` untuk informasi lebih lanjut.

---

## 📚 Documentation

Untuk panduan teknis yang lebih terperinci, silakan merujuk ke dokumentasi kami:
- [Architecture Overview](docs/overview.md)
- [Local Development & Setup Guide](docs/local-setup.md)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)

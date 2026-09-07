<div align="center">

# 🎓 LMS — Learning Management System

Modern, scalable, and secure full-stack Learning Management System designed for seamless virtual learning and classroom collaboration.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Nuxt 3](https://img.shields.io/badge/Nuxt-3-00DC82?logo=nuxt.js&logoColor=white)](https://nuxt.com/)
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
- **Framework:** [Nuxt 3](https://nuxt.com/) (Vue 3 + Composition API + Nitro Engine)
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

## 🚀 Quick Start

### 1. Prasyarat
- Go 1.21 atau lebih baru
- Node.js 18.x atau lebih baru (npm / pnpm / yarn)
- PostgreSQL (Opsional, fallback otomatis ke SQLite lokal)

### 2. Setup Backend
```bash
cd backend
cp .env.example .env
go run ./cmd/seed    # Jalankan database migration & seed awal
go run ./cmd/api     # Server aktif di http://localhost:8080
```

### 3. Setup Frontend
```bash
cd frontend
cp .env.example .env
npm install
npm run dev          # Aplikasi aktif di http://localhost:3000
```

### 4. Akun Bawaan Pengujian (Seed Data)
| Role | Email | Password |
| :--- | :--- | :--- |
| **Admin** | `admin@lms.local` | `admin123` |
| **Teacher** | `teacher1@lms.local` | `password123` |
| **Student** | `student1@lms.local` | `password123` |

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

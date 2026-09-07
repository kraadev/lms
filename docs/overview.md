# LMS System Architecture Overview

Dokumen ini menjelaskan arsitektur tingkat tinggi dan komponen utama dari platform **Learning Management System (LMS)**.

---

## 🏛️ Arsitektur Keseluruhan

Platform LMS dirancang menggunakan pola **Modular Monolith / Decoupled Client-Server**:
- **Frontend Layer**: Nuxt 4 (Vue 3 + Composition API + Nitro Engine) menyajikan Single Page Application (SPA) responsif dengan rendering optimal.
- **Backend API Layer**: Go (Golang) menggunakan router `chi/v5` dengan struktur *package-by-feature* modular.
- **Real-Time Communication**:
  - **Gorilla WebSocket**: Digunakan untuk chat room kelas instan, sinkronisasi state pesan, dan notifikasi live.
  - **LiveKit WebRTC**: Diterapkan untuk ruang tatap muka virtual audio/video real-time berlatensi ultra-rendah.
- **Persistence Layer**: PostgreSQL sebagai basis data utama production dengan fallback otomatis ke SQLite lokal saat pengujian/development.

---

## 📂 Struktur Modul Backend

Backend dibagi secara modular di dalam direktori `backend/internal/`:
- `auth`: Otentikasi pengguna, pembuatan JWT token, dan hashing bcrypt.
- `users`: Manajemen data pengguna dan profil.
- `classes`: Pengelolaan kelas, ruang ajar, dan keanggotaan mahasiswa (*membership*).
- `materials`: Distribusi silabus, dokumen materi, dan tautan belajar.
- `assignments`: Pengumpulan tugas, kalkulasi keterlambatan (*late submission*), dan penilaian pengajar.
- `quizzes`: Mesin kuis berbasis server dengan arsitektur *Zero-Leakage* (kunci jawaban tidak pernah bocor ke klien sebelum dievaluasi server).
- `meetings`: Manajemen ruang meeting virtual dan token generator LiveKit.
- `chat`: WebSocket hub terisolasi per ruang kelas dengan *thread-safe client management*.
- `notifications`: Sistem event notifikasi in-app untuk tugas, kuis, dan meeting.

---

## 💻 Struktur Modul Frontend

Frontend menggunakan Nuxt 4 dengan struktur direktori idiomatik:
- `pages/`: Rute berbasis file untuk dashboard, kelas, kuis, tugas, meeting, dan panel admin.
- `components/`: Komponen UI modular (tabs, badges, modals, avatars) dan komponen fitur per modul.
- `stores/`: Pengelolaan state global menggunakan Pinia (`auth`, `meeting`, `notifications`).
- `composables/`: Abstraksi logika reaktif (`useWebSocket`, `useLiveKit`, `useToast`, `useTheme`).
- `services/`: Lapisan abstraksi panggilan REST API ke backend.

---

## 🔐 Model Keamanan & Otorisasi

1. **Role-Based Access Control (RBAC)**: Validasi tiga tingkat hak akses: Admin, Teacher, dan Student.
2. **Server-Side Verification**: Seluruh aksi (pengumpulan tugas, pembukaan meeting, penilaian kuis) divalidasi mutlak di server, bukan hanya di UI klien.
3. **JWT Dual Mode**: Token disimpan aman di HttpOnly cookie serta didukung melalui header `Authorization: Bearer <token>`.

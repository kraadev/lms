# LMS Backend (Go + PostgreSQL + LiveKit + WebSocket)

Backend lokal production-ready untuk platform Learning Management System (LMS) modern, dibangun menggunakan arsitektur modular Go yang tangguh dan aman.

## 🚀 Fitur Utama

- **Role-Based Access Control (RBAC)**: Admin, Teacher, Student.
- **Strict Server-Side Access Enforcement**: Validasi kepemilikan kelas, keanggotaan mahasiswa (*membership*), materi, kuis, tugas, dan ruangan meeting diverifikasi di server.
- **Autentikasi Aman**: Hash password bcrypt, JWT token via HttpOnly Cookie (`lms_token`) dan header `Authorization: Bearer <token>`.
- **Database Migrations & Seeding Otomatis**: Mendukung PostgreSQL (default) serta fallback SQLite untuk local execution cepat tanpa ketergantungan Docker.
- **Realtime WebSocket Chat (`ws://localhost:8080/ws`)**: Concurrency-safe in-memory room hub, message validation, dead connection cleanup, dan pesan tersimpan ke PostgreSQL.
- **WebRTC Video & Audio Rooms (LiveKit)**: Pembuatan token room LiveKit dengan pembagian hak akses (admin room untuk guru, partisipan untuk siswa) dan pencegahan akses ke meeting yang telah berakhir (*ended*).
- **Sistem Kuis Anti-Bocor (Zero-Leakage)**: Kunci jawaban dan grading metadata tidak pernah diekspos ke siswa sebelum kuis selesai dikerjakan dan dinilai secara otomatis di server.
- **Tugas & Penilaian**: Unggah tugas/teks, kalkulasi otomatis keterlambatan (*late submission*), dan grading dengan feedback.
- **File Storage Aman**: Sanitasi nama file unik UUID dan proteksi path traversal `../`.
- **Sistem Notifikasi**: Notifikasi event tugas baru, penilaian, kuis, dan meeting.

---

## 🛠️ Prasyarat & Menjalankan Lokal

### 1. Konfigurasi Environment (`.env`)
Salin file `.env.example` menjadi `.env` jika belum ada:
```bash
cp .env.example .env
```

### 2. Menjalankan Backend
Jalankan server langsung dengan Go:
```bash
go run ./cmd/api
```
Atau cukup:
```bash
go run .
```

Server akan aktif di:
- **API Server**: `http://localhost:8080`
- **WebSocket Endpoint**: `ws://localhost:8080/ws`
- **Health Check**: `http://localhost:8080/health`

### 3. Menjalankan Database Seeder (Opsional / Terpisah)
```bash
go run ./cmd/seed
```

---

## 👥 Akun Demo / Seed Data

| Role | Nama | Email | Password |
| :--- | :--- | :--- | :--- |
| **Admin** | System Administrator | `admin@lms.local` | `admin123` |
| **Teacher** | Budi Santoso | `teacher1@lms.local` | `password123` |
| **Teacher** | Siti Rahma | `teacher2@lms.local` | `password123` |
| **Student** | Andi Pratama (Kelas 1) | `student1@lms.local` | `password123` |
| **Student** | Bunga Citra (Kelas 1) | `student2@lms.local` | `password123` |
| **Student** | Candra Wijaya (Kelas 2) | `student3@lms.local` | `password123` |
| **Student** | Dewi Lestari (Kelas 2) | `student4@lms.local` | `password123` |

---

## 🧪 Menjalankan Automated Tests

Jalankan rangkaian unit & integration tests:
```bash
go test -v ./internal/tests/...
```

---

## 📡 Daftar REST Endpoints Utama

### Auth
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/me`

### Classes
- `GET /api/classes`
- `POST /api/classes`
- `GET /api/classes/:id`
- `PATCH /api/classes/:id`
- `DELETE /api/classes/:id`
- `GET /api/classes/:id/members`
- `POST /api/classes/:id/members`
- `DELETE /api/classes/:id/members/:userId`

### Materials & Assignments
- `GET /api/classes/:classId/materials`
- `POST /api/classes/:classId/materials`
- `GET /api/classes/:classId/assignments`
- `POST /api/classes/:classId/assignments`
- `POST /api/assignments/:id/submissions`
- `GET /api/assignments/:id/my-submission`
- `PATCH /api/submissions/:id/grade`

### Quizzes
- `GET /api/classes/:classId/quizzes`
- `POST /api/classes/:classId/quizzes`
- `GET /api/quizzes/:id`
- `POST /api/quizzes/:id/attempts`
- `POST /api/quizzes/attempts/:id/submit`
- `GET /api/quizzes/attempts/:id`

### Video & Audio Meetings (LiveKit)
- `POST /api/classes/:classId/meetings`
- `GET /api/classes/:classId/meetings`
- `POST /api/meetings/:id/join`
- `POST /api/meetings/:id/end`

### Realtime Chat & History
- `ws://localhost:8080/ws` (WebSocket connection)
- `GET /api/classes/:classId/messages?limit=50&offset=0`

### Notifications
- `GET /api/notifications`
- `PATCH /api/notifications/:id/read`
- `PATCH /api/notifications/read-all`

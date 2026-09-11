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

---

## 🚀 Fase Pengembangan 2: Komponen Dasar & Data Stubs

Pada fase ini, fondasi struktur kode aplikasi diperkuat dengan menambahkan modul `src/` awal:
- **`src/utils/helpers.js`**: Utilitas bersama pemformatan tanggal, teks, dan mata uang.
- **`src/data/courses.json`**: Dataset kursus mock untuk pengujian tata letak katalog dan ruang ajar.
- **`src/components/Navbar.jsx`**: Kerangka komponen navigasi utama dan kontrol status autentikasi pengguna.
- **`src/constants/index.js`**: Registry konstanta peran (*roles*), status pertemuan (*meetings*), dan rute navigasi.
- **`src/routes/index.js`**: Pemetaan rute aplikasi terpusat dengan dukungan metadata hak akses.
- **`src/styles/global.css`**: Variabel token desain global untuk konsistensi UI antarmuka.

---

## 🎨 Fase Pengembangan 3: Design System & Semantic Tokens

Pada fase ini, sistem desain antarmuka diperkuat dengan komponen primitif modular dan token semantik:
- **Primitive UI Components**:
  - `UiCard.vue`: Kontainer kartu dengan varian `default`, `flat`, `outline`, dan `interactive` beserta header/footer slots.
  - `UiProgress.vue`: Indikator progres dengan varian semantik (`brand`, `success`, `warning`, `danger`), label terintegrasi, dan animasi pulsa.
  - `UiSelect.vue`: Input seleksi formulir aksesibel dengan opsi dinamis, status validasi error, dan adaptasi tema.
  - `UiDropdown.vue`: Menu popover mengambang dengan transisi animasi halus, deteksi *click-outside*, dan navigasi keyboard Esc.
- **Semantic Design Tokens**:
  - Palet warna semantik `brand`, `surface`, dan `status` (`success`, `warning`, `danger`, `info`) pada `tailwind.preset.js` dan konfigurasi Nuxt.
  - Variabel CSS adaptif di `frontend/assets/css/main.css` dan `src/styles/global.css` yang memenuhi standar kontras WCAG AA untuk light dan dark mode.

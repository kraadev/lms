# Contributing to LMS

Terima kasih atas minat Anda untuk berkontribusi pada proyek **LMS (Learning Management System)**! Kami menyambut segala bentuk kontribusi dari komunitas open source, baik berupa perbaikan bug, penambahan fitur baru, perbaikan dokumentasi, maupun pelaporan isu.

---

## 📜 Kode Etik (Code of Conduct)

Dengan berpartisipasi dalam proyek ini, Anda diharapkan menjunjung tinggi komunikasi yang profesional, inklusif, dan saling menghargai. Silakan merujuk ke berkas [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) untuk detail lengkap.

---

## 🔄 Alur Kontribusi (Pull Request Workflow)

Untuk menjaga riwayat commit yang rapi dan terverifikasi, setiap kontribusi dilakukan melalui mekanisme Pull Request (PR):

1. **Fork Repositori**: Buat salinan (*fork*) repositori ke akun GitHub Anda.
2. **Kloning Repositori**:
   ```bash
   git clone https://github.com/<username-anda>/lms.git
   cd lms
   ```
3. **Pastikan Branch `main` Sinkron**:
   ```bash
   git checkout main
   git pull origin main
   ```
4. **Buat Feature Branch**:
   Gunakan konvensi penamaan branch sesuai jenis perubahan:
   - `feat/<nama-fitur>` — Penambahan fitur baru
   - `fix/<nama-bug>` — Perbaikan bug / error
   - `docs/<deskripsi>` — Pembaruan dokumentasi
   - `refactor/<nama>` — Restrukturisasi kode tanpa mengubah fungsionalitas
   - `chore/<deskripsi>` — Pemeliharaan dependensi / konfigurasi

   Contoh:
   ```bash
   git checkout -b feat/quiz-timer
   ```

5. **Lakukan Perubahan & Pengujian**:
   Pastikan kode mematuhi standar proyek dan tidak menimbulkan *breaking changes*.

6. **Commit Perubahan**:
   Gunakan format Conventional Commits:
   ```bash
   git commit -m "feat(quizzes): add countdown timer for timed quiz sessions"
   ```

7. **Push ke Branch Anda**:
   ```bash
   git push -u origin feat/quiz-timer
   ```

8. **Buat Pull Request**:
   - Buka Pull Request ke branch `main` repositori `kraadev/lms`.
   - Berikan deskripsi yang jelas mengenai masalah yang diselesaikan dan perubahan yang dibuat.
   - Sertakan tangkapan layar (*screenshot*) jika berkaitan dengan antarmuka pengguna (UI).

---

## 🛠️ Standar Kode & Konvensi

- **Backend (Go)**:
  - Format kode menggunakan `gofmt` atau `goimports`.
  - Jalankan pengujian unit sebelum membuka PR: `go test -v ./...`.
  - Tangani setiap potensi error secara eksplisit.

- **Frontend (Nuxt 4 / Vue 3)**:
  - Gunakan Vue 3 Composition API dengan `<script setup lang="ts">`.
  - Gunakan Tailwind CSS untuk styling utilitas.
  - Pastikan tidak ada error TypeScript: `npm run build` atau `npx vue-tsc --noEmit`.

---

## 🐛 Melaporkan Bug & Mengajukan Fitur

Jika Anda menemukan kendala atau memiliki ide pengembangan:
1. Periksa tab [Issues](https://github.com/kraadev/lms/issues) untuk memastikan topik belum pernah dibahas sebelumnya.
2. Buat Issue baru dengan memilih template yang sesuai dan isi form secara lengkap.

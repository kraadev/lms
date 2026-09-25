# 🚀 LMS Platform v1.0.0 Release Notes

Kami dengan bangga merilis versi **v1.0.0** platform Learning Management System (LMS) modular berskala penuh!

## 🌟 Fitur Utama Rilis v1.0.0:
1. **Frontend Modern (Nuxt 4 + Vue 3):**
   - Dribbble-quality UI dengan token warna semantik (`brand`, `surface`, `status`).
   - Dark/Light theme switcher dengan persistence lokal.
   - Interactive Speed-Grader drawer dan Zen Exam Mode.
2. **Backend Berperforma Tinggi (Go Chi v5):**
   - Zero-Leakage quiz evaluation engine.
   - Dynamic Role-Based Access Control (Admin, Teacher, Student).
   - WebSocket Chat Hub terisolasi dan LiveKit video conferencing.
3. **Infrastruktur & Keamanan:**
   - Proteksi OWASP security headers, brute-force rate limiter, dan audit trail log.
   - Kubernetes liveness/readiness probes (`/healthz`, `/readyz`).
   - Orkestrasi Docker Compose multi-container satu perintah (`install.sh`, `install.ps1`).

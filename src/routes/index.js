/**
 * Basic Application Route Configuration Registry
 * Mapping page routes to corresponding views and access restrictions.
 */

import { USER_ROLES } from '../constants/index.js'

export const routes = [
  {
    path: '/',
    name: 'Home',
    meta: { public: true, title: 'Selamat Datang di LMS' }
  },
  {
    path: '/login',
    name: 'Login',
    meta: { guestOnly: true, title: 'Masuk ke Akun LMS' }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    meta: { requiresAuth: true, title: 'Ringkasan Belajar' }
  },
  {
    path: '/classes',
    name: 'ClassList',
    meta: { requiresAuth: true, title: 'Daftar Kelas Pembelajaran' }
  },
  {
    path: '/classes/:id',
    name: 'ClassDetail',
    meta: { requiresAuth: true, title: 'Detail Ruang Kelas' }
  },
  {
    path: '/assignments',
    name: 'Assignments',
    meta: { requiresAuth: true, title: 'Daftar Tugas & Pengumpulan' }
  },
  {
    path: '/quizzes',
    name: 'Quizzes',
    meta: { requiresAuth: true, title: 'Evaluasi Pembelajaran' }
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    meta: { requiresAuth: true, allowedRoles: [USER_ROLES.ADMIN], title: 'Manajemen Pengguna' }
  }
]

export default routes

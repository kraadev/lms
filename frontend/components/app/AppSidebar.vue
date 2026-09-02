<script setup lang="ts">
import {
  LayoutDashboard, BookOpen, ClipboardList, FileQuestion,
  Video, MessageSquare, Bell, Users, GraduationCap, Settings,
  ShieldCheck, BookMarked, Home
} from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'

const auth = useAuthStore()

interface NavItem {
  label: string
  href: string
  icon: any
  exact?: boolean
}

const studentNav: NavItem[] = [
  { label: 'Dashboard', href: '/dashboard', icon: Home, exact: true },
  { label: 'Kelas Saya', href: '/classes', icon: BookOpen },
  { label: 'Tugas', href: '/assignments', icon: ClipboardList },
  { label: 'Kuis', href: '/quizzes', icon: FileQuestion },
  { label: 'Pesan', href: '/messages', icon: MessageSquare },
  { label: 'Notifikasi', href: '/notifications', icon: Bell }
]

const teacherNav: NavItem[] = [
  { label: 'Dashboard', href: '/dashboard', icon: Home, exact: true },
  { label: 'Kelas Saya', href: '/classes', icon: BookOpen },
  { label: 'Tugas', href: '/assignments', icon: ClipboardList },
  { label: 'Kuis', href: '/quizzes', icon: FileQuestion },
  { label: 'Meeting', href: '/meetings', icon: Video },
  { label: 'Pesan', href: '/messages', icon: MessageSquare },
  { label: 'Siswa', href: '/students', icon: GraduationCap },
  { label: 'Notifikasi', href: '/notifications', icon: Bell }
]

const adminNav: NavItem[] = [
  { label: 'Dashboard', href: '/dashboard', icon: LayoutDashboard, exact: true },
  { label: 'Pengguna', href: '/admin/users', icon: Users },
  { label: 'Kelas', href: '/admin/classes', icon: BookMarked },
  { label: 'Notifikasi', href: '/notifications', icon: Bell }
]

const navItems = computed((): NavItem[] => {
  if (auth.isAdmin) return adminNav
  if (auth.isTeacher) return teacherNav
  return studentNav
})

const route = useRoute()

function isActive(item: NavItem): boolean {
  if (item.exact) return route.path === item.href
  return route.path.startsWith(item.href)
}
</script>

<template>
  <nav class="flex flex-col gap-0.5" aria-label="Main navigation">
    <NuxtLink
      v-for="item in navItems"
      :key="item.href"
      :to="item.href"
      :class="[
        'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all group',
        isActive(item)
          ? 'bg-brand-50 dark:bg-brand-950/50 text-brand-700 dark:text-brand-300'
          : 'text-surface-600 dark:text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800/70 hover:text-surface-900 dark:hover:text-surface-100'
      ]"
    >
      <component
        :is="item.icon"
        :class="[
          'w-4.5 h-4.5 shrink-0 transition-transform group-hover:scale-105',
          isActive(item) ? 'text-brand-600 dark:text-brand-400' : 'text-surface-400 dark:text-surface-500'
        ]"
      />
      {{ item.label }}
    </NuxtLink>
  </nav>
</template>

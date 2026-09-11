<script setup lang="ts">
import {
  LayoutDashboard, BookOpen, ClipboardList, FileQuestion,
  Video, MessageSquare, Bell, Users, GraduationCap,
  BookMarked, Home
} from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useNotificationsStore } from '~/stores/notifications'

const auth = useAuthStore()
const notifStore = useNotificationsStore()

interface NavItem {
  label: string
  href: string
  icon: any
  exact?: boolean
  badge?: string | number
}

interface NavSection {
  title: string
  items: NavItem[]
}

const studentSections = computed<NavSection[]>(() => [
  {
    title: 'Utama',
    items: [
      { label: 'Dashboard', href: '/dashboard', icon: Home, exact: true }
    ]
  },
  {
    title: 'Akademik',
    items: [
      { label: 'Kelas Saya', href: '/classes', icon: BookOpen },
      { label: 'Tugas', href: '/assignments', icon: ClipboardList },
      { label: 'Kuis', href: '/quizzes', icon: FileQuestion }
    ]
  },
  {
    title: 'Komunikasi',
    items: [
      { label: 'Pesan', href: '/messages', icon: MessageSquare },
      {
        label: 'Notifikasi',
        href: '/notifications',
        icon: Bell,
        badge: notifStore.unreadCount > 0 ? notifStore.unreadCount : undefined
      }
    ]
  }
])

const teacherSections = computed<NavSection[]>(() => [
  {
    title: 'Utama',
    items: [
      { label: 'Dashboard', href: '/dashboard', icon: Home, exact: true }
    ]
  },
  {
    title: 'Akademik',
    items: [
      { label: 'Kelas Saya', href: '/classes', icon: BookOpen },
      { label: 'Tugas', href: '/assignments', icon: ClipboardList },
      { label: 'Kuis', href: '/quizzes', icon: FileQuestion },
      { label: 'Siswa', href: '/students', icon: GraduationCap }
    ]
  },
  {
    title: 'Kolaborasi',
    items: [
      { label: 'Meeting', href: '/meetings', icon: Video },
      { label: 'Pesan', href: '/messages', icon: MessageSquare },
      {
        label: 'Notifikasi',
        href: '/notifications',
        icon: Bell,
        badge: notifStore.unreadCount > 0 ? notifStore.unreadCount : undefined
      }
    ]
  }
])

const adminSections = computed<NavSection[]>(() => [
  {
    title: 'Utama',
    items: [
      { label: 'Dashboard', href: '/dashboard', icon: LayoutDashboard, exact: true }
    ]
  },
  {
    title: 'Manajemen Sistem',
    items: [
      { label: 'Pengguna', href: '/admin/users', icon: Users },
      { label: 'Kelas', href: '/admin/classes', icon: BookMarked }
    ]
  },
  {
    title: 'Komunikasi',
    items: [
      {
        label: 'Notifikasi',
        href: '/notifications',
        icon: Bell,
        badge: notifStore.unreadCount > 0 ? notifStore.unreadCount : undefined
      }
    ]
  }
])

const navSections = computed((): NavSection[] => {
  if (auth.isAdmin) return adminSections.value
  if (auth.isTeacher) return teacherSections.value
  return studentSections.value
})

const route = useRoute()

function isActive(item: NavItem): boolean {
  if (item.exact) return route.path === item.href
  return route.path.startsWith(item.href)
}
</script>

<template>
  <nav class="flex flex-col gap-4" aria-label="Main navigation">
    <div v-for="section in navSections" :key="section.title" class="flex flex-col gap-0.5">
      <!-- Section Header -->
      <h4 class="text-[10px] font-bold uppercase tracking-wider text-surface-400 dark:text-surface-500 px-3 py-1 select-none">
        {{ section.title }}
      </h4>

      <!-- Section Nav Items -->
      <NuxtLink
        v-for="item in section.items"
        :key="item.href"
        :to="item.href"
        :class="[
          'flex items-center justify-between px-3 py-2 rounded-xl text-xs font-semibold transition-all group relative',
          isActive(item)
            ? 'bg-brand-50 dark:bg-brand-950/60 text-brand-600 dark:text-brand-300 shadow-soft'
            : 'text-surface-600 dark:text-surface-400 hover:bg-surface-100/80 dark:hover:bg-surface-800/60 hover:text-surface-900 dark:hover:text-surface-100'
        ]"
      >
        <div class="flex items-center gap-2.5 min-w-0">
          <component
            :is="item.icon"
            :class="[
              'w-4 h-4 shrink-0 transition-transform group-hover:scale-110',
              isActive(item) ? 'text-brand-600 dark:text-brand-400' : 'text-surface-400 dark:text-surface-500'
            ]"
          />
          <span class="truncate">{{ item.label }}</span>
        </div>

        <!-- Optional Unread Badge -->
        <span
          v-if="item.badge"
          class="flex items-center justify-center min-w-[18px] h-4 px-1 rounded-full text-[10px] font-bold bg-rose-500 text-white shrink-0 ml-1"
        >
          {{ Number(item.badge) > 9 ? '9+' : item.badge }}
        </span>
      </NuxtLink>
    </div>
  </nav>
</template>

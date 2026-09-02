<script setup lang="ts">
import { Bell, Sun, Moon, Monitor, ChevronDown, LogOut, User, Settings } from 'lucide-vue-next'
import { useAuthStore } from '~/stores/auth'
import { useNotificationsStore } from '~/stores/notifications'
import type { ThemeMode } from '~/composables/useTheme'

const auth = useAuthStore()
const notifStore = useNotificationsStore()
const { themeMode, applyTheme } = useTheme()

const showUserMenu = ref(false)
const showThemeMenu = ref(false)
const showNotifPanel = ref(false)

const themes: { key: ThemeMode; label: string; icon: any }[] = [
  { key: 'light', label: 'Terang', icon: Sun },
  { key: 'dark', label: 'Gelap', icon: Moon },
  { key: 'system', label: 'Sistem', icon: Monitor }
]

const currentThemeIcon = computed(() => {
  if (themeMode.value === 'dark') return Moon
  if (themeMode.value === 'light') return Sun
  return Monitor
})

function closeAll() {
  showUserMenu.value = false
  showThemeMenu.value = false
  showNotifPanel.value = false
}

async function handleLogout() {
  closeAll()
  await auth.logout()
}

// Click outside directive
function useClickOutside(el: Ref<HTMLElement | null>, cb: () => void) {
  function onClick(e: MouseEvent) {
    if (el.value && !el.value.contains(e.target as Node)) cb()
  }
  onMounted(() => document.addEventListener('click', onClick))
  onUnmounted(() => document.removeEventListener('click', onClick))
}
</script>

<template>
  <header class="h-14 flex items-center gap-3 px-4 lg:px-6 border-b border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-950 shrink-0">
    <!-- Slot for mobile menu button (injected from layout) -->
    <slot name="mobile-trigger" />

    <!-- Page title -->
    <div class="flex-1 min-w-0">
      <slot />
    </div>

    <div class="flex items-center gap-1.5">
      <!-- Theme toggle -->
      <div class="relative">
        <button
          type="button"
          class="p-2 rounded-lg text-surface-500 dark:text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800 hover:text-surface-800 dark:hover:text-surface-100 transition-colors"
          :aria-label="`Tema saat ini: ${themeMode}`"
          @click="showThemeMenu = !showThemeMenu; showUserMenu = false; showNotifPanel = false"
        >
          <component :is="currentThemeIcon" class="w-4.5 h-4.5" />
        </button>
        <Transition enter-active-class="transition duration-100 ease-out" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100" leave-active-class="transition duration-75 ease-in" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95">
          <div v-if="showThemeMenu" class="absolute right-0 mt-1 top-full w-36 bg-white dark:bg-surface-900 rounded-xl shadow-elevated border border-surface-200 dark:border-surface-800 py-1 z-30">
            <button
              v-for="theme in themes"
              :key="theme.key"
              type="button"
              :class="['w-full flex items-center gap-2.5 px-3 py-2 text-sm transition-colors', themeMode === theme.key ? 'text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-950/50' : 'text-surface-700 dark:text-surface-300 hover:bg-surface-50 dark:hover:bg-surface-800']"
              @click="applyTheme(theme.key); showThemeMenu = false"
            >
              <component :is="theme.icon" class="w-4 h-4 shrink-0" />
              {{ theme.label }}
            </button>
          </div>
        </Transition>
      </div>

      <!-- Notifications -->
      <div class="relative">
        <button
          type="button"
          class="relative p-2 rounded-lg text-surface-500 dark:text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800 hover:text-surface-800 dark:hover:text-surface-100 transition-colors"
          aria-label="Notifikasi"
          @click="showNotifPanel = !showNotifPanel; showUserMenu = false; showThemeMenu = false; if (showNotifPanel) notifStore.fetchNotifications()"
        >
          <Bell class="w-4.5 h-4.5" />
          <span
            v-if="notifStore.unreadCount > 0"
            class="absolute top-1 right-1 flex items-center justify-center min-w-[16px] h-4 px-1 rounded-full text-[9px] font-bold bg-rose-500 text-white"
          >{{ notifStore.unreadCount > 9 ? '9+' : notifStore.unreadCount }}</span>
        </button>

        <Transition enter-active-class="transition duration-100 ease-out" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100" leave-active-class="transition duration-75 ease-in" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95">
          <div v-if="showNotifPanel" class="absolute right-0 top-full mt-1 w-80 bg-white dark:bg-surface-900 rounded-xl shadow-elevated border border-surface-200 dark:border-surface-800 z-30 overflow-hidden">
            <div class="flex items-center justify-between px-4 py-3 border-b border-surface-100 dark:border-surface-800">
              <span class="text-sm font-semibold text-surface-900 dark:text-surface-100">Notifikasi</span>
              <button v-if="notifStore.unreadCount > 0" type="button" class="text-xs text-brand-600 dark:text-brand-400 hover:underline" @click="notifStore.markAllAsRead()">Tandai semua dibaca</button>
            </div>
            <div class="max-h-80 overflow-y-auto divide-y divide-surface-100 dark:divide-surface-800">
              <div v-if="notifStore.notifications.length === 0" class="py-10 text-center text-sm text-surface-400">Tidak ada notifikasi</div>
              <NuxtLink
                v-for="notif in notifStore.notifications.slice(0, 8)"
                :key="notif.id"
                :to="notif.link || '/notifications'"
                :class="['flex items-start gap-3 px-4 py-3 hover:bg-surface-50 dark:hover:bg-surface-800/60 transition-colors', !notif.is_read ? 'bg-brand-50/40 dark:bg-brand-950/20' : '']"
                @click="notifStore.markAsRead(notif.id); showNotifPanel = false"
              >
                <span :class="['w-2 h-2 rounded-full shrink-0 mt-1.5', notif.is_read ? 'bg-transparent' : 'bg-brand-500']" />
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-surface-800 dark:text-surface-200 truncate">{{ notif.title }}</p>
                  <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5 line-clamp-2">{{ notif.message }}</p>
                </div>
              </NuxtLink>
            </div>
            <div class="border-t border-surface-100 dark:border-surface-800">
              <NuxtLink to="/notifications" class="flex items-center justify-center py-2.5 text-sm text-brand-600 dark:text-brand-400 hover:underline" @click="showNotifPanel = false">Lihat semua notifikasi</NuxtLink>
            </div>
          </div>
        </Transition>
      </div>

      <!-- User menu -->
      <div class="relative">
        <button
          type="button"
          class="flex items-center gap-2 pl-2 pr-2 py-1.5 rounded-lg text-sm font-medium hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
          @click="showUserMenu = !showUserMenu; showThemeMenu = false; showNotifPanel = false"
        >
          <UiAvatar :name="auth.user?.name" size="sm" />
          <span class="hidden sm:block text-surface-700 dark:text-surface-300 max-w-[120px] truncate">{{ auth.user?.name }}</span>
          <ChevronDown class="w-4 h-4 text-surface-400 shrink-0" />
        </button>

        <Transition enter-active-class="transition duration-100 ease-out" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100" leave-active-class="transition duration-75 ease-in" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95">
          <div v-if="showUserMenu" class="absolute right-0 top-full mt-1 w-52 bg-white dark:bg-surface-900 rounded-xl shadow-elevated border border-surface-200 dark:border-surface-800 py-1.5 z-30">
            <div class="px-3 py-2 border-b border-surface-100 dark:border-surface-800 mb-1">
              <p class="text-sm font-medium text-surface-900 dark:text-surface-100 truncate">{{ auth.user?.name }}</p>
              <p class="text-xs text-surface-500 dark:text-surface-400 truncate">{{ auth.user?.email }}</p>
            </div>
            <button type="button" class="w-full flex items-center gap-2.5 px-3 py-2 text-sm text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors" @click="handleLogout">
              <LogOut class="w-4 h-4" />
              Keluar
            </button>
          </div>
        </Transition>
      </div>
    </div>
  </header>
</template>

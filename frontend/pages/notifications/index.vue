<script setup lang="ts">
import { Bell, CheckCheck, BookOpen, ClipboardList, FileQuestion, Video, MessageSquare } from 'lucide-vue-next'
import { notificationsService } from '~/services/notifications'
import type { Notification } from '~/types'
import { formatRelativeTime } from '~/utils/formatters'

definePageMeta({ middleware: 'auth' })
useSeoMeta({ title: 'Notifikasi' })

const notifStore = useNotificationsStore()
const isLoading = ref(true)
const isMarkingAll = ref(false)

onMounted(async () => {
  isLoading.value = true
  await notifStore.fetchNotifications()
  isLoading.value = false
})

async function handleMarkAllAsRead() {
  isMarkingAll.value = true
  await notifStore.markAllAsRead()
  isMarkingAll.value = false
}

async function handleClickNotif(notif: Notification) {
  if (!notif.read) {
    await notifStore.markAsRead(notif.id)
  }
  if (notif.link) {
    navigateTo(notif.link)
  }
}

function getIcon(type: string) {
  switch (type) {
    case 'assignment': return ClipboardList
    case 'quiz': return FileQuestion
    case 'meeting': return Video
    case 'announcement': return BookOpen
    case 'chat': return MessageSquare
    default: return Bell
  }
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-3xl mx-auto space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">Notifikasi</h1>
        <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Pembaruan tugas, pengumuman, dan aktivitas kelas.</p>
      </div>

      <UiButton
        v-if="notifStore.unreadCount > 0"
        variant="ghost"
        size="xs"
        :loading="isMarkingAll"
        @click="handleMarkAllAsRead"
      >
        <CheckCheck class="w-4 h-4 mr-1 text-brand-600" />
        Tandai Semua Dibaca
      </UiButton>
    </div>

    <div v-if="isLoading">
      <UiSkeleton :rows="5" />
    </div>

    <UiEmptyState
      v-else-if="!notifStore.notifications.length"
      :icon="Bell"
      title="Belum ada notifikasi"
      description="Semua kabar terbaru seputar aktivitas Anda akan muncul di sini."
    />

    <div v-else class="divide-y divide-surface-100 dark:divide-surface-800 rounded-2xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft">
      <div
        v-for="notif in notifStore.notifications"
        :key="notif.id"
        class="p-4 flex items-start gap-3.5 hover:bg-surface-50 dark:hover:bg-surface-800/50 cursor-pointer transition-colors"
        :class="{ 'bg-brand-50/30 dark:bg-brand-950/20': !notif.read }"
        @click="handleClickNotif(notif)"
      >
        <div
          class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
          :class="notif.read
            ? 'bg-surface-100 dark:bg-surface-800 text-surface-500'
            : 'bg-brand-100 dark:bg-brand-900/60 text-brand-600 dark:text-brand-400'"
        >
          <component :is="getIcon(notif.type)" class="w-5 h-5" />
        </div>

        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between gap-2">
            <p class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate" :class="{ 'font-bold text-brand-700 dark:text-brand-300': !notif.read }">
              {{ notif.title }}
            </p>
            <span class="text-[11px] text-surface-400 shrink-0">{{ formatRelativeTime(notif.created_at) }}</span>
          </div>
          <p class="text-xs text-surface-600 dark:text-surface-300 mt-0.5 leading-relaxed">{{ notif.message }}</p>
        </div>

        <span v-if="!notif.read" class="w-2.5 h-2.5 rounded-full bg-brand-500 shrink-0 mt-1.5" />
      </div>
    </div>
  </div>
</template>

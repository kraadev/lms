<script setup lang="ts">
import { ref } from 'vue'
import { Maximize2, Minimize2, Clock, ShieldCheck } from 'lucide-vue-next'

interface Props {
  quizTitle: string
  durationMinutes: number
  remainingSeconds: number
}

const props = withDefaults(defineProps<Props>(), {
  quizTitle: 'Ujian Berbatas Waktu',
  durationMinutes: 60,
  remainingSeconds: 3600
})

const isFullscreen = ref(false)

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen().catch(() => {})
    isFullscreen.value = true
  } else {
    document.exitFullscreen().catch(() => {})
    isFullscreen.value = false
  }
}

function formatTime(totalSec: number): string {
  const m = Math.floor(totalSec / 60)
  const s = totalSec % 60
  return `${m}:${s < 10 ? '0' : ''}${s}`
}
</script>

<template>
  <div class="min-h-screen bg-surface-50 dark:bg-surface-950 flex flex-col">
    <!-- Zen Mode Header -->
    <header class="h-14 px-6 border-b border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 flex items-center justify-between shrink-0">
      <div class="flex items-center gap-3">
        <ShieldCheck class="w-5 h-5 text-brand-600" />
        <h2 class="text-sm font-bold text-surface-900 dark:text-surface-100 truncate">{{ quizTitle }}</h2>
      </div>

      <!-- Floating Timer -->
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2 px-3 py-1 rounded-full bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-800/60 text-xs font-bold text-amber-700 dark:text-amber-300">
          <Clock class="w-3.5 h-3.5 animate-pulse" />
          <span>Sisa Waktu: {{ formatTime(remainingSeconds) }}</span>
        </div>

        <button
          type="button"
          class="p-2 rounded-lg text-surface-500 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
          :title="isFullscreen ? 'Keluar Fullscreen' : 'Layar Penuh'"
          @click="toggleFullscreen"
        >
          <component :is="isFullscreen ? Minimize2 : Maximize2" class="w-4 h-4" />
        </button>
      </div>
    </header>

    <!-- Content Slot -->
    <main class="flex-1 max-w-4xl w-full mx-auto p-6">
      <slot />
    </main>
  </div>
</template>

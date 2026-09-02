<script setup lang="ts">
import type { Toast } from '~/composables/useToast'

const { toasts, remove } = useToast()

const iconByType: Record<Toast['type'], string> = {
  success: 'M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
  error: 'M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z',
  warning: 'M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z',
  info: 'M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z'
}

const colorByType: Record<Toast['type'], string> = {
  success: 'text-emerald-500',
  error: 'text-rose-500',
  warning: 'text-amber-500',
  info: 'text-sky-500'
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed bottom-4 right-4 z-[100] flex flex-col gap-2 max-w-sm w-full pointer-events-none" aria-live="polite">
      <TransitionGroup
        tag="div"
        class="flex flex-col gap-2"
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-active-class="transition duration-200 ease-in absolute w-full"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="translate-x-full opacity-0"
      >
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="pointer-events-auto flex items-start gap-3 bg-white dark:bg-surface-900 rounded-xl shadow-elevated border border-surface-200 dark:border-surface-800 p-4"
          role="alert"
        >
          <svg :class="['w-5 h-5 shrink-0 mt-0.5', colorByType[toast.type]]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" :d="iconByType[toast.type]" />
          </svg>
          <div class="flex-1 min-w-0">
            <p v-if="toast.title" class="text-sm font-semibold text-surface-900 dark:text-surface-100">{{ toast.title }}</p>
            <p class="text-sm text-surface-600 dark:text-surface-400">{{ toast.message }}</p>
          </div>
          <button
            type="button"
            class="shrink-0 p-1 text-surface-400 hover:text-surface-600 dark:hover:text-surface-200 rounded transition-colors"
            :aria-label="`Tutup notifikasi`"
            @click="remove(toast.id)"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
interface Props {
  message?: string
  status?: number | null
  retry?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  retry: true
})

const emit = defineEmits<{ retry: [] }>()

const title = computed(() => {
  if (props.status === 403) return 'Akses Ditolak'
  if (props.status === 404) return 'Tidak Ditemukan'
  if (props.status === 0) return 'Tidak Dapat Terhubung'
  return 'Terjadi Kesalahan'
})

const icon = computed(() => {
  if (props.status === 403) return 'lock'
  if (props.status === 404) return 'search'
  if (props.status === 0) return 'wifi-off'
  return 'alert-circle'
})
</script>

<template>
  <div class="flex flex-col items-center justify-center py-14 px-6 text-center">
    <div class="mb-4 p-3 rounded-2xl bg-rose-50 dark:bg-rose-900/20 text-rose-500">
      <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
        <path v-if="status === 403" stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
        <path v-else-if="status === 404" stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 15.803M10.5 7.5v6m3-3h-6" />
        <path v-else-if="status === 0" stroke-linecap="round" stroke-linejoin="round" d="M3 3l18 18M8.111 8.111A7.5 7.5 0 0021 12M12.889 12.889A7.5 7.5 0 013 12m6-1.5h.01M12 12h.01M15 13.5h.01" />
        <path v-else stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
      </svg>
    </div>
    <h3 class="text-base font-semibold text-surface-800 dark:text-surface-200 mb-1">{{ title }}</h3>
    <p class="text-sm text-surface-500 dark:text-surface-400 max-w-sm">{{ message || 'Sesuatu yang tidak terduga terjadi.' }}</p>
    <button
      v-if="retry && status !== 403 && status !== 404"
      type="button"
      class="mt-5 px-4 py-2 text-sm font-medium rounded-lg bg-surface-100 dark:bg-surface-800 text-surface-700 dark:text-surface-200 hover:bg-surface-200 dark:hover:bg-surface-700 transition-colors"
      @click="emit('retry')"
    >
      Coba Lagi
    </button>
  </div>
</template>

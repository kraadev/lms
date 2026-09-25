<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { WifiOff, Wifi } from 'lucide-vue-next'

const isOnline = ref(true)
const showRestored = ref(false)

function handleOnline() {
  isOnline.value = true
  showRestored.value = true
  setTimeout(() => { showRestored.value = false }, 3000)
}

function handleOffline() {
  isOnline.value = false
  showRestored.value = false
}

onMounted(() => {
  if (import.meta.client) {
    isOnline.value = navigator.onLine
    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)
  }
})

onUnmounted(() => {
  if (import.meta.client) {
    window.removeEventListener('online', handleOnline)
    window.removeEventListener('offline', handleOffline)
  }
})
</script>

<template>
  <div
    v-if="!isOnline"
    class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50 px-4 py-2 rounded-xl bg-rose-600 text-white text-xs font-semibold shadow-elevated flex items-center gap-2 animate-bounce"
  >
    <WifiOff class="w-4 h-4" />
    <span>Koneksi internet terputus. Menunggu koneksi kembali...</span>
  </div>
  <div
    v-else-if="showRestored"
    class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50 px-4 py-2 rounded-xl bg-emerald-600 text-white text-xs font-semibold shadow-elevated flex items-center gap-2 transition-opacity"
  >
    <Wifi class="w-4 h-4" />
    <span>Koneksi internet berhasil dipulihkan!</span>
  </div>
</template>

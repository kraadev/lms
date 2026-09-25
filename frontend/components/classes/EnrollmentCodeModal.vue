<script setup lang="ts">
import { ref } from 'vue'
import { QrCode, Copy, Check, X } from 'lucide-vue-next'

interface Props {
  isOpen: boolean
  classTitle: string
  code: string
}
const props = withDefaults(defineProps<Props>(), { isOpen: false, code: '' })
const emit = defineEmits<{ (e: 'close'): void }>()
const copied = ref(false)

function copyCode() {
  if (navigator.clipboard) navigator.clipboard.writeText(props.code)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}
</script>
<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
    <div class="w-full max-w-sm bg-white dark:bg-surface-900 rounded-2xl p-6 shadow-elevated border border-surface-200 dark:border-surface-800 text-center">
      <div class="flex justify-end"><button @click="emit('close')"><X class="w-4 h-4 text-surface-400" /></button></div>
      <QrCode class="w-12 h-12 text-brand-600 mx-auto mb-3" />
      <h3 class="text-sm font-bold text-surface-900 dark:text-surface-100">Kode Masuk Kelas</h3>
      <p class="text-xs text-surface-400 mb-4">{{ classTitle }}</p>
      <div class="p-3 bg-surface-50 dark:bg-surface-800/80 rounded-xl font-mono text-lg font-bold tracking-widest text-brand-600 mb-4">{{ code }}</div>
      <UiButton variant="primary" size="sm" class="w-full" @click="copyCode">
        <component :is="copied ? Check : Copy" class="w-3.5 h-3.5 mr-1" />
        {{ copied ? 'Tersalin' : 'Salin Kode' }}
      </UiButton>
    </div>
  </div>
</template>

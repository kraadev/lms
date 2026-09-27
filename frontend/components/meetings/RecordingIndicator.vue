<script setup lang="ts">
import { Disc, Download } from 'lucide-vue-next'

interface Props {
  isRecording?: boolean
  recordingUrl?: string
}
withDefaults(defineProps<Props>(), { isRecording: false, recordingUrl: '' })
</script>
<template>
  <div class="inline-flex items-center gap-2">
    <div v-if="isRecording" class="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-rose-500/10 border border-rose-500/30 text-rose-600 text-xs font-bold animate-pulse">
      <Disc class="w-3.5 h-3.5" />
      <span>REC</span>
    </div>
    <a v-if="recordingUrl" :href="recordingUrl" download class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-surface-100 dark:bg-surface-800 text-surface-700 dark:text-surface-300 text-xs hover:bg-surface-200">
      <Download class="w-3.5 h-3.5" />
      <span>Unduh Rekaman</span>
    </a>
  </div>
</template>

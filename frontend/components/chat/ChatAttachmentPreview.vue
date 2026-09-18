<script setup lang="ts">
import { File, Download, Image, ExternalLink } from 'lucide-vue-next'

interface Props {
  fileName: string
  fileUrl: string
  fileSizeBytes?: number
  isImage?: boolean
}

withDefaults(defineProps<Props>(), {
  fileSizeBytes: 0,
  isImage: false
})
</script>

<template>
  <div class="inline-block max-w-sm rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 p-2.5 shadow-soft">
    <!-- Image thumbnail preview -->
    <div v-if="isImage" class="rounded-lg overflow-hidden mb-2 max-h-48 bg-surface-100 dark:bg-surface-800">
      <img :src="fileUrl" :alt="fileName" class="w-full h-full object-cover" />
    </div>

    <div class="flex items-center justify-between gap-3">
      <div class="flex items-center gap-2 min-w-0">
        <component :is="isImage ? Image : File" class="w-4 h-4 text-brand-600 shrink-0" />
        <div class="min-w-0">
          <p class="text-xs font-semibold text-surface-800 dark:text-surface-200 truncate">{{ fileName }}</p>
          <p v-if="fileSizeBytes" class="text-[10px] text-surface-400">{{ (fileSizeBytes / 1024).toFixed(0) }} KB</p>
        </div>
      </div>

      <a :href="fileUrl" target="_blank" download class="p-1.5 rounded-lg hover:bg-surface-100 dark:hover:bg-surface-800 text-surface-500 transition-colors">
        <Download class="w-3.5 h-3.5" />
      </a>
    </div>
  </div>
</template>

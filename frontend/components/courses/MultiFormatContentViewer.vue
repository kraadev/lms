<script setup lang="ts">
import { ExternalLink, Play, FileText, BookOpen } from 'lucide-vue-next'

interface Props {
  contentType?: 'markdown' | 'video' | 'pdf'
  content?: string
  contentUrl?: string
  title?: string
}

withDefaults(defineProps<Props>(), {
  contentType: 'markdown',
  content: '',
  contentUrl: '',
  title: 'Materi Belajar'
})
</script>

<template>
  <div class="rounded-2xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft">
    <!-- Header -->
    <div class="flex items-center justify-between px-5 py-3.5 border-b border-surface-100 dark:border-surface-800">
      <div class="flex items-center gap-2">
        <Video v-if="contentType === 'video'" class="w-4.5 h-4.5 text-brand-600 dark:text-brand-400" />
        <FileText v-else-if="contentType === 'pdf'" class="w-4.5 h-4.5 text-amber-600 dark:text-amber-400" />
        <BookOpen v-else class="w-4.5 h-4.5 text-emerald-600 dark:text-emerald-400" />
        <h3 class="text-sm font-semibold text-surface-900 dark:text-surface-100">{{ title }}</h3>
      </div>

      <a
        v-if="contentUrl && contentType === 'pdf'"
        :href="contentUrl"
        target="_blank"
        rel="noopener noreferrer"
        class="inline-flex items-center gap-1 text-xs font-medium text-brand-600 dark:text-brand-400 hover:underline"
      >
        <span>Buka Berkas</span>
        <ExternalLink class="w-3 h-3" />
      </a>
    </div>

    <!-- Content Body based on type -->
    <div class="p-6">
      <!-- Video Stream / Player Mode -->
      <div v-if="contentType === 'video'" class="aspect-video w-full rounded-xl bg-surface-950 flex flex-col items-center justify-center relative overflow-hidden group">
        <video
          v-if="contentUrl"
          :src="contentUrl"
          controls
          class="w-full h-full object-contain"
        />
        <div v-else class="text-center p-6 space-y-2">
          <div class="w-12 h-12 rounded-full bg-brand-600/30 flex items-center justify-center mx-auto text-brand-400 group-hover:scale-110 transition-transform">
            <Play class="w-6 h-6 fill-current" />
          </div>
          <p class="text-xs text-surface-400">Pratinjau Video Materi Pelajaran</p>
        </div>
      </div>

      <!-- PDF Viewer Mode -->
      <div v-else-if="contentType === 'pdf'" class="w-full min-h-[420px] rounded-xl border border-surface-200 dark:border-surface-800 bg-surface-50 dark:bg-surface-950/40 p-6 flex flex-col items-center justify-center text-center">
        <FileText class="w-12 h-12 text-amber-500/80 mb-3" />
        <h4 class="text-sm font-semibold text-surface-800 dark:text-surface-200 mb-1">Dokumen Silabus PDF</h4>
        <p class="text-xs text-surface-500 dark:text-surface-400 max-w-sm mb-4">Materi disajikan dalam format portabel dokumen.</p>
        <a
          v-if="contentUrl"
          :href="contentUrl"
          target="_blank"
          rel="noopener noreferrer"
        >
          <UiButton size="sm" variant="outline">
            <ExternalLink class="w-3.5 h-3.5 mr-1.5" /> Unduh Dokumen
          </UiButton>
        </a>
      </div>

      <!-- Markdown Reading Mode -->
      <div v-else class="prose dark:prose-invert max-w-none text-sm text-surface-700 dark:text-surface-300 leading-relaxed whitespace-pre-line">
        {{ content || 'Materi pembelajaran teks belum tersedia.' }}
      </div>
    </div>
  </div>
</template>

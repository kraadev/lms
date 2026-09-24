<script setup lang="ts">
import { Pin, MessageCircle, Clock } from 'lucide-vue-next'

interface Props {
  authorName: string
  title: string
  content: string
  postedAt: string
  isPinned?: boolean
  commentCount?: number
}

withDefaults(defineProps<Props>(), {
  isPinned: false,
  commentCount: 0
})
</script>

<template>
  <div :class="['rounded-2xl border p-5 transition-all shadow-soft bg-white dark:bg-surface-900', isPinned ? 'border-brand-300 dark:border-brand-700/60 bg-brand-50/20 dark:bg-brand-950/10' : 'border-surface-200 dark:border-surface-800']">
    <div class="flex items-center justify-between mb-3">
      <div class="flex items-center gap-2.5">
        <UiAvatar :name="authorName" size="sm" />
        <div>
          <p class="text-xs font-bold text-surface-900 dark:text-surface-100">{{ authorName }}</p>
          <p class="text-[10px] text-surface-400 flex items-center gap-1 mt-0.5">
            <Clock class="w-3 h-3" /> {{ postedAt }}
          </p>
        </div>
      </div>

      <span v-if="isPinned" class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold bg-brand-100 dark:bg-brand-900/40 text-brand-700 dark:text-brand-300">
        <Pin class="w-3 h-3" /> Disematkan
      </span>
    </div>

    <h4 class="text-sm font-bold text-surface-900 dark:text-surface-100 mb-1.5">{{ title }}</h4>
    <p class="text-xs text-surface-600 dark:text-surface-300 leading-relaxed whitespace-pre-line">{{ content }}</p>

    <div class="mt-4 pt-3 border-t border-surface-100 dark:border-surface-800/60 flex items-center gap-1.5 text-xs text-surface-400">
      <MessageCircle class="w-3.5 h-3.5" />
      <span>{{ commentCount }} komentar</span>
    </div>
  </div>
</template>

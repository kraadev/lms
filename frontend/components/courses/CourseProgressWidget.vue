<script setup lang="ts">
import { Trophy, CheckCircle } from 'lucide-vue-next'

interface Props {
  progressPercent: number
  completedCount: number
  totalCount: number
  milestoneLabel: string
}

withDefaults(defineProps<Props>(), {
  progressPercent: 0,
  completedCount: 0,
  totalCount: 0,
  milestoneLabel: 'Pemula'
})
</script>

<template>
  <div class="p-4 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 shadow-soft">
    <div class="flex items-center justify-between mb-2">
      <div class="flex items-center gap-2">
        <Trophy class="w-4 h-4 text-amber-500" />
        <span class="text-xs font-semibold text-surface-800 dark:text-surface-200">Progres Belajar</span>
      </div>
      <span class="text-xs font-bold text-brand-600 dark:text-brand-400">{{ milestoneLabel }}</span>
    </div>

    <!-- Progress Bar -->
    <UiProgress :value="progressPercent" :max="100" size="sm" variant="brand" />

    <div class="flex items-center justify-between mt-2 text-[11px] text-surface-500 dark:text-surface-400">
      <span class="flex items-center gap-1">
        <CheckCircle class="w-3 h-3 text-emerald-500" />
        {{ completedCount }} dari {{ totalCount }} materi
      </span>
      <span class="font-bold text-surface-700 dark:text-surface-300">{{ progressPercent }}%</span>
    </div>
  </div>
</template>

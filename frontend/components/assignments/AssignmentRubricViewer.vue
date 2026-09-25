<script setup lang="ts">
export interface RubricItem {
  id: number
  criterion: string
  maxPoints: number
  description: string
}
interface Props {
  rubrics?: RubricItem[]
}
withDefaults(defineProps<Props>(), { rubrics: () => [] })
</script>
<template>
  <div class="rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 p-4 shadow-soft">
    <h4 class="text-xs font-bold uppercase tracking-wider text-surface-500 mb-3">Rubrik Penilaian</h4>
    <div v-if="!rubrics.length" class="text-xs text-surface-400">Tidak ada rubrik khusus.</div>
    <div v-else class="space-y-2">
      <div v-for="r in rubrics" :key="r.id" class="p-2.5 rounded-lg bg-surface-50 dark:bg-surface-800/50 flex justify-between items-start">
        <div>
          <p class="text-xs font-semibold text-surface-900 dark:text-surface-100">{{ r.criterion }}</p>
          <p class="text-[11px] text-surface-500">{{ r.description }}</p>
        </div>
        <span class="text-xs font-bold text-brand-600 bg-brand-50 dark:bg-brand-950/60 px-2 py-0.5 rounded-md">{{ r.maxPoints }} pts</span>
      </div>
    </div>
  </div>
</template>

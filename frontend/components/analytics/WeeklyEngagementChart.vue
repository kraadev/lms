<script setup lang="ts">
import { BarChart3 } from 'lucide-vue-next'
interface Props {
  weeklyHours?: number[]
}
withDefaults(defineProps<Props>(), {
  weeklyHours: () => [2.5, 3.8, 1.2, 4.5, 3.0, 5.2, 2.0]
})
const days = ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min']
</script>
<template>
  <div class="rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 p-4 shadow-soft">
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-2">
        <BarChart3 class="w-4 h-4 text-brand-600" />
        <h4 class="text-xs font-bold text-surface-900 dark:text-surface-100">Aktivitas Belajar Mingguan</h4>
      </div>
      <span class="text-xs text-surface-400">Total: 22.2 Jam</span>
    </div>
    <div class="flex items-end justify-between gap-2 h-28 pt-4">
      <div v-for="(h, i) in weeklyHours" :key="i" class="flex-1 flex flex-col items-center gap-1.5 h-full justify-end">
        <div class="w-full bg-brand-500/80 rounded-t-md hover:bg-brand-600 transition-all" :style="{ height: `${(h / 6) * 100}%` }" />
        <span class="text-[10px] text-surface-400">{{ days[i] }}</span>
      </div>
    </div>
  </div>
</template>

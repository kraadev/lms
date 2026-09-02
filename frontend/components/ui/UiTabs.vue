<script setup lang="ts">
interface Tab {
  key: string
  label: string
  icon?: any
  count?: number
}

interface Props {
  modelValue: string
  tabs: Tab[]
}

const props = defineProps<Props>()
const emit = defineEmits<{ 'update:modelValue': [key: string] }>()
</script>

<template>
  <div class="flex items-center gap-0.5 border-b border-surface-200 dark:border-surface-800 overflow-x-auto no-scrollbar">
    <button
      v-for="tab in tabs"
      :key="tab.key"
      type="button"
      :class="[
        'flex items-center gap-1.5 px-3 py-2.5 text-sm font-medium whitespace-nowrap border-b-2 transition-all -mb-px',
        modelValue === tab.key
          ? 'border-brand-600 text-brand-600 dark:text-brand-400 dark:border-brand-400'
          : 'border-transparent text-surface-500 dark:text-surface-400 hover:text-surface-800 dark:hover:text-surface-200 hover:border-surface-300 dark:hover:border-surface-600'
      ]"
      @click="emit('update:modelValue', tab.key)"
    >
      <component :is="tab.icon" v-if="tab.icon" class="w-4 h-4 shrink-0" />
      {{ tab.label }}
      <span
        v-if="tab.count !== undefined"
        :class="[
          'inline-flex items-center justify-center min-w-[18px] h-4.5 px-1.5 rounded-full text-[10px] font-semibold',
          modelValue === tab.key
            ? 'bg-brand-100 text-brand-700 dark:bg-brand-900/50 dark:text-brand-300'
            : 'bg-surface-200 text-surface-600 dark:bg-surface-700 dark:text-surface-400'
        ]"
      >{{ tab.count }}</span>
    </button>
  </div>
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
</style>

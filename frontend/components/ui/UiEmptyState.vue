<script setup lang="ts">
interface Props {
  icon?: any
  title?: string
  description?: string
}

defineProps<Props>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-14 px-6 text-center">
    <div v-if="icon || $slots.icon" class="mb-4 p-3 rounded-2xl bg-surface-100 dark:bg-surface-800 text-surface-400 dark:text-surface-500">
      <slot name="icon">
        <component :is="icon" class="w-8 h-8" />
      </slot>
    </div>
    <h3 class="text-base font-semibold text-surface-700 dark:text-surface-300 mb-1">
      <slot name="title">{{ title || 'Tidak ada data' }}</slot>
    </h3>
    <p v-if="description || $slots.description" class="text-sm text-surface-400 dark:text-surface-500 max-w-sm">
      <slot name="description">{{ description }}</slot>
    </p>
    <div v-if="$slots.action" class="mt-5">
      <slot name="action" />
    </div>
  </div>
</template>

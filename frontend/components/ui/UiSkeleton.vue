<script setup lang="ts">
interface Props {
  rows?: number
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  rows: 3
})
</script>

<template>
  <div :class="['animate-pulse', props.class]">
    <slot>
      <div v-for="i in rows" :key="i" :class="['h-4 bg-surface-200 dark:bg-surface-800 rounded-md mb-2.5', i === rows ? 'w-3/4' : 'w-full']" />
    </slot>
  </div>
</template>

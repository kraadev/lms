<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  icon?: any
  title?: string
  description?: string
  variant?: 'default' | 'card' | 'dashed'
  size?: 'sm' | 'md' | 'lg'
  actionLabel?: string
  actionTo?: string
  secondaryActionLabel?: string
  secondaryActionTo?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'md'
})

const emit = defineEmits<{
  (e: 'action'): void
  (e: 'secondaryAction'): void
}>()

const containerClasses = computed(() => {
  const base = 'flex flex-col items-center justify-center text-center transition-colors'
  const paddingMap = {
    sm: 'py-8 px-4',
    md: 'py-14 px-6',
    lg: 'py-20 px-8'
  }
  const variantMap = {
    default: '',
    card: 'bg-white dark:bg-surface-900 border border-surface-200/80 dark:border-surface-800 rounded-2xl shadow-soft',
    dashed: 'border-2 border-dashed border-surface-200 dark:border-surface-800 rounded-2xl bg-surface-50/50 dark:bg-surface-900/30'
  }
  return [base, paddingMap[props.size], variantMap[props.variant]].filter(Boolean).join(' ')
})

const iconWrapperClasses = computed(() => {
  const sizeMap = {
    sm: 'mb-3 p-2 rounded-xl',
    md: 'mb-4 p-3 rounded-2xl',
    lg: 'mb-5 p-4 rounded-3xl'
  }
  return `${sizeMap[props.size]} bg-surface-100 dark:bg-surface-800/80 text-surface-400 dark:text-surface-500`
})

const iconSizeClasses = computed(() => {
  const sizeMap = {
    sm: 'w-6 h-6',
    md: 'w-8 h-8',
    lg: 'w-10 h-10'
  }
  return sizeMap[props.size]
})

const titleClasses = computed(() => {
  const sizeMap = {
    sm: 'text-sm font-semibold text-surface-800 dark:text-surface-200 mb-0.5',
    md: 'text-base font-semibold text-surface-800 dark:text-surface-200 mb-1',
    lg: 'text-lg font-semibold text-surface-900 dark:text-surface-100 mb-1.5'
  }
  return sizeMap[props.size]
})

const descriptionClasses = computed(() => {
  const sizeMap = {
    sm: 'text-xs text-surface-500 dark:text-surface-400 max-w-xs',
    md: 'text-sm text-surface-500 dark:text-surface-400 max-w-sm',
    lg: 'text-base text-surface-500 dark:text-surface-400 max-w-md'
  }
  return sizeMap[props.size]
})
</script>

<template>
  <div :class="containerClasses">
    <!-- Icon or Illustration Container -->
    <div v-if="icon || $slots.icon" :class="iconWrapperClasses">
      <slot name="icon">
        <component :is="icon" :class="iconSizeClasses" />
      </slot>
    </div>

    <!-- Title -->
    <h3 :class="titleClasses">
      <slot name="title">{{ title || 'Tidak ada data' }}</slot>
    </h3>

    <!-- Description -->
    <p v-if="description || $slots.description" :class="descriptionClasses">
      <slot name="description">{{ description }}</slot>
    </p>

    <!-- Actions Area -->
    <div
      v-if="$slots.action || actionLabel || $slots.secondaryAction || secondaryActionLabel"
      class="mt-5 flex flex-wrap items-center justify-center gap-2.5"
    >
      <slot name="action">
        <UiButton
          v-if="actionLabel"
          :to="actionTo"
          variant="primary"
          size="sm"
          @click="emit('action')"
        >
          {{ actionLabel }}
        </UiButton>
      </slot>

      <slot name="secondaryAction">
        <UiButton
          v-if="secondaryActionLabel"
          :to="secondaryActionTo"
          variant="outline"
          size="sm"
          @click="emit('secondaryAction')"
        >
          {{ secondaryActionLabel }}
        </UiButton>
      </slot>
    </div>
  </div>
</template>

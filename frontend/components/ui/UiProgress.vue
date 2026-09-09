<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  value?: number
  max?: number
  size?: 'xs' | 'sm' | 'md' | 'lg'
  variant?: 'brand' | 'success' | 'warning' | 'danger'
  showLabel?: boolean
  labelPosition?: 'top' | 'right'
  animated?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  value: 0,
  max: 100,
  size: 'md',
  variant: 'brand',
  showLabel: false,
  labelPosition: 'top',
  animated: false
})

const percentage = computed(() => {
  if (props.max <= 0) return 0
  const pct = Math.round((props.value / props.max) * 100)
  return Math.min(100, Math.max(0, pct))
})

const sizeClasses = {
  xs: 'h-1',
  sm: 'h-1.5',
  md: 'h-2.5',
  lg: 'h-4'
}

const variantClasses = {
  brand: 'bg-brand-600 dark:bg-brand-500',
  success: 'bg-emerald-600 dark:bg-emerald-500',
  warning: 'bg-amber-500 dark:bg-amber-400',
  danger: 'bg-rose-600 dark:bg-rose-500'
}
</script>

<template>
  <div class="w-full">
    <!-- Top Label Option -->
    <div
      v-if="showLabel && labelPosition === 'top'"
      class="flex justify-between items-center mb-1.5 text-xs font-medium text-surface-600 dark:text-surface-300"
    >
      <slot name="label">
        <span>Progress</span>
      </slot>
      <span class="font-semibold text-surface-800 dark:text-surface-100">{{ percentage }}%</span>
    </div>

    <!-- Progress Track & Bar -->
    <div class="flex items-center gap-3">
      <div
        :class="[
          'w-full bg-surface-100 dark:bg-surface-800 rounded-full overflow-hidden',
          sizeClasses[size]
        ]"
        role="progressbar"
        :aria-valuenow="value"
        :aria-valuemin="0"
        :aria-valuemax="max"
      >
        <div
          :class="[
            'h-full rounded-full transition-all duration-300 ease-out',
            variantClasses[variant],
            animated ? 'animate-pulse' : ''
          ]"
          :style="{ width: `${percentage}%` }"
        />
      </div>

      <!-- Right Label Option -->
      <span
        v-if="showLabel && labelPosition === 'right'"
        class="text-xs font-semibold text-surface-700 dark:text-surface-300 shrink-0 w-9 text-right"
      >
        {{ percentage }}%
      </span>
    </div>
  </div>
</template>

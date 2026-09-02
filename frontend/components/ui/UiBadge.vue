<script setup lang="ts">
type BadgeVariant = 'default' | 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'purple'
type BadgeSize = 'sm' | 'md'

interface Props {
  variant?: BadgeVariant
  size?: BadgeSize
  dot?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'md',
  dot: false
})

const variantClasses: Record<BadgeVariant, string> = {
  default: 'bg-surface-100 text-surface-700 dark:bg-surface-800 dark:text-surface-300',
  primary: 'bg-brand-100 text-brand-700 dark:bg-brand-900/50 dark:text-brand-300',
  success: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300',
  warning: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300',
  danger: 'bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-300',
  info: 'bg-sky-100 text-sky-700 dark:bg-sky-900/40 dark:text-sky-300',
  purple: 'bg-violet-100 text-violet-700 dark:bg-violet-900/40 dark:text-violet-300'
}

const dotClasses: Record<BadgeVariant, string> = {
  default: 'bg-surface-400',
  primary: 'bg-brand-500',
  success: 'bg-emerald-500',
  warning: 'bg-amber-500',
  danger: 'bg-rose-500',
  info: 'bg-sky-500',
  purple: 'bg-violet-500'
}

const sizeClasses: Record<BadgeSize, string> = {
  sm: 'text-[10px] px-1.5 py-0.5 rounded font-medium',
  md: 'text-xs px-2 py-0.5 rounded-md font-medium'
}
</script>

<template>
  <span :class="['inline-flex items-center gap-1.5 whitespace-nowrap', variantClasses[variant], sizeClasses[size]]">
    <span v-if="dot" :class="['w-1.5 h-1.5 rounded-full shrink-0', dotClasses[variant]]" />
    <slot />
  </span>
</template>

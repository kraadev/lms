<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger' | 'success'
  size?: 'xs' | 'sm' | 'md' | 'lg'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  block?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  block: false
})

const sizeClasses = {
  xs: 'px-2.5 py-1 text-xs gap-1.5 rounded-md font-medium',
  sm: 'px-3 py-1.5 text-xs sm:text-sm gap-1.5 rounded-lg font-medium',
  md: 'px-4 py-2 text-sm gap-2 rounded-lg font-medium',
  lg: 'px-5 py-2.5 text-base gap-2.5 rounded-xl font-semibold'
}

const variantClasses = {
  primary: 'bg-brand-600 hover:bg-brand-700 text-white shadow-sm active:bg-brand-800 focus-visible:ring-brand-500 dark:bg-brand-600 dark:hover:bg-brand-500',
  secondary: 'bg-surface-100 hover:bg-surface-200 text-surface-800 dark:bg-surface-800 dark:hover:bg-surface-700 dark:text-surface-100 focus-visible:ring-surface-400',
  outline: 'border border-surface-300 dark:border-surface-700 text-surface-700 dark:text-surface-200 hover:bg-surface-50 dark:hover:bg-surface-800/60 focus-visible:ring-brand-500',
  ghost: 'text-surface-600 dark:text-surface-300 hover:bg-surface-100 dark:hover:bg-surface-800/70 hover:text-surface-900 dark:hover:text-surface-100 focus-visible:ring-surface-400',
  danger: 'bg-rose-600 hover:bg-rose-700 text-white shadow-sm active:bg-rose-800 focus-visible:ring-rose-500',
  success: 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-sm active:bg-emerald-800 focus-visible:ring-emerald-500'
}
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'inline-flex items-center justify-center transition-all duration-150 select-none cursor-pointer',
      sizeClasses[size],
      variantClasses[variant],
      block ? 'w-full' : '',
      disabled || loading ? 'opacity-60 cursor-not-allowed pointer-events-none' : ''
    ]"
  >
    <svg
      v-if="loading"
      class="animate-spin -ml-0.5 h-4 w-4 text-current"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    <slot />
  </button>
</template>

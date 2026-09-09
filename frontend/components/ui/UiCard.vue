<script setup lang="ts">
interface Props {
  as?: string
  variant?: 'default' | 'flat' | 'outline' | 'interactive'
  padding?: 'none' | 'sm' | 'md' | 'lg'
  rounded?: 'sm' | 'md' | 'lg' | 'xl' | '2xl'
  hover?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  as: 'div',
  variant: 'default',
  padding: 'md',
  rounded: 'xl',
  hover: false
})

const variantClasses = {
  default: 'bg-white dark:bg-surface-900 border border-surface-200/80 dark:border-surface-800 shadow-soft',
  flat: 'bg-surface-50 dark:bg-surface-900/60 border border-transparent',
  outline: 'bg-transparent border border-surface-200 dark:border-surface-800',
  interactive: 'bg-white dark:bg-surface-900 border border-surface-200 dark:border-surface-800 shadow-soft hover:shadow-elevated hover:border-brand-300 dark:hover:border-brand-700/60 transition-all duration-200 cursor-pointer'
}

const paddingClasses = {
  none: 'p-0',
  sm: 'p-3 sm:p-4',
  md: 'p-4 sm:p-6',
  lg: 'p-6 sm:p-8'
}

const roundedClasses = {
  sm: 'rounded-md',
  md: 'rounded-lg',
  lg: 'rounded-xl',
  xl: 'rounded-2xl',
  '2xl': 'rounded-3xl'
}
</script>

<template>
  <component
    :is="as"
    :class="[
      'relative overflow-hidden transition-colors',
      roundedClasses[rounded],
      variantClasses[variant],
      hover && variant !== 'interactive' ? 'hover:shadow-elevated hover:border-surface-300 dark:hover:border-surface-700 transition-all duration-200' : ''
    ]"
  >
    <!-- Card Header (Optional Slot) -->
    <div
      v-if="$slots.header || $slots.actions"
      class="flex items-center justify-between px-4 sm:px-6 py-4 border-b border-surface-100 dark:border-surface-800/80"
    >
      <div class="font-semibold text-surface-900 dark:text-surface-50 text-base">
        <slot name="header" />
      </div>
      <div v-if="$slots.actions" class="flex items-center space-x-2">
        <slot name="actions" />
      </div>
    </div>

    <!-- Card Main Body -->
    <div :class="paddingClasses[padding]">
      <slot />
    </div>

    <!-- Card Footer (Optional Slot) -->
    <div
      v-if="$slots.footer"
      class="px-4 sm:px-6 py-3 bg-surface-50/50 dark:bg-surface-900/50 border-t border-surface-100 dark:border-surface-800/80 text-sm text-surface-500 dark:text-surface-400"
    >
      <slot name="footer" />
    </div>
  </component>
</template>

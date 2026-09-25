<script setup lang="ts">
interface Props {
  count?: number
  max?: number
  pulse?: boolean
  variant?: 'danger' | 'brand' | 'success'
}
withDefaults(defineProps<Props>(), {
  count: 0,
  max: 99,
  pulse: true,
  variant: 'danger'
})
const variantClasses = {
  danger: 'bg-rose-500 text-white',
  brand: 'bg-brand-600 text-white',
  success: 'bg-emerald-500 text-white'
}
</script>
<template>
  <span
    v-if="count > 0"
    :class="[
      'inline-flex items-center justify-center min-w-[18px] h-4.5 px-1 rounded-full text-[10px] font-bold shadow-xs',
      variantClasses[variant],
      pulse ? 'animate-pulse' : ''
    ]"
  >
    {{ count > max ? `${max}+` : count }}
  </span>
</template>

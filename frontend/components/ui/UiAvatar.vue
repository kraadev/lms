<script setup lang="ts">
interface Props {
  name?: string
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  src?: string | null
  color?: 'indigo' | 'emerald' | 'sky' | 'amber' | 'rose' | 'violet' | 'auto'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  color: 'auto'
})

const sizeClasses = {
  xs: 'w-6 h-6 text-[9px]',
  sm: 'w-8 h-8 text-xs',
  md: 'w-9 h-9 text-sm',
  lg: 'w-11 h-11 text-base',
  xl: 'w-14 h-14 text-lg'
}

const colorMap = ['indigo', 'emerald', 'sky', 'amber', 'rose', 'violet']

const bgColor = computed(() => {
  if (props.color !== 'auto') return props.color
  if (!props.name) return 'indigo'
  const index = props.name.charCodeAt(0) % colorMap.length
  return colorMap[index]
})

const colorClasses: Record<string, string> = {
  indigo: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/50 dark:text-indigo-300',
  emerald: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/50 dark:text-emerald-300',
  sky: 'bg-sky-100 text-sky-700 dark:bg-sky-900/50 dark:text-sky-300',
  amber: 'bg-amber-100 text-amber-700 dark:bg-amber-900/50 dark:text-amber-300',
  rose: 'bg-rose-100 text-rose-700 dark:bg-rose-900/50 dark:text-rose-300',
  violet: 'bg-violet-100 text-violet-700 dark:bg-violet-900/50 dark:text-violet-300'
}

function getInitials(name?: string): string {
  if (!name) return '?'
  const parts = name.trim().split(/\s+/)
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase()
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
}
</script>

<template>
  <div
    :class="[
      'inline-flex items-center justify-center shrink-0 rounded-full overflow-hidden font-semibold',
      sizeClasses[size],
      !src ? colorClasses[bgColor] : ''
    ]"
    :aria-label="name ? `Avatar for ${name}` : 'User avatar'"
    role="img"
  >
    <img v-if="src" :src="src" :alt="name || 'avatar'" class="w-full h-full object-cover" />
    <span v-else>{{ getInitials(name) }}</span>
  </div>
</template>

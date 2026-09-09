<script setup lang="ts">
import { computed } from 'vue'

export interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
  description?: string
}

interface Props {
  modelValue?: string | number | null
  options: (SelectOption | string)[]
  label?: string
  placeholder?: string
  error?: string
  hint?: string
  disabled?: boolean
  required?: boolean
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  options: () => [],
  placeholder: 'Pilih salah satu...',
  disabled: false,
  required: false,
  size: 'md'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void
  (e: 'change', value: string | number): void
}>()

const normalizedOptions = computed<SelectOption[]>(() => {
  return props.options.map((opt) => {
    if (typeof opt === 'string' || typeof opt === 'number') {
      return { label: String(opt), value: opt }
    }
    return opt
  })
})

const sizeClasses = {
  sm: 'py-1.5 pl-3 pr-8 text-xs rounded-lg',
  md: 'py-2 pl-3.5 pr-10 text-sm rounded-lg',
  lg: 'py-2.5 pl-4 pr-11 text-base rounded-xl'
}

const onChange = (e: Event) => {
  const target = e.target as HTMLSelectElement
  emit('update:modelValue', target.value)
  emit('change', target.value)
}
</script>

<template>
  <div class="w-full">
    <!-- Label -->
    <label
      v-if="label"
      class="block text-xs font-semibold text-surface-700 dark:text-surface-200 mb-1.5 select-none"
    >
      {{ label }}
      <span v-if="required" class="text-rose-500 ml-0.5">*</span>
    </label>

    <!-- Select Field Container -->
    <div class="relative">
      <select
        :value="modelValue"
        :disabled="disabled"
        :required="required"
        :class="[
          'w-full appearance-none bg-white dark:bg-surface-900 border transition-all duration-150',
          'text-surface-900 dark:text-surface-50 focus:outline-none focus:ring-2',
          sizeClasses[size],
          error
            ? 'border-rose-400 dark:border-rose-500/80 focus:border-rose-500 focus:ring-rose-500/20'
            : 'border-surface-200 dark:border-surface-800 focus:border-brand-500 focus:ring-brand-500/20 dark:focus:ring-brand-500/30',
          disabled ? 'opacity-60 cursor-not-allowed bg-surface-50 dark:bg-surface-800/50' : 'cursor-pointer'
        ]"
        @change="onChange"
      >
        <option v-if="placeholder" value="" disabled :selected="!modelValue">
          {{ placeholder }}
        </option>
        <option
          v-for="opt in normalizedOptions"
          :key="opt.value"
          :value="opt.value"
          :disabled="opt.disabled"
        >
          {{ opt.label }}
        </option>
      </select>

      <!-- Chevron Down Icon -->
      <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2.5 text-surface-400 dark:text-surface-500">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </div>

    <!-- Error Message -->
    <p v-if="error" class="mt-1 text-xs text-rose-500 dark:text-rose-400 font-medium">
      {{ error }}
    </p>

    <!-- Hint Helper -->
    <p v-else-if="hint" class="mt-1 text-xs text-surface-500 dark:text-surface-400">
      {{ hint }}
    </p>
  </div>
</template>

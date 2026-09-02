<script setup lang="ts">
interface Props {
  modelValue?: string | null
  placeholder?: string
  label?: string
  hint?: string
  error?: string | null
  disabled?: boolean
  required?: boolean
  rows?: number
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  required: false,
  rows: 4
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const inputId = computed(() => props.id || `textarea-${Math.random().toString(36).substring(2, 7)}`)
</script>

<template>
  <div class="w-full">
    <label v-if="label" :for="inputId" class="block mb-1.5 text-sm font-medium text-surface-700 dark:text-surface-300">
      {{ label }} <span v-if="required" class="text-rose-500">*</span>
    </label>
    <textarea
      :id="inputId"
      :value="modelValue ?? ''"
      :placeholder="placeholder"
      :disabled="disabled"
      :required="required"
      :rows="rows"
      :class="[
        'w-full rounded-lg border bg-white dark:bg-surface-900 text-surface-900 dark:text-surface-100 placeholder:text-surface-400 transition-colors outline-none py-2 px-3.5 text-sm resize-vertical min-h-[80px]',
        error
          ? 'border-rose-400 dark:border-rose-600 focus:ring-2 focus:ring-rose-400/30 focus:border-rose-400'
          : 'border-surface-300 dark:border-surface-700 focus:ring-2 focus:ring-brand-500/30 focus:border-brand-500 dark:focus:border-brand-400',
        disabled ? 'opacity-60 cursor-not-allowed' : ''
      ]"
      @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    />
    <p v-if="error" class="mt-1.5 text-xs text-rose-600 dark:text-rose-400">{{ error }}</p>
    <p v-else-if="hint" class="mt-1.5 text-xs text-surface-500">{{ hint }}</p>
  </div>
</template>

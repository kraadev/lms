<script setup lang="ts">
interface Props {
  show: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  closeOnBackdrop?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  closeOnBackdrop: true
})

const emit = defineEmits<{ close: [] }>()

const sizeClasses = {
  sm: 'max-w-sm',
  md: 'max-w-lg',
  lg: 'max-w-2xl',
  xl: 'max-w-4xl',
  full: 'max-w-full mx-4'
}

function onBackdropClick() {
  if (props.closeOnBackdrop) emit('close')
}

// Lock body scroll
watch(() => props.show, (val) => {
  if (import.meta.client) {
    document.body.style.overflow = val ? 'hidden' : ''
  }
}, { immediate: true })

onUnmounted(() => {
  if (import.meta.client) {
    document.body.style.overflow = ''
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/50 dark:bg-black/70"
          @click="onBackdropClick"
          aria-hidden="true"
        />
        <!-- Panel -->
        <Transition
          enter-active-class="transition duration-150 ease-out"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition duration-100 ease-in"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
        >
          <div
            v-if="show"
            :class="['relative w-full bg-white dark:bg-surface-900 rounded-xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]', sizeClasses[size]]"
          >
            <!-- Header -->
            <div v-if="title || $slots.header" class="flex items-center justify-between px-5 py-4 border-b border-surface-200 dark:border-surface-800 shrink-0">
              <slot name="header">
                <h2 class="text-base font-semibold text-surface-900 dark:text-surface-100">{{ title }}</h2>
              </slot>
              <button
                type="button"
                class="p-1 rounded-lg text-surface-400 hover:text-surface-600 dark:hover:text-surface-200 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
                aria-label="Tutup modal"
                @click="emit('close')"
              >
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <!-- Body -->
            <div class="flex-1 overflow-y-auto px-5 py-4">
              <slot />
            </div>
            <!-- Footer -->
            <div v-if="$slots.footer" class="px-5 py-4 border-t border-surface-200 dark:border-surface-800 shrink-0">
              <slot name="footer" />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

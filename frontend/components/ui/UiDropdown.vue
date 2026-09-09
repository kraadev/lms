<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface Props {
  align?: 'left' | 'right'
  width?: 'sm' | 'md' | 'lg' | 'auto'
  closeOnClick?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  align: 'right',
  width: 'md',
  closeOnClick: true
})

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

const toggle = () => {
  isOpen.value = !isOpen.value
}

const close = () => {
  isOpen.value = false
}

const handleClickOutside = (event: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    close()
  }
}

const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && isOpen.value) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})

const widthClasses = {
  sm: 'w-36',
  md: 'w-48',
  lg: 'w-64',
  auto: 'w-auto min-w-[8rem]'
}

const alignClasses = {
  left: 'left-0 origin-top-left',
  right: 'right-0 origin-top-right'
}
</script>

<template>
  <div ref="dropdownRef" class="relative inline-block text-left">
    <!-- Trigger slot with exposed toggle/isOpen state -->
    <div @click="toggle">
      <slot name="trigger" :is-open="isOpen" :toggle="toggle" />
    </div>

    <!-- Dropdown Content Popover -->
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <div
        v-if="isOpen"
        :class="[
          'absolute z-50 mt-1.5 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 p-1 shadow-elevated focus:outline-none',
          widthClasses[width],
          alignClasses[align]
        ]"
        @click="closeOnClick ? close() : null"
      >
        <slot :close="close" />
      </div>
    </Transition>
  </div>
</template>

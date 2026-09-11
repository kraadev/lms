<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Sun, Moon, Monitor, Check } from 'lucide-vue-next'
import { useTheme, type ThemeMode } from '~/composables/useTheme'

interface Props {
  variant?: 'dropdown' | 'segmented' | 'icon'
  size?: 'sm' | 'md'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'dropdown',
  size: 'md'
})

const { themeMode, applyTheme } = useTheme()
const isOpen = ref(false)
const containerRef = ref<HTMLElement | null>(null)

const options = [
  { key: 'light' as ThemeMode, label: 'Terang', icon: Sun },
  { key: 'dark' as ThemeMode, label: 'Gelap', icon: Moon },
  { key: 'system' as ThemeMode, label: 'Sistem', icon: Monitor }
]

const currentIcon = computed(() => {
  if (themeMode.value === 'dark') return Moon
  if (themeMode.value === 'light') return Sun
  return Monitor
})

const currentLabel = computed(() => {
  const found = options.find((o) => o.key === themeMode.value)
  return found ? found.label : 'Sistem'
})

function cycleTheme() {
  if (themeMode.value === 'light') {
    applyTheme('dark')
  } else if (themeMode.value === 'dark') {
    applyTheme('system')
  } else {
    applyTheme('light')
  }
}

function selectTheme(mode: ThemeMode) {
  applyTheme(mode)
  isOpen.value = false
}

function handleClickOutside(e: MouseEvent) {
  if (containerRef.value && !containerRef.value.contains(e.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div ref="containerRef" class="relative inline-flex items-center">
    <!-- Icon-only direct cycle button -->
    <button
      v-if="variant === 'icon'"
      type="button"
      :class="[
        'p-2 rounded-lg text-surface-500 dark:text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800 hover:text-surface-800 dark:text-surface-100 transition-colors focus-ring',
        size === 'sm' ? 'p-1.5' : 'p-2'
      ]"
      :aria-label="`Tema saat ini: ${currentLabel}. Klik untuk ganti.`"
      :title="`Tema: ${currentLabel}`"
      @click="cycleTheme"
    >
      <component :is="currentIcon" :class="size === 'sm' ? 'w-4 h-4' : 'w-4.5 h-4.5'" />
    </button>

    <!-- Segmented 3-way toggle -->
    <div
      v-else-if="variant === 'segmented'"
      class="inline-flex p-0.5 rounded-xl bg-surface-100 dark:bg-surface-800/80 border border-surface-200/80 dark:border-surface-700/60"
      role="radiogroup"
      aria-label="Pilih mode tema"
    >
      <button
        v-for="opt in options"
        :key="opt.key"
        type="button"
        role="radio"
        :aria-checked="themeMode === opt.key"
        :class="[
          'flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-medium transition-all focus-ring',
          themeMode === opt.key
            ? 'bg-white dark:bg-surface-900 text-brand-600 dark:text-brand-400 shadow-soft'
            : 'text-surface-500 dark:text-surface-400 hover:text-surface-800 dark:hover:text-surface-200'
        ]"
        @click="applyTheme(opt.key)"
      >
        <component :is="opt.icon" class="w-3.5 h-3.5" />
        <span>{{ opt.label }}</span>
      </button>
    </div>

    <!-- Dropdown Menu Mode (Default for Topbar) -->
    <div v-else class="relative">
      <button
        type="button"
        :class="[
          'rounded-lg text-surface-500 dark:text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800 hover:text-surface-800 dark:hover:text-surface-100 transition-colors focus-ring',
          size === 'sm' ? 'p-1.5' : 'p-2'
        ]"
        :aria-label="`Pilih tema (saat ini ${currentLabel})`"
        :title="`Tema: ${currentLabel}`"
        @click="isOpen = !isOpen"
      >
        <component :is="currentIcon" :class="size === 'sm' ? 'w-4 h-4' : 'w-4.5 h-4.5'" />
      </button>

      <Transition
        enter-active-class="transition duration-100 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-75 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div
          v-if="isOpen"
          class="absolute right-0 mt-1.5 top-full w-36 bg-white dark:bg-surface-900 rounded-xl shadow-elevated border border-surface-200 dark:border-surface-800 py-1 z-50 focus:outline-none"
        >
          <button
            v-for="opt in options"
            :key="opt.key"
            type="button"
            :class="[
              'w-full flex items-center justify-between px-3 py-2 text-xs font-medium transition-colors',
              themeMode === opt.key
                ? 'text-brand-600 dark:text-brand-400 bg-brand-50/70 dark:bg-brand-950/50'
                : 'text-surface-700 dark:text-surface-300 hover:bg-surface-50 dark:hover:bg-surface-800/60'
            ]"
            @click="selectTheme(opt.key)"
          >
            <div class="flex items-center gap-2">
              <component :is="opt.icon" class="w-4 h-4 shrink-0" />
              <span>{{ opt.label }}</span>
            </div>
            <Check v-if="themeMode === opt.key" class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
          </button>
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Palette, Check, RefreshCw } from 'lucide-vue-next'

const platformName = ref('LMS Portal')
const primaryColor = ref('#4f46e5')
const isSaved = ref(false)

const presetColors = ['#4f46e5', '#2563eb', '#0891b2', '#059669', '#d97706', '#dc2626', '#7c3aed']

function saveBranding() {
  if (import.meta.client) {
    document.documentElement.style.setProperty('--color-brand-primary', primaryColor.value)
  }
  isSaved.value = true
  setTimeout(() => { isSaved.value = false }, 2500)
}

function resetDefault() {
  primaryColor.value = '#4f46e5'
  saveBranding()
}
</script>

<template>
  <div class="rounded-2xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 p-6 shadow-soft space-y-6">
    <div class="flex items-center justify-between border-b border-surface-100 dark:border-surface-800 pb-4">
      <div>
        <h3 class="text-base font-bold text-surface-900 dark:text-surface-100">Kustomisasi Branding & Warna</h3>
        <p class="text-xs text-surface-500 dark:text-surface-400">Atur identitas visual platform secara langsung</p>
      </div>
      <Palette class="w-5 h-5 text-brand-600" />
    </div>

    <!-- Platform Name -->
    <div>
      <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-1.5">Nama Platform LMS</label>
      <input
        v-model="platformName"
        type="text"
        class="w-full sm:w-80 px-3 py-2 rounded-xl text-xs border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-950 text-surface-900 dark:text-surface-100"
      />
    </div>

    <!-- Color Presets -->
    <div>
      <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-2">Warna Utama Tema (Brand Accent)</label>
      <div class="flex items-center gap-2">
        <button
          v-for="color in presetColors"
          :key="color"
          type="button"
          class="w-8 h-8 rounded-full border-2 border-white dark:border-surface-900 shadow-sm flex items-center justify-center transition-transform hover:scale-110"
          :style="{ backgroundColor: color }"
          @click="primaryColor = color"
        >
          <Check v-if="primaryColor === color" class="w-4 h-4 text-white stroke-[3]" />
        </button>
      </div>
    </div>

    <!-- Actions -->
    <div class="pt-4 border-t border-surface-100 dark:border-surface-800 flex items-center justify-between">
      <button type="button" class="text-xs text-surface-400 hover:text-surface-600 flex items-center gap-1" @click="resetDefault">
        <RefreshCw class="w-3.5 h-3.5" /> Reset Default
      </button>

      <div class="flex items-center gap-2">
        <span v-if="isSaved" class="text-xs font-medium text-emerald-600 flex items-center gap-1">
          <Check class="w-3.5 h-3.5" /> Tersimpan!
        </span>
        <UiButton variant="primary" size="sm" @click="saveBranding">Terapkan Perubahan</UiButton>
      </div>
    </div>
  </div>
</template>

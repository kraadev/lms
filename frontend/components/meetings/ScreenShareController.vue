<script setup lang="ts">
import { ref } from 'vue'
import { ScreenShare, StopCircle, Tv, Signal } from 'lucide-vue-next'

export type QualityProfile = '1080p' | '720p' | '360p'

interface Props {
  isSharing?: boolean
  currentQuality?: QualityProfile
}

const props = withDefaults(defineProps<Props>(), {
  isSharing: false,
  currentQuality: '720p'
})

const emit = defineEmits<{
  (e: 'start-share', quality: QualityProfile): void
  (e: 'stop-share'): void
  (e: 'change-quality', quality: QualityProfile): void
}>()

const selectedQuality = ref<QualityProfile>(props.currentQuality)

const profiles = [
  { key: '1080p' as QualityProfile, label: 'Full HD (1080p 30fps)', desc: 'Kecepatan tinggi' },
  { key: '720p' as QualityProfile, label: 'Standard HD (720p 30fps)', desc: 'Rekomendasi' },
  { key: '360p' as QualityProfile, label: 'Hemat Kuota (360p 15fps)', desc: 'Koneksi lambat' }
]

function toggle() {
  if (props.isSharing) {
    emit('stop-share')
  } else {
    emit('start-share', selectedQuality.value)
  }
}
</script>

<template>
  <div class="inline-flex items-center gap-2 p-1.5 rounded-xl bg-surface-900/80 backdrop-blur-md border border-surface-800 text-white shadow-elevated">
    <!-- Toggle Share Screen -->
    <button
      type="button"
      :class="[
        'flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all',
        isSharing ? 'bg-rose-600 hover:bg-rose-700' : 'bg-brand-600 hover:bg-brand-700'
      ]"
      @click="toggle"
    >
      <component :is="isSharing ? StopCircle : ScreenShare" class="w-4 h-4" />
      <span>{{ isSharing ? 'Hentikan Share' : 'Bagikan Layar' }}</span>
    </button>

    <!-- Quality Selector -->
    <div class="relative">
      <select
        v-model="selectedQuality"
        :disabled="isSharing"
        class="bg-surface-800 text-surface-200 text-xs rounded-lg px-2 py-1.5 border border-surface-700 focus:outline-none cursor-pointer disabled:opacity-60"
        @change="emit('change-quality', selectedQuality)"
      >
        <option v-for="p in profiles" :key="p.key" :value="p.key">
          {{ p.label }}
        </option>
      </select>
    </div>
  </div>
</template>

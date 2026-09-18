<script setup lang="ts">
import { MicOff, Lock, Unlock, UserX } from 'lucide-vue-next'

interface Props {
  isLocked?: boolean
  participantCount?: number
}

withDefaults(defineProps<Props>(), {
  isLocked: false,
  participantCount: 0
})

const emit = defineEmits<{
  (e: 'mute-all'): void
  (e: 'toggle-lock'): void
}>()
</script>

<template>
  <div class="inline-flex items-center gap-2 p-1.5 rounded-xl bg-surface-900/80 backdrop-blur-md border border-surface-800 text-white shadow-elevated">
    <!-- Mute All -->
    <button
      type="button"
      class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold bg-rose-600/80 hover:bg-rose-600 transition-colors"
      title="Bisukan semua mikrofon siswa"
      @click="emit('mute-all')"
    >
      <MicOff class="w-3.5 h-3.5" />
      <span>Mute All</span>
    </button>

    <!-- Lock Room -->
    <button
      type="button"
      :class="[
        'flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors',
        isLocked ? 'bg-amber-600/80 hover:bg-amber-600' : 'bg-surface-800 hover:bg-surface-700'
      ]"
      @click="emit('toggle-lock')"
    >
      <component :is="isLocked ? Lock : Unlock" class="w-3.5 h-3.5" />
      <span>{{ isLocked ? 'Ruang Terkunci' : 'Kunci Ruangan' }}</span>
    </button>
  </div>
</template>

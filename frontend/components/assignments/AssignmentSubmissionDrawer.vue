<script setup lang="ts">
import { ref, computed } from 'vue'
import { X, UploadCloud, File, AlertCircle, CheckCircle, Clock } from 'lucide-vue-next'

interface Props {
  isOpen: boolean
  assignmentTitle: string
  dueDate?: string
  maxSizeBytes?: number
}

const props = withDefaults(defineProps<Props>(), {
  isOpen: false,
  assignmentTitle: 'Pengumpulan Tugas',
  maxSizeBytes: 50 * 1024 * 1024 // 50MB default
})

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', payload: { file: File; notes: string }): void
}>()

const selectedFile = ref<File | null>(null)
const notes = ref('')
const error = ref('')

const isLate = computed(() => {
  if (!props.dueDate) return false
  return new Date() > new Date(props.dueDate)
})

function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    if (file.size > props.maxSizeBytes) {
      error.value = 'Ukuran berkas melebihi batas maksimal 50MB'
      selectedFile.value = null
      return
    }
    error.value = ''
    selectedFile.value = file
  }
}

function submit() {
  if (!selectedFile.value) {
    error.value = 'Silakan pilih berkas pengumpulan'
    return
  }
  emit('submit', { file: selectedFile.value, notes: notes.value })
  emit('close')
}
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex justify-end bg-black/40 backdrop-blur-xs transition-opacity">
    <div class="w-full max-w-md bg-white dark:bg-surface-900 h-full shadow-2xl flex flex-col p-6 overflow-y-auto">
      <!-- Drawer Header -->
      <div class="flex items-center justify-between pb-4 border-b border-surface-100 dark:border-surface-800">
        <div>
          <h3 class="text-base font-bold text-surface-900 dark:text-surface-100">Kirimkan Tugas</h3>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5 truncate">{{ assignmentTitle }}</p>
        </div>
        <button type="button" class="p-1.5 rounded-lg text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800" @click="emit('close')">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Late status warning -->
      <div v-if="isLate" class="mt-4 p-3 rounded-xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-900/50 flex items-center gap-2.5 text-xs text-rose-700 dark:text-rose-300">
        <Clock class="w-4 h-4 shrink-0 text-rose-600 dark:text-rose-400" />
        <span>Batas waktu terlewati. Pengumpulan ini akan ditandai sebagai terlambat.</span>
      </div>

      <!-- Upload Zone -->
      <div class="mt-6">
        <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-2">Unggah Berkas (PDF / ZIP maks 50MB)</label>
        <div class="border-2 border-dashed border-surface-200 dark:border-surface-800 rounded-2xl p-6 text-center hover:border-brand-400 transition-colors bg-surface-50/50 dark:bg-surface-950/30 relative cursor-pointer">
          <input type="file" class="absolute inset-0 opacity-0 cursor-pointer" @change="handleFileSelect" />
          <UploadCloud class="w-8 h-8 text-surface-400 mx-auto mb-2" />
          <p class="text-xs font-medium text-surface-700 dark:text-surface-300">Tarik berkas ke sini atau klik untuk memilih</p>
          <p class="text-[11px] text-surface-400 mt-1">Dokumen tugas akan disimpan dengan aman</p>
        </div>

        <div v-if="selectedFile" class="mt-3 flex items-center justify-between p-2.5 rounded-xl bg-surface-100 dark:bg-surface-800 text-xs">
          <div class="flex items-center gap-2 truncate">
            <File class="w-4 h-4 text-brand-600 shrink-0" />
            <span class="font-medium text-surface-800 dark:text-surface-200 truncate">{{ selectedFile.name }}</span>
          </div>
          <span class="text-surface-400 shrink-0 text-[11px]">{{ (selectedFile.size / 1024 / 1024).toFixed(2) }} MB</span>
        </div>

        <p v-if="error" class="mt-2 text-xs text-rose-500 flex items-center gap-1">
          <AlertCircle class="w-3.5 h-3.5" /> {{ error }}
        </p>
      </div>

      <!-- Notes -->
      <div class="mt-6 flex-1">
        <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-1.5">Catatan Tambahan (Opsional)</label>
        <textarea
          v-model="notes"
          rows="3"
          placeholder="Tulis pesan atau catatan untuk pengajar..."
          class="w-full px-3 py-2 rounded-xl text-xs border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-950 text-surface-800 dark:text-surface-200 focus:outline-none focus:ring-2 focus:ring-brand-500/20"
        />
      </div>

      <!-- Footer CTA -->
      <div class="pt-4 border-t border-surface-100 dark:border-surface-800 flex gap-2.5">
        <UiButton variant="outline" size="sm" class="flex-1" @click="emit('close')">Batal</UiButton>
        <UiButton variant="primary" size="sm" class="flex-1" @click="submit">Kirim Tugas</UiButton>
      </div>
    </div>
  </div>
</template>

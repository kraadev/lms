<script setup lang="ts">
import { ref } from 'vue'
import { X, CheckCircle, ChevronLeft, ChevronRight, Download, User } from 'lucide-vue-next'

export interface SubmissionItem {
  id: number
  studentName: string
  submittedAt: string
  fileUrl?: string
  score?: number
  feedback?: string
}

interface Props {
  isOpen: boolean
  submissions: SubmissionItem[]
}

const props = withDefaults(defineProps<Props>(), {
  isOpen: false,
  submissions: () => []
})

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'grade', payload: { submissionId: number; score: number; feedback: string }): void
}>()

const currentIndex = ref(0)
const scoreInput = ref<number | ''>('')
const feedbackInput = ref('')

function saveGrade() {
  const current = props.submissions[currentIndex.value]
  if (!current || scoreInput.value === '') return
  emit('grade', {
    submissionId: current.id,
    score: Number(scoreInput.value),
    feedback: feedbackInput.value
  })
  if (currentIndex.value < props.submissions.length - 1) {
    currentIndex.value++
    scoreInput.value = ''
    feedbackInput.value = ''
  }
}
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex justify-end bg-black/40 backdrop-blur-xs transition-opacity">
    <div class="w-full max-w-lg bg-white dark:bg-surface-900 h-full shadow-2xl flex flex-col p-6 overflow-y-auto">
      <div class="flex items-center justify-between pb-4 border-b border-surface-100 dark:border-surface-800">
        <div>
          <h3 class="text-base font-bold text-surface-900 dark:text-surface-100">Speed-Grader</h3>
          <p class="text-xs text-surface-500 dark:text-surface-400">Penilaian cepat tugas siswa</p>
        </div>
        <button type="button" class="p-1.5 rounded-lg text-surface-400 hover:bg-surface-100 dark:hover:bg-surface-800" @click="emit('close')">
          <X class="w-5 h-5" />
        </button>
      </div>

      <div v-if="submissions.length" class="mt-4 flex items-center justify-between p-3 rounded-xl bg-surface-50 dark:bg-surface-800/50">
        <div class="flex items-center gap-2">
          <User class="w-4 h-4 text-brand-600" />
          <span class="text-xs font-semibold text-surface-900 dark:text-surface-100">{{ submissions[currentIndex]?.studentName }}</span>
        </div>
        <div class="flex items-center gap-1.5">
          <button
            :disabled="currentIndex === 0"
            class="p-1 rounded text-surface-400 disabled:opacity-30 hover:bg-surface-200 dark:hover:bg-surface-700"
            @click="currentIndex--"
          >
            <ChevronLeft class="w-4 h-4" />
          </button>
          <span class="text-[11px] text-surface-500">{{ currentIndex + 1 }} / {{ submissions.length }}</span>
          <button
            :disabled="currentIndex === submissions.length - 1"
            class="p-1 rounded text-surface-400 disabled:opacity-30 hover:bg-surface-200 dark:hover:bg-surface-700"
            @click="currentIndex++"
          >
            <ChevronRight class="w-4 h-4" />
          </button>
        </div>
      </div>

      <div class="mt-6 flex-1 space-y-4">
        <div>
          <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-1">Nilai (Skala 0 - 100)</label>
          <input
            v-model="scoreInput"
            type="number"
            min="0"
            max="100"
            placeholder="Misal: 85"
            class="w-full px-3 py-2 rounded-xl text-sm font-bold border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-950 text-surface-900 dark:text-surface-100"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-1">Komentar & Masukan</label>
          <textarea
            v-model="feedbackInput"
            rows="4"
            placeholder="Berikan masukan yang membangun untuk siswa..."
            class="w-full px-3 py-2 rounded-xl text-xs border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-950 text-surface-800 dark:text-surface-200"
          />
        </div>
      </div>

      <div class="pt-4 border-t border-surface-100 dark:border-surface-800 flex gap-2.5">
        <UiButton variant="outline" size="sm" class="flex-1" @click="emit('close')">Selesai</UiButton>
        <UiButton variant="primary" size="sm" class="flex-1" @click="saveGrade">Simpan & Lanjut</UiButton>
      </div>
    </div>
  </div>
</template>

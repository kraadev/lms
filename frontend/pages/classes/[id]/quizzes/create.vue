<script setup lang="ts">
import {
  FileQuestion, Plus, Trash2, ArrowLeft, Check, AlertCircle,
  HelpCircle, Sparkles
} from 'lucide-vue-next'
import { quizzesService, type CreateQuizPayload } from '~/services/quizzes'

definePageMeta({ middleware: 'auth' })

const route = useRoute()
const classId = computed(() => route.params.id as string)
const auth = useAuthStore()
const toast = useToast()

useSeoMeta({ title: 'Buat Kuis Baru' })

const form = reactive<CreateQuizPayload>({
  title: '',
  description: '',
  duration_minutes: 30,
  attempts_allowed: 1,
  due_date: '',
  questions: [
    {
      question_text: '',
      question_type: 'multiple_choice',
      points: 10,
      order: 1,
      options: [
        { option_text: '', order: 1, is_correct: true },
        { option_text: '', order: 2, is_correct: false },
        { option_text: '', order: 3, is_correct: false },
        { option_text: '', order: 4, is_correct: false }
      ]
    }
  ]
})

const isSaving = ref(false)

function addQuestion(type: 'multiple_choice' | 'essay' = 'multiple_choice') {
  form.questions.push({
    question_text: '',
    question_type: type,
    points: 10,
    order: form.questions.length + 1,
    options: type === 'multiple_choice' ? [
      { option_text: '', order: 1, is_correct: true },
      { option_text: '', order: 2, is_correct: false },
      { option_text: '', order: 3, is_correct: false },
      { option_text: '', order: 4, is_correct: false }
    ] : undefined
  })
}

function removeQuestion(index: number) {
  if (form.questions.length <= 1) {
    toast.error('Kuis harus memiliki minimal 1 soal')
    return
  }
  form.questions.splice(index, 1)
  form.questions.forEach((q, i) => { q.order = i + 1 })
}

function addOption(questionIndex: number) {
  const q = form.questions[questionIndex]
  if (!q.options) q.options = []
  q.options.push({
    option_text: '',
    order: q.options.length + 1,
    is_correct: false
  })
}

function removeOption(questionIndex: number, optionIndex: number) {
  const q = form.questions[questionIndex]
  if (!q.options || q.options.length <= 2) {
    toast.error('Pilihan ganda minimal memiliki 2 opsi')
    return
  }
  const wasCorrect = q.options[optionIndex].is_correct
  q.options.splice(optionIndex, 1)
  q.options.forEach((opt, idx) => { opt.order = idx + 1 })
  if (wasCorrect && q.options.length > 0) {
    q.options[0].is_correct = true
  }
}

function setCorrectOption(questionIndex: number, optionIndex: number) {
  const q = form.questions[questionIndex]
  if (!q.options) return
  q.options.forEach((opt, idx) => {
    opt.is_correct = idx === optionIndex
  })
}

const totalPoints = computed(() => {
  return form.questions.reduce((sum, q) => sum + (Number(q.points) || 0), 0)
})

function validate(): boolean {
  if (!form.title.trim()) {
    toast.error('Judul kuis wajib diisi')
    return false
  }

  for (let i = 0; i < form.questions.length; i++) {
    const q = form.questions[i]
    if (!q.question_text.trim()) {
      toast.error(`Pertanyaan #${i + 1} belum diisi`)
      return false
    }

    if (q.question_type === 'multiple_choice') {
      if (!q.options || q.options.length < 2) {
        toast.error(`Pertanyaan #${i + 1} harus memiliki minimal 2 opsi`)
        return false
      }
      for (let j = 0; j < q.options.length; j++) {
        if (!q.options[j].option_text.trim()) {
          toast.error(`Opsi pilihan ${String.fromCharCode(65 + j)} pada soal #${i + 1} belum diisi`)
          return false
        }
      }
      const hasCorrect = q.options.some(o => o.is_correct)
      if (!hasCorrect) {
        toast.error(`Pilih 1 kunci jawaban benar pada soal #${i + 1}`)
        return false
      }
    }
  }

  return true
}

async function handleSaveQuiz() {
  if (!validate()) return
  isSaving.value = true
  try {
    const created = await quizzesService.create(classId.value, form)
    toast.success('Kuis berhasil dibuat!')
    await navigateTo(`/classes/${classId.value}`)
  } catch (err: any) {
    toast.error(err?.message || 'Gagal membuat kuis')
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-4xl mx-auto space-y-6">
    <div>
      <NuxtLink
        :to="`/classes/${classId}`"
        class="inline-flex items-center gap-1.5 text-xs font-medium text-surface-500 hover:text-surface-900 dark:hover:text-surface-100 transition-colors"
      >
        <ArrowLeft class="w-4 h-4" />
        Kembali ke Kelas
      </NuxtLink>
    </div>

    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">Pembuat Kuis Baru</h1>
        <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Rancang kuis pilihan ganda atau essay secara terstruktur.</p>
      </div>

      <div class="flex items-center gap-3">
        <div class="px-3 py-1.5 rounded-xl bg-brand-50 dark:bg-brand-950/40 border border-brand-200 dark:border-brand-800 text-xs font-semibold text-brand-700 dark:text-brand-300">
          Total: {{ totalPoints }} Poin
        </div>
        <UiButton :loading="isSaving" @click="handleSaveQuiz">
          Publikasikan Kuis
        </UiButton>
      </div>
    </div>

    <!-- Quiz Settings Card -->
    <div class="p-6 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft space-y-4">
      <h2 class="text-sm font-bold uppercase tracking-wider text-surface-500 mb-2">Informasi Umum</h2>

      <div class="space-y-4">
        <UiInput
          v-model="form.title"
          label="Judul Kuis"
          placeholder="Contoh: Kuis 1 - Aljabar Linear"
          required
        />

        <UiTextarea
          v-model="form.description"
          label="Deskripsi / Petunjuk Umum"
          placeholder="Tuliskan petunjuk untuk siswa..."
          :rows="2"
        />

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <UiInput
            v-model.number="form.duration_minutes"
            type="number"
            label="Durasi (Menit)"
            min="1"
            required
          />

          <UiInput
            v-model.number="form.attempts_allowed"
            type="number"
            label="Maksimal Percobaan"
            min="1"
            required
          />

          <UiInput
            v-model="form.due_date"
            type="datetime-local"
            label="Tenggat Waktu (Opsional)"
          />
        </div>
      </div>
    </div>

    <!-- Questions Builder List -->
    <div class="space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-sm font-bold uppercase tracking-wider text-surface-500">
          Daftar Pertanyaan ({{ form.questions.length }})
        </h2>

        <div class="flex items-center gap-2">
          <UiButton size="xs" variant="outline" @click="addQuestion('multiple_choice')">
            <Plus class="w-3.5 h-3.5 mr-1" />
            + Pilihan Ganda
          </UiButton>
          <UiButton size="xs" variant="outline" @click="addQuestion('essay')">
            <Plus class="w-3.5 h-3.5 mr-1" />
            + Essay
          </UiButton>
        </div>
      </div>

      <!-- Question Cards -->
      <div
        v-for="(q, qIdx) in form.questions"
        :key="qIdx"
        class="p-6 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft space-y-4"
      >
        <div class="flex items-center justify-between pb-3 border-b border-surface-100 dark:border-surface-800">
          <div class="flex items-center gap-3">
            <span class="w-7 h-7 rounded-lg bg-brand-100 dark:bg-brand-950/60 text-brand-700 dark:text-brand-300 font-bold text-xs flex items-center justify-center">
              {{ qIdx + 1 }}
            </span>
            <span class="text-xs font-semibold text-surface-700 dark:text-surface-300">
              {{ q.question_type === 'multiple_choice' ? 'Pilihan Ganda' : 'Essay' }}
            </span>
          </div>

          <div class="flex items-center gap-3">
            <div class="flex items-center gap-1.5">
              <label class="text-xs text-surface-400">Poin:</label>
              <input
                v-model.number="q.points"
                type="number"
                min="1"
                class="w-16 px-2 py-1 bg-surface-50 dark:bg-surface-800 border border-surface-200 dark:border-surface-700 rounded-lg text-xs font-semibold text-center text-surface-900 dark:text-surface-100 focus:outline-none focus:border-brand-500"
              />
            </div>

            <button
              type="button"
              class="p-1.5 text-surface-400 hover:text-rose-600 dark:hover:text-rose-400 rounded-lg hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors"
              title="Hapus Soal"
              @click="removeQuestion(qIdx)"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Question Text -->
        <UiTextarea
          v-model="q.question_text"
          label="Teks Pertanyaan"
          placeholder="Tuliskan butir soal di sini..."
          :rows="3"
          required
        />

        <!-- Multiple Choice Options Editor -->
        <div v-if="q.question_type === 'multiple_choice' && q.options" class="space-y-3 pt-2">
          <label class="block text-xs font-medium text-surface-600 dark:text-surface-400">
            Pilihan Jawaban (Pilih radio untuk menentukan kunci jawaban yang benar):
          </label>

          <div class="space-y-2.5">
            <div
              v-for="(opt, optIdx) in q.options"
              :key="optIdx"
              class="flex items-center gap-3 p-2.5 rounded-xl border transition-all"
              :class="opt.is_correct
                ? 'border-emerald-500 bg-emerald-50/40 dark:bg-emerald-950/20'
                : 'border-surface-200 dark:border-surface-700 bg-surface-50/50 dark:bg-surface-800/40'"
            >
              <input
                type="radio"
                :name="`correct_${qIdx}`"
                :checked="opt.is_correct"
                class="w-4 h-4 text-emerald-600 focus:ring-emerald-500 cursor-pointer"
                @change="setCorrectOption(qIdx, optIdx)"
              />

              <span class="text-xs font-bold text-surface-500 w-4">{{ String.fromCharCode(65 + optIdx) }}.</span>

              <input
                v-model="opt.option_text"
                type="text"
                placeholder="Isi teks pilihan..."
                class="flex-1 bg-transparent border-0 text-sm text-surface-900 dark:text-surface-100 placeholder-surface-400 focus:outline-none focus:ring-0"
                required
              />

              <button
                type="button"
                class="text-surface-400 hover:text-rose-500 p-1 rounded transition-colors"
                title="Hapus Opsi"
                @click="removeOption(qIdx, optIdx)"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <div class="pt-1">
            <button
              type="button"
              class="text-xs font-medium text-brand-600 dark:text-brand-400 hover:underline flex items-center gap-1"
              @click="addOption(qIdx)"
            >
              <Plus class="w-3.5 h-3.5" /> Tambah Opsi Pilihan
            </button>
          </div>
        </div>
      </div>

      <!-- Add Question Button Box -->
      <div class="flex items-center justify-center gap-3 p-6 border-2 border-dashed border-surface-200 dark:border-surface-800 rounded-2xl">
        <UiButton variant="outline" size="sm" @click="addQuestion('multiple_choice')">
          <Plus class="w-4 h-4 mr-1.5" />
          Tambah Soal Pilihan Ganda
        </UiButton>
        <UiButton variant="outline" size="sm" @click="addQuestion('essay')">
          <Plus class="w-4 h-4 mr-1.5" />
          Tambah Soal Essay
        </UiButton>
      </div>
    </div>
  </div>
</template>

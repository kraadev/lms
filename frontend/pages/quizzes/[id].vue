<script setup lang="ts">
import {
  FileQuestion, Clock, AlertTriangle, ArrowLeft, ArrowRight,
  CheckCircle2, Award, Check, Play, RefreshCw
} from 'lucide-vue-next'
import { quizzesService } from '~/services/quizzes'
import type { Quiz, QuizQuestion, QuizAttempt, QuizAttemptAnswer } from '~/types'
import { formatDate } from '~/utils/formatters'

definePageMeta({ middleware: 'auth' })

const route = useRoute()
const quizId = computed(() => route.params.id as string)
const auth = useAuthStore()
const toast = useToast()

const quiz = ref<(Quiz & { questions?: QuizQuestion[] }) | null>(null)
const currentAttempt = ref<(QuizAttempt & { questions: QuizQuestion[] }) | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

// Taking quiz state
const activeQuestionIndex = ref(0)
const answers = reactive<Record<number, { option_id?: number; essay_answer?: string }>>({})
const remainingSeconds = ref(0)
let timerInterval: any = null
const isSubmitting = ref(false)
const showConfirmSubmitModal = ref(false)

useSeoMeta({ title: computed(() => quiz.value?.title || 'Kuis') })

async function loadQuiz() {
  isLoading.value = true
  error.value = null
  try {
    quiz.value = await quizzesService.getById(quizId.value)
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat kuis'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadQuiz)

onBeforeUnmount(() => {
  if (timerInterval) clearInterval(timerInterval)
})

function startTimer(durationMinutes: number) {
  remainingSeconds.value = durationMinutes * 60
  if (timerInterval) clearInterval(timerInterval)

  timerInterval = setInterval(() => {
    if (remainingSeconds.value > 0) {
      remainingSeconds.value -= 1
    } else {
      clearInterval(timerInterval)
      toast.warning('Waktu pengerjaan habis! Kuis akan dikumpulkan otomatis.')
      handleFinishQuiz()
    }
  }, 1000)
}

const formattedTimer = computed(() => {
  const m = Math.floor(remainingSeconds.value / 60)
  const s = remainingSeconds.value % 60
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
})

async function handleStartAttempt() {
  isLoading.value = true
  try {
    const attempt = await quizzesService.startAttempt(quizId.value)
    currentAttempt.value = attempt
    activeQuestionIndex.value = 0

    // Initialize answer model
    if (attempt.questions) {
      attempt.questions.forEach(q => {
        answers[q.id] = {}
      })
    }

    startTimer(quiz.value?.duration_minutes || 30)
  } catch (err: any) {
    toast.error(err?.message || 'Gagal memulai kuis')
  } finally {
    isLoading.value = false
  }
}

const currentQuestion = computed<QuizQuestion | null>(() => {
  if (!currentAttempt.value?.questions?.length) return null
  return currentAttempt.value.questions[activeQuestionIndex.value] || null
})

const answeredCount = computed(() => {
  if (!currentAttempt.value?.questions) return 0
  return currentAttempt.value.questions.filter(q => {
    const a = answers[q.id]
    return a && (a.option_id !== undefined || (a.essay_answer && a.essay_answer.trim()))
  }).length
})

async function handleFinishQuiz() {
  if (!currentAttempt.value || isSubmitting.value) return
  isSubmitting.value = true
  showConfirmSubmitModal.value = false

  if (timerInterval) clearInterval(timerInterval)

  try {
    const formattedAnswers: QuizAttemptAnswer[] = Object.entries(answers).map(([qId, ans]) => ({
      question_id: Number(qId),
      option_id: ans.option_id,
      essay_answer: ans.essay_answer
    }))

    const result = await quizzesService.submitAttempt(currentAttempt.value.id, formattedAnswers)
    currentAttempt.value = null
    toast.success('Kuis berhasil dikumpulkan!')
    await loadQuiz()
  } catch (err: any) {
    toast.error(err?.message || 'Gagal mengumpulkan kuis')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-5xl mx-auto space-y-6">
    <!-- Back link -->
    <div v-if="!currentAttempt">
      <NuxtLink
        :to="quiz?.class_id ? `/classes/${quiz.class_id}` : '/quizzes'"
        class="inline-flex items-center gap-1.5 text-xs font-medium text-surface-500 hover:text-surface-900 dark:hover:text-surface-100 transition-colors"
      >
        <ArrowLeft class="w-4 h-4" />
        Kembali ke {{ quiz?.class_title ? quiz.class_title : 'Daftar' }}
      </NuxtLink>
    </div>

    <div v-if="isLoading">
      <UiSkeleton class="h-40 rounded-2xl mb-4" />
      <UiSkeleton :rows="5" />
    </div>

    <UiErrorState v-else-if="error" :message="error" @retry="loadQuiz" />

    <!-- ================= 1. QUIZ TAKING IN PROGRESS ================= -->
    <div v-else-if="currentAttempt && currentAttempt.questions?.length" class="space-y-6">
      <!-- Sticky Quiz Header with Timer -->
      <div class="sticky top-2 z-20 p-4 bg-white/95 dark:bg-surface-900/95 backdrop-blur-md rounded-2xl border border-surface-200 dark:border-surface-800 shadow-elevated flex items-center justify-between gap-4">
        <div>
          <h2 class="text-sm font-bold text-surface-900 dark:text-surface-100 line-clamp-1">{{ quiz?.title }}</h2>
          <p class="text-xs text-surface-500">Soal {{ activeQuestionIndex + 1 }} dari {{ currentAttempt.questions.length }} &middot; Terjawab {{ answeredCount }}/{{ currentAttempt.questions.length }}</p>
        </div>

        <div class="flex items-center gap-3 shrink-0">
          <div
            class="flex items-center gap-2 px-3.5 py-1.5 rounded-xl font-mono text-sm font-bold shadow-sm"
            :class="remainingSeconds < 180
              ? 'bg-rose-100 text-rose-700 dark:bg-rose-950/60 dark:text-rose-300 animate-pulse'
              : 'bg-surface-100 text-surface-800 dark:bg-surface-800 dark:text-surface-200'"
          >
            <Clock class="w-4 h-4" />
            <span>{{ formattedTimer }}</span>
          </div>

          <UiButton size="sm" variant="danger" @click="showConfirmSubmitModal = true">
            Selesai
          </UiButton>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <!-- Main Question Area (3 cols) -->
        <div class="lg:col-span-3 space-y-6">
          <div v-if="currentQuestion" class="p-6 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft space-y-6">
            <div class="flex items-center justify-between pb-3 border-b border-surface-100 dark:border-surface-800 text-xs">
              <span class="font-bold text-brand-600 dark:text-brand-400">Pertanyaan #{{ activeQuestionIndex + 1 }}</span>
              <span class="text-surface-400">{{ currentQuestion.points }} Poin</span>
            </div>

            <!-- Question Text -->
            <p class="text-base text-surface-900 dark:text-surface-100 font-medium leading-relaxed whitespace-pre-line">
              {{ currentQuestion.question_text }}
            </p>

            <!-- Multiple Choice Options -->
            <div v-if="currentQuestion.question_type === 'multiple_choice' && currentQuestion.options?.length" class="space-y-3 pt-2">
              <label
                v-for="(opt, idx) in currentQuestion.options"
                :key="opt.id"
                class="flex items-start gap-3.5 p-4 rounded-xl border cursor-pointer transition-all"
                :class="answers[currentQuestion.id]?.option_id === opt.id
                  ? 'border-brand-500 bg-brand-50/50 dark:bg-brand-950/30 ring-2 ring-brand-500/20 text-surface-900 dark:text-white'
                  : 'border-surface-200 dark:border-surface-700 bg-surface-50/50 dark:bg-surface-800/40 hover:bg-surface-100/70 text-surface-700 dark:text-surface-300'"
              >
                <input
                  type="radio"
                  :name="`q_${currentQuestion.id}`"
                  :value="opt.id"
                  :checked="answers[currentQuestion.id]?.option_id === opt.id"
                  class="mt-1 w-4 h-4 text-brand-600 focus:ring-brand-500"
                  @change="answers[currentQuestion.id] = { option_id: opt.id }"
                />
                <div class="flex-1 text-sm leading-relaxed">
                  <span class="font-semibold mr-1.5">{{ String.fromCharCode(65 + idx) }}.</span>
                  {{ opt.option_text }}
                </div>
              </label>
            </div>

            <!-- Essay Response -->
            <div v-else-if="currentQuestion.question_type === 'essay'" class="pt-2">
              <UiTextarea
                v-model="answers[currentQuestion.id].essay_answer"
                label="Jawaban Anda"
                placeholder="Tuliskan jawaban essay selengkapnya di sini..."
                :rows="6"
              />
            </div>

            <!-- Question Navigation Bottom Bar -->
            <div class="flex items-center justify-between pt-4 border-t border-surface-100 dark:border-surface-800">
              <UiButton
                variant="outline"
                size="sm"
                :disabled="activeQuestionIndex === 0"
                @click="activeQuestionIndex -= 1"
              >
                <ArrowLeft class="w-4 h-4 mr-1" />
                Sebelumnya
              </UiButton>

              <UiButton
                v-if="activeQuestionIndex < currentAttempt.questions.length - 1"
                size="sm"
                @click="activeQuestionIndex += 1"
              >
                Berikutnya
                <ArrowRight class="w-4 h-4 ml-1" />
              </UiButton>

              <UiButton
                v-else
                variant="success"
                size="sm"
                @click="showConfirmSubmitModal = true"
              >
                Kumpulkan Kuis
              </UiButton>
            </div>
          </div>
        </div>

        <!-- Question Number Navigator Sidebar (1 col) -->
        <div class="space-y-4">
          <div class="p-5 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft">
            <h3 class="text-xs font-bold uppercase tracking-wider text-surface-500 mb-3">Navigasi Soal</h3>
            <div class="grid grid-cols-5 gap-2">
              <button
                v-for="(q, idx) in currentAttempt.questions"
                :key="q.id"
                type="button"
                class="w-full aspect-square rounded-lg text-xs font-bold transition-all flex items-center justify-center"
                :class="{
                  'ring-2 ring-brand-500 ring-offset-2 dark:ring-offset-surface-900': activeQuestionIndex === idx,
                  'bg-brand-600 text-white': answers[q.id]?.option_id !== undefined || (answers[q.id]?.essay_answer && answers[q.id]?.essay_answer.trim()),
                  'bg-surface-100 dark:bg-surface-800 text-surface-700 dark:text-surface-300': !(answers[q.id]?.option_id !== undefined || (answers[q.id]?.essay_answer && answers[q.id]?.essay_answer.trim()))
                }"
                @click="activeQuestionIndex = idx"
              >
                {{ idx + 1 }}
              </button>
            </div>

            <div class="mt-4 pt-4 border-t border-surface-100 dark:border-surface-800 space-y-2 text-[11px] text-surface-500">
              <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded bg-brand-600" />
                <span>Sudah Dijawab</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded bg-surface-200 dark:bg-surface-700" />
                <span>Belum Dijawab</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ================= 2. QUIZ OVERVIEW & SUMMARY (Not Taking) ================= -->
    <div v-else-if="quiz" class="space-y-6">
      <div class="p-6 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-surface-100 dark:border-surface-800">
          <div>
            <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100 mb-1">{{ quiz.title }}</h1>
            <p class="text-xs text-surface-500 dark:text-surface-400">
              {{ quiz.class_title || quiz.class?.title }} &middot; Dibuat {{ formatDate(quiz.created_at, { day: 'numeric', month: 'short', year: 'numeric' }) }}
            </p>
          </div>

          <div class="flex flex-wrap items-center gap-3">
            <div class="px-3 py-1.5 rounded-xl bg-surface-50 dark:bg-surface-800/80 border border-surface-200 dark:border-surface-700 text-xs">
              <span class="text-surface-400">Total Nilai: </span>
              <strong class="text-surface-900 dark:text-surface-100">{{ quiz.total_points }} Poin</strong>
            </div>
            <div class="px-3 py-1.5 rounded-xl bg-violet-50 dark:bg-violet-950/30 border border-violet-200 dark:border-violet-900/40 text-xs text-violet-800 dark:text-violet-200 flex items-center gap-1.5">
              <Clock class="w-3.5 h-3.5 text-violet-600 dark:text-violet-400" />
              <span>Durasi: <strong>{{ quiz.duration_minutes }} Menit</strong></span>
            </div>
          </div>
        </div>

        <div class="pt-5 space-y-4">
          <p class="text-sm text-surface-600 dark:text-surface-300 leading-relaxed whitespace-pre-line">
            {{ quiz.description || 'Kuis untuk menguji pemahaman materi.' }}
          </p>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 pt-2 text-xs">
            <div class="p-3 bg-surface-50 dark:bg-surface-800/50 rounded-xl">
              <span class="text-surface-400 block mb-0.5">Jumlah Soal</span>
              <strong class="text-surface-800 dark:text-surface-200 text-sm">{{ quiz.total_questions || quiz.questions?.length || 0 }} Butir</strong>
            </div>
            <div class="p-3 bg-surface-50 dark:bg-surface-800/50 rounded-xl">
              <span class="text-surface-400 block mb-0.5">Kesempatan Mengerjakan</span>
              <strong class="text-surface-800 dark:text-surface-200 text-sm">{{ quiz.attempts_allowed || 1 }} Kali</strong>
            </div>
            <div class="p-3 bg-surface-50 dark:bg-surface-800/50 rounded-xl">
              <span class="text-surface-400 block mb-0.5">Tenggat Waktu</span>
              <strong class="text-surface-800 dark:text-surface-200 text-sm">
                {{ quiz.due_date ? formatDate(quiz.due_date, { day: 'numeric', month: 'short' }) : 'Tidak ada tenggat' }}
              </strong>
            </div>
          </div>

          <!-- Latest Result if Student already took the quiz -->
          <div
            v-if="auth.isStudent && quiz.my_latest_attempt?.score !== null && quiz.my_latest_attempt?.score !== undefined"
            class="mt-4 p-5 bg-gradient-to-br from-emerald-50 to-teal-50 dark:from-emerald-950/30 dark:to-teal-950/20 border border-emerald-200 dark:border-emerald-800/50 rounded-2xl flex items-center justify-between"
          >
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-xl bg-emerald-500 text-white flex items-center justify-center shadow-md shadow-emerald-500/20">
                <Award class="w-6 h-6" />
              </div>
              <div>
                <p class="text-xs font-bold text-emerald-800 dark:text-emerald-300 uppercase tracking-wider">Hasil Pengerjaan Terakhir</p>
                <p class="text-xs text-emerald-600 dark:text-emerald-400 mt-0.5">
                  Selesai pada {{ formatDate(quiz.my_latest_attempt.submitted_at, { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' }) }}
                </p>
              </div>
            </div>

            <div class="text-right">
              <div class="text-3xl font-black text-emerald-700 dark:text-emerald-300">
                {{ quiz.my_latest_attempt.score }} <span class="text-sm font-normal text-emerald-600/80">/ {{ quiz.total_points }}</span>
              </div>
            </div>
          </div>

          <!-- Action Button -->
          <div v-if="auth.isStudent" class="pt-4 flex justify-end">
            <UiButton size="lg" @click="handleStartAttempt">
              <Play class="w-4.5 h-4.5 mr-2 fill-current" />
              {{ quiz.my_latest_attempt ? 'Kerjakan Ulang' : 'Mulai Kerjakan Kuis' }}
            </UiButton>
          </div>
        </div>
      </div>

      <!-- Question Preview for Teachers -->
      <div v-if="(auth.isTeacher || auth.isAdmin) && quiz.questions?.length" class="space-y-4">
        <h3 class="text-base font-semibold text-surface-900 dark:text-surface-100">
          Daftar Soal Kuis ({{ quiz.questions.length }})
        </h3>

        <div class="space-y-3">
          <div
            v-for="(q, idx) in quiz.questions"
            :key="q.id"
            class="p-5 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 space-y-3"
          >
            <div class="flex items-center justify-between text-xs text-surface-500">
              <span class="font-bold text-brand-600">Soal #{{ idx + 1 }} &middot; {{ q.question_type === 'multiple_choice' ? 'Pilihan Ganda' : 'Essay' }}</span>
              <span>{{ q.points }} Poin</span>
            </div>

            <p class="text-sm font-medium text-surface-900 dark:text-surface-100">{{ q.question_text }}</p>

            <div v-if="q.options?.length" class="space-y-1.5 pl-2">
              <div
                v-for="(opt, optIdx) in q.options"
                :key="opt.id"
                class="text-xs flex items-center gap-2 p-2 rounded-lg"
                :class="opt.is_correct ? 'bg-emerald-50 dark:bg-emerald-950/30 text-emerald-800 dark:text-emerald-300 font-semibold' : 'text-surface-600 dark:text-surface-400'"
              >
                <span>{{ String.fromCharCode(65 + optIdx) }}.</span>
                <span>{{ opt.option_text }}</span>
                <Check v-if="opt.is_correct" class="w-3.5 h-3.5 text-emerald-600 ml-auto" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Confirm Submit Modal -->
    <UiModal :show="showConfirmSubmitModal" title="Kumpulkan Kuis?" size="sm" @close="showConfirmSubmitModal = false">
      <div class="space-y-3 text-sm text-surface-600 dark:text-surface-400">
        <p>Anda telah menjawab <strong>{{ answeredCount }}</strong> dari <strong>{{ currentAttempt?.questions.length }}</strong> soal.</p>
        <p v-if="answeredCount < (currentAttempt?.questions.length || 0)" class="text-amber-600 dark:text-amber-400 text-xs font-medium">
          ⚠️ Masih ada soal yang belum dijawab. Apakah Anda yakin ingin menyelesaikan sekarang?
        </p>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="showConfirmSubmitModal = false">Kembali</UiButton>
          <UiButton variant="danger" :loading="isSubmitting" @click="handleFinishQuiz">Ya, Kumpulkan</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

<script setup lang="ts">
import {
  ClipboardList, Clock, ArrowLeft, Download, FileText,
  CheckCircle2, AlertCircle, Send, Award, MessageSquare, ChevronRight
} from 'lucide-vue-next'
import { assignmentsService } from '~/services/assignments'
import type { Assignment, AssignmentSubmission, SubmissionStatus } from '~/types'
import { formatDate, formatFileSize, formatRelativeTime } from '~/utils/formatters'

definePageMeta({ middleware: 'auth' })

const route = useRoute()
const assignmentId = computed(() => route.params.id as string)
const auth = useAuthStore()
const toast = useToast()

const assignment = ref<Assignment | null>(null)
const submissions = ref<AssignmentSubmission[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

// Student submission state
const submissionText = ref('')
const selectedFile = ref<File | null>(null)
const isSubmitting = ref(false)
const showResubmit = ref(false)

// Teacher grading state
const gradingSubmission = ref<AssignmentSubmission | null>(null)
const gradeForm = reactive({ score: 0, feedback: '' })
const isGrading = ref(false)

useSeoMeta({ title: computed(() => assignment.value?.title || 'Detail Tugas') })

async function loadData() {
  isLoading.value = true
  error.value = null
  try {
    assignment.value = await assignmentsService.getById(assignmentId.value)
    
    if (auth.isTeacher || auth.isAdmin) {
      const subData: any = await assignmentsService.getSubmissions(assignmentId.value)
      submissions.value = Array.isArray(subData) ? subData : (subData?.submissions || [])
    }
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat detail tugas'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadData)

function onFileSelect(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (files && files[0]) {
    selectedFile.value = files[0]
  }
}

async function handleStudentSubmit() {
  if (!submissionText.value.trim() && !selectedFile.value) {
    toast.error('Masukkan jawaban teks atau unggah file')
    return
  }

  isSubmitting.value = true
  try {
    const formData = new FormData()
    if (submissionText.value.trim()) formData.append('text_response', submissionText.value.trim())
    if (selectedFile.value) formData.append('file', selectedFile.value)

    const sub = await assignmentsService.submit(assignmentId.value, formData)
    if (assignment.value) {
      assignment.value.my_submission = sub
    }
    showResubmit.value = false
    submissionText.value = ''
    selectedFile.value = null
    toast.success('Tugas berhasil dikumpulkan!')
  } catch (err: any) {
    toast.error(err?.message || 'Gagal mengumpulkan tugas')
  } finally {
    isSubmitting.value = false
  }
}

function openGradeModal(sub: AssignmentSubmission) {
  gradingSubmission.value = sub
  gradeForm.score = sub.score || 0
  gradeForm.feedback = sub.feedback || ''
}

async function handleSaveGrade() {
  if (!gradingSubmission.value) return
  isGrading.value = true
  try {
    const updated = await assignmentsService.gradeSubmission(gradingSubmission.value.id, {
      score: Number(gradeForm.score),
      feedback: gradeForm.feedback
    })

    const idx = submissions.value.findIndex(s => s.id === gradingSubmission.value!.id)
    if (idx !== -1) submissions.value[idx] = updated

    gradingSubmission.value = null
    toast.success('Nilai berhasil disimpan!')
  } catch (err: any) {
    toast.error(err?.message || 'Gagal menyimpan nilai')
  } finally {
    isGrading.value = false
  }
}

const statusBadge: Record<SubmissionStatus, { variant: string; label: string }> = {
  not_submitted: { variant: 'warning', label: 'Belum Dikumpulkan' },
  submitted: { variant: 'info', label: 'Sudah Dikumpulkan' },
  late: { variant: 'danger', label: 'Terlambat' },
  graded: { variant: 'success', label: 'Dinilai' }
}

function isPastDue(dateStr: string): boolean {
  return new Date(dateStr) < new Date()
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-5xl mx-auto space-y-6">
    <!-- Back button -->
    <div>
      <NuxtLink
        :to="assignment?.class_id ? `/classes/${assignment.class_id}` : '/assignments'"
        class="inline-flex items-center gap-1.5 text-xs font-medium text-surface-500 hover:text-surface-900 dark:hover:text-surface-100 transition-colors"
      >
        <ArrowLeft class="w-4 h-4" />
        Kembali ke {{ assignment?.class_title ? assignment.class_title : 'Daftar' }}
      </NuxtLink>
    </div>

    <div v-if="isLoading">
      <UiSkeleton class="h-28 rounded-xl mb-4" />
      <UiSkeleton :rows="5" />
    </div>

    <UiErrorState v-else-if="error" :message="error" @retry="loadData" />

    <div v-else-if="assignment" class="space-y-6">
      <!-- Assignment Header Card -->
      <div class="p-6 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-surface-100 dark:border-surface-800">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">{{ assignment.title }}</h1>
            </div>
            <p class="text-xs text-surface-500 dark:text-surface-400">
              {{ assignment.class_title || assignment.class?.title }} &middot; Dibuat {{ formatDate(assignment.created_at, { day: 'numeric', month: 'short', year: 'numeric' }) }}
            </p>
          </div>

          <div class="flex flex-wrap items-center gap-3 shrink-0">
            <div class="px-3 py-1.5 rounded-xl bg-surface-50 dark:bg-surface-800/80 border border-surface-200 dark:border-surface-700 text-xs">
              <span class="text-surface-400">Nilai Maksimal: </span>
              <strong class="text-surface-900 dark:text-surface-100">{{ assignment.points }} Poin</strong>
            </div>
            <div class="px-3 py-1.5 rounded-xl bg-amber-50 dark:bg-amber-950/30 border border-amber-200 dark:border-amber-900/40 text-xs text-amber-800 dark:text-amber-200 flex items-center gap-1.5">
              <Clock class="w-3.5 h-3.5 text-amber-600 dark:text-amber-400" />
              <span>Tenggat: <strong>{{ formatDate(assignment.due_date, { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' }) }}</strong></span>
            </div>
          </div>
        </div>

        <!-- Instructions -->
        <div class="pt-5 space-y-3">
          <h2 class="text-sm font-semibold text-surface-800 dark:text-surface-200">Petunjuk Pengerjaan:</h2>
          <div class="text-sm text-surface-600 dark:text-surface-300 whitespace-pre-line leading-relaxed">
            {{ assignment.instructions || 'Tidak ada instruksi khusus.' }}
          </div>

          <!-- Assignment Attachment if any -->
          <div v-if="assignment.attachment_url" class="mt-4 pt-4 border-t border-surface-100 dark:border-surface-800">
            <h3 class="text-xs font-semibold text-surface-500 mb-2">Lampiran Soal:</h3>
            <a
              :href="assignment.attachment_url"
              target="_blank"
              rel="noopener"
              class="inline-flex items-center gap-2 px-3 py-2 rounded-xl bg-surface-50 dark:bg-surface-800/80 border border-surface-200 dark:border-surface-700 text-xs font-medium text-brand-600 dark:text-brand-400 hover:bg-surface-100 transition-colors"
            >
              <Download class="w-4 h-4" />
              <span>Unduh Berkas Lampiran</span>
            </a>
          </div>
        </div>
      </div>

      <!-- ================= STUDENT SECTION ================= -->
      <div v-if="auth.isStudent" class="space-y-6">
        <!-- Submission Status Summary if already submitted -->
        <div
          v-if="assignment.my_submission && !showResubmit"
          class="p-6 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft space-y-4"
        >
          <div class="flex items-center justify-between">
            <h3 class="text-base font-semibold text-surface-900 dark:text-surface-100 flex items-center gap-2">
              <CheckCircle2 class="w-5 h-5 text-emerald-600 dark:text-emerald-400" />
              Status Pengumpulan Tugas
            </h3>
            <UiBadge :variant="statusBadge[assignment.my_submission.status]?.variant as any">
              {{ statusBadge[assignment.my_submission.status]?.label }}
            </UiBadge>
          </div>

          <!-- Submission details -->
          <div class="p-4 rounded-xl bg-surface-50 dark:bg-surface-800/50 border border-surface-200/60 dark:border-surface-700/60 space-y-2 text-xs">
            <div class="flex justify-between">
              <span class="text-surface-500">Waktu Pengumpulan:</span>
              <span class="text-surface-800 dark:text-surface-200 font-medium">
                {{ formatDate(assignment.my_submission.submitted_at, { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) }}
              </span>
            </div>

            <div v-if="assignment.my_submission.file_url" class="flex justify-between items-center pt-2 border-t border-surface-200/50 dark:border-surface-700/50">
              <span class="text-surface-500">Berkas Terlampir:</span>
              <a :href="assignment.my_submission.file_url" target="_blank" rel="noopener" class="text-brand-600 dark:text-brand-400 font-medium hover:underline flex items-center gap-1">
                <Download class="w-3.5 h-3.5" />
                {{ assignment.my_submission.file_name || 'Unduh Berkas' }}
              </a>
            </div>

            <div v-if="assignment.my_submission.text_response" class="pt-2 border-t border-surface-200/50 dark:border-surface-700/50">
              <span class="text-surface-500 block mb-1">Jawaban Teks:</span>
              <p class="text-surface-800 dark:text-surface-200 text-sm whitespace-pre-wrap bg-white dark:bg-surface-900 p-3 rounded-lg border border-surface-200 dark:border-surface-700">
                {{ assignment.my_submission.text_response }}
              </p>
            </div>
          </div>

          <!-- Feedback & Score if graded -->
          <div
            v-if="assignment.my_submission.status === 'graded'"
            class="p-5 bg-gradient-to-br from-emerald-50 to-teal-50 dark:from-emerald-950/30 dark:to-teal-950/20 border border-emerald-200 dark:border-emerald-800/50 rounded-xl space-y-3"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2 text-emerald-800 dark:text-emerald-200 font-semibold text-sm">
                <Award class="w-5 h-5 text-emerald-600 dark:text-emerald-400" />
                Hasil Penilaian
              </div>
              <div class="text-2xl font-black text-emerald-700 dark:text-emerald-300">
                {{ assignment.my_submission.score }} <span class="text-xs font-normal text-emerald-600/80">/ {{ assignment.points }} Poin</span>
              </div>
            </div>

            <div v-if="assignment.my_submission.feedback" class="pt-3 border-t border-emerald-200/50 dark:border-emerald-800/50 text-xs">
              <span class="font-semibold text-emerald-900 dark:text-emerald-200 block mb-1">Catatan dari Guru:</span>
              <p class="text-emerald-800 dark:text-emerald-300 whitespace-pre-wrap">
                {{ assignment.my_submission.feedback }}
              </p>
            </div>
          </div>

          <!-- Resubmit button if not graded yet -->
          <div v-if="assignment.my_submission.status !== 'graded'" class="flex justify-end pt-2">
            <UiButton variant="outline" size="sm" @click="showResubmit = true">
              Kirim Ulang Jawaban
            </UiButton>
          </div>
        </div>

        <!-- Submission Form (Not submitted or Resubmitting) -->
        <div
          v-else
          class="p-6 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft space-y-4"
        >
          <div class="flex items-center justify-between">
            <h3 class="text-base font-semibold text-surface-900 dark:text-surface-100 flex items-center gap-2">
              <Send class="w-4.5 h-4.5 text-brand-600" />
              Kumpulkan Tugas
            </h3>
            <UiButton v-if="showResubmit" variant="ghost" size="xs" @click="showResubmit = false">Batal</UiButton>
          </div>

          <form @submit.prevent="handleStudentSubmit" class="space-y-4">
            <UiTextarea
              v-model="submissionText"
              label="Jawaban Teks (Opsional)"
              placeholder="Tuliskan jawaban langsung atau catatan untuk pengajar..."
              :rows="4"
            />

            <div>
              <label class="block mb-1.5 text-sm font-medium text-surface-700 dark:text-surface-300">Unggah Berkas Tugas (PDF, DOCX, ZIP, Gambar)</label>
              <input
                type="file"
                class="block w-full text-sm text-surface-600 dark:text-surface-400 file:mr-3 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-xs file:font-semibold file:bg-brand-50 file:text-brand-700 dark:file:bg-brand-950/50 dark:file:text-brand-300 hover:file:bg-brand-100"
                @change="onFileSelect"
              />
              <p v-if="selectedFile" class="mt-1 text-xs text-brand-600 dark:text-brand-400">
                Berkas dipilih: <strong>{{ selectedFile.name }}</strong> ({{ formatFileSize(selectedFile.size) }})
              </p>
            </div>

            <div class="pt-3 flex justify-end">
              <UiButton type="submit" :loading="isSubmitting" size="md">
                {{ isSubmitting ? 'Mengirim...' : 'Kirim Tugas Sekarang' }}
              </UiButton>
            </div>
          </form>
        </div>
      </div>

      <!-- ================= TEACHER SECTION (Grading & Submissions) ================= -->
      <div v-if="auth.isTeacher || auth.isAdmin" class="space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-base font-semibold text-surface-900 dark:text-surface-100">
            Daftar Pengumpulan Siswa ({{ (submissions || []).length }})
          </h3>
        </div>

        <div v-if="!submissions || !submissions.length">
          <UiEmptyState
            :icon="ClipboardList"
            title="Belum ada pengumpulan"
            description="Belum ada siswa yang mengumpulkan tugas ini."
          />
        </div>

        <div v-else class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft">
          <div
            v-for="sub in (submissions || [])"
            :key="sub.id"
            class="p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 hover:bg-surface-50 dark:hover:bg-surface-800/40 transition-colors"
          >
            <div class="flex items-center gap-3 min-w-0">
              <UiAvatar :name="sub.student?.name || 'Siswa'" :src="sub.student?.avatar" size="md" />
              <div class="min-w-0">
                <p class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate">{{ sub.student?.name }}</p>
                <div class="flex flex-wrap items-center gap-2 text-xs text-surface-500 dark:text-surface-400 mt-0.5">
                  <span>{{ formatRelativeTime(sub.submitted_at) }}</span>
                  <a v-if="sub.file_url" :href="sub.file_url" target="_blank" rel="noopener" class="text-brand-600 dark:text-brand-400 font-medium hover:underline flex items-center gap-0.5">
                    <Download class="w-3 h-3" />
                    Unduh Berkas
                  </a>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-3 w-full sm:w-auto justify-between sm:justify-end">
              <div class="text-right">
                <UiBadge :variant="statusBadge[sub.status]?.variant as any" size="sm">
                  {{ statusBadge[sub.status]?.label }}
                </UiBadge>
                <p v-if="sub.score !== null && sub.score !== undefined" class="text-xs font-bold text-emerald-600 dark:text-emerald-400 mt-0.5">
                  Nilai: {{ sub.score }}/{{ assignment.points }}
                </p>
              </div>

              <UiButton size="sm" variant="outline" @click="openGradeModal(sub)">
                {{ sub.status === 'graded' ? 'Ubah Nilai' : 'Beri Nilai' }}
              </UiButton>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Teacher Grade Modal -->
    <UiModal :show="!!gradingSubmission" title="Beri Nilai Tugas" size="md" @close="gradingSubmission = null">
      <div v-if="gradingSubmission" class="space-y-4">
        <div class="p-3 bg-surface-50 dark:bg-surface-800/60 rounded-xl text-xs space-y-1">
          <p>Siswa: <strong>{{ gradingSubmission.student?.name }}</strong></p>
          <p v-if="gradingSubmission.text_response" class="whitespace-pre-wrap pt-1 text-surface-600 dark:text-surface-300">
            <strong>Jawaban:</strong> {{ gradingSubmission.text_response }}
          </p>
          <div v-if="gradingSubmission.file_url" class="pt-1">
            <a :href="gradingSubmission.file_url" target="_blank" rel="noopener" class="text-brand-600 font-medium flex items-center gap-1 hover:underline">
              <Download class="w-3.5 h-3.5" /> Unduh Berkas Siswa
            </a>
          </div>
        </div>

        <UiInput
          v-model.number="gradeForm.score"
          type="number"
          label="Nilai Angka"
          :max="assignment?.points || 100"
          min="0"
          placeholder="0 - 100"
          required
        />

        <UiTextarea
          v-model="gradeForm.feedback"
          label="Catatan / Feedback untuk Siswa"
          placeholder="Bagus, perbaiki nomor 2..."
          :rows="3"
        />
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="gradingSubmission = null">Batal</UiButton>
          <UiButton :loading="isGrading" @click="handleSaveGrade">Simpan Nilai</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

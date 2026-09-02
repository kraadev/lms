<script setup lang="ts">
import { ClipboardList, Clock, Plus } from 'lucide-vue-next'
import { assignmentsService } from '~/services/assignments'
import type { Assignment, SubmissionStatus } from '~/types'
import { formatDate } from '~/utils/formatters'

const props = defineProps<{ classId: string }>()
const auth = useAuthStore()
const toast = useToast()

const assignments = ref<Assignment[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)
const showCreateModal = ref(false)
const createForm = reactive({ title: '', instructions: '', due_date: '', points: 100 })
const isSaving = ref(false)

async function load() {
  isLoading.value = true
  error.value = null
  try {
    assignments.value = await assignmentsService.getByClass(props.classId)
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat tugas'
  } finally {
    isLoading.value = false
  }
}

onMounted(load)

async function createAssignment() {
  if (!createForm.title || !createForm.due_date) return
  isSaving.value = true
  try {
    const a = await assignmentsService.create(props.classId, { ...createForm })
    assignments.value.unshift(a)
    showCreateModal.value = false
    Object.assign(createForm, { title: '', instructions: '', due_date: '', points: 100 })
    toast.success('Tugas berhasil dibuat!')
  } catch (err: any) {
    toast.error(err?.message || 'Gagal membuat tugas')
  } finally {
    isSaving.value = false
  }
}

const statusBadge: Record<SubmissionStatus, { variant: string; label: string }> = {
  not_submitted: { variant: 'warning', label: 'Belum Dikumpulkan' },
  submitted: { variant: 'info', label: 'Dikumpulkan' },
  late: { variant: 'danger', label: 'Terlambat' },
  graded: { variant: 'success', label: 'Dinilai' }
}

function isPastDue(dateStr: string): boolean {
  return new Date(dateStr) < new Date()
}
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="flex items-center justify-between mb-5">
      <h2 class="text-base font-semibold text-surface-800 dark:text-surface-200">Tugas</h2>
      <UiButton v-if="auth.isTeacher || auth.isAdmin" size="sm" @click="showCreateModal = true">
        <Plus class="w-4 h-4" />
        Buat Tugas
      </UiButton>
    </div>

    <div v-if="isLoading"><UiSkeleton :rows="4" /></div>
    <UiErrorState v-else-if="error" :message="error" @retry="load" />
    <UiEmptyState v-else-if="!assignments.length" :icon="ClipboardList" title="Belum ada tugas" :description="auth.isTeacher ? 'Buat tugas pertama untuk kelas ini.' : 'Tidak ada tugas yang diberikan.'" />

    <div v-else class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden">
      <NuxtLink
        v-for="a in assignments"
        :key="a.id"
        :to="`/assignments/${a.id}`"
        class="flex items-center gap-4 p-4 hover:bg-surface-50 dark:hover:bg-surface-800/60 transition-colors"
      >
        <div class="w-9 h-9 rounded-lg bg-amber-50 dark:bg-amber-950/30 flex items-center justify-center shrink-0">
          <ClipboardList class="w-4.5 h-4.5 text-amber-600 dark:text-amber-400" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate">{{ a.title }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 flex items-center gap-1 mt-0.5">
            <Clock class="w-3 h-3" />
            Tenggat: {{ formatDate(a.due_date, { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) }}
            <span v-if="isPastDue(a.due_date) && !a.my_submission?.status" class="text-rose-500 font-medium">· Lewat</span>
          </p>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <!-- Student view -->
          <template v-if="auth.isStudent && a.my_submission">
            <UiBadge :variant="statusBadge[a.my_submission.status]?.variant as any" size="sm">
              {{ statusBadge[a.my_submission.status]?.label }}
            </UiBadge>
            <span v-if="a.my_submission.score !== null && a.my_submission.score !== undefined" class="text-sm font-semibold text-emerald-600 dark:text-emerald-400">
              {{ a.my_submission.score }}/{{ a.points }}
            </span>
          </template>
          <UiBadge v-else-if="auth.isStudent" variant="warning" size="sm">Belum Dikumpulkan</UiBadge>
          <!-- Teacher view -->
          <template v-if="auth.isTeacher || auth.isAdmin">
            <span class="text-xs text-surface-500 dark:text-surface-400">{{ a.submissions_count || 0 }}/{{ a.total_students || 0 }} dikumpulkan</span>
          </template>
        </div>
      </NuxtLink>
    </div>

    <!-- Create Assignment Modal -->
    <UiModal :show="showCreateModal" title="Buat Tugas Baru" @close="showCreateModal = false">
      <form @submit.prevent="createAssignment" class="space-y-4">
        <UiInput v-model="createForm.title" label="Judul Tugas" placeholder="Nama tugas" required />
        <UiTextarea v-model="createForm.instructions" label="Instruksi" placeholder="Deskripsi detail tugas..." :rows="4" />
        <div class="grid grid-cols-2 gap-4">
          <UiInput v-model="createForm.due_date" type="datetime-local" label="Tenggat Waktu" required />
          <UiInput v-model.number="createForm.points" type="number" label="Nilai Maksimal" />
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="showCreateModal = false">Batal</UiButton>
          <UiButton :loading="isSaving" @click="createAssignment">Buat Tugas</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

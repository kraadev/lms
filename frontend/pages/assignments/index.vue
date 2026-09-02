<script setup lang="ts">
import { ClipboardList, Clock, CheckCircle2, AlertCircle, Filter } from 'lucide-vue-next'
import { assignmentsService } from '~/services/assignments'
import type { Assignment, SubmissionStatus } from '~/types'
import { formatDate } from '~/utils/formatters'

definePageMeta({ middleware: 'auth' })
useSeoMeta({ title: 'Daftar Tugas' })

const auth = useAuthStore()
const assignments = ref<Assignment[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)
const selectedFilter = ref<'all' | 'pending' | 'submitted' | 'graded'>('all')

async function loadAssignments() {
  isLoading.value = true
  error.value = null
  try {
    assignments.value = await assignmentsService.getAll()
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat daftar tugas'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadAssignments)

const filteredAssignments = computed(() => {
  if (selectedFilter.value === 'all') return assignments.value

  return assignments.value.filter(a => {
    if (auth.isStudent) {
      const status = a.my_submission?.status || 'not_submitted'
      if (selectedFilter.value === 'pending') return status === 'not_submitted' || status === 'late'
      if (selectedFilter.value === 'submitted') return status === 'submitted'
      if (selectedFilter.value === 'graded') return status === 'graded'
    }
    return true
  })
})

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
  <div class="p-4 md:p-6 max-w-5xl mx-auto space-y-6">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">Semua Tugas</h1>
        <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Daftar tugas dari seluruh kelas yang Anda ikuti.</p>
      </div>

      <!-- Filter tabs for students -->
      <div v-if="auth.isStudent" class="flex items-center p-1 bg-surface-100 dark:bg-surface-800/80 rounded-xl text-xs font-medium">
        <button
          v-for="tab in [
            { key: 'all', label: 'Semua' },
            { key: 'pending', label: 'Belum Selesai' },
            { key: 'submitted', label: 'Dikumpulkan' },
            { key: 'graded', label: 'Dinilai' }
          ]"
          :key="tab.key"
          type="button"
          class="px-3 py-1.5 rounded-lg transition-colors"
          :class="selectedFilter === tab.key
            ? 'bg-white dark:bg-surface-900 text-surface-900 dark:text-white shadow-sm'
            : 'text-surface-500 hover:text-surface-900 dark:hover:text-white'"
          @click="selectedFilter = tab.key as any"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div v-if="isLoading">
      <UiSkeleton :rows="5" />
    </div>

    <UiErrorState v-else-if="error" :message="error" @retry="loadAssignments" />

    <UiEmptyState
      v-else-if="!filteredAssignments.length"
      :icon="ClipboardList"
      title="Tidak ada tugas"
      description="Tidak ada tugas yang sesuai dengan filter saat ini."
    />

    <div v-else class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft">
      <NuxtLink
        v-for="a in filteredAssignments"
        :key="a.id"
        :to="`/assignments/${a.id}`"
        class="flex items-center gap-4 p-4 hover:bg-surface-50 dark:hover:bg-surface-800/60 transition-colors"
      >
        <div class="w-10 h-10 rounded-xl bg-amber-50 dark:bg-amber-950/30 flex items-center justify-center shrink-0">
          <ClipboardList class="w-5 h-5 text-amber-600 dark:text-amber-400" />
        </div>

        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <p class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate">{{ a.title }}</p>
            <span v-if="a.class_title || a.class?.title" class="px-2 py-0.5 rounded text-[11px] bg-surface-100 dark:bg-surface-800 text-surface-600 dark:text-surface-300 font-medium">
              {{ a.class_title || a.class?.title }}
            </span>
          </div>
          <p class="text-xs text-surface-500 dark:text-surface-400 flex items-center gap-1.5 mt-1">
            <Clock class="w-3.5 h-3.5" />
            Tenggat: {{ formatDate(a.due_date, { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) }}
            <span v-if="isPastDue(a.due_date) && (!a.my_submission || a.my_submission.status === 'not_submitted')" class="text-rose-500 font-medium">
              &middot; Lewat Tenggat
            </span>
          </p>
        </div>

        <div class="flex items-center gap-2.5 shrink-0">
          <!-- Student status -->
          <template v-if="auth.isStudent">
            <UiBadge
              v-if="a.my_submission"
              :variant="statusBadge[a.my_submission.status]?.variant as any"
              size="sm"
            >
              {{ statusBadge[a.my_submission.status]?.label }}
            </UiBadge>
            <UiBadge v-else variant="warning" size="sm">Belum Dikumpulkan</UiBadge>

            <span v-if="a.my_submission?.score !== null && a.my_submission?.score !== undefined" class="text-sm font-bold text-emerald-600 dark:text-emerald-400">
              {{ a.my_submission.score }}/{{ a.points }}
            </span>
          </template>

          <!-- Teacher stats -->
          <template v-else>
            <span class="text-xs text-surface-500 dark:text-surface-400">
              {{ a.submissions_count || 0 }}/{{ a.total_students || 0 }} terkumpul
            </span>
          </template>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>

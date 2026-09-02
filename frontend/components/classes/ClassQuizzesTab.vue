<script setup lang="ts">
import { FileQuestion, Clock, Plus } from 'lucide-vue-next'
import { quizzesService } from '~/services/quizzes'
import type { Quiz } from '~/types'
import { formatDate } from '~/utils/formatters'

const props = defineProps<{ classId: string }>()
const auth = useAuthStore()
const toast = useToast()

const quizzes = ref<Quiz[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

async function load() {
  isLoading.value = true
  error.value = null
  try {
    quizzes.value = await quizzesService.getByClass(props.classId)
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat kuis'
  } finally {
    isLoading.value = false
  }
}

onMounted(load)

const statusBadge: Record<string, { variant: string; label: string }> = {
  upcoming: { variant: 'info', label: 'Akan Datang' },
  active: { variant: 'success', label: 'Sedang Berlangsung' },
  closed: { variant: 'default', label: 'Ditutup' }
}
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="flex items-center justify-between mb-5">
      <h2 class="text-base font-semibold text-surface-800 dark:text-surface-200">Kuis</h2>
      <NuxtLink v-if="auth.isTeacher || auth.isAdmin" :to="`/classes/${classId}/quizzes/create`">
        <UiButton size="sm">
          <Plus class="w-4 h-4" />
          Buat Kuis
        </UiButton>
      </NuxtLink>
    </div>

    <div v-if="isLoading"><UiSkeleton :rows="4" /></div>
    <UiErrorState v-else-if="error" :message="error" @retry="load" />
    <UiEmptyState v-else-if="!quizzes.length" :icon="FileQuestion" title="Belum ada kuis" :description="auth.isTeacher ? 'Buat kuis untuk menguji pemahaman siswa.' : 'Belum ada kuis di kelas ini.'" />

    <div v-else class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden">
      <NuxtLink
        v-for="q in quizzes"
        :key="q.id"
        :to="`/quizzes/${q.id}`"
        class="flex items-center gap-4 p-4 hover:bg-surface-50 dark:hover:bg-surface-800/60 transition-colors"
      >
        <div class="w-9 h-9 rounded-lg bg-violet-50 dark:bg-violet-950/30 flex items-center justify-center shrink-0">
          <FileQuestion class="w-4.5 h-4.5 text-violet-600 dark:text-violet-400" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate">{{ q.title }}</p>
          <div class="flex items-center gap-3 text-xs text-surface-500 dark:text-surface-400 mt-0.5">
            <span><Clock class="w-3 h-3 inline mr-0.5" />{{ q.duration_minutes }} menit</span>
            <span>{{ q.total_questions }} soal</span>
            <span v-if="q.due_date">Tenggat: {{ formatDate(q.due_date, { day: 'numeric', month: 'short' }) }}</span>
          </div>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <UiBadge :variant="statusBadge[q.status]?.variant as any" size="sm">{{ statusBadge[q.status]?.label }}</UiBadge>
          <div v-if="auth.isStudent && q.my_latest_attempt?.score !== null && q.my_latest_attempt?.score !== undefined" class="text-sm font-bold text-emerald-600 dark:text-emerald-400">
            {{ q.my_latest_attempt.score }}/{{ q.total_points }}
          </div>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>

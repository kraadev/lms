<script setup lang="ts">
import { FileQuestion, Clock, CheckCircle2, AlertCircle, ArrowRight } from 'lucide-vue-next'
import { quizzesService } from '~/services/quizzes'
import type { Quiz } from '~/types'
import { formatDate } from '~/utils/formatters'

definePageMeta({ middleware: 'auth' })
useSeoMeta({ title: 'Daftar Kuis' })

const auth = useAuthStore()
const quizzes = ref<Quiz[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

async function loadQuizzes() {
  isLoading.value = true
  error.value = null
  try {
    quizzes.value = await quizzesService.getAll()
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat kuis'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadQuizzes)

const statusBadge: Record<string, { variant: string; label: string }> = {
  upcoming: { variant: 'info', label: 'Akan Datang' },
  active: { variant: 'success', label: 'Sedang Berlangsung' },
  closed: { variant: 'default', label: 'Ditutup' }
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-5xl mx-auto space-y-6">
    <div>
      <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">Semua Kuis</h1>
      <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Daftar evaluasi dan kuis dari kelas Anda.</p>
    </div>

    <div v-if="isLoading">
      <UiSkeleton :rows="5" />
    </div>

    <UiErrorState v-else-if="error" :message="error" @retry="loadQuizzes" />

    <UiEmptyState
      v-else-if="!quizzes.length"
      :icon="FileQuestion"
      title="Belum ada kuis"
      description="Kuis yang dibuat pengajar akan muncul di sini."
    />

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <NuxtLink
        v-for="q in quizzes"
        :key="q.id"
        :to="`/quizzes/${q.id}`"
        class="group p-5 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 hover:border-brand-300 dark:hover:border-brand-700 shadow-soft hover:shadow-elevated transition-all flex flex-col justify-between"
      >
        <div>
          <div class="flex items-start justify-between gap-2 mb-3">
            <div class="w-10 h-10 rounded-xl bg-violet-50 dark:bg-violet-950/40 flex items-center justify-center shrink-0">
              <FileQuestion class="w-5 h-5 text-violet-600 dark:text-violet-400" />
            </div>
            <UiBadge :variant="statusBadge[q.status]?.variant as any" size="sm">
              {{ statusBadge[q.status]?.label }}
            </UiBadge>
          </div>

          <h3 class="text-sm font-semibold text-surface-900 dark:text-surface-100 group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors mb-1 line-clamp-1">
            {{ q.title }}
          </h3>
          <p class="text-xs text-surface-500 dark:text-surface-400 mb-3">
            {{ q.class_title || q.class?.title || 'Kelas' }}
          </p>
        </div>

        <div class="pt-3 border-t border-surface-100 dark:border-surface-800 flex items-center justify-between text-xs text-surface-500">
          <div class="flex items-center gap-1.5">
            <Clock class="w-3.5 h-3.5" />
            <span>{{ q.duration_minutes }} menit</span>
          </div>

          <div v-if="auth.isStudent && q.my_latest_attempt?.score !== null && q.my_latest_attempt?.score !== undefined" class="font-bold text-emerald-600 dark:text-emerald-400">
            Nilai: {{ q.my_latest_attempt.score }}/{{ q.total_points }}
          </div>
          <div v-else class="flex items-center text-brand-600 dark:text-brand-400 font-medium group-hover:translate-x-0.5 transition-transform">
            Buka <ArrowRight class="w-3.5 h-3.5 ml-1" />
          </div>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>

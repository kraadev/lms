<script setup lang="ts">
import { BookOpen, ClipboardList, FileQuestion, Video, AlertCircle, Clock, Sparkles, ArrowRight, Users, GraduationCap } from 'lucide-vue-next'
import { dashboardService } from '~/services/dashboard'
import type { StudentDashboardData, TeacherDashboardData, AdminDashboardData } from '~/types'
import { formatDate, formatRelativeTime } from '~/utils/formatters'

definePageMeta({ middleware: 'auth' })

useSeoMeta({ title: 'Dashboard' })

const auth = useAuthStore()

const studentData = ref<StudentDashboardData | null>(null)
const teacherData = ref<TeacherDashboardData | null>(null)
const adminData = ref<AdminDashboardData | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

async function loadDashboard() {
  isLoading.value = true
  error.value = null
  try {
    if (auth.isStudent) {
      studentData.value = await dashboardService.getStudentDashboard()
    } else if (auth.isTeacher) {
      teacherData.value = await dashboardService.getTeacherDashboard()
    } else if (auth.isAdmin) {
      adminData.value = await dashboardService.getAdminDashboard()
    }
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat dashboard'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadDashboard)

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Selamat pagi'
  if (hour < 17) return 'Selamat siang'
  return 'Selamat malam'
})

function submissionStatusVariant(status: string) {
  const map: Record<string, any> = {
    not_submitted: 'warning',
    submitted: 'info',
    late: 'danger',
    graded: 'success'
  }
  return map[status] || 'default'
}

function submissionStatusLabel(status: string) {
  const map: Record<string, string> = {
    not_submitted: 'Belum Dikumpulkan',
    submitted: 'Dikumpulkan',
    late: 'Terlambat',
    graded: 'Dinilai'
  }
  return map[status] || status
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-6xl mx-auto">
    <!-- Loading -->
    <div v-if="isLoading" class="space-y-6">
      <UiSkeleton class="h-10 w-64 rounded-xl" />
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <UiSkeleton v-for="i in 6" :key="i" class="h-28 rounded-xl" />
      </div>
    </div>

    <!-- Error -->
    <UiErrorState v-else-if="error" :message="error" @retry="loadDashboard" />

    <!-- Student Dashboard -->
    <template v-else-if="auth.isStudent && studentData">
      <!-- Header greeting hero -->
      <div class="mb-6 relative overflow-hidden rounded-2xl p-6 bg-gradient-to-br from-brand-600/10 via-brand-500/5 to-surface-100/50 dark:from-brand-950/40 dark:via-surface-900/60 dark:to-surface-950 border border-brand-200/50 dark:border-brand-900/40 shadow-soft">
        <div class="relative z-10 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <div class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-brand-100/80 dark:bg-brand-950/80 text-brand-700 dark:text-brand-300 mb-2">
              <Sparkles class="w-3.5 h-3.5 text-brand-500" /> Portal Siswa
            </div>
            <h1 class="text-2xl font-bold tracking-tight text-surface-900 dark:text-surface-100">{{ greeting }}, {{ auth.user?.name?.split(' ')[0] }} 👋</h1>
            <p class="text-sm text-surface-600 dark:text-surface-400 mt-1">Ringkasan aktivitas belajar, tugas, dan jadwal kelas hari ini.</p>
          </div>
          <div class="flex items-center gap-3">
            <NuxtLink to="/classes">
              <UiButton variant="outline" size="sm" class="bg-white/80 dark:bg-surface-900/80 backdrop-blur-xs">
                Jelajahi Kelas
              </UiButton>
            </NuxtLink>
            <NuxtLink to="/assignments">
              <UiButton variant="primary" size="sm">
                Lihat Tugas
              </UiButton>
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Active Meeting Banner -->
      <div v-if="studentData.active_meeting" class="mb-6 flex items-center justify-between gap-4 p-4.5 bg-gradient-to-r from-emerald-500/15 via-emerald-500/5 to-transparent border border-emerald-500/30 rounded-2xl shadow-soft">
        <div class="flex items-center gap-3.5">
          <div class="w-10 h-10 rounded-xl bg-emerald-500 text-white flex items-center justify-center shrink-0 shadow-sm shadow-emerald-500/30 animate-pulse-subtle">
            <Video class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span class="inline-block w-2 h-2 rounded-full bg-emerald-500 animate-ping" />
              <p class="text-sm font-bold text-emerald-900 dark:text-emerald-200">Kelas online sedang berlangsung</p>
            </div>
            <p class="text-xs font-medium text-emerald-700 dark:text-emerald-400 mt-0.5">{{ studentData.active_meeting.title }}</p>
          </div>
        </div>
        <NuxtLink :to="`/meetings/${studentData.active_meeting.id}`">
          <UiButton size="sm" variant="success" class="shadow-sm">Gabung Sekarang</UiButton>
        </NuxtLink>
      </div>

      <!-- Continue Learning Hero Card (TSK-035) -->
      <UiCard
        v-if="studentData.current_classes?.length"
        variant="interactive"
        padding="lg"
        class="mb-6 bg-gradient-to-br from-brand-500/10 via-brand-500/5 to-transparent border-brand-200/80 dark:border-brand-900/60 shadow-soft hover:shadow-elevated glass-card-interactive"
      >
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-5">
          <div class="space-y-2 min-w-0">
            <div class="flex items-center gap-2">
              <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-md text-[11px] font-semibold bg-brand-100 text-brand-700 dark:bg-brand-950/70 dark:text-brand-300">
                <Sparkles class="w-3 h-3 text-brand-600 dark:text-brand-400" />
                Lanjutkan Belajar
              </span>
              <span class="text-xs text-surface-500 dark:text-surface-400">Aktivitas Terkini</span>
            </div>
            <h2 class="text-lg font-bold text-surface-900 dark:text-surface-50 truncate">
              {{ studentData.current_classes[0].title }}
            </h2>
            <p class="text-xs text-surface-600 dark:text-surface-400">
              Pengajar: <span class="font-medium text-surface-800 dark:text-surface-200">{{ studentData.current_classes[0].teacher?.name || 'Guru Pengampu' }}</span> &middot; {{ studentData.current_classes[0].academic_year || 'Tahun Ajaran Aktif' }}
            </p>
          </div>

          <div class="flex sm:flex-col sm:items-end justify-between items-center gap-3 shrink-0">
            <div class="w-36 sm:w-48">
              <UiProgress :value="68" :max="100" size="sm" variant="brand" show-label label-position="right" />
            </div>
            <NuxtLink :to="`/classes/${studentData.current_classes[0].id}`">
              <UiButton size="sm" variant="primary" class="shadow-sm">
                Masuk Kelas <ArrowRight class="w-3.5 h-3.5 ml-1" />
              </UiButton>
            </NuxtLink>
          </div>
        </div>
      </UiCard>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Left column -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Classes -->
          <section>
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300 flex items-center gap-1.5"><BookOpen class="w-4 h-4" /> Kelas Aktif</h2>
              <NuxtLink to="/classes" class="text-xs text-brand-600 dark:text-brand-400 hover:underline">Lihat semua</NuxtLink>
            </div>
            <div v-if="!studentData.current_classes?.length">
              <UiEmptyState title="Belum ada kelas" description="Hubungi admin atau guru untuk bergabung ke kelas." />
            </div>
            <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <NuxtLink
                v-for="cls in studentData.current_classes"
                :key="cls.id"
                :to="`/classes/${cls.id}`"
                class="block p-4 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 hover:border-brand-300 dark:hover:border-brand-700 hover:shadow-soft transition-all"
              >
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-lg bg-brand-100 dark:bg-brand-900/40 flex items-center justify-center shrink-0">
                    <BookOpen class="w-4.5 h-4.5 text-brand-600 dark:text-brand-400" />
                  </div>
                  <div class="min-w-0">
                    <p class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate">{{ cls.title }}</p>
                    <p class="text-xs text-surface-500 dark:text-surface-400">{{ cls.teacher?.name }}</p>
                  </div>
                </div>
              </NuxtLink>
            </div>
          </section>

          <!-- Upcoming Assignments -->
          <section>
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300 flex items-center gap-1.5"><ClipboardList class="w-4 h-4" /> Tugas Mendatang</h2>
              <NuxtLink to="/assignments" class="text-xs text-brand-600 dark:text-brand-400 hover:underline">Lihat semua</NuxtLink>
            </div>
            <div v-if="!studentData.upcoming_assignments?.length">
              <p class="text-sm text-surface-400 py-4">Tidak ada tugas mendatang. 🎉</p>
            </div>
            <div v-else class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden">
              <NuxtLink
                v-for="a in studentData.upcoming_assignments"
                :key="a.id"
                :to="`/assignments/${a.id}`"
                class="flex items-center gap-3 p-3.5 hover:bg-surface-50 dark:hover:bg-surface-800/60 transition-colors"
              >
                <div class="w-1.5 h-1.5 rounded-full bg-amber-400 shrink-0" />
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-surface-800 dark:text-surface-200 truncate">{{ a.title }}</p>
                  <p class="text-xs text-surface-500 dark:text-surface-400">{{ a.class_title }}</p>
                </div>
                <div class="text-right shrink-0">
                  <UiBadge v-if="a.my_submission" :variant="submissionStatusVariant(a.my_submission.status)" size="sm">{{ submissionStatusLabel(a.my_submission.status) }}</UiBadge>
                  <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5 flex items-center gap-0.5">
                    <Clock class="w-3 h-3" />
                    {{ formatDate(a.due_date, { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' }) }}
                  </p>
                </div>
              </NuxtLink>
            </div>
          </section>
        </div>

        <!-- Right column -->
        <div class="space-y-6">
          <!-- Upcoming Quizzes -->
          <section>
            <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300 flex items-center gap-1.5 mb-3"><FileQuestion class="w-4 h-4" /> Kuis Mendatang</h2>
            <div v-if="!studentData.upcoming_quizzes?.length">
              <p class="text-sm text-surface-400 py-2">Tidak ada kuis aktif.</p>
            </div>
            <div v-else class="space-y-2">
              <NuxtLink
                v-for="q in studentData.upcoming_quizzes"
                :key="q.id"
                :to="`/quizzes/${q.id}`"
                class="flex items-center gap-3 p-3 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 hover:border-brand-300 dark:hover:border-brand-700 transition-all"
              >
                <div class="w-8 h-8 rounded-lg bg-violet-100 dark:bg-violet-900/40 flex items-center justify-center shrink-0">
                  <FileQuestion class="w-4 h-4 text-violet-600 dark:text-violet-400" />
                </div>
                <div class="min-w-0">
                  <p class="text-sm font-medium text-surface-800 dark:text-surface-200 truncate">{{ q.title }}</p>
                  <p class="text-xs text-surface-500 dark:text-surface-400">{{ q.duration_minutes }} menit</p>
                </div>
              </NuxtLink>
            </div>
          </section>

          <!-- Recent Announcements -->
          <section v-if="studentData.recent_announcements?.length">
            <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">Pengumuman Terbaru</h2>
            <div class="space-y-2">
              <div
                v-for="ann in studentData.recent_announcements"
                :key="ann.id"
                class="p-3 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900"
              >
                <p class="text-sm font-medium text-surface-800 dark:text-surface-200">{{ ann.title }}</p>
                <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">{{ ann.class_title }} &middot; {{ formatRelativeTime(ann.created_at) }}</p>
              </div>
            </div>
          </section>

          <!-- Recent Grades -->
          <section v-if="studentData.recent_grades?.length">
            <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">Nilai Terbaru</h2>
            <div class="space-y-2">
              <div
                v-for="grade in studentData.recent_grades"
                :key="grade.id"
                class="flex items-center justify-between gap-3 p-3 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900"
              >
                <div class="min-w-0">
                  <p class="text-sm font-medium text-surface-800 dark:text-surface-200 truncate">{{ grade.title }}</p>
                  <p class="text-xs text-surface-500 dark:text-surface-400">{{ grade.class_title }}</p>
                </div>
                <div class="text-right shrink-0">
                  <p class="text-sm font-bold text-emerald-600 dark:text-emerald-400">{{ grade.score }}<span class="text-xs font-normal text-surface-400">/{{ grade.max_score }}</span></p>
                </div>
              </div>
            </div>
          </section>
        </div>
      </div>
    </template>

    <!-- Teacher Dashboard -->
    <template v-else-if="auth.isTeacher && teacherData">
      <div class="mb-6 relative overflow-hidden rounded-2xl p-6 bg-gradient-to-br from-brand-600/10 via-brand-500/5 to-surface-100/50 dark:from-brand-950/40 dark:via-surface-900/60 dark:to-surface-950 border border-brand-200/50 dark:border-brand-900/40 shadow-soft">
        <div class="relative z-10 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <div class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-brand-100/80 dark:bg-brand-950/80 text-brand-700 dark:text-brand-300 mb-2">
              <GraduationCap class="w-3.5 h-3.5 text-brand-500" /> Portal Pengajar
            </div>
            <h1 class="text-2xl font-bold tracking-tight text-surface-900 dark:text-surface-100">{{ greeting }}, {{ auth.user?.name?.split(' ')[0] }} 👋</h1>
            <p class="text-sm text-surface-600 dark:text-surface-400 mt-1">Kelola kelas, evaluasi tugas masuk, dan pantau perkembangan siswa.</p>
          </div>
          <div class="flex items-center gap-3">
            <NuxtLink to="/classes">
              <UiButton variant="primary" size="sm" class="shadow-sm">
                Kelola Kelas
              </UiButton>
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Stats row -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3.5 mb-6">
        <div class="p-4.5 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200/80 dark:border-surface-800 shadow-soft hover:shadow-elevated glass-card-interactive transition-all">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold uppercase tracking-wider text-surface-400 dark:text-surface-500">Kelas Diajar</span>
            <div class="w-8 h-8 rounded-xl bg-brand-50 dark:bg-brand-950/60 text-brand-600 dark:text-brand-400 flex items-center justify-center">
              <BookOpen class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-extrabold text-surface-900 dark:text-surface-100">{{ teacherData.classes_taught?.length || 0 }}</p>
          <p class="text-[11px] text-surface-500 dark:text-surface-400 mt-0.5">Total kelas aktif</p>
        </div>

        <div class="p-4.5 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200/80 dark:border-surface-800 shadow-soft hover:shadow-elevated glass-card-interactive transition-all">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold uppercase tracking-wider text-surface-400 dark:text-surface-500">Menunggu Penilaian</span>
            <div class="w-8 h-8 rounded-xl bg-amber-50 dark:bg-amber-950/60 text-amber-600 dark:text-amber-400 flex items-center justify-center">
              <ClipboardList class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-extrabold text-amber-600 dark:text-amber-400">{{ teacherData.pending_grading_count || 0 }}</p>
          <p class="text-[11px] text-surface-500 dark:text-surface-400 mt-0.5">Submission belum diperiksa</p>
        </div>

        <div class="p-4.5 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200/80 dark:border-surface-800 shadow-soft hover:shadow-elevated glass-card-interactive transition-all">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold uppercase tracking-wider text-surface-400 dark:text-surface-500">Kuis Aktif</span>
            <div class="w-8 h-8 rounded-xl bg-violet-50 dark:bg-violet-950/60 text-violet-600 dark:text-violet-400 flex items-center justify-center">
              <FileQuestion class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-extrabold text-violet-600 dark:text-violet-400">{{ teacherData.quiz_overview?.active_quizzes || 0 }}</p>
          <p class="text-[11px] text-surface-500 dark:text-surface-400 mt-0.5">Evaluasi berlangsung</p>
        </div>

        <div class="p-4.5 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200/80 dark:border-surface-800 shadow-soft hover:shadow-elevated glass-card-interactive transition-all">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold uppercase tracking-wider text-surface-400 dark:text-surface-500">Meeting Aktif</span>
            <div class="w-8 h-8 rounded-xl bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
              <Video class="w-4 h-4" />
            </div>
          </div>
          <p class="text-2xl font-extrabold text-emerald-600 dark:text-emerald-400">{{ teacherData.active_meetings?.length || 0 }}</p>
          <p class="text-[11px] text-surface-500 dark:text-surface-400 mt-0.5">Ruang vicon dibuka</p>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Classes -->
        <section>
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300">Kelas yang Diajar</h2>
            <NuxtLink to="/classes" class="text-xs text-brand-600 dark:text-brand-400 hover:underline">Lihat semua</NuxtLink>
          </div>
          <div class="space-y-2">
            <NuxtLink
              v-for="cls in teacherData.classes_taught"
              :key="cls.id"
              :to="`/classes/${cls.id}`"
              class="flex items-center gap-3 p-3.5 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 hover:border-brand-300 dark:hover:border-brand-700 transition-all"
            >
              <div class="w-9 h-9 rounded-lg bg-brand-100 dark:bg-brand-900/40 flex items-center justify-center shrink-0">
                <BookOpen class="w-4.5 h-4.5 text-brand-600 dark:text-brand-400" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate">{{ cls.title }}</p>
                <p class="text-xs text-surface-500 dark:text-surface-400">{{ cls.member_count || 0 }} siswa &middot; {{ cls.academic_year }}</p>
              </div>
              <UiB :variant="cls.status === 'active' ? 'success' : 'default'" size="sm">{{ cls.status === 'active' ? 'Aktif' : 'Arsip' }}</UiB>
            </NuxtLink>
          </div>
        </section>

        <!-- Pending grading -->
        <section>
          <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">Menunggu Penilaian</h2>
          <div v-if="!teacherData.pending_grading?.length">
            <p class="text-sm text-surface-400 py-2">Tidak ada tugas yang perlu dinilai. ✅</p>
          </div>
          <div v-else class="space-y-2">
            <NuxtLink
              v-for="item in teacherData.pending_grading"
              :key="item.assignment.id"
              :to="`/assignments/${item.assignment.id}`"
              class="flex items-center gap-3 p-3.5 rounded-xl border border-amber-200 dark:border-amber-900/40 bg-amber-50 dark:bg-amber-950/20 hover:border-amber-300 dark:hover:border-amber-800 transition-all"
            >
              <ClipboardList class="w-4.5 h-4.5 text-amber-600 dark:text-amber-400 shrink-0" />
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-surface-800 dark:text-surface-200 truncate">{{ item.assignment.title }}</p>
                <p class="text-xs text-surface-500 dark:text-surface-400">{{ item.submission_count }} pengumpulan menunggu</p>
              </div>
            </NuxtLink>
          </div>
        </section>
      </div>
    </template>

    <!-- Admin Dashboard -->
    <template v-else-if="auth.isAdmin && adminData">
      <div class="mb-6 relative overflow-hidden rounded-2xl p-6 bg-gradient-to-br from-brand-600/10 via-brand-500/5 to-surface-100/50 dark:from-brand-950/40 dark:via-surface-900/60 dark:to-surface-950 border border-brand-200/50 dark:border-brand-900/40 shadow-soft">
        <div class="relative z-10 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <div class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-brand-100/80 dark:bg-brand-950/80 text-brand-700 dark:text-brand-300 mb-2">
              <Users class="w-3.5 h-3.5 text-brand-500" /> Pusat Kontrol Administrator
            </div>
            <h1 class="text-2xl font-bold tracking-tight text-surface-900 dark:text-surface-100">Admin Dashboard</h1>
            <p class="text-sm text-surface-600 dark:text-surface-400 mt-1">Pantau statistik ekosistem, alur pengguna, kelas, dan integritas sistem LMS.</p>
          </div>
          <div class="flex items-center gap-3">
            <NuxtLink to="/admin/users">
              <UiButton variant="outline" size="sm" class="bg-white/80 dark:bg-surface-900/80 backdrop-blur-xs">
                Kelola User
              </UiButton>
            </NuxtLink>
            <NuxtLink to="/admin/classes">
              <UiButton variant="primary" size="sm" class="shadow-sm">
                Kelola Kelas
              </UiButton>
            </NuxtLink>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3.5 mb-6">
        <UiCard variant="default" padding="md" class="border border-surface-200/80 dark:border-surface-800 rounded-2xl shadow-soft hover:shadow-elevated glass-card-interactive transition-all">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wider text-surface-400 dark:text-surface-500">Total Siswa</p>
              <p class="text-2xl sm:text-3xl font-extrabold text-surface-900 dark:text-surface-50 mt-1">{{ adminData.total_students }}</p>
            </div>
            <div class="w-10 h-10 rounded-xl bg-brand-50 dark:bg-brand-950/60 flex items-center justify-center text-brand-600 dark:text-brand-400">
              <Users class="w-5 h-5" />
            </div>
          </div>
        </UiCard>

        <UiCard variant="default" padding="md" class="border border-surface-200/80 dark:border-surface-800 rounded-2xl shadow-soft hover:shadow-elevated glass-card-interactive transition-all">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wider text-surface-400 dark:text-surface-500">Total Guru</p>
              <p class="text-2xl sm:text-3xl font-extrabold text-amber-600 dark:text-amber-400 mt-1">{{ adminData.total_teachers }}</p>
            </div>
            <div class="w-10 h-10 rounded-xl bg-amber-50 dark:bg-amber-950/60 flex items-center justify-center text-amber-600 dark:text-amber-400">
              <GraduationCap class="w-5 h-5" />
            </div>
          </div>
        </UiCard>

        <UiCard variant="default" padding="md" class="border border-surface-200/80 dark:border-surface-800 rounded-2xl shadow-soft hover:shadow-elevated glass-card-interactive transition-all">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wider text-surface-400 dark:text-surface-500">Total Kelas</p>
              <p class="text-2xl sm:text-3xl font-extrabold text-violet-600 dark:text-violet-400 mt-1">{{ adminData.total_classes }}</p>
            </div>
            <div class="w-10 h-10 rounded-xl bg-violet-50 dark:bg-violet-950/60 flex items-center justify-center text-violet-600 dark:text-violet-400">
              <BookOpen class="w-5 h-5" />
            </div>
          </div>
        </UiCard>

        <UiCard variant="default" padding="md" class="border border-surface-200/80 dark:border-surface-800 rounded-2xl shadow-soft hover:shadow-elevated glass-card-interactive transition-all">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wider text-surface-400 dark:text-surface-500">Kelas Aktif</p>
              <p class="text-2xl sm:text-3xl font-extrabold text-emerald-600 dark:text-emerald-400 mt-1">{{ adminData.active_classes }}</p>
            </div>
            <div class="w-10 h-10 rounded-xl bg-emerald-50 dark:bg-emerald-950/60 flex items-center justify-center text-emerald-600 dark:text-emerald-400">
              <Sparkles class="w-5 h-5" />
            </div>
          </div>
        </UiCard>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <section>
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300">Pengguna Terbaru</h2>
            <NuxtLink to="/admin/users" class="text-xs font-medium text-brand-600 dark:text-brand-400 hover:underline">Kelola pengguna</NuxtLink>
          </div>
          <div class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft">
            <div v-if="!adminData.recent_users?.length" class="p-6 text-center text-sm text-surface-400">Belum ada pengguna baru</div>
            <div v-for="u in adminData.recent_users" :key="u.id" class="flex items-center gap-3 p-3.5 hover:bg-surface-50/70 dark:hover:bg-surface-800/40 transition-colors">
              <UiAvatar :name="u.name" size="sm" />
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-surface-800 dark:text-surface-200 truncate">{{ u.name }}</p>
                <p class="text-xs text-surface-500 dark:text-surface-400 truncate">{{ u.email }}</p>
              </div>
              <UiB :variant="u.role === 'admin' ? 'danger' : u.role === 'teacher' ? 'primary' : 'default'" size="sm">
                {{ u.role === 'admin' ? 'Admin' : u.role === 'teacher' ? 'Guru' : 'Siswa' }}
              </UiB>
            </div>
          </div>
        </section>

        <section>
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300">Aktivitas Sistem</h2>
          </div>
          <div class="rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft">
            <div v-if="!adminData.system_activity?.length" class="p-6">
              <UiEmptyState variant="default" size="sm" title="Belum Ada Aktivitas" description="Log dan aktivitas sistem akan tercatat di sini secara otomatis." />
            </div>
            <div v-else class="divide-y divide-surface-100 dark:divide-surface-800">
              <div v-for="act in adminData.system_activity" :key="act.id" class="flex items-start gap-3 p-3.5 hover:bg-surface-50/70 dark:hover:bg-surface-800/40 transition-colors">
                <div class="w-7 h-7 rounded-full bg-brand-50 dark:bg-brand-950/60 text-brand-600 dark:text-brand-400 flex items-center justify-center shrink-0 mt-0.5">
                  <div class="w-2 h-2 rounded-full bg-brand-500" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-sm text-surface-800 dark:text-surface-200">{{ act.description }}</p>
                  <p class="text-xs text-surface-400 mt-0.5">{{ formatRelativeTime(act.timestamp) }}</p>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
// Use UiB as alias to prevent component name conflict in template
import UiB from '~/components/ui/UiBadge.vue'
export default { components: { UiB } }
</script>

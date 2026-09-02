<script setup lang="ts">
import { BookOpen, ClipboardList, FileQuestion, Video, AlertCircle, Clock } from 'lucide-vue-next'
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
      <!-- Header greeting -->
      <div class="mb-6">
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">{{ greeting }}, {{ auth.user?.name?.split(' ')[0] }} 👋</h1>
        <p class="text-sm text-surface-500 dark:text-surface-400 mt-0.5">Ini ringkasan aktivitas belajar kamu hari ini.</p>
      </div>

      <!-- Active Meeting Banner -->
      <div v-if="studentData.active_meeting" class="mb-5 flex items-center justify-between gap-4 p-4 bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-200 dark:border-emerald-900/50 rounded-xl">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-full bg-emerald-100 dark:bg-emerald-900/50 flex items-center justify-center shrink-0">
            <Video class="w-4.5 h-4.5 text-emerald-600 dark:text-emerald-400" />
          </div>
          <div>
            <p class="text-sm font-semibold text-emerald-800 dark:text-emerald-200">Kelas online sedang berlangsung</p>
            <p class="text-xs text-emerald-600 dark:text-emerald-400">{{ studentData.active_meeting.title }}</p>
          </div>
        </div>
        <NuxtLink :to="`/meetings/${studentData.active_meeting.id}`">
          <UiButton size="sm" variant="success">Gabung</UiButton>
        </NuxtLink>
      </div>

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
                  <UiB v-if="a.my_submission" :variant="submissionStatusVariant(a.my_submission.status)" size="sm">{{ submissionStatusLabel(a.my_submission.status) }}</UiB>
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
      <div class="mb-6">
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">{{ greeting }}, {{ auth.user?.name?.split(' ')[0] }} 👋</h1>
        <p class="text-sm text-surface-500 dark:text-surface-400 mt-0.5">Kelola kelas dan pantau perkembangan siswa Anda.</p>
      </div>

      <!-- Stats row -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-6">
        <div class="p-4 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          <p class="text-2xl font-bold text-surface-900 dark:text-surface-100">{{ teacherData.classes_taught?.length || 0 }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Kelas Diajar</p>
        </div>
        <div class="p-4 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          <p class="text-2xl font-bold text-amber-600 dark:text-amber-400">{{ teacherData.pending_grading_count || 0 }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Menunggu Penilaian</p>
        </div>
        <div class="p-4 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          <p class="text-2xl font-bold text-brand-600 dark:text-brand-400">{{ teacherData.quiz_overview?.active_quizzes || 0 }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Kuis Aktif</p>
        </div>
        <div class="p-4 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          <p class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{{ teacherData.active_meetings?.length || 0 }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Meeting Aktif</p>
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
      <div class="mb-6">
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">Admin Dashboard</h1>
        <p class="text-sm text-surface-500 dark:text-surface-400 mt-0.5">Pantau statistik dan aktivitas sistem LMS.</p>
      </div>

      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-6">
        <div class="p-5 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          <p class="text-3xl font-bold text-surface-900 dark:text-surface-100">{{ adminData.total_students }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-1">Total Siswa</p>
        </div>
        <div class="p-5 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          <p class="text-3xl font-bold text-brand-600 dark:text-brand-400">{{ adminData.total_teachers }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-1">Total Guru</p>
        </div>
        <div class="p-5 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          <p class="text-3xl font-bold text-violet-600 dark:text-violet-400">{{ adminData.total_classes }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-1">Total Kelas</p>
        </div>
        <div class="p-5 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          <p class="text-3xl font-bold text-emerald-600 dark:text-emerald-400">{{ adminData.active_classes }}</p>
          <p class="text-xs text-surface-500 dark:text-surface-400 mt-1">Kelas Aktif</p>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <section>
          <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">Pengguna Terbaru</h2>
          <div class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden">
            <div v-if="!adminData.recent_users?.length" class="p-6 text-center text-sm text-surface-400">Belum ada pengguna baru</div>
            <div v-for="u in adminData.recent_users" :key="u.id" class="flex items-center gap-3 p-3.5">
              <UiAvatar :name="u.name" size="sm" />
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-surface-800 dark:text-surface-200 truncate">{{ u.name }}</p>
                <p class="text-xs text-surface-500 dark:text-surface-400 truncate">{{ u.email }}</p>
              </div>
              <UiB :variant="u.role === 'admin' ? 'danger' : u.role === 'teacher' ? 'primary' : 'default'" size="sm">{{ u.role }}</UiB>
            </div>
          </div>
        </section>
        <section>
          <h2 class="text-sm font-semibold text-surface-700 dark:text-surface-300 mb-3">Aktivitas Sistem</h2>
          <div class="space-y-2">
            <div v-if="!adminData.system_activity?.length" class="text-sm text-surface-400 py-4">Belum ada aktivitas</div>
            <div v-for="act in adminData.system_activity" :key="act.id" class="flex items-start gap-3 p-3 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900">
              <div class="w-7 h-7 rounded-full bg-surface-100 dark:bg-surface-800 flex items-center justify-center shrink-0 mt-0.5">
                <div class="w-2 h-2 rounded-full bg-brand-500" />
              </div>
              <div>
                <p class="text-sm text-surface-800 dark:text-surface-200">{{ act.description }}</p>
                <p class="text-xs text-surface-400 mt-0.5">{{ formatRelativeTime(act.timestamp) }}</p>
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

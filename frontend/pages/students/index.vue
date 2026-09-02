<script setup lang="ts">
import { GraduationCap, Search, Mail, BookOpen, Users, Phone } from 'lucide-vue-next'
import { adminService } from '~/services/admin'
import { classesService } from '~/services/classes'
import type { User, Class } from '~/types'

definePageMeta({ middleware: 'auth' })
useSeoMeta({ title: 'Daftar Siswa - LMS' })

const auth = useAuthStore()
const toast = useToast()

const students = ref<User[]>([])
const searchQuery = ref('')
const isLoading = ref(true)

async function loadStudents() {
  isLoading.value = true
  try {
    const res = await adminService.getUsers({ role: 'student', search: searchQuery.value })
    students.value = Array.isArray(res) ? res : []
  } catch (err: any) {
    toast.error('Gagal memuat siswa', err?.message)
  } finally {
    isLoading.value = false
  }
}

onMounted(loadStudents)
watch(searchQuery, () => {
  loadStudents()
})
</script>

<template>
  <div class="p-4 md:p-6 max-w-6xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <div>
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100 flex items-center gap-2">
          <GraduationCap class="w-6 h-6 text-brand-600 dark:text-brand-400" />
          Direktori Siswa
        </h1>
        <p class="text-sm text-surface-500 dark:text-surface-400 mt-0.5">
          Daftar seluruh siswa terdaftar dalam sistem akademik
        </p>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-4 mb-6 shadow-soft">
      <div class="relative">
        <Search class="w-4 h-4 text-surface-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari nama siswa atau email..."
          class="w-full pl-9 pr-4 py-2 text-sm rounded-lg border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500"
        />
      </div>
    </div>

    <!-- Students List -->
    <div v-if="isLoading" class="space-y-3">
      <UiSkeleton :rows="4" />
    </div>

    <div v-else-if="!students.length" class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 p-8 shadow-soft">
      <UiEmptyState
        :icon="Users"
        title="Tidak ada siswa ditemukan"
        description="Belum ada data siswa atau tidak sesuai dengan kata kunci pencarian."
      />
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="s in students"
        :key="s.id"
        class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 p-5 shadow-soft hover:border-brand-500/50 transition-all flex items-start gap-4"
      >
        <UiAvatar :name="s.name" :src="s.avatar" size="lg" />
        <div class="min-w-0 flex-1">
          <h3 class="text-sm font-bold text-surface-900 dark:text-surface-100 truncate">{{ s.name }}</h3>
          <p class="text-xs text-surface-500 dark:text-surface-400 flex items-center gap-1.5 mt-1 truncate">
            <Mail class="w-3.5 h-3.5 text-surface-400 shrink-0" />
            <span class="truncate">{{ s.email }}</span>
          </p>
          <div class="mt-3">
            <UiBadge variant="info" size="xs">Siswa Aktif</UiBadge>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

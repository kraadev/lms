<script setup lang="ts">
import { BookOpen, Plus } from 'lucide-vue-next'
import { classesService } from '~/services/classes'
import type { Class } from '~/types'

definePageMeta({ middleware: 'auth' })
useSeoMeta({ title: 'Kelas Saya' })

const auth = useAuthStore()

if (auth.isAdmin) {
  navigateTo('/admin/classes')
}

const classes = ref<Class[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)
const showJoinModal = ref(false)
const showCreateModal = ref(false)
const joinCode = ref('')
const joinLoading = ref(false)
const joinError = ref('')
const toast = useToast()

const createForm = reactive({ title: '', academic_year: '', description: '' })
const createLoading = ref(false)

async function load() {
  isLoading.value = true
  error.value = null
  try {
    classes.value = await classesService.getAll()
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat kelas'
  } finally {
    isLoading.value = false
  }
}

onMounted(load)

async function joinClass() {
  if (!joinCode.value.trim()) { joinError.value = 'Kode kelas wajib diisi'; return }
  joinLoading.value = true
  joinError.value = ''
  try {
    const res = await classesService.joinByCode(joinCode.value.trim())
    classes.value.unshift(res.class)
    showJoinModal.value = false
    joinCode.value = ''
    toast.success('Berhasil bergabung ke kelas!')
  } catch (err: any) {
    joinError.value = err?.message || 'Kode kelas tidak valid'
  } finally {
    joinLoading.value = false
  }
}

async function createClass() {
  if (!createForm.title || !createForm.academic_year) return
  createLoading.value = true
  try {
    const cls = await classesService.create({ ...createForm })
    classes.value.unshift(cls)
    showCreateModal.value = false
    Object.assign(createForm, { title: '', academic_year: '', description: '' })
    toast.success('Kelas berhasil dibuat!')
  } catch (err: any) {
    toast.error(err?.message || 'Gagal membuat kelas')
  } finally {
    createLoading.value = false
  }
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-6xl mx-auto">
    <!-- Header banner -->
    <div class="mb-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-5 bg-gradient-to-r from-brand-600/10 via-brand-500/5 to-surface-100/50 dark:from-brand-950/40 dark:via-surface-900/60 dark:to-surface-950 border border-brand-200/50 dark:border-brand-900/40 rounded-2xl shadow-soft">
      <div>
        <div class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-brand-100/80 dark:bg-brand-950/80 text-brand-700 dark:text-brand-300 mb-1">
          <BookOpen class="w-3.5 h-3.5 text-brand-500" /> Katalog Akademik
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-surface-900 dark:text-surface-100">Kelas Saya</h1>
        <p class="text-xs sm:text-sm text-surface-600 dark:text-surface-400 mt-0.5">Daftar kelas yang sedang kamu ikuti dan aktif semester ini.</p>
      </div>
      <div class="flex gap-2">
        <UiButton v-if="auth.isStudent" size="sm" variant="primary" class="shadow-sm" @click="showJoinModal = true">
          <Plus class="w-4 h-4" />
          Bergabung dengan Kode
        </UiButton>
        <UiButton v-if="auth.isTeacher || auth.isAdmin" size="sm" variant="primary" class="shadow-sm" @click="showCreateModal = true">
          <Plus class="w-4 h-4" />
          Buat Kelas Baru
        </UiButton>
      </div>
    </div>

    <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <UiSkeleton v-for="i in 6" :key="i" class="h-40 rounded-2xl" />
    </div>
    <UiErrorState v-else-if="error" :message="error" @retry="load" />
    <UiEmptyState v-else-if="!classes.length" :icon="BookOpen" title="Belum ada kelas" :description="auth.isStudent ? 'Bergabung ke kelas menggunakan kode yang diberikan guru.' : 'Buat kelas baru untuk memulai.'" >
      <template #action>
        <UiButton v-if="auth.isStudent" size="sm" @click="showJoinModal = true">Bergabung dengan Kode</UiButton>
        <UiButton v-else size="sm" @click="showCreateModal = true">Buat Kelas Baru</UiButton>
      </template>
    </UiEmptyState>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4.5">
      <NuxtLink
        v-for="cls in classes"
        :key="cls.id"
        :to="`/classes/${cls.id}`"
        class="group relative flex flex-col justify-between p-5 bg-white dark:bg-surface-900 rounded-2xl border border-surface-200/80 dark:border-surface-800 hover:border-brand-300 dark:hover:border-brand-700/80 shadow-soft hover:shadow-elevated glass-card-interactive transition-all"
      >
        <div>
          <div class="flex items-start justify-between gap-2 mb-3">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-brand-600 to-brand-400 text-white flex items-center justify-center shrink-0 shadow-sm shadow-brand-500/30 group-hover:scale-105 transition-transform">
              <BookOpen class="w-5 h-5" />
            </div>
            <UiBadge :variant="cls.status === 'active' ? 'success' : 'default'" size="sm">{{ cls.status === 'active' ? 'Aktif' : 'Arsip' }}</UiBadge>
          </div>
          <h3 class="text-base font-bold text-surface-900 dark:text-surface-100 mb-1.5 group-hover:text-brand-600 dark:group-hover:text-brand-300 transition-colors line-clamp-1">{{ cls.title }}</h3>
          <p class="text-xs text-surface-500 dark:text-surface-400 line-clamp-2 leading-relaxed mb-4">{{ cls.description || 'Tidak ada deskripsi tambahan untuk kelas ini.' }}</p>
        </div>

        <div class="pt-3 border-t border-surface-100 dark:border-surface-800/80 flex items-center justify-between text-xs text-surface-500 dark:text-surface-400 font-medium">
          <span class="truncate max-w-[130px]">{{ cls.teacher?.name || 'Guru Pengajar' }}</span>
          <span class="inline-flex items-center gap-1 bg-surface-100 dark:bg-surface-800/60 px-2 py-0.5 rounded-md text-[11px]">
            {{ cls.member_count || 0 }} anggota
          </span>
        </div>
      </NuxtLink>
    </div>

    <!-- Join Class Modal -->
    <UiModal :show="showJoinModal" title="Bergabung ke Kelas" size="sm" @close="showJoinModal = false">
      <div class="space-y-4">
        <UiInput v-model="joinCode" label="Kode Kelas" placeholder="Masukkan kode kelas" :error="joinError" required />
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="showJoinModal = false">Batal</UiButton>
          <UiButton :loading="joinLoading" @click="joinClass">Bergabung</UiButton>
        </div>
      </template>
    </UiModal>

    <!-- Create Class Modal -->
    <UiModal :show="showCreateModal" title="Buat Kelas Baru" size="sm" @close="showCreateModal = false">
      <form @submit.prevent="createClass" class="space-y-4">
        <UiInput v-model="createForm.title" label="Nama Kelas" placeholder="Matematika Kelas X" required />
        <UiInput v-model="createForm.academic_year" label="Tahun Ajaran" placeholder="2025/2026" required />
        <UiTextarea v-model="createForm.description" label="Deskripsi" placeholder="Opsional" :rows="3" />
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="showCreateModal = false">Batal</UiButton>
          <UiButton :loading="createLoading" @click="createClass">Buat Kelas</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

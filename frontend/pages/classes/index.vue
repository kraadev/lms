<script setup lang="ts">
import { BookOpen, Plus } from 'lucide-vue-next'
import { classesService } from '~/services/classes'
import type { Class } from '~/types'

definePageMeta({ middleware: 'auth' })
useSeoMeta({ title: 'Kelas Saya' })

const auth = useAuthStore()

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
  <div class="p-4 md:p-6 max-w-5xl mx-auto">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">Kelas Saya</h1>
      <div class="flex gap-2">
        <UiButton v-if="auth.isStudent" size="sm" variant="outline" @click="showJoinModal = true">
          <Plus class="w-4 h-4" />
          Bergabung
        </UiButton>
        <UiButton v-if="auth.isTeacher || auth.isAdmin" size="sm" @click="showCreateModal = true">
          <Plus class="w-4 h-4" />
          Buat Kelas
        </UiButton>
      </div>
    </div>

    <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <UiSkeleton v-for="i in 6" :key="i" class="h-32 rounded-xl" />
    </div>
    <UiErrorState v-else-if="error" :message="error" @retry="load" />
    <UiEmptyState v-else-if="!classes.length" :icon="BookOpen" title="Belum ada kelas" :description="auth.isStudent ? 'Bergabung ke kelas menggunakan kode yang diberikan guru.' : 'Buat kelas baru untuk memulai.'" >
      <template #action>
        <UiButton v-if="auth.isStudent" size="sm" @click="showJoinModal = true">Bergabung dengan Kode</UiButton>
        <UiButton v-else size="sm" @click="showCreateModal = true">Buat Kelas Baru</UiButton>
      </template>
    </UiEmptyState>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <NuxtLink
        v-for="cls in classes"
        :key="cls.id"
        :to="`/classes/${cls.id}`"
        class="group block p-5 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 hover:border-brand-300 dark:hover:border-brand-700 hover:shadow-soft transition-all"
      >
        <div class="flex items-start justify-between gap-2 mb-3">
          <div class="w-10 h-10 rounded-xl bg-brand-50 dark:bg-brand-950/50 flex items-center justify-center shrink-0">
            <BookOpen class="w-5 h-5 text-brand-600 dark:text-brand-400" />
          </div>
          <UiBadge :variant="cls.status === 'active' ? 'success' : 'default'" size="sm">{{ cls.status === 'active' ? 'Aktif' : 'Arsip' }}</UiBadge>
        </div>
        <h3 class="text-sm font-semibold text-surface-900 dark:text-surface-100 mb-1 group-hover:text-brand-700 dark:group-hover:text-brand-300 transition-colors">{{ cls.title }}</h3>
        <p class="text-xs text-surface-500 dark:text-surface-400 line-clamp-2 mb-2">{{ cls.description || cls.academic_year }}</p>
        <div class="flex items-center justify-between text-xs text-surface-400 dark:text-surface-500">
          <span>{{ cls.teacher?.name }}</span>
          <span>{{ cls.member_count || 0 }} anggota</span>
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

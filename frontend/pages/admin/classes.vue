<script setup lang="ts">
import { BookMarked, Plus, Search, Trash2, Edit2, Users, BookOpen } from 'lucide-vue-next'
import { adminService } from '~/services/admin'
import { classesService } from '~/services/classes'
import type { Class, User } from '~/types'

definePageMeta({ middleware: 'admin' })
useSeoMeta({ title: 'Manajemen Kelas - Admin' })

const toast = useToast()

const classes = ref<Class[]>([])
const teachers = ref<User[]>([])
const searchQuery = ref('')
const isLoading = ref(true)
const isSaving = ref(false)
const showModal = ref(false)
const editingClass = ref<Class | null>(null)

const form = reactive({
  title: '',
  code: '',
  teacher_id: '' as string | number,
  academic_year: '2024/2025',
  description: ''
})

async function loadData() {
  isLoading.value = true
  try {
    const [clsList, userList] = await Promise.all([
      adminService.getClasses({ search: searchQuery.value }),
      adminService.getUsers({ role: 'teacher' })
    ])
    classes.value = Array.isArray(clsList) ? clsList : []
    teachers.value = Array.isArray(userList) ? userList : []
  } catch (err: any) {
    toast.error('Gagal memuat kelas', err?.message)
  } finally {
    isLoading.value = false
  }
}

onMounted(loadData)
watch(searchQuery, () => loadData())

function openCreateModal() {
  editingClass.value = null
  form.title = ''
  form.code = ''
  form.teacher_id = teachers.value[0]?.id || ''
  form.academic_year = '2024/2025'
  form.description = ''
  showModal.value = true
}

function openEditModal(cls: Class) {
  editingClass.value = cls
  form.title = cls.title || cls.name || ''
  form.code = cls.code || ''
  form.teacher_id = cls.teacher_id || ''
  form.academic_year = cls.academic_year || '2024/2025'
  form.description = cls.description || ''
  showModal.value = true
}

async function handleSave() {
  if (!form.title.trim() || !form.teacher_id) {
    toast.error('Validasi Gagal', 'Judul dan guru pengajar wajib diisi.')
    return
  }

  isSaving.value = true
  try {
    if (editingClass.value) {
      await adminService.updateClass(editingClass.value.id, {
        title: form.title.trim(),
        description: form.description.trim(),
        academic_year: form.academic_year
      })
      toast.success('Berhasil', 'Data kelas berhasil diperbarui.')
    } else {
      await adminService.createClass({
        title: form.title.trim(),
        code: form.code.trim() || `CLS-${Math.floor(1000 + Math.random() * 9000)}`,
        teacher_id: Number(form.teacher_id),
        academic_year: form.academic_year,
        description: form.description.trim()
      })
      toast.success('Berhasil', 'Kelas baru berhasil dibuat.')
    }
    showModal.value = false
    await loadData()
  } catch (err: any) {
    toast.error('Gagal menyimpan kelas', err?.message)
  } finally {
    isSaving.value = false
  }
}

async function handleDelete(cls: Class) {
  if (!confirm(`Hapus kelas "${cls.title || cls.name}"?`)) return

  try {
    await adminService.deleteClass(cls.id)
    toast.success('Berhasil', 'Kelas berhasil dihapus.')
    await loadData()
  } catch (err: any) {
    toast.error('Gagal menghapus kelas', err?.message)
  }
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-6xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <div>
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100 flex items-center gap-2">
          <BookMarked class="w-6 h-6 text-brand-600 dark:text-brand-400" />
          Manajemen Kelas & Kurikulum
        </h1>
        <p class="text-sm text-surface-500 dark:text-surface-400 mt-0.5">
          Kelola seluruh kelas, penugasan guru, dan tahun akademik
        </p>
      </div>

      <UiButton @click="openCreateModal" size="sm" class="gap-1.5 self-start sm:self-auto">
        <Plus class="w-4 h-4" />
        Tambah Kelas Baru
      </UiButton>
    </div>

    <!-- Search Bar -->
    <div class="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-4 mb-6 shadow-soft">
      <div class="relative">
        <Search class="w-4 h-4 text-surface-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari kelas..."
          class="w-full pl-9 pr-4 py-2 text-sm rounded-lg border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500"
        />
      </div>
    </div>

    <!-- Classes List -->
    <div v-if="isLoading" class="space-y-3">
      <UiSkeleton :rows="4" />
    </div>

    <div v-else-if="!classes.length" class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 p-8 shadow-soft">
      <UiEmptyState
        :icon="BookMarked"
        title="Belum ada kelas"
        description="Belum ada data kelas yang terdaftar."
      >
        <template #action>
          <UiButton @click="openCreateModal" size="sm" class="gap-1.5 mt-3">
            <Plus class="w-4 h-4" />
            Buat Kelas Pertama
          </UiButton>
        </template>
      </UiEmptyState>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="c in classes"
        :key="c.id"
        class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 p-5 shadow-soft flex flex-col justify-between hover:border-brand-500/50 transition-all"
      >
        <div>
          <div class="flex items-start justify-between gap-3 mb-3">
            <div class="w-10 h-10 rounded-xl bg-brand-100 dark:bg-brand-950/50 flex items-center justify-center shrink-0">
              <BookOpen class="w-5 h-5 text-brand-600 dark:text-brand-400" />
            </div>
            <UiBadge :variant="c.status === 'active' ? 'success' : 'default'" size="sm">
              {{ c.status === 'active' ? 'Aktif' : 'Arsip' }}
            </UiBadge>
          </div>

          <h3 class="text-base font-bold text-surface-900 dark:text-surface-100 mb-1">
            {{ c.title || c.name }}
          </h3>
          <p class="text-xs text-surface-500 dark:text-surface-400 line-clamp-2 mb-3">
            {{ c.description || 'Tidak ada deskripsi' }}
          </p>
          <p class="text-xs text-surface-600 dark:text-surface-300 font-medium">
            Tahun Ajaran: {{ c.academic_year }}
          </p>
        </div>

        <div class="pt-4 border-t border-surface-100 dark:border-surface-800 flex items-center justify-between mt-4">
          <NuxtLink :to="`/classes/${c.id}`" class="text-xs font-semibold text-brand-600 dark:text-brand-400 hover:underline">
            Lihat Kelas &rarr;
          </NuxtLink>

          <div class="flex items-center gap-1">
            <button
              @click="openEditModal(c)"
              class="p-1.5 rounded-lg text-surface-500 hover:text-brand-600 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
              title="Edit Kelas"
            >
              <Edit2 class="w-4 h-4" />
            </button>
            <button
              @click="handleDelete(c)"
              class="p-1.5 rounded-lg text-surface-500 hover:text-red-600 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
              title="Hapus Kelas"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <UiModal v-model="showModal" :title="editingClass ? 'Edit Data Kelas' : 'Buat Kelas Baru'" size="md">
      <form @submit.prevent="handleSave" class="space-y-4">
        <UiInput
          v-model="form.title"
          label="Nama / Judul Kelas"
          placeholder="Contoh: Pemrograman Golang & Vue"
          required
        />

        <div>
          <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-1">Guru Pengajar</label>
          <select
            v-model="form.teacher_id"
            class="w-full px-3 py-2 text-sm rounded-lg border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500"
            required
          >
            <option value="" disabled>Pilih Guru Pengajar</option>
            <option v-for="t in teachers" :key="t.id" :value="t.id">
              {{ t.name }} ({{ t.email }})
            </option>
          </select>
        </div>

        <UiInput
          v-model="form.academic_year"
          label="Tahun Akademik"
          placeholder="2024/2025"
          required
        />

        <div>
          <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-1">Deskripsi Kelas</label>
          <textarea
            v-model="form.description"
            rows="3"
            class="w-full px-3 py-2 text-sm rounded-lg border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500 resize-none"
            placeholder="Keterangan materi pembelajaran..."
          />
        </div>

        <div class="pt-3 flex justify-end gap-2">
          <UiButton type="button" variant="secondary" @click="showModal = false">Batal</UiButton>
          <UiButton type="submit" :loading="isSaving">Simpan Data</UiButton>
        </div>
      </form>
    </UiModal>
  </div>
</template>

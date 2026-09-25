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

const showMembersModal = ref(false)
const selectedClassForMembers = ref<Class | null>(null)
const classMembers = ref<any[]>([])
const loadingMembers = ref(false)

const form = reactive({
  title: '',
  code: '',
  teacher_id: '' as string | number,
  academic_year: '2024/2025',
  description: ''
})

async function openMembersModal(cls: Class) {
  selectedClassForMembers.value = cls
  showMembersModal.value = true
  loadingMembers.value = true
  try {
    const res = await classesService.getMembers(cls.id)
    classMembers.value = Array.isArray(res) ? res : []
  } catch (err: any) {
    classMembers.value = []
    toast.error('Gagal memuat anggota', err?.message)
  } finally {
    loadingMembers.value = false
  }
}

async function toggleStatus(cls: Class) {
  const newStatus = cls.status === 'active' ? 'archived' : 'active'
  const label = newStatus === 'active' ? 'diaktifkan' : 'diarsipkan'
  try {
    await adminService.updateClass(cls.id, { status: newStatus })
    cls.status = newStatus
    toast.success('Status Berubah', `Kelas berhasil ${label}.`)
  } catch (err: any) {
    toast.error('Gagal mengubah status', err?.message)
  }
}

function copyCode(code?: string) {
  if (!code) return
  if (navigator.clipboard) navigator.clipboard.writeText(code)
  toast.success('Tersalin', `Kode kelas ${code} disalin ke clipboard.`)
}

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
          <!-- Header card with status toggle -->
          <div class="flex items-start justify-between gap-3 mb-3">
            <div class="w-10 h-10 rounded-xl bg-brand-50 dark:bg-brand-950/50 flex items-center justify-center shrink-0">
              <BookOpen class="w-5 h-5 text-brand-600 dark:text-brand-400" />
            </div>
            <button
              type="button"
              :class="[
                'px-2.5 py-0.5 rounded-full text-xs font-semibold transition-colors cursor-pointer',
                c.status === 'active'
                  ? 'bg-emerald-50 text-emerald-700 hover:bg-emerald-100 dark:bg-emerald-950/40 dark:text-emerald-300'
                  : 'bg-surface-100 text-surface-600 hover:bg-surface-200 dark:bg-surface-800 dark:text-surface-400'
              ]"
              :title="c.status === 'active' ? 'Klik untuk mengarsipkan kelas' : 'Klik untuk mengaktifkan kembali'"
              @click="toggleStatus(c)"
            >
              {{ c.status === 'active' ? 'Aktif' : 'Arsip' }}
            </button>
          </div>

          <h3 class="text-base font-bold text-surface-900 dark:text-surface-100 mb-1 truncate">
            {{ c.title || c.name }}
          </h3>
          <p class="text-xs text-surface-500 dark:text-surface-400 line-clamp-2 mb-3">
            {{ c.description || 'Tidak ada deskripsi' }}
          </p>

          <!-- Class CRUD metadata -->
          <div class="space-y-2 text-xs text-surface-600 dark:text-surface-300 pt-3 border-t border-surface-100 dark:border-surface-800/80">
            <div class="flex items-center justify-between">
              <span class="text-surface-400">Guru Pengajar:</span>
              <span class="font-semibold text-surface-800 dark:text-surface-200 truncate max-w-[150px]">
                {{ c.teacher?.name || 'Belum Ditugaskan' }}
              </span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-surface-400">Tahun Ajaran:</span>
              <span class="font-medium">{{ c.academic_year || '2024/2025' }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-surface-400">Kode Masuk:</span>
              <button
                type="button"
                class="inline-flex items-center gap-1 font-mono font-bold text-[11px] text-brand-600 dark:text-brand-400 hover:underline cursor-pointer"
                title="Klik untuk menyalin kode"
                @click="copyCode(c.code)"
              >
                {{ c.code || '-' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Action CRUD Buttons (Admin only) -->
        <div class="pt-4 border-t border-surface-100 dark:border-surface-800 flex items-center justify-between mt-4">
          <UiButton variant="outline" size="sm" class="gap-1.5 text-xs" @click="openMembersModal(c)">
            <Users class="w-3.5 h-3.5" />
            Anggota ({{ c.member_count || 0 }})
          </UiButton>

          <div class="flex items-center gap-1">
            <button
              @click="openEditModal(c)"
              class="p-2 rounded-lg text-surface-500 hover:text-brand-600 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
              title="Edit Data Kelas"
            >
              <Edit2 class="w-4 h-4" />
            </button>
            <button
              @click="handleDelete(c)"
              class="p-2 rounded-lg text-surface-500 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors"
              title="Hapus Kelas"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <UiModal :show="showModal" @close="showModal = false" :title="editingClass ? 'Edit Data Kelas' : 'Buat Kelas Baru'" size="md">
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

    <!-- Members Inspection Modal (Admin CRUD) -->
    <UiModal :show="showMembersModal" @close="showMembersModal = false" :title="`Daftar Anggota: ${selectedClassForMembers?.title || ''}`" size="md">
      <div v-if="loadingMembers" class="py-6">
        <UiSkeleton :rows="4" />
      </div>
      <div v-else-if="!classMembers.length" class="py-8 text-center text-xs text-surface-400">
        Belum ada siswa atau anggota yang terdaftar di kelas ini.
      </div>
      <div v-else class="max-h-80 overflow-y-auto divide-y divide-surface-100 dark:divide-surface-800">
        <div v-for="m in classMembers" :key="m.id" class="flex items-center justify-between py-2.5">
          <div class="flex items-center gap-2.5">
            <UiAvatar :name="m.user?.name || m.name" size="sm" />
            <div>
              <p class="text-xs font-semibold text-surface-900 dark:text-surface-100">{{ m.user?.name || m.name }}</p>
              <p class="text-[10px] text-surface-400">{{ m.user?.email || m.email }}</p>
            </div>
          </div>
          <UiBadge :variant="m.role === 'teacher' ? 'primary' : 'default'" size="sm">
            {{ m.role === 'teacher' ? 'Guru' : 'Siswa' }}
          </UiBadge>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end">
          <UiButton variant="outline" size="sm" @click="showMembersModal = false">Tutup</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

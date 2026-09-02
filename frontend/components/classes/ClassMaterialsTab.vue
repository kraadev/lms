<script setup lang="ts">
import { BookMarked, Plus, ExternalLink, Download, Pencil, Trash2 } from 'lucide-vue-next'
import { materialsService } from '~/services/materials'
import type { Material } from '~/types'
import { formatDate, formatFileSize, formatRelativeTime } from '~/utils/formatters'

const props = defineProps<{ classId: string }>()
const auth = useAuthStore()
const toast = useToast()

const materials = ref<Material[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)
const showModal = ref(false)
const editMaterial = ref<Material | null>(null)
const deleteConfirm = ref<Material | null>(null)

const form = reactive({ title: '', description: '', external_link: '' })
const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const isSaving = ref(false)
const isDeleting = ref(false)

async function load() {
  isLoading.value = true
  error.value = null
  try {
    materials.value = await materialsService.getByClass(props.classId)
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat materi'
  } finally {
    isLoading.value = false
  }
}

onMounted(load)

function openCreate() {
  editMaterial.value = null
  Object.assign(form, { title: '', description: '', external_link: '' })
  selectedFile.value = null
  showModal.value = true
}

function openEdit(mat: Material) {
  editMaterial.value = mat
  Object.assign(form, { title: mat.title, description: mat.description || '', external_link: mat.external_link || '' })
  selectedFile.value = null
  showModal.value = true
}

async function save() {
  if (!form.title.trim()) return
  isSaving.value = true
  try {
    const data = new FormData()
    data.append('title', form.title)
    data.append('description', form.description)
    data.append('external_link', form.external_link)
    if (selectedFile.value) data.append('file', selectedFile.value)

    if (editMaterial.value) {
      const updated = await materialsService.update(editMaterial.value.id, data)
      const idx = materials.value.findIndex(m => m.id === editMaterial.value!.id)
      if (idx !== -1) materials.value[idx] = updated
      toast.success('Materi berhasil diperbarui')
    } else {
      const created = await materialsService.create(props.classId, data)
      materials.value.unshift(created)
      toast.success('Materi berhasil ditambahkan')
    }
    showModal.value = false
  } catch (err: any) {
    toast.error(err?.message || 'Gagal menyimpan materi')
  } finally {
    isSaving.value = false
  }
}

async function deleteMaterial(mat: Material) {
  isDeleting.value = true
  try {
    await materialsService.delete(mat.id)
    materials.value = materials.value.filter(m => m.id !== mat.id)
    deleteConfirm.value = null
    toast.success('Materi dihapus')
  } catch (err: any) {
    toast.error(err?.message || 'Gagal menghapus materi')
  } finally {
    isDeleting.value = false
  }
}

function onFileChange(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (files && files[0]) selectedFile.value = files[0]
}
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="flex items-center justify-between mb-5">
      <h2 class="text-base font-semibold text-surface-800 dark:text-surface-200">Materi Pembelajaran</h2>
      <UiButton v-if="auth.isTeacher || auth.isAdmin" size="sm" @click="openCreate">
        <Plus class="w-4 h-4" />
        Tambah Materi
      </UiButton>
    </div>

    <div v-if="isLoading"><UiSkeleton :rows="4" /></div>
    <UiErrorState v-else-if="error" :message="error" @retry="load" />
    <UiEmptyState v-else-if="!materials.length" :icon="BookMarked" title="Belum ada materi" :description="auth.isTeacher ? 'Tambahkan materi untuk siswa Anda.' : 'Guru belum menambahkan materi.'">
      <template #action>
        <UiButton v-if="auth.isTeacher" size="sm" @click="openCreate">Tambah Materi Pertama</UiButton>
      </template>
    </UiEmptyState>

    <div v-else class="space-y-2.5">
      <div
        v-for="mat in materials"
        :key="mat.id"
        class="flex items-start gap-4 p-4 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 hover:border-surface-300 dark:hover:border-surface-700 transition-all"
      >
        <div class="w-9 h-9 rounded-lg bg-brand-50 dark:bg-brand-950/40 flex items-center justify-center shrink-0">
          <BookMarked class="w-4.5 h-4.5 text-brand-500 dark:text-brand-400" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-surface-900 dark:text-surface-100">{{ mat.title }}</p>
          <p v-if="mat.description" class="text-xs text-surface-500 dark:text-surface-400 mt-0.5 line-clamp-2">{{ mat.description }}</p>
          <div class="flex flex-wrap items-center gap-2 mt-2">
            <a v-if="mat.file_url" :href="mat.file_url" target="_blank" rel="noopener" class="inline-flex items-center gap-1 text-xs text-brand-600 dark:text-brand-400 hover:underline">
              <Download class="w-3.5 h-3.5" />
              {{ mat.file_name || 'Unduh' }}
              <span v-if="mat.file_size" class="text-surface-400">({{ formatFileSize(mat.file_size) }})</span>
            </a>
            <a v-if="mat.external_link" :href="mat.external_link" target="_blank" rel="noopener" class="inline-flex items-center gap-1 text-xs text-sky-600 dark:text-sky-400 hover:underline">
              <ExternalLink class="w-3.5 h-3.5" />
              Link Eksternal
            </a>
            <span class="text-xs text-surface-400 dark:text-surface-500">{{ formatRelativeTime(mat.created_at) }}</span>
          </div>
        </div>
        <div v-if="auth.isTeacher || auth.isAdmin" class="flex items-center gap-1 shrink-0">
          <button type="button" class="p-1.5 rounded-lg text-surface-400 hover:text-brand-600 dark:hover:text-brand-400 hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors" aria-label="Edit materi" @click="openEdit(mat)">
            <Pencil class="w-4 h-4" />
          </button>
          <button type="button" class="p-1.5 rounded-lg text-surface-400 hover:text-rose-600 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors" aria-label="Hapus materi" @click="deleteConfirm = mat">
            <Trash2 class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <UiModal :show="showModal" :title="editMaterial ? 'Edit Materi' : 'Tambah Materi'" @close="showModal = false">
      <form @submit.prevent="save" class="space-y-4">
        <UiInput v-model="form.title" label="Judul" placeholder="Nama materi" required />
        <UiTextarea v-model="form.description" label="Deskripsi" placeholder="Opsional" :rows="3" />
        <UiInput v-model="form.external_link" label="Link Eksternal" type="url" placeholder="https://..." />
        <div>
          <label class="block mb-1.5 text-sm font-medium text-surface-700 dark:text-surface-300">Upload File</label>
          <input ref="fileInput" type="file" class="block w-full text-sm text-surface-600 dark:text-surface-400 file:mr-3 file:py-1.5 file:px-3 file:rounded-lg file:border-0 file:text-sm file:font-medium file:bg-brand-50 file:text-brand-700 dark:file:bg-brand-950/50 dark:file:text-brand-300 hover:file:bg-brand-100" @change="onFileChange" />
          <p v-if="selectedFile" class="mt-1 text-xs text-surface-500">Dipilih: {{ selectedFile.name }}</p>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="showModal = false">Batal</UiButton>
          <UiButton :loading="isSaving" @click="save">{{ editMaterial ? 'Simpan' : 'Tambah' }}</UiButton>
        </div>
      </template>
    </UiModal>

    <!-- Delete confirm -->
    <UiModal :show="!!deleteConfirm" title="Hapus Materi?" size="sm" @close="deleteConfirm = null">
      <p class="text-sm text-surface-600 dark:text-surface-400">Apakah Anda yakin ingin menghapus materi <strong>{{ deleteConfirm?.title }}</strong>? Tindakan ini tidak dapat dibatalkan.</p>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="deleteConfirm = null">Batal</UiButton>
          <UiButton variant="danger" :loading="isDeleting" @click="deleteMaterial(deleteConfirm!)">Hapus</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

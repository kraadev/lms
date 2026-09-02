<script setup lang="ts">
import { BookOpen, Megaphone, Users, Key, Copy, Check, Plus, Calendar, UserCheck } from 'lucide-vue-next'
import type { Class } from '~/types'
import { formatRelativeTime } from '~/utils/formatters'
import { announcementsService } from '~/services/announcements'

const props = defineProps<{ classData: Class }>()
const auth = useAuthStore()
const toast = useToast()

const copied = ref(false)
const announcements = ref<any[]>([])
const isLoadingAnnouncements = ref(true)
const showCreateAnnouncement = ref(false)
const announcementForm = reactive({ title: '', content: '' })
const isSavingAnnouncement = ref(false)

async function copyClassCode() {
  if (!props.classData.code) return
  try {
    await navigator.clipboard.writeText(props.classData.code)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
    toast.success('Kode kelas disalin ke clipboard')
  } catch {
    toast.error('Gagal menyalin kode')
  }
}

async function loadAnnouncements() {
  isLoadingAnnouncements.value = true
  try {
    announcements.value = await announcementsService.getByClass(props.classData.id)
  } catch (err: any) {
    // silently catch or show empty
  } finally {
    isLoadingAnnouncements.value = false
  }
}

async function createAnnouncement() {
  if (!announcementForm.title.trim() || !announcementForm.content.trim()) return
  isSavingAnnouncement.value = true
  try {
    const ann = await announcementsService.create(props.classData.id, announcementForm)
    announcements.value.unshift(ann)
    showCreateAnnouncement.value = false
    Object.assign(announcementForm, { title: '', content: '' })
    toast.success('Pengumuman berhasil diposting!')
  } catch (err: any) {
    toast.error(err?.message || 'Gagal memposting pengumuman')
  } finally {
    isSavingAnnouncement.value = false
  }
}

onMounted(() => {
  loadAnnouncements()
})
</script>

<template>
  <div class="p-4 md:p-6 space-y-6">
    <!-- Class Info Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Description & Year -->
      <div class="md:col-span-2 p-5 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
        <h3 class="text-sm font-semibold text-surface-800 dark:text-surface-200 mb-2">Tentang Kelas</h3>
        <p class="text-sm text-surface-600 dark:text-surface-400 whitespace-pre-line leading-relaxed">
          {{ props.classData.description || 'Tidak ada deskripsi untuk kelas ini.' }}
        </p>
        <div class="flex items-center gap-4 mt-4 pt-4 border-t border-surface-100 dark:border-surface-800 text-xs text-surface-500 dark:text-surface-400">
          <div class="flex items-center gap-1.5">
            <Calendar class="w-4 h-4 text-surface-400" />
            <span>Tahun Ajaran: <strong>{{ props.classData.academic_year }}</strong></span>
          </div>
          <div class="flex items-center gap-1.5">
            <UserCheck class="w-4 h-4 text-surface-400" />
            <span>Pengajar: <strong>{{ props.classData.teacher?.name || '-' }}</strong></span>
          </div>
        </div>
      </div>

      <!-- Class Code & Stats (Teacher/Admin or Student) -->
      <div class="p-5 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 flex flex-col justify-between">
        <div>
          <h3 class="text-sm font-semibold text-surface-800 dark:text-surface-200 mb-1">Kode Kelas</h3>
          <p class="text-xs text-surface-500 dark:text-surface-400 mb-3">Bagikan kode ini kepada siswa untuk bergabung ke kelas.</p>
          <div class="flex items-center justify-between p-3 bg-surface-50 dark:bg-surface-800/80 rounded-lg border border-surface-200 dark:border-surface-700">
            <span class="text-base font-mono font-bold tracking-wider text-brand-600 dark:text-brand-400">
              {{ props.classData.code || '------' }}
            </span>
            <button
              v-if="props.classData.code"
              type="button"
              class="p-1.5 rounded-md hover:bg-surface-200 dark:hover:bg-surface-700 text-surface-600 dark:text-surface-300 transition-colors"
              title="Salin Kode"
              @click="copyClassCode"
            >
              <Check v-if="copied" class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
              <Copy v-else class="w-4 h-4" />
            </button>
          </div>
        </div>
        <div class="mt-4 pt-4 border-t border-surface-100 dark:border-surface-800 flex items-center justify-between text-xs text-surface-500">
          <span>Total Anggota</span>
          <span class="font-semibold text-surface-800 dark:text-surface-200">{{ props.classData.member_count || 0 }} Siswa</span>
        </div>
      </div>
    </div>

    <!-- Announcements Section -->
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <h3 class="text-base font-semibold text-surface-800 dark:text-surface-200 flex items-center gap-2">
          <Megaphone class="w-4.5 h-4.5 text-brand-600 dark:text-brand-400" />
          Pengumuman Kelas
        </h3>
        <UiButton v-if="auth.isTeacher || auth.isAdmin" size="sm" @click="showCreateAnnouncement = true">
          <Plus class="w-4 h-4" />
          Buat Pengumuman
        </UiButton>
      </div>

      <div v-if="isLoadingAnnouncements">
        <UiSkeleton :rows="3" />
      </div>
      <div v-else-if="!announcements.length">
        <UiEmptyState
          :icon="Megaphone"
          title="Belum ada pengumuman"
          description="Pengumuman penting terkait kelas akan muncul di sini."
        />
      </div>
      <div v-else class="space-y-3">
        <div
          v-for="ann in announcements"
          :key="ann.id"
          class="p-4 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 shadow-soft"
        >
          <div class="flex items-start justify-between gap-2 mb-2">
            <h4 class="text-sm font-semibold text-surface-900 dark:text-surface-100">{{ ann.title }}</h4>
            <span class="text-xs text-surface-400 shrink-0">{{ formatRelativeTime(ann.created_at) }}</span>
          </div>
          <p class="text-sm text-surface-600 dark:text-surface-400 whitespace-pre-line leading-relaxed">
            {{ ann.content }}
          </p>
          <div v-if="ann.author" class="mt-3 pt-3 border-t border-surface-100 dark:border-surface-800 text-xs text-surface-400 flex items-center gap-2">
            <UiAvatar :name="ann.author.name" size="xs" />
            <span>Diposting oleh <strong>{{ ann.author.name }}</strong></span>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Announcement Modal -->
    <UiModal :show="showCreateAnnouncement" title="Buat Pengumuman Baru" @close="showCreateAnnouncement = false">
      <form @submit.prevent="createAnnouncement" class="space-y-4">
        <UiInput v-model="announcementForm.title" label="Judul Pengumuman" placeholder="Contoh: Jadwal Ujian Tengah Semester" required />
        <UiTextarea v-model="announcementForm.content" label="Isi Pengumuman" placeholder="Tuliskan detail pengumuman..." :rows="4" required />
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="showCreateAnnouncement = false">Batal</UiButton>
          <UiButton :loading="isSavingAnnouncement" @click="createAnnouncement">Posting</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

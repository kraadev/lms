<script setup lang="ts">
import { Video, Plus, Calendar, Clock, User, ArrowRight, BookOpen } from 'lucide-vue-next'
import { classesService } from '~/services/classes'
import { meetingsService } from '~/services/meetings'
import type { Class, Meeting } from '~/types'
import { formatDate, formatRelativeTime } from '~/utils/formatters'

definePageMeta({ middleware: 'auth' })
useSeoMeta({ title: 'Ruang Kelas Online - LiveKit Meeting' })

const auth = useAuthStore()
const toast = useToast()

const classes = ref<Class[]>([])
const selectedClassId = ref<number | string>('')
const meetings = ref<Meeting[]>([])
const isLoading = ref(true)
const isCreating = ref(false)
const showCreateModal = ref(false)

const createForm = reactive({
  title: '',
  type: 'video' as 'video' | 'audio',
  class_id: ''
})

async function loadData() {
  isLoading.value = true
  try {
    const clsList = await classesService.getAll({ status: 'active' })
    classes.value = Array.isArray(clsList) ? clsList : []
    
    if (classes.value.length > 0 && !selectedClassId.value) {
      selectedClassId.value = classes.value[0].id
    }

    if (selectedClassId.value) {
      await loadMeetingsForClass(selectedClassId.value)
    }
  } catch (err: any) {
    toast.error('Gagal memuat kelas', err?.message)
  } finally {
    isLoading.value = false
  }
}

async function loadMeetingsForClass(classId: number | string) {
  try {
    const res = await meetingsService.getByClass(classId)
    meetings.value = Array.isArray(res) ? res : []
  } catch (err: any) {
    console.warn('Load meetings warning:', err)
    meetings.value = []
  }
}

watch(selectedClassId, (newId) => {
  if (newId) loadMeetingsForClass(newId)
})

onMounted(loadData)

async function handleCreateMeeting() {
  const targetClassId = createForm.class_id || selectedClassId.value
  if (!createForm.title.trim() || !targetClassId) return

  isCreating.value = true
  try {
    const created = await meetingsService.create(targetClassId, {
      title: createForm.title.trim(),
      type: createForm.type
    })
    toast.success('Meeting Dibuat', 'Ruang kelas online berhasil dibuka.')
    showCreateModal.value = false
    createForm.title = ''
    await loadMeetingsForClass(targetClassId)
    navigateTo(`/meetings/${created.id}`)
  } catch (err: any) {
    toast.error('Gagal membuat meeting', err?.message)
  } finally {
    isCreating.value = false
  }
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-6xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <div>
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100 flex items-center gap-2">
          <Video class="w-6 h-6 text-brand-600 dark:text-brand-400" />
          Ruang Kelas Online (WebRTC)
        </h1>
        <p class="text-sm text-surface-500 dark:text-surface-400 mt-0.5">
          Video conference interaktif bertenaga LiveKit WebRTC untuk tatap muka kelas
        </p>
      </div>

      <UiButton
        v-if="auth.isTeacher || auth.isAdmin"
        @click="showCreateModal = true"
        size="sm"
        class="gap-1.5 self-start sm:self-auto"
      >
        <Plus class="w-4 h-4" />
        Mulai Meeting Baru
      </UiButton>
    </div>

    <!-- Class Selector -->
    <div class="bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800 p-4 mb-6 shadow-soft">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <label class="text-sm font-medium text-surface-700 dark:text-surface-300">Pilih Kelas:</label>
        <select
          v-model="selectedClassId"
          class="px-3 py-2 text-sm rounded-lg border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500"
        >
          <option v-for="cls in classes" :key="cls.id" :value="cls.id">
            {{ cls.title || cls.name }} ({{ cls.academic_year }})
          </option>
        </select>
      </div>
    </div>

    <!-- Meetings List -->
    <div v-if="isLoading" class="space-y-3">
      <UiSkeleton :rows="3" />
    </div>

    <div v-else-if="!meetings.length" class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 p-8 shadow-soft">
      <UiEmptyState
        :icon="Video"
        title="Belum ada sesi meeting aktif"
        description="Belum ada ruang kelas video conference yang dibuat untuk kelas ini."
      >
        <template v-if="auth.isTeacher || auth.isAdmin" #action>
          <UiButton @click="showCreateModal = true" size="sm" class="gap-1.5 mt-3">
            <Plus class="w-4 h-4" />
            Buka Ruangan Sekarang
          </UiButton>
        </template>
      </UiEmptyState>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="m in meetings"
        :key="m.id"
        class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 p-5 shadow-soft flex flex-col justify-between hover:border-brand-500/50 transition-all"
      >
        <div>
          <div class="flex items-start justify-between gap-3 mb-3">
            <div class="w-10 h-10 rounded-xl bg-brand-100 dark:bg-brand-950/50 flex items-center justify-center shrink-0">
              <Video class="w-5 h-5 text-brand-600 dark:text-brand-400" />
            </div>
            <UiBadge :variant="m.status === 'active' ? 'success' : 'default'" size="sm">
              {{ m.status === 'active' ? 'Sedang Berlangsung' : (m.status === 'ended' ? 'Selesai' : 'Terjadwal') }}
            </UiBadge>
          </div>

          <h3 class="text-base font-bold text-surface-900 dark:text-surface-100 mb-1 line-clamp-1">
            {{ m.title }}
          </h3>
          <p class="text-xs text-surface-500 dark:text-surface-400 mb-4">
            Room: <code class="px-1.5 py-0.5 rounded bg-surface-100 dark:bg-surface-800 text-[11px] font-mono">{{ m.room_name }}</code>
          </p>
        </div>

        <div class="pt-3 border-t border-surface-100 dark:border-surface-800 flex items-center justify-between">
          <span class="text-xs text-surface-400 dark:text-surface-500">
            Dibuat {{ formatRelativeTime(m.created_at) }}
          </span>

          <NuxtLink v-if="m.status !== 'ended'" :to="`/meetings/${m.id}`">
            <UiButton size="sm" variant="primary" class="gap-1">
              Masuk Ruangan
              <ArrowRight class="w-3.5 h-3.5" />
            </UiButton>
          </NuxtLink>
          <span v-else class="text-xs font-medium text-surface-400 dark:text-surface-600">
            Sesi telah berakhir
          </span>
        </div>
      </div>
    </div>

    <!-- Create Modal -->
    <UiModal v-model="showCreateModal" title="Mulai Meeting Kelas Online" size="md">
      <form @submit.prevent="handleCreateMeeting" class="space-y-4">
        <div>
          <label class="block text-xs font-semibold text-surface-700 dark:text-surface-300 mb-1">Pilih Kelas</label>
          <select
            v-model="createForm.class_id"
            class="w-full px-3 py-2 text-sm rounded-lg border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500"
            required
          >
            <option value="" disabled>Pilih Kelas</option>
            <option v-for="cls in classes" :key="cls.id" :value="cls.id">
              {{ cls.title || cls.name }}
            </option>
          </select>
        </div>

        <UiInput
          v-model="createForm.title"
          label="Judul Sesi Meeting"
          placeholder="Contoh: Diskusi Materi & Tanya Jawab"
          required
        />

        <div class="pt-3 flex justify-end gap-2">
          <UiButton type="button" variant="secondary" @click="showCreateModal = false">Batal</UiButton>
          <UiButton type="submit" :loading="isCreating">Buka Ruangan</UiButton>
        </div>
      </form>
    </UiModal>
  </div>
</template>

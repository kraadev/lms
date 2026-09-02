<script setup lang="ts">
import { Video, Mic, Plus, Radio, Clock, CheckCircle2, AlertCircle } from 'lucide-vue-next'
import { meetingsService } from '~/services/meetings'
import type { Meeting, MeetingType, Class } from '~/types'
import { formatDate, formatRelativeTime } from '~/utils/formatters'

const props = defineProps<{ classId: string; classData?: Class }>()
const emit = defineEmits<{ (e: 'refresh'): void }>()
const auth = useAuthStore()
const toast = useToast()

const meetings = ref<Meeting[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

const showCreateModal = ref(false)
const createForm = reactive({
  title: '',
  type: 'video' as MeetingType
})
const isCreating = ref(false)

async function loadMeetings() {
  isLoading.value = true
  error.value = null
  try {
    meetings.value = await meetingsService.getByClass(props.classId)
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat daftar meeting'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadMeetings)

async function createMeeting() {
  if (!createForm.title.trim()) return
  isCreating.value = true
  try {
    const m = await meetingsService.create(props.classId, {
      title: createForm.title.trim(),
      type: createForm.type
    })
    showCreateModal.value = false
    Object.assign(createForm, { title: '', type: 'video' })
    toast.success('Kelas online berhasil dibuat!')
    emit('refresh')
    await navigateTo(`/meetings/${m.id}`)
  } catch (err: any) {
    toast.error(err?.message || 'Gagal membuat sesi meeting')
  } finally {
    isCreating.value = false
  }
}

const activeMeetings = computed(() => meetings.value.filter(m => m.status === 'active'))
const pastMeetings = computed(() => meetings.value.filter(m => m.status !== 'active'))
</script>

<template>
  <div class="p-4 md:p-6 space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-base font-semibold text-surface-800 dark:text-surface-200">Kelas Online & Video Conference</h2>
        <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Sesi tatap muka virtual interaktif dengan audio & video.</p>
      </div>
      <UiButton v-if="auth.isTeacher || auth.isAdmin" size="sm" @click="showCreateModal = true">
        <Plus class="w-4 h-4" />
        Mulai Sesi Baru
      </UiButton>
    </div>

    <div v-if="isLoading">
      <UiSkeleton :rows="4" />
    </div>

    <UiErrorState v-else-if="error" :message="error" @retry="loadMeetings" />

    <div v-else class="space-y-6">
      <!-- Active Meeting Card (if any) -->
      <div v-if="activeMeetings.length" class="space-y-3">
        <h3 class="text-xs font-bold uppercase tracking-wider text-emerald-600 dark:text-emerald-400 flex items-center gap-1.5">
          <Radio class="w-4 h-4 animate-pulse text-emerald-500" />
          Sedang Berlangsung (Live)
        </h3>

        <div
          v-for="m in activeMeetings"
          :key="m.id"
          class="p-5 bg-gradient-to-br from-emerald-50 to-teal-50 dark:from-emerald-950/40 dark:to-teal-950/30 border-2 border-emerald-400 dark:border-emerald-700/60 rounded-2xl shadow-elevated flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
        >
          <div class="flex items-center gap-3.5">
            <div class="w-12 h-12 rounded-xl bg-emerald-500 text-white flex items-center justify-center shrink-0 shadow-md shadow-emerald-500/20">
              <Video v-if="m.type === 'video'" class="w-6 h-6" />
              <Mic v-else class="w-6 h-6" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <h4 class="text-base font-bold text-surface-900 dark:text-white">{{ m.title }}</h4>
                <span class="px-2 py-0.5 rounded-full bg-emerald-500 text-white text-[10px] font-bold uppercase tracking-wider animate-pulse">
                  Live
                </span>
              </div>
              <p class="text-xs text-surface-600 dark:text-surface-300 mt-1">
                Dimulai {{ formatRelativeTime(m.created_at) }}
                <span v-if="m.host">&middot; Dipandu oleh {{ m.host.name }}</span>
              </p>
            </div>
          </div>

          <NuxtLink :to="`/meetings/${m.id}`" class="w-full sm:w-auto shrink-0">
            <UiButton variant="success" size="md" class="w-full shadow-md shadow-emerald-600/20">
              <Video class="w-4 h-4 mr-1.5" />
              Gabung Sekarang
            </UiButton>
          </NuxtLink>
        </div>
      </div>

      <!-- Past / Recorded Meetings -->
      <div class="space-y-3">
        <h3 class="text-xs font-bold uppercase tracking-wider text-surface-500 dark:text-surface-400">
          Riwayat Pertemuan ({{ pastMeetings.length }})
        </h3>

        <div v-if="!pastMeetings.length && !activeMeetings.length">
          <UiEmptyState
            :icon="Video"
            title="Belum ada sesi meeting"
            :description="auth.isTeacher ? 'Mulai sesi kelas online pertama untuk berinteraksi langsung.' : 'Belum ada sesi online yang dibuat pengajar.'"
          >
            <template #action>
              <UiButton v-if="auth.isTeacher" size="sm" @click="showCreateModal = true">Mulai Kelas Online</UiButton>
            </template>
          </UiEmptyState>
        </div>

        <div v-else-if="!pastMeetings.length" class="text-xs text-surface-400 p-4 bg-white dark:bg-surface-900 rounded-xl border border-surface-200 dark:border-surface-800">
          Belum ada riwayat pertemuan sebelumnya.
        </div>

        <div v-else class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden">
          <div
            v-for="m in pastMeetings"
            :key="m.id"
            class="flex items-center justify-between p-4 hover:bg-surface-50 dark:hover:bg-surface-800/40 transition-colors"
          >
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-lg bg-surface-100 dark:bg-surface-800 flex items-center justify-center shrink-0 text-surface-500">
                <Video v-if="m.type === 'video'" class="w-4.5 h-4.5" />
                <Mic v-else class="w-4.5 h-4.5" />
              </div>
              <div>
                <p class="text-sm font-semibold text-surface-800 dark:text-surface-200">{{ m.title }}</p>
                <div class="flex items-center gap-2 text-xs text-surface-400 mt-0.5">
                  <span>{{ formatDate(m.created_at, { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) }}</span>
                  <span v-if="m.host">&middot; {{ m.host.name }}</span>
                </div>
              </div>
            </div>

            <UiBadge variant="default" size="sm">Selesai</UiBadge>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Meeting Modal -->
    <UiModal :show="showCreateModal" title="Mulai Kelas Online Baru" size="sm" @close="showCreateModal = false">
      <form @submit.prevent="createMeeting" class="space-y-4">
        <UiInput
          v-model="createForm.title"
          label="Topik Pertemuan"
          placeholder="Contoh: Sesi Tanya Jawab Bab 4"
          required
        />
        
        <div>
          <label class="block mb-2 text-sm font-medium text-surface-700 dark:text-surface-300">Tipe Pertemuan</label>
          <div class="grid grid-cols-2 gap-3">
            <button
              type="button"
              class="p-3 rounded-xl border text-left flex flex-col items-start gap-1 transition-all"
              :class="createForm.type === 'video'
                ? 'border-brand-500 bg-brand-50/50 dark:bg-brand-950/30 text-brand-900 dark:text-brand-100 ring-2 ring-brand-500/20'
                : 'border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800/40 text-surface-600 dark:text-surface-400'"
              @click="createForm.type = 'video'"
            >
              <Video class="w-5 h-5 mb-1" :class="createForm.type === 'video' ? 'text-brand-600' : 'text-surface-400'" />
              <span class="text-xs font-bold">Video & Audio</span>
              <span class="text-[10px] text-surface-400">Kamera, mic & share screen</span>
            </button>

            <button
              type="button"
              class="p-3 rounded-xl border text-left flex flex-col items-start gap-1 transition-all"
              :class="createForm.type === 'audio'
                ? 'border-brand-500 bg-brand-50/50 dark:bg-brand-950/30 text-brand-900 dark:text-brand-100 ring-2 ring-brand-500/20'
                : 'border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800/40 text-surface-600 dark:text-surface-400'"
              @click="createForm.type = 'audio'"
            >
              <Mic class="w-5 h-5 mb-1" :class="createForm.type === 'audio' ? 'text-brand-600' : 'text-surface-400'" />
              <span class="text-xs font-bold">Audio Saja</span>
              <span class="text-[10px] text-surface-400">Hemat bandwidth</span>
            </button>
          </div>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="showCreateModal = false">Batal</UiButton>
          <UiButton :loading="isCreating" @click="createMeeting">Mulai Sekarang</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>

<script setup lang="ts">
import {
  Mic, MicOff, Video, VideoOff, ScreenShare, PhoneOff,
  Users, MessageSquare, Shield, Settings, Volume2, Maximize2
} from 'lucide-vue-next'
import { meetingsService } from '~/services/meetings'
import type { Meeting, LiveKitJoinResponse } from '~/types'

definePageMeta({
  layout: 'meeting',
  middleware: 'auth'
})

const route = useRoute()
const meetingId = computed(() => route.params.id as string)
const auth = useAuthStore()
const toast = useToast()

const meeting = ref<Meeting | null>(null)
const joinData = ref<LiveKitJoinResponse | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

// Media controls state
const isMicOn = ref(true)
const isCameraOn = ref(true)
const isScreenSharing = ref(false)
const activeSidebar = ref<'chat' | 'participants' | null>(null)
const isEnding = ref(false)

// In-meeting participants (local user + remote peers)
const participants = ref<any[]>([])

useSeoMeta({ title: computed(() => meeting.value?.title ? `${meeting.value.title} - Kelas Online` : 'Ruang Kelas Online') })

async function initMeeting() {
  isLoading.value = true
  error.value = null
  try {
    // 1. Fetch meeting info
    meeting.value = await meetingsService.getById(meetingId.value)
    
    // 2. Fetch LiveKit token & server url
    joinData.value = await meetingsService.join(meetingId.value)

    // Setup initial participants list
    participants.value = [
      {
        id: auth.user?.id,
        name: auth.user?.name || 'Anda',
        isLocal: true,
        isHost: meeting.value.host_id === auth.user?.id || meeting.value.host?.id === auth.user?.id,
        isAudioEnabled: isMicOn.value,
        isVideoEnabled: isCameraOn.value
      }
    ]
  } catch (err: any) {
    error.value = err?.message || 'Gagal terhubung ke ruang meeting'
  } finally {
    isLoading.value = false
  }
}

onMounted(initMeeting)

function toggleMic() {
  isMicOn.value = !isMicOn.value
  const local = participants.value.find(p => p.isLocal)
  if (local) local.isAudioEnabled = isMicOn.value
  toast.info(isMicOn.value ? 'Mikrofon aktif' : 'Mikrofon dibisukan')
}

function toggleCamera() {
  isCameraOn.value = !isCameraOn.value
  const local = participants.value.find(p => p.isLocal)
  if (local) local.isVideoEnabled = isCameraOn.value
  toast.info(isCameraOn.value ? 'Kamera aktif' : 'Kamera dinonaktifkan')
}

function toggleScreenShare() {
  isScreenSharing.value = !isScreenSharing.value
  toast.info(isScreenSharing.value ? 'Berbagi layar dimulai' : 'Berbagi layar dihentikan')
}

function leaveMeeting() {
  const returnUrl = meeting.value?.class_id ? `/classes/${meeting.value.class_id}` : '/dashboard'
  navigateTo(returnUrl)
}

async function endMeetingForAll() {
  isEnding.value = true
  try {
    await meetingsService.end(meetingId.value)
    toast.success('Sesi pertemuan diakhiri')
    const returnUrl = meeting.value?.class_id ? `/classes/${meeting.value.class_id}` : '/dashboard'
    await navigateTo(returnUrl)
  } catch (err: any) {
    toast.error(err?.message || 'Gagal mengakhiri meeting')
  } finally {
    isEnding.value = false
  }
}

const isHost = computed(() => {
  return meeting.value?.host_id === auth.user?.id || meeting.value?.host?.id === auth.user?.id || auth.isAdmin
})
</script>

<template>
  <div class="flex-1 flex flex-col h-screen overflow-hidden bg-surface-950 text-white select-none">
    <!-- Top Bar -->
    <header class="h-14 px-4 bg-surface-900/90 backdrop-blur-md border-b border-surface-800 flex items-center justify-between z-10">
      <div class="flex items-center gap-3">
        <div class="w-3 h-3 rounded-full bg-emerald-500 animate-pulse" />
        <h1 class="text-sm font-bold text-surface-100 truncate max-w-xs sm:max-w-md">
          {{ meeting?.title || 'Memuat Kelas Online...' }}
        </h1>
        <span v-if="meeting?.type" class="px-2 py-0.5 rounded text-[10px] uppercase font-bold bg-surface-800 text-surface-400">
          {{ meeting.type }}
        </span>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          class="p-2 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 transition-colors"
          :class="{ 'bg-surface-800 text-white': activeSidebar === 'participants' }"
          title="Daftar Peserta"
          @click="activeSidebar = activeSidebar === 'participants' ? null : 'participants'"
        >
          <Users class="w-5 h-5" />
        </button>

        <button
          type="button"
          class="p-2 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 transition-colors"
          :class="{ 'bg-surface-800 text-white': activeSidebar === 'chat' }"
          title="Obrolan Meeting"
          @click="activeSidebar = activeSidebar === 'chat' ? null : 'chat'"
        >
          <MessageSquare class="w-5 h-5" />
        </button>
      </div>
    </header>

    <!-- Main Content Stage -->
    <div class="flex-1 flex overflow-hidden relative">
      <!-- Loading State -->
      <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center gap-3">
        <div class="w-10 h-10 border-4 border-brand-500 border-t-transparent rounded-full animate-spin" />
        <p class="text-sm text-surface-400">Menyambungkan ke ruang LiveKit...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex-1 flex flex-col items-center justify-center p-6 text-center gap-4">
        <div class="w-14 h-14 rounded-2xl bg-rose-950/50 border border-rose-900/60 text-rose-500 flex items-center justify-center">
          <PhoneOff class="w-7 h-7" />
        </div>
        <p class="text-base font-semibold text-rose-400">{{ error }}</p>
        <UiButton variant="outline" size="sm" @click="leaveMeeting">Kembali ke Kelas</UiButton>
      </div>

      <!-- Video Grid -->
      <div v-else class="flex-1 p-4 overflow-y-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 items-center justify-center">
        <!-- Participant Tile -->
        <div
          v-for="p in participants"
          :key="p.id"
          class="relative aspect-video bg-surface-900 rounded-2xl overflow-hidden border border-surface-800 flex items-center justify-center shadow-elevated group"
        >
          <!-- Video active display (mock stream / camera) -->
          <div v-if="p.isVideoEnabled" class="w-full h-full bg-gradient-to-tr from-surface-800 to-surface-700 flex items-center justify-center">
            <UiAvatar :name="p.name" size="xl" />
          </div>

          <!-- Video disabled avatar fallback -->
          <div v-else class="flex flex-col items-center justify-center gap-3">
            <UiAvatar :name="p.name" size="xl" class="ring-4 ring-surface-800" />
            <p class="text-xs text-surface-400">Kamera dinonaktifkan</p>
          </div>

          <!-- Participant Label & Status Bottom overlay -->
          <div class="absolute bottom-3 left-3 right-3 flex items-center justify-between px-3 py-1.5 bg-black/60 backdrop-blur-md rounded-xl text-xs">
            <div class="flex items-center gap-2 truncate">
              <span class="font-medium truncate">{{ p.name }} {{ p.isLocal ? '(Anda)' : '' }}</span>
              <span v-if="p.isHost" class="px-1.5 py-0.5 rounded bg-brand-600 text-[10px] font-bold">Host</span>
            </div>

            <div class="flex items-center gap-1.5 shrink-0">
              <Mic v-if="p.isAudioEnabled" class="w-3.5 h-3.5 text-emerald-400" />
              <MicOff v-else class="w-3.5 h-3.5 text-rose-400" />
            </div>
          </div>
        </div>
      </div>

      <!-- Sidebars (Participants or Chat) -->
      <aside
        v-if="activeSidebar"
        class="w-80 bg-surface-900 border-l border-surface-800 flex flex-col z-10 transition-all"
      >
        <!-- Sidebar Header -->
        <div class="h-14 px-4 border-b border-surface-800 flex items-center justify-between">
          <h3 class="text-sm font-bold">
            {{ activeSidebar === 'participants' ? 'Daftar Peserta' : 'Obrolan Sesi' }}
          </h3>
          <button
            type="button"
            class="text-surface-400 hover:text-white text-xs"
            @click="activeSidebar = null"
          >
            Tutup
          </button>
        </div>

        <!-- Participants Content -->
        <div v-if="activeSidebar === 'participants'" class="flex-1 overflow-y-auto p-4 space-y-2">
          <div
            v-for="p in participants"
            :key="p.id"
            class="flex items-center justify-between p-2.5 rounded-xl bg-surface-800/60"
          >
            <div class="flex items-center gap-2.5">
              <UiAvatar :name="p.name" size="xs" />
              <div class="text-xs truncate">
                <p class="font-medium truncate">{{ p.name }}</p>
                <p v-if="p.isLocal" class="text-[10px] text-surface-400">Lokal</p>
              </div>
            </div>

            <div class="flex items-center gap-2 text-surface-400">
              <Mic v-if="p.isAudioEnabled" class="w-3.5 h-3.5 text-emerald-400" />
              <MicOff v-else class="w-3.5 h-3.5 text-rose-400" />
            </div>
          </div>
        </div>

        <!-- Chat Content -->
        <div v-else-if="activeSidebar === 'chat'" class="flex-1 flex flex-col p-4">
          <div class="flex-1 flex items-center justify-center text-center text-xs text-surface-400">
            Pesan dalam sesi live akan muncul di sini.
          </div>
        </div>
      </aside>
    </div>

    <!-- Bottom Control Bar -->
    <footer class="h-20 px-6 bg-surface-900 border-t border-surface-800 flex items-center justify-between z-10">
      <div class="hidden sm:flex items-center gap-2 text-xs text-surface-400">
        <Shield class="w-4 h-4 text-emerald-400" />
        <span>Tersambung Terenkripsi</span>
      </div>

      <!-- Main Action Controls -->
      <div class="flex items-center gap-3 mx-auto sm:mx-0">
        <!-- Mic Toggle -->
        <button
          type="button"
          class="w-12 h-12 rounded-2xl flex items-center justify-center transition-all shadow-md"
          :class="isMicOn
            ? 'bg-surface-800 hover:bg-surface-700 text-white'
            : 'bg-rose-600 hover:bg-rose-700 text-white'"
          title="Nyalakan/Matikan Mic"
          @click="toggleMic"
        >
          <Mic v-if="isMicOn" class="w-5 h-5" />
          <MicOff v-else class="w-5 h-5" />
        </button>

        <!-- Camera Toggle -->
        <button
          type="button"
          class="w-12 h-12 rounded-2xl flex items-center justify-center transition-all shadow-md"
          :class="isCameraOn
            ? 'bg-surface-800 hover:bg-surface-700 text-white'
            : 'bg-rose-600 hover:bg-rose-700 text-white'"
          title="Nyalakan/Matikan Kamera"
          @click="toggleCamera"
        >
          <Video v-if="isCameraOn" class="w-5 h-5" />
          <VideoOff v-else class="w-5 h-5" />
        </button>

        <!-- Screen Share -->
        <button
          type="button"
          class="w-12 h-12 rounded-2xl flex items-center justify-center transition-all shadow-md"
          :class="isScreenSharing
            ? 'bg-brand-600 hover:bg-brand-700 text-white'
            : 'bg-surface-800 hover:bg-surface-700 text-white'"
          title="Berbagi Layar"
          @click="toggleScreenShare"
        >
          <ScreenShare class="w-5 h-5" />
        </button>

        <!-- Leave Button -->
        <button
          type="button"
          class="px-5 h-12 rounded-2xl bg-rose-600 hover:bg-rose-700 text-white font-semibold text-xs flex items-center gap-2 shadow-md shadow-rose-600/20 transition-all"
          title="Tinggalkan Meeting"
          @click="leaveMeeting"
        >
          <PhoneOff class="w-4 h-4" />
          <span>Keluar</span>
        </button>
      </div>

      <!-- Host End Meeting Button -->
      <div class="hidden sm:flex items-center">
        <UiButton
          v-if="isHost"
          variant="danger"
          size="sm"
          :loading="isEnding"
          @click="endMeetingForAll"
        >
          Akhiri untuk Semua
        </UiButton>
      </div>
    </footer>
  </div>
</template>

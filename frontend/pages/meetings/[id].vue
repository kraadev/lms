<script setup lang="ts">
import {
  Mic, MicOff, Video, VideoOff, ScreenShare, PhoneOff,
  Users, MessageSquare, Shield, Settings, Volume2, Maximize2, Monitor, Activity, Copy
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

// Realtime Audio Volume Meter (0-100)
const micVolume = ref(0)
let audioContext: AudioContext | null = null
let analyser: AnalyserNode | null = null
let micSource: MediaStreamAudioSourceNode | null = null
let animFrame: number | null = null

// DOM Video Elements & Media Streams
const localVideoEl = ref<HTMLVideoElement | null>(null)
const screenVideoEl = ref<HTMLVideoElement | null>(null)
const localStream = ref<MediaStream | null>(null)
const screenStream = ref<MediaStream | null>(null)

// In-meeting participants (local user + remote peers)
const participants = ref<any[]>([])

useSeoMeta({ title: computed(() => meeting.value?.title ? `${meeting.value.title} - Kelas Online` : 'Ruang Kelas Online') })

function attachLocalVideo(el: HTMLVideoElement | null) {
  if (!el) return
  localVideoEl.value = el
  if (localStream.value && isCameraOn.value) {
    if (el.srcObject !== localStream.value) {
      el.srcObject = localStream.value
    }
    el.play().catch(err => console.warn('Play video catch:', err))
  }
}

function attachScreenVideo(el: HTMLVideoElement | null) {
  if (!el) return
  screenVideoEl.value = el
  if (screenStream.value && isScreenSharing.value) {
    if (el.srcObject !== screenStream.value) {
      el.srcObject = screenStream.value
    }
    el.play().catch(err => console.warn('Play screen catch:', err))
  }
}

function setupAudioAnalyser(stream: MediaStream) {
  if (typeof window === 'undefined') return
  const audioTracks = stream.getAudioTracks()
  if (!audioTracks.length) return

  try {
    const AudioCtx = window.AudioContext || (window as any).webkitAudioContext
    if (!AudioCtx) return

    if (!audioContext || audioContext.state === 'closed') {
      audioContext = new AudioCtx()
    }
    if (audioContext.state === 'suspended') {
      audioContext.resume()
    }

    analyser = audioContext.createAnalyser()
    analyser.fftSize = 64
    analyser.smoothingTimeConstant = 0.3

    micSource = audioContext.createMediaStreamSource(stream)
    micSource.connect(analyser)

    const dataArray = new Uint8Array(analyser.frequencyBinCount)

    function loop() {
      if (!isMicOn.value || !analyser) {
        micVolume.value = 0
      } else {
        analyser.getByteFrequencyData(dataArray)
        let sum = 0
        for (let i = 0; i < dataArray.length; i++) {
          sum += dataArray[i]
        }
        const avg = sum / dataArray.length
        micVolume.value = Math.min(100, Math.round((avg / 100) * 100))
      }
      animFrame = requestAnimationFrame(loop)
    }
    loop()
  } catch (err) {
    console.warn('Audio analyser setup error:', err)
  }
}

async function startMediaStreams() {
  if (typeof navigator === 'undefined' || !navigator.mediaDevices?.getUserMedia) {
    console.warn('getUserMedia not supported in this environment')
    return
  }

  try {
    const stream = await navigator.mediaDevices.getUserMedia({
      video: {
        width: { ideal: 1280 },
        height: { ideal: 720 },
        facingMode: 'user'
      },
      audio: true
    })
    localStream.value = stream
    isCameraOn.value = true
    isMicOn.value = true

    setupAudioAnalyser(stream)

    await nextTick()
    if (localVideoEl.value) {
      localVideoEl.value.srcObject = stream
      localVideoEl.value.play().catch(e => console.warn('Local video play catch:', e))
    }
  } catch (err: any) {
    console.warn('Full getUserMedia failed, trying audio-only fallback:', err)
    try {
      const audioStream = await navigator.mediaDevices.getUserMedia({ audio: true })
      localStream.value = audioStream
      isCameraOn.value = false
      isMicOn.value = true
      setupAudioAnalyser(audioStream)
    } catch (aErr) {
      console.warn('Audio fallback also failed:', aErr)
      isCameraOn.value = false
      isMicOn.value = false
    }
  }
}

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
        isHost: meeting.value.host_id === auth.user?.id || meeting.value.host?.id === auth.user?.id || auth.isAdmin,
        isAudioEnabled: isMicOn.value,
        isVideoEnabled: isCameraOn.value
      }
    ]

    // Start local camera/mic stream
    await startMediaStreams()
  } catch (err: any) {
    error.value = err?.message || 'Gagal terhubung ke ruang meeting'
  } finally {
    isLoading.value = false
  }
}

onMounted(initMeeting)

onBeforeUnmount(() => {
  if (animFrame) cancelAnimationFrame(animFrame)
  if (audioContext && audioContext.state !== 'closed') {
    audioContext.close()
  }
  if (localStream.value) {
    localStream.value.getTracks().forEach(t => t.stop())
  }
  if (screenStream.value) {
    screenStream.value.getTracks().forEach(t => t.stop())
  }
})

async function toggleMic() {
  if (localStream.value && localStream.value.getAudioTracks().length > 0) {
    const next = !isMicOn.value
    localStream.value.getAudioTracks().forEach(t => { t.enabled = next })
    isMicOn.value = next
    if (!next) micVolume.value = 0
  } else if (!isMicOn.value) {
    try {
      const audioStream = await navigator.mediaDevices.getUserMedia({ audio: true })
      if (localStream.value) {
        audioStream.getAudioTracks().forEach(t => localStream.value?.addTrack(t))
      } else {
        localStream.value = audioStream
      }
      setupAudioAnalyser(audioStream)
      isMicOn.value = true
    } catch (e: any) {
      toast.warning('Mikrofon', 'Tidak dapat mengakses mikrofon: ' + (e.message || 'Izin ditolak'))
    }
  } else {
    isMicOn.value = false
    micVolume.value = 0
  }

  const local = participants.value.find(p => p.isLocal)
  if (local) local.isAudioEnabled = isMicOn.value
  toast.info(isMicOn.value ? 'Mikrofon aktif' : 'Mikrofon dibisukan')
}

async function toggleCamera() {
  if (isCameraOn.value) {
    if (localStream.value) {
      localStream.value.getVideoTracks().forEach(t => {
        t.stop()
        localStream.value?.removeTrack(t)
      })
    }
    isCameraOn.value = false
    toast.info('Kamera dinonaktifkan')
  } else {
    try {
      const videoStream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: { ideal: 1280 },
          height: { ideal: 720 },
          facingMode: 'user'
        }
      })
      const vTrack = videoStream.getVideoTracks()[0]
      if (localStream.value) {
        localStream.value.addTrack(vTrack)
      } else {
        localStream.value = videoStream
      }
      isCameraOn.value = true

      await nextTick()
      if (localVideoEl.value && localStream.value) {
        localVideoEl.value.srcObject = localStream.value
        localVideoEl.value.play().catch(e => console.warn('Play camera catch:', e))
      }
      toast.success('Kamera aktif')
    } catch (e: any) {
      toast.warning('Kamera', 'Tidak dapat mengakses kamera: ' + (e.message || 'Izin ditolak'))
    }
  }

  const local = participants.value.find(p => p.isLocal)
  if (local) local.isVideoEnabled = isCameraOn.value
}

async function toggleScreenShare() {
  if (isScreenSharing.value) {
    if (screenStream.value) {
      screenStream.value.getTracks().forEach(t => t.stop())
      screenStream.value = null
    }
    isScreenSharing.value = false
    toast.info('Berbagi layar dihentikan')
  } else {
    try {
      if (!navigator.mediaDevices?.getDisplayMedia) {
        toast.error('Browser tidak mendukung screen share')
        return
      }
      const stream = await navigator.mediaDevices.getDisplayMedia({ video: true })
      screenStream.value = stream
      isScreenSharing.value = true

      stream.getVideoTracks()[0].onended = () => {
        isScreenSharing.value = false
        screenStream.value = null
      }

      await nextTick()
      if (screenVideoEl.value) {
        screenVideoEl.value.srcObject = stream
        screenVideoEl.value.play().catch(e => console.warn('Play screen catch:', e))
      }
      toast.success('Berbagi Layar Dimulai')
    } catch (e: any) {
      isScreenSharing.value = false
      if (e.name !== 'NotAllowedError') {
        toast.warning('Berbagi Layar', e.message || 'Gagal membagikan layar')
      }
    }
  }
}

function leaveMeeting() {
  if (animFrame) cancelAnimationFrame(animFrame)
  if (audioContext && audioContext.state !== 'closed') audioContext.close()
  if (localStream.value) localStream.value.getTracks().forEach(t => t.stop())
  if (screenStream.value) screenStream.value.getTracks().forEach(t => t.stop())
  const returnUrl = meeting.value?.class_id ? `/classes/${meeting.value.class_id}` : '/dashboard'
  navigateTo(returnUrl)
}

async function endMeetingForAll() {
  isEnding.value = true
  try {
    if (animFrame) cancelAnimationFrame(animFrame)
    if (audioContext && audioContext.state !== 'closed') audioContext.close()
    if (localStream.value) localStream.value.getTracks().forEach(t => t.stop())
    if (screenStream.value) screenStream.value.getTracks().forEach(t => t.stop())
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

function copyMeetingLink() {
  if (typeof window !== 'undefined') {
    navigator.clipboard.writeText(window.location.href)
    toast.success('Tautan Disalin', 'Link meeting berhasil disalin ke clipboard')
  }
}

const isHost = computed(() => {
  return meeting.value?.host_id === auth.user?.id || meeting.value?.host?.id === auth.user?.id || auth.isAdmin
})

// Dynamic Adaptive Grid Classes
const gridClasses = computed(() => {
  if (isScreenSharing.value) {
    return 'w-full lg:w-80 grid grid-cols-2 lg:grid-cols-1 gap-3 overflow-y-auto'
  }
  const count = participants.value.length
  if (count <= 1) {
    return 'flex-1 max-w-4xl w-full mx-auto flex items-center justify-center p-2'
  }
  if (count === 2) {
    return 'flex-1 max-w-5xl w-full mx-auto grid grid-cols-1 md:grid-cols-2 gap-4 items-center justify-center p-2'
  }
  if (count <= 4) {
    return 'flex-1 max-w-5xl w-full mx-auto grid grid-cols-2 gap-4 items-center justify-center p-2'
  }
  return 'flex-1 w-full mx-auto grid grid-cols-2 lg:grid-cols-3 gap-4 items-center justify-center p-2'
})

const tileClasses = computed(() => {
  if (isScreenSharing.value) {
    return 'aspect-video w-full'
  }
  const count = participants.value.length
  if (count <= 1) {
    return 'w-full max-w-3xl aspect-video max-h-[75vh]'
  }
  return 'aspect-video w-full'
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
          class="p-2 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 transition-colors cursor-pointer"
          :class="{ 'bg-surface-800 text-white': activeSidebar === 'participants' }"
          title="Daftar Peserta"
          @click="activeSidebar = activeSidebar === 'participants' ? null : 'participants'"
        >
          <Users class="w-5 h-5" />
        </button>

        <button
          type="button"
          class="p-2 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 transition-colors cursor-pointer"
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
      <!-- Waiting / Solo Banner Indicator -->
      <Transition enter-active-class="transition duration-200 ease-out" enter-from-class="opacity-0 -translate-y-4" enter-to-class="opacity-100 translate-y-0">
        <div v-if="participants.length <= 1 && !isScreenSharing && !isLoading && !error" class="absolute top-4 left-1/2 -translate-x-1/2 z-20 px-4 py-2 bg-surface-900/90 backdrop-blur-md border border-surface-700/60 rounded-full text-xs text-surface-300 flex items-center gap-2.5 shadow-xl">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          <span>Ruang kelas aktif · Menunggu siswa bergabung</span>
          <button @click="copyMeetingLink" class="ml-1 px-2.5 py-1 rounded-md bg-brand-600 hover:bg-brand-500 text-white font-semibold text-[11px] transition-colors flex items-center gap-1 cursor-pointer">
            <Copy class="w-3 h-3" />
            <span>Salin Link</span>
          </button>
        </div>
      </Transition>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center gap-3">
        <div class="w-10 h-10 border-4 border-brand-500 border-t-transparent rounded-full animate-spin" />
        <p class="text-sm text-surface-400">Menyambungkan ke ruang LiveKit WebRTC...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex-1 flex flex-col items-center justify-center p-6 text-center gap-4">
        <div class="w-14 h-14 rounded-2xl bg-rose-950/50 border border-rose-900/60 text-rose-500 flex items-center justify-center">
          <PhoneOff class="w-7 h-7" />
        </div>
        <p class="text-base font-semibold text-rose-400">{{ error }}</p>
        <UiButton variant="outline" size="sm" @click="leaveMeeting">Kembali ke Kelas</UiButton>
      </div>

      <!-- Live Meeting Stage -->
      <div v-else class="flex-1 flex flex-col lg:flex-row p-4 gap-4 overflow-hidden items-center justify-center">
        <!-- Screen Share Stage if active -->
        <div v-if="isScreenSharing" class="flex-1 w-full bg-black rounded-2xl border border-surface-800 overflow-hidden relative flex items-center justify-center">
          <video
            :ref="attachScreenVideo"
            autoplay
            playsinline
            class="w-full h-full object-contain"
          />
          <div class="absolute top-3 left-3 px-3 py-1 rounded-lg bg-black/70 backdrop-blur-sm text-xs flex items-center gap-2 text-emerald-400 border border-surface-800">
            <Monitor class="w-4 h-4" />
            <span>Layar yang Anda Bagikan</span>
          </div>
        </div>

        <!-- Video Grid (Centered & Adaptive) -->
        <div :class="gridClasses">
          <!-- Participant Tile -->
          <div
            v-for="p in participants"
            :key="p.id"
            :class="[
              tileClasses,
              'relative bg-surface-900 rounded-2xl overflow-hidden border transition-all duration-200 flex items-center justify-center shadow-elevated group',
              p.isLocal && micVolume > 8
                ? 'border-emerald-500 ring-2 ring-emerald-500/50 shadow-emerald-500/20'
                : 'border-surface-800'
            ]"
          >
            <!-- Real Local Camera Video Stream -->
            <video
              v-if="p.isLocal && isCameraOn"
              :ref="attachLocalVideo"
              autoplay
              playsinline
              muted
              class="w-full h-full object-cover -scale-x-100"
            />

            <!-- Video disabled avatar fallback -->
            <div v-else class="flex flex-col items-center justify-center gap-3">
              <UiAvatar :name="p.name" size="xl" class="ring-4 ring-surface-800" />
              <p class="text-xs text-surface-400">Kamera dinonaktifkan</p>
            </div>

            <!-- Participant Label & Status Bottom overlay -->
            <div class="absolute bottom-3 left-3 right-3 flex items-center justify-between px-3 py-1.5 bg-black/60 backdrop-blur-md rounded-xl text-xs z-10">
              <div class="flex items-center gap-2 truncate">
                <span class="font-medium truncate">{{ p.name }} {{ p.isLocal ? '(Anda)' : '' }}</span>
                <span v-if="p.isHost" class="px-1.5 py-0.5 rounded bg-brand-600 text-[10px] font-bold">Host</span>
              </div>

              <div class="flex items-center gap-2 shrink-0">
                <!-- Audio Equalizer Indicator -->
                <div v-if="p.isLocal && isMicOn" class="flex items-end gap-0.5 h-3">
                  <span class="w-1 bg-emerald-400 rounded-full transition-all duration-75" :style="{ height: `${Math.max(20, Math.min(100, micVolume * 1.5))}%` }" />
                  <span class="w-1 bg-emerald-400 rounded-full transition-all duration-75" :style="{ height: `${Math.max(30, Math.min(100, micVolume * 2.2))}%` }" />
                  <span class="w-1 bg-emerald-400 rounded-full transition-all duration-75" :style="{ height: `${Math.max(20, Math.min(100, micVolume * 1.1))}%` }" />
                </div>

                <Mic v-if="p.isAudioEnabled" class="w-3.5 h-3.5 text-emerald-400" />
                <MicOff v-else class="w-3.5 h-3.5 text-rose-400" />
              </div>
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
            class="text-surface-400 hover:text-white text-xs cursor-pointer"
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
            <div class="flex items-center gap-2.5 min-w-0">
              <UiAvatar :name="p.name" size="sm" />
              <div class="min-w-0">
                <p class="text-xs font-semibold truncate">{{ p.name }} {{ p.isLocal ? '(Anda)' : '' }}</p>
                <span v-if="p.isHost" class="text-[10px] text-brand-400 font-medium">Pengajar / Host</span>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <Mic v-if="p.isAudioEnabled" class="w-3.5 h-3.5 text-emerald-400" />
              <MicOff v-else class="w-3.5 h-3.5 text-rose-400" />
              <Video v-if="p.isVideoEnabled" class="w-3.5 h-3.5 text-emerald-400" />
              <VideoOff v-else class="w-3.5 h-3.5 text-surface-500" />
            </div>
          </div>
        </div>

        <!-- Chat Content -->
        <div v-else class="flex-1 flex flex-col p-4">
          <div class="flex-1 flex items-center justify-center text-center text-xs text-surface-400">
            Pesan dalam sesi live akan tersinkronisasi realtime di sini.
          </div>
        </div>
      </aside>
    </div>

    <!-- Bottom Control Bar -->
    <footer class="h-20 px-6 bg-surface-900 border-t border-surface-800 flex items-center justify-between z-10">
      <div class="hidden sm:flex items-center gap-3 text-xs text-surface-400">
        <div class="flex items-center gap-1.5">
          <Shield class="w-4 h-4 text-emerald-400" />
          <span>WebRTC Aktif</span>
        </div>

        <!-- Live Volume Visual Meter -->
        <div v-if="isMicOn" class="flex items-center gap-2 pl-3 border-l border-surface-800">
          <Volume2 class="w-3.5 h-3.5 text-emerald-400" />
          <div class="w-20 h-1.5 bg-surface-800 rounded-full overflow-hidden">
            <div
              class="h-full bg-emerald-400 transition-all duration-75 rounded-full"
              :style="{ width: `${micVolume}%` }"
            />
          </div>
          <span class="text-[10px] font-mono text-emerald-400">{{ micVolume }}%</span>
        </div>
      </div>

      <!-- Main Action Controls -->
      <div class="flex items-center gap-3 mx-auto sm:mx-0">
        <!-- Mic Toggle -->
        <button
          type="button"
          class="w-12 h-12 rounded-2xl flex items-center justify-center transition-all shadow-md cursor-pointer relative"
          :class="isMicOn
            ? 'bg-surface-800 hover:bg-surface-700 text-white'
            : 'bg-rose-600 hover:bg-rose-700 text-white'"
          title="Nyalakan/Matikan Mic"
          @click="toggleMic"
        >
          <Mic v-if="isMicOn" class="w-5 h-5" />
          <MicOff v-else class="w-5 h-5" />
          
          <!-- Mic Ping dot when speaking -->
          <span
            v-if="isMicOn && micVolume > 8"
            class="absolute -top-1 -right-1 w-3 h-3 rounded-full bg-emerald-400 ring-2 ring-surface-900 animate-ping"
          />
        </button>

        <!-- Camera Toggle -->
        <button
          type="button"
          class="w-12 h-12 rounded-2xl flex items-center justify-center transition-all shadow-md cursor-pointer"
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
          class="w-12 h-12 rounded-2xl flex items-center justify-center transition-all shadow-md cursor-pointer"
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
          class="px-5 h-12 rounded-2xl bg-rose-600 hover:bg-rose-700 text-white font-semibold text-xs flex items-center gap-2 shadow-md shadow-rose-600/20 transition-all cursor-pointer"
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

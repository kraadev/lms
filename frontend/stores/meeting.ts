import { defineStore } from 'pinia'
import type { Meeting } from '~/types'

export const useMeetingStore = defineStore('meeting', () => {
  const currentMeeting = ref<Meeting | null>(null)
  const isConnecting = ref(false)
  const isConnected = ref(false)
  const isEnded = ref(false)
  const error = ref<string | null>(null)

  // Local media states
  const isMicEnabled = ref(true)
  const isCamEnabled = ref(true)
  const isScreenSharing = ref(false)
  const selectedAudioInput = ref<string>('')
  const selectedAudioOutput = ref<string>('')
  const selectedVideoInput = ref<string>('')

  function setMeeting(meeting: Meeting | null) {
    currentMeeting.value = meeting
    if (meeting?.type === 'audio') {
      isCamEnabled.value = false
    }
  }

  function reset() {
    currentMeeting.value = null
    isConnecting.value = false
    isConnected.value = false
    isEnded.value = false
    error.value = null
    isScreenSharing.value = false
  }

  return {
    currentMeeting,
    isConnecting,
    isConnected,
    isEnded,
    error,
    isMicEnabled,
    isCamEnabled,
    isScreenSharing,
    selectedAudioInput,
    selectedAudioOutput,
    selectedVideoInput,
    setMeeting,
    reset
  }
})

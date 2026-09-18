import { ref } from 'vue'

export function useNotificationSound() {
  const isMuted = ref(false)

  if (import.meta.client) {
    isMuted.value = localStorage.getItem('lms_sound_muted') === 'true'
  }

  function toggleMute() {
    isMuted.value = !isMuted.value
    if (import.meta.client) {
      localStorage.setItem('lms_sound_muted', String(isMuted.value))
    }
  }

  // Play pleasant synthetic notification chime using Web Audio API (zero external assets required)
  function playChime() {
    if (isMuted.value || !import.meta.client) return

    try {
      const AudioContext = window.AudioContext || (window as any).webkitAudioContext
      if (!AudioContext) return

      const ctx = new AudioContext()
      const osc = ctx.createOscillator()
      const gain = ctx.createGain()

      osc.type = 'sine'
      osc.frequency.setValueAtTime(587.33, ctx.currentTime) // D5
      osc.frequency.exponentialRampToValueAtTime(880, ctx.currentTime + 0.15) // A5

      gain.gain.setValueAtTime(0.15, ctx.currentTime)
      gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.35)

      osc.connect(gain)
      gain.connect(ctx.destination)

      osc.start()
      osc.stop(ctx.currentTime + 0.35)
    } catch {
      // AudioContext policy suppression fallback
    }
  }

  return {
    isMuted,
    toggleMute,
    playChime
  }
}

import { ref } from 'vue'

export function useFeatureFlags() {
  const flags = ref<Record<string, boolean>>({
    video_meetings: true,
    chat_reactions: true,
    ai_grader: false,
    audit_logging: true
  })

  function isEnabled(key: string): boolean {
    return !!flags.value[key]
  }

  function setFlag(key: string, val: boolean) {
    flags.value[key] = val
  }

  return {
    flags,
    isEnabled,
    setFlag
  }
}

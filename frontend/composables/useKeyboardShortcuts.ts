import { onMounted, onUnmounted } from 'vue'

export function useKeyboardShortcuts(shortcuts: Record<string, () => void>) {
  function handleKeyDown(e: KeyboardEvent) {
    // Ignore input textfields
    const target = e.target as HTMLElement
    if (['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName)) return

    const isCtrlOrCmd = e.ctrlKey || e.metaKey
    if (isCtrlOrCmd && e.key.toLowerCase() === 'k' && shortcuts['ctrl+k']) {
      e.preventDefault()
      shortcuts['ctrl+k']()
    } else if (e.key === 'Escape' && shortcuts['escape']) {
      shortcuts['escape']()
    } else if (e.key === '?' && shortcuts['?']) {
      shortcuts['?']()
    }
  }

  onMounted(() => {
    if (import.meta.client) {
      window.addEventListener('keydown', handleKeyDown)
    }
  })

  onUnmounted(() => {
    if (import.meta.client) {
      window.removeEventListener('keydown', handleKeyDown)
    }
  })
}

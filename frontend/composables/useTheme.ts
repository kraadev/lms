export type ThemeMode = 'light' | 'dark' | 'system'

export function useTheme() {
  const themeMode = useState<ThemeMode>('theme-mode', () => 'system')
  const isDark = useState<boolean>('is-dark', () => false)

  function applyTheme(mode: ThemeMode) {
    if (import.meta.server) return

    themeMode.value = mode
    localStorage.setItem('lms_theme', mode)

    let dark = false
    if (mode === 'system') {
      dark = window.matchMedia('(prefers-color-scheme: dark)').matches
    } else {
      dark = mode === 'dark'
    }

    isDark.value = dark
    if (dark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  function initTheme() {
    if (import.meta.server) return

    const saved = localStorage.getItem('lms_theme') as ThemeMode | null
    const initialMode = saved || 'system'
    applyTheme(initialMode)

    // Listen for system theme changes if in system mode
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      if (themeMode.value === 'system') {
        isDark.value = e.matches
        if (e.matches) {
          document.documentElement.classList.add('dark')
        } else {
          document.documentElement.classList.remove('dark')
        }
      }
    })
  }

  return {
    themeMode,
    isDark,
    applyTheme,
    initTheme
  }
}

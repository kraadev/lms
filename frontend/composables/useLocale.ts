import { ref } from 'vue'

export type LocaleCode = 'id' | 'en'

export function useLocale() {
  const currentLocale = ref<LocaleCode>('id')

  if (import.meta.client) {
    const saved = localStorage.getItem('lms_locale') as LocaleCode
    if (saved === 'id' || saved === 'en') {
      currentLocale.value = saved
    }
  }

  function setLocale(loc: LocaleCode) {
    currentLocale.value = loc
    if (import.meta.client) {
      localStorage.setItem('lms_locale', loc)
    }
  }

  return {
    currentLocale,
    setLocale
  }
}

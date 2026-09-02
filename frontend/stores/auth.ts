import { defineStore } from 'pinia'
import { authService } from '~/services/auth'
import type { User, UserRole } from '~/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isInitialized = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value)
  const role = computed<UserRole | null>(() => user.value?.role || null)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isTeacher = computed(() => user.value?.role === 'teacher')
  const isStudent = computed(() => user.value?.role === 'student')

  async function login(credentials: { email: string; password: string }): Promise<User> {
    isLoading.value = true
    error.value = null
    try {
      const response: any = await authService.login(credentials)
      
      const userData = response?.user || response?.data?.user || response
      const token = response?.token || response?.data?.token

      if (token && typeof window !== 'undefined') {
        localStorage.setItem('lms_token', token)
        localStorage.setItem('token', token)
        document.cookie = `lms_token=${token}; path=/; max-age=86400`
      }

      user.value = userData
      isInitialized.value = true
      return userData
    } catch (err: any) {
      error.value = err.message || 'Login gagal. Periksa kembali email dan password Anda.'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function logout(): Promise<void> {
    isLoading.value = true
    try {
      await authService.logout()
    } catch (err) {
      console.warn('Logout error (ignoring client state reset):', err)
    } finally {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('lms_token')
        localStorage.removeItem('token')
        document.cookie = 'lms_token=; path=/; max-age=0'
      }
      user.value = null
      isInitialized.value = true
      isLoading.value = false
      navigateTo('/login')
    }
  }

  async function restoreSession(): Promise<User | null> {
    if (isInitialized.value && user.value) return user.value

    isLoading.value = true
    try {
      const token = typeof window !== 'undefined' ? (localStorage.getItem('lms_token') || localStorage.getItem('token')) : null
      if (!token) {
        user.value = null
        return null
      }

      const currentUser: any = await authService.getMe()
      const userData = currentUser?.user || currentUser?.data || currentUser
      if (userData && userData.id) {
        user.value = userData
        return userData
      }

      user.value = null
      return null
    } catch (err) {
      user.value = null
      return null
    } finally {
      isLoading.value = false
      isInitialized.value = true
    }
  }

  function setUser(newUser: User | null) {
    user.value = newUser
  }

  return {
    user,
    isInitialized,
    isLoading,
    error,
    isAuthenticated,
    role,
    isAdmin,
    isTeacher,
    isStudent,
    login,
    logout,
    restoreSession,
    setUser
  }
})

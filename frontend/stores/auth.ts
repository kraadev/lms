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
      const response = await authService.login(credentials)
      user.value = response.user
      return response.user
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
      user.value = null
      isLoading.value = false
      navigateTo('/login')
    }
  }

  async function restoreSession(): Promise<User | null> {
    if (isInitialized.value) return user.value

    isLoading.value = true
    try {
      const currentUser = await authService.getMe()
      user.value = currentUser
      return currentUser
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

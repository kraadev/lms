<script setup lang="ts">
definePageMeta({
  layout: 'auth',
  middleware: 'guest'
})

useSeoMeta({
  title: 'Login - LMS Portal',
  description: 'Masuk ke sistem Learning Management System'
})

const auth = useAuthStore()
const { push: toastPush } = useToast()
const route = useRoute()

const form = reactive({
  email: '',
  password: ''
})

const errors = reactive({
  email: '',
  password: ''
})

const isLoading = ref(false)
const serverError = ref('')

function validate(): boolean {
  errors.email = ''
  errors.password = ''
  let ok = true

  if (!form.email) {
    errors.email = 'Email wajib diisi'
    ok = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = 'Format email tidak valid'
    ok = false
  }

  if (!form.password) {
    errors.password = 'Password wajib diisi'
    ok = false
  } else if (form.password.length < 6) {
    errors.password = 'Password minimal 6 karakter'
    ok = false
  }

  return ok
}

async function handleLogin() {
  if (!validate()) return

  isLoading.value = true
  serverError.value = ''

  try {
    const user = await auth.login({ email: form.email, password: form.password })
    
    // Connect WebSocket after login
    const { connect } = useWebSocket()
    connect()

    // Fetch notifications
    const notifStore = useNotificationsStore()
    notifStore.fetchNotifications()

    const redirect = (route.query.redirect as string) || '/dashboard'
    await navigateTo(redirect)
  } catch (err: any) {
    if (err?.status === 401 || err?.status === 403) {
      serverError.value = 'Email atau password yang Anda masukkan salah.'
    } else {
      serverError.value = err?.message || 'Gagal masuk. Coba lagi nanti.'
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="w-full max-w-md">
    <!-- Logo / Brand -->
    <div class="flex flex-col items-center mb-8">
      <div class="w-12 h-12 rounded-2xl bg-brand-600 flex items-center justify-center mb-4 shadow-lg">
        <svg class="w-7 h-7 text-white" fill="currentColor" viewBox="0 0 20 20">
          <path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838l-2.727 1.17 1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3zM3.31 9.397L5 10.12v4.102a8.969 8.969 0 00-1.05-.174 1 1 0 01-.89-.89 11.115 11.115 0 01.25-3.762zm5.99 7.176A9.026 9.026 0 007 14.935v-3.957l1.818.78a3 3 0 002.364 0l5.508-2.361a11.026 11.026 0 01.25 3.762 1 1 0 01-.89.89 8.968 8.968 0 00-5.35 2.524 1 1 0 01-1.4 0zM6 18a1 1 0 001-1v-2.065a8.935 8.935 0 00-2-.712V17a1 1 0 001 1z"/>
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-surface-900 dark:text-white">LMS Portal</h1>
      <p class="mt-1 text-sm text-surface-500 dark:text-surface-400">Masuk ke akun Anda untuk melanjutkan</p>
    </div>

    <!-- Card -->
    <div class="bg-white dark:bg-surface-900 rounded-2xl shadow-elevated border border-surface-200 dark:border-surface-800 p-7">
      <!-- Server error alert -->
      <div v-if="serverError" class="mb-5 flex items-start gap-3 p-3.5 bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900 rounded-xl" role="alert">
        <svg class="w-4.5 h-4.5 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
        </svg>
        <p class="text-sm text-rose-700 dark:text-rose-300">{{ serverError }}</p>
      </div>

      <form @submit.prevent="handleLogin" novalidate>
        <div class="space-y-4">
          <UiInput
            id="email"
            v-model="form.email"
            type="email"
            label="Alamat Email"
            placeholder="nama@sekolah.id"
            autocomplete="email"
            :error="errors.email"
            required
          />
          <UiInput
            id="password"
            v-model="form.password"
            type="password"
            label="Password"
            placeholder="Masukkan password"
            autocomplete="current-password"
            :error="errors.password"
            required
          />
        </div>

        <UiButton
          type="submit"
          class="w-full mt-6"
          :loading="isLoading"
          block
          size="lg"
        >
          {{ isLoading ? 'Memverifikasi...' : 'Masuk' }}
        </UiButton>
      </form>
    </div>

    <p class="text-center mt-5 text-xs text-surface-400 dark:text-surface-600">
      Sistem Manajemen Pembelajaran &copy; {{ new Date().getFullYear() }}
    </p>
  </div>
</template>

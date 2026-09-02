<script setup lang="ts">
import { Menu, X } from 'lucide-vue-next'

const showMobileMenu = ref(false)

// Close on route change
const route = useRoute()
watch(() => route.path, () => {
  showMobileMenu.value = false
})
</script>

<template>
  <div class="flex min-h-screen bg-surface-50 dark:bg-surface-950">
    <!-- Desktop Sidebar -->
    <aside class="hidden lg:flex flex-col w-60 shrink-0 border-r border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 h-screen sticky top-0 overflow-y-auto">
      <!-- Logo -->
      <div class="flex items-center gap-2.5 px-5 h-14 border-b border-surface-200 dark:border-surface-800 shrink-0">
        <div class="w-7 h-7 rounded-lg bg-brand-600 flex items-center justify-center">
          <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20"><path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838l-2.727 1.17 1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3zM3.31 9.397L5 10.12v4.102a8.969 8.969 0 00-1.05-.174 1 1 0 01-.89-.89 11.115 11.115 0 01.25-3.762zm5.99 7.176A9.026 9.026 0 007 14.935v-3.957l1.818.78a3 3 0 002.364 0l5.508-2.361a11.026 11.026 0 01.25 3.762 1 1 0 01-.89.89 8.968 8.968 0 00-5.35 2.524 1 1 0 01-1.4 0zM6 18a1 1 0 001-1v-2.065a8.935 8.935 0 00-2-.712V17a1 1 0 001 1z"/></svg>
        </div>
        <span class="font-bold text-surface-900 dark:text-white text-sm">LMS Portal</span>
      </div>
      <!-- Navigation -->
      <div class="flex-1 px-3 py-4">
        <AppSidebar />
      </div>
    </aside>

    <!-- Mobile Sidebar Overlay -->
    <Transition enter-active-class="transition-opacity duration-200" enter-from-class="opacity-0" enter-to-class="opacity-100" leave-active-class="transition-opacity duration-150" leave-from-class="opacity-100" leave-to-class="opacity-0">
      <div v-if="showMobileMenu" class="fixed inset-0 bg-black/50 z-40 lg:hidden" @click="showMobileMenu = false" aria-hidden="true" />
    </Transition>
    <Transition enter-active-class="transition duration-200 ease-out" enter-from-class="-translate-x-full" enter-to-class="translate-x-0" leave-active-class="transition duration-150 ease-in" leave-from-class="translate-x-0" leave-to-class="-translate-x-full">
      <aside v-if="showMobileMenu" class="fixed left-0 top-0 h-full w-72 bg-white dark:bg-surface-900 border-r border-surface-200 dark:border-surface-800 z-50 flex flex-col lg:hidden overflow-y-auto">
        <div class="flex items-center justify-between px-5 h-14 border-b border-surface-200 dark:border-surface-800 shrink-0">
          <div class="flex items-center gap-2.5">
            <div class="w-7 h-7 rounded-lg bg-brand-600 flex items-center justify-center">
              <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20"><path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838l-2.727 1.17 1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3z"/></svg>
            </div>
            <span class="font-bold text-surface-900 dark:text-white text-sm">LMS Portal</span>
          </div>
          <button type="button" class="p-1.5 rounded-lg hover:bg-surface-100 dark:hover:bg-surface-800 text-surface-500" aria-label="Tutup menu" @click="showMobileMenu = false">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="flex-1 px-3 py-4">
          <AppSidebar />
        </div>
      </aside>
    </Transition>

    <!-- Main content area -->
    <div class="flex-1 flex flex-col min-w-0">
      <AppTopbar>
        <template #mobile-trigger>
          <button
            type="button"
            class="lg:hidden p-2 rounded-lg text-surface-500 hover:bg-surface-100 dark:hover:bg-surface-800 hover:text-surface-800 dark:hover:text-surface-200 transition-colors"
            aria-label="Buka menu"
            @click="showMobileMenu = true"
          >
            <Menu class="w-5 h-5" />
          </button>
        </template>
      </AppTopbar>

      <main class="flex-1 overflow-auto">
        <slot />
      </main>
    </div>

    <!-- Toast container -->
    <UiToastContainer />
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: [
    function () {
      const auth = useAuthStore()
      if (auth.isAuthenticated) {
        return navigateTo('/dashboard')
      }
      return navigateTo('/login')
    }
  ]
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="w-8 h-8 border-4 border-brand-500 border-t-transparent rounded-full animate-spin" />
  </div>
</template>

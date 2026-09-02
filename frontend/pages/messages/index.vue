<script setup lang="ts">
import { MessageSquare, Send, BookOpen, User, Hash } from 'lucide-vue-next'
import { classesService } from '~/services/classes'
import { messagesService } from '~/services/messages'
import type { Class, ChatMessage } from '~/types'
import { formatRelativeTime } from '~/utils/formatters'
import ClassChatTab from '~/components/classes/ClassChatTab.vue'

definePageMeta({ middleware: 'auth' })
useSeoMeta({ title: 'Pesan & Diskusi Kelas - LMS' })

const auth = useAuthStore()
const classes = ref<Class[]>([])
const selectedClassId = ref<number | string>('')
const isLoading = ref(true)

async function loadData() {
  isLoading.value = true
  try {
    const list = await classesService.getAll({ status: 'active' })
    classes.value = Array.isArray(list) ? list : []
    if (classes.value.length > 0) {
      selectedClassId.value = classes.value[0].id
    }
  } catch (e) {
    classes.value = []
  } finally {
    isLoading.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <div class="p-4 md:p-6 max-w-6xl mx-auto">
    <div class="mb-6">
      <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100 flex items-center gap-2">
        <MessageSquare class="w-6 h-6 text-brand-600 dark:text-brand-400" />
        Pesan & Diskusi Realtime Kelas
      </h1>
      <p class="text-sm text-surface-500 dark:text-surface-400 mt-0.5">
        Ruang obrolan interaktif bertenaga native WebSocket per kelas
      </p>
    </div>

    <!-- Class Selector Bar -->
    <div class="flex items-center gap-2 overflow-x-auto pb-2 mb-4 scrollbar-thin">
      <button
        v-for="cls in classes"
        :key="cls.id"
        @click="selectedClassId = cls.id"
        :class="[
          'px-4 py-2 rounded-xl text-sm font-semibold whitespace-nowrap transition-all flex items-center gap-2 border',
          selectedClassId === cls.id
            ? 'bg-brand-600 text-white border-brand-600 shadow-sm'
            : 'bg-white dark:bg-surface-900 text-surface-700 dark:text-surface-300 border-surface-200 dark:border-surface-800 hover:border-brand-500/50'
        ]"
      >
        <Hash class="w-4 h-4 opacity-70" />
        {{ cls.title || cls.name }}
      </button>
    </div>

    <!-- Active Class Chat Component -->
    <div v-if="selectedClassId" class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 p-4 shadow-soft">
      <ClassChatTab :class-id="String(selectedClassId)" />
    </div>

    <div v-else-if="!isLoading" class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 p-8 shadow-soft">
      <UiEmptyState
        :icon="MessageSquare"
        title="Belum ada kelas aktif"
        description="Anda belum terdaftar dalam kelas manapun untuk memulai diskusi."
      />
    </div>
  </div>
</template>

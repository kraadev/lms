<script setup lang="ts">
import {
  BookOpen, ClipboardList, FileQuestion, Users,
  MessageSquare, Video, BookMarked, AlertCircle
} from 'lucide-vue-next'
import { classesService } from '~/services/classes'
import type { Class } from '~/types'

definePageMeta({ middleware: 'auth' })

const route = useRoute()
const classId = computed(() => route.params.id as string)

const cls = ref<Class | null>(null)
const isLoading = ref(true)
const error = ref<{ message: string; status?: number } | null>(null)
const activeTab = ref('overview')

const tabs = computed(() => [
  { key: 'overview', label: 'Ringkasan', icon: BookOpen },
  { key: 'materials', label: 'Materi', icon: BookMarked },
  { key: 'assignments', label: 'Tugas', icon: ClipboardList },
  { key: 'quizzes', label: 'Kuis', icon: FileQuestion },
  { key: 'members', label: 'Anggota', icon: Users },
  { key: 'chat', label: 'Chat', icon: MessageSquare },
  { key: 'meeting', label: 'Meeting', icon: Video }
])

useSeoMeta({ title: computed(() => cls.value?.title || 'Kelas') })

async function load() {
  isLoading.value = true
  error.value = null
  try {
    cls.value = await classesService.getById(classId.value)
  } catch (err: any) {
    error.value = { message: err?.message || 'Gagal memuat kelas', status: err?.status }
  } finally {
    isLoading.value = false
  }
}

onMounted(load)
watch(classId, load)
</script>

<template>
  <div>
    <!-- Loading -->
    <div v-if="isLoading" class="p-4 md:p-6">
      <UiSkeleton class="h-28 rounded-xl mb-4" />
      <UiSkeleton class="h-10 rounded-xl mb-6" />
      <UiSkeleton :rows="5" />
    </div>

    <!-- Error -->
    <div v-else-if="error" class="p-4 md:p-6">
      <UiErrorState :message="error.message" :status="error.status" @retry="load" />
    </div>

    <!-- Content -->
    <div v-else-if="cls">
      <!-- Class Header -->
      <div class="border-b border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900">
        <div class="px-4 md:px-6 pt-5 pb-0 max-w-6xl mx-auto">
          <div class="flex items-start gap-4 mb-4">
            <div class="w-12 h-12 rounded-xl bg-brand-100 dark:bg-brand-950/50 flex items-center justify-center shrink-0">
              <BookOpen class="w-6 h-6 text-brand-600 dark:text-brand-400" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <h1 class="text-lg font-bold text-surface-900 dark:text-surface-100">{{ cls.title }}</h1>
                <UiBadge :variant="cls.status === 'active' ? 'success' : 'default'" size="sm">{{ cls.status === 'active' ? 'Aktif' : 'Arsip' }}</UiBadge>
              </div>
              <div class="flex flex-wrap items-center gap-3 mt-1 text-xs text-surface-500 dark:text-surface-400">
                <span>{{ cls.academic_year }}</span>
                <span v-if="cls.teacher">&middot; {{ cls.teacher.name }}</span>
                <span>&middot; {{ cls.member_count || 0 }} anggota</span>
              </div>
            </div>
          </div>

          <!-- Active Meeting Banner in class header -->
          <div v-if="cls.active_meeting" class="flex items-center justify-between gap-3 px-3 py-2 bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-200 dark:border-emerald-900/40 rounded-lg mb-4">
            <div class="flex items-center gap-2 text-sm">
              <Video class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
              <span class="text-emerald-800 dark:text-emerald-200 font-medium">Kelas online sedang berlangsung:</span>
              <span class="text-emerald-700 dark:text-emerald-300">{{ cls.active_meeting.title }}</span>
            </div>
            <NuxtLink :to="`/meetings/${cls.active_meeting.id}`">
              <UiButton size="xs" variant="success">Gabung</UiButton>
            </NuxtLink>
          </div>

          <!-- Tabs -->
          <UiTabs v-model="activeTab" :tabs="tabs" />
        </div>
      </div>

      <!-- Tab Content -->
      <div class="max-w-6xl mx-auto">
        <ClassOverviewTab v-if="activeTab === 'overview'" :class-data="cls" />
        <ClassMaterialsTab v-else-if="activeTab === 'materials'" :class-id="classId" />
        <ClassAssignmentsTab v-else-if="activeTab === 'assignments'" :class-id="classId" />
        <ClassQuizzesTab v-else-if="activeTab === 'quizzes'" :class-id="classId" />
        <ClassMembersTab v-else-if="activeTab === 'members'" :class-id="classId" />
        <ClassChatTab v-else-if="activeTab === 'chat'" :class-id="classId" />
        <ClassMeetingTab v-else-if="activeTab === 'meeting'" :class-id="classId" :class-data="cls" @refresh="load" />
      </div>
    </div>
  </div>
</template>

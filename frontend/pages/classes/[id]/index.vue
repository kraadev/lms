<script setup lang="ts">
import {
  BookOpen, ClipboardList, FileQuestion, Users,
  MessageSquare, Video, BookMarked, AlertCircle, Copy, Check
} from 'lucide-vue-next'
import { classesService } from '~/services/classes'
import type { Class } from '~/types'

import ClassOverviewTab from '~/components/classes/ClassOverviewTab.vue'
import ClassMaterialsTab from '~/components/classes/ClassMaterialsTab.vue'
import ClassAssignmentsTab from '~/components/classes/ClassAssignmentsTab.vue'
import ClassQuizzesTab from '~/components/classes/ClassQuizzesTab.vue'
import ClassMembersTab from '~/components/classes/ClassMembersTab.vue'
import ClassChatTab from '~/components/classes/ClassChatTab.vue'
import ClassMeetingTab from '~/components/classes/ClassMeetingTab.vue'

definePageMeta({ middleware: 'auth' })

const auth = useAuthStore()
if (auth.isAdmin) {
  navigateTo('/admin/classes')
}

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

const isCopied = ref(false)
const toast = useToast()

function copyClassCode() {
  if (!cls.value?.code) return
  if (navigator.clipboard) {
    navigator.clipboard.writeText(cls.value.code)
  }
  isCopied.value = true
  toast.success('Kode kelas berhasil disalin ke clipboard')
  setTimeout(() => {
    isCopied.value = false
  }, 2000)
}

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
      <div class="border-b border-surface-200/80 dark:border-surface-800/80 bg-white/80 dark:bg-surface-900/80 backdrop-blur-md">
        <div class="px-4 md:px-6 pt-6 pb-0 max-w-6xl mx-auto">
          <div class="flex items-start gap-4.5 mb-5">
            <div class="w-14 h-14 rounded-2xl bg-gradient-to-tr from-brand-600 to-brand-400 text-white flex items-center justify-center shrink-0 shadow-md shadow-brand-500/25">
              <BookOpen class="w-7 h-7" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2.5 flex-wrap">
                <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-surface-900 dark:text-surface-100">{{ cls.title }}</h1>
                <UiBadge :variant="cls.status === 'active' ? 'success' : 'default'" size="sm">{{ cls.status === 'active' ? 'Aktif' : 'Arsip' }}</UiBadge>
                <button
                  v-if="cls.code"
                  type="button"
                  class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-mono font-semibold bg-brand-50 hover:bg-brand-100/80 dark:bg-brand-950/60 dark:hover:bg-brand-900/70 text-brand-700 dark:text-brand-300 border border-brand-200/60 dark:border-brand-800/60 transition-colors focus-ring cursor-pointer shadow-xs"
                  title="Klik untuk menyalin kode kelas"
                  @click="copyClassCode"
                >
                  <span>Kode: {{ cls.code }}</span>
                  <component :is="isCopied ? Check : Copy" class="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
                </button>
              </div>
              <div class="flex flex-wrap items-center gap-3 mt-1.5 text-xs text-surface-500 dark:text-surface-400 font-medium">
                <span class="inline-flex items-center gap-1 bg-surface-100 dark:bg-surface-800/60 px-2 py-0.5 rounded-md text-[11px]">{{ cls.academic_year }}</span>
                <span v-if="cls.teacher">&middot; Pengajar: {{ cls.teacher.name }}</span>
                <span>&middot; {{ cls.member_count || 0 }} siswa terdaftar</span>
              </div>
            </div>
          </div>

          <!-- Active Meeting Banner in class header -->
          <div v-if="cls.active_meeting" class="flex items-center justify-between gap-3 px-4 py-2.5 bg-gradient-to-r from-emerald-500/15 via-emerald-500/5 to-transparent border border-emerald-500/30 rounded-xl mb-4 shadow-xs">
            <div class="flex items-center gap-2.5 text-sm">
              <span class="w-2 h-2 rounded-full bg-emerald-500 animate-ping" />
              <Video class="w-4.5 h-4.5 text-emerald-600 dark:text-emerald-400 shrink-0" />
              <span class="text-emerald-900 dark:text-emerald-200 font-bold">Kelas Online Aktif:</span>
              <span class="text-emerald-700 dark:text-emerald-300 font-medium truncate">{{ cls.active_meeting.title }}</span>
            </div>
            <NuxtLink :to="`/meetings/${cls.active_meeting.id}`">
              <UiButton size="xs" variant="success" class="shadow-xs">Gabung Sekarang</UiButton>
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

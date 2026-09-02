<script setup lang="ts">
import { Users, Search, GraduationCap, Shield } from 'lucide-vue-next'
import { classesService } from '~/services/classes'
import type { ClassMember } from '~/types'
import { formatDate } from '~/utils/formatters'

const props = defineProps<{ classId: string }>()
const auth = useAuthStore()

const members = ref<ClassMember[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)
const searchQuery = ref('')

async function loadMembers() {
  isLoading.value = true
  error.value = null
  try {
    members.value = await classesService.getMembers(props.classId)
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat anggota kelas'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadMembers)

const teachers = computed(() => {
  return members.value.filter(m => (m.user?.role === 'teacher' || m.role === 'teacher') && 
    (!searchQuery.value || m.user?.name?.toLowerCase().includes(searchQuery.value.toLowerCase()) || m.user?.email?.toLowerCase().includes(searchQuery.value.toLowerCase())))
})

const students = computed(() => {
  return members.value.filter(m => (m.user?.role === 'student' || m.role === 'student' || (!m.user?.role && m.role !== 'teacher')) && 
    (!searchQuery.value || m.user?.name?.toLowerCase().includes(searchQuery.value.toLowerCase()) || m.user?.email?.toLowerCase().includes(searchQuery.value.toLowerCase())))
})
</script>

<template>
  <div class="p-4 md:p-6 space-y-6">
    <!-- Header with Search -->
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
      <div class="flex items-center gap-2">
        <Users class="w-5 h-5 text-brand-600 dark:text-brand-400" />
        <h2 class="text-base font-semibold text-surface-800 dark:text-surface-200">Anggota Kelas</h2>
        <span class="text-xs px-2 py-0.5 rounded-full bg-surface-100 dark:bg-surface-800 text-surface-600 dark:text-surface-400 font-medium">
          {{ members.length }} orang
        </span>
      </div>
      <div class="w-full sm:w-64">
        <UiInput
          v-model="searchQuery"
          placeholder="Cari anggota..."
          size="sm"
        />
      </div>
    </div>

    <div v-if="isLoading">
      <UiSkeleton :rows="5" />
    </div>

    <UiErrorState v-else-if="error" :message="error" @retry="loadMembers" />

    <div v-else class="space-y-6">
      <!-- Teachers Group -->
      <section>
        <h3 class="text-xs font-bold uppercase tracking-wider text-surface-500 dark:text-surface-400 mb-3 flex items-center gap-1.5">
          <Shield class="w-3.5 h-3.5 text-brand-600" />
          Pengajar ({{ teachers.length }})
        </h3>
        <div class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden">
          <div v-if="!teachers.length" class="p-4 text-xs text-surface-400">Tidak ada pengajar ditemukan</div>
          <div
            v-for="m in teachers"
            :key="m.id"
            class="flex items-center gap-3.5 p-3.5 hover:bg-surface-50 dark:hover:bg-surface-800/50 transition-colors"
          >
            <UiAvatar :name="m.user?.name || 'Guru'" :src="m.user?.avatar" size="md" />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate">{{ m.user?.name }}</p>
              <p class="text-xs text-surface-500 dark:text-surface-400 truncate">{{ m.user?.email }}</p>
            </div>
            <UiBadge variant="primary" size="sm">Guru</UiBadge>
          </div>
        </div>
      </section>

      <!-- Students Group -->
      <section>
        <h3 class="text-xs font-bold uppercase tracking-wider text-surface-500 dark:text-surface-400 mb-3 flex items-center gap-1.5">
          <GraduationCap class="w-3.5 h-3.5 text-surface-500" />
          Siswa ({{ students.length }})
        </h3>
        <div class="divide-y divide-surface-100 dark:divide-surface-800 rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden">
          <div v-if="!students.length" class="p-6 text-center text-xs text-surface-400">
            {{ searchQuery ? 'Tidak ada siswa yang cocok dengan pencarian' : 'Belum ada siswa yang bergabung ke kelas ini' }}
          </div>
          <div
            v-for="m in students"
            :key="m.id"
            class="flex items-center gap-3.5 p-3.5 hover:bg-surface-50 dark:hover:bg-surface-800/50 transition-colors"
          >
            <UiAvatar :name="m.user?.name || 'Siswa'" :src="m.user?.avatar" size="md" />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-surface-900 dark:text-surface-100 truncate">{{ m.user?.name }}</p>
              <p class="text-xs text-surface-500 dark:text-surface-400 truncate">{{ m.user?.email }}</p>
            </div>
            <span class="text-xs text-surface-400 hidden sm:inline">
              Bergabung {{ formatDate(m.joined_at || m.created_at, { day: 'numeric', month: 'short', year: 'numeric' }) }}
            </span>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

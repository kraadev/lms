<script setup lang="ts">
import { Clock, Shield, Check, X, User } from 'lucide-vue-next'

export interface WaitingParticipant {
  id: string
  name: string
  requestedAt: string
}

interface Props {
  isTeacher?: boolean
  meetingTitle?: string
  waitingList?: WaitingParticipant[]
}

withDefaults(defineProps<Props>(), {
  isTeacher: false,
  meetingTitle: 'Pertemuan Online',
  waitingList: () => []
})

const emit = defineEmits<{
  (e: 'admit', id: string): void
  (e: 'reject', id: string): void
}>()
</script>

<template>
  <div class="rounded-2xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 p-6 shadow-soft">
    <!-- Student Waiting View -->
    <div v-if="!isTeacher" class="text-center py-8 space-y-3">
      <div class="w-12 h-12 rounded-full bg-brand-50 dark:bg-brand-950/60 text-brand-600 flex items-center justify-center mx-auto animate-pulse">
        <Clock class="w-6 h-6" />
      </div>
      <h3 class="text-base font-bold text-surface-900 dark:text-surface-100">Ruang Tunggu Kelas</h3>
      <p class="text-xs text-surface-500 dark:text-surface-400 max-w-sm mx-auto">
        Pengajar akan segera mengizinkan Anda masuk ke dalam pertemuan: <span class="font-semibold">{{ meetingTitle }}</span>.
      </p>
    </div>

    <!-- Teacher Admission Queue View -->
    <div v-else class="space-y-4">
      <div class="flex items-center justify-between border-b border-surface-100 dark:border-surface-800 pb-3">
        <div class="flex items-center gap-2">
          <Shield class="w-4 h-4 text-brand-600" />
          <h3 class="text-xs font-bold uppercase tracking-wider text-surface-700 dark:text-surface-300">Antrean Masuk Siswa</h3>
        </div>
        <span class="text-xs font-semibold text-brand-600 dark:text-brand-400">{{ waitingList.length }} menunggu</span>
      </div>

      <div v-if="!waitingList.length" class="py-6 text-center text-xs text-surface-400">
        Tidak ada siswa di ruang tunggu.
      </div>
      <div v-else class="divide-y divide-surface-100 dark:divide-surface-800">
        <div v-for="p in waitingList" :key="p.id" class="flex items-center justify-between py-2.5">
          <div class="flex items-center gap-2.5">
            <div class="w-7 h-7 rounded-full bg-surface-100 dark:bg-surface-800 flex items-center justify-center text-surface-500">
              <User class="w-3.5 h-3.5" />
            </div>
            <span class="text-xs font-medium text-surface-800 dark:text-surface-200">{{ p.name }}</span>
          </div>

          <div class="flex items-center gap-1.5">
            <button
              type="button"
              class="p-1 rounded-lg bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 hover:bg-emerald-100"
              title="Izinkan Masuk"
              @click="emit('admit', p.id)"
            >
              <Check class="w-4 h-4" />
            </button>
            <button
              type="button"
              class="p-1 rounded-lg bg-rose-50 dark:bg-rose-950/50 text-rose-600 hover:bg-rose-100"
              title="Tolak"
              @click="emit('reject', p.id)"
            >
              <X class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

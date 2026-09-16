<script setup lang="ts">
import { ref } from 'vue'
import { ChevronDown, ChevronRight, CheckCircle2, Circle, FileText, Video, BookOpen, Clock } from 'lucide-vue-next'

export interface LessonItem {
  id: number
  title: string
  contentType: 'markdown' | 'video' | 'pdf'
  durationSec?: number
  isCompleted?: boolean
}

export interface ModuleItem {
  id: number
  title: string
  lessons: LessonItem[]
}

interface Props {
  modules: ModuleItem[]
  activeLessonId?: number
}

const props = withDefaults(defineProps<Props>(), {
  modules: () => [],
  activeLessonId: undefined
})

const emit = defineEmits<{
  (e: 'select-lesson', lesson: LessonItem): void
  (e: 'toggle-complete', lessonId: number): void
}>()

// Open all modules by default
const openModules = ref<Record<number, boolean>>(
  props.modules.reduce((acc, m) => ({ ...acc, [m.id]: true }), {})
)

function toggleModule(id: number) {
  openModules.value[id] = !openModules.value[id]
}

function formatDuration(sec?: number): string {
  if (!sec) return ''
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}:${s < 10 ? '0' : ''}${s}`
}
</script>

<template>
  <div class="space-y-3">
    <div
      v-for="(mod, idx) in modules"
      :key="mod.id"
      class="rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft"
    >
      <!-- Module Header -->
      <button
        type="button"
        class="w-full flex items-center justify-between px-4 py-3 bg-surface-50/70 dark:bg-surface-800/50 hover:bg-surface-100/80 dark:hover:bg-surface-800 transition-colors text-left"
        @click="toggleModule(mod.id)"
      >
        <div class="flex items-center gap-2.5 min-w-0">
          <component
            :is="openModules[mod.id] ? ChevronDown : ChevronRight"
            class="w-4 h-4 text-surface-400 shrink-0"
          />
          <span class="text-xs font-bold uppercase tracking-wider text-brand-600 dark:text-brand-400">
            Modul {{ idx + 1 }}
          </span>
          <span class="text-sm font-semibold text-surface-900 dark:text-surface-100 truncate">
            {{ mod.title }}
          </span>
        </div>
        <span class="text-xs text-surface-500 dark:text-surface-400 shrink-0">
          {{ mod.lessons.length }} materi
        </span>
      </button>

      <!-- Lessons List -->
      <div v-show="openModules[mod.id]" class="divide-y divide-surface-100 dark:divide-surface-800/60">
        <div
          v-for="lesson in mod.lessons"
          :key="lesson.id"
          :class="[
            'flex items-center justify-between px-4 py-2.5 transition-colors cursor-pointer group',
            activeLessonId === lesson.id
              ? 'bg-brand-50/80 dark:bg-brand-950/60 text-brand-700 dark:text-brand-300 font-medium'
              : 'hover:bg-surface-50 dark:hover:bg-surface-800/40 text-surface-700 dark:text-surface-300'
          ]"
          @click="emit('select-lesson', lesson)"
        >
          <div class="flex items-center gap-3 min-w-0">
            <!-- Complete Checkbox Toggle -->
            <button
              type="button"
              class="text-surface-400 hover:text-emerald-600 dark:hover:text-emerald-400 transition-colors shrink-0"
              :title="lesson.isCompleted ? 'Tandai belum selesai' : 'Tandai sudah selesai'"
              @click.stop="emit('toggle-complete', lesson.id)"
            >
              <CheckCircle2 v-if="lesson.isCompleted" class="w-4.5 h-4.5 text-emerald-600 dark:text-emerald-400" />
              <Circle v-else class="w-4.5 h-4.5 text-surface-300 dark:text-surface-600" />
            </button>

            <!-- Icon indicator -->
            <Video v-if="lesson.contentType === 'video'" class="w-4 h-4 text-brand-500 shrink-0" />
            <FileText v-else-if="lesson.contentType === 'pdf'" class="w-4 h-4 text-amber-500 shrink-0" />
            <BookOpen v-else class="w-4 h-4 text-emerald-500 shrink-0" />

            <span class="text-xs truncate">
              {{ lesson.title }}
            </span>
          </div>

          <!-- Duration Badge -->
          <div v-if="lesson.durationSec" class="flex items-center gap-1 text-[11px] text-surface-400 shrink-0 ml-2">
            <Clock class="w-3 h-3" />
            <span>{{ formatDuration(lesson.durationSec) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

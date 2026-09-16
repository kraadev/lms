import { ref, computed } from 'vue'

export function useCourseProgress(courseId: number | string, totalLessons: number) {
  const storageKey = `lms_course_progress_${courseId}`
  const completedIds = ref<number[]>([])

  if (import.meta.client) {
    try {
      const saved = localStorage.getItem(storageKey)
      if (saved) {
        completedIds.value = JSON.parse(saved)
      }
    } catch {
      completedIds.value = []
    }
  }

  function save() {
    if (import.meta.client) {
      localStorage.setItem(storageKey, JSON.stringify(completedIds.value))
    }
  }

  function toggleLesson(lessonId: number) {
    const idx = completedIds.value.indexOf(lessonId)
    if (idx >= 0) {
      completedIds.value.splice(idx, 1)
    } else {
      completedIds.value.push(lessonId)
    }
    save()
  }

  function isLessonComplete(lessonId: number): boolean {
    return completedIds.value.includes(lessonId)
  }

  const completedCount = computed(() => completedIds.value.length)
  const progressPercent = computed(() => {
    if (totalLessons <= 0) return 0
    const pct = Math.round((completedCount.value / totalLessons) * 100)
    return Math.min(100, Math.max(0, pct))
  })

  const milestoneLabel = computed(() => {
    const p = progressPercent.value
    if (p >= 100) return 'Selesai 🏆'
    if (p >= 75) return 'Mahir ⭐'
    if (p >= 50) return 'Pertengahan 🚀'
    if (p >= 25) return 'Berkembang 💡'
    return 'Pemula 🌱'
  })

  return {
    completedIds,
    completedCount,
    progressPercent,
    milestoneLabel,
    toggleLesson,
    isLessonComplete
  }
}

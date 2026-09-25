import { ref, computed, type Ref } from 'vue'

export interface SearchableCourse {
  id: number
  title: string
  description?: string
  academicYear?: string
  teacherName?: string
}

export function useCourseSearch(courses: Ref<SearchableCourse[]>) {
  const query = ref('')

  const results = computed(() => {
    const q = query.value.toLowerCase().trim()
    if (!q) return courses.value

    return courses.value.filter((c) => {
      const matchTitle = c.title.toLowerCase().includes(q)
      const matchDesc = c.description?.toLowerCase().includes(q) || false
      const matchYear = c.academicYear?.toLowerCase().includes(q) || false
      const matchTeacher = c.teacherName?.toLowerCase().includes(q) || false
      return matchTitle || matchDesc || matchYear || matchTeacher
    })
  })

  return {
    query,
    results
  }
}

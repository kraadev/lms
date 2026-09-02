import { api } from './api'
import type { Class, ClassMember } from '~/types'

export const classesService = {
  getAll(params?: { status?: string; search?: string }): Promise<Class[]> {
    return api.get<Class[]>('/classes', params)
  },

  getById(id: number | string): Promise<Class> {
    return api.get<Class>(`/classes/${id}`)
  },

  create(payload: { title: string; description?: string; academic_year: string; semester?: string }): Promise<Class> {
    return api.post<Class>('/classes', payload)
  },

  update(id: number | string, payload: Partial<Class>): Promise<Class> {
    return api.put<Class>(`/classes/${id}`, payload)
  },

  delete(id: number | string): Promise<{ message: string }> {
    return api.delete<{ message: string }>(`/classes/${id}`)
  },

  getMembers(classId: number | string): Promise<ClassMember[]> {
    return api.get<ClassMember[]>(`/classes/${classId}/members`)
  },

  joinByCode(code: string): Promise<{ message: string; class: Class }> {
    return api.post<{ message: string; class: Class }>('/classes/join', { code })
  }
}

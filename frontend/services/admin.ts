import { api } from './api'
import type { User, UserRole, Class } from '~/types'

export const adminService = {
  getUsers(params?: { role?: UserRole; search?: string; status?: string }): Promise<User[]> {
    return api.get<User[]>('/users', params)
  },

  createUser(payload: { name: string; email: string; password?: string; role: UserRole; phone?: string }): Promise<User> {
    return api.post<User>('/users', payload)
  },

  updateUser(id: number | string, payload: Partial<User & { password?: string }>): Promise<User> {
    return api.put<User>(`/users/${id}`, payload)
  },

  deleteUser(id: number | string): Promise<{ message: string }> {
    return api.delete<{ message: string }>(`/users/${id}`)
  },

  getClasses(params?: { search?: string }): Promise<Class[]> {
    return api.get<Class[]>('/classes', params)
  },

  createClass(payload: { title: string; code: string; teacher_id: number; academic_year: string; description?: string }): Promise<Class> {
    return api.post<Class>('/admin/classes', payload)
  },

  updateClass(id: number | string, payload: Partial<Class>): Promise<Class> {
    return api.put<Class>(`/admin/classes/${id}`, payload)
  },

  deleteClass(id: number | string): Promise<{ message: string }> {
    return api.delete<{ message: string }>(`/admin/classes/${id}`)
  }
}

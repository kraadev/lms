import { api } from './api'
import type { Material } from '~/types'

export const materialsService = {
  getByClass(classId: number | string): Promise<Material[]> {
    return api.get<Material[]>(`/classes/${classId}/materials`)
  },

  getById(id: number | string): Promise<Material> {
    return api.get<Material>(`/materials/${id}`)
  },

  create(classId: number | string, data: FormData | { title: string; description?: string; external_link?: string }): Promise<Material> {
    return api.post<Material>(`/classes/${classId}/materials`, data)
  },

  update(id: number | string, data: FormData | Partial<Material>): Promise<Material> {
    return api.put<Material>(`/materials/${id}`, data)
  },

  delete(id: number | string): Promise<{ message: string }> {
    return api.delete<{ message: string }>(`/materials/${id}`)
  }
}

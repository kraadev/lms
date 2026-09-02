import { api } from './api'
import type { ChatMessage } from '~/types'

export const messagesService = {
  getByClass(classId: number | string, params?: { page?: number; limit?: number; before?: string }): Promise<{ messages: ChatMessage[]; has_more?: boolean; total?: number }> {
    return api.get<{ messages: ChatMessage[]; has_more?: boolean; total?: number }>(`/classes/${classId}/messages`, params)
  }
}

import { api } from './api'
import type { User, AuthResponse } from '~/types'

export const authService = {
  login(credentials: { email: string; password: string }): Promise<AuthResponse> {
    return api.post<AuthResponse>('/auth/login', credentials)
  },

  logout(): Promise<{ message: string }> {
    return api.post<{ message: string }>('/auth/logout')
  },

  getMe(): Promise<User> {
    return api.get<User>('/auth/me')
  }
}

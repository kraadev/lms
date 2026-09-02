import type { ApiError } from '~/types'

function normalizeData(obj: any): any {
  if (!obj || typeof obj !== 'object') return obj

  if (Array.isArray(obj)) {
    return obj.map(normalizeData)
  }

  const normalized: any = { ...obj }

  // Normalize Class: title <-> name
  if (normalized.name && !normalized.title) {
    normalized.title = normalized.name
  } else if (normalized.title && !normalized.name) {
    normalized.name = normalized.title
  }

  // Normalize Assignment: deadline <-> due_date
  if (normalized.deadline && !normalized.due_date) {
    normalized.due_date = normalized.deadline
  } else if (normalized.due_date && !normalized.deadline) {
    normalized.deadline = normalized.due_date
  }

  // Recursively normalize child objects
  for (const key of Object.keys(normalized)) {
    if (normalized[key] && typeof normalized[key] === 'object') {
      normalized[key] = normalizeData(normalized[key])
    }
  }

  return normalized
}

export class ApiClient {
  private getBaseUrl(): string {
    const config = useRuntimeConfig()
    return config.public.apiBaseUrl || 'http://localhost:8080/api'
  }

  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const baseUrl = this.getBaseUrl()
    const url = endpoint.startsWith('http') ? endpoint : `${baseUrl}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`

    const headers = new Headers(options.headers || {})
    
    // If not sending FormData, default to application/json
    if (!(options.body instanceof FormData) && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }
    headers.set('Accept', 'application/json')

    // Attach Bearer token from localStorage or cookie if present
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('lms_token') || localStorage.getItem('token')
      if (token && !headers.has('Authorization')) {
        headers.set('Authorization', `Bearer ${token}`)
      }
    }

    const config: RequestInit = {
      ...options,
      headers,
      credentials: 'include' // Support HttpOnly cookies
    }

    try {
      const response = await fetch(url, config)

      if (response.status === 204) {
        return {} as T
      }

      const isJson = response.headers.get('content-type')?.includes('application/json')
      const data = isJson ? await response.json() : await response.text()

      if (!response.ok) {
        const errorMsg = (isJson && (data?.error || data?.message)) || response.statusText || 'Terjadi kesalahan pada server'
        const errorObj: ApiError = {
          message: errorMsg,
          status: response.status,
          errors: isJson && data?.errors ? data.errors : undefined
        }
        throw errorObj
      }

      // If response is wrapped in standard { success: true, data: ... }
      if (data && typeof data === 'object' && 'data' in data && 'success' in data) {
        // Only persist JWT auth token for explicit auth endpoints (prevent LiveKit room tokens from overwriting user session)
        if (endpoint.includes('/auth/') && data.data && typeof data.data === 'object' && data.data.token) {
          if (typeof window !== 'undefined') {
            localStorage.setItem('lms_token', data.data.token)
            localStorage.setItem('token', data.data.token)
            document.cookie = `lms_token=${data.data.token}; path=/; max-age=86400`
          }
        }
        return normalizeData(data.data) as T
      }

      return normalizeData(data) as T
    } catch (err: any) {
      if (err?.status) {
        throw err
      }

      // Network / Connection Error
      const networkError: ApiError = {
        message: 'Gagal terhubung ke server LMS. Pastikan server backend (port 8080) sedang berjalan.',
        status: 0
      }
      throw networkError
    }
  }

  get<T>(endpoint: string, queryParams?: Record<string, any>, options: RequestInit = {}): Promise<T> {
    let url = endpoint
    if (queryParams) {
      const filteredParams = Object.entries(queryParams).filter(([_, v]) => v !== undefined && v !== null && v !== '')
      if (filteredParams.length > 0) {
        const search = new URLSearchParams(filteredParams.map(([k, v]) => [k, String(v)])).toString()
        url += `${url.includes('?') ? '&' : '?'}${search}`
      }
    }
    return this.request<T>(url, { ...options, method: 'GET' })
  }

  post<T>(endpoint: string, body?: any, options: RequestInit = {}): Promise<T> {
    const isFormData = body instanceof FormData
    return this.request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: isFormData ? body : (body ? JSON.stringify(body) : undefined)
    })
  }

  put<T>(endpoint: string, body?: any, options: RequestInit = {}): Promise<T> {
    const isFormData = body instanceof FormData
    return this.request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: isFormData ? body : (body ? JSON.stringify(body) : undefined)
    })
  }

  patch<T>(endpoint: string, body?: any, options: RequestInit = {}): Promise<T> {
    const isFormData = body instanceof FormData
    return this.request<T>(endpoint, {
      ...options,
      method: 'PATCH',
      body: isFormData ? body : (body ? JSON.stringify(body) : undefined)
    })
  }

  delete<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: 'DELETE' })
  }
}

export const api = new ApiClient()

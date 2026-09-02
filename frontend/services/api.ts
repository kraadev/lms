import type { ApiError } from '~/types'

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
        const errorObj: ApiError = {
          message: (isJson && data?.message) || response.statusText || 'Terjadi kesalahan pada server',
          status: response.status,
          errors: isJson && data?.errors ? data.errors : undefined
        }
        throw errorObj
      }

      return data as T
    } catch (err: any) {
      if (err?.status) {
        throw err
      }

      // Network / Connection Error
      const networkError: ApiError = {
        message: 'Gagal terhubung ke server LMS. Periksa koneksi atau pastikan server backend berjalan.',
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

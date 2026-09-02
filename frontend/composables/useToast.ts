export interface Toast {
  id: string
  title?: string
  message: string
  type: 'success' | 'error' | 'warning' | 'info'
  duration?: number
}

export function useToast() {
  const toasts = useState<Toast[]>('app_toasts', () => [])

  function add(toast: Omit<Toast, 'id'>) {
    const id = Math.random().toString(36).substring(2, 9)
    const duration = toast.duration ?? 4000
    const newToast: Toast = { ...toast, id, duration }

    toasts.value.push(newToast)

    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }

    return id
  }

  function remove(id: string) {
    const idx = toasts.value.findIndex(t => t.id === id)
    if (idx !== -1) {
      toasts.value.splice(idx, 1)
    }
  }

  function success(message: string, title?: string) {
    return add({ message, title, type: 'success' })
  }

  function error(message: string, title?: string) {
    return add({ message, title, type: 'error' })
  }

  function warning(message: string, title?: string) {
    return add({ message, title, type: 'warning' })
  }

  function info(message: string, title?: string) {
    return add({ message, title, type: 'info' })
  }

  return {
    toasts,
    add,
    remove,
    success,
    error,
    warning,
    info
  }
}

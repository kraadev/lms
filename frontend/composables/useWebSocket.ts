import type { WebSocketEvent } from '~/types'

export type WsStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting'

type EventCallback<T = any> = (payload: T) => void

// Global module-level subscriber registry (singleton across all components)
const listeners = new Map<string, Set<EventCallback>>()

export function useWebSocket() {
  const config = useRuntimeConfig()
  const wsUrl = config.public.wsUrl || 'ws://localhost:8080/ws'

  const socket = useState<WebSocket | null>('ws_socket', () => null)
  const status = useState<WsStatus>('ws_status', () => 'disconnected')
  const retryCount = useState<number>('ws_retry_count', () => 0)
  const activeClassId = useState<number | string | null>('ws_active_class_id', () => null)

  let reconnectTimer: any = null
  const MAX_RETRY_DELAY = 15000
  const BASE_DELAY = 1000

  function connect() {
    if (import.meta.server) return
    if (socket.value && (socket.value.readyState === WebSocket.OPEN || socket.value.readyState === WebSocket.CONNECTING)) {
      return
    }

    status.value = retryCount.value > 0 ? 'reconnecting' : 'connecting'

    try {
      const token = typeof window !== 'undefined' ? (localStorage.getItem('lms_token') || localStorage.getItem('token') || '') : ''
      const urlWithToken = token ? `${wsUrl}${wsUrl.includes('?') ? '&' : '?'}token=${encodeURIComponent(token)}` : wsUrl

      const ws = new WebSocket(urlWithToken)

      ws.onopen = () => {
        status.value = 'connected'
        retryCount.value = 0
        if (reconnectTimer) {
          clearTimeout(reconnectTimer)
          reconnectTimer = null
        }

        // Auto re-join active class room on open
        if (activeClassId.value) {
          ws.send(JSON.stringify({
            type: 'class.join',
            payload: { class_id: Number(activeClassId.value) }
          }))
        }
      }

      ws.onmessage = (event: MessageEvent) => {
        try {
          const data: WebSocketEvent = JSON.parse(event.data)
          if (data && data.type) {
            const callbacks = listeners.get(data.type)
            if (callbacks) {
              callbacks.forEach(cb => {
                try {
                  cb(data.payload)
                } catch (err) {
                  console.error(`Error in WS event listener for ${data.type}:`, err)
                }
              })
            }
            // Also notify wildcard listeners
            const allCallbacks = listeners.get('*')
            if (allCallbacks) {
              allCallbacks.forEach(cb => {
                try {
                  cb(data)
                } catch (err) {
                  console.error('Error in wildcard WS listener:', err)
                }
              })
            }
          }
        } catch (e) {
          console.warn('Non-JSON WebSocket message received:', event.data)
        }
      }

      ws.onerror = (err) => {
        console.warn('WebSocket error:', err)
      }

      ws.onclose = (event) => {
        status.value = 'disconnected'
        socket.value = null

        // Auto-reconnect with exponential backoff if closed unexpectedly
        if (event.code !== 1000) {
          const delay = Math.min(BASE_DELAY * Math.pow(1.5, retryCount.value), MAX_RETRY_DELAY)
          retryCount.value += 1
          status.value = 'reconnecting'
          reconnectTimer = setTimeout(() => {
            connect()
          }, delay)
        }
      }

      socket.value = ws
    } catch (err) {
      status.value = 'disconnected'
      console.error('Failed to create WebSocket:', err)
    }
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (socket.value) {
      socket.value.close(1000, 'User logged out or client disconnected')
      socket.value = null
      status.value = 'disconnected'
    }
  }

  function send<T = any>(type: string, payload: T): boolean {
    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
      // Connect if not connected yet
      connect()
      return false
    }

    const message: WebSocketEvent<T> = { type, payload }
    socket.value.send(JSON.stringify(message))
    return true
  }

  function on<T = any>(type: string, callback: EventCallback<T>): () => void {
    if (!listeners.has(type)) {
      listeners.set(type, new Set())
    }
    listeners.get(type)!.add(callback)

    // Return unsubscribe function
    return () => {
      const callbacks = listeners.get(type)
      if (callbacks) {
        callbacks.delete(callback)
        if (callbacks.size === 0) {
          listeners.delete(type)
        }
      }
    }
  }

  // Helper to join a class chat room
  function joinClassRoom(classId: number | string) {
    activeClassId.value = classId
    connect()
    return send('class.join', { class_id: Number(classId) })
  }

  // Helper to leave a class chat room
  function leaveClassRoom(classId: number | string) {
    if (activeClassId.value === classId) {
      activeClassId.value = null
    }
    return send('class.leave', { class_id: Number(classId) })
  }

  // Helper to send a chat message
  function sendChatMessage(classId: number | string, message: string) {
    return send('chat.send', {
      class_id: Number(classId),
      message
    })
  }

  return {
    socket,
    status,
    retryCount,
    connect,
    disconnect,
    send,
    on,
    joinClassRoom,
    leaveClassRoom,
    sendChatMessage
  }
}

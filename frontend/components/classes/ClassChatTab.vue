<script setup lang="ts">
import { Send, MessageSquare, AlertCircle, RefreshCw, Clock } from 'lucide-vue-next'
import { messagesService } from '~/services/messages'
import type { ChatMessage } from '~/types'
import { formatDate, formatRelativeTime } from '~/utils/formatters'

const props = defineProps<{ classId: string }>()
const auth = useAuthStore()
const { status: wsStatus, on: onWs, joinClassRoom, leaveClassRoom, sendChatMessage } = useWebSocket()

const messages = ref<ChatMessage[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)
const inputText = ref('')
const isSending = ref(false)
const chatContainer = ref<HTMLDivElement | null>(null)

let unsubWs: (() => void) | null = null

function scrollToBottom(smooth = false) {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTo({
        top: chatContainer.value.scrollHeight,
        behavior: smooth ? 'smooth' : 'auto'
      })
    }
  })
}

async function loadHistory() {
  if (!props.classId) return
  isLoading.value = true
  error.value = null
  try {
    const res: any = await messagesService.getByClass(props.classId)
    messages.value = Array.isArray(res) ? res : (res?.messages || [])
    scrollToBottom(false)
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat riwayat obrolan'
    messages.value = []
  } finally {
    isLoading.value = false
  }
}

function handleIncomingMessage(payload: any) {
  if (!payload) return
  
  // Verify target class matches
  const targetClassId = payload.class_id ?? payload.classId
  if (String(targetClassId) === String(props.classId)) {
    const msgId = payload.id
    const exists = messages.value.some(m => String(m.id) === String(msgId))
    if (!exists) {
      // Normalize incoming structure
      const formatted: ChatMessage = {
        id: msgId,
        class_id: Number(props.classId),
        user_id: payload.user?.id ?? payload.user_id ?? payload.sender_id,
        user_name: payload.user?.name ?? payload.user_name ?? payload.sender?.name,
        user: payload.user ?? payload.sender,
        message: payload.message ?? payload.content ?? '',
        created_at: payload.created_at ?? new Date().toISOString()
      }
      messages.value = [...messages.value, formatted]
      scrollToBottom(true)
    }
  }
}

function setupRoom() {
  loadHistory()
  joinClassRoom(props.classId)
  if (unsubWs) unsubWs()
  unsubWs = onWs('chat.message', handleIncomingMessage)
}

watch(() => props.classId, (newId, oldId) => {
  if (oldId) leaveClassRoom(oldId)
  if (newId) setupRoom()
})

onMounted(() => {
  setupRoom()
})

onBeforeUnmount(() => {
  leaveClassRoom(props.classId)
  if (unsubWs) unsubWs()
})

async function submitMessage() {
  const text = inputText.value.trim()
  if (!text || isSending.value) return

  isSending.value = true
  try {
    // 1. Send through WebSocket
    const wsSent = sendChatMessage(props.classId, text)
    if (wsSent) {
      inputText.value = ''
      scrollToBottom(true)
    } else {
      // 2. Fallback to HTTP REST API if WebSocket is offline
      const newMsg: any = await messagesService.send(props.classId, text)
      inputText.value = ''
      handleIncomingMessage(newMsg)
    }
  } catch (e: any) {
    console.error('Failed to send message:', e)
  } finally {
    isSending.value = false
  }
}

function isMyMessage(msg: ChatMessage) {
  const senderId = msg.user_id ?? msg.user?.id ?? msg.sender_id ?? msg.sender?.id
  return senderId === auth.user?.id
}

function getSenderName(msg: ChatMessage) {
  return msg.user_name || msg.user?.name || msg.sender?.name || 'Anggota Kelas'
}
</script>

<template>
  <div class="p-4 md:p-6 flex flex-col h-[calc(100vh-14rem)] min-h-[500px]">
    <!-- WS Connection Status Indicator -->
    <div class="flex items-center justify-between pb-3 border-b border-surface-200 dark:border-surface-800 text-xs">
      <div class="flex items-center gap-2">
        <span
          class="w-2 h-2 rounded-full"
          :class="{
            'bg-emerald-500': wsStatus === 'connected',
            'bg-amber-500 animate-pulse': wsStatus === 'connecting' || wsStatus === 'reconnecting',
            'bg-rose-500': wsStatus === 'disconnected'
          }"
        />
        <span class="text-surface-500 dark:text-surface-400">
          {{ wsStatus === 'connected' ? 'Terhubung realtime (WebSocket)' : wsStatus === 'reconnecting' ? 'Menyambung kembali...' : 'Mode Offline (REST)' }}
        </span>
      </div>
      <button
        type="button"
        class="text-surface-400 hover:text-surface-600 dark:hover:text-surface-300 flex items-center gap-1 cursor-pointer transition-colors"
        @click="loadHistory"
      >
        <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
        <span>Segarkan</span>
      </button>
    </div>

    <!-- Messages list -->
    <div
      ref="chatContainer"
      class="flex-1 overflow-y-auto py-4 space-y-3 pr-2 scrollbar-thin"
    >
      <div v-if="isLoading" class="space-y-3">
        <UiSkeleton :rows="5" />
      </div>

      <div v-else-if="error" class="text-center py-8">
        <p class="text-sm text-rose-500">{{ error }}</p>
        <UiButton variant="outline" size="sm" class="mt-2" @click="loadHistory">
          Coba Lagi
        </UiButton>
      </div>

      <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center h-full text-center py-12">
        <div class="w-12 h-12 rounded-2xl bg-surface-100 dark:bg-surface-800 flex items-center justify-center text-surface-400 mb-3">
          <MessageSquare class="w-6 h-6" />
        </div>
        <p class="text-sm font-semibold text-surface-700 dark:text-surface-300">Belum ada obrolan</p>
        <p class="text-xs text-surface-400 max-w-xs mt-1">
          Mulai diskusi dengan guru dan teman sekelas Anda di ruangan ini!
        </p>
      </div>

      <!-- Chat Bubbles -->
      <template v-else>
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="[
            'flex flex-col max-w-[80%] sm:max-w-[70%]',
            isMyMessage(msg) ? 'ml-auto items-end' : 'mr-auto items-start'
          ]"
        >
          <!-- Sender name for other users -->
          <span
            v-if="!isMyMessage(msg)"
            class="text-[11px] font-medium text-surface-500 dark:text-surface-400 mb-1 px-1"
          >
            {{ getSenderName(msg) }}
          </span>

          <!-- Bubble content -->
          <div
            :class="[
              'px-4 py-2.5 rounded-2xl text-sm leading-relaxed break-words shadow-sm',
              isMyMessage(msg)
                ? 'bg-brand-600 text-white rounded-br-xs'
                : 'bg-surface-100 dark:bg-surface-800 text-surface-900 dark:text-surface-100 rounded-bl-xs'
            ]"
          >
            {{ msg.message || (msg as any).content }}
          </div>

          <!-- Timestamp -->
          <span class="text-[10px] text-surface-400 mt-1 px-1 flex items-center gap-1">
            <Clock class="w-2.5 h-2.5" />
            {{ formatRelativeTime(msg.created_at) }}
          </span>
        </div>
      </template>
    </div>

    <!-- Input Form -->
    <div class="pt-3 border-t border-surface-200 dark:border-surface-800">
      <form class="flex items-center gap-2" @submit.prevent="submitMessage">
        <input
          v-model="inputText"
          type="text"
          placeholder="Tulis pesan ke kelas..."
          class="flex-1 px-4 py-2.5 text-sm rounded-xl border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100 focus:outline-none focus:ring-2 focus:ring-brand-500 transition-all placeholder:text-surface-400"
          :disabled="isSending"
          @keydown.enter.prevent="submitMessage"
        />
        <UiButton
          type="submit"
          variant="primary"
          size="md"
          :disabled="!inputText.trim() || isSending"
          :loading="isSending"
        >
          <Send class="w-4 h-4" />
          <span class="hidden sm:inline">Kirim</span>
        </UiButton>
      </form>
    </div>
  </div>
</template>

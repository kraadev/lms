<script setup lang="ts">
import { Send, MessageSquare, AlertCircle, RefreshCw } from 'lucide-vue-next'
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
  isLoading.value = true
  error.value = null
  try {
    const res = await messagesService.getByClass(props.classId)
    messages.value = res.messages || []
    scrollToBottom(false)
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat riwayat obrolan'
  } finally {
    isLoading.value = false
  }
}

function handleIncomingMessage(payload: any) {
  if (payload && (payload.class_id === Number(props.classId) || payload.class_id === props.classId)) {
    // Avoid duplicate if we receive our own message echo
    const exists = messages.value.some(m => m.id === payload.id)
    if (!exists) {
      messages.value.push(payload)
      scrollToBottom(true)
    }
  }
}

onMounted(() => {
  loadHistory()
  joinClassRoom(props.classId)
  unsubWs = onWs('chat.message', handleIncomingMessage)
})

onBeforeUnmount(() => {
  leaveClassRoom(props.classId)
  if (unsubWs) unsubWs()
})

async function submitMessage() {
  const text = inputText.value.trim()
  if (!text || isSending.value) return

  isSending.value = true
  const ok = sendChatMessage(props.classId, text)
  if (ok) {
    inputText.value = ''
    scrollToBottom(true)
  } else {
    // WebSocket not open fallback or notification
  }
  isSending.value = false
}

function isMyMessage(msg: ChatMessage) {
  return msg.user_id === auth.user?.id || msg.user?.id === auth.user?.id
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
          {{ wsStatus === 'connected' ? 'Terhubung realtime' : wsStatus === 'reconnecting' ? 'Menyambung kembali...' : 'Offline' }}
        </span>
      </div>
      <button
        type="button"
        class="text-surface-400 hover:text-surface-600 dark:hover:text-surface-300 flex items-center gap-1"
        @click="loadHistory"
      >
        <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
        Segarkan
      </button>
    </div>

    <!-- Messages list -->
    <div
      ref="chatContainer"
      class="flex-1 overflow-y-auto py-4 space-y-3 pr-2 scrollbar-thin"
    >
      <div v-if="isLoading">
        <UiSkeleton :rows="5" />
      </div>

      <UiErrorState v-else-if="error" :message="error" @retry="loadHistory" />

      <div v-else-if="!messages.length" class="h-full flex flex-col items-center justify-center text-center p-6 text-surface-400">
        <MessageSquare class="w-10 h-10 mb-2 stroke-1" />
        <p class="text-sm font-medium text-surface-700 dark:text-surface-300">Belum ada pesan di kelas ini</p>
        <p class="text-xs text-surface-400 mt-0.5">Kirim pesan pertama untuk memulai percakapan kelas.</p>
      </div>

      <template v-else>
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="flex items-end gap-2.5 max-w-[85%]"
          :class="isMyMessage(msg) ? 'ml-auto flex-row-reverse' : 'mr-auto'"
        >
          <UiAvatar
            v-if="!isMyMessage(msg)"
            :name="msg.user?.name || 'User'"
            :src="msg.user?.avatar"
            size="xs"
            class="shrink-0 mb-1"
          />

          <div class="flex flex-col" :class="isMyMessage(msg) ? 'items-end' : 'items-start'">
            <span v-if="!isMyMessage(msg)" class="text-[11px] font-medium text-surface-500 dark:text-surface-400 mb-1 px-1">
              {{ msg.user?.name }}
            </span>

            <div
              class="px-4 py-2.5 rounded-2xl text-sm leading-relaxed"
              :class="isMyMessage(msg)
                ? 'bg-brand-600 text-white rounded-br-none'
                : 'bg-surface-100 dark:bg-surface-800 text-surface-900 dark:text-surface-100 rounded-bl-none border border-surface-200/50 dark:border-surface-700/50'"
            >
              <p class="whitespace-pre-wrap break-words">{{ msg.message || msg.content }}</p>
            </div>

            <span class="text-[10px] text-surface-400 mt-1 px-1">
              {{ formatDate(msg.created_at, { hour: '2-digit', minute: '2-digit' }) }}
            </span>
          </div>
        </div>
      </template>
    </div>

    <!-- Message Input Bar -->
    <div class="pt-3 border-t border-surface-200 dark:border-surface-800">
      <form @submit.prevent="submitMessage" class="flex items-center gap-2">
        <input
          v-model="inputText"
          type="text"
          placeholder="Tulis pesan..."
          class="flex-1 px-4 py-2.5 bg-surface-50 dark:bg-surface-800/80 border border-surface-200 dark:border-surface-700 rounded-xl text-sm text-surface-900 dark:text-surface-100 placeholder-surface-400 focus:outline-none focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-colors"
          :disabled="isSending"
          @keydown.enter.exact.prevent="submitMessage"
        />
        <UiButton
          type="submit"
          :disabled="!inputText.trim() || isSending"
          class="shrink-0 px-4"
        >
          <Send class="w-4 h-4" />
        </UiButton>
      </form>
    </div>
  </div>
</template>

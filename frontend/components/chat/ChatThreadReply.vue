<script setup lang="ts">
import { CornerDownRight, Smile } from 'lucide-vue-next'

interface Props {
  replyToName?: string
  replyToMessage?: string
  reactions?: Record<string, number>
}

withDefaults(defineProps<Props>(), {
  replyToName: '',
  replyToMessage: '',
  reactions: () => ({})
})

const emit = defineEmits<{
  (e: 'react', emoji: string): void
  (e: 'cancel-reply'): void
}>()

const availableEmojis = ['👍', '❤️', '👏', '🔥', '💡']
</script>

<template>
  <div class="space-y-1.5">
    <!-- Active reply preview banner -->
    <div v-if="replyToName" class="flex items-center justify-between px-3 py-1.5 rounded-lg bg-surface-100 dark:bg-surface-800 text-xs border-l-2 border-brand-500">
      <div class="flex items-center gap-2 truncate">
        <CornerDownRight class="w-3.5 h-3.5 text-brand-600 shrink-0" />
        <span class="font-semibold text-surface-800 dark:text-surface-200">{{ replyToName }}:</span>
        <span class="text-surface-500 dark:text-surface-400 truncate">{{ replyToMessage }}</span>
      </div>
      <button type="button" class="text-surface-400 hover:text-surface-600 text-xs ml-2" @click="emit('cancel-reply')">&times;</button>
    </div>

    <!-- Quick Emoji Reaction Bar -->
    <div class="flex items-center gap-1">
      <button
        v-for="emoji in availableEmojis"
        :key="emoji"
        type="button"
        class="px-2 py-0.5 rounded-full text-xs hover:bg-surface-100 dark:hover:bg-surface-800 transition-colors"
        @click="emit('react', emoji)"
      >
        {{ emoji }}
        <span v-if="reactions[emoji]" class="ml-1 text-[10px] font-bold text-surface-600 dark:text-surface-400">{{ reactions[emoji] }}</span>
      </button>
    </div>
  </div>
</template>

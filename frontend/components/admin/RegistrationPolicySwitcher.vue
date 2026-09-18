<script setup lang="ts">
import { ref } from 'vue'
import { Globe, Mail, ShieldAlert, Check } from 'lucide-vue-next'

export type PolicyMode = 'public' | 'invite_only' | 'admin_approval'

interface Props {
  modelValue?: PolicyMode
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: 'public'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: PolicyMode): void
}>()

const policies = [
  {
    key: 'public' as PolicyMode,
    title: 'Pendaftaran Terbuka',
    desc: 'Siapa saja dapat membuat akun dan langsung bergabung ke platform.',
    icon: Globe
  },
  {
    key: 'invite_only' as PolicyMode,
    title: 'Khusus Undangan',
    desc: 'Pendaftaran membutuhkan tautan atau kode undangan khusus dari pengajar/admin.',
    icon: Mail
  },
  {
    key: 'admin_approval' as PolicyMode,
    title: 'Persetujuan Admin',
    desc: 'Akun baru harus ditinjau dan disetujui administrator sebelum dapat login.',
    icon: ShieldAlert
  }
]
</script>

<template>
  <div class="space-y-3">
    <div
      v-for="p in policies"
      :key="p.key"
      :class="[
        'flex items-start gap-4 p-4 rounded-xl border transition-all cursor-pointer select-none',
        modelValue === p.key
          ? 'border-brand-500 bg-brand-50/50 dark:bg-brand-950/30 shadow-soft'
          : 'border-surface-200 dark:border-surface-800 hover:border-surface-300 dark:hover:border-surface-700 bg-white dark:bg-surface-900'
      ]"
      @click="emit('update:modelValue', p.key)"
    >
      <div
        :class="[
          'w-9 h-9 rounded-xl flex items-center justify-center shrink-0 mt-0.5',
          modelValue === p.key ? 'bg-brand-600 text-white' : 'bg-surface-100 dark:bg-surface-800 text-surface-500'
        ]"
      >
        <component :is="p.icon" class="w-4.5 h-4.5" />
      </div>

      <div class="flex-1 min-w-0">
        <p class="text-sm font-semibold text-surface-900 dark:text-surface-100">{{ p.title }}</p>
        <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5 leading-relaxed">{{ p.desc }}</p>
      </div>

      <div v-if="modelValue === p.key" class="w-5 h-5 rounded-full bg-brand-600 text-white flex items-center justify-center shrink-0 mt-1">
        <Check class="w-3 h-3 stroke-[3]" />
      </div>
    </div>
  </div>
</template>

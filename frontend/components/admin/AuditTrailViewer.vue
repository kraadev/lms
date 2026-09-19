<script setup lang="ts">
import { ref, computed } from 'vue'
import { ShieldAlert, Search, Filter } from 'lucide-vue-next'

export interface AuditRecord {
  id: number
  actorEmail: string
  action: string
  method: string
  path: string
  statusCode: number
  durationMs: number
  timestamp: string
}

interface Props {
  logs: AuditRecord[]
}

const props = withDefaults(defineProps<Props>(), {
  logs: () => []
})

const searchQuery = ref('')
const selectedMethod = ref<string>('all')

const filteredLogs = computed(() => {
  return props.logs.filter((item) => {
    const matchQuery =
      item.actorEmail.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      item.path.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      item.action.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchMethod = selectedMethod.value === 'all' || item.method === selectedMethod.value
    return matchQuery && matchMethod
  })
})

function methodBadgeVariant(method: string) {
  switch (method) {
    case 'POST': return 'primary'
    case 'DELETE': return 'danger'
    case 'PUT':
    case 'PATCH': return 'warning'
    default: return 'default'
  }
}
</script>

<template>
  <div class="space-y-4">
    <!-- Filter bar -->
    <div class="flex flex-col sm:flex-row gap-3 items-center justify-between">
      <div class="relative w-full sm:w-64">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari aktor, path, aksi..."
          class="w-full pl-9 pr-3 py-2 rounded-xl text-xs border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 text-surface-900 dark:text-surface-100 focus:outline-none"
        />
        <Search class="w-4 h-4 text-surface-400 absolute left-3 top-2.5" />
      </div>

      <div class="flex items-center gap-2 self-end sm:self-auto">
        <Filter class="w-4 h-4 text-surface-400" />
        <select
          v-model="selectedMethod"
          class="px-2.5 py-1.5 rounded-xl text-xs border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 text-surface-700 dark:text-surface-300 focus:outline-none cursor-pointer"
        >
          <option value="all">Semua Method</option>
          <option value="POST">POST</option>
          <option value="PUT">PUT</option>
          <option value="PATCH">PATCH</option>
          <option value="DELETE">DELETE</option>
        </select>
      </div>
    </div>

    <!-- Table Container -->
    <div class="rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft">
      <div v-if="!filteredLogs.length" class="p-8 text-center text-xs text-surface-400">
        Tidak ada log audit yang sesuai kriteria.
      </div>
      <table v-else class="w-full text-left text-xs">
        <thead class="bg-surface-50 dark:bg-surface-800/50 border-b border-surface-100 dark:border-surface-800 text-surface-500">
          <tr>
            <th class="px-4 py-3 font-semibold">Method</th>
            <th class="px-4 py-3 font-semibold">Aktor</th>
            <th class="px-4 py-3 font-semibold">Endpoint Path</th>
            <th class="px-4 py-3 font-semibold">Status</th>
            <th class="px-4 py-3 font-semibold text-right">Durasi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-surface-100 dark:divide-surface-800 text-surface-700 dark:text-surface-300">
          <tr v-for="l in filteredLogs" :key="l.id" class="hover:bg-surface-50/60 dark:hover:bg-surface-800/30">
            <td class="px-4 py-3 font-mono font-bold">
              <UiBadge :variant="methodBadgeVariant(l.method)" size="sm">{{ l.method }}</UiBadge>
            </td>
            <td class="px-4 py-3 truncate max-w-xs">{{ l.actorEmail }}</td>
            <td class="px-4 py-3 font-mono text-[11px] truncate max-w-xs">{{ l.path }}</td>
            <td class="px-4 py-3">
              <span :class="l.statusCode < 400 ? 'text-emerald-600 font-semibold' : 'text-rose-500 font-semibold'">
                {{ l.statusCode }}
              </span>
            </td>
            <td class="px-4 py-3 text-right text-surface-400">{{ l.durationMs }}ms</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, Key, ShieldCheck, UserMinus } from 'lucide-vue-next'

export interface UserRow {
  id: number
  name: string
  email: string
  role: 'admin' | 'teacher' | 'student'
  status: 'active' | 'suspended'
}

interface Props {
  users?: UserRow[]
}

const props = withDefaults(defineProps<Props>(), {
  users: () => []
})

const emit = defineEmits<{
  (e: 'reset-password', user: UserRow): void
  (e: 'toggle-suspend', user: UserRow): void
}>()

const search = ref('')
const roleFilter = ref('all')

const filteredUsers = computed(() => {
  return props.users.filter((u) => {
    const matchSearch = u.name.toLowerCase().includes(search.value.toLowerCase()) ||
                        u.email.toLowerCase().includes(search.value.toLowerCase())
    const matchRole = roleFilter.value === 'all' || u.role === roleFilter.value
    return matchSearch && matchRole
  })
})
</script>

<template>
  <div class="space-y-4">
    <!-- Search & Filter Controls -->
    <div class="flex flex-col sm:flex-row items-center justify-between gap-3">
      <div class="relative w-full sm:w-72">
        <input
          v-model="search"
          type="text"
          placeholder="Cari nama atau email..."
          class="w-full pl-9 pr-3 py-2 rounded-xl text-xs border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 text-surface-900 dark:text-surface-100 focus:outline-none"
        />
        <Search class="w-4 h-4 text-surface-400 absolute left-3 top-2.5" />
      </div>

      <div class="flex items-center gap-1.5 self-end sm:self-auto">
        <button
          v-for="r in ['all', 'admin', 'teacher', 'student']"
          :key="r"
          type="button"
          :class="[
            'px-3 py-1.5 rounded-lg text-xs font-semibold capitalize transition-colors',
            roleFilter === r ? 'bg-brand-600 text-white' : 'bg-surface-100 dark:bg-surface-800 text-surface-600 dark:text-surface-300'
          ]"
          @click="roleFilter = r"
        >
          {{ r === 'all' ? 'Semua' : r }}
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="rounded-xl border border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900 overflow-hidden shadow-soft">
      <table class="w-full text-left text-xs">
        <thead class="bg-surface-50 dark:bg-surface-800/50 border-b border-surface-100 dark:border-surface-800 text-surface-500">
          <tr>
            <th class="px-4 py-3 font-semibold">Pengguna</th>
            <th class="px-4 py-3 font-semibold">Peran</th>
            <th class="px-4 py-3 font-semibold">Status</th>
            <th class="px-4 py-3 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-surface-100 dark:divide-surface-800 text-surface-700 dark:text-surface-300">
          <tr v-for="u in filteredUsers" :key="u.id" class="hover:bg-surface-50/60 dark:hover:bg-surface-800/30">
            <td class="px-4 py-3 flex items-center gap-3">
              <UiAvatar :name="u.name" size="sm" />
              <div>
                <p class="font-semibold text-surface-900 dark:text-surface-100">{{ u.name }}</p>
                <p class="text-[11px] text-surface-400">{{ u.email }}</p>
              </div>
            </td>
            <td class="px-4 py-3">
              <UiBadge :variant="u.role === 'admin' ? 'danger' : u.role === 'teacher' ? 'primary' : 'default'" size="sm">
                {{ u.role }}
              </UiBadge>
            </td>
            <td class="px-4 py-3">
              <span :class="u.status === 'active' ? 'text-emerald-600 font-semibold' : 'text-rose-500 font-semibold'">
                {{ u.status === 'active' ? 'Aktif' : 'Suspended' }}
              </span>
            </td>
            <td class="px-4 py-3 text-right">
              <div class="flex items-center justify-end gap-1.5">
                <button type="button" class="p-1.5 rounded-lg hover:bg-surface-100 dark:hover:bg-surface-800 text-surface-500" title="Reset Password" @click="emit('reset-password', u)">
                  <Key class="w-3.5 h-3.5" />
                </button>
                <button type="button" class="p-1.5 rounded-lg hover:bg-surface-100 dark:hover:bg-surface-800 text-rose-500" title="Suspend User" @click="emit('toggle-suspend', u)">
                  <UserMinus class="w-3.5 h-3.5" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

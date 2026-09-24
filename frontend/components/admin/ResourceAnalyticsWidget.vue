<script setup lang="ts">
import { Cpu, HardDrive, Activity, Radio } from 'lucide-vue-next'

interface Props {
  cpuPercent?: number
  ramPercent?: number
  storageUsedGb?: number
  storageMaxGb?: number
  activeSockets?: number
}

withDefaults(defineProps<Props>(), {
  cpuPercent: 24,
  ramPercent: 42,
  storageUsedGb: 12.4,
  storageMaxGb: 50.0,
  activeSockets: 8
})
</script>

<template>
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
    <!-- CPU Card -->
    <UiCard variant="default" padding="md">
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs font-semibold text-surface-500">CPU Usage</span>
        <Cpu class="w-4 h-4 text-brand-600" />
      </div>
      <p class="text-2xl font-bold text-surface-900 dark:text-surface-100">{{ cpuPercent }}%</p>
      <UiProgress :value="cpuPercent" :max="100" size="xs" variant="brand" class="mt-2" />
    </UiCard>

    <!-- RAM Card -->
    <UiCard variant="default" padding="md">
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs font-semibold text-surface-500">Memory</span>
        <Activity class="w-4 h-4 text-emerald-600" />
      </div>
      <p class="text-2xl font-bold text-surface-900 dark:text-surface-100">{{ ramPercent }}%</p>
      <UiProgress :value="ramPercent" :max="100" size="xs" variant="success" class="mt-2" />
    </UiCard>

    <!-- Storage Card -->
    <UiCard variant="default" padding="md">
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs font-semibold text-surface-500">Storage</span>
        <HardDrive class="w-4 h-4 text-amber-600" />
      </div>
      <p class="text-2xl font-bold text-surface-900 dark:text-surface-100">{{ storageUsedGb }} <span class="text-xs font-normal text-surface-400">/ {{ storageMaxGb }} GB</span></p>
      <UiProgress :value="storageUsedGb" :max="storageMaxGb" size="xs" variant="warning" class="mt-2" />
    </UiCard>

    <!-- WebSocket Card -->
    <UiCard variant="default" padding="md">
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs font-semibold text-surface-500">WS Conns</span>
        <Radio class="w-4 h-4 text-violet-600" />
      </div>
      <p class="text-2xl font-bold text-surface-900 dark:text-surface-100">{{ activeSockets }}</p>
      <p class="text-[11px] text-surface-400 mt-2">Active WebSockets</p>
    </UiCard>
  </div>
</template>

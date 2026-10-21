<template>
  <span class="adm-status" :class="`adm-status--${meta.tone}`">
    <i />
    {{ meta.label }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PageStatus } from '@/types/admin'

const props = defineProps<{
  status?: PageStatus
}>()

const map: Record<PageStatus, { label: string; tone: string }> = {
  published: { label: '已发布', tone: 'success' },
  draft: { label: '未发布', tone: 'neutral' }
}

const meta = computed(() => map[props.status ?? 'draft'])
</script>

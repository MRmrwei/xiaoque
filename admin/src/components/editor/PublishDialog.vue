<template>
  <el-dialog
    :model-value="open"
    title="发布页面"
    width="600px"
    @update:model-value="emit('close')"
  >
    <div v-if="page" class="publish-body">
      <div class="publish-target">
        <div>
          <span class="mono publish-eyebrow">TARGET</span>
          <strong>{{ page.name }}</strong>
          <span class="mono">{{ fullPath }}</span>
        </div>
        <span class="lang-badge" :class="{ 'is-en': page.locale === 'en' }">
          {{ page.locale === 'en' ? 'EN' : '中' }}
        </span>
      </div>

      <ul class="publish-checks">
        <li
          v-for="check in checks"
          :key="check.key"
          :class="check.ok ? 'is-ok' : check.warn ? 'is-warn' : 'is-fail'"
        >
          <span class="check-icon">{{ check.ok ? '✓' : check.warn ? '!' : '×' }}</span>
          <span class="check-label">{{ check.label }}</span>
          <span class="check-state">{{ check.ok ? '通过' : check.warn ? '建议' : '需修复' }}</span>
        </li>
      </ul>

      <div class="local-note">
发布会用当前草稿覆盖线上内容，官网访客会看到最新发布的版本。
      </div>
    </div>

    <template #footer>
      <el-button @click="emit('close')">取消</el-button>
      <el-button type="primary" :disabled="!ready" @click="confirm">确认发布</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PageRecord } from '~/domain/types'
import type { ValidationCheck } from '@/composables/usePublishValidation'

const props = defineProps<{
  open: boolean
  page: PageRecord | null
  checks: ValidationCheck[]
  ready: boolean
}>()

const emit = defineEmits<{
  close: []
  confirm: []
}>()

const fullPath = computed(() => {
  if (!props.page) return ''
  const prefix = props.page.locale === 'en' ? '/en' : ''
  return props.page.path === '/' ? prefix || '/' : `${prefix}${props.page.path}`
})

function confirm() {
  if (!props.ready) return
  emit('confirm')
}
</script>

<style scoped>
.publish-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.publish-target {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 14px;
  border-radius: 9px;
  background: var(--adm-surface);
}

.publish-target > div {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.publish-eyebrow {
  font-size: 7.5px;
  letter-spacing: 0.12em;
  color: var(--adm-muted-2);
}

.publish-target strong {
  font-size: 13px;
}

.publish-target .mono {
  font-size: 9.5px;
  color: var(--adm-muted);
}

.publish-checks {
  display: flex;
  flex-direction: column;
  gap: 3px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.publish-checks li {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 8px 10px;
  border-radius: 7px;
  background: var(--adm-surface);
  font-size: 11.5px;
}

.check-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  font-size: 11px;
  font-weight: 700;
  flex: none;
}

.is-ok .check-icon {
  background: var(--adm-green-soft);
  color: var(--adm-green);
}

.is-warn .check-icon {
  background: var(--adm-amber-soft);
  color: var(--adm-amber);
}

.is-fail .check-icon {
  background: var(--adm-red-soft);
  color: #b23a25;
}

.check-label {
  flex: 1;
}

.check-state {
  font-size: 9px;
  color: var(--adm-muted-2);
}

.is-ok .check-state {
  color: var(--adm-green);
}

.is-warn .check-state {
  color: var(--adm-amber);
}

.is-fail .check-state {
  color: #b23a25;
}

.local-note {
  padding: 9px 11px;
  border-radius: 7px;
  background: var(--adm-amber-soft);
  color: var(--adm-amber);
  font-size: 10.5px;
  line-height: 1.55;
}
</style>

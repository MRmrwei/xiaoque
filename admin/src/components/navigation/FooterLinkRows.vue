<template>
  <div class="footer-link-editor">
    <div v-if="showHead" class="footer-link-head">
      <span />
      <span>中文文案</span>
      <span>English</span>
      <span>跳转路径</span>
      <span class="footer-link-head-actions">显示</span>
    </div>

    <VueDraggable
      :model-value="modelValue"
      :animation="150"
      handle=".footer-link-drag"
      class="footer-link-rows"
      @update:model-value="onReorder"
    >
      <div
        v-for="(link, index) in modelValue"
        :key="link.key"
        class="footer-link-row"
        :class="{ 'is-hidden': !link.visible }"
      >
        <span class="footer-link-drag" title="拖动排序">
          <UiIcon name="drag" :size="12" />
        </span>
        <el-input
          v-model="link.zh"
          size="small"
          class="footer-link-field footer-link-field--zh"
          placeholder="中文文案"
        />
        <el-input
          v-model="link.en"
          size="small"
          class="footer-link-field footer-link-field--en"
          placeholder="English label"
        />
        <div class="footer-link-path">
          <UiIcon name="link" :size="12" />
          <el-input v-model="link.path" size="small" placeholder="/path" />
        </div>
        <div class="footer-link-actions">
          <el-switch
            v-model="link.visible"
            size="small"
            :title="link.visible ? '显示中，点击隐藏' : '已隐藏，点击显示'"
          />
          <button class="footer-link-remove" title="删除链接" @click="remove(index)">
            <UiIcon name="trash" :size="12" />
          </button>
        </div>
      </div>
    </VueDraggable>

    <div v-if="!modelValue.length" class="footer-link-empty">
      还没有链接，点下方按钮添加。
    </div>
  </div>
</template>

<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import type { FooterLink } from '~/domain/types'
import UiIcon from '~/components/UiIcon.vue'

const props = withDefaults(
  defineProps<{
    modelValue: FooterLink[]
    showHead?: boolean
  }>(),
  { showHead: true }
)

const emit = defineEmits<{
  'update:modelValue': [value: FooterLink[]]
}>()

function onReorder(value: FooterLink[]) {
  emit('update:modelValue', value)
}

function remove(index: number) {
  const next = props.modelValue.slice()
  next.splice(index, 1)
  emit('update:modelValue', next)
}
</script>

<style scoped>
.footer-link-editor {
  border: 1px solid var(--adm-line-soft);
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
}

.footer-link-head {
  display: grid;
  grid-template-columns:
    30px minmax(0, 1.05fr) minmax(0, 1.05fr) minmax(0, 0.9fr) 76px;
  gap: 0;
  align-items: stretch;
  padding: 0;
  background: #f7f3ec;
  border-bottom: 1px solid var(--adm-line);
  font-size: 11.5px;
  font-weight: 600;
  color: var(--adm-ink-soft);
}

.footer-link-head > span {
  display: flex;
  align-items: center;
  min-width: 0;
  padding: 8px;
  border-left: 1px solid var(--adm-line-soft);
}

.footer-link-head > span:first-child {
  border-left: 0;
}

.footer-link-head-actions {
  justify-content: flex-end;
  padding-right: 10px;
}

.footer-link-rows {
  display: flex;
  flex-direction: column;
}

.footer-link-row {
  display: grid;
  grid-template-columns:
    30px minmax(0, 1.05fr) minmax(0, 1.05fr) minmax(0, 0.9fr) 76px;
  gap: 0;
  align-items: stretch;
  padding: 0;
  border-bottom: 1px solid var(--adm-line-soft);
  background: #fff;
  transition: background 0.15s ease;
}

.footer-link-row > * {
  display: flex;
  align-items: center;
  min-width: 0;
  padding: 4px 8px;
  border-left: 1px solid var(--adm-line-soft);
}

.footer-link-row > *:first-child {
  justify-content: center;
  padding: 4px 0;
  border-left: 0;
}

.footer-link-row:last-child {
  border-bottom: 0;
}

.footer-link-row:hover {
  background: #faf8f4;
}

.footer-link-row.is-hidden {
  opacity: 0.45;
}

.footer-link-drag {
  display: grid;
  place-items: center;
  height: 24px;
  color: var(--adm-muted-2);
  cursor: grab;
}

.footer-link-drag:active {
  cursor: grabbing;
}

.footer-link-field :deep(.el-input__wrapper) {
  padding: 2px 6px;
  background: transparent;
  box-shadow: none;
  border-radius: 5px;
  transition: background 0.15s ease, box-shadow 0.15s ease;
}

.footer-link-field :deep(.el-input__wrapper:hover) {
  background: rgba(20, 21, 24, 0.045);
}

.footer-link-field :deep(.el-input__wrapper.is-focus) {
  background: #fff;
  box-shadow: inset 0 0 0 1px var(--adm-line);
}

.footer-link-field--zh :deep(input) {
  color: var(--adm-ink);
  font-weight: 500;
}

.footer-link-field--en :deep(input) {
  color: var(--adm-muted);
}

.footer-link-path {
  display: flex;
  align-items: center;
  gap: 5px;
  min-width: 0;
  color: var(--adm-muted-2);
}

.footer-link-path .ui-icon {
  opacity: 0.5;
}

.footer-link-path :deep(.el-input) {
  flex: 1;
  min-width: 0;
}

.footer-link-path :deep(.el-input__wrapper) {
  padding: 2px 4px;
  background: transparent;
  box-shadow: none;
}

.footer-link-path :deep(.el-input__wrapper:hover) {
  background: rgba(20, 21, 24, 0.045);
}

.footer-link-path :deep(.el-input__wrapper.is-focus) {
  background: #fff;
  box-shadow: inset 0 0 0 1px var(--adm-line);
}

.footer-link-path :deep(input) {
  font-family: var(--adm-mono);
  font-size: 10px;
  color: var(--adm-muted);
}

.footer-link-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  justify-content: flex-end;
}

.footer-link-remove {
  width: 24px;
  height: 24px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 5px;
  background: transparent;
  color: var(--adm-muted-2);
  cursor: pointer;
  opacity: 0.35;
  transition: opacity 0.15s ease, background 0.15s ease, color 0.15s ease;
}

.footer-link-row:hover .footer-link-remove {
  opacity: 1;
}

.footer-link-remove:hover {
  background: var(--adm-red-soft);
  color: var(--adm-accent);
}

.footer-link-empty {
  margin: 0;
  padding: 14px;
  text-align: center;
  font-size: 10.5px;
  color: var(--adm-muted-2);
}

@media (max-width: 760px) {
  .footer-link-head {
    display: none;
  }
  .footer-link-row {
    grid-template-columns: 22px minmax(0, 1fr) auto;
    grid-template-areas:
      'drag zh actions'
      'drag en en'
      'drag path path';
    align-items: start;
    gap: 6px;
    padding: 8px;
  }
  .footer-link-row > * {
    padding: 0;
    border-left: 0;
  }
  .footer-link-drag {
    grid-area: drag;
    align-self: start;
  }
  .footer-link-field--zh {
    grid-area: zh;
  }
  .footer-link-field--en {
    grid-area: en;
  }
  .footer-link-path {
    grid-area: path;
  }
  .footer-link-actions {
    grid-area: actions;
    justify-content: flex-end;
  }
}
</style>

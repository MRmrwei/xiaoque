<template>
  <div class="editor-panel">
    <div class="panel-head">
      <div>
        <p class="adm-kicker">STRUCTURE</p>
        <h2>页面结构 <span class="count-chip">{{ page.blocks.length }}</span></h2>
      </div>
      <button class="panel-add-btn" title="添加区块" @click="emit('add')">
        <UiIcon name="plus" :size="15" />
      </button>
    </div>

    <div class="block-list">
      <VueDraggable
        :model-value="page.blocks"
        :animation="150"
        handle=".block-drag"
        class="block-drag-list"
        @update:model-value="onReorder"
      >
        <div
          v-for="block in page.blocks"
          :key="block.id"
          class="block-row"
          :class="{ active: selectedId === block.id && mode === 'block' }"
          @click="emit('select', block)"
        >
          <span class="block-drag">
            <UiIcon name="drag" :size="12" />
          </span>
          <span class="block-icon">
            <UiIcon :name="blockSchemas[block.type]?.icon || 'blocks'" :size="14" />
          </span>
          <span class="block-text">
            <strong>{{ blockSchemas[block.type]?.label || block.type }}</strong>
            <small>{{ blockSummary(block) }}</small>
          </span>
          <span class="block-actions" @click.stop>
            <button title="复制区块" @click="emit('duplicate', block)">
              <UiIcon name="copy" :size="12" />
            </button>
            <button class="delete" title="删除区块" @click="emit('remove', block)">
              <UiIcon name="trash" :size="12" />
            </button>
          </span>
        </div>
      </VueDraggable>

      <button class="add-block-btn" @click="emit('add')">
        <UiIcon name="plus" :size="13" /> 添加区块
        <span class="mono add-types">{{ blockList.length }} TYPES</span>
      </button>
    </div>

    <div class="panel-foot">
      <button
        class="settings-link"
        :class="{ active: mode === 'page' }"
        @click="emit('page-settings')"
      >
        <UiIcon name="settings" :size="14" />
        页面设置与 SEO
        <UiIcon name="chevron-right" :size="12" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import type { PageRecord, PageBlock } from '~/domain/types'
import { blockList, blockSchemas, blockSummary } from '@/domain/blockSchema'
import UiIcon from '~/components/UiIcon.vue'

defineProps<{
  page: PageRecord
  selectedId: string | null
  mode: 'block' | 'page'
}>()

const emit = defineEmits<{
  select: [block: PageBlock]
  add: []
  duplicate: [block: PageBlock]
  remove: [block: PageBlock]
  reorder: [blocks: PageBlock[]]
  'page-settings': []
}>()

function onReorder(value: PageBlock[]) {
  emit('reorder', value)
}
</script>

<style scoped>
.block-drag-list {
  min-height: 20px;
}

.add-types {
  font-size: 7px;
  color: var(--adm-muted-2);
  margin-left: 2px;
}
</style>

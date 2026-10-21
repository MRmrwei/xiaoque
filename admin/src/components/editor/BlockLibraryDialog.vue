<template>
  <el-dialog
    :model-value="open"
    title="区块库"
    width="880px"
    top="6vh"
    class="library-dialog"
    @update:model-value="emit('close')"
  >
    <div class="lib-search">
      <el-input v-model="search" placeholder="搜索区块名称或用途…" clearable>
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <span class="lib-count mono">{{ count }} 个区块</span>
    </div>

    <div v-for="group in filteredGroups" :key="group.id" class="lib-group">
      <div class="lib-group-head">
        <span>{{ group.label }}</span>
        <span class="mono">{{ group.hint }}</span>
      </div>
      <div class="lib-grid">
        <button
          v-for="block in group.items"
          :key="block.type"
          type="button"
          class="lib-item"
          @click="add(block.type)"
        >
          <span class="lib-icon">
            <UiIcon :name="block.icon" :size="17" />
          </span>
          <span class="lib-text">
            <span class="lib-title">{{ block.label }} <em>{{ block.labelEn }}</em></span>
            <span class="lib-desc">{{ block.description }}</span>
          </span>
          <span class="lib-plus">＋</span>
        </button>
      </div>
    </div>

    <el-empty v-if="!count" description="没有匹配的区块，换个关键词试试" :image-size="60" />

    <p class="lib-foot">新增区块会追加到当前页面末尾，可随时拖拽调整顺序。</p>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { blockGroups, blockList } from '@/domain/blockSchema'
import type { BlockType } from '~/domain/blocks'
import UiIcon from '~/components/UiIcon.vue'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  close: []
  add: [type: BlockType]
}>()

const search = ref('')

watch(
  () => props.open,
  (open) => {
    if (open) search.value = ''
  }
)

const filteredGroups = computed(() => {
  const keyword = search.value.trim().toLowerCase()
  return blockGroups
    .map((group) => ({
      ...group,
      items: blockList.filter((block) => {
        if (block.group !== group.id) return false
        if (!keyword) return true
        return `${block.label} ${block.labelEn} ${block.description} ${block.type}`
          .toLowerCase()
          .includes(keyword)
      })
    }))
    .filter((group) => group.items.length)
})

const count = computed(() =>
  filteredGroups.value.reduce((total, group) => total + group.items.length, 0)
)

function add(type: string) {
  emit('add', type as BlockType)
}
</script>

<style scoped>
.lib-search {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 18px;
}

.lib-search .el-input {
  max-width: 360px;
}

.lib-count {
  font-size: 9px;
  color: var(--adm-muted-2);
}

.lib-group + .lib-group {
  margin-top: 20px;
}

.lib-group-head {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 9px;
  padding-bottom: 7px;
  border-bottom: 1px solid var(--adm-line-soft);
}

.lib-group-head > span:first-child {
  font-size: 12px;
  font-weight: 600;
}

.lib-group-head .mono {
  font-size: 8.5px;
  color: var(--adm-muted-2);
}

.lib-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.lib-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px;
  border: 1px solid var(--adm-line);
  border-radius: 9px;
  background: #fff;
  text-align: left;
  transition: border-color 0.18s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.lib-item:hover {
  border-color: #cfc6b7;
  box-shadow: 0 4px 14px rgba(20, 21, 24, 0.06);
  transform: translateY(-1px);
}

.lib-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: grid;
  place-items: center;
  background: #f3efe7;
  color: var(--adm-muted);
  flex: none;
}

.lib-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}

.lib-title {
  font-size: 11.5px;
  font-weight: 600;
  white-space: nowrap;
}

.lib-title em {
  font-style: normal;
  font-family: var(--adm-mono);
  font-size: 7.5px;
  color: var(--adm-muted-2);
  margin-left: 3px;
}

.lib-desc {
  font-size: 10px;
  line-height: 1.45;
  color: var(--adm-muted);
}

.lib-plus {
  color: var(--adm-muted-2);
  font-size: 16px;
  line-height: 1;
}

.lib-foot {
  margin-top: 18px;
  padding-top: 12px;
  border-top: 1px solid var(--adm-line-soft);
  font-size: 10px;
  color: var(--adm-muted-2);
}
</style>

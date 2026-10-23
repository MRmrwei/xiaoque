<template>
  <div class="repeater">
    <div
      v-for="(item, index) in items"
      :key="index"
      class="repeater-item"
      :class="{ open: isOpen(index) }"
    >
      <div class="repeater-head">
        <button type="button" class="repeater-toggle" @click="toggle(index)">
          <span class="repeater-index mono">{{ String(index + 1).padStart(2, '0') }}</span>
          <span class="repeater-summary">{{ summary(item) }}</span>
          <el-icon class="repeater-chevron"><ArrowDown /></el-icon>
        </button>
        <div class="repeater-actions">
          <button type="button" title="上移" :disabled="index === 0" @click="move(index, -1)">↑</button>
          <button
            type="button"
            title="下移"
            :disabled="index === items.length - 1"
            @click="move(index, 1)"
          >
            ↓
          </button>
          <button type="button" class="danger" title="删除" @click="remove(index)">×</button>
        </div>
      </div>

      <div v-if="isOpen(index)" class="repeater-body">
        <div
          v-for="itemField in visibleFields(item)"
          :key="itemField.key"
          class="repeater-field"
          :class="{ 'repeater-field--switch': itemField.type === 'switch' }"
        >
          <span class="prop-label">
            {{ itemField.label }}
            <span v-if="itemField.required" class="req">*</span>
          </span>
          <RepeaterField
            v-if="itemField.type === 'repeater'"
            :field="itemField"
            :model-value="item[itemField.key]"
            @update:model-value="update(index, itemField.key, $event)"
          />
          <FieldControl
            v-else
            :field="itemField"
            :model-value="item[itemField.key]"
            @update:model-value="update(index, itemField.key, $event)"
          />
          <span v-if="itemField.help" class="prop-help">{{ itemField.help }}</span>
        </div>
      </div>
    </div>

    <button type="button" class="repeater-add" @click="add">
      <el-icon><Plus /></el-icon> 添加{{ field.itemLabel || '条目' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { ArrowDown, Plus } from '@element-plus/icons-vue'
import { isFieldVisible, type BlockField } from '@/domain/blockSchema'
import FieldControl from './FieldControl.vue'

const props = defineProps<{
  field: BlockField
  modelValue: unknown
}>()

const emit = defineEmits<{
  'update:modelValue': [value: unknown]
}>()

const items = computed<Record<string, unknown>[]>(() =>
  Array.isArray(props.modelValue) ? (props.modelValue as Record<string, unknown>[]) : []
)
const openIndexes = ref<Set<number>>(new Set([0]))

function visibleFields(item: Record<string, unknown>): BlockField[] {
  return (props.field.itemFields ?? []).filter((field) => isFieldVisible(field, item))
}

function isOpen(index: number) {
  return openIndexes.value.has(index)
}

function toggle(index: number) {
  const next = new Set(openIndexes.value)
  if (next.has(index)) next.delete(index)
  else next.add(index)
  openIndexes.value = next
}

// 折叠时优先显示「标题」字段；标题为空再回退到第一个有内容的文本字段
// （新闻 / 产品 / 案例卡片的第一项是「分类」「角标」之类的短标签，不适合当标题）。
function textValue(itemField: BlockField, item: Record<string, unknown>): string {
  if (itemField.type !== 'text' && itemField.type !== 'textarea') return ''
  const value = item[itemField.key]
  return typeof value === 'string' ? value.trim() : ''
}

function summary(item: Record<string, unknown>): string {
  const fields = props.field.itemFields ?? []
  const titleField = fields.find((itemField) => itemField.key === 'title')
  const title = titleField ? textValue(titleField, item) : ''
  if (title) return title
  for (const itemField of fields) {
    const value = textValue(itemField, item)
    if (value) return value
  }
  return props.field.itemLabel || '条目'
}

function defaultValue(field: BlockField): unknown {
  if (field.type === 'switch') return false
  if (field.type === 'number') return field.min ?? 0
  if (field.type === 'select' || field.type === 'image' || field.type === 'images') {
    return field.options?.[0]?.value ?? ''
  }
  if (field.type === 'repeater') return []
  return ''
}

function update(index: number, key: string, value: unknown) {
  const next = items.value.map((item, itemIndex) =>
    itemIndex === index ? { ...item, [key]: value } : item
  )
  emit('update:modelValue', next)
}

function add() {
  const item: Record<string, unknown> = {}
  for (const itemField of props.field.itemFields ?? []) {
    item[itemField.key] = defaultValue(itemField)
  }
  const next = [...items.value, item]
  emit('update:modelValue', next)
  openIndexes.value = new Set([...openIndexes.value, next.length - 1])
}

function remove(index: number) {
  const next = items.value.filter((_, itemIndex) => itemIndex !== index)
  emit('update:modelValue', next)
  openIndexes.value = new Set()
}

function move(index: number, direction: number) {
  const target = index + direction
  if (target < 0 || target >= items.value.length) return
  const next = [...items.value]
  const [item] = next.splice(index, 1)
  next.splice(target, 0, item)
  emit('update:modelValue', next)
  openIndexes.value = new Set([target])
}
</script>

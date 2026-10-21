<template>
  <div class="field-control">
    <el-input
      v-if="field.type === 'text' || field.type === 'link'"
      :model-value="stringValue"
      :placeholder="field.type === 'link' ? '例如 /contact 或 https://…' : ''"
      clearable
      @update:model-value="update"
    />
    <el-input
      v-else-if="field.type === 'textarea'"
      type="textarea"
      :rows="field.rows || 3"
      :model-value="stringValue"
      @update:model-value="update"
    />
    <RichTextEditor
      v-else-if="field.type === 'richtext'"
      :model-value="stringValue"
      @update:model-value="update"
    />
    <el-input
      v-else-if="field.type === 'html'"
      type="textarea"
      :rows="field.rows || 8"
      class="code-input"
      :model-value="stringValue"
      @update:model-value="update"
    />
    <el-input-number
      v-else-if="field.type === 'number'"
      :model-value="numberValue"
      :min="field.min"
      :max="field.max"
      controls-position="right"
      class="full-control"
      @update:model-value="update"
    />
    <el-switch
      v-else-if="field.type === 'switch'"
      :model-value="Boolean(modelValue)"
      @update:model-value="update"
    />
    <el-select
      v-else-if="field.type === 'select'"
      :model-value="stringValue"
      class="full-control"
      @update:model-value="update"
    >
      <el-option
        v-for="option in normalizedOptions"
        :key="option.value"
        :label="option.label"
        :value="option.value"
      />
    </el-select>
    <ImagePicker
      v-else-if="field.type === 'image' && field.upload"
      :model-value="stringValue"
      :options="field.options"
      :upload-only="field.uploadOnly"
      @update:model-value="update"
    />
    <ArtPicker
      v-else-if="field.type === 'image' || field.type === 'images'"
      :model-value="stringValue"
      :options="field.options"
      @update:model-value="update"
    />
    <el-color-picker
      v-else-if="field.type === 'color'"
      :model-value="stringValue"
      @update:model-value="update"
    />
    <el-date-picker
      v-else-if="field.type === 'date'"
      type="date"
      value-format="YYYY-MM-DD"
      class="full-control"
      :model-value="stringValue"
      @update:model-value="update"
    />
    <div v-else-if="field.type === 'collection'" class="collection-note">
      数据由内容管理自动提供
    </div>
    <el-input v-else :model-value="stringValue" @update:model-value="update" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { BlockField, FieldOption } from '@/domain/blockSchema'
import ArtPicker from './ArtPicker.vue'
import ImagePicker from './ImagePicker.vue'
import RichTextEditor from './RichTextEditor.vue'

const props = defineProps<{
  field: BlockField
  modelValue: unknown
}>()

const emit = defineEmits<{
  'update:modelValue': [value: unknown]
}>()

const stringValue = computed(() => String(props.modelValue ?? ''))
const numberValue = computed(() => {
  const value = Number(props.modelValue)
  return Number.isFinite(value) ? value : 0
})
const normalizedOptions = computed<FieldOption[]>(() => {
  const options = props.field.options ?? []
  return options.map((option) =>
    typeof option === 'string' ? { value: option, label: option } : option
  )
})

function update(value: unknown) {
  emit('update:modelValue', value)
}
</script>

<style scoped>
.field-control {
  width: 100%;
}

.full-control {
  width: 100%;
}

.code-input :deep(textarea) {
  font-family: var(--adm-mono);
  font-size: 10.5px;
  line-height: 1.7;
}

.collection-note {
  font-size: 10.5px;
  color: var(--adm-muted-2);
  background: var(--adm-surface);
  border: 1px dashed var(--adm-line);
  border-radius: 7px;
  padding: 9px 11px;
}
</style>

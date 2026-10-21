<template>
  <div class="prop-form">
    <div
      v-for="field in visibleFields"
      :key="field.key"
      class="prop-field"
      :class="{ 'prop-field--switch': field.type === 'switch' }"
    >
      <span class="prop-label">
        {{ field.label }}
        <span v-if="field.required" class="req">*</span>
      </span>

      <RepeaterField
        v-if="field.type === 'repeater'"
        :field="field"
        :model-value="model[field.key]"
        @update:model-value="emit('update', field.key, $event)"
      />
      <FieldControl
        v-else
        :field="field"
        :model-value="model[field.key]"
        @update:model-value="emit('update', field.key, $event)"
      />

      <span v-if="field.help" class="prop-help">{{ field.help }}</span>
    </div>

    <div v-if="!schema.fields.length" class="props-empty">该区块没有可配置的属性。</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { isFieldVisible, type BlockField, type BlockSchema } from '@/domain/blockSchema'
import FieldControl from './FieldControl.vue'
import RepeaterField from './RepeaterField.vue'

const props = defineProps<{
  schema: BlockSchema
  model: Record<string, unknown>
}>()

const visibleFields = computed(() =>
  props.schema.fields.filter((field) => isFieldVisible(field, props.model))
)

const emit = defineEmits<{
  update: [key: BlockField['key'], value: unknown]
}>()
</script>

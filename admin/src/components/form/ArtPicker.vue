<template>
  <div class="art-picker">
    <button
      v-for="option in normalizedOptions"
      :key="option.value"
      type="button"
      class="art-option"
      :class="{ active: modelValue === option.value }"
      @click="$emit('update:modelValue', option.value)"
    >
      <span class="art-thumb">
        <VisualArt
          v-if="option.value !== 'builder'"
          :variant="option.value"
          tone="cream"
        />
        <span v-else class="builder-thumb">
          <span class="builder-mark">鹊</span>
          <span class="builder-lines"><i /><i /><i /></span>
        </span>
      </span>
      <span class="art-label">{{ option.label }}</span>
      <span v-if="modelValue === option.value" class="art-check">
        <el-icon><Check /></el-icon>
      </span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Check } from '@element-plus/icons-vue'
import type { FieldOption } from '@/domain/blockSchema'
import VisualArt from '~/components/VisualArt.vue'

const props = defineProps<{
  modelValue: string
  options?: FieldOption[]
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const fallback: FieldOption[] = [
  { value: 'orbit', label: '轨道 Orbit' },
  { value: 'arch', label: '拱门 Arch' },
  { value: 'grid', label: '网格 Grid' },
  { value: 'dunes', label: '沙丘 Dunes' },
  { value: 'waves', label: '波浪 Waves' },
  { value: 'type', label: '字体 Type' },
  { value: 'stack', label: '层叠 Stack' },
  { value: 'bloom', label: '绽放 Bloom' }
]

const normalizedOptions = computed<FieldOption[]>(() =>
  props.options?.length ? props.options : fallback
)
</script>

<style scoped>
.art-thumb {
  background: #f3efe7;
}

.builder-thumb {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  background: linear-gradient(145deg, #efe9dd, #f7f3eb);
}

.builder-mark {
  width: 18px;
  height: 18px;
  border-radius: 5px;
  display: grid;
  place-items: center;
  background: var(--adm-ink);
  color: #f6f2ea;
  font-family: var(--adm-serif);
  font-size: 10px;
}

.builder-lines {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.builder-lines i {
  display: block;
  width: 22px;
  height: 2px;
  border-radius: 2px;
  background: rgba(20, 21, 24, 0.2);
}

.builder-lines i:nth-child(2) {
  width: 16px;
  background: var(--adm-accent);
}

.builder-lines i:nth-child(3) {
  width: 10px;
}
</style>

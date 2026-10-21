<template>
  <div class="image-picker">
    <div v-if="!uploadOnly" class="image-tabs">
      <button
        type="button"
        class="image-tab"
        :class="{ active: mode === 'art' }"
        @click="mode = 'art'"
      >
        内置图形
      </button>
      <button
        type="button"
        class="image-tab"
        :class="{ active: mode === 'upload' }"
        @click="mode = 'upload'"
      >
        上传图片
      </button>
    </div>

    <ArtPicker
      v-if="!uploadOnly && mode === 'art'"
      :model-value="modelValue"
      :options="options"
      @update:model-value="emit('update:modelValue', $event)"
    />

    <div v-else class="upload-panel">
      <div v-if="custom" class="upload-preview">
        <img :src="modelValue" alt="已上传图片预览" />
        <div class="upload-foot">
          <span class="upload-meta">{{ sizeLabel }}</span>
          <div class="upload-actions">
            <el-button size="small" @click="pick">更换</el-button>
            <el-button size="small" @click="remove">移除</el-button>
          </div>
        </div>
      </div>

      <button
        v-else
        type="button"
        class="upload-drop"
        :class="{ busy: loading }"
        @click="pick"
        @dragover.prevent
        @drop.prevent="onDrop"
      >
        <span class="upload-icon"><UiIcon name="image" :size="20" /></span>
        <span class="upload-title">{{ loading ? '正在处理图片…' : '点击或拖拽图片到这里' }}</span>
        <span class="upload-hint">JPG / PNG / WebP / GIF / SVG，自动压缩到 1600px 以内</span>
      </button>

      <p class="upload-note">
        图片会压缩后保存在浏览器本地演示数据中，不会上传到服务器。
      </p>

      <input
        ref="fileInput"
        class="upload-input"
        type="file"
        accept="image/*"
        @change="onPick"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FieldOption } from '@/domain/blockSchema'
import { estimateDataUrlBytes, fileToImageDataUrl, formatBytes, isCustomImage } from '@/utils/image'
import UiIcon from '~/components/UiIcon.vue'
import ArtPicker from './ArtPicker.vue'

const props = defineProps<{
  modelValue: string
  options?: FieldOption[]
  /** 只能上传图片，不展示内置图形选项。 */
  uploadOnly?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const mode = ref<'art' | 'upload'>(
  props.uploadOnly || isCustomImage(props.modelValue) ? 'upload' : 'art'
)
const loading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const custom = computed(() => isCustomImage(props.modelValue))
const sizeLabel = computed(() => {
  if (!custom.value) return ''
  return formatBytes(estimateDataUrlBytes(props.modelValue))
})

watch(
  () => props.modelValue,
  (value) => {
    if (isCustomImage(value)) mode.value = 'upload'
  }
)

watch(
  () => props.options,
  () => {
    mode.value =
      props.uploadOnly || isCustomImage(props.modelValue) ? 'upload' : 'art'
  }
)

function pick() {
  fileInput.value?.click()
}

function remove() {
  emit('update:modelValue', props.options?.[0]?.value ?? '')
  mode.value = props.uploadOnly ? 'upload' : 'art'
}

function onPick(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (file) void applyFile(file)
}

function onDrop(event: DragEvent) {
  const file = event.dataTransfer?.files?.[0]
  if (file) void applyFile(file)
}

async function applyFile(file: File) {
  if (loading.value) return
  loading.value = true
  try {
    const dataUrl = await fileToImageDataUrl(file)
    emit('update:modelValue', dataUrl)
    mode.value = 'upload'
    ElMessage.success(`图片已更新（${formatBytes(estimateDataUrlBytes(dataUrl))}）`)
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '图片处理失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.image-picker {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.image-tabs {
  display: inline-flex;
  align-self: flex-start;
  padding: 2px;
  border: 1px solid var(--adm-line);
  border-radius: 7px;
  background: var(--adm-surface);
}

.image-tab {
  border: 0;
  background: transparent;
  padding: 4px 10px;
  border-radius: 5px;
  font-size: 10px;
  color: var(--adm-muted);
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.image-tab.active {
  background: #fff;
  color: var(--adm-ink);
  box-shadow: 0 1px 2px rgba(20, 21, 24, 0.12);
}

.upload-panel {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.upload-drop {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  padding: 20px 14px;
  border: 1px dashed var(--adm-line-strong, var(--adm-line));
  border-radius: 8px;
  background: var(--adm-surface);
  color: var(--adm-muted);
  cursor: pointer;
  text-align: center;
  transition: border-color 0.15s ease, background 0.15s ease;
}

.upload-drop:hover,
.upload-drop.busy {
  border-color: var(--adm-accent);
  background: rgba(232, 80, 58, 0.04);
}

.upload-icon {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: #fff;
  color: var(--adm-accent);
  border: 1px solid var(--adm-line);
}

.upload-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--adm-ink-soft, var(--adm-ink));
}

.upload-hint {
  font-size: 9.5px;
  color: var(--adm-muted-2);
}

.upload-note {
  margin: 0;
  font-size: 9.5px;
  line-height: 1.6;
  color: var(--adm-muted-2);
}

.upload-preview {
  border: 1px solid var(--adm-line);
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
}

.upload-preview img {
  display: block;
  width: 100%;
  aspect-ratio: 4 / 3;
  object-fit: cover;
  background: #f3efe7;
}

.upload-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 6px 8px;
}

.upload-meta {
  font-family: var(--adm-mono);
  font-size: 9.5px;
  color: var(--adm-muted);
}

.upload-actions {
  display: flex;
  gap: 4px;
}

.upload-input {
  display: none;
}
</style>

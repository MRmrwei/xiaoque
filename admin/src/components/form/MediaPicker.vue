<template>
  <div class="media-picker">
    <div class="media-picker-row">
      <el-input
        :model-value="modelValue"
        :placeholder="kind === 'video' ? '视频地址，或从素材库选择' : '图片地址，或从素材库选择'"
        clearable
        @update:model-value="update"
      />
      <el-button @click="visible = true">素材库</el-button>
    </div>
    <p class="media-picker-hint">{{ hint }}</p>

    <template v-if="modelValue && (kind !== 'image' || isCustomImage(modelValue))">
      <img
        v-if="kind === 'image'"
        class="media-picker-preview is-image"
        :src="modelValue"
        alt="素材预览"
      />
      <video
        v-else
        class="media-picker-preview"
        :src="modelValue"
        controls
        playsinline
        preload="metadata"
      />
    </template>

    <MediaLibraryDialog
      v-model:visible="visible"
      :kind="kind"
      :selected="modelValue"
      @select="choose"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import MediaLibraryDialog from './MediaLibraryDialog.vue'
import type { ApiMediaItem, MediaKind } from '@/services/mediaApi'
import { isCustomImage } from '@/utils/image'

const props = withDefaults(
  defineProps<{
    modelValue: string
    kind?: MediaKind
  }>(),
  { kind: 'video' }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const visible = ref(false)

const hint = computed(() =>
  props.kind === 'image'
    ? '支持 JPG / PNG / WebP / GIF / AVIF，单张不超过 10MB。'
    : '支持 MP4 / WebM / MOV，单个视频不超过 200MB。'
)

function update(value: string) {
  emit('update:modelValue', value ?? '')
}

function choose(item: ApiMediaItem) {
  update(item.url)
}
</script>

<style scoped>
.media-picker {
  width: 100%;
}

.media-picker-row {
  display: flex;
  gap: 8px;
}

.media-picker-hint {
  margin-top: 6px;
  font-size: 11px;
  color: var(--adm-muted-2);
}

.media-picker-preview {
  display: block;
  width: 100%;
  max-height: 220px;
  margin-top: 10px;
  border-radius: 8px;
  background: #141518;
  object-fit: contain;
}

.media-picker-preview.is-image {
  object-fit: contain;
  background: var(--adm-surface);
}
</style>

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
        :class="{ active: mode === 'library' }"
        @click="mode = 'library'"
      >
        素材库
      </button>
    </div>

    <ArtPicker
      v-if="!uploadOnly && mode === 'art'"
      :model-value="modelValue"
      :options="options"
      @update:model-value="emit('update:modelValue', $event)"
    />

    <div v-else class="library-panel">
      <div v-if="custom" class="library-preview">
        <img :src="modelValue" alt="已选图片预览" />
        <div class="library-foot">
          <span class="library-meta">已选图片</span>
          <div class="library-actions">
            <el-button size="small" @click="libraryOpen = true">更换</el-button>
            <el-button size="small" @click="remove">移除</el-button>
          </div>
        </div>
      </div>

      <template v-else>
        <el-button type="primary" plain @click="libraryOpen = true">
          打开图片素材库
        </el-button>
        <p class="library-note">
          从已上传的图片素材中选择；需要新图片时，在素材库弹窗里上传即可。
        </p>
      </template>
    </div>

    <MediaLibraryDialog
      v-model:visible="libraryOpen"
      kind="image"
      :selected="modelValue"
      @select="onLibrarySelect"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { FieldOption } from '@/domain/blockSchema'
import { isCustomImage } from '@/utils/image'
import ArtPicker from './ArtPicker.vue'
import MediaLibraryDialog from './MediaLibraryDialog.vue'
import type { ApiMediaItem } from '@/services/mediaApi'

const props = defineProps<{
  modelValue: string
  options?: FieldOption[]
  /** 只能从素材库选择图片，不展示内置图形选项。 */
  uploadOnly?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const mode = ref<'art' | 'library'>(
  props.uploadOnly || isCustomImage(props.modelValue) ? 'library' : 'art'
)
const libraryOpen = ref(false)

const custom = computed(() => isCustomImage(props.modelValue))

watch(
  () => props.modelValue,
  (value) => {
    if (isCustomImage(value)) mode.value = 'library'
    else if (!value && !props.uploadOnly) mode.value = 'art'
  }
)

watch(
  () => props.options,
  () => {
    mode.value =
      props.uploadOnly || isCustomImage(props.modelValue) ? 'library' : 'art'
  }
)

function remove() {
  emit('update:modelValue', props.uploadOnly ? '' : props.options?.[0]?.value ?? '')
  mode.value = props.uploadOnly ? 'library' : 'art'
}

function onLibrarySelect(item: ApiMediaItem) {
  emit('update:modelValue', item.url)
  mode.value = 'library'
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

.library-panel {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.library-preview {
  width: 100%;
  border: 1px solid var(--adm-line);
  border-radius: 8px;
  overflow: hidden;
  background: var(--adm-surface);
}

.library-preview img {
  display: block;
  width: 100%;
  max-height: 220px;
  object-fit: contain;
  background: var(--adm-surface);
}

.library-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 6px 8px;
}

.library-meta {
  font-size: 11px;
  color: var(--adm-muted);
}

.library-actions {
  display: flex;
  gap: 4px;
}

.library-note {
  margin: 0;
  font-size: 11px;
  line-height: 1.6;
  color: var(--adm-muted-2);
}
</style>

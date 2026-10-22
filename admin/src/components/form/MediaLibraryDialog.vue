<template>
  <el-dialog
    :model-value="visible"
    title="素材库"
    width="min(880px, 94vw)"
    append-to-body
    destroy-on-close
    @update:model-value="emit('update:visible', $event)"
  >
    <div class="media-tabs">
      <button
        v-for="tab in MEDIA_KINDS"
        :key="tab.value"
        type="button"
        class="media-tab"
        :class="{ active: activeKind === tab.value }"
        @click="selectKind(tab.value)"
      >
        {{ tab.label }}
      </button>
    </div>

    <div class="media-toolbar">
      <el-button type="primary" :loading="uploading" @click="pickFile">
        上传{{ kindLabel }}
      </el-button>
      <el-button :loading="loading" @click="loadItems">刷新</el-button>
      <span class="media-note">{{ hint }}</span>
    </div>

    <div v-loading="loading" class="media-grid">
      <div
        v-for="item in items"
        :key="item.id"
        class="media-card"
        :class="{ active: item.url === selected }"
        role="button"
        tabindex="0"
        @click="choose(item)"
        @keydown.enter="choose(item)"
      >
        <img
          v-if="activeKind === 'image'"
          :src="item.url"
          :alt="item.original_name"
          loading="lazy"
        />
        <video v-else :src="item.url" preload="metadata" muted playsinline />
        <span class="media-name">{{ item.original_name }}</span>
        <span class="media-meta">{{ formatBytes(item.size) }}</span>
        <button
          type="button"
          class="media-delete"
          title="删除素材"
          @click.stop="remove(item)"
        >
          删除
        </button>
      </div>
      <p v-if="!loading && !items.length" class="media-empty">
        还没有上传过{{ kindLabel }}素材，先上传一个吧。
      </p>
    </div>

    <input
      ref="fileInput"
      type="file"
      :accept="accept"
      hidden
      @change="onFileChange"
    />
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatBytes } from '@/utils/image'
import {
  deleteMedia,
  listMedia,
  uploadMedia,
  type ApiMediaItem,
  type MediaKind
} from '@/services/mediaApi'
import { MEDIA_RULES } from '@/utils/media'

const props = defineProps<{
  visible: boolean
  kind: MediaKind
  selected?: string
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  select: [item: ApiMediaItem]
}>()

/** 素材分类；后续新增类型时在这里补一项即可。 */
const MEDIA_KINDS: Array<{ value: MediaKind; label: string }> = [
  { value: 'video', label: '视频' },
  { value: 'image', label: '图片' }
]

const activeKind = ref<MediaKind>(props.kind)
const rule = computed(() => MEDIA_RULES[activeKind.value])
const kindLabel = computed(() => rule.value.label)
const hint = computed(() => rule.value.hint)
const accept = computed(() => rule.value.accept)

const loading = ref(false)
const uploading = ref(false)
const items = ref<ApiMediaItem[]>([])
const fileInput = ref<HTMLInputElement | null>(null)

watch(
  () => props.visible,
  (visible) => {
    if (!visible) return
    activeKind.value = props.kind
    void loadItems()
  }
)

function selectKind(kind: MediaKind) {
  if (activeKind.value === kind) return
  activeKind.value = kind
  void loadItems()
}

async function loadItems() {
  loading.value = true
  try {
    const result = await listMedia(activeKind.value, 1, 60)
    items.value = result.items
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '素材加载失败')
  } finally {
    loading.value = false
  }
}

function pickFile() {
  if (uploading.value) return
  fileInput.value?.click()
}

async function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  if (!rule.value.pattern.test(file.name)) {
    ElMessage.error(`仅支持 ${rule.value.label}支持的格式`)
    return
  }
  if (file.size > rule.value.maxBytes) {
    ElMessage.error(`${rule.value.label}不能超过 ${formatBytes(rule.value.maxBytes)}`)
    return
  }

  uploading.value = true
  try {
    const uploadedKind = activeKind.value
    const item = await uploadMedia(file, uploadedKind)
    items.value = [item, ...items.value.filter((entry) => entry.id !== item.id)]
    if (uploadedKind === props.kind) {
      emit('select', item)
      emit('update:visible', false)
      ElMessage.success(`${rule.value.label}已上传并选用`)
    } else {
      ElMessage.success(`${rule.value.label}已上传到素材库`)
    }
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '上传失败')
  } finally {
    uploading.value = false
  }
}

function choose(item: ApiMediaItem) {
  // 分类页签允许浏览全部素材，但字段有类型要求时不能选错类型。
  if (item.kind !== props.kind) {
    ElMessage.warning(
      `当前字段只能选择${props.kind === 'video' ? '视频' : '图片'}素材`
    )
    return
  }
  emit('select', item)
  emit('update:visible', false)
}

async function remove(item: ApiMediaItem) {
  try {
    await ElMessageBox.confirm(
      `确定删除「${item.original_name}」吗？已引用该素材的页面不会受影响，但素材库中将不再显示。`,
      '删除素材',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  try {
    await deleteMedia(item.id)
    items.value = items.value.filter((entry) => entry.id !== item.id)
    ElMessage.success('素材已删除')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '删除失败')
  }
}
</script>

<style scoped>
.media-tabs {
  display: flex;
  gap: 6px;
  margin-bottom: 14px;
  border-bottom: 1px solid var(--adm-line);
}

.media-tab {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  margin-bottom: -1px;
  border: 0;
  border-bottom: 2px solid transparent;
  background: none;
  font-size: 12px;
  color: var(--adm-muted);
  cursor: pointer;
}

.media-tab.active {
  color: var(--adm-ink);
  border-bottom-color: var(--adm-accent);
  font-weight: 600;
}

.media-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
}

.media-note {
  font-size: 11px;
  color: var(--adm-muted-2);
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
  gap: 12px;
  min-height: 120px;
  max-height: 56vh;
  overflow-y: auto;
}

.media-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  border: 1px solid var(--adm-line);
  border-radius: 9px;
  background: var(--adm-surface);
  text-align: left;
  cursor: pointer;
  transition: border-color 0.16s ease, box-shadow 0.16s ease;
}

.media-card:hover,
.media-card.active {
  border-color: var(--adm-accent);
  box-shadow: 0 0 0 2px rgba(232, 80, 58, 0.12);
}

.media-card img,
.media-card video {
  width: 100%;
  aspect-ratio: 16 / 9;
  border-radius: 6px;
  object-fit: cover;
  background: #141518;
  pointer-events: none;
}

.media-delete {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 2px 8px;
  border: 0;
  border-radius: 999px;
  background: rgba(20, 21, 24, 0.72);
  color: #f6f2ea;
  font-size: 11px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.16s ease, background 0.16s ease;
}

.media-card:hover .media-delete,
.media-card:focus-within .media-delete {
  opacity: 1;
}

.media-delete:hover {
  background: var(--adm-accent);
}

.media-name {
  font-size: 12px;
  color: var(--adm-ink);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.media-meta {
  font-size: 11px;
  color: var(--adm-muted-2);
}

.media-empty {
  grid-column: 1 / -1;
  padding: 40px 0;
  text-align: center;
  font-size: 12px;
  color: var(--adm-muted-2);
}
</style>

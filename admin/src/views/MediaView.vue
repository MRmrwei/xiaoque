<template>
  <div class="adm-page">
    <header class="adm-page-head">
      <div>
        <p class="adm-kicker">MEDIA LIBRARY</p>
        <h1 class="adm-title">素材管理</h1>
        <p class="adm-sub">统一管理官网使用的视频与图片素材，支持上传、重命名和删除。</p>
      </div>
      <div class="adm-actions">
        <el-dropdown :disabled="uploading" @command="pickUpload">
          <el-button type="primary" :loading="uploading">上传素材</el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="video">上传视频</el-dropdown-item>
              <el-dropdown-item command="image">上传图片</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <section class="adm-card media-panel" v-loading="loading">
      <div class="adm-toolbar">
        <el-radio-group v-model="kindFilter" @change="reload">
          <el-radio-button value="all">全部</el-radio-button>
          <el-radio-button value="video">视频</el-radio-button>
          <el-radio-button value="image">图片</el-radio-button>
        </el-radio-group>

        <el-input
          v-model="keyword"
          class="adm-search"
          placeholder="搜索素材名称…"
          clearable
          @keyup.enter="reload"
          @clear="reload"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button @click="reload">搜索</el-button>

        <template v-if="items.length">
          <el-checkbox
            :model-value="allSelected"
            :indeterminate="someSelected && !allSelected"
            @change="toggleAll"
          >
            全选
          </el-checkbox>
          <el-button
            type="danger"
            plain
            :disabled="!selectedIds.length"
            @click="batchRemove"
          >
            删除选中<span v-if="selectedIds.length">（{{ selectedIds.length }}）</span>
          </el-button>
        </template>

        <span class="toolbar-spacer" />
        <span class="media-total">共 {{ total }} 个素材</span>
      </div>

      <div v-if="items.length" class="media-grid">
        <article
          v-for="item in items"
          :key="item.id"
          class="media-card"
          :class="{ selected: isSelected(item.id) }"
        >
          <el-checkbox
            class="media-check"
            :model-value="isSelected(item.id)"
            @click.stop
            @change="(checked: string | number | boolean) => toggleSelect(item.id, Boolean(checked))"
          />
          <button type="button" class="media-thumb" @click="openPreview(item)">
            <img
              v-if="item.kind === 'image'"
              :src="item.url"
              :alt="item.original_name"
              loading="lazy"
            />
            <video v-else :src="item.url" preload="metadata" muted playsinline />
            <span class="media-kind">{{ item.kind === 'video' ? '视频' : '图片' }}</span>
          </button>

          <div class="media-info">
            <strong :title="item.original_name">{{ item.original_name }}</strong>
            <span>{{ formatBytes(item.size) }} · {{ formatTime(item.created_at) }}</span>
          </div>

          <div class="media-actions">
            <el-button link @click="openPreview(item)">预览</el-button>
            <el-button link @click="rename(item)">重命名</el-button>
            <el-button link type="danger" @click="remove(item)">删除</el-button>
          </div>
        </article>
      </div>

      <el-empty v-else-if="!loading" description="还没有素材，点击右上角上传。" />

      <div v-if="total > pageSize" class="media-pager">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          background
          @current-change="load"
        />
      </div>
    </section>

    <el-dialog
      v-model="previewOpen"
      :title="previewItem?.original_name || '素材预览'"
      width="min(760px, 92vw)"
      append-to-body
      destroy-on-close
    >
      <img
        v-if="previewItem?.kind === 'image'"
        class="preview-media"
        :src="previewItem.url"
        :alt="previewItem.original_name"
      />
      <video
        v-else-if="previewItem"
        class="preview-media"
        :src="previewItem.url"
        controls
        autoplay
        playsinline
      />
    </el-dialog>

    <input
      ref="fileInput"
      type="file"
      :accept="uploadAccept"
      hidden
      @change="onUploadPick"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { formatBytes } from '@/utils/image'
import { MEDIA_RULES } from '@/utils/media'
import {
  deleteMedia,
  deleteMediaMany,
  listMedia,
  updateMedia,
  uploadMedia,
  type ApiMediaItem,
  type MediaKind
} from '@/services/mediaApi'

const kindFilter = ref<MediaKind | 'all'>('all')
const keyword = ref('')
const page = ref(1)
const pageSize = 24
const total = ref(0)
const items = ref<ApiMediaItem[]>([])
const loading = ref(false)

const uploading = ref(false)
const uploadKind = ref<MediaKind>('video')
const fileInput = ref<HTMLInputElement | null>(null)
const uploadAccept = computed(() => MEDIA_RULES[uploadKind.value].accept)

const previewOpen = ref(false)
const previewItem = ref<ApiMediaItem | null>(null)

const selectedIds = ref<string[]>([])
const allSelected = computed(
  () => items.value.length > 0 && selectedIds.value.length === items.value.length
)
const someSelected = computed(() => selectedIds.value.length > 0)

onMounted(() => {
  void load()
})

async function load() {
  loading.value = true
  try {
    const result = await listMedia(
      kindFilter.value,
      page.value,
      pageSize,
      keyword.value
    )
    items.value = result.items
    total.value = result.total
    selectedIds.value = []
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '素材加载失败')
  } finally {
    loading.value = false
  }
}

function reload() {
  page.value = 1
  void load()
}

async function pickUpload(kind: MediaKind) {
  if (uploading.value) return
  uploadKind.value = kind
  await nextTick()
  fileInput.value?.click()
}

async function onUploadPick(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  const rule = MEDIA_RULES[uploadKind.value]
  if (!rule.pattern.test(file.name)) {
    ElMessage.error(`仅支持${rule.label}格式：${rule.hint}`)
    return
  }
  if (file.size > rule.maxBytes) {
    ElMessage.error(`${rule.label}不能超过 ${formatBytes(rule.maxBytes)}`)
    return
  }

  uploading.value = true
  try {
    await uploadMedia(file, uploadKind.value)
    ElMessage.success(`${rule.label}已上传`)
    page.value = 1
    await load()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '上传失败')
  } finally {
    uploading.value = false
  }
}

function isSelected(id: string): boolean {
  return selectedIds.value.includes(id)
}

function toggleSelect(id: string, checked: boolean) {
  if (checked) {
    if (!selectedIds.value.includes(id)) {
      selectedIds.value = [...selectedIds.value, id]
    }
    return
  }
  selectedIds.value = selectedIds.value.filter((entry) => entry !== id)
}

function toggleAll(checked: unknown) {
  selectedIds.value = checked ? items.value.map((item) => item.id) : []
}

async function batchRemove() {
  if (!selectedIds.value.length) return

  try {
    await ElMessageBox.confirm(
      `确定删除选中的 ${selectedIds.value.length} 个素材吗？已引用这些素材的页面不会受影响，但素材库中将不再显示。`,
      '批量删除素材',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  const count = selectedIds.value.length
  try {
    await deleteMediaMany(selectedIds.value)
    selectedIds.value = []
    ElMessage.success(`已删除 ${count} 个素材`)
    if (count >= items.value.length && page.value > 1) page.value -= 1
    await load()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '批量删除失败')
  }
}

function openPreview(item: ApiMediaItem) {
  previewItem.value = item
  previewOpen.value = true
}

async function rename(item: ApiMediaItem) {
  try {
    const result = await ElMessageBox.prompt('请输入新的素材名称', '重命名素材', {
      inputValue: item.original_name,
      inputValidator: (value) => {
        const name = (value ?? '').trim()
        if (!name) return '请输入素材名称'
        if (name.length > 255) return '名称不能超过 255 个字符'
        return true
      },
      confirmButtonText: '保存',
      cancelButtonText: '取消'
    })
    const updated = await updateMedia(item.id, result.value.trim())
    const index = items.value.findIndex((entry) => entry.id === item.id)
    if (index >= 0) items.value[index] = updated
    ElMessage.success('素材名称已更新')
  } catch (error) {
    if (error === 'cancel' || error === 'close') return
    ElMessage.error(error instanceof Error ? error.message : '重命名失败')
  }
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
    ElMessage.success('素材已删除')
    if (items.value.length === 1 && page.value > 1) page.value -= 1
    await load()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '删除失败')
  }
}

function formatTime(seconds: number): string {
  if (!seconds) return '—'
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date(seconds * 1000))
}
</script>

<style scoped>
.media-panel {
  padding: 16px;
}

.media-total {
  font-size: 12px;
  color: var(--adm-muted-2);
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 14px;
  margin-top: 16px;
}

.media-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 10px;
  border: 1px solid var(--adm-line);
  border-radius: 10px;
  background: var(--adm-surface);
}

.media-card.selected {
  border-color: var(--adm-accent);
  box-shadow: 0 0 0 2px rgba(232, 80, 58, 0.12);
}

.media-check {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 1;
  padding: 0 5px;
  border-radius: 5px;
  background: rgba(255, 255, 255, 0.9);
}

.media-thumb {
  position: relative;
  display: block;
  width: 100%;
  aspect-ratio: 16 / 9;
  padding: 0;
  border: 0;
  border-radius: 8px;
  overflow: hidden;
  background: #141518;
  cursor: pointer;
}

.media-thumb img,
.media-thumb video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  pointer-events: none;
}

.media-kind {
  position: absolute;
  right: 8px;
  top: 8px;
  padding: 2px 7px;
  border-radius: 999px;
  background: rgba(20, 21, 24, 0.72);
  color: #f6f2ea;
  font-size: 11px;
}

.media-info {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.media-info strong {
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.media-info span {
  font-size: 11px;
  color: var(--adm-muted-2);
}

.media-actions {
  display: flex;
  gap: 4px;
}

.media-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.preview-media {
  display: block;
  width: 100%;
  max-height: 70vh;
  object-fit: contain;
  border-radius: 8px;
  background: #141518;
}
</style>

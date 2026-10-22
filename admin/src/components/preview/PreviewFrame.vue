<template>
  <div class="preview-frame">
    <div class="preview-toolbar">
      <div class="device-switch">
        <button
          v-for="option in devices"
          :key="option.id"
          :class="{ active: device === option.id }"
          :title="option.label"
          @click="device = option.id"
        >
          <UiIcon :name="option.icon" :size="15" />
        </button>
      </div>
      <span class="preview-draft">草稿预览 · noindex</span>
      <span class="preview-ready" :class="{ 'is-ready': ready }">
        <i /> {{ ready ? '已连接' : '等待预览' }}
      </span>
    </div>

    <div class="preview-stage">
      <div class="preview-browser" :style="{ maxWidth: currentDevice.width }">
        <div class="browser-bar">
          <span class="browser-dots"><i /><i /><i /></span>
          <span class="browser-url">{{ urlLabel }}</span>
        </div>
        <div class="preview-scroll">
          <iframe
            ref="frameRef"
            :src="previewUrl"
            title="页面草稿预览"
            @load="onLoad"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { Locale, PageRecord, SiteCollections, SiteSettings } from '~/domain/types'
import type { DeviceMode, PreviewUpdateMessage } from '@/types/admin'
import UiIcon from '~/components/UiIcon.vue'

const props = defineProps<{
  page: PageRecord
  locale: Locale
  collections: SiteCollections
  settings: SiteSettings
  /** 当前选中的区块与点击序号；序号变化时预览滚动到对应区块。 */
  focusRequest?: { id: string; token: number } | null
}>()

const devices: Array<{ id: DeviceMode; label: string; icon: string; width: string }> = [
  { id: 'desktop', label: '桌面', icon: 'monitor', width: '100%' },
  { id: 'tablet', label: '平板', icon: 'tablet', width: '760px' },
  { id: 'mobile', label: '手机', icon: 'phone', width: '390px' }
]

const device = ref<DeviceMode>('desktop')
const currentDevice = computed(
  () => devices.find((item) => item.id === device.value) ?? devices[0]
)
const frameRef = ref<HTMLIFrameElement | null>(null)
const ready = ref(false)
let revision = 0
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const previewUrl = new URL('preview.html', window.location.href).toString()

const urlLabel = computed(() => {
  const prefix = props.locale === 'en' ? '/en' : ''
  const path = props.page.path === '/' ? '' : props.page.path
  return `xiaoque.com${prefix}${path}` || 'xiaoque.com'
})

function sendUpdate() {
  const frame = frameRef.value
  if (!frame?.contentWindow) return
  revision += 1
  const payload: PreviewUpdateMessage = {
    type: 'xq-preview:update',
    revision,
    page: JSON.parse(JSON.stringify(props.page)),
    locale: props.locale,
    collections: JSON.parse(JSON.stringify(props.collections)),
    settings: JSON.parse(JSON.stringify(props.settings)),
    focusBlockId: props.focusRequest?.id,
    focusToken: props.focusRequest?.token
  }
  frame.contentWindow.postMessage(payload, '*')
}

function scheduleUpdate() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(sendUpdate, 150)
}

function onMessage(event: MessageEvent) {
  if (event.source !== frameRef.value?.contentWindow) return
  if (event.data?.type === 'xq-preview:ready') {
    ready.value = true
    sendUpdate()
  }
}

function onLoad() {
  // The preview app also posts a ready message; this is a fallback.
  setTimeout(sendUpdate, 60)
}

watch(
  () => props.page,
  () => scheduleUpdate(),
  { deep: true }
)
watch(
  () => props.locale,
  () => scheduleUpdate()
)
watch(
  () => props.collections,
  () => scheduleUpdate(),
  { deep: true }
)
watch(
  () => props.settings,
  () => scheduleUpdate(),
  { deep: true }
)
// 点击左侧区块时立即推送，让预览尽快滚到对应位置。
watch(
  () => props.focusRequest?.token,
  () => sendUpdate()
)

onMounted(() => {
  window.addEventListener('message', onMessage)
})

onBeforeUnmount(() => {
  window.removeEventListener('message', onMessage)
  if (debounceTimer) clearTimeout(debounceTimer)
})
</script>

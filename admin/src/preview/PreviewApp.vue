<template>
  <div class="preview-app">
    <template v-if="page && settings">
      <PreviewChrome :settings="settings" :locale="locale" />
      <main class="preview-page page-canvas">
        <div
          v-for="block in page.blocks"
          :key="block.id"
          :data-preview-block="block.id"
        >
          <BlockRenderer
            :block="block"
            :locale="locale"
            :collections="collections"
            preview
          />
        </div>
        <div v-if="!page.blocks.length" class="preview-empty">
          这个页面还是空的，在左侧添加一个区块试试。
        </div>
      </main>
      <PreviewFooter :settings="settings" :locale="locale" />
    </template>
    <div v-else class="preview-loading">
      <span class="mono">DRAFT PREVIEW</span>
      <p>正在加载草稿预览…</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watchEffect } from 'vue'
import type { Locale, PageRecord, SiteCollections, SiteSettings } from '~/domain/types'
import BlockRenderer from '~/components/blocks/BlockRenderer.vue'
import PreviewChrome from './PreviewChrome.vue'
import PreviewFooter from './PreviewFooter.vue'
import { setPreviewLocale } from '~/composables/useLocalePath'
import { brandStyle } from '@/utils/brand'

const page = ref<PageRecord | null>(null)
const locale = ref<Locale>('zh')
const collections = ref<SiteCollections>({ news: [], products: [], cases: [] })
const settings = ref<SiteSettings | null>(null)
let lastRevision = -1
let lastFocusToken = 0

// 品牌变量写到根节点，让 body、页脚与所有区块共享同一套配色。
watchEffect(() => {
  const vars = brandStyle(settings.value?.brand)
  for (const [key, value] of Object.entries(vars)) {
    document.documentElement.style.setProperty(key, value)
  }
})

function onMessage(event: MessageEvent) {
  if (event.source !== window.parent) return
  const data = event.data as
    | {
        type?: string
        revision?: number
        page?: PageRecord
        locale?: Locale
        collections?: SiteCollections
        settings?: SiteSettings
        focusBlockId?: string
        focusToken?: number
      }
    | undefined
  if (!data || data.type !== 'xq-preview:update') return
  if (typeof data.revision === 'number' && data.revision < lastRevision) return
  if (typeof data.revision === 'number') lastRevision = data.revision

  if (data.page) page.value = data.page
  if (data.locale === 'zh' || data.locale === 'en') {
    locale.value = data.locale
    setPreviewLocale(data.locale)
  }
  if (data.collections) collections.value = data.collections
  if (data.settings) settings.value = data.settings

  // 点击左侧区块时滚动到对应位置；token 相同说明只是内容更新，不重复滚动。
  const token = typeof data.focusToken === 'number' ? data.focusToken : 0
  if (
    data.focusBlockId &&
    typeof data.focusToken === 'number' &&
    token !== lastFocusToken
  ) {
    lastFocusToken = token
    nextTick(() => scrollToBlock(data.focusBlockId as string))
  }
}

function scrollToBlock(blockId: string) {
  const target = document.querySelector<HTMLElement>(
    `[data-preview-block="${blockId}"]`
  )
  target?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

onMounted(() => {
  window.addEventListener('message', onMessage)
  window.parent.postMessage({ type: 'xq-preview:ready' }, '*')
})

onBeforeUnmount(() => {
  window.removeEventListener('message', onMessage)
})
</script>

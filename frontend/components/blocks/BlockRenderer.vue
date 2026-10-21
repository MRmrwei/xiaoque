<template>
  <component :is="resolved" v-if="resolved" v-bind="blockProps(block)" />
  <div v-else class="unknown-block">
    <span class="unknown-tag mono">UNKNOWN BLOCK</span>
    <p>{{ block.type }}</p>
    <span>该区块类型未注册，发布时会被跳过。</span>
  </div>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import type { BlockType, PageBlock } from '~/domain/blocks'
import type { Locale, SiteCollections } from '~/domain/types'
import HeroBlock from './HeroBlock.vue'
import HeroCarouselBlock from './HeroCarouselBlock.vue'
import PartnersBlock from './PartnersBlock.vue'
import FeaturesBlock from './FeaturesBlock.vue'
import ImageTextBlock from './ImageTextBlock.vue'
import StatsBlock from './StatsBlock.vue'
import ProductsBlock from './ProductsBlock.vue'
import NewsListBlock from './NewsListBlock.vue'
import CasesBlock from './CasesBlock.vue'
import FaqBlock from './FaqBlock.vue'
import CtaBlock from './CtaBlock.vue'
import ContactFormBlock from './ContactFormBlock.vue'
import RichTextBlock from './RichTextBlock.vue'
import VideoBlock from './VideoBlock.vue'
import HtmlBlock from './HtmlBlock.vue'

const props = withDefaults(
  defineProps<{
    block: PageBlock
    locale: Locale
    collections: SiteCollections
    preview?: boolean
  }>(),
  { preview: false }
)

const registry: Record<BlockType, Component> = {
  hero: HeroBlock,
  'hero-carousel': HeroCarouselBlock,
  partners: PartnersBlock,
  features: FeaturesBlock,
  'image-text': ImageTextBlock,
  stats: StatsBlock,
  products: ProductsBlock,
  'news-list': NewsListBlock,
  cases: CasesBlock,
  faq: FaqBlock,
  cta: CtaBlock,
  'contact-form': ContactFormBlock,
  'rich-text': RichTextBlock,
  video: VideoBlock,
  html: HtmlBlock
}

const resolved = computed<Component | null>(() => registry[props.block.type] ?? null)

function blockProps(block: PageBlock): Record<string, unknown> {
  const common: Record<string, unknown> = {
    propsData: block.props,
    locale: props.locale,
    preview: props.preview
  }

  if (block.type === 'products') {
    return { ...common, products: props.collections.products }
  }
  if (block.type === 'news-list') {
    return { ...common, news: props.collections.news }
  }
  if (block.type === 'cases') {
    return { ...common, cases: props.collections.cases }
  }
  return common
}
</script>

<style scoped>
.unknown-block {
  max-width: 720px;
  margin: 40px auto;
  padding: 28px;
  border: 1px dashed var(--line-strong);
  border-radius: var(--radius-md);
  text-align: center;
  color: var(--muted);
  background: rgba(251, 249, 245, 0.5);
}

.unknown-tag {
  font-size: 9px;
  letter-spacing: 0.16em;
  color: var(--accent);
}

.unknown-block p {
  font-size: 18px;
  margin: 10px 0 4px;
  color: var(--ink);
}

.unknown-block span:last-child {
  font-size: 12px;
}
</style>

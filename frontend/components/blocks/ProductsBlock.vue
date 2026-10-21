<template>
  <section class="products-block section">
    <div class="container">
      <header class="products-head">
        <div>
          <p v-if="data.eyebrow" class="eyebrow">{{ data.eyebrow }}</p>
          <h2 v-html="title" />
        </div>
        <div class="products-head-right">
          <p v-if="data.subtitle">{{ data.subtitle }}</p>
          <BlockLink
            v-if="data.moreLabel"
            :href="data.moreHref"
            :preview="preview"
            class="products-more"
          >
            {{ data.moreLabel }} <span class="arr">→</span>
          </BlockLink>
        </div>
      </header>

      <div class="product-grid">
        <template v-if="customCards.length">
          <BlockLink
            v-for="(card, index) in customCards"
            :key="index"
            v-reveal="{ delay: index * 70 }"
            :href="card.href"
            :preview="preview"
            class="product-card"
            :class="{ 'product-card--link': card.href }"
          >
            <span class="product-cover product-cover--custom">
              <BlockImage :src="card.image" :alt="card.title" />
              <span v-if="card.label" class="product-label mono">{{ card.label }}</span>
              <span class="product-open">
                <UiIcon name="arrow-up-right" :size="17" />
              </span>
            </span>
            <span class="product-body">
              <span class="product-title">{{ card.title }}</span>
              <span v-if="card.summary" class="product-summary">{{ card.summary }}</span>
              <span v-if="parseBullets(card.bullets).length" class="product-points">
                <span v-for="(point, i) in parseBullets(card.bullets)" :key="i" class="product-point">
                  <UiIcon name="check" :size="11" /> {{ point }}
                </span>
              </span>
              <span v-if="card.linkLabel || defaultCardLabel" class="product-link product-link--static">
                {{ card.linkLabel || defaultCardLabel }}
                <UiIcon name="arrow" :size="14" />
              </span>
            </span>
          </BlockLink>
        </template>

        <article
          v-for="(product, index) in items"
          v-else
          :key="product.id"
          v-reveal="{ delay: index * 70 }"
          class="product-card"
        >
          <button class="product-cover" @click="openProduct(product)">
            <VisualArt :variant="product.cover || 'orbit'" :tone="product.tone" />
            <span class="product-label mono">{{ product.label }}</span>
            <span class="product-open">
              <UiIcon name="arrow-up-right" :size="17" />
            </span>
          </button>
          <div class="product-body">
            <h3>{{ product.title }}</h3>
            <p class="product-summary">{{ product.summary }}</p>
            <ul class="product-points">
              <li v-for="(point, i) in (product.bullets || []).slice(0, 3)" :key="i">
                <UiIcon name="check" :size="11" /> {{ point }}
              </li>
            </ul>
            <button class="product-link" @click="openProduct(product)">
              {{ defaultCardLabel }}
              <UiIcon name="arrow" :size="14" />
            </button>
          </div>
        </article>
      </div>
    </div>

    <DetailModal :open="Boolean(active)" @close="active = null">
      <div v-if="active" class="product-modal">
        <div class="product-modal-cover">
          <VisualArt :variant="active.cover || 'orbit'" :tone="active.tone" />
        </div>
        <div class="product-modal-body">
          <span class="tag tag--accent">{{ active.label }}</span>
          <h3>{{ active.title }}</h3>
          <p class="product-modal-desc">{{ active.description || active.summary }}</p>
          <ul class="product-points product-points--modal">
            <li v-for="(point, i) in active.bullets || []" :key="i">
              <UiIcon name="check" :size="12" /> {{ point }}
            </li>
          </ul>
          <BlockLink :href="data.modalCtaHref || localePath('/contact')" class="btn btn--primary">
            {{ data.modalCtaLabel || (locale === 'en' ? 'Book a demo' : '预约演示') }}
            <span class="arr">↗</span>
          </BlockLink>
        </div>
      </div>
    </DetailModal>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import BlockImage from '~/components/BlockImage.vue'
import BlockLink from '~/components/BlockLink.vue'
import DetailModal from '~/components/DetailModal.vue'
import UiIcon from '~/components/UiIcon.vue'
import VisualArt from '~/components/VisualArt.vue'
import type { ProductCardItem, ProductsProps } from '~/domain/blocks'
import type { Locale, ProductRecord } from '~/domain/types'
import { titleHtml } from '~/utils/format'

const props = withDefaults(
  defineProps<{
    propsData: ProductsProps
    locale: Locale
    products: ProductRecord[]
    preview?: boolean
  }>(),
  { preview: false }
)

const localePath = useLocalePath()
const data = computed(() => props.propsData)
const title = computed(() => titleHtml(data.value.title || ''))
const items = computed(() => props.products.slice(0, data.value.limit || 3))
const active = ref<ProductRecord | null>(null)
const defaultCardLabel = computed(() =>
  data.value.cardLabel || (props.locale === 'en' ? 'Explore product' : '了解产品')
)

const customCards = computed<ProductCardItem[]>(() =>
  (data.value.items ?? [])
    .filter((item) => item && (item.title || item.summary || item.image || item.href))
    .slice(0, data.value.limit || 3)
)

function parseBullets(value: string): string[] {
  return String(value || '')
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter(Boolean)
    .slice(0, 3)
}

function openProduct(product: ProductRecord) {
  active.value = product
}
</script>

<style scoped>
.products-block {
  background: var(--paper);
  border-top: 1px solid var(--line-soft);
}

.products-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 40px;
  margin-bottom: clamp(34px, 4.4vw, 58px);
}

.products-head h2 {
  margin-top: 18px;
  font-size: clamp(30px, 3.8vw, 52px);
  max-width: 100%;
}

.products-head-right {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 16px;
  max-width: 34ch;
}

.products-head-right p {
  color: var(--muted);
  font-size: 14.5px;
  line-height: 1.7;
}

.products-more {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 13px;
  font-weight: 700;
  border-bottom: 1px solid var(--ink);
  padding-bottom: 3px;
  transition: color 0.2s ease, border-color 0.2s ease;
}

.products-more:hover {
  color: var(--accent);
  border-color: var(--accent);
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.product-card {
  display: flex;
  flex-direction: column;
  background: color-mix(in srgb, var(--paper-light) 62%, transparent);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  overflow: hidden;
  transition: transform 0.45s var(--ease), box-shadow 0.3s ease, border-color 0.3s ease;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
  border-color: var(--line-strong);
}

.product-cover {
  position: relative;
  display: block;
  width: 100%;
  aspect-ratio: 5 / 3.6;
  padding: 0;
  border: 0;
  overflow: hidden;
  background: var(--paper-deep);
  text-align: left;
}

.product-cover :deep(.visual-art) {
  transition: transform 0.8s var(--ease);
}

.product-card:hover .product-cover :deep(.visual-art) {
  transform: scale(1.035);
}

.product-label {
  position: absolute;
  left: 14px;
  top: 14px;
  font-size: 8.5px;
  letter-spacing: 0.14em;
  color: var(--paper);
  background: rgba(20, 21, 24, 0.45);
  backdrop-filter: blur(6px);
  padding: 5px 9px;
  border-radius: 999px;
}

.product-open {
  position: absolute;
  right: 14px;
  bottom: 14px;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--paper);
  color: var(--ink);
  transform: translateY(4px);
  opacity: 0;
  transition: opacity 0.3s ease, transform 0.4s var(--ease);
}

.product-card:hover .product-open {
  opacity: 1;
  transform: none;
}

.product-body {
  padding: 22px 22px 24px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.product-body h3,
.product-title {
  display: block;
  font-size: clamp(19px, 1.8vw, 23px);
  letter-spacing: -0.03em;
  line-height: 1.25;
}

.product-card--link {
  cursor: pointer;
}

.product-summary {
  display: block;
  margin-top: 9px;
  font-size: 13.5px;
  line-height: 1.65;
  color: var(--muted);
}

.product-points {
  margin-top: 18px;
  padding-top: 16px;
  border-top: 1px solid var(--line-soft);
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.product-points li,
.product-point {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
  color: var(--ink-soft);
}

.product-points .ui-icon {
  color: var(--accent);
}

.product-link {
  margin-top: 20px;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  border: 0;
  background: transparent;
  padding: 0;
  font-size: 12.5px;
  font-weight: 700;
  color: var(--ink);
  align-self: flex-start;
  transition: color 0.2s ease, gap 0.3s var(--ease);
}

.product-link:hover {
  color: var(--accent);
  gap: 10px;
}

.product-modal {
  display: grid;
  grid-template-columns: 0.9fr 1.1fr;
}

.product-modal-cover {
  min-height: 420px;
  overflow: hidden;
}

.product-modal-cover :deep(.visual-art) {
  height: 100%;
}

.product-modal-body {
  padding: clamp(30px, 4vw, 52px);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.product-modal-body h3 {
  margin-top: 18px;
  font-size: clamp(26px, 3vw, 38px);
  letter-spacing: -0.035em;
}

.product-modal-desc {
  margin-top: 14px;
  color: var(--muted);
  font-size: 14.5px;
  line-height: 1.75;
}

.product-points--modal {
  margin-top: 22px;
  padding-top: 18px;
  border-top: 1px solid var(--line);
  width: 100%;
  flex: none;
}

.product-points--modal li {
  font-size: 13.5px;
}

.product-modal-body .btn {
  margin-top: 28px;
}

@container canvas (max-width: 960px) {
  .product-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .products-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 18px;
  }
}

@container canvas (max-width: 640px) {
  .product-grid {
    grid-template-columns: 1fr;
  }
  .product-modal {
    grid-template-columns: 1fr;
  }
  .product-modal-cover {
    min-height: 220px;
    aspect-ratio: 16 / 9;
  }
}
</style>

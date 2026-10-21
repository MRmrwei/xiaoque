<template>
  <section class="news-list-block section">
    <div class="container">
      <header v-if="data.eyebrow || data.title || data.subtitle" class="news-head">
        <div>
          <p v-if="data.eyebrow" class="eyebrow">{{ data.eyebrow }}</p>
          <h2 v-if="data.title" v-html="title" />
        </div>
        <p v-if="data.subtitle" class="news-sub">{{ data.subtitle }}</p>
      </header>

      <div v-if="data.showFilters && categories.length > 1" class="news-filters">
        <button :class="{ active: selected === 'all' }" @click="selected = 'all'">
          {{ locale === 'en' ? 'All' : '全部' }}
        </button>
        <button
          v-for="category in categories"
          :key="category"
          :class="{ active: selected === category }"
          @click="selected = category"
        >
          {{ category }}
        </button>
      </div>

      <div
        v-if="customCards.length || filtered.length"
        class="news-grid"
        :class="{ 'news-grid--feature': featureMode }"
      >
        <template v-if="customCards.length">
          <BlockLink
            v-for="(item, index) in filteredCustom"
            :key="index"
            v-reveal="{ delay: Math.min(index, 5) * 60 }"
            :href="item.href"
            :preview="preview"
            class="news-card"
            :class="{ 'news-card--link': item.href }"
          >
            <span class="news-cover">
              <BlockImage :src="item.image" alt="" fallback="grid" />
              <span class="news-arrow"><UiIcon name="arrow-up-right" :size="17" /></span>
            </span>
            <span class="news-body">
              <span class="news-meta">
                <span class="news-category">{{ item.category }}</span>
                <span class="news-date">{{ formatDate(item.date, locale) }}</span>
              </span>
              <span class="news-title">{{ item.title }}</span>
              <span class="news-summary">{{ item.summary }}</span>
              <span v-if="item.readMinutes" class="news-foot">
                <span class="mono">{{ item.readMinutes }} {{ locale === 'en' ? 'min read' : '分钟阅读' }}</span>
                <UiIcon name="arrow" :size="15" />
              </span>
            </span>
          </BlockLink>
        </template>

        <BlockLink
          v-for="(item, index) in filtered"
          v-else
          :key="item.id"
          v-reveal="{ delay: Math.min(index, 5) * 60 }"
          :href="detailPath(item)"
          :preview="preview"
          class="news-card"
        >
          <span class="news-cover">
            <VisualArt :variant="item.cover || 'grid'" :tone="item.tone" />
            <span v-if="item.featured" class="news-featured mono">
              {{ locale === 'en' ? 'FEATURED' : '置顶' }}
            </span>
            <span class="news-arrow"><UiIcon name="arrow-up-right" :size="17" /></span>
          </span>
          <span class="news-body">
            <span class="news-meta">
              <span class="news-category">{{ item.category }}</span>
              <span class="news-date">{{ formatDate(item.publishedAt, locale) }}</span>
            </span>
            <span class="news-title">{{ item.title }}</span>
            <span class="news-summary">{{ item.summary }}</span>
            <span class="news-foot">
              <span class="mono">{{ item.readMinutes }} {{ locale === 'en' ? 'min read' : '分钟阅读' }}</span>
              <UiIcon name="arrow" :size="15" />
            </span>
          </span>
        </BlockLink>
      </div>

      <div v-else class="news-empty">
        <UiIcon name="book" :size="26" />
        <p>{{ locale === 'en' ? 'No stories in this category yet' : '该分类下暂无文章' }}</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import BlockImage from '~/components/BlockImage.vue'
import BlockLink from '~/components/BlockLink.vue'
import UiIcon from '~/components/UiIcon.vue'
import VisualArt from '~/components/VisualArt.vue'
import type { NewsCardItem, NewsListProps } from '~/domain/blocks'
import type { Locale, NewsRecord } from '~/domain/types'
import { formatDate, titleHtml } from '~/utils/format'

const props = withDefaults(
  defineProps<{
    propsData: NewsListProps
    locale: Locale
    news: NewsRecord[]
    preview?: boolean
  }>(),
  { preview: false }
)

const localePath = useLocalePath()
const data = computed(() => props.propsData)
const title = computed(() => titleHtml(data.value.title || ''))
const selected = ref<string>('all')

const items = computed(() => props.news.slice(0, data.value.limit || 3))
const customCards = computed<NewsCardItem[]>(() =>
  (data.value.items ?? [])
    .filter((item) => item && (item.title || item.summary || item.image || item.href))
    .slice(0, data.value.limit || 3)
)
const categories = computed<string[]>(() => [
  ...new Set(
    (customCards.value.length ? customCards.value : items.value).map((item) => item.category)
  )
])

const filtered = computed<NewsRecord[]>(() =>
  selected.value === 'all'
    ? items.value
    : items.value.filter((item) => item.category === selected.value)
)

const filteredCustom = computed<NewsCardItem[]>(() =>
  selected.value === 'all'
    ? customCards.value
    : customCards.value.filter((item) => item.category === selected.value)
)

const featureMode = computed(
  () =>
    Boolean(data.value.showFilters) &&
    (customCards.value.length ? filteredCustom.value.length : filtered.value.length) > 3
)

function detailPath(item: NewsRecord): string {
  return localePath(`/news/${item.slug}`)
}

watch(
  () => props.locale,
  () => {
    selected.value = 'all'
  }
)
</script>

<style scoped>
.news-list-block {
  background: var(--paper);
}

.news-card--link {
  cursor: pointer;
}

.news-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 40px;
  margin-bottom: clamp(28px, 3.6vw, 44px);
}

.news-head h2 {
  margin-top: 18px;
  font-size: clamp(30px, 3.8vw, 52px);
  max-width: 100%;
}

.news-sub {
  max-width: 36ch;
  color: var(--muted);
  font-size: 14.5px;
  line-height: 1.7;
  padding-bottom: 6px;
}

.news-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--line-soft);
}

.news-filters button {
  height: 34px;
  padding: 0 16px;
  border: 1px solid transparent;
  border-radius: 999px;
  background: transparent;
  font-size: 12px;
  color: var(--muted);
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}

.news-filters button:hover {
  color: var(--ink);
  background: rgba(20, 21, 24, 0.04);
}

.news-filters button.active {
  background: var(--ink);
  color: var(--paper);
}

.news-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: clamp(16px, 2vw, 22px);
}

.news-card {
  display: flex;
  flex-direction: column;
  border-top: 1px solid var(--line);
  padding-top: 16px;
  transition: opacity 0.25s ease;
}

.news-card:hover .news-title {
  color: var(--accent);
}

.news-cover {
  position: relative;
  display: block;
  aspect-ratio: 3 / 2;
  border-radius: 12px;
  overflow: hidden;
  background: var(--paper-deep);
  min-width: 0;
}

.news-cover :deep(.visual-art) {
  transition: transform 0.8s var(--ease);
}

.news-card:hover .news-cover :deep(.visual-art) {
  transform: scale(1.035);
}

.news-featured {
  position: absolute;
  top: 12px;
  left: 12px;
  font-size: 8px;
  letter-spacing: 0.14em;
  color: var(--paper);
  background: var(--accent);
  padding: 5px 9px;
  border-radius: 999px;
}

.news-arrow {
  position: absolute;
  right: 12px;
  bottom: 12px;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: color-mix(in srgb, var(--paper-light) 92%, transparent);
  color: var(--ink);
  opacity: 0;
  transform: translateY(4px);
  transition: opacity 0.3s ease, transform 0.4s var(--ease);
}

.news-card:hover .news-arrow {
  opacity: 1;
  transform: none;
}

.news-body {
  display: flex;
  flex-direction: column;
  padding-top: 16px;
  flex: 1;
  min-width: 0;
}

.news-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.news-category {
  font-size: 10.5px;
  font-weight: 600;
  letter-spacing: 0.03em;
  color: var(--accent);
}

.news-date {
  font-size: 11px;
  color: var(--muted-2);
  position: relative;
  padding-left: 10px;
}

.news-date::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--line-strong);
  transform: translateY(-50%);
}

.news-title {
  display: block;
  margin-top: 11px;
  font-size: clamp(17px, 1.6vw, 20px);
  font-weight: 600;
  line-height: 1.4;
  letter-spacing: -0.025em;
  transition: color 0.2s ease;
}

.news-summary {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-top: 9px;
  font-size: 13px;
  line-height: 1.65;
  color: var(--muted);
}

.news-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
  padding-top: 16px;
  color: var(--muted-2);
}

.news-foot .mono {
  font-size: 8.5px;
  letter-spacing: 0.1em;
}

.news-foot .ui-icon {
  transition: transform 0.35s var(--ease), color 0.2s ease;
}

.news-card:hover .news-foot .ui-icon {
  transform: translateX(3px);
  color: var(--accent);
}

/* Featured layout for the journal index */
.news-grid--feature .news-card:first-child {
  grid-column: span 2;
  display: grid;
  grid-template-columns: 1.06fr 0.94fr;
  gap: 22px;
  padding-top: 18px;
}

.news-grid--feature .news-card:first-child .news-cover {
  aspect-ratio: auto;
  width: 100%;
  height: 100%;
  min-height: 280px;
  align-self: stretch;
}

.news-grid--feature .news-card:first-child .news-body {
  padding-top: 4px;
  justify-content: center;
}

.news-grid--feature .news-card:first-child .news-title {
  font-size: clamp(24px, 2.5vw, 34px);
  letter-spacing: -0.035em;
}

.news-grid--feature .news-card:first-child .news-summary {
  -webkit-line-clamp: 3;
  font-size: 14px;
}

.news-empty {
  padding: 70px 20px;
  text-align: center;
  border: 1px dashed var(--line-strong);
  border-radius: var(--radius-md);
  color: var(--muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.news-empty p {
  font-size: 13px;
}

@container canvas (max-width: 960px) {
  .news-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .news-grid--feature .news-card:first-child {
    grid-column: span 2;
  }
  .news-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 18px;
  }
}

@container canvas (max-width: 620px) {
  .news-grid {
    grid-template-columns: 1fr;
  }
  .news-grid--feature .news-card:first-child {
    grid-column: span 1;
    grid-template-columns: 1fr;
  }
  .news-grid--feature .news-card:first-child .news-cover {
    min-height: 0;
    aspect-ratio: 3 / 2;
  }
}
</style>

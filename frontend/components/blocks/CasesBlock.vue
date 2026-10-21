<template>
  <section class="cases-block section">
    <div class="container">
      <header class="cases-head">
        <div>
          <p v-if="data.eyebrow" class="eyebrow">{{ data.eyebrow }}</p>
          <h2 v-html="title" />
        </div>
        <p v-if="data.subtitle" class="cases-sub">{{ data.subtitle }}</p>
      </header>

      <div v-if="data.showFilters && industries.length > 1" class="case-filters">
        <button :class="{ active: selected === 'all' }" @click="selected = 'all'">
          {{ locale === 'en' ? 'All industries' : '全部行业' }}
        </button>
        <button
          v-for="industry in industries"
          :key="industry"
          :class="{ active: selected === industry }"
          @click="selected = industry"
        >
          {{ industry }}
        </button>
      </div>

      <div class="case-grid">
        <template v-if="customCards.length">
          <button
            v-for="(item, index) in filteredCustom"
            :key="index"
            type="button"
            v-reveal="{ delay: index * 60 }"
            class="case-card"
            @click="activeCustom = item"
          >
            <span class="case-cover">
              <BlockImage :src="item.image" :alt="item.title || item.client" />
              <span class="case-metric">
                <span class="case-metric-value">{{ item.metric }}<em>{{ item.metricSuffix }}</em></span>
                <span class="case-metric-label">{{ item.metricLabel }}</span>
              </span>
              <span class="case-open"><UiIcon name="arrow-up-right" :size="17" /></span>
            </span>
            <span class="case-body">
              <span class="case-meta">
                <span class="case-industry">{{ item.industry }}</span>
                <span class="case-year mono">{{ item.year }}</span>
              </span>
              <span class="case-title-row">
                <span class="case-client">{{ item.client || item.title }}</span>
                <UiIcon name="arrow" :size="16" />
              </span>
              <span class="case-summary">{{ item.summary }}</span>
            </span>
          </button>
        </template>

        <button
          v-for="(item, index) in filtered"
          v-else
          :key="item.id"
          v-reveal="{ delay: index * 60 }"
          class="case-card"
          @click="active = item"
        >
          <span class="case-cover">
            <VisualArt :variant="item.cover || 'orbit'" :tone="item.tone" />
            <span class="case-metric">
              <span class="case-metric-value">{{ item.metric }}<em>{{ item.metricSuffix }}</em></span>
              <span class="case-metric-label">{{ item.metricLabel }}</span>
            </span>
            <span class="case-open"><UiIcon name="arrow-up-right" :size="17" /></span>
          </span>
          <span class="case-body">
            <span class="case-meta">
              <span class="case-industry">{{ item.industry }}</span>
              <span class="case-year mono">{{ String(item.publishedAt).slice(0, 4) }}</span>
            </span>
            <span class="case-title-row">
              <span class="case-client">{{ item.client }}</span>
              <UiIcon name="arrow" :size="16" />
            </span>
            <span class="case-summary">{{ item.summary }}</span>
          </span>
        </button>
      </div>
    </div>

    <DetailModal :open="Boolean(active || activeCustom)" @close="closeDetail">
      <div v-if="active" class="case-modal">
        <div class="case-modal-cover">
          <VisualArt :variant="active.cover || 'orbit'" :tone="active.tone" />
          <div class="case-modal-metric">
            <span class="case-metric-value">{{ active.metric }}<em>{{ active.metricSuffix }}</em></span>
            <span class="case-metric-label">{{ active.metricLabel }}</span>
          </div>
        </div>
        <div class="case-modal-body">
          <span class="tag">{{ active.industry }}</span>
          <h3>{{ active.title }}</h3>
          <p class="case-modal-summary">{{ active.summary }}</p>
          <p class="case-modal-extra">{{ detailNote }}</p>
          <div class="case-modal-foot">
            <span class="mono">CASE / {{ String(active.publishedAt).slice(0, 4) }}</span>
            <BlockLink
              :href="data.ctaHref || localePath('/contact')"
              :preview="preview"
              class="btn btn--primary btn--sm"
            >
              {{ detailCtaLabel }}
              <span class="arr">↗</span>
            </BlockLink>
          </div>
        </div>
      </div>

      <div v-else-if="activeCustom" class="case-modal">
        <div class="case-modal-cover">
          <BlockImage
            :src="activeCustom.image"
            :alt="activeCustom.title || activeCustom.client"
            fallback="orbit"
          />
          <div class="case-modal-metric">
            <span class="case-metric-value">
              {{ activeCustom.metric }}<em>{{ activeCustom.metricSuffix }}</em>
            </span>
            <span class="case-metric-label">{{ activeCustom.metricLabel }}</span>
          </div>
        </div>
        <div class="case-modal-body">
          <span class="tag">{{ activeCustom.industry }}</span>
          <h3>{{ activeCustom.title || activeCustom.client }}</h3>
          <p v-if="activeCustom.summary" class="case-modal-summary">{{ activeCustom.summary }}</p>
          <p class="case-modal-extra">{{ detailNote }}</p>
          <div class="case-modal-foot">
            <span class="mono">CASE / {{ activeCustom.year || '—' }}</span>
            <div class="case-modal-actions">
              <BlockLink
                v-if="activeCustom.href"
                :href="activeCustom.href"
                :preview="preview"
                class="btn btn--ghost btn--sm"
              >
                {{ locale === 'en' ? 'View case' : '查看完整案例' }}
                <span class="arr">↗</span>
              </BlockLink>
              <BlockLink
                :href="data.ctaHref || localePath('/contact')"
                :preview="preview"
                class="btn btn--primary btn--sm"
              >
                {{ detailCtaLabel }}
                <span class="arr">↗</span>
              </BlockLink>
            </div>
          </div>
        </div>
      </div>
    </DetailModal>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import BlockImage from '~/components/BlockImage.vue'
import BlockLink from '~/components/BlockLink.vue'
import DetailModal from '~/components/DetailModal.vue'
import UiIcon from '~/components/UiIcon.vue'
import VisualArt from '~/components/VisualArt.vue'
import type { CaseCardItem, CasesProps } from '~/domain/blocks'
import type { CaseRecord, Locale } from '~/domain/types'
import { titleHtml } from '~/utils/format'

const props = withDefaults(
  defineProps<{
    propsData: CasesProps
    locale: Locale
    cases: CaseRecord[]
    preview?: boolean
  }>(),
  { preview: false }
)

const localePath = useLocalePath()
const data = computed(() => props.propsData)
const title = computed(() => titleHtml(data.value.title || ''))
const selected = ref<string>('all')
const active = ref<CaseRecord | null>(null)
const activeCustom = ref<CaseCardItem | null>(null)

const detailNote = computed(
  () =>
    data.value.modalNote ||
    (props.locale === 'en'
      ? 'A focused rollout, a bilingual content model and a publishing workflow the team can own. That is how the numbers above happened.'
      : '一次聚焦的上线、一套双语内容模型，以及一个团队真正能自己运营的发布流程——这就是上面这些数字的来处。')
)
const detailCtaLabel = computed(
  () =>
    data.value.ctaLabel ||
    (props.locale === 'en' ? 'Start a similar project' : '聊聊类似的项目')
)

const showLimit = computed(() => Math.max(1, Number(data.value.limit) || 3))
const allItems = computed(() => props.cases.slice(0, showLimit.value))
const customCards = computed<CaseCardItem[]>(() =>
  (data.value.items ?? [])
    .filter(
      (item) => item && (item.title || item.client || item.summary || item.image || item.href)
    )
    .slice(0, showLimit.value)
)
const industries = computed<string[]>(() => [
  ...new Set(
    (customCards.value.length ? customCards.value : allItems.value).map((item) => item.industry)
  )
])

const filtered = computed<CaseRecord[]>(() =>
  selected.value === 'all'
    ? allItems.value
    : allItems.value.filter((item) => item.industry === selected.value)
)

const filteredCustom = computed<CaseCardItem[]>(() =>
  selected.value === 'all'
    ? customCards.value
    : customCards.value.filter((item) => item.industry === selected.value)
)

function closeDetail() {
  active.value = null
  activeCustom.value = null
}

watch(
  () => props.locale,
  () => {
    selected.value = 'all'
    closeDetail()
  }
)

// 编辑器改动数据或条数时，关闭旧弹层，避免展示过期内容。
watch(
  () => [props.cases, props.propsData],
  () => closeDetail(),
  { deep: true }
)
</script>

<style scoped>
.cases-block {
  background: var(--paper);
}

.cases-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 40px;
  margin-bottom: clamp(30px, 4vw, 48px);
}

.cases-head h2 {
  margin-top: 18px;
  font-size: clamp(30px, 3.8vw, 52px);
  max-width: 100%;
}

.cases-sub {
  max-width: 36ch;
  color: var(--muted);
  font-size: 14.5px;
  line-height: 1.7;
  padding-bottom: 6px;
}

.case-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 30px;
}

.case-filters button {
  height: 34px;
  padding: 0 16px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: transparent;
  font-size: 12px;
  font-weight: 500;
  color: var(--muted);
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}

.case-filters button:hover {
  border-color: var(--line-strong);
  color: var(--ink);
}

.case-filters button.active {
  background: var(--ink);
  border-color: var(--ink);
  color: var(--paper);
}

.case-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.case-card {
  display: flex;
  flex-direction: column;
  padding: 0;
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--paper-light) 60%, transparent);
  overflow: hidden;
  text-align: left;
  transition: transform 0.45s var(--ease), box-shadow 0.3s ease, border-color 0.3s ease;
}

.case-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
  border-color: var(--line-strong);
}

.case-cover {
  position: relative;
  display: block;
  aspect-ratio: 4 / 2.9;
  overflow: hidden;
  background: var(--paper-deep);
}

.case-cover :deep(.visual-art) {
  transition: transform 0.8s var(--ease);
}

.case-card:hover .case-cover :deep(.visual-art) {
  transform: scale(1.035);
}

.case-cover::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(20, 21, 24, 0.52), transparent 46%);
  pointer-events: none;
}

.case-metric {
  position: absolute;
  left: 18px;
  bottom: 15px;
  z-index: 1;
  display: flex;
  flex-direction: column;
  color: var(--paper);
}

.case-metric-value {
  font-family: var(--font-serif);
  font-size: clamp(30px, 3vw, 40px);
  line-height: 1;
  letter-spacing: -0.03em;
  display: flex;
  align-items: baseline;
}

.case-metric-value em {
  font-family: var(--font-sans);
  font-style: normal;
  font-size: 0.4em;
  font-weight: 600;
  margin-left: 2px;
  color: #ffab99;
}

.case-metric-label {
  font-size: 10.5px;
  letter-spacing: 0.04em;
  opacity: 0.78;
  margin-top: 3px;
}

.case-open {
  position: absolute;
  right: 14px;
  bottom: 14px;
  z-index: 1;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--paper);
  color: var(--ink);
  opacity: 0;
  transform: translateY(4px);
  transition: opacity 0.3s ease, transform 0.4s var(--ease);
}

.case-card:hover .case-open {
  opacity: 1;
  transform: none;
}

.case-body {
  display: flex;
  flex-direction: column;
  padding: 20px 20px 22px;
  flex: 1;
}

.case-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.case-industry {
  font-size: 10.5px;
  letter-spacing: 0.05em;
  color: var(--accent);
  font-weight: 600;
}

.case-year {
  font-size: 9px;
  letter-spacing: 0.1em;
  color: var(--muted-2);
}

.case-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 12px;
}

.case-client {
  font-size: clamp(17px, 1.7vw, 21px);
  font-weight: 600;
  letter-spacing: -0.025em;
}

.case-title-row .ui-icon {
  color: var(--muted-2);
  transition: transform 0.35s var(--ease), color 0.2s ease;
}

.case-card:hover .case-title-row .ui-icon {
  transform: translateX(3px);
  color: var(--accent);
}

.case-summary {
  margin-top: 8px;
  font-size: 13px;
  line-height: 1.65;
  color: var(--muted);
}

.case-modal-cover {
  position: relative;
  height: clamp(220px, 34vw, 340px);
  overflow: hidden;
}

.case-modal-cover :deep(.visual-art) {
  height: 100%;
}

.case-modal-cover::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(20, 21, 24, 0.6), transparent 55%);
}

.case-modal-metric {
  position: absolute;
  left: clamp(24px, 4vw, 44px);
  bottom: 24px;
  z-index: 1;
  color: var(--paper);
  display: flex;
  flex-direction: column;
}

.case-modal-metric .case-metric-value {
  font-size: clamp(46px, 7vw, 72px);
}

.case-modal-metric .case-metric-label {
  font-size: 12px;
}

.case-modal-body {
  padding: clamp(26px, 4vw, 44px);
}

.case-modal-body h3 {
  margin-top: 16px;
  font-size: clamp(28px, 3.4vw, 42px);
  letter-spacing: -0.04em;
}

.case-modal-summary {
  margin-top: 14px;
  font-size: 16px;
  line-height: 1.7;
  color: var(--ink-soft);
}

.case-modal-extra {
  margin-top: 12px;
  font-size: 14px;
  line-height: 1.75;
  color: var(--muted);
}

.case-modal-foot {
  margin-top: 28px;
  padding-top: 20px;
  border-top: 1px solid var(--line);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
}

.case-modal-foot .mono {
  font-size: 9px;
  letter-spacing: 0.14em;
  color: var(--muted-2);
}

.case-modal-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

@container canvas (max-width: 960px) {
  .case-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .cases-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 18px;
  }
}

@container canvas (max-width: 620px) {
  .case-grid {
    grid-template-columns: 1fr;
  }
}
</style>

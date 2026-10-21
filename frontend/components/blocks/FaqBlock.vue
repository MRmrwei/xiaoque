<template>
  <section class="faq-block section">
    <div class="container faq-grid">
      <div class="faq-intro">
        <p v-if="data.eyebrow" class="eyebrow">{{ data.eyebrow }}</p>
        <h2 v-html="title" />
        <p v-if="data.subtitle" class="faq-sub">{{ data.subtitle }}</p>
        <div class="faq-contact">
          <span v-if="contactKicker" class="mono">{{ contactKicker }}</span>
          <p v-if="contactText">{{ contactText }}</p>
          <BlockLink
            v-if="contactLabel"
            :href="contactHref || localePath('/contact')"
            class="faq-link"
          >
            {{ contactLabel }} <UiIcon name="arrow-up-right" :size="14" />
          </BlockLink>
        </div>
      </div>

      <div class="faq-list">
        <div
          v-for="(item, index) in data.items"
          :key="index"
          class="faq-item"
          :class="{ open: openIndex === index }"
        >
          <button class="faq-question" :aria-expanded="openIndex === index" @click="toggle(index)">
            <span class="faq-index mono">{{ String(index + 1).padStart(2, '0') }}</span>
            <span class="faq-q">{{ item.q }}</span>
            <span class="faq-toggle"><UiIcon name="plus" :size="16" /></span>
          </button>
          <div class="faq-answer">
            <div class="faq-answer-inner">
              <p>{{ item.a }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import BlockLink from '~/components/BlockLink.vue'
import UiIcon from '~/components/UiIcon.vue'
import type { FaqProps } from '~/domain/blocks'
import type { Locale } from '~/domain/types'
import { titleHtml } from '~/utils/format'

const props = defineProps<{
  propsData: FaqProps
  locale: Locale
}>()

const localePath = useLocalePath()
const data = computed(() => props.propsData)
const title = computed(() => titleHtml(data.value.title || ''))
const contactKicker = computed(() =>
  data.value.contactKicker === undefined ? 'STILL CURIOUS?' : data.value.contactKicker
)
const contactText = computed(() =>
  data.value.contactText === undefined
    ? props.locale === 'en'
      ? 'We are happy to answer directly.'
      : '还有疑问？我们很乐意直接回答。'
    : data.value.contactText
)
const contactLabel = computed(() =>
  data.value.contactLabel === undefined
    ? props.locale === 'en'
      ? 'Contact us'
      : '联系我们'
    : data.value.contactLabel
)
const contactHref = computed(() => data.value.contactHref || '')
const openIndex = ref(0)

function toggle(index: number) {
  openIndex.value = openIndex.value === index ? -1 : index
}

watch(
  () => props.propsData.items,
  () => {
    openIndex.value = 0
  }
)
</script>

<style scoped>
.faq-block {
  background: var(--paper);
  border-top: 1px solid var(--line-soft);
}

.faq-grid {
  display: grid;
  grid-template-columns: minmax(0, 0.8fr) minmax(0, 1.2fr);
  gap: clamp(40px, 7vw, 100px);
  align-items: start;
}

.faq-intro {
  position: sticky;
  top: 110px;
}

.faq-intro h2 {
  margin-top: 18px;
  font-size: clamp(30px, 3.8vw, 52px);
  max-width: 12ch;
}

.faq-sub {
  margin-top: 16px;
  color: var(--muted);
  font-size: 14.5px;
  max-width: 34ch;
}

.faq-contact {
  margin-top: 40px;
  padding-top: 22px;
  border-top: 1px solid var(--line);
  max-width: 300px;
}

.faq-contact .mono {
  font-size: 8.5px;
  letter-spacing: 0.16em;
  color: var(--muted-2);
}

.faq-contact p {
  margin-top: 8px;
  font-size: 13px;
  color: var(--muted);
}

.faq-link {
  margin-top: 12px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  font-weight: 700;
  border-bottom: 1px solid var(--ink);
  padding-bottom: 2px;
  transition: color 0.2s ease, border-color 0.2s ease;
}

.faq-link:hover {
  color: var(--accent);
  border-color: var(--accent);
}

.faq-list {
  border-top: 1px solid var(--line);
}

.faq-item {
  border-bottom: 1px solid var(--line);
}

.faq-question {
  width: 100%;
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 22px 0;
  background: transparent;
  border: 0;
  text-align: left;
}

.faq-index {
  font-size: 9px;
  letter-spacing: 0.1em;
  color: var(--muted-2);
  padding-top: 5px;
  flex: none;
}

.faq-q {
  flex: 1;
  font-size: clamp(15px, 1.5vw, 17.5px);
  font-weight: 600;
  letter-spacing: -0.02em;
  line-height: 1.45;
  transition: color 0.2s ease;
}

.faq-question:hover .faq-q {
  color: var(--accent);
}

.faq-toggle {
  flex: none;
  width: 30px;
  height: 30px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 50%;
  color: var(--muted);
  transition: background 0.3s ease, color 0.3s ease, border-color 0.3s ease,
    transform 0.4s var(--ease);
}

.faq-item.open .faq-toggle {
  background: var(--ink);
  border-color: var(--ink);
  color: var(--paper);
  transform: rotate(45deg);
}

.faq-answer {
  display: grid;
  grid-template-rows: 0fr;
  transition: grid-template-rows 0.45s var(--ease);
}

.faq-item.open .faq-answer {
  grid-template-rows: 1fr;
}

.faq-answer-inner {
  overflow: hidden;
}

.faq-answer-inner p {
  padding: 0 46px 0 31px;
  max-width: 62ch;
  font-size: 14px;
  line-height: 1.8;
  color: var(--muted);
}

.faq-item.open .faq-answer-inner p {
  padding-bottom: 24px;
}

@container canvas (max-width: 900px) {
  .faq-grid {
    grid-template-columns: 1fr;
    gap: 36px;
  }
  .faq-intro {
    position: static;
  }
  .faq-intro h2 {
    max-width: 100%;
  }
  .faq-contact {
    margin-top: 28px;
  }
}

@container canvas (max-width: 620px) {
  .faq-answer-inner p {
    padding-left: 0;
    padding-right: 20px;
  }
}
</style>

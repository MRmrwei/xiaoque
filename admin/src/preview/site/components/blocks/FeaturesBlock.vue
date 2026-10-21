<template>
  <section class="features-block section">
    <div class="container">
      <header class="features-head">
        <div class="features-head-left">
          <p v-if="data.eyebrow" class="eyebrow">{{ data.eyebrow }}</p>
          <h2 v-html="title" />
        </div>
        <p v-if="data.subtitle" class="features-sub">{{ data.subtitle }}</p>
      </header>

      <div
        class="feature-grid"
        :class="[
          `feature-grid--${data.layout || 'grid'}`,
          `feature-grid--n${(data.items || []).length}`
        ]"
      >
        <article
          v-for="(item, index) in data.items"
          :key="index"
          class="feature-card"
          :class="`feature-card--${index + 1}`"
        >
          <span class="feature-index mono">{{ String(index + 1).padStart(2, '0') }}</span>
          <div class="feature-icon">
            <UiIcon :name="item.icon || 'spark'" :size="22" />
          </div>
          <div class="feature-body">
            <div class="feature-title-row">
              <h3>{{ item.title }}</h3>
              <span v-if="item.tag" class="tag">{{ item.tag }}</span>
            </div>
            <p>{{ item.desc }}</p>
          </div>
          <span class="feature-watermark" aria-hidden="true">
            <UiIcon :name="item.icon || 'spark'" :size="110" :stroke-width="0.7" />
          </span>
        </article>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import UiIcon from '~/components/UiIcon.vue'
import type { FeaturesProps } from '~/domain/blocks'
import { titleHtml } from '~/utils/format'

const props = defineProps<{
  propsData: FeaturesProps
}>()

const data = computed(() => props.propsData)
const title = computed(() => titleHtml(data.value.title || ''))
</script>

<style scoped>
.features-block {
  background: var(--paper);
}

.features-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 40px;
  margin-bottom: clamp(34px, 4.5vw, 60px);
}

.features-head h2 {
  margin-top: 18px;
  font-size: clamp(30px, 3.8vw, 52px);
  max-width: 100%;
}

.features-sub {
  max-width: 38ch;
  color: var(--muted);
  font-size: 14.5px;
  line-height: 1.75;
  padding-bottom: 6px;
}

.feature-grid {
  display: grid;
  gap: 14px;
}

.feature-grid--grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.feature-grid--grid.feature-grid--n4 {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.feature-grid--bento {
  grid-template-columns: repeat(3, minmax(0, 1fr));
  grid-auto-rows: minmax(190px, auto);
}

.feature-grid--bento .feature-card--1 {
  grid-row: span 2;
}

.feature-grid--bento .feature-card--2 {
  grid-column: span 2;
}

.feature-grid--bento .feature-card--6 {
  grid-column: span 2;
}

.feature-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 22px;
  padding: clamp(22px, 2.6vw, 32px);
  background: color-mix(in srgb, var(--paper-light) 58%, transparent);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  overflow: hidden;
  transition: transform 0.45s var(--ease), border-color 0.25s ease, background 0.25s ease,
    box-shadow 0.25s ease;
}

.feature-card:hover {
  transform: translateY(-4px);
  border-color: var(--line-strong);
  background: var(--paper-light);
  box-shadow: var(--shadow-sm);
}

.feature-index {
  position: absolute;
  top: 20px;
  right: 22px;
  font-size: 9px;
  letter-spacing: 0.14em;
  color: var(--muted-2);
}

.feature-icon {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border-radius: 13px;
  background: var(--ink);
  color: var(--paper);
  position: relative;
  z-index: 1;
  transition: background 0.3s ease, color 0.3s ease, transform 0.45s var(--ease);
}

.feature-card:hover .feature-icon {
  background: var(--accent);
  transform: rotate(-4deg);
}

.feature-body {
  margin-top: auto;
  position: relative;
  z-index: 1;
}

.feature-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.feature-body h3 {
  font-size: clamp(17px, 1.6vw, 20px);
  letter-spacing: -0.025em;
}

.feature-body p {
  margin-top: 10px;
  font-size: 13.5px;
  line-height: 1.7;
  color: var(--muted);
  max-width: 42ch;
}

.feature-watermark {
  position: absolute;
  right: -20px;
  bottom: -26px;
  color: var(--ink);
  opacity: 0.04;
  transition: opacity 0.3s ease, transform 0.5s var(--ease);
}

.feature-card:hover .feature-watermark {
  opacity: 0.07;
  transform: rotate(-6deg);
}

.feature-grid--bento .feature-card--1 {
  background: var(--ink);
  color: var(--paper);
  border-color: var(--ink);
}

.feature-grid--bento .feature-card--1 h3 {
  color: var(--paper);
  font-size: clamp(22px, 2.2vw, 28px);
}

.feature-grid--bento .feature-card--1 .feature-body p {
  color: rgba(246, 242, 234, 0.58);
  font-size: 14px;
}

.feature-grid--bento .feature-card--1 .feature-index {
  color: rgba(246, 242, 234, 0.35);
}

.feature-grid--bento .feature-card--1 .feature-icon {
  background: var(--accent);
  color: #fff;
}

.feature-grid--bento .feature-card--1 .tag {
  background: rgba(246, 242, 234, 0.08);
  border-color: rgba(246, 242, 234, 0.15);
  color: rgba(246, 242, 234, 0.7);
}

.feature-grid--bento .feature-card--1 .feature-watermark {
  color: var(--paper);
  opacity: 0.05;
  transform: scale(1.7);
  right: 6px;
  bottom: -10px;
}

.feature-grid--bento .feature-card--1:hover .feature-watermark {
  opacity: 0.09;
}

.feature-grid:not(.feature-grid--bento) .feature-card--1 {
  background: var(--paper-light);
}

@container canvas (max-width: 1000px) {
  .feature-grid--grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .feature-grid--bento {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .feature-grid--bento .feature-card--1 {
    grid-row: span 1;
  }
  .feature-grid--bento .feature-card--2,
  .feature-grid--bento .feature-card--6 {
    grid-column: span 1;
  }
  .features-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 18px;
  }
  .features-sub {
    padding-bottom: 0;
  }
}

@container canvas (max-width: 620px) {
  .feature-grid--grid,
  .feature-grid--bento {
    grid-template-columns: 1fr;
    grid-auto-rows: auto;
  }
  .feature-card {
    min-height: 170px;
  }
  .features-head h2 {
    font-size: clamp(28px, 8vw, 38px);
  }
}
</style>

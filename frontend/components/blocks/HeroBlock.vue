<template>
  <section
    class="hero-block"
    :class="[
      `hero-block--${data.theme}`,
      `hero-block--${data.variant}`,
      { 'no-visual': !hasVisual, 'hero-block--en': locale === 'en' }
    ]"
  >
    <div class="hero-grid-bg" aria-hidden="true" />
    <div class="hero-glow" aria-hidden="true" />
    <div class="container hero-inner">
      <div class="hero-copy">
        <p v-if="data.eyebrow" class="eyebrow hero-eyebrow">{{ data.eyebrow }}</p>
        <h1 v-html="title" />
        <p v-if="data.subtitle" class="hero-sub">{{ data.subtitle }}</p>
        <div v-if="data.primaryLabel || data.secondaryLabel" class="btn-row hero-actions">
          <BlockLink
            v-if="data.primaryLabel"
            :href="data.primaryHref"
            :preview="preview"
            class="btn btn--primary btn--lg"
          >
            {{ data.primaryLabel }} <span class="arr">↗</span>
          </BlockLink>
          <BlockLink
            v-if="data.secondaryLabel"
            :href="data.secondaryHref"
            :preview="preview"
            class="btn btn--ghost btn--lg"
          >
            {{ data.secondaryLabel }}
          </BlockLink>
        </div>
        <p v-if="data.note" class="hero-note">
          <span class="note-dot" /> {{ data.note }}
        </p>
        <div v-if="data.variant === 'home' && scrollLabel" class="hero-scroll" aria-hidden="true">
          <span class="mono">{{ scrollLabel }}</span>
          <span class="scroll-line"><i /></span>
        </div>
      </div>

      <div v-if="hasVisual" class="hero-visual">
        <div class="hero-visual-frame">
          <HeroBuilderMock v-if="data.visual === 'builder'" :locale="locale" />
          <BlockImage v-else :src="data.visual" fallback="orbit" />
        </div>
        <span class="visual-corner visual-corner--tl" aria-hidden="true" />
        <span class="visual-corner visual-corner--br" aria-hidden="true" />
      </div>
    </div>
    <div v-if="sideNote" class="hero-side-note mono" aria-hidden="true">{{ sideNote }}</div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BlockImage from '~/components/BlockImage.vue'
import BlockLink from '~/components/BlockLink.vue'
import HeroBuilderMock from '~/components/blocks/HeroBuilderMock.vue'
import type { HeroProps } from '~/domain/blocks'
import type { Locale } from '~/domain/types'
import { titleHtml } from '~/utils/format'

const props = withDefaults(
  defineProps<{
    propsData: HeroProps
    locale: Locale
    preview?: boolean
  }>(),
  { preview: false }
)

const data = computed(() => props.propsData)
const title = computed(() => titleHtml(data.value.title || ''))
const hasVisual = computed(() => data.value.visual && data.value.visual !== 'none')
const sideNote = computed(() =>
  data.value.sideNote === undefined ? 'EST. 2019 / SHANGHAI' : data.value.sideNote
)
const scrollLabel = computed(() =>
  data.value.scrollLabel === undefined ? 'SCROLL' : data.value.scrollLabel
)
</script>

<style scoped>
.hero-block {
  position: relative;
  overflow: hidden;
  background: var(--paper);
  padding-block: clamp(48px, 7vw, 104px) clamp(56px, 8vw, 112px);
  border-bottom: 1px solid var(--line-soft);
}

.hero-block--home {
  min-height: min(84vh, 860px);
  display: flex;
  align-items: center;
}

.hero-block--page {
  padding-block: clamp(48px, 6vw, 92px) clamp(52px, 6vw, 90px);
}

.hero-grid-bg {
  position: absolute;
  inset: 0;
  background-image: linear-gradient(to right, rgba(20, 21, 24, 0.035) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(20, 21, 24, 0.035) 1px, transparent 1px);
  background-size: 88px 88px;
  mask-image: radial-gradient(ellipse 85% 70% at 60% 35%, #000 10%, transparent 75%);
  -webkit-mask-image: radial-gradient(ellipse 85% 70% at 60% 35%, #000 10%, transparent 75%);
  pointer-events: none;
}

.hero-glow {
  position: absolute;
  width: 720px;
  height: 720px;
  right: -180px;
  top: -280px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(232, 80, 58, 0.1), transparent 66%);
  pointer-events: none;
}

.hero-inner {
  position: relative;
  z-index: 2;
  display: grid;
  grid-template-columns: minmax(0, 1.04fr) minmax(0, 0.96fr);
  align-items: center;
  gap: clamp(36px, 5vw, 72px);
}

.hero-block.no-visual .hero-inner {
  grid-template-columns: minmax(0, 880px);
}

.hero-copy {
  position: relative;
  z-index: 2;
}

.hero-eyebrow {
  margin-bottom: 22px;
}

.hero-block h1 {
  font-size: clamp(44px, 5.9vw, 84px);
  line-height: 1.03;
  letter-spacing: -0.05em;
  font-weight: 600;
  max-width: 100%;
}

.hero-block h1 :deep(em) {
  font-family: var(--font-serif);
  font-style: normal;
  font-weight: 500;
  color: var(--accent);
  letter-spacing: -0.015em;
}

.hero-block--page h1 {
  font-size: clamp(36px, 4.3vw, 60px);
}

.hero-block--en h1 {
  font-size: clamp(38px, 4.4vw, 64px);
  letter-spacing: -0.045em;
}

.hero-block--en.hero-block--page h1 {
  font-size: clamp(34px, 3.8vw, 54px);
}

.hero-sub {
  margin-top: 24px;
  max-width: 48ch;
  font-size: clamp(15px, 1.35vw, 17.5px);
  line-height: 1.75;
  color: var(--muted);
}

.hero-block--page .hero-sub {
  max-width: 52ch;
}

.hero-actions {
  margin-top: 34px;
}

.hero-note {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 20px;
  font-size: 12px;
  color: var(--muted-2);
}

.note-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--accent);
}

.hero-scroll {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: clamp(40px, 5vw, 70px);
  color: var(--muted-2);
}

.hero-scroll .mono {
  font-size: 8.5px;
  letter-spacing: 0.22em;
}

.scroll-line {
  width: 54px;
  height: 1px;
  background: var(--line-strong);
  position: relative;
  overflow: hidden;
}

.scroll-line i {
  position: absolute;
  inset: 0;
  background: var(--accent);
  transform-origin: left;
  animation: scrollPulse 2.6s var(--ease) infinite;
}

@keyframes scrollPulse {
  0% {
    transform: scaleX(0);
    opacity: 0;
  }
  35% {
    transform: scaleX(1);
    opacity: 1;
  }
  70% {
    transform: scaleX(1);
    opacity: 0;
  }
  100% {
    transform: scaleX(1);
    opacity: 0;
  }
}

.hero-visual {
  position: relative;
  min-width: 0;
}

.hero-visual-frame {
  position: relative;
  border-radius: clamp(18px, 2vw, 26px);
  overflow: hidden;
  box-shadow: var(--shadow-lg);
  background: var(--paper-light);
  border: 1px solid rgba(20, 21, 24, 0.08);
  transform: rotate(0.6deg);
  transition: transform 0.6s var(--ease);
}

.hero-visual:hover .hero-visual-frame {
  transform: rotate(0deg) translateY(-4px);
}

.hero-block:not(.hero-block--home) .hero-visual-frame {
  aspect-ratio: 1.15 / 1;
}

.visual-corner {
  position: absolute;
  width: 16px;
  height: 16px;
  border-color: var(--accent);
  border-style: solid;
  opacity: 0.7;
}

.visual-corner--tl {
  top: -10px;
  left: -10px;
  border-width: 1px 0 0 1px;
}

.visual-corner--br {
  bottom: -10px;
  right: -10px;
  border-width: 0 1px 1px 0;
}

.hero-side-note {
  position: absolute;
  right: 18px;
  bottom: 90px;
  writing-mode: vertical-rl;
  font-size: 8px;
  letter-spacing: 0.24em;
  color: var(--muted-2);
  opacity: 0.7;
}

.hero-block--dark {
  background: var(--ink);
  color: var(--paper);
  border-bottom-color: rgba(246, 242, 234, 0.08);
}

.hero-block--dark .hero-grid-bg {
  background-image: linear-gradient(to right, rgba(246, 242, 234, 0.05) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(246, 242, 234, 0.05) 1px, transparent 1px);
}

.hero-block--dark .hero-glow {
  background: radial-gradient(circle, rgba(232, 80, 58, 0.22), transparent 66%);
}

.hero-block--dark h1 {
  color: var(--paper);
}

.hero-block--dark h1 :deep(em) {
  color: #ff8b73;
}

.hero-block--dark .hero-sub {
  color: rgba(246, 242, 234, 0.6);
}

.hero-block--dark .eyebrow {
  color: rgba(246, 242, 234, 0.5);
}

.hero-block--dark .hero-note {
  color: rgba(246, 242, 234, 0.4);
}

.hero-block--dark .btn--ghost {
  --btn-fg: var(--paper);
  border-color: rgba(246, 242, 234, 0.25);
}

.hero-block--dark .btn--ghost:hover {
  --btn-bg: var(--paper);
  --btn-fg: var(--ink);
  border-color: var(--paper);
}

@container canvas (max-width: 980px) {
  .hero-inner {
    grid-template-columns: 1fr;
    gap: 46px;
  }
  .hero-block h1 {
    max-width: 100%;
  }
  .hero-visual {
    max-width: 640px;
  }
  .hero-side-note {
    display: none;
  }
  .hero-block--home {
    min-height: auto;
  }
}

@container canvas (max-width: 620px) {
  .hero-block {
    padding-block: 40px 56px;
  }
  .hero-block h1 {
    font-size: clamp(36px, 11vw, 52px);
    max-width: 100%;
    letter-spacing: -0.045em;
  }
  .hero-sub {
    margin-top: 18px;
    font-size: 15px;
  }
  .hero-actions {
    margin-top: 26px;
    gap: 10px;
  }
  .hero-actions .btn {
    height: 48px;
    padding: 0 20px;
    font-size: 13.5px;
  }
  .hero-scroll {
    display: none;
  }
  .hero-eyebrow {
    margin-bottom: 16px;
  }
}
</style>

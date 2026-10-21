<template>
  <section
    class="cta-block"
    :class="[`cta-block--${data.theme || 'accent'}`, { 'cta-block--image': hasImage }]"
  >
    <div class="container">
      <div class="cta-panel">
        <div
          v-if="hasImage"
          class="cta-bg"
          :style="{ backgroundImage: `url(${data.backgroundImage})` }"
          aria-hidden="true"
        />
        <div
          v-if="hasImage"
          class="cta-bg-overlay"
          :style="{ opacity: overlayOpacity }"
          aria-hidden="true"
        />
        <div v-if="!hasImage" class="cta-rings" aria-hidden="true"><i /><i /><i /></div>
        <span v-if="!hasImage" class="cta-glyph" aria-hidden="true">鹊</span>

        <div class="cta-content">
          <div v-if="kicker" class="cta-kicker">
            <span class="dot" />
            <span class="mono">{{ kicker }}</span>
          </div>
          <h2 v-html="title" />
          <p v-if="data.text" class="cta-text">{{ data.text }}</p>
          <div class="btn-row cta-actions">
            <BlockLink
              v-if="data.primaryLabel"
              :href="data.primaryHref"
              :preview="preview"
              class="btn btn--light btn--lg"
            >
              {{ data.primaryLabel }} <span class="arr">↗</span>
            </BlockLink>
            <BlockLink
              v-if="data.secondaryLabel"
              :href="data.secondaryHref"
              :preview="preview"
              class="btn btn--outline-light btn--lg"
            >
              {{ data.secondaryLabel }}
            </BlockLink>
          </div>
        </div>

        <div v-if="sideNote" class="cta-side">
          <span class="cta-side-line" />
          <p class="mono cta-side-note">{{ sideNote }}</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BlockLink from '~/components/BlockLink.vue'
import type { CtaProps } from '~/domain/blocks'
import { titleHtml } from '~/utils/format'

const props = withDefaults(
  defineProps<{
    propsData: CtaProps
    preview?: boolean
  }>(),
  { preview: false }
)

const data = computed(() => props.propsData)
const kicker = computed(() =>
  data.value.kicker === undefined ? "LET'S TALK" : data.value.kicker
)
const sideNote = computed(() =>
  data.value.sideNote === undefined ? '30 MIN\nPRODUCT DEMO' : data.value.sideNote
)
const hasImage = computed(() => Boolean(data.value.backgroundImage?.trim()))
const overlayOpacity = computed(() => {
  const value = Number(data.value.backgroundOverlay)
  if (!Number.isFinite(value)) return 0.65
  return Math.min(0.9, Math.max(0, value / 100))
})
const title = computed(() => titleHtml(data.value.title || ''))
</script>

<style scoped>
.cta-block {
  padding-block: clamp(50px, 7vw, 96px);
  background: var(--paper);
}

.cta-panel {
  position: relative;
  border-radius: clamp(22px, 3vw, 40px);
  overflow: hidden;
  padding: clamp(40px, 6vw, 84px);
  min-height: 380px;
  display: flex;
  align-items: center;
}

.cta-block--accent .cta-panel {
  background: var(--accent);
  color: var(--paper);
}

.cta-block--dark .cta-panel {
  background: var(--ink);
  color: var(--paper);
}

.cta-block--paper .cta-panel {
  background: var(--paper-deep);
  color: var(--ink);
  border: 1px solid var(--line);
}

.cta-bg {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.cta-bg-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    100deg,
    rgba(20, 21, 24, 0.95) 0%,
    rgba(20, 21, 24, 0.72) 45%,
    rgba(20, 21, 24, 0.42) 75%,
    rgba(20, 21, 24, 0.28) 100%
  );
}

.cta-block--image .cta-panel {
  color: var(--paper);
}

.cta-block--image .cta-content,
.cta-block--image .cta-side {
  text-shadow: 0 1px 18px rgba(0, 0, 0, 0.35);
}

.cta-rings {
  position: absolute;
  right: -8%;
  top: 50%;
  transform: translateY(-50%);
  width: 460px;
  height: 460px;
  pointer-events: none;
}

.cta-rings i {
  position: absolute;
  border-radius: 50%;
  border: 1px solid currentColor;
  opacity: 0.16;
}

.cta-rings i:nth-child(1) {
  inset: 0;
}

.cta-rings i:nth-child(2) {
  inset: 17%;
}

.cta-rings i:nth-child(3) {
  inset: 34%;
  background: currentColor;
  opacity: 0.05;
}

.cta-glyph {
  position: absolute;
  right: 4%;
  top: 50%;
  transform: translateY(-54%);
  font-family: var(--font-serif);
  font-size: clamp(180px, 25vw, 360px);
  line-height: 1;
  opacity: 0.09;
  user-select: none;
  pointer-events: none;
}

.cta-content {
  position: relative;
  z-index: 2;
  max-width: 720px;
}

.cta-kicker {
  display: flex;
  align-items: center;
  gap: 9px;
  opacity: 0.7;
}

.cta-kicker .mono {
  font-size: 9px;
  letter-spacing: 0.18em;
}

.cta-kicker .dot {
  background: currentColor;
}

.cta-content h2 {
  margin-top: 20px;
  font-size: clamp(34px, 4.8vw, 68px);
  line-height: 1.04;
  letter-spacing: -0.045em;
  max-width: 100%;
}

.cta-text {
  margin-top: 20px;
  max-width: 46ch;
  font-size: 14.5px;
  line-height: 1.75;
  opacity: 0.72;
}

.cta-actions {
  margin-top: 32px;
}

.cta-block--dark .cta-actions .btn--light {
  --btn-bg: var(--accent);
  --btn-fg: #fff;
}

.cta-block--paper .cta-actions .btn--light {
  --btn-bg: var(--ink);
  --btn-fg: var(--paper);
}

.cta-block--paper .cta-actions .btn--outline-light {
  --btn-fg: var(--ink);
  border-color: var(--line-strong);
}

.cta-block--paper .cta-actions .btn--outline-light:hover {
  --btn-bg: var(--ink);
  --btn-fg: var(--paper);
}

.cta-block--image .cta-actions .btn--light {
  --btn-bg: var(--accent);
  --btn-fg: #fff;
}

.cta-block--image .cta-actions .btn--outline-light {
  --btn-fg: var(--paper);
  border-color: rgba(246, 242, 234, 0.4);
}

.cta-block--image .cta-actions .btn--outline-light:hover {
  --btn-bg: rgba(246, 242, 234, 0.14);
  --btn-fg: var(--paper);
  border-color: rgba(246, 242, 234, 0.7);
}

.cta-side {
  position: absolute;
  right: clamp(28px, 4vw, 56px);
  bottom: clamp(28px, 4vw, 48px);
  display: flex;
  align-items: flex-end;
  gap: 12px;
  opacity: 0.55;
  z-index: 2;
}

.cta-side-line {
  width: 44px;
  height: 1px;
  background: currentColor;
  margin-bottom: 7px;
}

.cta-side p {
  font-size: 8.5px;
  line-height: 1.7;
  letter-spacing: 0.14em;
  text-align: right;
}

.cta-side-note {
  white-space: pre-line;
}

@container canvas (max-width: 900px) {
  .cta-panel {
    padding: clamp(34px, 6vw, 56px);
    min-height: 0;
  }
  .cta-glyph {
    font-size: 240px;
    right: -4%;
    opacity: 0.06;
  }
  .cta-side {
    display: none;
  }
}

@container canvas (max-width: 620px) {
  .cta-content h2 {
    font-size: clamp(31px, 9.6vw, 42px);
  }
  .cta-actions {
    flex-direction: column;
    align-items: stretch;
  }
  .cta-actions .btn {
    width: 100%;
  }
  .cta-rings {
    right: -40%;
    opacity: 0.7;
  }
  .cta-glyph {
    font-size: 180px;
    right: -10%;
  }
}
</style>

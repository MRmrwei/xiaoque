<template>
  <section class="image-text-block section">
    <div class="container it-grid" :class="{ 'it-grid--reverse': data.imageSide === 'left' }">
      <div class="it-copy">
        <p v-if="data.eyebrow" class="eyebrow">{{ data.eyebrow }}</p>
        <h2 v-html="title" />
        <p v-if="data.body" class="it-body">{{ data.body }}</p>
        <ul v-if="data.bullets?.length" class="it-list">
          <li v-for="(bullet, index) in data.bullets" :key="index">
            <span class="it-check"><UiIcon name="check" :size="12" /></span>
            {{ bullet.text }}
          </li>
        </ul>
        <BlockLink
          v-if="data.ctaLabel"
          :href="data.ctaHref"
          :preview="preview"
          class="it-cta"
        >
          {{ data.ctaLabel }} <span class="arr">↗</span>
        </BlockLink>
      </div>

      <div class="it-visual">
        <div class="it-frame">
          <img v-if="customImage" class="it-img" :src="data.image" alt="" />
          <VisualArt v-else :variant="artVariant" />
        </div>
        <span v-if="data.badge" class="it-badge mono">{{ data.badge }}</span>
        <span class="it-cross it-cross--tl" aria-hidden="true" />
        <span class="it-cross it-cross--br" aria-hidden="true" />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BlockLink from '~/components/BlockLink.vue'
import UiIcon from '~/components/UiIcon.vue'
import VisualArt from '~/components/VisualArt.vue'
import type { ImageTextProps, VisualKey } from '~/domain/blocks'
import { titleHtml } from '~/utils/format'
import { isCustomImage } from '@/utils/image'

const props = withDefaults(
  defineProps<{
    propsData: ImageTextProps
    preview?: boolean
  }>(),
  { preview: false }
)

const data = computed(() => props.propsData)
const title = computed(() => titleHtml(data.value.title || ''))
const customImage = computed(() => isCustomImage(data.value.image))
const artVariant = computed<VisualKey>(() => {
  const image = data.value.image
  return customImage.value || !image ? 'orbit' : (image as VisualKey)
})
</script>

<style scoped>
.image-text-block {
  background: var(--paper);
}

.it-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 0.92fr);
  align-items: center;
  gap: clamp(40px, 7vw, 100px);
}

.it-grid--reverse .it-copy {
  order: 2;
}

.it-grid--reverse .it-visual {
  order: 1;
}

.it-copy h2 {
  margin-top: 18px;
  font-size: clamp(30px, 3.4vw, 46px);
  max-width: 100%;
}

.it-body {
  margin-top: 20px;
  color: var(--muted);
  font-size: 15px;
  line-height: 1.8;
  max-width: 46ch;
}

.it-list {
  margin-top: 26px;
  border-top: 1px solid var(--line);
  max-width: 48ch;
}

.it-list li {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding-block: 13px;
  border-bottom: 1px solid var(--line-soft);
  font-size: 14px;
  color: var(--ink-soft);
}

.it-check {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--accent-soft);
  color: var(--accent-deep);
  display: grid;
  place-items: center;
  flex: none;
  margin-top: 1px;
}

.it-cta {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  margin-top: 28px;
  font-size: 13.5px;
  font-weight: 700;
  color: var(--ink);
  border-bottom: 1px solid var(--ink);
  padding-bottom: 3px;
  transition: color 0.2s ease, border-color 0.2s ease;
}

.it-cta:hover {
  color: var(--accent);
  border-color: var(--accent);
}

.it-visual {
  position: relative;
}

.it-frame {
  aspect-ratio: 4 / 4.4;
  border-radius: clamp(20px, 2.6vw, 32px);
  overflow: hidden;
  border: 1px solid var(--line-soft);
  box-shadow: var(--shadow-md);
  background: var(--paper-deep);
}

.it-img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.it-badge {
  position: absolute;
  left: 18px;
  bottom: 18px;
  font-size: 9px;
  letter-spacing: 0.16em;
  color: var(--paper);
  background: rgba(20, 21, 24, 0.62);
  backdrop-filter: blur(8px);
  padding: 8px 12px;
  border-radius: 999px;
  border: 1px solid rgba(246, 242, 234, 0.15);
}

.it-cross {
  position: absolute;
  width: 18px;
  height: 18px;
  pointer-events: none;
}

.it-cross::before,
.it-cross::after {
  content: '';
  position: absolute;
  background: var(--ink);
  opacity: 0.25;
}

.it-cross::before {
  left: 50%;
  top: 0;
  width: 1px;
  height: 100%;
}

.it-cross::after {
  top: 50%;
  left: 0;
  height: 1px;
  width: 100%;
}

.it-cross--tl { top: -28px; left: -28px; }
.it-cross--br { bottom: -28px; right: -28px; }

@container canvas (max-width: 900px) {
  .it-grid {
    grid-template-columns: 1fr;
    gap: 42px;
  }
  .it-grid--reverse .it-copy {
    order: 1;
  }
  .it-grid--reverse .it-visual {
    order: 2;
  }
  .it-frame {
    aspect-ratio: 16 / 11;
  }
  .it-cross {
    display: none;
  }
  .it-copy h2 {
    max-width: 100%;
  }
}

@container canvas (max-width: 620px) {
  .it-frame {
    aspect-ratio: 4 / 3.4;
  }
  .it-copy h2 {
    font-size: clamp(28px, 8.4vw, 38px);
  }
}
</style>

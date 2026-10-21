<template>
  <section class="video-block section">
    <div class="container">
      <header v-if="data.title" class="video-head">
        <h2>{{ data.title }}</h2>
        <span v-if="data.caption" class="mono video-caption">{{ data.caption }}</span>
      </header>

      <component
        :is="data.url && !preview ? 'a' : 'div'"
        :href="data.url && !preview ? data.url : undefined"
        :target="data.url && !preview ? '_blank' : undefined"
        class="video-frame"
      >
        <BlockImage :src="data.visual" fallback="bloom" />
        <span class="video-veil" aria-hidden="true" />
        <span class="play-button">
          <UiIcon name="play" :size="26" />
        </span>
        <span v-if="data.label" class="video-label mono">{{ data.label }}</span>
      </component>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BlockImage from '~/components/BlockImage.vue'
import UiIcon from '~/components/UiIcon.vue'
import type { VideoProps } from '~/domain/blocks'

const props = withDefaults(
  defineProps<{
    propsData: VideoProps
    preview?: boolean
  }>(),
  { preview: false }
)

const data = computed(() => props.propsData)
</script>

<style scoped>
.video-block {
  background: var(--paper);
}

.video-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 28px;
}

.video-head h2 {
  font-size: clamp(28px, 3.4vw, 44px);
  max-width: 18ch;
}

.video-caption {
  font-size: 9px;
  letter-spacing: 0.14em;
  color: var(--muted-2);
  white-space: nowrap;
}

.video-frame {
  position: relative;
  display: block;
  aspect-ratio: 16 / 8.4;
  border-radius: clamp(18px, 2.4vw, 30px);
  overflow: hidden;
  border: 1px solid var(--line-soft);
  max-width: 1080px;
  margin-inline: auto;
}

.video-veil {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(20, 21, 24, 0.5), rgba(20, 21, 24, 0.08) 55%, transparent);
  transition: background 0.3s ease;
}

.video-frame:hover .video-veil {
  background: linear-gradient(to top, rgba(20, 21, 24, 0.58), rgba(20, 21, 24, 0.16) 55%, transparent);
}

.play-button {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: clamp(64px, 8vw, 92px);
  height: clamp(64px, 8vw, 92px);
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: rgba(246, 242, 234, 0.92);
  color: var(--ink);
  backdrop-filter: blur(10px);
  transition: transform 0.4s var(--ease), background 0.2s ease, color 0.2s ease;
}

.video-frame:hover .play-button {
  transform: translate(-50%, -50%) scale(1.08);
  background: var(--accent);
  color: #fff;
}

.video-label {
  position: absolute;
  left: 22px;
  bottom: 18px;
  font-size: 8.5px;
  letter-spacing: 0.16em;
  color: rgba(246, 242, 234, 0.75);
}

@container canvas (max-width: 620px) {
  .video-frame {
    aspect-ratio: 4 / 3;
  }
  .video-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>

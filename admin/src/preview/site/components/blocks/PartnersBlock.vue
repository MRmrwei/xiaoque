<template>
  <section class="partners-block">
    <div class="container">
      <div class="partners-head">
        <p class="partners-title">{{ data.title }}</p>
        <p v-if="data.note" class="partners-note mono">{{ data.note }}</p>
      </div>
    </div>

    <div class="marquee" :class="{ 'marquee--static': !data.marquee }">
      <div class="marquee-track" :class="{ 'is-marquee': data.marquee }">
        <div v-for="pass in passes" :key="pass" class="marquee-group" :aria-hidden="pass > 1">
          <span
            v-for="(item, index) in data.items"
            :key="`${pass}-${index}`"
            class="partner"
            :class="`partner--${index % 3}`"
          >
            <span class="partner-mark" />
            {{ item.name }}
          </span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PartnersProps } from '~/domain/blocks'

const props = defineProps<{
  propsData: PartnersProps
}>()

const data = computed(() => props.propsData)
const passes = computed(() => (data.value.marquee ? [1, 2] : [1]))
</script>

<style scoped>
.partners-block {
  padding-block: clamp(34px, 4vw, 50px) clamp(40px, 4.5vw, 56px);
  border-bottom: 1px solid var(--line-soft);
  background: var(--paper);
  overflow: hidden;
}

.partners-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 28px;
}

.partners-title {
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.02em;
  color: var(--muted);
}

.partners-note {
  font-size: 9.5px;
  letter-spacing: 0.1em;
  color: var(--muted-2);
  text-transform: uppercase;
}

.marquee {
  overflow: hidden;
  mask-image: linear-gradient(to right, transparent, #000 9%, #000 91%, transparent);
  -webkit-mask-image: linear-gradient(to right, transparent, #000 9%, #000 91%, transparent);
}

.marquee-track {
  display: flex;
  width: max-content;
}

.marquee-track.is-marquee {
  animation: marquee 34s linear infinite;
}

.marquee:hover .marquee-track {
  animation-play-state: paused;
}

.marquee-group {
  display: flex;
  align-items: center;
  flex: none;
}

.partner {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding-inline: clamp(22px, 3vw, 46px);
  font-size: clamp(17px, 1.8vw, 23px);
  color: var(--ink);
  white-space: nowrap;
  opacity: 0.52;
  transition: opacity 0.3s ease, color 0.3s ease;
}

.partner:hover {
  opacity: 1;
  color: var(--accent);
}

.partner-mark {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--accent);
  opacity: 0.6;
}

.partner--0 {
  font-family: var(--font-serif);
  letter-spacing: -0.02em;
}

.partner--1 {
  font-weight: 700;
  letter-spacing: -0.045em;
}

.partner--2 {
  font-family: var(--font-mono);
  font-size: clamp(13px, 1.3vw, 16px);
  letter-spacing: 0.06em;
}

.marquee--static .marquee-track {
  width: 100%;
  justify-content: center;
}

.marquee--static .marquee-group {
  flex-wrap: wrap;
  justify-content: center;
}

@keyframes marquee {
  from { transform: translateX(0); }
  to { transform: translateX(-50%); }
}

@container canvas (max-width: 700px) {
  .partners-head {
    flex-direction: column;
    gap: 6px;
  }
}
</style>

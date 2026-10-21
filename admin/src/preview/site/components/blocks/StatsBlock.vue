<template>
  <section class="stats-block" :class="`stats-block--${data.theme || 'dark'}`">
    <div class="stats-grid-lines" aria-hidden="true" />
    <div class="container">
      <div class="stats-head">
        <span v-if="kicker" class="mono stats-kicker">{{ kicker }}</span>
        <span v-if="data.note" class="stats-note">{{ data.note }}</span>
      </div>
      <div class="stats-grid">
        <div v-for="(item, index) in data.items" :key="index" class="stat">
          <span class="stat-index mono">{{ String(index + 1).padStart(2, '0') }}</span>
          <p class="stat-value">
            {{ item.value }}<span v-if="item.suffix" class="stat-suffix">{{ item.suffix }}</span>
          </p>
          <p class="stat-label">{{ item.label }}</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { StatsProps } from '~/domain/blocks'

const props = defineProps<{
  propsData: StatsProps
}>()

const data = computed(() => props.propsData)
const kicker = computed(() =>
  data.value.kicker === undefined ? 'BY THE NUMBERS' : data.value.kicker
)
</script>

<style scoped>
.stats-block {
  position: relative;
  overflow: hidden;
  padding-block: clamp(52px, 6vw, 80px);
}

.stats-block--dark {
  background: var(--ink);
  color: var(--paper);
}

.stats-block--paper {
  background: var(--paper-deep);
  color: var(--ink);
  border-block: 1px solid var(--line);
}

.stats-grid-lines {
  position: absolute;
  inset: 0;
  background-image: linear-gradient(to right, rgba(246, 242, 234, 0.05) 1px, transparent 1px);
  background-size: 25% 100%;
  pointer-events: none;
  opacity: 0.6;
}

.stats-block--paper .stats-grid-lines {
  background-image: linear-gradient(to right, rgba(20, 21, 24, 0.06) 1px, transparent 1px);
}

.stats-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: clamp(26px, 3vw, 40px);
  position: relative;
  z-index: 1;
}

.stats-kicker {
  font-size: 9px;
  letter-spacing: 0.2em;
  opacity: 0.4;
}

.stats-note {
  font-size: 11px;
  opacity: 0.4;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  position: relative;
  z-index: 1;
}

.stat {
  position: relative;
  padding: 8px clamp(14px, 2.4vw, 34px) 10px;
  border-left: 1px solid rgba(246, 242, 234, 0.12);
}

.stats-block--paper .stat {
  border-left-color: rgba(20, 21, 24, 0.13);
}

.stat:first-child {
  border-left: 0;
  padding-left: 0;
}

.stat-index {
  position: absolute;
  top: 0;
  right: clamp(14px, 2.4vw, 34px);
  font-size: 8px;
  letter-spacing: 0.12em;
  opacity: 0.3;
}

.stat-value {
  font-family: var(--font-serif);
  font-size: clamp(44px, 5.4vw, 76px);
  line-height: 1;
  letter-spacing: -0.04em;
  font-weight: 500;
  display: flex;
  align-items: baseline;
}

.stat-suffix {
  font-family: var(--font-sans);
  font-size: 0.42em;
  font-weight: 600;
  color: var(--accent);
  margin-left: 3px;
  letter-spacing: -0.02em;
}

.stats-block--paper .stat-suffix {
  color: var(--accent-deep);
}

.stat-label {
  margin-top: 12px;
  font-size: 12px;
  letter-spacing: 0.02em;
  opacity: 0.55;
}

@container canvas (max-width: 800px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 26px 0;
  }
  .stat:nth-child(3) {
    border-left: 0;
    padding-left: 0;
  }
  .stat:nth-child(3),
  .stat:nth-child(4) {
    padding-top: 24px;
    border-top: 1px solid rgba(246, 242, 234, 0.12);
  }
  .stats-block--paper .stat:nth-child(3),
  .stats-block--paper .stat:nth-child(4) {
    border-top-color: rgba(20, 21, 24, 0.13);
  }
  .stats-head {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}

@container canvas (max-width: 460px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  .stat {
    padding-top: 20px;
    padding-left: 0;
    border-left: 0;
    border-top: 1px solid rgba(246, 242, 234, 0.12);
  }
  .stat:first-child {
    border-top: 0;
    padding-top: 8px;
  }
  .stats-block--paper .stat {
    border-top-color: rgba(20, 21, 24, 0.13);
  }
}
</style>

<template>
  <div class="builder-mock" aria-hidden="true">
    <div class="browser">
      <div class="browser-bar">
        <div class="browser-dots"><i /><i /><i /></div>
        <div class="browser-address">
          <span class="lock" />
          <span class="address-text mono">app.xiaoque.com/pages/home</span>
          <span class="live-chip mono">LIVE</span>
        </div>
        <div class="browser-tools"><i /><i /></div>
      </div>

      <div class="builder-body">
        <aside class="rail-left">
          <div class="rail-head">
            <span class="rail-dot" />
            <span class="rail-title">区块</span>
          </div>
          <div class="rail-item rail-item--active">
            <span class="rail-icon" />
            <span class="rail-bar" />
          </div>
          <div class="rail-item">
            <span class="rail-icon" />
            <span class="rail-bar rail-bar--short" />
          </div>
          <div class="rail-item">
            <span class="rail-icon" />
            <span class="rail-bar" />
          </div>
          <div class="rail-item">
            <span class="rail-icon" />
            <span class="rail-bar rail-bar--short" />
          </div>
          <div class="rail-item">
            <span class="rail-icon" />
            <span class="rail-bar" />
          </div>
          <div class="rail-add mono">+</div>
        </aside>

        <div class="canvas-area">
          <div class="canvas-page">
            <div class="mini-nav">
              <span class="mini-logo" />
              <span class="mini-links"><i /><i /><i /></span>
              <span class="mini-nav-btn" />
            </div>

            <div class="mini-hero">
              <span class="selection-chip mono">HERO</span>
              <span class="handle handle--tl" />
              <span class="handle handle--tr" />
              <span class="handle handle--bl" />
              <span class="handle handle--br" />
              <div class="mini-copy">
                <span class="mini-eyebrow" />
                <span class="mini-heading mini-heading--1" />
                <span class="mini-heading mini-heading--2" />
                <span class="mini-heading mini-heading--3" />
                <span class="mini-sub" />
                <span class="mini-sub mini-sub--short" />
                <span class="mini-buttons"><i class="mini-btn mini-btn--dark" /><i class="mini-btn mini-btn--ghost" /></span>
              </div>
              <div class="mini-art">
                <VisualArt variant="orbit" tone="accent" />
              </div>
            </div>

            <div class="mini-strip">
              <div class="mini-card">
                <span class="mini-card-art mini-card-art--a" />
                <span class="mini-line" />
                <span class="mini-line mini-line--short" />
              </div>
              <div class="mini-card">
                <span class="mini-card-art mini-card-art--b" />
                <span class="mini-line" />
                <span class="mini-line mini-line--short" />
              </div>
              <div class="mini-card">
                <span class="mini-card-art mini-card-art--c" />
                <span class="mini-line" />
                <span class="mini-line mini-line--short" />
              </div>
            </div>
          </div>
          <span class="canvas-label mono">CANVAS · 1440</span>
        </div>

        <aside class="rail-right">
          <div class="panel-title mono">属性</div>
          <div class="panel-group">
            <span class="panel-label" />
            <span class="panel-input" />
          </div>
          <div class="panel-group">
            <span class="panel-label panel-label--short" />
            <span class="panel-input panel-input--tall" />
          </div>
          <div class="panel-group">
            <span class="panel-label panel-label--short" />
            <span class="panel-swatches"><i class="swatch swatch--ink" /><i class="swatch swatch--red" /><i class="swatch swatch--cream" /></span>
          </div>
          <div class="panel-toggle">
            <span class="panel-label panel-label--short" />
            <span class="toggle"><i /></span>
          </div>
        </aside>
      </div>
    </div>

    <div class="float-card float-status">
      <span class="status-dot" />
      <span>{{ locale === 'en' ? 'Published' : '已发布' }}</span>
      <b class="mono">v12</b>
    </div>

    <div class="float-card float-lang">
      <span class="lang-opt is-active">中</span>
      <span class="lang-sep" />
      <span class="lang-opt">EN</span>
      <span class="sync-dot" />
      <em>{{ locale === 'en' ? 'Synced' : '已同步' }}</em>
    </div>

    <div class="float-seal">鹊</div>
  </div>
</template>

<script setup lang="ts">
import VisualArt from '~/components/VisualArt.vue'
import type { Locale } from '~/domain/types'

withDefaults(
  defineProps<{
    locale?: Locale
  }>(),
  { locale: 'zh' }
)
</script>

<style scoped>
.builder-mock {
  position: relative;
  width: 100%;
  padding: clamp(22px, 3vw, 34px) clamp(16px, 2.4vw, 28px);
  background:
    radial-gradient(circle at 82% 12%, rgba(232, 80, 58, 0.14), transparent 42%),
    linear-gradient(145deg, #eee7da 0%, #f6f2ea 55%, #efe9dd 100%);
}

.builder-mock::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: radial-gradient(rgba(20, 21, 24, 0.14) 0.7px, transparent 0.7px);
  background-size: 17px 17px;
  opacity: 0.35;
  mask-image: linear-gradient(135deg, #000, transparent 55%);
  pointer-events: none;
}

.browser {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  aspect-ratio: 1.17 / 1;
  border-radius: 13px;
  overflow: hidden;
  border: 1px solid rgba(20, 21, 24, 0.16);
  background: #fbf9f5;
  box-shadow: 0 2px 3px rgba(20, 21, 24, 0.06), 0 34px 70px -20px rgba(20, 21, 24, 0.4);
}

.browser-bar {
  height: clamp(24px, 4.6%, 34px);
  flex: none;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 10px;
  background: #fbf9f5;
  border-bottom: 1px solid rgba(20, 21, 24, 0.09);
}

.browser-dots {
  display: flex;
  gap: 3.5px;
}

.browser-dots i {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: rgba(20, 21, 24, 0.18);
}

.browser-dots i:first-child {
  background: #e8503a;
}

.browser-address {
  flex: 1;
  min-width: 0;
  height: 17px;
  border-radius: 999px;
  background: rgba(20, 21, 24, 0.045);
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 0 7px;
}

.lock {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #7ba273;
}

.address-text {
  font-size: clamp(4.5px, 0.72vw, 6.6px);
  letter-spacing: 0.02em;
  color: rgba(20, 21, 24, 0.45);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.live-chip {
  margin-left: auto;
  font-size: clamp(4px, 0.6vw, 5.5px);
  letter-spacing: 0.12em;
  color: #3f7a4f;
  background: rgba(93, 145, 103, 0.12);
  padding: 1px 4px;
  border-radius: 999px;
}

.browser-tools {
  display: flex;
  gap: 4px;
}

.browser-tools i {
  width: 7px;
  height: 7px;
  border: 1px solid rgba(20, 21, 24, 0.16);
  border-radius: 2px;
}

.builder-body {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: 13.5% minmax(0, 1fr) 16.5%;
}

.rail-left {
  background: #f3eee5;
  border-right: 1px solid rgba(20, 21, 24, 0.07);
  padding: 8px 6px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow: hidden;
}

.rail-head {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 1px 2px 6px;
}

.rail-dot {
  width: 5px;
  height: 5px;
  border-radius: 2px;
  background: var(--accent);
}

.rail-title {
  font-size: clamp(4.5px, 0.62vw, 6px);
  font-weight: 700;
  letter-spacing: 0.08em;
  color: rgba(20, 21, 24, 0.45);
}

.rail-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 4px;
  border-radius: 4px;
}

.rail-item--active {
  background: #fbf9f5;
  box-shadow: 0 1px 2px rgba(20, 21, 24, 0.08);
  position: relative;
}

.rail-item--active::before {
  content: '';
  position: absolute;
  left: -6px;
  top: 4px;
  bottom: 4px;
  width: 2px;
  background: var(--accent);
  border-radius: 2px;
}

.rail-icon {
  width: clamp(6px, 1vw, 9px);
  height: clamp(6px, 1vw, 9px);
  border-radius: 2px;
  background: rgba(20, 21, 24, 0.14);
  flex: none;
}

.rail-item--active .rail-icon {
  background: var(--accent);
}

.rail-bar {
  height: 2.5px;
  flex: 1;
  border-radius: 2px;
  background: rgba(20, 21, 24, 0.16);
}

.rail-bar--short {
  max-width: 62%;
}

.rail-add {
  margin-top: auto;
  height: 16px;
  border: 1px dashed rgba(20, 21, 24, 0.2);
  border-radius: 4px;
  display: grid;
  place-items: center;
  font-size: 9px;
  color: rgba(20, 21, 24, 0.3);
}

.canvas-area {
  position: relative;
  background: #eae4d9;
  background-image: radial-gradient(rgba(20, 21, 24, 0.14) 0.6px, transparent 0.6px);
  background-size: 11px 11px;
  padding: 12px 12px 18px;
  display: flex;
  min-height: 0;
}

.canvas-page {
  flex: 1;
  min-height: 0;
  border-radius: 6px;
  overflow: hidden;
  background: #fbf9f5;
  box-shadow: 0 6px 24px rgba(20, 21, 24, 0.13);
  display: flex;
  flex-direction: column;
}

.mini-nav {
  height: 8.5%;
  min-height: 18px;
  flex: none;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 9px;
  border-bottom: 1px solid rgba(20, 21, 24, 0.06);
}

.mini-logo {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--ink);
  position: relative;
  flex: none;
}

.mini-logo::after {
  content: '';
  position: absolute;
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--accent);
  right: -1px;
  top: -1px;
}

.mini-links {
  display: flex;
  gap: 5px;
  margin-left: 4px;
}

.mini-links i {
  width: 12px;
  height: 1.6px;
  border-radius: 2px;
  background: rgba(20, 21, 24, 0.16);
}

.mini-nav-btn {
  margin-left: auto;
  width: 24px;
  height: 8px;
  border-radius: 999px;
  background: var(--ink);
}

.mini-hero {
  position: relative;
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: 1.18fr 0.82fr;
  align-items: center;
  gap: 8px;
  padding: 8% 7% 6%;
  background: linear-gradient(115deg, rgba(232, 80, 58, 0.06), transparent 55%);
  outline: 1px solid rgba(232, 80, 58, 0.55);
  outline-offset: -2px;
}

.selection-chip {
  position: absolute;
  top: 2px;
  left: 3px;
  font-size: clamp(4px, 0.55vw, 5.2px);
  letter-spacing: 0.12em;
  color: #fff;
  background: var(--accent);
  padding: 1.5px 4px;
  border-radius: 0 0 3px 0;
}

.handle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: #fff;
  border: 1px solid var(--accent);
}

.handle--tl { top: -2px; left: -2px; }
.handle--tr { top: -2px; right: -2px; }
.handle--bl { bottom: -2px; left: -2px; }
.handle--br { bottom: -2px; right: -2px; }

.mini-copy {
  display: flex;
  flex-direction: column;
  gap: clamp(3px, 0.6vw, 5px);
  min-width: 0;
}

.mini-eyebrow {
  width: 32%;
  height: 2px;
  border-radius: 2px;
  background: rgba(232, 80, 58, 0.55);
  margin-bottom: 2px;
}

.mini-heading {
  display: block;
  height: clamp(5px, 1.1vw, 9px);
  border-radius: 2px;
  background: var(--ink);
}

.mini-heading--1 { width: 94%; }
.mini-heading--2 { width: 78%; background: var(--accent); }
.mini-heading--3 { width: 56%; }

.mini-sub {
  width: 84%;
  height: 2px;
  border-radius: 2px;
  background: rgba(20, 21, 24, 0.2);
  margin-top: 2px;
}

.mini-sub--short { width: 62%; }

.mini-buttons {
  display: flex;
  gap: 4px;
  margin-top: 4px;
}

.mini-btn {
  display: block;
  width: 30px;
  height: 9px;
  border-radius: 999px;
}

.mini-btn--dark { background: var(--ink); }
.mini-btn--ghost { border: 0.7px solid rgba(20, 21, 24, 0.25); }

.mini-art {
  height: 100%;
  max-height: 100%;
  border-radius: 4px;
  overflow: hidden;
  transform: rotate(2deg);
  box-shadow: 0 5px 16px rgba(20, 21, 24, 0.16);
  min-height: 0;
}

.mini-strip {
  flex: none;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 5px;
  padding: 0 7% 6%;
  height: 20%;
  min-height: 40px;
}

.mini-card {
  border: 1px solid rgba(20, 21, 24, 0.09);
  border-radius: 4px;
  padding: 4px;
  display: flex;
  flex-direction: column;
  gap: 3px;
  background: #fbf9f5;
}

.mini-card-art {
  flex: 1;
  min-height: 10px;
  border-radius: 2px;
}

.mini-card-art--a { background: linear-gradient(120deg, #e8503a, #f3b5a7); }
.mini-card-art--b { background: linear-gradient(120deg, #141518, #3f574c); }
.mini-card-art--c { background: linear-gradient(120deg, #d8cdb8, #efe6d6); }

.mini-line {
  height: 2px;
  border-radius: 2px;
  background: rgba(20, 21, 24, 0.2);
}

.mini-line--short { width: 60%; background: rgba(20, 21, 24, 0.12); }

.canvas-label {
  position: absolute;
  left: 12px;
  bottom: 5px;
  font-size: clamp(4px, 0.5vw, 5px);
  letter-spacing: 0.1em;
  color: rgba(20, 21, 24, 0.35);
}

.rail-right {
  background: #fbf9f5;
  border-left: 1px solid rgba(20, 21, 24, 0.07);
  padding: 8px 7px;
  display: flex;
  flex-direction: column;
  gap: 9px;
  overflow: hidden;
}

.panel-title {
  font-size: clamp(4.5px, 0.62vw, 6px);
  letter-spacing: 0.12em;
  color: rgba(20, 21, 24, 0.4);
  padding-bottom: 4px;
  border-bottom: 1px solid rgba(20, 21, 24, 0.07);
}

.panel-group {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.panel-label {
  width: 66%;
  height: 2.5px;
  border-radius: 2px;
  background: rgba(20, 21, 24, 0.18);
}

.panel-label--short { width: 42%; }

.panel-input {
  height: 12px;
  border-radius: 3px;
  border: 0.7px solid rgba(20, 21, 24, 0.12);
  background: #fff;
}

.panel-input--tall { height: 22px; }

.panel-swatches {
  display: flex;
  gap: 3px;
}

.swatch {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  border: 1px solid rgba(20, 21, 24, 0.1);
}

.swatch--ink { background: var(--ink); }
.swatch--red { background: var(--accent); }
.swatch--cream { background: #efe6d6; }

.panel-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 2px;
}

.toggle {
  width: 15px;
  height: 8px;
  border-radius: 999px;
  background: var(--accent);
  position: relative;
}

.toggle i {
  position: absolute;
  right: 1.5px;
  top: 1.5px;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #fff;
}

.float-card {
  position: absolute;
  z-index: 3;
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(251, 249, 245, 0.94);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(20, 21, 24, 0.1);
  border-radius: 999px;
  box-shadow: 0 10px 30px rgba(20, 21, 24, 0.16);
  font-size: clamp(7px, 0.78vw, 9px);
  color: var(--ink-soft);
  white-space: nowrap;
  padding: 5px 8px;
}

.float-status {
  top: 9px;
  right: 8px;
  animation: floaty 7s ease-in-out infinite;
}

.float-status b {
  font-weight: 500;
  color: var(--muted);
  font-size: 0.9em;
}

.status-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #6aa374;
  box-shadow: 0 0 0 2.5px rgba(106, 163, 116, 0.18);
}

.float-lang {
  bottom: 8px;
  left: 8px;
  gap: 4px;
  animation: floaty 8s ease-in-out 0.8s infinite;
}

.lang-opt {
  font-size: 0.92em;
  color: var(--muted-2);
  font-weight: 600;
}

.lang-opt.is-active {
  color: var(--ink);
}

.lang-sep {
  width: 1px;
  height: 8px;
  background: var(--line);
}

.sync-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #6aa374;
  margin-left: 2px;
}

.float-lang em {
  font-style: normal;
  color: var(--muted);
  font-size: 0.9em;
}

.float-seal {
  position: absolute;
  z-index: 3;
  right: 16px;
  bottom: 7px;
  width: clamp(22px, 3vw, 34px);
  height: clamp(22px, 3vw, 34px);
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--accent);
  color: #fff;
  font-family: var(--font-serif);
  font-size: clamp(11px, 1.5vw, 17px);
  box-shadow: 0 8px 20px rgba(201, 59, 38, 0.35);
  transform: rotate(8deg);
}

@keyframes floaty {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-3px); }
}

@media (prefers-reduced-motion: reduce) {
  .float-status,
  .float-lang {
    animation: none;
  }
}
</style>

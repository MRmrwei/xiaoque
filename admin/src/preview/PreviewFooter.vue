<template>
  <footer class="pv-footer">
    <div class="pv-footer-inner">
      <div class="pv-footer-top">
        <div class="pv-footer-brand">
          <div class="pv-footer-logo">
            <BrandMark :size="38" inverse :image="settings.brand.logo" />
            <span class="pv-logo-text">
              <strong>{{ settings.siteName[locale] }}</strong>
              <span>{{ settings.siteName.en }}</span>
            </span>
          </div>
          <p class="pv-footer-tagline">{{ settings.tagline[locale] }}</p>
        </div>

        <div class="pv-footer-cols" :class="{ 'is-contactless': !footer.showContact }">
          <div v-for="column in visibleColumns" :key="column.key" class="pv-footer-col">
            <h4>{{ column.title[locale] }}</h4>
            <span v-for="link in visibleLinks(column)" :key="link.key" :title="link.path">
              {{ link[locale] }}
            </span>
          </div>

          <div v-if="footer.showContact" class="pv-footer-col pv-footer-col--contact">
            <h4>{{ locale === 'en' ? 'Contact' : '联系' }}</h4>
            <span class="pv-line">
              <UiIcon name="mail" :size="14" /> {{ settings.email }}
            </span>
            <span class="pv-line">
              <UiIcon name="phone" :size="14" />
              {{ locale === 'en' ? settings.phoneEn : settings.phone }}
            </span>
            <span class="pv-line">
              <UiIcon name="pin" :size="14" /> {{ settings.address[locale] }}
            </span>
            <span class="pv-contact-cta" :title="footer.cta.path">
              {{ footer.cta.label[locale] }} <UiIcon name="arrow-up-right" :size="14" />
            </span>
          </div>
        </div>
      </div>

      <div v-if="footer.showWordmark" class="pv-wordmark" aria-hidden="true">
        <span>{{ footer.wordmark }}</span>
        <span class="pv-wordmark-dot">.</span>
      </div>

      <div class="pv-footer-bottom">
        <p>{{ settings.copyright[locale] }} <span class="pv-sep">·</span> {{ settings.icp }}</p>
        <div class="pv-bottom-links">
          <span v-for="link in visibleBottomLinks" :key="link.key" :title="link.path">
            {{ link[locale] }}
          </span>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { FooterColumn, FooterLink, Locale, SiteSettings } from '~/domain/types'
import BrandMark from '~/components/BrandMark.vue'
import UiIcon from '~/components/UiIcon.vue'

const props = defineProps<{
  settings: SiteSettings
  locale: Locale
}>()

const footer = computed(() => props.settings.footer)
const visibleColumns = computed(() => footer.value.columns.filter((column) => column.visible))
const visibleBottomLinks = computed(() =>
  footer.value.bottomLinks.filter((link) => link.visible)
)

function visibleLinks(column: FooterColumn): FooterLink[] {
  return column.links.filter((link) => link.visible)
}
</script>

<style scoped>
.pv-footer {
  position: relative;
  overflow: hidden;
  container-type: inline-size;
  background: var(--ink);
  color: rgba(246, 242, 234, 0.72);
}

.pv-footer::before {
  content: '';
  position: absolute;
  width: 520px;
  height: 520px;
  right: -200px;
  top: -240px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(232, 80, 58, 0.16), transparent 65%);
  pointer-events: none;
}

.pv-footer-inner {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 1240px;
  margin: 0 auto;
  padding-top: clamp(40px, 7cqw, 84px);
  padding-inline: clamp(16px, 3cqw, 26px);
}

.pv-footer-top {
  display: grid;
  grid-template-columns: 1.1fr 2fr;
  gap: clamp(32px, 6cqw, 90px);
}

.pv-footer-logo {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.pv-logo-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
  line-height: 1;
}

.pv-logo-text strong {
  color: #f6f2ea;
  font-size: 18px;
  letter-spacing: -0.03em;
  word-spacing: 0.06em;
}

.pv-logo-text span {
  font-family: var(--font-mono);
  font-size: 7.5px;
  letter-spacing: 0.2em;
  color: rgba(246, 242, 234, 0.4);
}

.pv-footer-tagline {
  margin-top: 20px;
  font-size: 13.5px;
  line-height: 1.8;
  max-width: 36ch;
  color: rgba(246, 242, 234, 0.55);
}

.pv-footer-cols {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 22px;
}

.pv-footer-cols.is-contactless {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.pv-footer-col {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 0;
}

.pv-footer-col h4 {
  margin: 0 0 6px;
  font-family: var(--font-mono);
  font-size: 9.5px;
  font-weight: 500;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: rgba(246, 242, 234, 0.36);
}

.pv-footer-col > span {
  display: inline-flex;
  align-items: flex-start;
  gap: 7px;
  font-size: 12.5px;
  line-height: 1.5;
  color: rgba(246, 242, 234, 0.66);
}

.pv-footer-col .ui-icon {
  margin-top: 2px;
  opacity: 0.5;
}

.pv-footer-col--contact .pv-contact-cta {
  margin-top: 6px;
  color: #f6f2ea;
  font-weight: 600;
  gap: 5px;
}

.pv-wordmark {
  margin-top: clamp(40px, 7cqw, 84px);
  display: flex;
  align-items: baseline;
  font-family: var(--font-serif, 'Instrument Serif', serif);
  font-size: clamp(56px, 16cqw, 200px);
  line-height: 0.72;
  letter-spacing: -0.04em;
  color: rgba(246, 242, 234, 0.09);
  user-select: none;
  white-space: nowrap;
  border-bottom: 1px solid rgba(246, 242, 234, 0.1);
  padding-bottom: 18px;
}

.pv-wordmark-dot {
  color: rgba(232, 80, 58, 0.55);
}

.pv-footer-bottom {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px 24px;
  padding-block: 22px 28px;
  font-size: 11px;
  color: rgba(246, 242, 234, 0.38);
}

.pv-sep {
  margin-inline: 5px;
  opacity: 0.5;
}

.pv-bottom-links {
  display: flex;
  gap: 20px;
}

.pv-bottom-links span {
  color: rgba(246, 242, 234, 0.48);
}

@container (max-width: 1020px) {
  .pv-footer-top {
    grid-template-columns: 1fr;
  }
  .pv-footer-cols,
  .pv-footer-cols.is-contactless {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 30px 22px;
  }
}

@container (max-width: 620px) {
  .pv-footer-cols,
  .pv-footer-cols.is-contactless {
    grid-template-columns: 1fr 1fr;
  }
  .pv-footer-bottom {
    flex-direction: column;
    align-items: flex-start;
  }
  .pv-wordmark {
    font-size: 18cqw;
  }
}
</style>

<template>
  <footer class="site-footer">
    <div class="container">
      <div class="footer-top">
        <div class="footer-brand">
          <NuxtLink :to="localePath('/')" class="footer-logo">
            <BrandMark :size="38" inverse :image="settings.brand.logo" />
            <span class="footer-logo-text">
              <strong>{{ settings.siteName[locale] }}</strong>
              <span>{{ settings.siteName.en }}</span>
            </span>
          </NuxtLink>
          <p class="footer-tagline">{{ settings.tagline[locale] }}</p>
        </div>

        <div class="footer-cols" :class="{ 'is-contactless': !footer.showContact }">
          <div v-for="column in visibleColumns" :key="column.key" class="footer-col">
            <h4>{{ column.title[locale] }}</h4>
            <BlockLink
              v-for="link in visibleLinks(column)"
              :key="link.key"
              :href="link.path"
            >
              {{ link[locale] }}
            </BlockLink>
          </div>
          <div v-if="footer.showContact" class="footer-col footer-col--contact">
            <h4>{{ t('footer.contact') }}</h4>
            <a :href="`mailto:${settings.email}`">
              <UiIcon name="mail" :size="14" /> {{ settings.email }}
            </a>
            <a :href="`tel:${locale === 'en' ? settings.phoneEn : settings.phone}`">
              <UiIcon name="phone" :size="14" /> {{ locale === 'en' ? settings.phoneEn : settings.phone }}
            </a>
            <span class="contact-line">
              <UiIcon name="pin" :size="14" /> {{ settings.address[locale] }}
            </span>
            <BlockLink :href="footer.cta.path" class="footer-cta">
              {{ footer.cta.label[locale] }} <UiIcon name="arrow-up-right" :size="14" />
            </BlockLink>
          </div>
        </div>
      </div>

      <div v-if="footer.showWordmark" class="footer-wordmark" aria-hidden="true">
        <span>{{ footer.wordmark }}</span>
        <span class="wordmark-dot">.</span>
      </div>

      <div class="footer-bottom">
        <p>{{ settings.copyright[locale] }} <span class="dot-sep">·</span> {{ settings.icp }}</p>
        <div class="footer-bottom-links">
          <BlockLink
            v-for="link in visibleBottomLinks"
            :key="link.key"
            :href="link.path"
          >
            {{ link[locale] }}
          </BlockLink>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BlockLink from '~/components/BlockLink.vue'
import BrandMark from '~/components/BrandMark.vue'
import UiIcon from '~/components/UiIcon.vue'
import { useSiteChrome } from '~/composables/useSiteChrome'
import { useUiText } from '~/composables/useUiText'
import type { FooterColumn, FooterLink } from '~/domain/types'

const localePath = useLocalePath()
const { t, locale } = useUiText()
const { settings } = useSiteChrome()

const footer = computed(() => settings.value.footer)
const visibleColumns = computed(() => footer.value.columns.filter((column) => column.visible))
const visibleBottomLinks = computed(() =>
  footer.value.bottomLinks.filter((link) => link.visible)
)

function visibleLinks(column: FooterColumn): FooterLink[] {
  return column.links.filter((link) => link.visible)
}
</script>

<style scoped>
.site-footer {
  background: var(--ink);
  color: rgba(246, 242, 234, 0.72);
  padding-top: clamp(56px, 7vw, 90px);
  overflow: hidden;
  position: relative;
}

.site-footer::before {
  content: '';
  position: absolute;
  width: 560px;
  height: 560px;
  right: -220px;
  top: -260px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(232, 80, 58, 0.16), transparent 65%);
  pointer-events: none;
}

.footer-top {
  display: grid;
  grid-template-columns: 1.1fr 2fr;
  gap: clamp(40px, 7vw, 110px);
  position: relative;
  z-index: 1;
}

.footer-logo {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.footer-logo-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
  line-height: 1;
}

.footer-logo-text strong {
  color: var(--paper);
  font-size: 18px;
  letter-spacing: -0.03em;
  word-spacing: 0.06em;
}

.footer-logo-text span {
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: rgba(246, 242, 234, 0.45);
}

.footer-tagline {
  margin-top: 22px;
  font-size: 14px;
  line-height: 1.8;
  max-width: 36ch;
  color: rgba(246, 242, 234, 0.55);
}

.footer-cols {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 24px;
}

.footer-cols.is-contactless {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.footer-col {
  display: flex;
  flex-direction: column;
  gap: 11px;
}

.footer-col h4 {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 500;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: rgba(246, 242, 234, 0.36);
  margin-bottom: 8px;
}

.footer-col a,
.footer-col span.contact-line {
  display: inline-flex;
  align-items: flex-start;
  gap: 7px;
  font-size: 13px;
  color: rgba(246, 242, 234, 0.66);
  transition: color 0.2s ease;
  line-height: 1.5;
}

.footer-col a:hover {
  color: var(--paper);
}

.footer-col .ui-icon {
  margin-top: 2px;
  opacity: 0.5;
}

.footer-col--contact .footer-cta {
  margin-top: 8px;
  color: var(--paper);
  font-weight: 600;
  gap: 5px;
}

.footer-col--contact .footer-cta:hover {
  color: #ff9a85;
}

.footer-wordmark {
  margin-top: clamp(56px, 8vw, 100px);
  display: flex;
  align-items: baseline;
  font-family: var(--font-serif);
  font-size: clamp(72px, 17vw, 230px);
  line-height: 0.72;
  letter-spacing: -0.04em;
  color: rgba(246, 242, 234, 0.09);
  user-select: none;
  white-space: nowrap;
  border-bottom: 1px solid rgba(246, 242, 234, 0.1);
  padding-bottom: 22px;
}

.wordmark-dot {
  color: rgba(232, 80, 58, 0.55);
}

.footer-bottom {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 16px 28px;
  padding-block: 26px 34px;
  font-size: 11.5px;
  color: rgba(246, 242, 234, 0.38);
}

.dot-sep {
  margin-inline: 6px;
  opacity: 0.5;
}

.footer-bottom-links {
  display: flex;
  gap: 22px;
}

.footer-bottom-links a {
  color: rgba(246, 242, 234, 0.48);
  transition: color 0.2s ease;
}

.footer-bottom-links a:hover {
  color: var(--paper);
}

@media (max-width: 1020px) {
  .footer-top {
    grid-template-columns: 1fr;
  }
  .footer-cols {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 34px 24px;
  }
}

@media (max-width: 620px) {
  .footer-cols {
    grid-template-columns: 1fr 1fr;
  }
  .footer-bottom {
    flex-direction: column;
    align-items: flex-start;
  }
  .footer-wordmark {
    font-size: 20vw;
  }
}
</style>

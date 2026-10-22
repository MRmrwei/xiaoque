<template>
  <div class="pv-chrome">
    <div v-if="showAnnouncement" class="pv-announcement">
      <div class="container pv-announcement-inner">
        <p class="pv-announcement-text">
          <span class="pv-announcement-dot" />
          <span>{{ announcement.text[locale] }}</span>
          <span class="pv-announcement-sep">—</span>
          <em>{{ announcement.highlight[locale] }}</em>
        </p>
        <span class="pv-announcement-link" :title="announcement.linkHref">
          {{ announcement.linkLabel[locale] }} ↗
        </span>
      </div>
    </div>

    <header class="pv-header" :class="{ 'pv-header--static': !sticky }">
      <div class="container pv-header-inner">
        <div class="pv-brand">
          <BrandMark :size="28" :image="settings.brand.logo" />
          <span class="pv-brand-text">
            <strong>{{ settings.siteName[locale] }}</strong>
            <span>{{ settings.siteName.en }}</span>
          </span>
        </div>
        <nav class="pv-nav">
          <span v-for="item in navItems" :key="item.key" :title="item.path">
            {{ item[locale] }}
          </span>
        </nav>
        <div class="pv-actions">
          <span v-if="showLanguage" class="pv-lang">
            <b>{{ locale === 'en' ? 'EN' : '中' }}</b><i>|</i>{{ locale === 'en' ? '中' : 'EN' }}
          </span>
          <span
            v-for="cta in visibleCtas"
            :key="cta.key"
            class="pv-cta"
            :class="`pv-cta--${cta.variant}`"
            :title="cta.href"
          >
            {{ cta.label[locale] }}
          </span>
        </div>
      </div>
    </header>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Locale, SiteSettings } from '~/domain/types'
import BrandMark from '~/components/BrandMark.vue'

const props = defineProps<{
  settings: SiteSettings
  locale: Locale
}>()

const navItems = computed(() => props.settings.nav.filter((item) => item.visible))
const header = computed(() => props.settings.header)
const sticky = computed(() => header.value?.sticky !== false)
const showLanguage = computed(() => header.value?.showLanguage !== false)
const visibleCtas = computed(() => (header.value?.ctas ?? []).filter((cta) => cta.visible))
const announcement = computed(() => header.value.announcement)
const showAnnouncement = computed(
  () =>
    header.value.showAnnouncement !== false &&
    Boolean(
      announcement.value?.text?.[props.locale] ||
        announcement.value?.highlight?.[props.locale] ||
        announcement.value?.linkLabel?.[props.locale]
    )
)
</script>

<style scoped>
.pv-announcement {
  background: var(--ink);
  color: rgba(246, 242, 234, 0.8);
}

.pv-announcement-inner {
  min-height: 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  font-size: 10px;
}

.pv-announcement-text {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pv-announcement-text em {
  font-style: normal;
  font-weight: 600;
  color: #f6f2ea;
}

.pv-announcement-dot {
  width: 5px;
  height: 5px;
  flex: none;
  border-radius: 50%;
  background: var(--accent);
}

.pv-announcement-sep {
  opacity: 0.35;
}

.pv-announcement-link {
  flex: none;
  color: #f6f2ea;
  font-weight: 600;
}

.pv-header {
  position: sticky;
  top: 0;
  z-index: 20;
  container-type: inline-size;
  background: color-mix(in srgb, var(--paper) 94%, transparent);
  border-bottom: 1px solid rgba(20, 21, 24, 0.08);
  backdrop-filter: blur(14px);
}

.pv-header--static {
  position: relative;
}

.pv-header-inner {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.pv-brand {
  display: flex;
  align-items: center;
  gap: 9px;
  flex: none;
}

.pv-brand-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  line-height: 1;
}

.pv-brand-text strong {
  font-size: 14px;
  letter-spacing: -0.03em;
}

.pv-brand-text span {
  font-family: var(--font-mono);
  font-size: 7px;
  letter-spacing: 0.15em;
  color: var(--muted-2);
}

.pv-nav {
  display: flex;
  align-items: center;
  gap: 24px;
  font-size: 12px;
  color: var(--ink-soft);
}

.pv-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: none;
}

.pv-lang {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 9px;
  color: var(--muted-2);
  border: 1px solid var(--line);
  border-radius: 999px;
  padding: 3px 8px;
}

.pv-lang b {
  color: var(--ink);
}

.pv-lang i {
  font-style: normal;
  opacity: 0.4;
}

.pv-cta {
  border-radius: 999px;
  padding: 7px 14px;
  font-size: 10.5px;
  font-weight: 600;
  white-space: nowrap;
}

.pv-cta--primary {
  background: var(--accent);
  color: #fff;
}

.pv-cta--secondary {
  border: 1px solid var(--line-strong);
  color: var(--ink-soft);
  padding: 6px 13px;
}

.pv-cta--ghost {
  color: var(--ink-soft);
  font-weight: 500;
  padding: 7px 4px;
}

@container (max-width: 860px) {
  .pv-actions {
    gap: 7px;
  }
  .pv-cta {
    padding: 6px 11px;
    font-size: 10px;
  }
  .pv-cta--secondary {
    padding: 5px 10px;
  }
  .pv-cta--ghost {
    padding: 6px 3px;
  }
}

@container (max-width: 760px) {
  .pv-nav {
    display: none;
  }
  .pv-header-inner {
    height: 54px;
  }
  .pv-cta--ghost {
    display: none;
  }
}

@container (max-width: 620px) {
  .pv-lang {
    display: none;
  }
  .pv-cta--secondary {
    display: none;
  }
}
</style>

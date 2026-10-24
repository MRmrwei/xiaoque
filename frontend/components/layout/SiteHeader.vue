<template>
  <div v-if="settings.header.showAnnouncement && announcementVisible" class="announcement">
    <div class="container announcement-inner">
      <p class="announcement-text">
        <span class="pulse-dot" />
        <span>{{ announcement.text[locale] }}</span>
        <span class="announcement-sep">—</span>
        <em>{{ announcement.highlight[locale] }}</em>
      </p>
      <BlockLink
        v-if="announcement.linkLabel[locale]"
        :href="announcement.linkHref"
        class="announcement-link"
      >
        {{ announcement.linkLabel[locale] }} <UiIcon name="arrow-up-right" :size="13" />
      </BlockLink>
      <button class="announcement-close" :aria-label="t('nav.close')" @click="announcementVisible = false">
        <UiIcon name="close" :size="14" />
      </button>
    </div>
  </div>

  <header
    class="site-header"
    :class="{
      'is-scrolled': scrolled,
      'menu-open': menuOpen,
      'is-static': !settings.header.sticky
    }"
  >
    <div class="container header-inner">
      <NuxtLink :to="localePath('/')" class="brand" @click="menuOpen = false">
        <BrandMark :size="34" :image="settings.brand.logo" />
        <span class="brand-text">
          <span class="brand-name">{{ siteName }}</span>
          <span class="brand-sub">{{ settings.siteName.en }}</span>
        </span>
      </NuxtLink>

      <nav class="desktop-nav" :aria-label="t('nav.menu')">
        <NuxtLink
          v-for="item in nav"
          :key="item.path"
          :to="localePath(item.path)"
          class="nav-link"
          :class="{ active: isActive(item.path) }"
        >
          {{ item.label }}
        </NuxtLink>
      </nav>

      <div class="header-actions">
        <div v-if="settings.header.showLanguage" class="lang-switch" role="group" aria-label="Language">
          <button :class="{ active: locale === 'zh' }" @click="switchTo('zh')">中</button>
          <button :class="{ active: locale === 'en' }" @click="switchTo('en')">EN</button>
        </div>
        <BlockLink
          v-for="cta in visibleCtas"
          :key="cta.key"
          :href="cta.href"
          class="header-cta"
          :class="[`header-cta--${cta.variant}`, ctaClass(cta.variant)]"
        >
          {{ cta.label[locale] }}
          <span v-if="cta.variant === 'primary'" class="arr">↗</span>
        </BlockLink>
        <button
          class="mobile-toggle"
          :aria-expanded="menuOpen"
          :aria-label="menuOpen ? t('nav.close') : t('nav.menu')"
          @click="menuOpen = !menuOpen"
        >
          <UiIcon :name="menuOpen ? 'close' : 'menu'" :size="22" />
        </button>
      </div>
    </div>

    <transition name="drop">
      <div v-if="menuOpen" class="mobile-menu">
        <div class="container mobile-menu-inner">
          <nav class="mobile-nav">
            <NuxtLink
              v-for="(item, index) in nav"
              :key="item.path"
              :to="localePath(item.path)"
              class="mobile-nav-link"
              :class="{ active: isActive(item.path) }"
              @click="menuOpen = false"
            >
              <span class="mobile-index">0{{ index + 1 }}</span>
              {{ item.label }}
              <UiIcon name="arrow-up-right" :size="17" />
            </NuxtLink>
          </nav>
          <div class="mobile-menu-foot">
            <div v-if="settings.header.showLanguage" class="mobile-lang">
              <button :class="{ active: locale === 'zh' }" @click="switchTo('zh')">{{ t('lang.zh') }}</button>
              <button :class="{ active: locale === 'en' }" @click="switchTo('en')">{{ t('lang.en') }}</button>
            </div>
            <BlockLink
              v-if="primaryCta"
              :href="primaryCta.href"
              class="btn btn--primary"
              @click="menuOpen = false"
            >
              {{ primaryCta.label[locale] }} <span class="arr">↗</span>
            </BlockLink>
          </div>
        </div>
      </div>
    </transition>
  </header>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import BlockLink from '~/components/BlockLink.vue'
import BrandMark from '~/components/BrandMark.vue'
import UiIcon from '~/components/UiIcon.vue'
import { useLocaleSwitch } from '~/composables/useLocaleSwitch'
import { useSiteChrome } from '~/composables/useSiteChrome'
import { useUiText } from '~/composables/useUiText'
import type { HeaderCtaVariant, Locale } from '~/domain/types'

const route = useRoute()
const localePath = useLocalePath()
const { t, locale } = useUiText()
const { nav, siteName, settings } = useSiteChrome()
const { switchTo: switchLocale } = useLocaleSwitch()

const announcement = computed(() => settings.value.header.announcement)
const visibleCtas = computed(() => settings.value.header.ctas.filter((cta) => cta.visible))
const primaryCta = computed(() =>
  visibleCtas.value.find((cta) => cta.variant === 'primary')
)
const ctaClassMap: Record<HeaderCtaVariant, string> = {
  primary: 'btn btn--primary btn--sm',
  secondary: 'btn btn--ghost btn--sm',
  ghost: 'header-cta--text'
}

function ctaClass(variant: HeaderCtaVariant): string {
  return ctaClassMap[variant]
}

const announcementVisible = ref(true)
const menuOpen = ref(false)
const scrolled = ref(false)

const basePath = computed(() => {
  let path = route.path
  if (path === '/en') return '/'
  if (path.startsWith('/en/')) path = path.slice(3)
  return path || '/'
})

function isActive(path: string) {
  if (path === '/') return basePath.value === '/'
  return basePath.value === path || basePath.value.startsWith(`${path}/`)
}

async function switchTo(target: Locale) {
  await switchLocale(target)
  menuOpen.value = false
}

function onScroll() {
  scrolled.value = window.scrollY > 8
}

watch(
  () => route.fullPath,
  () => {
    menuOpen.value = false
  }
)

onMounted(() => {
  onScroll()
  window.addEventListener('scroll', onScroll, { passive: true })
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<style scoped>
.announcement {
  background: var(--ink);
  color: rgba(246, 242, 234, 0.82);
  height: 38px;
  display: flex;
  align-items: center;
  font-size: 12px;
}

.announcement-inner {
  display: flex;
  align-items: center;
  gap: 18px;
}

.announcement-text {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.announcement-text em {
  font-style: normal;
  color: var(--paper);
  font-weight: 600;
}

.announcement-sep {
  color: rgba(246, 242, 234, 0.35);
}

.pulse-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 0 4px rgba(232, 80, 58, 0.16);
  flex: none;
}

.announcement-link {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--paper);
  font-weight: 600;
  font-size: 12px;
  white-space: nowrap;
  transition: color 0.2s ease;
}

.announcement-link:hover {
  color: #ff9a85;
}

.announcement-close {
  width: 26px;
  height: 26px;
  display: grid;
  place-items: center;
  border: 0;
  background: transparent;
  color: rgba(246, 242, 234, 0.45);
  border-radius: 50%;
  transition: color 0.2s ease, background 0.2s ease;
}

.announcement-close:hover {
  color: var(--paper);
  background: rgba(246, 242, 234, 0.1);
}

/* 吸顶依赖 header 的包含块是整页高度的 .site-shell；
   不要再套一层高度等于头部的包装元素，否则 sticky 没有可滚动空间。 */
.site-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: color-mix(in srgb, var(--paper) 82%, transparent);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  border-bottom: 1px solid transparent;
  transition: background 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease;
}

.site-header.is-static {
  position: relative;
}

.site-header.is-scrolled {
  background: color-mix(in srgb, var(--paper) 94%, transparent);
  border-bottom-color: var(--line-soft);
  box-shadow: 0 6px 30px rgba(20, 21, 24, 0.03);
}

.site-header.menu-open {
  background: var(--paper);
}

.header-inner {
  height: 76px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 11px;
  flex: none;
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1;
  gap: 4px;
}

.brand-name {
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.03em;
  word-spacing: 0.08em;
}

.brand-sub {
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--muted-2);
}

.desktop-nav {
  display: flex;
  align-items: center;
  gap: 28px;
}

.nav-link {
  position: relative;
  font-size: 14px;
  font-weight: 500;
  color: var(--ink-soft);
  padding: 8px 0;
  transition: color 0.2s ease;
}

.nav-link::after {
  content: '';
  position: absolute;
  left: 50%;
  bottom: 0;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--accent);
  transform: translateX(-50%) scale(0);
  transition: transform 0.3s var(--ease);
}

.nav-link:hover {
  color: var(--accent);
}

.nav-link.active::after {
  transform: translateX(-50%) scale(1);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: none;
}

.header-cta--text {
  padding: 8px 0;
  font-size: 13px;
  font-weight: 500;
  color: var(--ink-soft);
  white-space: nowrap;
  transition: color 0.2s ease;
}

.header-cta--text:hover {
  color: var(--accent);
}

.lang-switch {
  display: flex;
  align-items: center;
  padding: 3px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: rgba(251, 249, 245, 0.6);
}

.lang-switch button {
  min-width: 34px;
  height: 28px;
  padding: 0 8px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  font-size: 11px;
  font-weight: 700;
  color: var(--muted);
  transition: background 0.2s ease, color 0.2s ease;
}

.lang-switch button.active {
  background: var(--ink);
  color: var(--paper);
}

.lang-switch button:not(.active):hover {
  color: var(--ink);
}

.admin-link {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  color: var(--muted);
  border-radius: 50%;
  transition: color 0.2s ease, background 0.2s ease;
}

.admin-link:hover {
  color: var(--ink);
  background: rgba(20, 21, 24, 0.06);
}

.mobile-toggle {
  display: none;
  width: 38px;
  height: 38px;
  place-items: center;
  border: 1px solid var(--line);
  background: transparent;
  border-radius: 50%;
  color: var(--ink);
}

.mobile-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--paper);
  border-top: 1px solid var(--line-soft);
  border-bottom: 1px solid var(--line);
  box-shadow: 0 30px 60px rgba(20, 21, 24, 0.12);
  max-height: calc(100vh - 76px);
  overflow-y: auto;
}

.mobile-menu-inner {
  padding-block: 20px 32px;
}

.mobile-nav {
  display: flex;
  flex-direction: column;
}

.mobile-nav-link {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px 2px;
  border-bottom: 1px solid var(--line-soft);
  font-size: 22px;
  font-weight: 600;
  letter-spacing: -0.03em;
}

.mobile-nav-link .ui-icon {
  margin-left: auto;
  opacity: 0.35;
}

.mobile-nav-link.active {
  color: var(--accent);
}

.mobile-index {
  font-family: var(--font-mono);
  font-size: 10px;
  letter-spacing: 0.1em;
  color: var(--muted-2);
}

.mobile-menu-foot {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 14px;
  margin-top: 26px;
}

.mobile-lang {
  display: flex;
  gap: 4px;
  border: 1px solid var(--line);
  border-radius: 999px;
  padding: 3px;
}

.mobile-lang button {
  height: 32px;
  padding: 0 14px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  font-size: 12px;
  font-weight: 600;
  color: var(--muted);
}

.mobile-lang button.active {
  background: var(--ink);
  color: var(--paper);
}

.mobile-admin {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--muted);
}

.drop-enter-active,
.drop-leave-active {
  transition: opacity 0.25s ease, transform 0.3s var(--ease);
}

.drop-enter-from,
.drop-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 1080px) {
  .desktop-nav {
    gap: 20px;
  }
  .nav-link {
    font-size: 13px;
  }
  .admin-link {
    display: none;
  }
}

@media (max-width: 940px) {
  .desktop-nav {
    display: none;
  }
  .mobile-toggle {
    display: grid;
  }
  .header-inner {
    height: 68px;
  }
}

@media (max-width: 760px) {
  .header-cta--ghost {
    display: none;
  }
}

@media (max-width: 620px) {
  .announcement {
    height: 34px;
    font-size: 11px;
  }
  .announcement-text em,
  .announcement-sep,
  .announcement-link {
    display: none;
  }
  .announcement-close {
    margin-left: auto;
  }
  .lang-switch,
  .header-cta--secondary {
    display: none;
  }
  .brand-name {
    font-size: 15px;
  }
  .mobile-menu {
    max-height: calc(100vh - 68px);
  }
  .mobile-nav-link {
    font-size: 19px;
  }
}
</style>

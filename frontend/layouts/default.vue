<template>
  <div class="site-shell">
    <div class="scroll-progress" aria-hidden="true">
      <span :style="{ transform: `scaleX(${progress})` }" />
    </div>
    <SiteHeader />
    <main class="site-main">
      <slot />
    </main>
    <SiteFooter />
    <ToastHost />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import SiteFooter from '~/components/layout/SiteFooter.vue'
import SiteHeader from '~/components/layout/SiteHeader.vue'
import ToastHost from '~/components/ToastHost.vue'
import { getSiteSettings } from '~/services/siteService'
import { brandStyle } from '~/utils/brand'

const themeStyle = brandStyle(getSiteSettings().brand)
const themeCss = Object.entries(themeStyle)
  .map(([key, value]) => `${key}:${value}`)
  .join(';')

useHead({
  htmlAttrs: {
    style: themeCss
  }
})

const progress = ref(0)

function updateProgress() {
  const doc = document.documentElement
  const max = doc.scrollHeight - window.innerHeight
  progress.value = max > 0 ? Math.min(1, Math.max(0, window.scrollY / max)) : 0
}

onMounted(() => {
  updateProgress()
  window.addEventListener('scroll', updateProgress, { passive: true })
  window.addEventListener('resize', updateProgress)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updateProgress)
  window.removeEventListener('resize', updateProgress)
})
</script>

<style scoped>
.site-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.site-main {
  flex: 1;
}

.scroll-progress {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  z-index: 300;
  pointer-events: none;
}

.scroll-progress span {
  display: block;
  height: 100%;
  background: var(--accent);
  transform-origin: left center;
  transition: transform 0.1s linear;
}
</style>

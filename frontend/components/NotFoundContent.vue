<template>
  <section class="not-found">
    <div class="container not-found-inner">
      <div class="not-found-code mono">{{ t('notfound.code') }}</div>
      <div class="not-found-art">
        <BrandMark :size="72" :image="settings.brand.logo" />
        <span class="flight-path" aria-hidden="true" />
      </div>
      <h1>{{ t('notfound.title') }}</h1>
      <p>{{ t('notfound.text') }}</p>
      <div class="btn-row">
        <NuxtLink :to="localePath('/')" class="btn btn--primary">
          {{ t('common.backHome') }} <span class="arr">↗</span>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import BrandMark from '~/components/BrandMark.vue'
import { useSiteChrome } from '~/composables/useSiteChrome'
import { useUiText } from '~/composables/useUiText'

const { t } = useUiText()
const localePath = useLocalePath()
const { settings } = useSiteChrome()
</script>

<style scoped>
.not-found {
  min-height: 64vh;
  display: grid;
  place-items: center;
  padding-block: clamp(70px, 11vw, 150px);
  position: relative;
  overflow: hidden;
}

.not-found::before {
  content: '';
  position: absolute;
  top: 10%;
  right: -10%;
  width: 560px;
  height: 560px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(232, 80, 58, 0.12), transparent 65%);
  pointer-events: none;
}

.not-found-inner {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  z-index: 1;
}

.not-found-code {
  font-size: 11px;
  letter-spacing: 0.2em;
  color: var(--muted-2);
}

.not-found-art {
  margin: 26px 0 20px;
  position: relative;
}

.flight-path {
  position: absolute;
  left: 70px;
  top: 8px;
  width: 110px;
  height: 34px;
  border-top: 1px dashed rgba(20, 21, 24, 0.22);
  border-radius: 50%;
  transform: rotate(-9deg);
}

.flight-path::after {
  content: '';
  position: absolute;
  right: 3px;
  top: -3.5px;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
}

.not-found h1 {
  font-size: clamp(38px, 5.4vw, 68px);
  margin-top: 4px;
  letter-spacing: -0.04em;
}

.not-found p {
  margin-top: 16px;
  color: var(--muted);
  max-width: 44ch;
}

.not-found .btn-row {
  justify-content: center;
  margin-top: 30px;
}

</style>

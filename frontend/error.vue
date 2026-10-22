<template>
  <NuxtLayout>
    <NotFoundContent v-if="isNotFound" />
    <section v-else class="error-fallback">
      <div class="container">
        <p class="mono error-code">{{ error.statusCode }}</p>
        <h1>{{ isEnglish ? 'Something went wrong.' : '页面出了点问题。' }}</h1>
        <p>{{ error.message }}</p>
        <button class="btn btn--primary" @click="handleError">
          {{ isEnglish ? 'Back home' : '返回首页' }}
        </button>
      </div>
    </section>
  </NuxtLayout>
</template>

<script setup lang="ts">
import type { NuxtError } from '#app'
import NotFoundContent from '~/components/NotFoundContent.vue'
import { useSiteChrome } from '~/composables/useSiteChrome'

const props = defineProps<{
  error: NuxtError
}>()

const isNotFound = computed(() => props.error.statusCode === 404)
const isEnglish = computed(() => useRoute().path.startsWith('/en'))
const { settings } = useSiteChrome()
const siteName = computed(() => settings.value.siteName[isEnglish.value ? 'en' : 'zh'])

useHead({
  title: computed(() =>
    isNotFound.value
      ? isEnglish.value
        ? `Page not found · ${siteName.value}`
        : `页面未找到 · ${siteName.value}`
      : isEnglish.value
        ? `Something went wrong · ${siteName.value}`
        : `页面出错了 · ${siteName.value}`
  )
})

function handleError() {
  clearError({ redirect: isEnglish.value ? '/en' : '/' })
}
</script>

<style scoped>
.error-fallback {
  min-height: 60vh;
  display: grid;
  place-items: center;
  text-align: center;
  padding-block: 90px;
}

.error-code {
  font-size: 11px;
  letter-spacing: 0.16em;
  color: var(--muted-2);
}

.error-fallback h1 {
  margin-top: 16px;
  font-size: clamp(34px, 5vw, 60px);
  letter-spacing: -0.04em;
}

.error-fallback p {
  margin-top: 12px;
  color: var(--muted);
}

.error-fallback .btn {
  margin-top: 26px;
}
</style>

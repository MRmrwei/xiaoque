<template>
  <PageRenderer v-if="currentPage" :page="currentPage" :collections="collections" />
  <NotFoundContent v-else />
</template>

<script setup lang="ts">
import NotFoundContent from '~/components/NotFoundContent.vue'
import PageRenderer from '~/components/PageRenderer.vue'
import { usePageContent } from '~/composables/usePageContent'
import { usePageSeo } from '~/composables/usePageSeo'
import { useUiText } from '~/composables/useUiText'

const props = defineProps<{
  path: string
}>()

const { locale } = useUiText()
const { page, collections } = await usePageContent(props.path)
const currentPage = page.value

if (currentPage) {
  usePageSeo({ page: currentPage, path: currentPage.path })
} else {
  if (import.meta.server) {
    const event = useRequestEvent()
    if (event) setResponseStatus(event, 404)
  }
  usePageSeo({
    title: locale.value === 'en' ? 'Page not found · Magpie Cloud' : '页面未找到 · 小鹊云',
    description:
      locale.value === 'en'
        ? 'The page you are looking for may have moved, or never existed.'
        : '你要找的页面可能已经迁移，或者从未存在过。',
    path: props.path,
    noindex: true
  })
}
</script>

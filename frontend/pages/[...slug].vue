<template>
  <SitePageView :path="path" />
</template>

<script setup lang="ts">
import SitePageView from '~/components/SitePageView.vue'

const route = useRoute()
const path = computed(() => {
  const stripped = route.path.replace(/^\/en(?=\/|$)/, '')
  return stripped || '/'
})
</script>

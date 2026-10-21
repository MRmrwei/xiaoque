<template>
  <div class="page-canvas">
    <BlockRenderer
      v-for="block in page.blocks"
      :key="block.id"
      :block="block"
      :locale="page.locale"
      :collections="collections"
    />
  </div>
</template>

<script setup lang="ts">
import BlockRenderer from '~/components/blocks/BlockRenderer.vue'
import type { PublishedPage, SiteCollections } from '~/domain/types'

defineProps<{
  page: PublishedPage
  collections: SiteCollections
}>()
</script>

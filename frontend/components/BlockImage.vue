<template>
  <img v-if="custom" class="block-image" :src="src" :alt="alt" loading="lazy" />
  <VisualArt v-else :variant="artVariant" :tone="tone" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { isCustomImage } from '~/utils/image'

const props = withDefaults(
  defineProps<{
    src?: string
    alt?: string
    tone?: string
    fallback?: string
  }>(),
  { src: '', alt: '', tone: '', fallback: 'orbit' }
)

const custom = computed(() => isCustomImage(props.src))
const artVariant = computed(() => (custom.value || !props.src ? props.fallback : props.src))
</script>

<style scoped>
.block-image {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>

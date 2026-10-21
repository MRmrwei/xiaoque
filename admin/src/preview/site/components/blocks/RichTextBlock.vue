<template>
  <section class="rich-text-block">
    <div class="container">
      <div class="rt-inner" :class="`rt-inner--${data.width || 'narrow'}`">
        <span class="rt-rule" aria-hidden="true" />
        <div class="prose" v-html="data.content" />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RichTextProps } from '~/domain/blocks'

const props = defineProps<{
  propsData: RichTextProps
}>()

const data = computed(() => props.propsData)
</script>

<style scoped>
.rich-text-block {
  padding-block: clamp(52px, 7vw, 100px);
  background: var(--paper);
}

.rt-inner {
  margin-inline: auto;
  position: relative;
}

.rt-inner--narrow {
  max-width: 760px;
}

.rt-inner--wide {
  max-width: 980px;
}

.rt-rule {
  display: block;
  width: 40px;
  height: 2px;
  background: var(--accent);
  margin-bottom: 30px;
}
</style>

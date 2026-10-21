<template>
  <section class="html-block">
    <div class="container">
      <div v-if="preview" class="html-badge mono">
        <UiIcon name="code" :size="12" /> CUSTOM HTML · 发布前会进行安全清洗
      </div>
      <div class="html-content" v-html="data.code" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import UiIcon from '~/components/UiIcon.vue'
import type { HtmlProps } from '~/domain/blocks'

const props = withDefaults(
  defineProps<{
    propsData: HtmlProps
    preview?: boolean
  }>(),
  { preview: false }
)

const data = computed(() => props.propsData)
</script>

<style scoped>
.html-block {
  padding-block: 40px;
  background: var(--paper);
}

.html-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 8.5px;
  letter-spacing: 0.12em;
  color: var(--accent-deep);
  background: var(--accent-soft);
  padding: 5px 10px;
  border-radius: 999px;
  margin-bottom: 16px;
}
</style>

<template>
  <component :is="tag" v-bind="linkAttrs">
    <slot />
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    href?: string
    preview?: boolean
  }>(),
  { href: '', preview: false }
)

const isInternal = computed(
  () => !props.href || props.href.startsWith('/') || props.href.startsWith('#')
)

// Preview canvas is inert: internal actions render as spans, external links open safely.
const tag = computed(() => (props.preview || isInternal.value ? 'span' : 'a'))

function computeLinkAttrs(): Record<string, string> {
  if (props.preview || isInternal.value) return {}
  return { href: props.href, target: '_blank', rel: 'noopener noreferrer' }
}

const linkAttrs = computed(computeLinkAttrs)
</script>

<template>
  <component :is="tag" v-bind="linkAttrs" @click="onClick">
    <slot />
  </component>
</template>

<script setup lang="ts">
import { computed, resolveComponent } from 'vue'

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

const nuxtLink = resolveComponent('NuxtLink')

const tag = computed(() => (props.preview ? 'span' : isInternal.value ? nuxtLink : 'a'))

function computeLinkAttrs(): Record<string, string> {
  if (props.preview) return {}
  if (isInternal.value) return { to: props.href || '/' }
  return { href: props.href, target: '_blank', rel: 'noopener noreferrer' }
}

const linkAttrs = computed(computeLinkAttrs)

function onClick(event: MouseEvent) {
  if (props.preview) event.preventDefault()
}
</script>

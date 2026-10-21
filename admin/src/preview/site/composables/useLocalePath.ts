import { ref } from 'vue'
import type { Locale } from '~/domain/types'

export const previewLocale = ref<Locale>('zh')

export function setPreviewLocale(locale: Locale) {
  previewLocale.value = locale
}

/** Minimal localePath shim for the preview; links are inert in the canvas. */
export function useLocalePath(): (path: string, target?: Locale) => string {
  return (path, target) => {
    const locale = target ?? previewLocale.value
    if (locale !== 'en') return path
    return path === '/' ? '/en' : `/en${path}`
  }
}

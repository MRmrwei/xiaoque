import { computed } from 'vue'
import type { Locale } from '~/domain/types'
import zh from '../locales/zh-CN'
import en from '../locales/en'
import { previewLocale } from './useLocalePath'

export function useUiText() {
  const dictionary = computed<Record<string, string>>(() =>
    previewLocale.value === 'en' ? (en as Record<string, string>) : (zh as Record<string, string>)
  )

  return {
    locale: computed<Locale>(() => previewLocale.value),
    t: (key: string) => dictionary.value[key] ?? key
  }
}

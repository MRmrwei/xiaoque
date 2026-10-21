import type { Locale } from '~/domain/types'
import type { MessageKey } from '~/locales/zh-CN'

export function useUiText() {
  const { t, locale } = useI18n()
  const currentLocale = computed<Locale>(() => (locale.value === 'en' ? 'en' : 'zh'))

  const translate = (key: MessageKey): string => String(t(key))

  return { locale: currentLocale, t: translate }
}

import type { Locale } from '~/domain/types'
import { getNewsBySlug, getNewsTranslation } from '~/services/siteService'

export function useLocaleSwitch() {
  const { locale } = useUiText()
  const route = useRoute()
  const localePath = useLocalePath()
  const switchLocalePath = useSwitchLocalePath()

  async function switchTo(target: Locale) {
    if (target === locale.value) return

    // Article slugs are localized, so the generic switchLocalePath cannot be used.
    // @nuxtjs/i18n appends the locale to generated page names (news-slug___zh).
    if (String(route.name).startsWith('news-slug')) {
      const slug = String(route.params.slug ?? '')
      const current = await getNewsBySlug(locale.value, slug)
      const counterpart = current ? await getNewsTranslation(current) : null
      await navigateTo(
        counterpart
          ? localePath(`/news/${counterpart.slug}`, target)
          : localePath('/news', target)
      )
      return
    }

    await navigateTo(switchLocalePath(target))
  }

  return { switchTo }
}

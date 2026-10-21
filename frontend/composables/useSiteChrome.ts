import { getNavLinks, getSiteSettings } from '~/services/siteService'

export function useSiteChrome() {
  const { locale } = useUiText()
  const settings = getSiteSettings()
  const nav = computed(() => getNavLinks(locale.value))
  const siteName = computed(() => settings.siteName[locale.value])

  return { locale, settings, nav, siteName }
}

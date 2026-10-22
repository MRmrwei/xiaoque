import type { SiteSettings } from '~/domain/types'
import { settings as defaultSettings } from '~/data/settings'
import { getSiteSettingsFromApi } from '~/services/renderService'

/** 站点配置在 SSR 内先加载，头部 / 页脚与页面 SEO 共用同一份状态。 */
function defaultSiteSettings(): SiteSettings {
  return JSON.parse(JSON.stringify(defaultSettings)) as SiteSettings
}

export function useSiteChrome() {
  const { locale } = useUiText()
  const settings = useState<SiteSettings>('site-settings', defaultSiteSettings)

  const nav = computed(() =>
    settings.value.nav
      .filter((item) => item.visible)
      .map((item) => ({ label: item[locale.value], path: item.path }))
  )
  const siteName = computed(() => settings.value.siteName[locale.value])

  return { locale, settings, nav, siteName }
}

/**
 * 加载站点配置与导航（locale 变化时自动重新拉取）。
 * 单文件分享版没有服务端，继续使用内置演示数据。
 */
export async function loadSiteChrome(): Promise<void> {
  const config = useRuntimeConfig()
  if (config.public.useMockData) return

  const { locale } = useUiText()
  const settings = useState<SiteSettings>('site-settings', defaultSiteSettings)

  const { data } = await useAsyncData(
    `site-chrome:${locale.value}`,
    async () => (await getSiteSettingsFromApi(locale.value)) ?? defaultSiteSettings(),
    { watch: [locale] }
  )

  watch(
    data,
    (value) => {
      if (value) settings.value = value
    },
    { immediate: true }
  )
}

import type { Locale, PublishedPage } from '~/domain/types'
import { useSiteChrome } from '~/composables/useSiteChrome'

interface PageSeoOptions {
  page?: PublishedPage | null
  title?: string
  description?: string
  path?: string
  /** 已发布语言对应的路径；未发布的语言不要传，避免 hreflang 指向 404。 */
  alternates?: Partial<Record<Locale, string>>
  type?: 'website' | 'article'
  noindex?: boolean
}

export function usePageSeo(options: PageSeoOptions = {}) {
  const { locale } = useUiText()
  const localePath = useLocalePath()
  const origin = useRequestURL().origin
  const { settings } = useSiteChrome()

  const fallbackPath = options.path ?? options.page?.path ?? '/'
  const alternates = options.alternates
  // 没有 alternates 时沿用旧行为：两种语言都指向当前路径。
  const zhRawPath = alternates?.zh ?? (alternates ? undefined : fallbackPath)
  const enRawPath = alternates?.en ?? (alternates ? undefined : fallbackPath)
  const currentRawPath =
    locale.value === 'en'
      ? (enRawPath ?? zhRawPath ?? fallbackPath)
      : (zhRawPath ?? enRawPath ?? fallbackPath)

  const fallbackTitle = locale.value === 'en' ? 'Website' : '官网'
  const title =
    options.title ??
    options.page?.seo.title ??
    (settings.value.siteName[locale.value] || fallbackTitle)
  const description = options.description ?? options.page?.seo.description ?? ''

  const canonical = new URL(localePath(currentRawPath), origin).href
  const zhHref = zhRawPath ? new URL(localePath(zhRawPath, 'zh'), origin).href : ''
  const enHref = enRawPath ? new URL(localePath(enRawPath, 'en'), origin).href : ''
  const defaultHref = zhHref || enHref

  const alternateLinks = [
    ...(zhHref ? [{ rel: 'alternate', hreflang: 'zh-CN', href: zhHref }] : []),
    ...(enHref ? [{ rel: 'alternate', hreflang: 'en', href: enHref }] : []),
    ...(defaultHref
      ? [{ rel: 'alternate', hreflang: 'x-default', href: defaultHref }]
      : [])
  ]

  useSeoMeta({
    title,
    description,
    ogTitle: title,
    ogDescription: description,
    ogType: options.type ?? 'website'
  })

  useHead({
    htmlAttrs: {
      lang: locale.value === 'en' ? 'en' : 'zh-CN'
    },
    meta: [
      {
        name: 'robots',
        content: computed(() =>
          options.noindex || !settings.value.allowIndex
            ? 'noindex, nofollow'
            : 'index, follow'
        )
      }
    ],
    link: [{ rel: 'canonical', href: canonical }, ...alternateLinks]
  })
}

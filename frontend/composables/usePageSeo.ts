import type { PublishedPage } from '~/domain/types'
import { getSiteSettings } from '~/services/siteService'

interface PageSeoOptions {
  page?: PublishedPage | null
  title?: string
  description?: string
  path?: string
  alternates?: { zh: string; en: string }
  type?: 'website' | 'article'
  noindex?: boolean
}

export function usePageSeo(options: PageSeoOptions = {}) {
  const { locale } = useUiText()
  const localePath = useLocalePath()
  const origin = useRequestURL().origin
  const siteSettings = getSiteSettings()

  const fallbackPath = options.path ?? options.page?.path ?? '/'
  const zhRawPath = options.alternates?.zh ?? fallbackPath
  const enRawPath = options.alternates?.en ?? fallbackPath
  const currentRawPath = locale.value === 'en' ? enRawPath : zhRawPath

  const title = options.title ?? options.page?.seo.title ?? '小鹊云'
  const description = options.description ?? options.page?.seo.description ?? ''

  const canonical = new URL(localePath(currentRawPath), origin).href
  const zhHref = new URL(localePath(zhRawPath, 'zh'), origin).href
  const enHref = new URL(localePath(enRawPath, 'en'), origin).href

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
        content:
          options.noindex || !siteSettings.allowIndex
            ? 'noindex, nofollow'
            : 'index, follow'
      }
    ],
    link: [
      { rel: 'canonical', href: canonical },
      { rel: 'alternate', hreflang: 'zh-CN', href: zhHref },
      { rel: 'alternate', hreflang: 'en', href: enHref },
      { rel: 'alternate', hreflang: 'x-default', href: zhHref }
    ]
  })
}

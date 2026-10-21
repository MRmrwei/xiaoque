import type { CaseRecord, Locale, NewsRecord, ProductRecord, SiteCollections } from '~/domain/types'
import { seedContents, seedSettings } from '@/data/seed/seed'

export const siteSettings = seedSettings

export function collectionsFor(locale: Locale): SiteCollections {
  const published = seedContents.filter(
    (item) => item.locale === locale && item.status === 'published'
  )
  return {
    news: published.filter((item): item is NewsRecord => item.type === 'news'),
    products: published.filter((item): item is ProductRecord => item.type === 'product'),
    cases: published.filter((item): item is CaseRecord => item.type === 'case')
  }
}

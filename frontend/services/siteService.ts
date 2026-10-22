import type {
  CaseRecord,
  ContentRecord,
  Locale,
  NewsRecord,
  PageRecord,
  ProductRecord,
  PublishedPage,
  SiteCollections
} from '~/domain/types'
import { pages } from '~/data/pages'
import { contents } from '~/data/content'

/**
 * Public render data boundary.
 *
 * Today this reads from local seed data. When the Go render API lands, these
 * functions become HTTP calls without touching pages or UI components.
 */

export function normalizePath(path: string): string {
  if (!path) return '/'
  let normalized = path.trim()
  if (!normalized.startsWith('/')) normalized = `/${normalized}`
  if (normalized.length > 1 && normalized.endsWith('/')) {
    normalized = normalized.replace(/\/+$/, '')
  }
  return normalized || '/'
}

function toPublishedPage(page: PageRecord): PublishedPage | null {
  if (!page.published) return null
  return {
    id: page.id,
    groupId: page.groupId,
    locale: page.locale,
    path: page.path,
    name: page.name,
    seo: page.published.seo,
    blocks: page.published.blocks,
    version: page.published.version,
    publishedAt: page.published.publishedAt,
    publishedBy: page.published.publishedBy
  }
}

function isPublished(item: ContentRecord): boolean {
  return item.status === 'published'
}

function isNews(item: ContentRecord): item is NewsRecord {
  return item.type === 'news'
}

function isProduct(item: ContentRecord): item is ProductRecord {
  return item.type === 'product'
}

function isCase(item: ContentRecord): item is CaseRecord {
  return item.type === 'case'
}

export async function getPublishedPage(
  locale: Locale,
  rawPath: string
): Promise<PublishedPage | null> {
  const path = normalizePath(rawPath)
  const page = pages.find((item) => item.locale === locale && normalizePath(item.path) === path)
  return page ? toPublishedPage(page) : null
}

export async function getPageTranslation(page: PublishedPage): Promise<PublishedPage | null> {
  const counterpart = pages.find(
    (item) => item.groupId === page.groupId && item.locale !== page.locale
  )
  return counterpart ? toPublishedPage(counterpart) : null
}

export async function getNews(locale: Locale, limit?: number): Promise<NewsRecord[]> {
  const items = contents
    .filter(isNews)
    .filter((item) => item.locale === locale && isPublished(item))
  return typeof limit === 'number' ? items.slice(0, limit) : items
}

export async function getNewsBySlug(locale: Locale, slug: string): Promise<NewsRecord | null> {
  return (
    contents
      .filter(isNews)
      .find((item) => item.locale === locale && item.slug === slug && isPublished(item)) ?? null
  )
}

export async function getNewsTranslation(item: NewsRecord): Promise<NewsRecord | null> {
  return (
    contents
      .filter(isNews)
      .find(
        (candidate) =>
          candidate.groupId === item.groupId &&
          candidate.locale !== item.locale &&
          isPublished(candidate)
      ) ?? null
  )
}

export async function getRelatedNews(
  locale: Locale,
  excludeId: string,
  limit = 2
): Promise<NewsRecord[]> {
  const items = await getNews(locale)
  return items.filter((item) => item.id !== excludeId).slice(0, limit)
}

export async function getProducts(locale: Locale, limit?: number): Promise<ProductRecord[]> {
  const items = contents
    .filter(isProduct)
    .filter((item) => item.locale === locale && isPublished(item))
  return typeof limit === 'number' ? items.slice(0, limit) : items
}

export async function getCases(locale: Locale, limit?: number): Promise<CaseRecord[]> {
  const items = contents
    .filter(isCase)
    .filter((item) => item.locale === locale && isPublished(item))
  return typeof limit === 'number' ? items.slice(0, limit) : items
}

export async function getSiteCollections(locale: Locale): Promise<SiteCollections> {
  const [news, products, cases] = await Promise.all([
    getNews(locale),
    getProducts(locale),
    getCases(locale)
  ])
  return { news, products, cases }
}

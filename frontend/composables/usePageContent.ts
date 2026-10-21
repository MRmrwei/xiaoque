import type { SiteCollections } from '~/domain/types'
import { getPublishedPage, getSiteCollections, normalizePath } from '~/services/siteService'

export async function usePageContent(rawPath: string) {
  const { locale } = useUiText()
  const path = normalizePath(rawPath)

  const { data: loaded } = await useAsyncData(
    `page:${locale.value}:${path}`,
    async () => ({ page: await getPublishedPage(locale.value, path) }),
    { watch: [locale] }
  )

  const { data: collections } = await useAsyncData(
    `collections:${locale.value}`,
    () => getSiteCollections(locale.value),
    { watch: [locale] }
  )

  const page = computed(() => loaded.value?.page ?? null)
  const resolvedCollections = computed<SiteCollections>(
    () => collections.value ?? { news: [], products: [], cases: [] }
  )

  return { page, collections: resolvedCollections, path }
}

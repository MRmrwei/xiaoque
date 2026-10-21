import type { NewsRecord } from '~/domain/types'
import { getNewsBySlug, getNewsTranslation, getRelatedNews } from '~/services/siteService'

export async function useNewsArticle(slug: string) {
  const { locale } = useUiText()

  const { data: loadedArticle } = await useAsyncData(
    `news:${locale.value}:${slug}`,
    async () => ({ article: await getNewsBySlug(locale.value, slug) })
  )

  const article = computed(() => loadedArticle.value?.article ?? null)

  const { data: loadedTranslation } = await useAsyncData(
    `news-translation:${locale.value}:${slug}`,
    async () => {
      const current = article.value
      return { translation: current ? await getNewsTranslation(current) : null }
    },
    { watch: [article] }
  )

  const translation = computed(() => loadedTranslation.value?.translation ?? null)

  const { data: related } = await useAsyncData(
    `news-related:${locale.value}:${slug}`,
    async () => {
      const current = article.value
      return current ? await getRelatedNews(locale.value, current.id, 2) : []
    },
    { watch: [article] }
  )

  const relatedNews = computed<NewsRecord[]>(() => related.value ?? [])

  return { article, translation, related: relatedNews }
}

import type { SiteCollections } from '~/domain/types'
import { getPublishedPage, getSiteCollections, normalizePath } from '~/services/siteService'
import { getRenderPage } from '~/services/renderService'

const EMPTY_COLLECTIONS: SiteCollections = { news: [], products: [], cases: [] }

export async function usePageContent(rawPath: string) {
  const { locale } = useUiText()
  const path = normalizePath(rawPath)
  const config = useRuntimeConfig()
  const alternates = useState<Record<string, string>>('page-alternates', () => ({}))

  // 单文件分享版没有服务端接口，继续读取内置演示数据。
  if (config.public.useMockData) {
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
    alternates.value = {}
    return {
      page,
      collections: computed<SiteCollections>(() => collections.value ?? EMPTY_COLLECTIONS),
      path
    }
  }

  const { data } = await useAsyncData(
    `render:${locale.value}:${path}`,
    () => getRenderPage(locale.value, path),
    { watch: [locale] }
  )

  const page = computed(() => data.value?.page ?? null)
  const collections = computed<SiteCollections>(
    () => data.value?.collections ?? EMPTY_COLLECTIONS
  )

  watch(
    page,
    (value) => {
      alternates.value = value?.alternates
        ? {
            zh: value.alternates.zh ?? '',
            en: value.alternates.en ?? ''
          }
        : {}
    },
    { immediate: true }
  )

  return { page, collections, path }
}

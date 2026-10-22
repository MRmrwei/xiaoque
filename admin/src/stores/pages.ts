import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { Locale, PageRecord, SiteCollections, SiteSettings } from '~/domain/types'
import type {
  ContactSubmissionRecord,
  LoginSettings,
  PageStatus,
  PersistedDemo,
  SubmissionStatus
} from '@/types/admin'
import { clone, createBlock } from '@/domain/blockSchema'
import { seedLoginSettings, seedSettings } from '@/data/seed/seed'
import { loadDemo, saveDemo } from '@/services/storage'
import { collectionsFor } from '@/services/contentService'
import { normalizeSiteSettings } from '@/services/settingsService'
import {
  duplicatePageRecord,
  formatClock,
  normalizePath,
  pageStatus
} from '@/services/pageService'
import {
  createPage as createPageApi,
  deletePage as deletePageApi,
  getPage,
  getSettings,
  getMenu,
  listMenus,
  listPages,
  publishPage as publishPageApi,
  saveMenu,
  savePageDraft,
  saveSettings as saveSettingsApi,
  translatePage as translatePageApi,
  unpublishPage as unpublishPageApi,
  type ApiMenuDetail,
  type ApiPageDraft,
  type ApiPageSummary
} from '@/services/siteApi'
import {
  bottomMenuPayload,
  footerMenuPayload,
  formatApiTime,
  headerMenuPayload,
  pageBlocksToApi,
  pageFromApi,
  siteSettingsFromApi,
  siteSettingsPayload,
  toApiLocale
} from '@/services/siteMapping'
import {
  listSubmissions,
  submissionFromApi,
  updateSubmissionStatus,
  type ApiSubmissionQuery,
  type ApiSubmissionStats
} from '@/services/submissionApi'

export interface SubmissionQuery {
  status?: SubmissionStatus | 'all'
  locale?: 'zh' | 'en' | 'all'
  keyword?: string
  page?: number
  pageSize?: number
}

export const usePagesStore = defineStore('pages', () => {
  const persisted = loadDemo()

  const pages = ref<PageRecord[]>([])
  const settings = ref<SiteSettings>(normalizeSiteSettings())
  const submissions = ref<ContactSubmissionRecord[]>([])
  const submissionStats = ref<ApiSubmissionStats>({
    total: 0,
    new: 0,
    processing: 0,
    done: 0,
    spam: 0
  })
  const submissionTotal = ref(0)
  const submissionsLoading = ref(false)
  const submissionsError = ref('')
  const loginSettings = ref<LoginSettings>(
    clone({ ...seedLoginSettings, ...(persisted?.loginSettings ?? {}) })
  )

  const loading = ref(false)
  const loaded = ref(false)
  const loadError = ref('')
  const storageError = ref('')
  const lastSavedAt = ref(formatClock())
  let loadTask: Promise<void> | null = null
  let submissionRequestId = 0
  let lastSubmissionQuery: SubmissionQuery = { page: 1, pageSize: 20 }

  /** 页面与站点设置改由 Go API 保存；localStorage 只保留表单留言与登录页设置。 */
  function persist(): void {
    const demo: PersistedDemo = {
      schemaVersion: 1,
      savedAt: new Date().toISOString(),
      loginSettings: clone(loginSettings.value)
    }
    if (saveDemo(demo)) {
      storageError.value = ''
    } else {
      storageError.value = '当前浏览器无法保存设置，刷新后可能丢失。'
    }
  }

  async function fetchPageRecords(): Promise<PageRecord[]> {
    const summaries: ApiPageSummary[] = []
    let page = 1
    for (;;) {
      const result = await listPages({ page, page_size: 100 })
      summaries.push(...result.items)
      if (result.items.length === 0 || summaries.length >= result.total || page >= 50) break
      page += 1
    }

    // 列表接口不返回区块，编辑器与预览需要完整草稿，因此逐页取详情。
    const settled = await Promise.allSettled(
      summaries.map((item) => getPage(item.id))
    )
    return settled
      .filter((item): item is PromiseFulfilledResult<ApiPageDraft> => item.status === 'fulfilled')
      .map((item) => pageFromApi(item.value))
  }

  async function fetchSiteSettings(): Promise<SiteSettings> {
    const [raw, menuList] = await Promise.all([getSettings('zh-CN'), listMenus()])
    const codes = new Set(menuList.items.map((item) => item.code))

    const [header, footer, footerBottom] = await Promise.all([
      codes.has('header') ? getMenu('header') : Promise.resolve(undefined),
      codes.has('footer') ? getMenu('footer') : Promise.resolve(undefined),
      codes.has('footer_bottom') ? getMenu('footer_bottom') : Promise.resolve(undefined)
    ])

    return siteSettingsFromApi(raw as unknown as Record<string, unknown>, {
      header: header as ApiMenuDetail | undefined,
      footer: footer as ApiMenuDetail | undefined,
      footerBottom: footerBottom as ApiMenuDetail | undefined
    })
  }

  /** 从 Go API 加载页面与站点配置；重复调用会复用同一次请求。 */
  async function loadSubmissions(query: SubmissionQuery = {}): Promise<void> {
    lastSubmissionQuery = query
    const requestId = ++submissionRequestId
    submissionsLoading.value = true
    submissionsError.value = ''
    try {
      const result = await listSubmissions(toApiSubmissionQuery(query))
      if (requestId !== submissionRequestId) return
      submissions.value = result.items.map(submissionFromApi)
      submissionStats.value = result.stats
      submissionTotal.value = result.total
      // 只有未筛选的全局查询才更新侧边栏红点，避免筛选语言 / 关键词时数量变小。
      if (isGlobalSubmissionQuery(query)) {
        newSubmissionCount.value = result.stats.new
      }
    } catch (error) {
      if (requestId !== submissionRequestId) return
      submissionsError.value =
        error instanceof Error ? error.message : '加载留言失败，请稍后重试'
    } finally {
      if (requestId === submissionRequestId) submissionsLoading.value = false
    }
  }

  /** 只刷新待处理留言数，供侧边栏红点使用；不覆盖当前列表、筛选与统计卡片。 */
  async function refreshSubmissionStats(): Promise<void> {
    try {
      const result = await listSubmissions({ page: 1, page_size: 1 })
      newSubmissionCount.value = result.stats.new
    } catch {
      // 静默失败：保留上一次统计，不打断当前操作。
    }
  }

  function isGlobalSubmissionQuery(query: SubmissionQuery): boolean {
    return (
      (!query.status || query.status === 'all') &&
      (!query.locale || query.locale === 'all') &&
      !query.keyword?.trim()
    )
  }

  function toApiSubmissionQuery(query: SubmissionQuery): ApiSubmissionQuery {
    const result: ApiSubmissionQuery = {
      page: query.page && query.page > 0 ? query.page : 1,
      page_size: query.pageSize && query.pageSize > 0 ? query.pageSize : 20
    }
    if (query.status && query.status !== 'all') result.status = query.status
    if (query.locale && query.locale !== 'all') {
      result.locale = query.locale === 'en' ? 'en' : 'zh-CN'
    }
    if (query.keyword?.trim()) result.keyword = query.keyword.trim()
    return result
  }

  async function ensureLoaded(force = false): Promise<void> {
    if (loaded.value && !force) return
    if (loadTask) return loadTask

    loading.value = true
    loadError.value = ''
    loadTask = (async () => {
      try {
        const [pageRecords, nextSettings] = await Promise.all([
          fetchPageRecords(),
          fetchSiteSettings(),
          loadSubmissions(lastSubmissionQuery)
        ])
        pages.value = pageRecords
        settings.value = nextSettings
        loaded.value = true
        lastSavedAt.value = formatClock()
        // 清掉旧版本地演示缓存里的页面与站点设置，只保留留言 / 登录设置。
        persist()
      } catch (error) {
        loadError.value =
          error instanceof Error ? error.message : '加载站点数据失败，请稍后重试'
      } finally {
        loading.value = false
        loadTask = null
      }
    })()
    return loadTask
  }

  async function reloadAll(): Promise<void> {
    await ensureLoaded(true)
  }

  function getById(id: string): PageRecord | undefined {
    return pages.value.find((page) => page.id === id)
  }

  function getByPath(locale: Locale, path: string): PageRecord | undefined {
    return pages.value.find((page) => page.locale === locale && page.path === path)
  }

  function pairOf(page: PageRecord | null | undefined): PageRecord | undefined {
    if (!page) return undefined
    return pages.value.find(
      (candidate) => candidate.groupId === page.groupId && candidate.locale !== page.locale
    )
  }

  function statusOf(page: PageRecord): PageStatus {
    return pageStatus(page)
  }

  async function createPage(input: {
    name: string
    path: string
    locale: Locale
    template?: 'blank' | 'home'
  }): Promise<PageRecord> {
    let copyFromId = ''
    if (input.template === 'home') {
      copyFromId = getByPath(input.locale, '/')?.id ?? ''
    }

    const hero = copyFromId ? null : createBlock('hero')
    const blocks = hero ? [hero] : []
    const draft = await createPageApi({
      locale: toApiLocale(input.locale),
      path: normalizePath(input.path),
      name: input.name,
      seo: { title: input.name, description: '', keywords: [] },
      blocks: pageBlocksToApi(blocks),
      ...(copyFromId ? { copy_from_id: copyFromId } : {})
    })

    const page = pageFromApi(draft)
    pages.value.unshift(page)
    return page
  }

  async function duplicatePage(page: PageRecord): Promise<PageRecord> {
    const template = duplicatePageRecord(page, pages.value)
    const draft = await createPageApi({
      locale: toApiLocale(page.locale),
      path: template.path,
      name: template.name,
      copy_from_id: page.id
    })
    const copy = pageFromApi(draft)
    pages.value.unshift(copy)
    return copy
  }

  async function deletePage(pageId: string): Promise<void> {
    await deletePageApi(pageId)
    const index = pages.value.findIndex((page) => page.id === pageId)
    if (index >= 0) pages.value.splice(index, 1)
  }

  /** 保存草稿。返回后只同步服务端确定的元数据，避免覆盖请求期间的新编辑。 */
  async function savePage(page: PageRecord): Promise<void> {
    const draft = await savePageDraft(page.id, {
      path: normalizePath(page.path),
      name: page.name,
      seo: page.seo,
      blocks: pageBlocksToApi(page.blocks)
    })

    page.path = draft.path
    page.name = draft.name
    page.status = draft.status === 'published' ? 'published' : 'draft'
    page.updatedAt = formatApiTime(draft.updated_at)
    if (page.status === 'published') {
      page.published = {
        version: draft.version,
        blocks: page.published?.blocks ?? [],
        seo: page.published?.seo ?? clone(page.seo),
        publishedAt: page.published?.publishedAt || formatApiTime(draft.updated_at),
        publishedBy: page.published?.publishedBy ?? ''
      }
    } else {
      page.published = null
    }
  }

  /** 发布 = 用当前草稿覆盖线上内容，不保留历史版本。 */
  async function publishPage(page: PageRecord): Promise<void> {
    const result = await publishPageApi(page.id)
    page.status = 'published'
    page.published = {
      version: result.version,
      blocks: clone(page.blocks),
      seo: clone(page.seo),
      publishedAt: formatApiTime(result.published_at),
      publishedBy: ''
    }
    page.updatedAt = formatApiTime(result.published_at)
  }

  /** 下架：状态回到未发布，并清空线上内容；草稿保留。 */
  async function unpublishPage(page: PageRecord): Promise<void> {
    await unpublishPageApi(page.id)
    page.status = 'draft'
    page.published = null
    page.updatedAt = new Date().toISOString()
  }

  async function createTranslation(page: PageRecord): Promise<PageRecord> {
    const existing = pairOf(page)
    if (existing) return existing

    const targetLocale: Locale = page.locale === 'zh' ? 'en' : 'zh'
    const draft = await translatePageApi(page.id, {
      to_locale: toApiLocale(targetLocale),
      path: page.path
    })
    const translated = pageFromApi(draft)
    pages.value.push(translated)
    return translated
  }

  /** 保存站点设置 + 顶部 / 底部导航到 Go API。 */
  async function saveSettings(): Promise<void> {
    // 菜单是整体替换（先删后插），串行提交，避免三个菜单事务并发时互相锁等待。
    await saveSettingsApi('zh-CN', siteSettingsPayload(settings.value))
    await saveMenu('header', headerMenuPayload(settings.value))
    await saveMenu('footer', footerMenuPayload(settings.value))
    await saveMenu('footer_bottom', bottomMenuPayload(settings.value))
    lastSavedAt.value = formatClock()
  }

  /** 保存后台登录页设置（验证码、展示文案，仍为本地演示）。 */
  function saveLoginSettings(next: LoginSettings): void {
    loginSettings.value = clone(next)
    persist()
  }

  /** 只把顶部导航与头部配置恢复为种子数据（仍需手动保存）。 */
  function resetNavigation(): void {
    settings.value.nav = clone(seedSettings.nav)
    settings.value.header = clone(seedSettings.header)
    settings.value.siteName = clone(seedSettings.siteName)
  }

  /** 只把页脚编排（含联系信息、标语、版权、备案号）恢复为种子数据（仍需手动保存）。 */
  function resetFooter(): void {
    settings.value.footer = clone(seedSettings.footer)
    settings.value.tagline = clone(seedSettings.tagline)
    settings.value.copyright = clone(seedSettings.copyright)
    settings.value.icp = seedSettings.icp
    settings.value.email = seedSettings.email
    settings.value.phone = seedSettings.phone
    settings.value.phoneEn = seedSettings.phoneEn
    settings.value.address = clone(seedSettings.address)
  }

  /** 把品牌外观与全站配置恢复为种子数据（仍需手动保存）。 */
  function resetSite(): void {
    settings.value.brand = clone(seedSettings.brand)
    settings.value.defaultLocale = seedSettings.defaultLocale
    settings.value.allowIndex = seedSettings.allowIndex
    settings.value.headCode = seedSettings.headCode
  }

  function collections(locale: Locale): SiteCollections {
    return collectionsFor(locale)
  }

  async function setSubmissionStatus(
    id: string,
    status: SubmissionStatus
  ): Promise<void> {
    const updated = await updateSubmissionStatus(id, status)
    const record = submissionFromApi(updated)
    const index = submissions.value.findIndex((item) => item.id === id)
    if (index >= 0) submissions.value[index] = record
    await loadSubmissions(lastSubmissionQuery)
    // 状态变更可能让待处理数量变化，立即同步红点。
    void refreshSubmissionStats()
  }

  const pageCount = computed(() => pages.value.length)
  const submissionCount = computed(() => submissionStats.value.total)
  /** 全站待处理留言数，只由全局查询或定时刷新更新。 */
  const newSubmissionCount = ref(0)

  return {
    pages,
    settings,
    submissions,
    submissionStats,
    submissionTotal,
    submissionsLoading,
    submissionsError,
    loginSettings,
    loading,
    loaded,
    loadError,
    storageError,
    lastSavedAt,
    pageCount,
    submissionCount,
    newSubmissionCount,
    ensureLoaded,
    reloadAll,
    loadSubmissions,
    refreshSubmissionStats,
    setSubmissionStatus,
    getById,
    getByPath,
    pairOf,
    statusOf,
    createPage,
    duplicatePage,
    deletePage,
    savePage,
    publishPage,
    unpublishPage,
    createTranslation,
    saveSettings,
    saveLoginSettings,
    resetNavigation,
    resetFooter,
    resetSite,
    collections
  }
})

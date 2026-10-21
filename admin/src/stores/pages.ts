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
import { clone, uid } from '@/domain/blockSchema'
import {
  createSeedPersisted,
  seedLoginSettings,
  seedSettings,
  seedSubmissions
} from '@/data/seed/seed'
import { clearDemo, loadDemo, saveDemo } from '@/services/storage'
import { collectionsFor } from '@/services/contentService'
import { normalizeSiteSettings } from '@/services/settingsService'
import {
  clearTranslatableBlocks,
  createBlankPage,
  duplicatePageRecord,
  emptySeo,
  formatClock,
  migrateLegacyBlocks,
  nowStamp,
  pageStatus
} from '@/services/pageService'

function migratePage(page: PageRecord): PageRecord {
  return {
    ...page,
    blocks: migrateLegacyBlocks(page.blocks),
    published: page.published
      ? { ...page.published, blocks: migrateLegacyBlocks(page.published.blocks) }
      : null
  }
}

export const usePagesStore = defineStore('pages', () => {
  const seeded = createSeedPersisted()
  const persisted = loadDemo()

  const pages = ref<PageRecord[]>(
    (persisted?.pages?.length ? persisted.pages : seeded.pages).map(migratePage)
  )
  const storageError = ref('')
  const lastSavedAt = ref(
    persisted?.savedAt
      ? new Intl.DateTimeFormat('zh-CN', { hour: '2-digit', minute: '2-digit' }).format(
          new Date(persisted.savedAt)
        )
      : formatClock()
  )
  const settings = ref<SiteSettings>(normalizeSiteSettings(persisted?.settings))
  const submissions = ref<ContactSubmissionRecord[]>(
    clone(persisted?.submissions ?? seeded.submissions ?? seedSubmissions)
  )
  const loginSettings = ref<LoginSettings>(
    clone({ ...seedLoginSettings, ...(persisted?.loginSettings ?? {}) })
  )

  function persist(): void {
    const demo: PersistedDemo = {
      schemaVersion: 1,
      savedAt: new Date().toISOString(),
      pages: clone(pages.value),
      submissions: clone(submissions.value),
      loginSettings: clone(loginSettings.value),
      settings: clone(settings.value)
    }
    if (saveDemo(demo)) {
      storageError.value = ''
      lastSavedAt.value = formatClock()
    } else {
      storageError.value = '浏览器本地存储不可用，修改仅保留在当前页面。'
    }
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

  function createPage(input: {
    name: string
    path: string
    locale: Locale
    template?: 'blank' | 'home'
  }): PageRecord {
    const page = createBlankPage(input.name, input.path, input.locale)
    if (input.template === 'home') {
      const home = getByPath(input.locale, '/')
      if (home) {
        page.blocks = clone(home.blocks).map(
          (block) => ({ ...block, id: uid('blk') }) as PageRecord['blocks'][number]
        )
      }
    }
    pages.value.push(page)
    persist()
    return page
  }

  function duplicatePage(page: PageRecord): PageRecord {
    const copy = duplicatePageRecord(page, pages.value)
    pages.value.push(copy)
    persist()
    return copy
  }

  function deletePage(pageId: string): void {
    const index = pages.value.findIndex((page) => page.id === pageId)
    if (index >= 0) pages.value.splice(index, 1)
    persist()
  }

  function savePage(page: PageRecord): void {
    page.updatedAt = nowStamp()
    persist()
  }

  /** 发布 = 用当前草稿覆盖线上内容，不保留历史版本。 */
  function publishPage(page: PageRecord): void {
    const versionNumber = (page.published?.version ?? 0) + 1
    const at = nowStamp()
    page.published = {
      version: versionNumber,
      blocks: clone(page.blocks),
      seo: clone(page.seo),
      publishedAt: at,
      publishedBy: '林小鹊'
    }
    page.status = 'published'
    page.updatedAt = at
    persist()
  }

  /** 下架：状态回到未发布，并清空线上内容；草稿保留。 */
  function unpublishPage(page: PageRecord): void {
    page.status = 'draft'
    page.published = null
    page.updatedAt = nowStamp()
    persist()
  }

  function createTranslation(page: PageRecord): PageRecord {
    const existing = pairOf(page)
    if (existing) return existing
    const targetLocale: Locale = page.locale === 'zh' ? 'en' : 'zh'
    const translated: PageRecord = {
      id: uid('pg'),
      groupId: page.groupId,
      locale: targetLocale,
      path: page.path,
      name: `${page.name} ${targetLocale === 'en' ? '(EN)' : '（中文）'}`,
      status: 'draft',
      updatedAt: nowStamp(),
      seo: emptySeo(),
      blocks: clearTranslatableBlocks(page.blocks),
      published: null
    }
    pages.value.push(translated)
    persist()
    return translated
  }

  function resetDemo(): void {
    clearDemo()
    const fresh = createSeedPersisted()
    pages.value = fresh.pages.map(migratePage)
    submissions.value = clone(fresh.submissions ?? seedSubmissions)
    loginSettings.value = clone(fresh.loginSettings ?? seedLoginSettings)
    settings.value = normalizeSiteSettings()
    storageError.value = ''
    lastSavedAt.value = formatClock()
  }

  /** 保存站点级设置（顶部导航 / 底部导航等）。 */
  function saveSettings(): void {
    persist()
  }

  /** 保存后台登录页设置（验证码、展示文案）。 */
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

  function getSubmissionById(id: string): ContactSubmissionRecord | undefined {
    return submissions.value.find((item) => item.id === id)
  }

  function setSubmissionStatus(id: string, status: SubmissionStatus): void {
    const item = submissions.value.find((entry) => entry.id === id)
    if (!item) return
    item.status = status
    item.updatedAt = new Date().toISOString()
    persist()
  }

  const pageCount = computed(() => pages.value.length)
  const submissionCount = computed(() => submissions.value.length)
  const newSubmissionCount = computed(
    () => submissions.value.filter((item) => item.status === 'new').length
  )

  return {
    pages,
    settings,
    submissions,
    loginSettings,
    storageError,
    lastSavedAt,
    pageCount,
    submissionCount,
    newSubmissionCount,
    getSubmissionById,
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
    resetDemo,
    saveSettings,
    saveLoginSettings,
    resetNavigation,
    resetFooter,
    resetSite,
    collections
  }
})

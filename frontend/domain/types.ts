import type { PageBlock } from './blocks'

export type Locale = 'zh' | 'en'

export interface SeoMeta {
  title: string
  description: string
  keywords: string[]
}

export interface PagePublishedSnapshot {
  version: number
  blocks: PageBlock[]
  seo: SeoMeta
  publishedAt: string
  publishedBy: string
}

export interface PageRecord {
  id: string
  groupId: string
  locale: Locale
  path: string
  name: string
  status: 'published' | 'draft'
  updatedAt: string
  seo: SeoMeta
  blocks: PageBlock[]
  published: PagePublishedSnapshot | null
}

/** The immutable view model returned by the public render API. */
export interface PublishedPage {
  id: string
  groupId: string
  locale: Locale
  path: string
  name: string
  seo: SeoMeta
  blocks: PageBlock[]
  version: number
  publishedAt: string
  publishedBy: string
  /** 已发布语言对应的路径，用于 hreflang 与语言切换；未发布语言不出现。 */
  alternates?: Partial<Record<Locale, string>>
}

export type ContentType = 'news' | 'product' | 'case'

export interface ContentBase {
  id: string
  groupId: string
  locale: Locale
  slug: string
  title: string
  summary: string
  cover: string
  tone: string
  status: 'published' | 'draft'
  publishedAt: string
  seoTitle?: string
  seoDescription?: string
}

export interface NewsRecord extends ContentBase {
  type: 'news'
  category: string
  body: string
  readMinutes: number
  featured: boolean
}

export interface ProductRecord extends ContentBase {
  type: 'product'
  label: string
  description: string
  bullets: string[]
}

export interface CaseRecord extends ContentBase {
  type: 'case'
  client: string
  industry: string
  metric: string
  metricSuffix: string
  metricLabel: string
}

export type ContentRecord = NewsRecord | ProductRecord | CaseRecord

export interface SiteCollections {
  news: NewsRecord[]
  products: ProductRecord[]
  cases: CaseRecord[]
}

export interface NavItem {
  key: string
  path: string
  zh: string
  en: string
  visible: boolean
}

export interface NavLink {
  label: string
  path: string
}

export interface SocialLink {
  name: string
  handle: string
}

export type HeaderCtaVariant = 'primary' | 'secondary' | 'ghost'

/** 顶部右侧按钮；建议 1 个主按钮 + 最多 2 个辅助按钮。 */
export interface HeaderCta {
  key: string
  label: Record<Locale, string>
  href: string
  variant: HeaderCtaVariant
  visible: boolean
}

export interface HeaderAnnouncement {
  /** 公告主文案。 */
  text: Record<Locale, string>
  /** 高亮文案，接在主文案后面。 */
  highlight: Record<Locale, string>
  /** 右侧链接文案。 */
  linkLabel: Record<Locale, string>
  linkHref: string
}

export interface HeaderSettings {
  /** 是否显示右侧语言切换。 */
  showLanguage: boolean
  /** 页面滚动时是否吸顶。 */
  sticky: boolean
  /** 是否显示顶部公告条。 */
  showAnnouncement: boolean
  announcement: HeaderAnnouncement
  /** 右侧按钮，按数组顺序渲染。 */
  ctas: HeaderCta[]
}

/** 页脚栏目或底部栏中的一条链接。 */
export interface FooterLink {
  key: string
  path: string
  zh: string
  en: string
  visible: boolean
}

/** 页脚的一个链接栏目。 */
export interface FooterColumn {
  key: string
  title: Record<Locale, string>
  visible: boolean
  links: FooterLink[]
}

export interface FooterCta {
  label: Record<Locale, string>
  path: string
}

export interface FooterSettings {
  /** 联系列：邮箱 / 电话 / 地址沿用站点设置，CTA 单独配置。 */
  showContact: boolean
  cta: FooterCta
  /** 底部大字标。 */
  showWordmark: boolean
  wordmark: string
  columns: FooterColumn[]
  bottomLinks: FooterLink[]
}

export interface BrandSettings {
  /** 官网 Logo，data URL 或远程地址；为空时使用内置品牌图标。 */
  logo: string
  /** 品牌主色，十六进制色值。 */
  accent: string
  /** 页面整体底色，十六进制色值。 */
  background: string
  /** 页脚与深色区块底色，十六进制色值。 */
  darkBackground: string
}

export interface SiteSettings {
  nav: NavItem[]
  header: HeaderSettings
  footer: FooterSettings
  brand: BrandSettings
  siteName: Record<Locale, string>
  companyName: Record<Locale, string>
  tagline: Record<Locale, string>
  email: string
  phone: string
  phoneEn: string
  address: Record<Locale, string>
  workHours: Record<Locale, string>
  icp: string
  copyright: Record<Locale, string>
  socials: SocialLink[]
  defaultLocale: Locale
  analytics: string
  headCode: string
  allowIndex: boolean
}

export interface ContactInput {
  name: string
  email: string
  company: string
  phone: string
  message: string
}

export interface ContactSubmission extends ContactInput {
  createdAt: string
  page: string
}

export interface ToastItem {
  id: number
  text: string
  type: 'success' | 'error'
}

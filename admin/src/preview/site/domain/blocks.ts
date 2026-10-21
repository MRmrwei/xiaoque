export type VisualKey =
  | 'builder'
  | 'orbit'
  | 'arch'
  | 'grid'
  | 'dunes'
  | 'waves'
  | 'type'
  | 'stack'
  | 'bloom'
  | 'none'

export interface HeroProps {
  variant: 'home' | 'page'
  eyebrow: string
  title: string
  subtitle: string
  primaryLabel: string
  primaryHref: string
  secondaryLabel: string
  secondaryHref: string
  note: string
  /** 内置图形 key，或自定义图片 URL。 */
  visual: VisualKey | (string & {})
  theme: 'light' | 'dark'
  /** 右下角竖排小字，留空则不显示。 */
  sideNote?: string
  /** 左下角滚动提示，留空则不显示。 */
  scrollLabel?: string
}

export interface HeroCarouselSlide {
  id: string
  mediaType: 'image' | 'video'
  /** 图片模式的主视觉；视频模式下作为旧数据兜底封面。 */
  media: VisualKey | (string & {})
  videoUrl: string
  /** 视频模式的上传封面，作为 <video poster>。 */
  videoPoster?: string
  eyebrow: string
  title: string
  subtitle: string
  primaryLabel: string
  primaryHref: string
  secondaryLabel: string
  secondaryHref: string
}

export interface HeroCarouselProps {
  slides: HeroCarouselSlide[]
  autoplay: boolean
  interval: number
  note: string
}

export interface PartnerItem {
  name: string
}

export interface PartnersProps {
  title: string
  note: string
  marquee: boolean
  items: PartnerItem[]
}

export interface FeatureItem {
  icon: string
  title: string
  desc: string
  tag: string
}

export interface FeaturesProps {
  eyebrow: string
  title: string
  subtitle: string
  layout: 'bento' | 'grid'
  items: FeatureItem[]
}

export interface ImageTextProps {
  eyebrow: string
  title: string
  body: string
  bullets: { text: string }[]
  /** 内置图形 key，或自定义图片的 data / remote URL。 */
  image: VisualKey | (string & {})
  imageSide: 'left' | 'right'
  badge: string
  ctaLabel: string
  ctaHref: string
}

export interface StatItem {
  value: string
  suffix: string
  label: string
}

export interface StatsProps {
  items: StatItem[]
  theme: 'dark' | 'paper'
  note: string
  /** 左上角小标题，留空则不显示。 */
  kicker?: string
}

export interface ProductCardItem {
  image: string
  label: string
  title: string
  summary: string
  /** 每行一个要点 */
  bullets: string
  linkLabel: string
  href: string
}

export interface ProductsProps {
  eyebrow: string
  title: string
  subtitle: string
  limit: number
  moreLabel: string
  moreHref: string
  /** 自定义卡片；为空时自动读取已发布产品。 */
  items?: ProductCardItem[]
  cardLabel?: string
  modalCtaLabel?: string
  modalCtaHref?: string
}

export interface NewsCardItem {
  image: string
  category: string
  date: string
  title: string
  summary: string
  readMinutes: number
  href: string
}

export interface NewsListProps {
  eyebrow: string
  title: string
  subtitle: string
  limit: number
  showFilters: boolean
  /** 自定义卡片；为空时自动读取已发布新闻。 */
  items?: NewsCardItem[]
}

export interface CaseCardItem {
  image: string
  industry: string
  client: string
  year: string
  metric: string
  metricSuffix: string
  metricLabel: string
  title: string
  summary: string
  href: string
}

export interface CasesProps {
  eyebrow: string
  title: string
  subtitle: string
  limit: number
  showFilters: boolean
  /** 自定义卡片；为空时自动读取已发布案例。 */
  items?: CaseCardItem[]
  modalNote?: string
  ctaLabel?: string
  ctaHref?: string
}

export interface FaqItem {
  q: string
  a: string
}

export interface FaqProps {
  eyebrow: string
  title: string
  subtitle: string
  items: FaqItem[]
  contactKicker?: string
  contactText?: string
  contactLabel?: string
  contactHref?: string
}

export interface CtaProps {
  title: string
  text: string
  primaryLabel: string
  primaryHref: string
  secondaryLabel: string
  secondaryHref: string
  theme: 'accent' | 'dark' | 'paper'
  kicker?: string
  /** 右侧竖排小字，可换行。 */
  sideNote?: string
  /** 自定义背景图；为空时使用主题底色。 */
  backgroundImage?: string
  /** 背景图遮罩强度，0–100。 */
  backgroundOverlay?: number
}

export interface ContactFormProps {
  title: string
  subtitle: string
  submitLabel: string
  successTitle: string
  successText: string
  email: string
  phone: string
  address: string
  hours: string
  /** 内置图形 key，或自定义图片 URL。 */
  mapStyle: VisualKey | (string & {})
  kicker?: string
  /** 图片上的坐标文案，留空则不显示。 */
  mapNote?: string
  formKicker?: string
  replyNote?: string
}

export interface RichTextProps {
  content: string
  width: 'narrow' | 'wide'
}

export interface VideoProps {
  title: string
  caption: string
  /** 内置图形 key，或自定义封面图片 URL。 */
  visual: VisualKey | (string & {})
  url: string
  /** 封面上方的播放标签，留空则不显示。 */
  label?: string
}

export interface HtmlProps {
  code: string
}

export interface BlockPropsMap {
  hero: HeroProps
  'hero-carousel': HeroCarouselProps
  partners: PartnersProps
  features: FeaturesProps
  'image-text': ImageTextProps
  stats: StatsProps
  products: ProductsProps
  'news-list': NewsListProps
  cases: CasesProps
  faq: FaqProps
  cta: CtaProps
  'contact-form': ContactFormProps
  'rich-text': RichTextProps
  video: VideoProps
  html: HtmlProps
}

export type BlockType = keyof BlockPropsMap

export type BlockInput = {
  [K in BlockType]: { type: K; props: BlockPropsMap[K] }
}[BlockType]

export type PageBlock = {
  [K in BlockType]: { id: string; type: K; props: BlockPropsMap[K] }
}[BlockType]

export type BlockOfType<K extends BlockType> = Extract<PageBlock, { type: K }>

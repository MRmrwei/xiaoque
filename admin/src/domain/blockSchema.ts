import type { BlockType, PageBlock } from '~/domain/blocks'

export type FieldType =
  | 'text'
  | 'textarea'
  | 'richtext'
  | 'number'
  | 'switch'
  | 'select'
  | 'color'
  | 'image'
  | 'images'
  | 'link'
  | 'date'
  | 'repeater'
  | 'group'
  | 'collection'
  | 'html'

export interface FieldOption {
  value: string
  label: string
}

export interface FieldCondition {
  /** 同级字段的 key。 */
  key: string
  equals: string | number | boolean
}

export interface BlockField {
  key: string
  label: string
  type: FieldType
  options?: FieldOption[]
  rows?: number
  help?: string
  itemLabel?: string
  itemFields?: BlockField[]
  min?: number
  max?: number
  required?: boolean
  default?: unknown
  /** 图片字段是否开放本地上传（默认只能选内置图形）。 */
  upload?: boolean
  /** 图片字段只允许上传，不展示内置图形。 */
  uploadOnly?: boolean
  /** 仅当同级字段等于指定值时显示。 */
  visibleWhen?: FieldCondition
}

export function isFieldVisible(
  field: BlockField,
  model: Record<string, unknown>
): boolean {
  if (!field.visibleWhen) return true
  return model[field.visibleWhen.key] === field.visibleWhen.equals
}

export interface BlockSchema {
  label: string
  labelEn: string
  icon: string
  group: string
  description: string
  /** 仅用于兼容旧数据，不在区块库里展示。 */
  hidden?: boolean
  defaults: Record<string, unknown>
  fields: BlockField[]
}

export interface BlockGroup {
  id: string
  label: string
  hint: string
}

let seq = 0
export function uid(prefix = 'b') {
  seq += 1
  return `${prefix}_${Date.now().toString(36)}_${seq}`
}

export function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

const field = (
  key: string,
  label: string,
  type: FieldType = 'text',
  extra: Partial<Omit<BlockField, 'key' | 'label' | 'type'>> = {}
): BlockField => ({ key, label, type, ...extra })

const iconOptions: FieldOption[] = [
  { value: 'spark', label: '星芒 / Spark' },
  { value: 'blocks', label: '区块 / Blocks' },
  { value: 'globe', label: '全球 / Globe' },
  { value: 'search', label: '搜索 / Search' },
  { value: 'history', label: '历史 / History' },
  { value: 'image', label: '图片 / Image' },
  { value: 'shield', label: '安全 / Shield' },
  { value: 'zap', label: '闪电 / Zap' },
  { value: 'brush', label: '画笔 / Brush' },
  { value: 'layers', label: '图层 / Layers' },
  { value: 'book', label: '书本 / Book' },
  { value: 'users', label: '用户 / Users' }
]

const artOptions: FieldOption[] = [
  { value: 'orbit', label: '轨道 Orbit' },
  { value: 'arch', label: '拱门 Arch' },
  { value: 'grid', label: '网格 Grid' },
  { value: 'dunes', label: '沙丘 Dunes' },
  { value: 'waves', label: '波浪 Waves' },
  { value: 'type', label: '字体 Type' },
  { value: 'stack', label: '层叠 Stack' },
  { value: 'bloom', label: '绽放 Bloom' }
]

export const artOptionsList = artOptions

export const blockGroups: BlockGroup[] = [
  { id: 'basic', label: '基础区块', hint: '页面骨架与叙事' },
  { id: 'content', label: '内容区块', hint: '数据与内容流' },
  { id: 'conversion', label: '转化区块', hint: '引导访客行动' }
]

export const blockSchemas: Record<BlockType, BlockSchema> = {
  hero: {
    label: '首屏横幅',
    labelEn: 'Hero',
    icon: 'spark',
    group: 'basic',
    description: '大标题 + 行动按钮，决定第一印象',
    defaults: {
      variant: 'page',
      eyebrow: '新的页面',
      title: '把想法\n变成<em>看得见</em>的页面。',
      subtitle: '在这里补充一段简洁有力的介绍文字，说明你提供的价值。',
      primaryLabel: '开始体验',
      primaryHref: '/contact',
      secondaryLabel: '了解更多',
      secondaryHref: '/products',
      note: '',
      visual: 'orbit',
      theme: 'light',
      sideNote: 'EST. 2019 / SHANGHAI',
      scrollLabel: 'SCROLL'
    },
    fields: [
      field('variant', '版式', 'select', {
        options: [
          { value: 'home', label: '首页版（含产品视觉）' },
          { value: 'page', label: '内页版（紧凑）' }
        ]
      }),
      field('eyebrow', '眉标', 'text'),
      field('title', '主标题', 'textarea', {
        rows: 3,
        help: '支持换行。用 <em>强调词</em> 包裹需要高亮的文字。'
      }),
      field('subtitle', '副标题', 'textarea', { rows: 3 }),
      field('primaryLabel', '主按钮文案'),
      field('primaryHref', '主按钮链接'),
      field('secondaryLabel', '次按钮文案'),
      field('secondaryHref', '次按钮链接'),
      field('note', '按钮下方小字'),
      field('visual', '右侧视觉', 'image', {
        options: artOptions.concat([
          { value: 'builder', label: '搭建器界面' },
          { value: 'none', label: '不显示' }
        ]),
        upload: true
      }),
      field('theme', '主题', 'select', {
        options: [
          { value: 'light', label: '暖白' },
          { value: 'dark', label: '墨黑' }
        ]
      }),
      field('sideNote', '右侧竖排小字', 'text', { help: '留空则不显示。' }),
      field('scrollLabel', '滚动提示文案', 'text', { help: '首页版左下角显示，留空则不显示。' })
    ]
  },

  'hero-carousel': {
    label: '头部轮播',
    labelEn: 'Hero carousel',
    icon: 'layers',
    group: 'basic',
    description: '图片/视频轮播首屏，支持自动播放与手动切换',
    defaults: {
      autoplay: true,
      interval: 7000,
      note: '左右滑动或按 ← → 切换 · 支持图片与视频',
      slides: [
        {
          id: 'slide_1',
          mediaType: 'image',
          media: 'builder',
          videoUrl: '',
          eyebrow: '小鹊云 3.0 · 企业内容云',
          title: '让每一次表达\n都成为<em>增长</em>。',
          subtitle: '在这里补充一段简洁有力的介绍文字。',
          primaryLabel: '预约演示',
          primaryHref: '/contact',
          secondaryLabel: '了解更多',
          secondaryHref: '/products'
        },
        {
          id: 'slide_2',
          mediaType: 'video',
          media: 'orbit',
          videoUrl: '/media/hero-demo.mp4',
          eyebrow: '视频 · VIDEO',
          title: '每一次发布，\n都能<em>安心上线</em>。',
          subtitle: '草稿与线上隔离，发布前可完整预览，确认后立即生效。',
          primaryLabel: '了解详情',
          primaryHref: '/products',
          secondaryLabel: '',
          secondaryHref: ''
        }
      ]
    },
    fields: [
      field('slides', '轮播内容', 'repeater', {
        itemLabel: '幻灯片',
        itemFields: [
          field('mediaType', '内容类型', 'select', {
            options: [
              { value: 'image', label: '图片' },
              { value: 'video', label: '视频' }
            ]
          }),
          field('media', '图片', 'image', {
            options: [{ value: 'builder', label: '搭建器界面' }, ...artOptions],
            upload: true,
            help: '仅在内容类型为图片时显示，作为轮播主视觉。',
            visibleWhen: { key: 'mediaType', equals: 'image' }
          }),
          field('videoUrl', '视频地址', 'text', {
            help: 'MP4 直链；仅在内容类型为视频时显示。',
            visibleWhen: { key: 'mediaType', equals: 'video' }
          }),
          field('videoPoster', '视频封面（可选）', 'image', {
            upload: true,
            uploadOnly: true,
            help: '视频播放前显示的封面图，上传后作为 video poster。',
            visibleWhen: { key: 'mediaType', equals: 'video' }
          }),
          field('eyebrow', '眉标'),
          field('title', '标题', 'textarea', { rows: 3, required: true }),
          field('subtitle', '副标题', 'textarea', { rows: 2 }),
          field('primaryLabel', '主按钮文案'),
          field('primaryHref', '主按钮链接'),
          field('secondaryLabel', '次按钮文案'),
          field('secondaryHref', '次按钮链接')
        ]
      }),
      field('autoplay', '自动播放', 'switch'),
      field('interval', '切换间隔（毫秒）', 'number', { min: 3000, max: 20000 }),
      field('note', '底部提示', 'text')
    ]
  },

  'rich-text': {
    label: '富文本',
    labelEn: 'Rich text',
    icon: 'type',
    group: 'basic',
    description: '自由排版的长文本内容',
    defaults: {
      content:
        '<p>在这里开始写作。富文本支持标题、列表、引用与链接，适合承载品牌故事、服务说明或政策条款。</p>',
      width: 'narrow'
    },
    fields: [
      field('content', '正文内容', 'richtext', {
        rows: 12,
        help: '可视化编辑标题、列表、引用、代码块与链接，发布后按文章样式渲染。'
      }),
      field('width', '内容宽度', 'select', {
        options: [
          { value: 'narrow', label: '窄栏（阅读）' },
          { value: 'wide', label: '宽栏（页面）' }
        ]
      })
    ]
  },

  'image-text': {
    label: '图文组合',
    labelEn: 'Image + text',
    icon: 'image',
    group: 'basic',
    description: '一张视觉、一段叙述，左右布局',
    defaults: {
      eyebrow: '我们的方式',
      title: '把复杂留给自己，\n把简单留给团队。',
      body: '用一段话讲清楚这件事为什么重要，以及你如何做到。',
      bullets: [{ text: '第一个要点，短而具体' }, { text: '第二个要点，给出证据' }],
      image: 'arch',
      imageSide: 'right',
      badge: '',
      ctaLabel: '',
      ctaHref: '/contact'
    },
    fields: [
      field('eyebrow', '眉标'),
      field('title', '标题', 'textarea', { rows: 2 }),
      field('body', '正文', 'textarea', { rows: 4 }),
      field('bullets', '要点列表', 'repeater', {
        itemLabel: '要点',
        itemFields: [field('text', '文本')]
      }),
      field('image', '视觉图片', 'image', {
        options: artOptions,
        upload: true,
        help: '可选用内置图形，或上传自己的图片；上传后自动压缩，仅保存在本地。'
      }),
      field('imageSide', '图片位置', 'select', {
        options: [
          { value: 'right', label: '右侧' },
          { value: 'left', label: '左侧' }
        ]
      }),
      field('badge', '图片角标文案'),
      field('ctaLabel', '按钮文案'),
      field('ctaHref', '按钮链接')
    ]
  },

  features: {
    label: '特性列表',
    labelEn: 'Features',
    icon: 'blocks',
    group: 'basic',
    description: '用网格展示核心能力',
    defaults: {
      eyebrow: '核心能力',
      title: '让每一个能力都恰到好处。',
      subtitle: '从搭建到发布，每一步都为运营团队设计。',
      layout: 'grid',
      items: [
        { icon: 'blocks', title: '区块化搭建', desc: '拖拽组合预置区块，无需代码。', tag: '' },
        { icon: 'globe', title: '中英双语', desc: '一套后台，两种语言独立发布。', tag: '' },
        { icon: 'search', title: 'SEO 就绪', desc: 'SSR 与 hreflang 开箱即用。', tag: '' }
      ]
    },
    fields: [
      field('eyebrow', '眉标'),
      field('title', '标题', 'textarea', { rows: 2 }),
      field('subtitle', '副标题', 'textarea', { rows: 2 }),
      field('layout', '布局', 'select', {
        options: [
          { value: 'bento', label: 'Bento 混排' },
          { value: 'grid', label: '等宽网格' }
        ]
      }),
      field('items', '特性条目', 'repeater', {
        itemLabel: '特性',
        itemFields: [
          field('icon', '图标', 'select', { options: iconOptions }),
          field('title', '标题'),
          field('desc', '描述', 'textarea', { rows: 2 }),
          field('tag', '角标（选填）')
        ]
      })
    ]
  },

  stats: {
    label: '数据统计',
    labelEn: 'Stats',
    icon: 'chart',
    group: 'basic',
    description: '用数字建立信任',
    defaults: {
      items: [
        { value: '60', suffix: 's', label: '发布生效' },
        { value: '98', suffix: '%', label: '客户续约率' },
        { value: '3', suffix: '×', label: '内容效率提升' },
        { value: '14', suffix: '', label: '预置区块' }
      ],
      theme: 'dark',
      note: '数据来自 2025–2026 年客户调研',
      kicker: 'BY THE NUMBERS'
    },
    fields: [
      field('kicker', '左上角小标题', 'text', { help: '留空则不显示。' }),
      field('items', '统计项', 'repeater', {
        itemLabel: '统计',
        itemFields: [
          field('value', '数值'),
          field('suffix', '后缀'),
          field('label', '说明')
        ]
      }),
      field('theme', '主题', 'select', {
        options: [
          { value: 'dark', label: '墨黑' },
          { value: 'paper', label: '暖白' }
        ]
      }),
      field('note', '底部注释')
    ]
  },

  products: {
    label: '产品列表',
    labelEn: 'Products',
    icon: 'layers',
    group: 'content',
    description: '自动读取已发布的产品内容',
    defaults: {
      eyebrow: '产品矩阵',
      title: '一整套内容能力。',
      subtitle: '按需组合，随业务生长。',
      items: [],
      limit: 3,
      cardLabel: '了解产品',
      moreLabel: '查看全部产品',
      moreHref: '/products',
      modalCtaLabel: '预约演示',
      modalCtaHref: '/contact'
    },
    fields: [
      field('eyebrow', '眉标'),
      field('title', '标题', 'textarea', { rows: 2 }),
      field('subtitle', '副标题', 'textarea', { rows: 2 }),
      field('items', '自定义卡片', 'repeater', {
        itemLabel: '产品卡片',
        help: '留空时自动读取已发布的产品内容；添加后按这里的图片、文案和链接展示。',
        itemFields: [
          field('image', '卡片图片', 'image', { options: artOptions, upload: true }),
          field('label', '角标 / 分类'),
          field('title', '标题'),
          field('summary', '简介', 'textarea', { rows: 2 }),
          field('bullets', '要点（每行一个）', 'textarea', { rows: 3 }),
          field('linkLabel', '按钮文案', 'text', { help: '留空使用上方的「卡片按钮文案」。' }),
          field('href', '跳转链接', 'link')
        ]
      }),
      field('limit', '显示条数', 'number', {
        min: 1,
        max: 6,
        help: '自动读取和自定义卡片都按此数量显示，取排序后的前 N 条。'
      }),
      field('cardLabel', '卡片按钮文案'),
      field('moreLabel', '更多按钮文案'),
      field('moreHref', '更多按钮链接'),
      field('modalCtaLabel', '详情按钮文案（自动读取时）'),
      field('modalCtaHref', '详情按钮链接（自动读取时）')
    ]
  },

  'news-list': {
    label: '新闻列表',
    labelEn: 'News list',
    icon: 'book',
    group: 'content',
    description: '按语言自动读取新闻内容',
    defaults: {
      eyebrow: '新闻动态',
      title: '小鹊最近在做什么。',
      subtitle: '',
      items: [],
      limit: 3,
      showFilters: false
    },
    fields: [
      field('eyebrow', '眉标'),
      field('title', '标题', 'textarea', { rows: 2 }),
      field('subtitle', '副标题', 'textarea', { rows: 2 }),
      field('items', '自定义卡片', 'repeater', {
        itemLabel: '新闻卡片',
        help: '留空时自动读取已发布的新闻内容；添加后按这里的图片、文案和链接展示。',
        itemFields: [
          field('image', '卡片图片', 'image', { options: artOptions, upload: true }),
          field('category', '分类'),
          field('date', '日期', 'date'),
          field('title', '标题'),
          field('summary', '摘要', 'textarea', { rows: 2 }),
          field('readMinutes', '阅读时长（分钟）', 'number', { min: 1, max: 60 }),
          field('href', '跳转链接', 'link')
        ]
      }),
      field('limit', '显示条数', 'number', {
        min: 1,
        max: 12,
        help: '自动读取和自定义卡片都按此数量显示，取排序后的前 N 条。'
      }),
      field('showFilters', '显示分类筛选', 'switch')
    ]
  },

  cases: {
    label: '客户案例',
    labelEn: 'Case studies',
    icon: 'star',
    group: 'content',
    description: '用客户成果证明价值',
    defaults: {
      eyebrow: '客户案例',
      title: '真实增长的数字。',
      subtitle: '',
      items: [],
      limit: 3,
      showFilters: false,
      modalNote: '',
      ctaLabel: '聊聊类似的项目',
      ctaHref: '/contact'
    },
    fields: [
      field('eyebrow', '眉标'),
      field('title', '标题', 'textarea', { rows: 2 }),
      field('subtitle', '副标题', 'textarea', { rows: 2 }),
      field('items', '自定义卡片', 'repeater', {
        itemLabel: '案例卡片',
        help: '留空时自动读取已发布的客户案例；添加后按这里的图片、文案和链接展示，点击卡片可打开详情弹层。',
        itemFields: [
          field('image', '卡片图片', 'image', { options: artOptions, upload: true }),
          field('industry', '行业'),
          field('client', '客户'),
          field('year', '年份'),
          field('metric', '核心数字'),
          field('metricSuffix', '数字后缀', 'text', { help: '如 %、×、s' }),
          field('metricLabel', '数字说明'),
          field('title', '案例标题', 'text', { help: '详情弹层标题；留空时使用客户名称。' }),
          field('summary', '摘要', 'textarea', { rows: 2 }),
          field('href', '案例链接', 'link', { help: '详情弹层中的「查看完整案例」按钮；可留空。' })
        ]
      }),
      field('limit', '显示条数', 'number', {
        min: 1,
        max: 9,
        help: '自动读取和自定义卡片都按此数量显示，取排序后的前 N 条。'
      }),
      field('showFilters', '显示行业筛选', 'switch'),
      field('modalNote', '详情补充说明', 'textarea', {
        rows: 3,
        help: '点击案例卡片打开的详情弹层中显示，两种数据来源都生效。'
      }),
      field('ctaLabel', '详情按钮文案', 'text', { help: '详情弹层底部主按钮的文案。' }),
      field('ctaHref', '详情按钮链接', 'link', { help: '详情弹层底部主按钮的链接。' })
    ]
  },

  partners: {
    label: '合作伙伴',
    labelEn: 'Partners',
    icon: 'globe',
    group: 'content',
    description: '滚动展示客户与伙伴 Logo',
    defaults: {
      title: '他们选择与小鹊一起生长',
      note: '260+ 家企业正在使用小鹊云',
      marquee: true,
      items: [{ name: 'NOVA' }, { name: '山丘制造' }, { name: 'LUMEN' }, { name: '栖木设计' }, { name: 'HORIZON' }, { name: '云启资本' }]
    },
    fields: [
      field('title', '标题'),
      field('note', '注释'),
      field('marquee', '滚动展示', 'switch'),
      field('items', '名称', 'repeater', {
        itemLabel: '伙伴',
        itemFields: [field('name', '名称')]
      })
    ]
  },

  faq: {
    label: '常见问题',
    labelEn: 'FAQ',
    icon: 'quote',
    group: 'content',
    description: '回答访客最关心的问题',
    defaults: {
      eyebrow: '常见问题',
      title: '你可能想知道。',
      subtitle: '',
      items: [
        { q: '我的团队不会写代码，能用吗？', a: '可以。所有区块都通过表单配置，运营人员无需接触代码。' },
        { q: '可以从现有官网迁移吗？', a: '可以。我们会协助梳理栏目与内容，并提供批量导入工具。' }
      ],
      contactKicker: 'STILL CURIOUS?',
      contactText: '还有疑问？我们很乐意直接回答。',
      contactLabel: '联系我们',
      contactHref: '/contact'
    },
    fields: [
      field('eyebrow', '眉标'),
      field('title', '标题', 'textarea', { rows: 2 }),
      field('subtitle', '副标题', 'textarea', { rows: 2 }),
      field('items', '问题列表', 'repeater', {
        itemLabel: '问题',
        itemFields: [
          field('q', '问题'),
          field('a', '回答', 'textarea', { rows: 3 })
        ]
      }),
      field('contactKicker', '侧栏小标题'),
      field('contactText', '侧栏说明', 'textarea', { rows: 2 }),
      field('contactLabel', '侧栏按钮文案'),
      field('contactHref', '侧栏按钮链接')
    ]
  },

  cta: {
    label: '行动号召',
    labelEn: 'Call to action',
    icon: 'zap',
    group: 'conversion',
    description: '页面末尾的强引导',
    defaults: {
      title: '准备好，让官网\n开始为你工作了吗？',
      text: '预约一次 30 分钟的产品演示，看看小鹊云如何适配你的团队。',
      primaryLabel: '预约演示',
      primaryHref: '/contact',
      secondaryLabel: '查看产品',
      secondaryHref: '/products',
      theme: 'accent',
      kicker: "LET'S TALK",
      sideNote: '30 MIN\nPRODUCT DEMO',
      backgroundImage: '',
      backgroundOverlay: 65
    },
    fields: [
      field('title', '标题', 'textarea', { rows: 2 }),
      field('text', '说明', 'textarea', { rows: 2 }),
      field('primaryLabel', '主按钮文案'),
      field('primaryHref', '主按钮链接'),
      field('secondaryLabel', '次按钮文案'),
      field('secondaryHref', '次按钮链接'),
      field('theme', '主题底色', 'select', {
        options: [
          { value: 'accent', label: '朱砂红' },
          { value: 'dark', label: '墨黑' },
          { value: 'paper', label: '暖白' }
        ]
      }),
      field('kicker', '上方小标题', 'text', { help: '留空则不显示。' }),
      field('sideNote', '右侧竖排小字', 'textarea', { rows: 2, help: '支持换行，留空则不显示。' }),
      field('backgroundImage', '背景图片', 'image', {
        upload: true,
        uploadOnly: true,
        help: '上传后作为底图并覆盖主题底色；留空则使用主题色。'
      }),
      field('backgroundOverlay', '图片遮罩', 'number', {
        min: 0,
        max: 90,
        help: '0–90，数值越大文字越清晰。'
      })
    ]
  },

  'contact-form': {
    label: '联系表单',
    labelEn: 'Contact form',
    icon: 'form',
    group: 'conversion',
    description: '收集留言与线索',
    defaults: {
      title: '聊聊你的官网。',
      subtitle: '留下联系方式，我们会在一个工作日内回复。',
      submitLabel: '发送消息',
      successTitle: '消息已经飞过来了',
      successText: '感谢你的联系，我们会在一个工作日内回复。',
      email: 'hello@xiaoque.com',
      phone: '400-860-2026',
      address: '上海市徐汇区西岸智塔 18F',
      hours: '周一至周五 9:00 – 18:00',
      mapStyle: 'orbit',
      kicker: 'CONTACT / 01',
      mapNote: 'SHANGHAI · 31.2304° N, 121.4737° E',
      formKicker: 'PROJECT INQUIRY'
    },
    fields: [
      field('title', '标题', 'textarea', { rows: 2 }),
      field('subtitle', '副标题', 'textarea', { rows: 2 }),
      field('submitLabel', '提交按钮文案'),
      field('successTitle', '成功标题'),
      field('successText', '成功说明', 'textarea', { rows: 2 }),
      field('email', '邮箱'),
      field('phone', '电话'),
      field('address', '地址'),
      field('hours', '工作时间'),
      field('kicker', '左上角小标题', 'text', { help: '留空则不显示。' }),
      field('mapStyle', '右侧图片', 'image', { options: artOptions, upload: true }),
      field('mapNote', '图片上的坐标文案', 'text', { help: '留空则不显示。' }),
      field('formKicker', '表单上方小标题', 'text', { help: '留空则不显示。' }),
      field('replyNote', '回复时效提示', 'text', { help: '留空则不显示。' })
    ]
  },

  video: {
    label: '视频',
    labelEn: 'Video',
    icon: 'play',
    group: 'content',
    description: '品牌片或产品演示',
    defaults: {
      title: '一分钟，看懂小鹊云。',
      caption: '产品演示 · 01:12',
      visual: 'bloom',
      url: '',
      label: 'PLAY FILM · 01:12'
    },
    fields: [
      field('title', '标题'),
      field('caption', '说明文字'),
      field('visual', '封面图片', 'image', { options: artOptions, upload: true }),
      field('label', '封面播放标签', 'text', { help: '留空则不显示。' }),
      field('url', '视频链接')
    ]
  },

  html: {
    label: '自定义 HTML（旧版）',
    labelEn: 'Legacy HTML',
    icon: 'edit',
    group: 'conversion',
    description: '旧版 HTML 区块，已并入富文本；仅用于兼容历史数据',
    hidden: true,
    defaults: {
      code: '<p>旧版自定义 HTML</p>'
    },
    fields: [field('code', '正文内容', 'richtext', { rows: 12 })]
  }
}

export const blockList = Object.entries(blockSchemas)
  .filter(([, schema]) => !schema.hidden)
  .map(([type, schema]) => ({
    type,
    ...schema
  }))

export function createBlock(type: BlockType): PageBlock | null {
  const schema = blockSchemas[type]
  if (!schema) return null
  return {
    id: uid('blk'),
    type,
    props: clone(schema.defaults)
  } as unknown as PageBlock
}

export function blockLabel(type: BlockType | string) {
  return blockSchemas[type as BlockType]?.label || type
}

export function blockSummary(block: { type?: string; props?: object }): string {
  const p = (block?.props ?? {}) as Record<string, unknown>
  const items = Array.isArray(p.items) ? (p.items as Array<{ label?: string }>) : []
  const candidate = p.title || p.eyebrow || p.content || p.q || items[0]?.label || ''
  if (!candidate) return blockLabel(block?.type ?? '')
  return String(candidate)
    .replace(/<[^>]+>/g, '')
    .replace(/\n/g, ' ')
    .slice(0, 34)
}

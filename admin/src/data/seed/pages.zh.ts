import { b, makePage } from './page-factory'

const home = makePage({
  id: 'pg_home_zh',
  group: 'tg_home',
  locale: 'zh',
  path: '/',
  name: '首页',
  version: 12,
  updatedAt: '2026-10-14T16:20:00+08:00',
  publishedAt: '2026-10-14T16:20:00+08:00',
  seo: {
    title: '小鹊云 · 让品牌官网成为增长引擎',
    description:
      '小鹊云是面向成长型企业的区块化官网搭建与内容云平台。运营自主搭建、中英双语同步、SEO 就绪，发布即刻生效。',
    keywords: ['企业官网', '内容管理平台', '中英双语官网', '区块化建站', '小鹊云']
  },
  blocks: [
    b('hero-carousel', {
      autoplay: true,
      interval: 7000,
      note: '左右滑动或按 ← → 切换 · 支持图片与视频',
      slides: [
        {
          id: 'zh_carousel_1',
          mediaType: 'image',
          media: 'builder',
          videoUrl: '',
          eyebrow: '小鹊云 3.0 · 企业内容云',
          title: '让每一次表达\n都成为<em>增长</em>。',
          subtitle:
            '小鹊云把品牌官网变成一套会生长的内容系统：运营自主搭建、中英双语同步、发布即刻生效。',
          primaryLabel: '预约产品演示',
          primaryHref: '/contact',
          secondaryLabel: '探索产品能力',
          secondaryHref: '/products'
        },
        {
          id: 'zh_carousel_2',
          mediaType: 'image',
          media: 'arch',
          videoUrl: '',
          eyebrow: '一套内容系统 · CONTENT CLOUD',
          title: '从内容到发布，\n只隔一个<em>念头</em>。',
          subtitle: '页面、内容、媒体与发布流程，在同一个系统里完成。',
          primaryLabel: '看看产品能力',
          primaryHref: '/products',
          secondaryLabel: '预约演示',
          secondaryHref: '/contact'
        },
        {
          id: 'zh_carousel_3',
          mediaType: 'video',
          media: 'orbit',
          videoUrl:
            '/media/hero-demo.mp4',
          eyebrow: '发布与生效 · SAFE PUBLISH',
          title: '每一次发布，\n都能<em>安心上线</em>。',
          subtitle: '草稿与线上隔离，发布前可完整预览，确认后立即生效。',
          primaryLabel: '了解发布机制',
          primaryHref: '/products',
          secondaryLabel: '查看客户案例',
          secondaryHref: '/cases'
        }
      ]
    }),
    b('partners', {
      title: '他们选择与小鹊一起生长',
      note: '260+ 家企业正在使用小鹊云',
      marquee: true,
      items: [
        { name: 'NOVA' },
        { name: '山丘制造' },
        { name: 'LUMEN' },
        { name: '栖木设计' },
        { name: 'HORIZON' },
        { name: '云启资本' },
        { name: '未迟' },
        { name: 'ARK' }
      ]
    }),
    b('features', {
      eyebrow: '为什么选择小鹊',
      title: '把复杂留给我们，\n把简单留给团队。',
      subtitle: '从内容搭建到全球发布，每一个细节都为「运营能自己完成」而设计。',
      layout: 'bento',
      items: [
        {
          icon: 'blocks',
          title: '区块化搭建',
          desc: '14 种预置区块，像搭积木一样组合页面。响应式与品牌规范由组件保证。',
          tag: '核心能力'
        },
        {
          icon: 'globe',
          title: '中英双语',
          desc: '一套后台，两种语言。独立草稿、独立发布，互不覆盖。',
          tag: ''
        },
        {
          icon: 'search',
          title: 'SEO 就绪',
          desc: 'SSR 直出、hreflang、双语 sitemap 与结构化数据，开箱即用。',
          tag: ''
        },
        {
          icon: 'history',
          title: '发布与生效',
          desc: '草稿与线上隔离，发布前完整预览，确认后立即生效。',
          tag: ''
        },
        {
          icon: 'image',
          title: '媒体资产库',
          desc: '上传、压缩、双语 alt 与引用检查，一站式管理。',
          tag: ''
        },
        {
          icon: 'shield',
          title: '权限与审计',
          desc: '按钮级权限、预览令牌、操作全留痕。',
          tag: ''
        }
      ]
    }),
    b('image-text', {
      eyebrow: '给运营的自由',
      title: '改了十几版？\n这一次，不用再等排期。',
      body:
        '可视化编辑器让运营团队直接完成页面搭建。三栏布局、实时预览、自动保存，调整文案与图片不再需要排进开发周期。',
      bullets: [
        { text: '拖拽排序与区块库，像整理文档一样搭页面' },
        { text: '桌面 / 平板 / 手机预览，所见即所得' },
        { text: '草稿与线上隔离，确认后再发布' }
      ],
      image: 'stack',
      imageSide: 'right',
      badge: 'LIVE PREVIEW',
      ctaLabel: '看看怎么搭',
      ctaHref: '/products'
    }),
    b('stats', {
      theme: 'dark',
      note: '数据来自 2025–2026 年客户调研',
      items: [
        { value: '60', suffix: 's', label: '发布生效' },
        { value: '98', suffix: '%', label: '客户续约率' },
        { value: '3', suffix: '×', label: '内容效率提升' },
        { value: '14', suffix: '', label: '预置区块' }
      ]
    }),
    b('products', {
      eyebrow: '产品矩阵',
      title: '不止于建站，\n是一整套内容能力。',
      subtitle: '从页面搭建到内容洞察，按需组合，随业务生长。',
      limit: 3,
      moreLabel: '查看全部产品',
      moreHref: '/products'
    }),
    b('cases', {
      eyebrow: '客户案例',
      title: '真实增长，\n写在客户的数字里。',
      subtitle: '覆盖智能制造、专业服务与消费品牌。',
      limit: 3,
      showFilters: false
    }),
    b('news-list', {
      eyebrow: '新闻动态',
      title: '小鹊最近在做什么。',
      subtitle: '',
      limit: 3,
      showFilters: false
    }),
    b('faq', {
      eyebrow: '常见问题',
      title: '你可能想知道。',
      subtitle: '',
      items: [
        {
          q: '我的团队不会写代码，能用吗？',
          a: '可以。所有区块都通过表单配置，运营人员无需接触代码。拖拽、填写、预览、发布，就是全部操作。'
        },
        {
          q: '可以从现有官网迁移吗？',
          a: '可以。我们会协助梳理栏目与内容结构，并提供批量导入工具。通常 3–7 天可以完成一次完整迁移。'
        },
        {
          q: '中英文内容必须同时上线吗？',
          a: '不需要。每种语言都有独立草稿与发布版本，可以先上线中文，再完善英文。未翻译的页面不会出现语言切换入口。'
        },
        {
          q: '支持自定义域名和 ICP 备案吗？',
          a: '支持绑定自有域名，并提供备案指引。小鹊云使用国内节点，也可按需部署到海外区域。'
        },
        {
          q: '数据安全如何保障？',
          a: '内容每日备份，媒体资源统一管理；后台支持按钮级权限与操作审计，发布操作全程可追溯。'
        }
      ]
    }),
    b('cta', {
      title: '准备好，让官网\n开始为你工作了吗？',
      text: '预约一次 30 分钟的产品演示，看看小鹊云如何适配你的团队。',
      primaryLabel: '预约演示',
      primaryHref: '/contact',
      secondaryLabel: '查看产品',
      secondaryHref: '/products',
      theme: 'accent'
    })
  ]
})

const products = makePage({
  id: 'pg_products_zh',
  group: 'tg_products',
  locale: 'zh',
  path: '/products',
  name: '产品',
  version: 7,
  updatedAt: '2026-10-09T11:40:00+08:00',
  publishedAt: '2026-10-09T11:40:00+08:00',
  seo: {
    title: '产品能力 · 小鹊云',
    description: '区块化页面搭建、双语内容管理、版本发布与媒体资产库，一套后台管理整座官网。',
    keywords: ['页面搭建', '内容管理', '双语官网', '小鹊云产品']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: '产品 / PRODUCTS',
      title: '一套后台，\n装下你的整座官网。',
      subtitle:
        '小鹊云由三个相互协作的产品组成：搭建、内容与洞察。可以一起用，也可以从最需要的那一块开始。',
      primaryLabel: '预约演示',
      primaryHref: '/contact',
      secondaryLabel: '查看客户案例',
      secondaryHref: '/cases',
      note: '',
      visual: 'orbit',
      theme: 'light'
    }),
    b('features', {
      eyebrow: '能力全景',
      title: '从第一行文案，\n到每一次发布。',
      subtitle: '把重复劳动交给系统，把判断留给人。',
      layout: 'grid',
      items: [
        {
          icon: 'blocks',
          title: '区块化页面搭建',
          desc: '拖拽组合 14 种预置区块，属性通过表单配置，响应式由组件保证。',
          tag: ''
        },
        {
          icon: 'globe',
          title: '双语内容管理',
          desc: '页面、内容、导航与站点设置均按语言独立管理，翻译状态一目了然。',
          tag: ''
        },
        {
          icon: 'history',
          title: '版本与发布',
          desc: '草稿隔离、发布前预览、直接上线，发布不再是一场冒险。',
          tag: ''
        },
        {
          icon: 'image',
          title: '媒体资产库',
          desc: '拖拽上传、自动压缩、双语 alt、引用检查，图片不再散落各处。',
          tag: ''
        },
        {
          icon: 'form',
          title: '表单与线索',
          desc: '内置联系表单，提交入库并支持邮件通知，数据可导出 CSV。',
          tag: ''
        },
        {
          icon: 'shield',
          title: '权限与审计',
          desc: '编辑与发布角色分离，按钮级权限，关键操作全量留痕。',
          tag: ''
        }
      ]
    }),
    b('products', {
      eyebrow: '三个产品',
      title: '选择适合你的起点。',
      subtitle: '每个产品都可以独立使用，也能无缝组合。',
      limit: 3,
      moreLabel: '',
      moreHref: ''
    }),
    b('image-text', {
      eyebrow: '发布流程',
      title: '发布，\n不再是一场冒险。',
      body:
        '草稿与线上完全隔离。发布前可以完整预览，确认后当前内容立即覆盖线上；每次发布都会写入操作审计。',
      bullets: [
        { text: '一键发布当前语言，或同时发布中英文' },
        { text: '每次发布记录操作人与时间' },
        { text: '发布前自动校验必填项、失效图片与重复路径' }
      ],
      image: 'arch',
      imageSide: 'right',
      badge: 'SAFE PUBLISH',
      ctaLabel: '了解发布机制',
      ctaHref: '/contact'
    }),
    b('faq', {
      eyebrow: '产品问答',
      title: '上线前，\n你可能还会问。',
      subtitle: '',
      items: [
        {
          q: '可以只购买其中一部分能力吗？',
          a: '可以。小鹊建站适合只需要品牌官网的团队；需要内容中台、多语言或表单线索管理时，再升级到小鹊云 CMS。'
        },
        {
          q: '现有的内容可以导入吗？',
          a: '支持从 Excel、Word 与常见 CMS 批量导入，并提供字段映射工具，迁移过程不影响旧站运行。'
        },
        {
          q: '发布后多久生效？',
          a: '发布后 60 秒内全站与 CDN 生效。服务端使用版本号 + ETag 缓存，并在发布时主动清理。'
        },
        {
          q: '能否接入我们自己的设计规范？',
          a: '可以。区块库支持品牌色、字体与圆角等设计变量；更深的定制可由我们基于区块协议扩展专属组件。'
        }
      ]
    }),
    b('cta', {
      title: '想看看它\n在你的业务里怎么用？',
      text: '我们会结合你的栏目结构，准备一场有针对性的演示。',
      primaryLabel: '预约演示',
      primaryHref: '/contact',
      secondaryLabel: '查看客户案例',
      secondaryHref: '/cases',
      theme: 'dark'
    })
  ]
})

const solutions = makePage({
  id: 'pg_solutions_zh',
  group: 'tg_solutions',
  locale: 'zh',
  path: '/solutions',
  name: '解决方案',
  version: 5,
  updatedAt: '2026-10-08T14:10:00+08:00',
  publishedAt: '2026-10-08T14:10:00+08:00',
  seo: {
    title: '行业解决方案 · 小鹊云',
    description: '面向智能制造、专业服务、消费品牌与内容团队的官网与内容解决方案。',
    keywords: ['官网解决方案', '行业网站', '智能制造官网', '品牌官网']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: '解决方案 / SOLUTIONS',
      title: '为你的行业，\n配一套合适的官网。',
      subtitle:
        '我们服务制造、专业服务与消费品牌，沉淀了可复用的栏目结构、内容模型与翻译工作流。',
      primaryLabel: '预约演示',
      primaryHref: '/contact',
      secondaryLabel: '查看客户案例',
      secondaryHref: '/cases',
      note: '',
      visual: 'grid',
      theme: 'light'
    }),
    b('features', {
      eyebrow: '按行业',
      title: '懂行业，\n才谈得上合适。',
      subtitle: '不是套模板，而是从真实业务出发的页面结构。',
      layout: 'grid',
      items: [
        {
          icon: 'layers',
          title: '智能制造',
          desc: '产品参数、解决方案、经销网络与海外询盘，支持多语言产品目录。',
          tag: '制造'
        },
        {
          icon: 'users',
          title: '专业服务',
          desc: '合伙人介绍、观点文章与信任背书，让专业能力被准确看见。',
          tag: '服务'
        },
        {
          icon: 'spark',
          title: '消费品牌',
          desc: '新品发布、品牌故事与会员触点，保持每季更新的节奏感。',
          tag: '品牌'
        },
        {
          icon: 'book',
          title: '内容团队',
          desc: '多语言内容工作流、媒体资产与发布日历，让协作有章可循。',
          tag: '内容'
        }
      ]
    }),
    b('image-text', {
      eyebrow: '可复用的内容模型',
      title: '结构可以复制，\n表达不必重复。',
      body:
        '页面、内容与分类都按语言成对存储。创建英文版本时只复制结构与不可翻译字段，文案留空交给译者，避免中英内容互相覆盖。',
      bullets: [
        { text: '同一条内容的中英文版本互为翻译组' },
        { text: '数据源区块自动按当前语言过滤' },
        { text: '语言切换只指向已翻译且已发布的页面' }
      ],
      image: 'waves',
      imageSide: 'left',
      badge: 'ZH / EN',
      ctaLabel: '了解多语言能力',
      ctaHref: '/products'
    }),
    b('cases', {
      eyebrow: '落地成果',
      title: '已经发生的改变。',
      subtitle: '',
      limit: 3,
      showFilters: true
    }),
    b('stats', {
      theme: 'paper',
      note: '',
      items: [
        { value: '12', suffix: '', label: '覆盖行业' },
        { value: '260', suffix: '+', label: '服务客户' },
        { value: '8', suffix: '', label: '城市团队' },
        { value: '99.98', suffix: '%', label: '服务可用性' }
      ]
    }),
    b('cta', {
      title: '你的行业，\n值得一套更好的官网。',
      text: '告诉我们你的业务与目标，我们会在一个工作日内给出初步建议。',
      primaryLabel: '联系我们',
      primaryHref: '/contact',
      secondaryLabel: '查看产品能力',
      secondaryHref: '/products',
      theme: 'accent'
    })
  ]
})

const cases = makePage({
  id: 'pg_cases_zh',
  group: 'tg_cases',
  locale: 'zh',
  path: '/cases',
  name: '客户案例',
  version: 9,
  updatedAt: '2026-10-11T09:30:00+08:00',
  publishedAt: '2026-10-11T09:30:00+08:00',
  seo: {
    title: '客户案例 · 小鹊云',
    description: '智能制造、专业服务、消费品牌等行业的官网与内容平台实践。',
    keywords: ['客户案例', '官网改版案例', '企业官网案例']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: '客户案例 / CUSTOMERS',
      title: '增长不是形容词，\n是一组数字。',
      subtitle: '看看不同行业的团队，如何用小鹊云把官网变成业务的一部分。',
      primaryLabel: '预约演示',
      primaryHref: '/contact',
      secondaryLabel: '了解产品',
      secondaryHref: '/products',
      note: '',
      visual: 'bloom',
      theme: 'light'
    }),
    b('cases', {
      eyebrow: '按行业浏览',
      title: '每个行业，\n都有相似的答案。',
      subtitle: '',
      limit: 6,
      showFilters: true
    }),
    b('stats', {
      theme: 'dark',
      note: '客户上线后 6 个月平均数据',
      items: [
        { value: '186', suffix: '%', label: '海外询盘增长' },
        { value: '72', suffix: '%', label: '品牌搜索增长' },
        { value: '3.2', suffix: '×', label: '内容产出效率' },
        { value: '20', suffix: 'min', label: '平均发布周期' }
      ]
    }),
    b('cta', {
      title: '下一个案例，\n也许就是你的。',
      text: '聊聊你的目标，我们来看看能一起做点什么。',
      primaryLabel: '预约演示',
      primaryHref: '/contact',
      secondaryLabel: '返回产品',
      secondaryHref: '/products',
      theme: 'accent'
    })
  ]
})

const news = makePage({
  id: 'pg_news_zh',
  group: 'tg_news',
  locale: 'zh',
  path: '/news',
  name: '新闻动态',
  version: 4,
  updatedAt: '2026-10-13T10:15:00+08:00',
  publishedAt: '2026-10-13T10:15:00+08:00',
  seo: {
    title: '新闻动态 · 小鹊云',
    description: '小鹊科技的产品更新、客户故事，以及关于内容、设计与增长的思考。',
    keywords: ['小鹊科技新闻', '产品更新', '内容运营']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: '新闻动态 / JOURNAL',
      title: '小鹊在做什么，\n以及我们在想什么。',
      subtitle: '产品更新、客户故事，以及关于内容、设计与增长的一些思考。',
      primaryLabel: '',
      primaryHref: '',
      secondaryLabel: '',
      secondaryHref: '',
      note: '',
      visual: 'none',
      theme: 'light'
    }),
    b('news-list', {
      eyebrow: '所有文章',
      title: '',
      subtitle: '',
      limit: 9,
      showFilters: true
    })
  ]
})

const about = makePage({
  id: 'pg_about_zh',
  group: 'tg_about',
  locale: 'zh',
  path: '/about',
  name: '关于我们',
  version: 6,
  updatedAt: '2026-10-06T17:00:00+08:00',
  publishedAt: '2026-10-06T17:00:00+08:00',
  seo: {
    title: '关于小鹊 · 小鹊科技',
    description: '小鹊科技成立于 2019 年，为成长型企业打造更轻、更自由的内容基础设施。',
    keywords: ['关于小鹊', '小鹊科技', '公司介绍']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: '关于小鹊 / ABOUT',
      title: '我们相信，\n好内容值得被看见。',
      subtitle:
        '小鹊科技成立于 2019 年。我们为成长型企业打造更轻、更自由的内容基础设施，让好内容不被流程耽误。',
      primaryLabel: '加入我们',
      primaryHref: '/contact',
      secondaryLabel: '看看客户',
      secondaryHref: '/cases',
      note: '',
      visual: 'arch',
      theme: 'light'
    }),
    b('stats', {
      theme: 'paper',
      note: '',
      items: [
        { value: '2019', suffix: '', label: '成立年份' },
        { value: '42', suffix: '', label: '团队成员' },
        { value: '12', suffix: '', label: '覆盖行业' },
        { value: '260', suffix: '+', label: '服务客户' }
      ]
    }),
    b('image-text', {
      eyebrow: '我们的故事',
      title: '从一只喜鹊开始。',
      body:
        '「鹊」在中文里意味着好消息。我们喜欢这个意象：把重要的信息，准确地送到该看见它的人眼前。六年来，我们从三个人、一间小办公室开始，把这件小事做得更认真了一点。',
      bullets: [
        { text: '2019 年，小鹊科技在上海成立' },
        { text: '2022 年，发布第一版区块化建站产品' },
        { text: '2025 年，小鹊云服务客户突破 200 家' }
      ],
      image: 'type',
      imageSide: 'left',
      badge: 'EST. 2019',
      ctaLabel: '',
      ctaHref: ''
    }),
    b('features', {
      eyebrow: '我们的价值观',
      title: '我们如何做事。',
      subtitle: '四条简单但不变的原则。',
      layout: 'grid',
      items: [
        {
          icon: 'clock',
          title: '长期主义',
          desc: '不追风口，把时间花在会被长期使用的东西上。',
          tag: ''
        },
        {
          icon: 'book',
          title: '先想清楚',
          desc: '写下来、画出来、争论过，再开始做。',
          tag: ''
        },
        {
          icon: 'brush',
          title: '对细节偏执',
          desc: '一个 1px 的错位，也值得被修正。',
          tag: ''
        },
        {
          icon: 'users',
          title: '和客户站在一起',
          desc: '客户的成功，是唯一真实的指标。',
          tag: ''
        }
      ]
    }),
    b('cta', {
      title: '想认识我们，\n或者加入我们？',
      text: '我们一直在寻找对内容和产品有耐心的人。写信给我们，聊聊你的想法。',
      primaryLabel: '联系我们',
      primaryHref: '/contact',
      secondaryLabel: '查看开放职位',
      secondaryHref: '/contact',
      theme: 'dark'
    })
  ]
})

const contact = makePage({
  id: 'pg_contact_zh',
  group: 'tg_contact',
  locale: 'zh',
  path: '/contact',
  name: '联系我们',
  version: 3,
  updatedAt: '2026-09-30T13:25:00+08:00',
  publishedAt: '2026-09-30T13:25:00+08:00',
  seo: {
    title: '联系我们 · 小鹊云',
    description: '预约产品演示，或与我们聊聊你的官网与内容需求。',
    keywords: ['联系小鹊', '预约演示', '产品咨询']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: '联系我们 / CONTACT',
      title: '说声你好，\n故事就开始了。',
      subtitle: '无论是产品咨询、合作，还是想聊聊你的官网，都欢迎联系我们。',
      primaryLabel: '',
      primaryHref: '',
      secondaryLabel: '',
      secondaryHref: '',
      note: '通常在一个工作日内回复',
      visual: 'none',
      theme: 'light'
    }),
    b('contact-form', {
      title: '聊聊你的官网。',
      subtitle: '填写下面的表单，我们会尽快联系你。',
      submitLabel: '发送消息',
      successTitle: '消息已经飞过来了',
      successText: '感谢你的联系，我们会在一个工作日内回复你。',
      email: 'hello@xiaoque.com',
      phone: '400-860-2026',
      address: '上海市徐汇区西岸智塔 18F',
      hours: '周一至周五 9:00 – 18:00',
      mapStyle: 'orbit'
    })
  ]
})

const investors = makePage({
  id: 'pg_investors_zh',
  group: 'tg_investors',
  locale: 'zh',
  path: '/investors',
  name: '投资者关系',
  status: 'draft',
  version: 0,
  updatedAt: '2026-10-17T09:10:00+08:00',
  seo: {
    title: '投资者关系 · 小鹊科技',
    description: '小鹊科技的投资者关系页面。',
    keywords: []
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: '投资者关系 / INVESTORS',
      title: '与长期主义者\n一起生长。',
      subtitle: '关于公司治理、财务信息与投资者沟通的最新资料。',
      primaryLabel: '联系投资者关系',
      primaryHref: '/contact',
      secondaryLabel: '',
      secondaryHref: '',
      note: '',
      visual: 'waves',
      theme: 'light'
    }),
    b('rich-text', {
      width: 'narrow',
      content:
        '<h2>投资者沟通</h2><p>我们重视与投资者的长期沟通。如需获取最新的公司资料，请联系 <a href="mailto:ir@xiaoque.com">ir@xiaoque.com</a>。</p><p>本页面正在完善中，将在下一次发布时上线。</p>'
    })
  ]
})

export default [home, products, solutions, cases, news, about, contact, investors]

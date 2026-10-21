import type { ContentRecord } from '~/domain/types'

export const contents: ContentRecord[] = [
  /* ------------------------------- News · ZH ------------------------------- */
  {
    id: 'ct_n1_zh',
    groupId: 'tg_news_1',
    type: 'news',
    locale: 'zh',
    slug: 'magpie-cloud-3',
    title: '小鹊云 3.0 发布：把 AI 带进内容工作流',
    summary:
      '新的版本带来 AI 辅助写作、区块智能填充与双语翻译草稿。不是替代编辑，而是把重复劳动交还给机器。',
    category: '产品更新',
    cover: 'type',
    tone: 'accent',
    status: 'published',
    publishedAt: '2026-09-28',
    readMinutes: 5,
    featured: true,
    body: `
      <p>今天，小鹊云 3.0 正式发布。这个版本的核心只有一句话：让 AI 待在它该待的位置——处理重复，而不是替你做判断。</p>
      <h2>AI 辅助写作，从空白页开始</h2>
      <p>在富文本区块里，编辑可以基于页面上下文生成初稿、续写或改写语气。生成内容始终以草稿形式出现，保留完整的撤销记录，并且不会自动发布。</p>
      <blockquote>我们不想做“一键生成官网”。我们想让你在卡住的时候，更快地往前走一步。</blockquote>
      <h2>区块智能填充</h2>
      <p>当你新增一个特性列表或 FAQ 区块时，系统会根据页面已有的内容推荐结构与文案。所有建议都可以逐条接受或忽略，不会破坏现有排版。</p>
      <h2>双语翻译草稿</h2>
      <p>中文页面可以一键生成结构一致的英文草稿，文案由翻译模型填入并标记为“待校对”。人工确认之前，英文版本不会出现在语言切换菜单里。</p>
      <p>3.0 已向所有付费客户开放，无需额外配置。你可以在「站点设置 → 实验性功能」中控制 AI 能力的开关。</p>
    `
  },
  {
    id: 'ct_n2_zh',
    groupId: 'tg_news_2',
    type: 'news',
    locale: 'zh',
    slug: 'qimu-website-rework',
    title: '从 2 天到 20 分钟：栖木设计的官网改版复盘',
    summary:
      '一个 8 人设计团队，如何把官网发布周期从两天压缩到二十分钟，同时把内容产出提升了 3.2 倍。',
    category: '客户故事',
    cover: 'stack',
    tone: 'sage',
    status: 'published',
    publishedAt: '2026-09-12',
    readMinutes: 6,
    featured: false,
    body: `
      <p>栖木设计是一家位于杭州的品牌设计工作室。2025 年底，他们决定重做已经用了四年的官网——问题不在于不好看，而在于“改不动”。</p>
      <h2>问题：每次改版都是一次项目</h2>
      <p>旧官网是一个定制开发的单页应用。换一张首图、调整一段文案，都需要设计师提需求、开发排期、测试上线。一年下来，官网只更新了三次。</p>
      <h2>改变：把页面交还给内容团队</h2>
      <p>迁移到小鹊云之后，栖木把 12 个页面拆成了区块组合。设计团队自己维护一套品牌组件，运营同学只需要在既定区块里填写内容。</p>
      <ul>
        <li>首页改版从 2 天缩短到 20 分钟</li>
        <li>半年内发布 38 次，是过去两年的 6 倍</li>
        <li>案例页平均停留时长提升 41%</li>
      </ul>
      <p>“最大的变化不是速度，而是心态。”栖木的创始人说，“我们现在愿意做小步的尝试，因为改错的成本足够低。”</p>
    `
  },
  {
    id: 'ct_n3_zh',
    groupId: 'tg_news_3',
    type: 'news',
    locale: 'zh',
    slug: 'bilingual-seo-checklist',
    title: '中英双语官网的 7 个 SEO 细节',
    summary:
      'hreflang、canonical、独立 slug、语言切换的收录策略——双语站点最容易踩的坑，一次讲清楚。',
    category: '实践指南',
    cover: 'grid',
    tone: 'cream',
    status: 'published',
    publishedAt: '2026-08-30',
    readMinutes: 8,
    featured: false,
    body: `
      <p>双语官网不只是把中文翻译成英文。如果 SEO 没有处理好，两个语言版本可能互相竞争，甚至都不被收录。</p>
      <h2>1. 每种语言都要有自引用 canonical</h2>
      <p>中文页面 canonical 指向中文 URL，英文页面指向英文 URL。不要把所有语言都指向默认语言，那会让英文页面失去索引资格。</p>
      <h2>2. hreflang 必须完整且互指</h2>
      <p>中文、英文与 x-default 需要互相声明。任何一方缺失，搜索引擎都可能忽略整组关系。</p>
      <h2>3. 英文使用独立、可读的 slug</h2>
      <p><code>/en/about-us</code> 比 <code>/en/guanyu-women</code> 更容易理解，也更利于点击。</p>
      <h2>4. 语言切换不要强制跳转</h2>
      <p>根据浏览器语言做 302 会让爬虫只看到一种版本。更好的做法是显示推荐条，让用户自己选择。</p>
      <h2>5. 未翻译页面不要出现在 sitemap</h2>
      <p>只提交已完成翻译并发布的 URL，避免大量 404 或半成品页面影响站点质量。</p>
      <h2>6. 结构化数据要按语言输出</h2>
      <p>Article、Product 与 Organization 的 name、description 都应该是当前语言。</p>
      <h2>7. 图片 alt 也要翻译</h2>
      <p>这既是无障碍要求，也是图片搜索的流量入口。小鹊云的媒体库支持为每张图分别填写中英文 alt。</p>
    `
  },
  {
    id: 'ct_n4_zh',
    groupId: 'tg_news_4',
    type: 'news',
    locale: 'zh',
    slug: 'series-a',
    title: '小鹊科技完成 A 轮融资，继续深耕内容基础设施',
    summary: '本轮融资将用于产品研发、行业解决方案与客户成功团队建设。',
    category: '公司动态',
    cover: 'waves',
    tone: 'ink',
    status: 'published',
    publishedAt: '2026-07-18',
    readMinutes: 3,
    featured: false,
    body: `
      <p>我们很高兴地宣布，小鹊科技已完成 A 轮融资。本轮资金将主要用于三件事：产品研发、行业解决方案沉淀，以及客户成功团队的建设。</p>
      <p>过去一年，小鹊云的客户数量增长了 140%。我们服务了制造、专业服务、消费品牌等行业的 260 多家企业，也收到了很多真诚的批评与建议。</p>
      <p>我们会继续把精力放在同一件小事上：让内容团队不被工具和流程拖慢。</p>
    `
  },

  {
    id: 'ct_n5_zh',
    groupId: 'tg_news_5',
    type: 'news',
    locale: 'zh',
    slug: 'calm-content-team',
    title: '内容团队的 5 个发布习惯',
    summary: '少一点等待，多一点节奏。把发布变成日常动作，而不是一次次临时项目。',
    category: '实践指南',
    cover: 'bloom',
    tone: 'cream',
    status: 'published',
    publishedAt: '2026-07-02',
    readMinutes: 4,
    featured: false,
    body: `
      <p>高效的内容团队并不是发布得更频繁，而是更少被发布本身打断。以下五个习惯，来自我们与 200 多家客户的观察。</p>
      <h2>1. 固定发布窗口</h2>
      <p>与其随时改、随时发，不如约定每周二与周四下午各一次发布窗口。编辑知道什么时候必须收尾，其他角色也知道什么时候需要就位。</p>
      <h2>2. 一次只改一件事</h2>
      <p>把「改文案」和「结构改版」分开上线。改动范围更清晰，定位问题也更快。</p>
      <h2>3. 预发布检查清单</h2>
      <p>链接、图片 alt、SEO 标题与英文版本，用同一份清单逐项确认。</p>
      <h2>4. 让数据进入例会</h2>
      <p>每次例会花五分钟看页面表现，而不只在季度复盘时打开报表。</p>
      <h2>5. 为失败留出空间</h2>
      <p>发布流程足够轻，团队才愿意做小步的尝试。节奏感来自确定性。</p>
    `
  },

  /* ------------------------------- News · EN ------------------------------- */
  {
    id: 'ct_n1_en',
    groupId: 'tg_news_1',
    type: 'news',
    locale: 'en',
    slug: 'magpie-cloud-3-0',
    title: 'Magpie Cloud 3.0: bringing AI into the content workflow',
    summary:
      'AI-assisted writing, smart block suggestions and bilingual draft translations. Not a replacement for editors — a way to give the busywork back to the machine.',
    category: 'Product',
    cover: 'type',
    tone: 'accent',
    status: 'published',
    publishedAt: '2026-09-28',
    readMinutes: 5,
    featured: true,
    body: `
      <p>Magpie Cloud 3.0 is here. The idea behind it is simple: let AI do what AI should do — handle repetition, not make judgment calls for you.</p>
      <h2>Writing assistance, starting from a blank page</h2>
      <p>Inside any rich-text block, editors can generate a first draft, continue a paragraph, or adjust tone based on page context. Suggestions always arrive as drafts, with full undo history, and are never published automatically.</p>
      <blockquote>We are not building “one-click websites.” We want to help you take one step forward when you are stuck.</blockquote>
      <h2>Smarter block filling</h2>
      <p>When you add a feature list or FAQ block, the system recommends structure and copy based on the rest of the page. Accept or ignore suggestions one by one — your layout stays intact.</p>
      <h2>Bilingual translation drafts</h2>
      <p>A Chinese page can generate a matching English draft with translated copy, clearly marked “needs review.” Until a human confirms it, the English version stays out of the language switcher.</p>
      <p>3.0 is available to all paid customers today. AI features can be toggled in Site Settings → Experimental.</p>
    `
  },
  {
    id: 'ct_n2_en',
    groupId: 'tg_news_2',
    type: 'news',
    locale: 'en',
    slug: 'qimu-case-study',
    title: 'From two days to twenty minutes: how Qimu rebuilt its site',
    summary:
      'An eight-person design studio cut its publishing cycle from two days to twenty minutes — and tripled its content output.',
    category: 'Customer story',
    cover: 'stack',
    tone: 'sage',
    status: 'published',
    publishedAt: '2026-09-12',
    readMinutes: 6,
    featured: false,
    body: `
      <p>Qimu Design is a brand studio based in Hangzhou. At the end of 2025, they decided to rebuild a website they had used for four years. The problem was not that it looked bad — it was that it could not be changed.</p>
      <h2>The problem: every update was a project</h2>
      <p>The old site was a custom single-page app. Changing a hero image or a paragraph meant a designer filing a ticket, a developer scheduling work, and QA before launch. In a full year, the site was updated three times.</p>
      <h2>The change: give the page back to the content team</h2>
      <p>After moving to Magpie Cloud, Qimu split twelve pages into block compositions. Design owns the brand components; the marketing team fills content inside approved blocks.</p>
      <ul>
        <li>Homepage redesigns went from 2 days to 20 minutes</li>
        <li>38 releases in six months — six times the previous two years</li>
        <li>Average time on case study pages up 41%</li>
      </ul>
      <p>“The biggest change is not speed, it is mindset,” says Qimu’s founder. “We are willing to try small things now, because being wrong is cheap.”</p>
    `
  },
  {
    id: 'ct_n3_en',
    groupId: 'tg_news_3',
    type: 'news',
    locale: 'en',
    slug: 'bilingual-seo-checklist',
    title: 'Seven SEO details for bilingual websites',
    summary:
      'hreflang, canonical, separate slugs and language-switcher indexing — the traps bilingual sites fall into, explained.',
    category: 'Guide',
    cover: 'grid',
    tone: 'cream',
    status: 'published',
    publishedAt: '2026-08-30',
    readMinutes: 8,
    featured: false,
    body: `
      <p>A bilingual website is not just a translation. Without the right SEO, two language versions can compete with each other — and neither gets indexed well.</p>
      <h2>1. Self-referencing canonical per language</h2>
      <p>The Chinese page points to the Chinese URL, English to English. Pointing everything at the default language can disqualify the English pages from indexing.</p>
      <h2>2. Complete, reciprocal hreflang</h2>
      <p>Chinese, English and x-default must declare each other. If one is missing, search engines may ignore the entire group.</p>
      <h2>3. Readable slugs for English</h2>
      <p><code>/en/about-us</code> is easier to understand and click than a transliterated slug.</p>
      <h2>4. Do not force redirects by browser language</h2>
      <p>A 302 based on Accept-Language means crawlers only ever see one version. Show a recommendation bar instead.</p>
      <h2>5. Keep untranslated pages out of the sitemap</h2>
      <p>Only submit translated and published URLs. Half-finished pages create 404s and drag down site quality.</p>
      <h2>6. Localize structured data</h2>
      <p>The name and description in Article, Product and Organization schema should match the current language.</p>
      <h2>7. Translate image alt text</h2>
      <p>This is both an accessibility requirement and a traffic source. Magpie’s media library stores alt text per language.</p>
    `
  },
  {
    id: 'ct_n4_en',
    groupId: 'tg_news_4',
    type: 'news',
    locale: 'en',
    slug: 'series-a',
    title: 'Magpie Tech raises Series A to keep building content infrastructure',
    summary: 'The round will fund product development, industry solutions and our customer success team.',
    category: 'Company',
    cover: 'waves',
    tone: 'ink',
    status: 'published',
    publishedAt: '2026-07-18',
    readMinutes: 3,
    featured: false,
    body: `
      <p>We are happy to share that Magpie Tech has closed its Series A. The funding goes to three things: product development, reusable industry solutions, and our customer success team.</p>
      <p>Over the past year our customer base grew 140%. We now serve more than 260 companies across manufacturing, professional services and consumer brands — and we have received plenty of honest criticism along the way.</p>
      <p>We will keep focusing on the same small thing: making sure content teams are never slowed down by their tools.</p>
    `
  },

  {
    id: 'ct_n5_en',
    groupId: 'tg_news_5',
    type: 'news',
    locale: 'en',
    slug: 'calm-content-teams',
    title: 'Five publishing habits of calm content teams',
    summary: 'Less waiting, more rhythm. Make publishing a daily habit instead of a series of last-minute projects.',
    category: 'Guide',
    cover: 'bloom',
    tone: 'cream',
    status: 'published',
    publishedAt: '2026-07-02',
    readMinutes: 4,
    featured: false,
    body: `
      <p>Effective content teams do not publish more often. They are simply interrupted less by the act of publishing. Here are five habits we have seen across more than 200 customers.</p>
      <h2>1. Keep a fixed release window</h2>
      <p>Rather than shipping whenever something changes, agree on Tuesday and Thursday afternoons. Editors know when to finish; everyone else knows when to be available.</p>
      <h2>2. Change one thing at a time</h2>
      <p>Separate copy edits from structural redesigns. Problems are easier to isolate, and easier to roll back.</p>
      <h2>3. Use a pre-flight checklist</h2>
      <p>Links, image alt text, SEO titles and the English version — confirmed against the same list every time.</p>
      <h2>4. Bring data into the weekly meeting</h2>
      <p>Spend five minutes on page performance each week, not only in quarterly reports.</p>
      <h2>5. Leave room to fail</h2>
      <p>When rollback is cheap, teams are willing to try small things. Rhythm comes from safety.</p>
    `
  },

  /* ------------------------------ Products · ZH ----------------------------- */
  {
    id: 'ct_p1_zh',
    groupId: 'tg_prod_1',
    type: 'product',
    locale: 'zh',
    slug: 'magpie-cms',
    title: '小鹊云 CMS',
    label: '旗舰产品',
    summary: '从页面搭建到多语言内容的完整内容云。',
    description: '面向需要持续运营官网的团队，把页面、内容、媒体与发布流程放在同一个系统里。',
    bullets: ['区块化页面搭建', '中英双语内容管理', '草稿预览与直接发布', '表单与线索管理'],
    cover: 'orbit',
    tone: 'accent',
    status: 'published',
    publishedAt: '2026-03-01'
  },
  {
    id: 'ct_p2_zh',
    groupId: 'tg_prod_2',
    type: 'product',
    locale: 'zh',
    slug: 'magpie-sites',
    title: '小鹊建站',
    label: '快速上线',
    summary: '为成长型企业交付一个有品牌感的官网。',
    description: '从信息架构到视觉落地，用标准化流程在两周内交付一个可持续运营的官网。',
    bullets: ['品牌官网设计与搭建', '最快 10 天上线', '响应式与 SEO 就绪', '托管、证书与运维'],
    cover: 'arch',
    tone: 'ink',
    status: 'published',
    publishedAt: '2026-04-15'
  },
  {
    id: 'ct_p3_zh',
    groupId: 'tg_prod_3',
    type: 'product',
    locale: 'zh',
    slug: 'magpie-insights',
    title: '小鹊洞察',
    label: '数据洞察',
    summary: '让内容决策有据可依。',
    description: '把页面访问、内容表现与转化路径放在一张图里，帮助内容团队判断下一步该做什么。',
    bullets: ['页面与内容表现分析', '搜索表现追踪', '转化路径与表单漏斗', '双周洞察报告'],
    cover: 'bloom',
    tone: 'cream',
    status: 'published',
    publishedAt: '2026-06-10'
  },

  /* ------------------------------ Products · EN ----------------------------- */
  {
    id: 'ct_p1_en',
    groupId: 'tg_prod_1',
    type: 'product',
    locale: 'en',
    slug: 'magpie-cms',
    title: 'Magpie CMS',
    label: 'FLAGSHIP',
    summary: 'A complete content cloud, from page building to multilingual content.',
    description: 'For teams that run their website every day — pages, content, media and release workflow in one system.',
    bullets: ['Block-based page builder', 'Bilingual content management', 'Version snapshots & rollback', 'Forms and lead management'],
    cover: 'orbit',
    tone: 'accent',
    status: 'published',
    publishedAt: '2026-03-01'
  },
  {
    id: 'ct_p2_en',
    groupId: 'tg_prod_2',
    type: 'product',
    locale: 'en',
    slug: 'magpie-sites',
    title: 'Magpie Sites',
    label: 'LAUNCH FAST',
    summary: 'A brand-grade website, delivered for growing companies.',
    description: 'From information architecture to visual design, a website you can keep operating — delivered in two weeks.',
    bullets: ['Brand website design & build', 'Live in as little as 10 days', 'Responsive and SEO-ready', 'Hosting, SSL and maintenance'],
    cover: 'arch',
    tone: 'ink',
    status: 'published',
    publishedAt: '2026-04-15'
  },
  {
    id: 'ct_p3_en',
    groupId: 'tg_prod_3',
    type: 'product',
    locale: 'en',
    slug: 'magpie-insights',
    title: 'Magpie Insights',
    label: 'ANALYTICS',
    summary: 'Make content decisions with evidence.',
    description: 'Page traffic, content performance and conversion paths in one view — so your team knows what to do next.',
    bullets: ['Page and content analytics', 'Search performance tracking', 'Conversion and form funnels', 'Biweekly insight reports'],
    cover: 'bloom',
    tone: 'cream',
    status: 'published',
    publishedAt: '2026-06-10'
  },

  /* ------------------------------- Cases · ZH ------------------------------- */
  {
    id: 'ct_c1_zh',
    groupId: 'tg_case_1',
    type: 'case',
    locale: 'zh',
    slug: 'hilltop',
    title: '山丘制造',
    client: '山丘制造',
    industry: '智能制造',
    summary: '用一套双语产品目录，把海外询盘从每月 43 条提升到 123 条。',
    metric: '186',
    metricSuffix: '%',
    metricLabel: '海外询盘增长',
    cover: 'arch',
    tone: 'ink',
    status: 'published',
    publishedAt: '2026-09-20'
  },
  {
    id: 'ct_c2_zh',
    groupId: 'tg_case_2',
    type: 'case',
    locale: 'zh',
    slug: 'nova-capital',
    title: 'NOVA 资本',
    client: 'NOVA Capital',
    industry: '专业服务',
    summary: '中英双语官网上线 6 个月，品牌搜索量增长 72%。',
    metric: '72',
    metricSuffix: '%',
    metricLabel: '品牌搜索增长',
    cover: 'orbit',
    tone: 'accent',
    status: 'published',
    publishedAt: '2026-08-16'
  },
  {
    id: 'ct_c3_zh',
    groupId: 'tg_case_3',
    type: 'case',
    locale: 'zh',
    slug: 'qimu',
    title: '栖木设计',
    client: '栖木设计',
    industry: '消费品牌',
    summary: '把发布周期从 2 天压缩到 20 分钟，内容产出提升 3.2 倍。',
    metric: '3.2',
    metricSuffix: '×',
    metricLabel: '内容效率',
    cover: 'grid',
    tone: 'cream',
    status: 'published',
    publishedAt: '2026-08-02'
  },
  {
    id: 'ct_c4_zh',
    groupId: 'tg_case_4',
    type: 'case',
    locale: 'zh',
    slug: 'lumen',
    title: 'LUMEN',
    client: 'LUMEN',
    industry: '科技',
    summary: '用区块化官网上线 3 个产品站，开发投入减少 60%。',
    metric: '60',
    metricSuffix: '%',
    metricLabel: '开发投入降低',
    cover: 'stack',
    tone: 'sage',
    status: 'published',
    publishedAt: '2026-06-25'
  },
  {
    id: 'ct_c5_zh',
    groupId: 'tg_case_5',
    type: 'case',
    locale: 'zh',
    slug: 'cloudrise',
    title: '云启资本',
    client: '云启资本',
    industry: '专业服务',
    summary: '让合伙人观点成为品牌资产，文章阅读量提升 140%。',
    metric: '140',
    metricSuffix: '%',
    metricLabel: '内容阅读增长',
    cover: 'type',
    tone: 'accent',
    status: 'published',
    publishedAt: '2026-05-18'
  },
  {
    id: 'ct_c6_zh',
    groupId: 'tg_case_6',
    type: 'case',
    locale: 'zh',
    slug: 'weichi',
    title: '未迟',
    client: '未迟',
    industry: '消费品牌',
    summary: '新品发布从 3 周缩短到 3 天，首月转化率提升 28%。',
    metric: '28',
    metricSuffix: '%',
    metricLabel: '转化率提升',
    cover: 'dunes',
    tone: 'sand',
    status: 'published',
    publishedAt: '2026-04-09'
  },

  /* ------------------------------- Cases · EN ------------------------------- */
  {
    id: 'ct_c1_en',
    groupId: 'tg_case_1',
    type: 'case',
    locale: 'en',
    slug: 'hilltop',
    title: 'Hilltop Manufacturing',
    client: 'Hilltop',
    industry: 'Manufacturing',
    summary: 'A bilingual product catalogue lifted overseas inquiries from 43 to 123 per month.',
    metric: '186',
    metricSuffix: '%',
    metricLabel: 'More inquiries',
    cover: 'arch',
    tone: 'ink',
    status: 'published',
    publishedAt: '2026-09-20'
  },
  {
    id: 'ct_c2_en',
    groupId: 'tg_case_2',
    type: 'case',
    locale: 'en',
    slug: 'nova-capital',
    title: 'NOVA Capital',
    client: 'NOVA Capital',
    industry: 'Professional services',
    summary: 'Six months after launching its bilingual site, brand search grew 72%.',
    metric: '72',
    metricSuffix: '%',
    metricLabel: 'Brand search growth',
    cover: 'orbit',
    tone: 'accent',
    status: 'published',
    publishedAt: '2026-08-16'
  },
  {
    id: 'ct_c3_en',
    groupId: 'tg_case_3',
    type: 'case',
    locale: 'en',
    slug: 'qimu',
    title: 'Qimu Design',
    client: 'Qimu',
    industry: 'Consumer brands',
    summary: 'Publishing went from two days to twenty minutes; content output grew 3.2×.',
    metric: '3.2',
    metricSuffix: '×',
    metricLabel: 'Content velocity',
    cover: 'grid',
    tone: 'cream',
    status: 'published',
    publishedAt: '2026-08-02'
  },
  {
    id: 'ct_c4_en',
    groupId: 'tg_case_4',
    type: 'case',
    locale: 'en',
    slug: 'lumen',
    title: 'LUMEN',
    client: 'LUMEN',
    industry: 'Technology',
    summary: 'Three product sites launched on blocks, with 60% less engineering effort.',
    metric: '60',
    metricSuffix: '%',
    metricLabel: 'Less dev effort',
    cover: 'stack',
    tone: 'sage',
    status: 'published',
    publishedAt: '2026-06-25'
  },
  {
    id: 'ct_c5_en',
    groupId: 'tg_case_5',
    type: 'case',
    locale: 'en',
    slug: 'cloudrise',
    title: 'Cloudrise Capital',
    client: 'Cloudrise',
    industry: 'Professional services',
    summary: 'Partner perspectives became a brand asset, lifting article reads by 140%.',
    metric: '140',
    metricSuffix: '%',
    metricLabel: 'More reads',
    cover: 'type',
    tone: 'accent',
    status: 'published',
    publishedAt: '2026-05-18'
  },
  {
    id: 'ct_c6_en',
    groupId: 'tg_case_6',
    type: 'case',
    locale: 'en',
    slug: 'weichi',
    title: 'Weichi',
    client: 'Weichi',
    industry: 'Consumer brands',
    summary: 'Product launches went from three weeks to three days; first-month conversion rose 28%.',
    metric: '28',
    metricSuffix: '%',
    metricLabel: 'More conversion',
    cover: 'dunes',
    tone: 'sand',
    status: 'published',
    publishedAt: '2026-04-09'
  }
]

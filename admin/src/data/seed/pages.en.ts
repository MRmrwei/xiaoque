import { b, makePage } from './page-factory'

const home = makePage({
  id: 'pg_home_en',
  group: 'tg_home',
  locale: 'en',
  path: '/',
  name: 'Home',
  version: 8,
  updatedAt: '2026-10-14T16:20:00+08:00',
  publishedAt: '2026-10-14T16:20:00+08:00',
  seo: {
    title: 'Magpie Cloud — Turn your website into a growth engine',
    description:
      'Magpie Cloud is a block-based website builder and content cloud for growing companies. Built by your team, bilingual by default, live in seconds.',
    keywords: ['company website', 'content management', 'bilingual website', 'block based builder']
  },
  blocks: [
    b('hero-carousel', {
      autoplay: true,
      interval: 7000,
      note: 'Swipe or use ← → · Images and video supported',
      slides: [
        {
          id: 'en_carousel_1',
          mediaType: 'image',
          media: 'builder',
          videoUrl: '',
          eyebrow: 'MAGPIE CLOUD 3.0 · CONTENT CLOUD',
          title: 'Every story\ndeserves a\n<em>beautiful</em> stage.',
          subtitle:
            'Magpie Cloud turns your website into a living content system — built by your team, bilingual by default, live in seconds.',
          primaryLabel: 'Book a demo',
          primaryHref: '/contact',
          secondaryLabel: 'Explore the platform',
          secondaryHref: '/products'
        },
        {
          id: 'en_carousel_2',
          mediaType: 'image',
          media: 'arch',
          videoUrl: '',
          eyebrow: 'ONE CONTENT SYSTEM',
          title: 'From first draft\nto live page,\nin one <em>place</em>.',
          subtitle: 'Pages, content, media and releases — one system for the whole site.',
          primaryLabel: 'Explore the platform',
          primaryHref: '/products',
          secondaryLabel: 'Book a demo',
          secondaryHref: '/contact'
        },
        {
          id: 'en_carousel_3',
          mediaType: 'video',
          media: 'orbit',
          videoUrl:
            '/media/hero-demo.mp4',
          eyebrow: 'SAFE PUBLISH',
          title: 'Every release.\nA <em>safe</em> way\nback.',
          subtitle: 'Immutable version snapshots. Any change can be traced and rolled back.',
          primaryLabel: 'How publishing works',
          primaryHref: '/products',
          secondaryLabel: 'Customer stories',
          secondaryHref: '/cases'
        }
      ]
    }),
    b('partners', {
      title: 'Teams growing with Magpie',
      note: '260+ companies build with Magpie Cloud',
      marquee: true,
      items: [
        { name: 'NOVA' },
        { name: 'Hilltop' },
        { name: 'LUMEN' },
        { name: 'Qimu' },
        { name: 'HORIZON' },
        { name: 'Cloudrise' },
        { name: 'Weichi' },
        { name: 'ARK' }
      ]
    }),
    b('features', {
      eyebrow: 'WHY MAGPIE',
      title: 'We handle the complex.\nYour team keeps it simple.',
      subtitle: 'Every detail is designed so marketing can ship without waiting for engineering.',
      layout: 'bento',
      items: [
        {
          icon: 'blocks',
          title: 'Block-based builder',
          desc: '14 ready-made blocks, composed like building blocks. Responsive and on-brand by default.',
          tag: 'CORE'
        },
        {
          icon: 'globe',
          title: 'Bilingual by design',
          desc: 'One console, two languages. Separate drafts and releases that never overwrite each other.',
          tag: ''
        },
        {
          icon: 'search',
          title: 'SEO, built in',
          desc: 'SSR, hreflang, bilingual sitemaps and structured data — ready out of the box.',
          tag: ''
        },
        {
          icon: 'history',
          title: 'Safe publishing',
          desc: 'Every release is an immutable snapshot with one-click rollback.',
          tag: ''
        },
        {
          icon: 'image',
          title: 'Media library',
          desc: 'Upload, compress, bilingual alt text and reference checks in one place.',
          tag: ''
        },
        {
          icon: 'shield',
          title: 'Roles & audit',
          desc: 'Button-level permissions, preview tokens, every action on record.',
          tag: ''
        }
      ]
    }),
    b('image-text', {
      eyebrow: 'FOR MARKETING TEAMS',
      title: 'Ten revisions?\nNo more waiting on dev.',
      body:
        'A three-panel editor lets your team build pages directly. Drag, edit, preview, autosave — copy and imagery changes no longer sit in a sprint backlog.',
      bullets: [
        { text: 'Drag-and-drop blocks, organized like a document' },
        { text: 'Desktop, tablet and mobile previews — exactly what you see' },
        { text: 'Drafts isolated from production until you hit publish' }
      ],
      image: 'stack',
      imageSide: 'right',
      badge: 'LIVE PREVIEW',
      ctaLabel: 'See how it works',
      ctaHref: '/products'
    }),
    b('stats', {
      theme: 'dark',
      note: 'Source: 2025–2026 customer survey',
      items: [
        { value: '60', suffix: 's', label: 'To go live' },
        { value: '98', suffix: '%', label: 'Net revenue retention' },
        { value: '3', suffix: '×', label: 'Content velocity' },
        { value: '14', suffix: '', label: 'Ready-made blocks' }
      ]
    }),
    b('products', {
      eyebrow: 'PRODUCT SUITE',
      title: 'More than a builder.\nA complete content stack.',
      subtitle: 'Start with what you need. Add the rest as you grow.',
      limit: 3,
      moreLabel: 'View all products',
      moreHref: '/products'
    }),
    b('cases', {
      eyebrow: 'CUSTOMERS',
      title: 'Real growth,\nin real numbers.',
      subtitle: 'Manufacturing, professional services and consumer brands.',
      limit: 3,
      showFilters: false
    }),
    b('news-list', {
      eyebrow: 'JOURNAL',
      title: 'What we have been up to.',
      subtitle: '',
      limit: 3,
      showFilters: false
    }),
    b('faq', {
      eyebrow: 'FAQ',
      title: 'Good to know.',
      subtitle: '',
      items: [
        {
          q: 'Our team does not write code. Can we use it?',
          a: 'Absolutely. Every block is configured through forms. Drag, fill in, preview, publish — that is the entire workflow.'
        },
        {
          q: 'Can we migrate from our current site?',
          a: 'Yes. We help map your information architecture and provide bulk import tools. A full migration usually takes 3–7 days.'
        },
        {
          q: 'Do Chinese and English need to launch together?',
          a: 'No. Each language has its own draft and release. Launch Chinese first and complete English later — untranslated pages simply hide the language switcher.'
        },
        {
          q: 'Can we use our own domain?',
          a: 'Yes. Bring your own domain with full DNS guidance. We run on regional nodes and can deploy closer to your audience.'
        },
        {
          q: 'How is our data protected?',
          a: 'Daily backups, versioned media storage, button-level permissions, and a complete audit trail for every publish and rollback.'
        }
      ]
    }),
    b('cta', {
      title: 'Ready to put your website\nto work?',
      text: 'Book a 30-minute demo and see how Magpie Cloud fits your team.',
      primaryLabel: 'Book a demo',
      primaryHref: '/contact',
      secondaryLabel: 'View product',
      secondaryHref: '/products',
      theme: 'accent'
    })
  ]
})

const products = makePage({
  id: 'pg_products_en',
  group: 'tg_products',
  locale: 'en',
  path: '/products',
  name: 'Product',
  version: 6,
  updatedAt: '2026-10-09T11:40:00+08:00',
  publishedAt: '2026-10-09T11:40:00+08:00',
  seo: {
    title: 'Platform · Magpie Cloud',
    description:
      'A block-based page builder, bilingual content management, versioned publishing and a media library — one console for your entire website.',
    keywords: ['page builder', 'content management', 'bilingual website', 'Magpie Cloud']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: 'PRODUCT / PLATFORM',
      title: 'One console for\nyour whole website.',
      subtitle:
        'Magpie Cloud is three products working together: build, content and insight. Use them as a suite, or start with the piece you need most.',
      primaryLabel: 'Book a demo',
      primaryHref: '/contact',
      secondaryLabel: 'See customer stories',
      secondaryHref: '/cases',
      note: '',
      visual: 'orbit',
      theme: 'light'
    }),
    b('features', {
      eyebrow: 'CAPABILITIES',
      title: 'From the first sentence\nto every release.',
      subtitle: 'Leave the repetitive work to the system. Keep the judgment for people.',
      layout: 'grid',
      items: [
        {
          icon: 'blocks',
          title: 'Block-based pages',
          desc: 'Compose 14 ready-made blocks. Configure through forms, and let components handle responsiveness.',
          tag: ''
        },
        {
          icon: 'globe',
          title: 'Bilingual content',
          desc: 'Pages, content, navigation and settings managed per language, with translation status at a glance.',
          tag: ''
        },
        {
          icon: 'history',
          title: 'Versions & releases',
          desc: 'Immutable snapshots, isolated drafts and one-click rollback. Publishing stops being a gamble.',
          tag: ''
        },
        {
          icon: 'image',
          title: 'Media library',
          desc: 'Drag to upload, automatic compression, bilingual alt text and reference checks.',
          tag: ''
        },
        {
          icon: 'form',
          title: 'Forms & leads',
          desc: 'Built-in contact forms with stored submissions, email notifications and CSV export.',
          tag: ''
        },
        {
          icon: 'shield',
          title: 'Roles & audit',
          desc: 'Separate editor and publisher roles, button-level permissions and a full audit trail.',
          tag: ''
        }
      ]
    }),
    b('products', {
      eyebrow: 'THREE PRODUCTS',
      title: 'Choose your starting point.',
      subtitle: 'Each product works on its own, and better together.',
      limit: 3,
      moreLabel: '',
      moreHref: ''
    }),
    b('image-text', {
      eyebrow: 'PUBLISHING',
      title: 'Publishing should not\nfeel like a gamble.',
      body:
        'Drafts are fully isolated from production. Every release creates an immutable snapshot, so any change can be traced, compared and rolled back.',
      bullets: [
        { text: 'Publish one language, or all languages at once' },
        { text: 'Keep the last 100 version snapshots' },
        { text: 'Automatic checks for required fields, broken images and duplicate paths' }
      ],
      image: 'arch',
      imageSide: 'right',
      badge: 'SAFE PUBLISH',
      ctaLabel: 'How publishing works',
      ctaHref: '/contact'
    }),
    b('faq', {
      eyebrow: 'PRODUCT FAQ',
      title: 'Before you go live,\nyou might ask.',
      subtitle: '',
      items: [
        {
          q: 'Can we buy only part of the platform?',
          a: 'Yes. Magpie Sites is enough for a brand website. Upgrade to Magpie CMS when you need a content hub, multiple languages, or lead management.'
        },
        {
          q: 'Can we import our existing content?',
          a: 'We support bulk import from Excel, Word and common CMS platforms, with field mapping. Your current site keeps running during the migration.'
        },
        {
          q: 'How fast do changes go live?',
          a: 'Within 60 seconds across the site and CDN. The server caches by version and ETag, and purges proactively on publish.'
        },
        {
          q: 'Can it match our design system?',
          a: 'Yes. Blocks expose brand color, type and radius tokens. For deeper needs, we can extend custom components on the same block protocol.'
        }
      ]
    }),
    b('cta', {
      title: 'Curious how it fits\nyour business?',
      text: 'We will prepare a demo based on your actual information architecture.',
      primaryLabel: 'Book a demo',
      primaryHref: '/contact',
      secondaryLabel: 'See customer stories',
      secondaryHref: '/cases',
      theme: 'dark'
    })
  ]
})

const solutions = makePage({
  id: 'pg_solutions_en',
  group: 'tg_solutions',
  locale: 'en',
  path: '/solutions',
  name: 'Solutions',
  version: 4,
  updatedAt: '2026-10-08T14:10:00+08:00',
  publishedAt: '2026-10-08T14:10:00+08:00',
  seo: {
    title: 'Industry solutions · Magpie Cloud',
    description:
      'Website and content solutions for manufacturing, professional services, consumer brands and content teams.',
    keywords: ['industry website', 'manufacturing website', 'brand website']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: 'SOLUTIONS',
      title: 'A website shaped\nfor your industry.',
      subtitle:
        'We work with manufacturers, professional services and consumer brands — and bring reusable structures, content models and translation workflows.',
      primaryLabel: 'Book a demo',
      primaryHref: '/contact',
      secondaryLabel: 'See customer stories',
      secondaryHref: '/cases',
      note: '',
      visual: 'grid',
      theme: 'light'
    }),
    b('features', {
      eyebrow: 'BY INDUSTRY',
      title: 'Know the industry.\nThen talk about fit.',
      subtitle: 'Not templates — page structures that start from real business needs.',
      layout: 'grid',
      items: [
        {
          icon: 'layers',
          title: 'Manufacturing',
          desc: 'Product specs, solutions, dealer networks and overseas inquiries, with multilingual catalogues.',
          tag: 'MFG'
        },
        {
          icon: 'users',
          title: 'Professional services',
          desc: 'Partner profiles, thought leadership and trust signals that make expertise visible.',
          tag: 'SERVICES'
        },
        {
          icon: 'spark',
          title: 'Consumer brands',
          desc: 'Product drops, brand stories and member touchpoints on a seasonal publishing rhythm.',
          tag: 'BRANDS'
        },
        {
          icon: 'book',
          title: 'Content teams',
          desc: 'Multilingual workflows, media assets and publishing calendars that keep everyone aligned.',
          tag: 'CONTENT'
        }
      ]
    }),
    b('image-text', {
      eyebrow: 'REUSABLE CONTENT MODEL',
      title: 'Structure can be copied.\nVoice should not.',
      body:
        'Pages, content and categories are stored per language. Creating an English version copies structure and non-translatable fields only, leaving copy for translators — no overwrites.',
      bullets: [
        { text: 'Chinese and English records are paired as translation groups' },
        { text: 'Data-source blocks filter automatically by current language' },
        { text: 'Language switching only points to translated, published pages' }
      ],
      image: 'waves',
      imageSide: 'left',
      badge: 'ZH / EN',
      ctaLabel: 'Explore multilingual support',
      ctaHref: '/products'
    }),
    b('cases', {
      eyebrow: 'OUTCOMES',
      title: 'Changes that already happened.',
      subtitle: '',
      limit: 3,
      showFilters: true
    }),
    b('stats', {
      theme: 'paper',
      note: '',
      items: [
        { value: '12', suffix: '', label: 'Industries' },
        { value: '260', suffix: '+', label: 'Customers' },
        { value: '8', suffix: '', label: 'Cities' },
        { value: '99.98', suffix: '%', label: 'Uptime' }
      ]
    }),
    b('cta', {
      title: 'Your industry deserves\na better website.',
      text: 'Tell us about your business and goals. We will share a first perspective within one business day.',
      primaryLabel: 'Contact us',
      primaryHref: '/contact',
      secondaryLabel: 'View the platform',
      secondaryHref: '/products',
      theme: 'accent'
    })
  ]
})

const cases = makePage({
  id: 'pg_cases_en',
  group: 'tg_cases',
  locale: 'en',
  path: '/cases',
  name: 'Customers',
  version: 7,
  updatedAt: '2026-10-11T09:30:00+08:00',
  publishedAt: '2026-10-11T09:30:00+08:00',
  seo: {
    title: 'Customer stories · Magpie Cloud',
    description: 'Website and content platform stories from manufacturing, professional services and consumer brands.',
    keywords: ['customer stories', 'website redesign', 'case studies']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: 'CUSTOMERS',
      title: 'Growth is not\nan adjective.\nIt is numbers.',
      subtitle: 'See how teams across industries turn their website into part of the business.',
      primaryLabel: 'Book a demo',
      primaryHref: '/contact',
      secondaryLabel: 'View the platform',
      secondaryHref: '/products',
      note: '',
      visual: 'bloom',
      theme: 'light'
    }),
    b('cases', {
      eyebrow: 'BROWSE BY INDUSTRY',
      title: 'Every industry finds\nits own answer.',
      subtitle: '',
      limit: 6,
      showFilters: true
    }),
    b('stats', {
      theme: 'dark',
      note: 'Average results, six months after launch',
      items: [
        { value: '186', suffix: '%', label: 'More overseas inquiries' },
        { value: '72', suffix: '%', label: 'Brand search growth' },
        { value: '3.2', suffix: '×', label: 'Content output' },
        { value: '20', suffix: 'min', label: 'Average release time' }
      ]
    }),
    b('cta', {
      title: 'The next story\ncould be yours.',
      text: 'Tell us where you want to go. We will figure out what we can build together.',
      primaryLabel: 'Book a demo',
      primaryHref: '/contact',
      secondaryLabel: 'Back to product',
      secondaryHref: '/products',
      theme: 'accent'
    })
  ]
})

const news = makePage({
  id: 'pg_news_en',
  group: 'tg_news',
  locale: 'en',
  path: '/news',
  name: 'Journal',
  version: 3,
  updatedAt: '2026-10-13T10:15:00+08:00',
  publishedAt: '2026-10-13T10:15:00+08:00',
  seo: {
    title: 'Journal · Magpie Cloud',
    description: 'Product updates, customer stories and notes on content, design and growth.',
    keywords: ['Magpie news', 'product updates', 'content operations']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: 'JOURNAL',
      title: 'What we’re building,\nand what we’re thinking.',
      subtitle: 'Product updates, customer stories, and notes on content, design and growth.',
      primaryLabel: '',
      primaryHref: '',
      secondaryLabel: '',
      secondaryHref: '',
      note: '',
      visual: 'none',
      theme: 'light'
    }),
    b('news-list', {
      eyebrow: 'ALL STORIES',
      title: '',
      subtitle: '',
      limit: 9,
      showFilters: true
    })
  ]
})

const about = makePage({
  id: 'pg_about_en',
  group: 'tg_about',
  locale: 'en',
  path: '/about',
  name: 'About',
  version: 5,
  updatedAt: '2026-10-06T17:00:00+08:00',
  publishedAt: '2026-10-06T17:00:00+08:00',
  seo: {
    title: 'About Magpie · Magpie Tech',
    description:
      'Founded in 2019, Magpie Tech builds lighter, freer content infrastructure for growing companies.',
    keywords: ['about Magpie', 'Magpie Tech', 'company']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: 'ABOUT / MAGPIE',
      title: 'We believe good\nstories deserve\nto be seen.',
      subtitle:
        'Founded in 2019, Magpie Tech builds lighter, freer content infrastructure for growing companies — so good work never waits on process.',
      primaryLabel: 'Join us',
      primaryHref: '/contact',
      secondaryLabel: 'Meet our customers',
      secondaryHref: '/cases',
      note: '',
      visual: 'arch',
      theme: 'light'
    }),
    b('stats', {
      theme: 'paper',
      note: '',
      items: [
        { value: '2019', suffix: '', label: 'Founded' },
        { value: '42', suffix: '', label: 'Team members' },
        { value: '12', suffix: '', label: 'Industries' },
        { value: '260', suffix: '+', label: 'Customers' }
      ]
    }),
    b('image-text', {
      eyebrow: 'OUR STORY',
      title: 'It started with a magpie.',
      body:
        'In Chinese culture, the magpie carries good news. We like that image: delivering something important, precisely, to the people who need to see it. Six years on, we are still working on that small idea — a little more carefully each year.',
      bullets: [
        { text: '2019 — Magpie Tech is founded in Shanghai' },
        { text: '2022 — The first block-based builder ships' },
        { text: '2025 — More than 200 companies run on Magpie Cloud' }
      ],
      image: 'type',
      imageSide: 'left',
      badge: 'EST. 2019',
      ctaLabel: '',
      ctaHref: ''
    }),
    b('features', {
      eyebrow: 'OUR VALUES',
      title: 'How we work.',
      subtitle: 'Four simple principles we keep coming back to.',
      layout: 'grid',
      items: [
        {
          icon: 'clock',
          title: 'Play the long game',
          desc: 'Ignore the hype. Spend time on things people will still use in ten years.',
          tag: ''
        },
        {
          icon: 'book',
          title: 'Think it through',
          desc: 'Write it, sketch it, argue about it — then start building.',
          tag: ''
        },
        {
          icon: 'brush',
          title: 'Obsess over details',
          desc: 'A single misplaced pixel is worth fixing.',
          tag: ''
        },
        {
          icon: 'users',
          title: 'Stand with customers',
          desc: 'Customer success is the only metric that truly matters.',
          tag: ''
        }
      ]
    }),
    b('cta', {
      title: 'Want to know us,\nor join us?',
      text: 'We are always looking for people with patience for content and product. Write to us.',
      primaryLabel: 'Contact us',
      primaryHref: '/contact',
      secondaryLabel: 'See open roles',
      secondaryHref: '/contact',
      theme: 'dark'
    })
  ]
})

const contact = makePage({
  id: 'pg_contact_en',
  group: 'tg_contact',
  locale: 'en',
  path: '/contact',
  name: 'Contact',
  version: 2,
  updatedAt: '2026-09-30T13:25:00+08:00',
  publishedAt: '2026-09-30T13:25:00+08:00',
  seo: {
    title: 'Contact · Magpie Cloud',
    description: 'Book a product demo, or talk to us about your website and content needs.',
    keywords: ['contact Magpie', 'book a demo', 'product inquiry']
  },
  blocks: [
    b('hero', {
      variant: 'page',
      eyebrow: 'CONTACT',
      title: 'Say hello.\nStart a story.',
      subtitle: 'Product questions, partnerships, or just a conversation about your website — we are here.',
      primaryLabel: '',
      primaryHref: '',
      secondaryLabel: '',
      secondaryHref: '',
      note: 'We reply within one business day',
      visual: 'none',
      theme: 'light'
    }),
    b('contact-form', {
      title: 'Tell us about your site.',
      subtitle: 'Fill in the form and we will get back to you shortly.',
      submitLabel: 'Send message',
      successTitle: 'Your message is on its way',
      successText: 'Thanks for reaching out. We will reply within one business day.',
      email: 'hello@xiaoque.com',
      phone: '+86 400 860 2026',
      address: '18F, West Bund AI Tower, Shanghai',
      hours: 'Mon – Fri, 9:00 – 18:00 CST',
      mapStyle: 'orbit'
    })
  ]
})

export default [home, products, solutions, cases, news, about, contact]

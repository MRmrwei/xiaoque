const share = process.env.NUXT_SHARE === '1'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  // The share build is a client-only SPA that is later inlined into one HTML file.
  ssr: !share,
  buildDir: share ? '.nuxt-share' : '.nuxt',
  modules: ['@nuxtjs/i18n'],
  css: ['~/assets/css/base.css'],
  app: {
    head: {
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1',
      title: '小鹊云 · 企业官网与内容云平台',
      htmlAttrs: { lang: 'zh-CN' },
      meta: [
        { name: 'theme-color', content: '#F6F2EA' },
        {
          name: 'description',
          content:
            '小鹊云 — 让品牌官网成为增长引擎。区块化搭建、中英双语、SEO 就绪的企业内容云平台。'
        }
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: 'anonymous' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Instrument+Serif:ital@0;1&family=Manrope:wght@400;500;600;700;800&family=Noto+Serif+SC:wght@500;600;700&display=swap'
        },
      ]
    }
  },
  i18n: {
    locales: [
      { code: 'zh', language: 'zh-CN', name: '中文' },
      { code: 'en', language: 'en', name: 'English' }
    ],
    defaultLocale: 'zh',
    strategy: 'prefix_except_default',
    detectBrowserLanguage: false,
    restructureDir: false,
    vueI18n: './i18n.config.ts'
  },
  typescript: {
    strict: true
  },
  runtimeConfig: {
    // 服务端 SSR 直接请求 Go 后端。
    apiBase: process.env.NUXT_API_BASE || 'http://127.0.0.1:8888/api/v1',
    public: {
      // 浏览器侧通过 Nitro 代理访问，避免跨域。
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api/backend',
      // 单文件分享版没有服务端，继续使用内置演示数据。
      useMockData: share
    }
  },
  routeRules: {
    '/api/backend/**': {
      proxy: `${process.env.NUXT_API_TARGET || 'http://127.0.0.1:8888'}/api/v1/**`
    }
  },
  ...(share
    ? {
        nitro: {
          output: { dir: 'share-output' }
        },
        vite: {
          build: {
            cssCodeSplit: false,
            assetsInlineLimit: 100000000,
            rollupOptions: {
              output: { inlineDynamicImports: true }
            }
          }
        }
      }
    : {})
})

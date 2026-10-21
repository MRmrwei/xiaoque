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
        { rel: 'icon', type: 'image/svg+xml', href: 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgdmlld0JveD0iMCAwIDQwIDQwIj4KICA8cmVjdCB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHJ4PSIxMiIgZmlsbD0iIzE0MTUxOCIvPgogIDxwYXRoIGQ9Ik0xMS4yIDI2LjhjLTEuOS0xLjItMi44LTMuMS0yLjQtNS40LjMtMS45IDEuNC0zLjYgMy00LjcgMS44LTEuMyA0LTEuOCA2LjItMS4zIDEuOS40IDMuNCAxLjYgNC4zIDMuMi44IDEuNCAxIDMuMS41IDQuOC0uNiAyLjEtMi4yIDMuNy00LjMgNC40LTIuNC44LTUgLjYtNy4zLS42WiIgZmlsbD0iI0Y2RjJFQSIvPgogIDxjaXJjbGUgY3g9IjI1LjciIGN5PSIxNS4zIiByPSI1LjEiIGZpbGw9IiNGNkYyRUEiLz4KICA8cGF0aCBkPSJNMzAuNSAxNC4yIDM2IDE2bC01LjIgMi4zLjYtMi4xLS45LTJaIiBmaWxsPSIjRTg1MDNBIi8+CiAgPHBhdGggZD0iTTEyLjMgMjUuNmMxLjQtNCA1LTYuNiA5LjMtNi42IDIuNiAwIDQuOCAxIDYuNCAyLjctMy4zLTEtNi44IDAtOS4xIDIuNS0uOSAxLTEuNSAyLjMtMS43IDMuNi0yLjItLjEtNC0uOS00LjktMi4yWiIgZmlsbD0iI0U4NTAzQSIvPgogIDxwYXRoIGQ9Ik0xMS42IDI1LjggNS4yIDMzLjJsMS41LTguNWMxLjQuOCAzLjEgMS4yIDQuOSAxLjFaIiBmaWxsPSIjRjZGMkVBIi8+CiAgPGNpcmNsZSBjeD0iMjcuMiIgY3k9IjE0LjMiIHI9IjEuMDUiIGZpbGw9IiMxNDE1MTgiLz4KPC9zdmc+Cg==' }
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

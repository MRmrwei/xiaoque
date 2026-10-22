# 小鹊云官网前台

基于 Nuxt 3 的双语企业官网，使用 SSR（服务端渲染）输出页面，面向搜索引擎和首屏体验优化。

官网的页面、导航与站点配置来自 Go 后端；页面内容由管理后台编排并发布。

## 技术栈

- Nuxt 3 + Vue 3 + TypeScript
- Vite
- `@nuxtjs/i18n`（中英文）
- Nitro（SSR 服务端）
- 原生 CSS / 自定义设计系统

> 官网不引入 Element Plus，Element Plus 仅用于 `admin/` 管理后台。

## 功能

- **SSR 渲染**：服务器生成完整 HTML，利于 SEO 与首屏展示
- **中英文**：默认中文无前缀，英文使用 `/en`
- **页面渲染**：通过 `render/page` 获取已发布页面、导航、站点配置
- **动态页面**：管理后台新增的页面路径由 `pages/[...slug].vue` 承载
- **新闻详情**：`/news/:slug`，中英文 slug 独立
- **SEO**：title、description、canonical、hreflang、`<html lang>`、robots
- **404**：未发布或不存在的页面输出 404 + `noindex`
- **联系表单**：提交到 Go 后端，写入后台「表单留言」
- **素材**：图片 / 视频使用素材库 URL
- **品牌配置**：Logo、主色、底色、favicon 跟随站点配置
- **单文件分享版**：可选生成纯静态分享包，仅用于临时预览

## 目录

```text
domain/       领域类型、区块协议
data/         本地演示数据（仅在分享版或内容未接入时使用）
services/     数据访问边界（renderService / contactService / siteService）
composables/  页面加载、SEO、语言切换、表单、Toast
layouts/      站点外壳（头部、页脚、全局配置）
pages/        路由与页面组装
components/   展示组件（区块、布局、基础 UI）
locales/      中英文 UI 词条
assets/css/   设计系统与基础样式
```

规则：

- 页面不直接读取演示数据，通过 `composables` → `services` 获取。
- 区块组件不直接调用接口，数据由 `PageRenderer` 传入。
- 业务规则放在 services / composables，不散落在页面模板里。

## 开发

```bash
npm install
npm run dev        # http://localhost:3000
npm run typecheck  # Nuxt 类型检查
npm run build      # SSR 生产构建，产物在 .output/
npm run preview    # 本地预览生产构建
```

开发时需要 Go 后端，默认地址：

```text
http://127.0.0.1:8888/api/v1
```

## 环境变量

| 变量 | 作用 | 默认值 |
| --- | --- | --- |
| `NUXT_API_BASE` | 服务端 SSR 请求 Go API 的地址 | `http://127.0.0.1:8888/api/v1` |
| `NUXT_PUBLIC_API_BASE` | 浏览器侧 API 前缀 | `/api/backend` |
| `NUXT_API_TARGET` | Nitro `/api/backend/**` 代理目标，构建时写入 | `http://127.0.0.1:8888` |
| `NUXT_SHARE` | 设为 `1` 时使用内置演示数据并生成分享版 | 未设置 |

浏览器请求 `/api/backend/**`，由 Nitro 代理到 Go 后端，避免跨域。

## 路由与多语言

- 中文（默认）：`/`、`/about`、`/contact` 等由后台页面路径决定
- 英文：`/en`、`/en/about`、`/en/contact`
- 新闻详情：`/news/:slug`、`/en/news/:slug`
- 管理后台新增页面时，路径会自动匹配到动态页面路由

## 构建与部署

SSR 生产构建：

```bash
npm run build
node .output/server/index.mjs
```

需要 Node.js 22 运行时，不能把 `.output` 当静态文件直接交给 Nginx。

生产环境推荐：

```text
Nginx / 网关
├── /            → Nuxt SSR（Node）
├── /admin/      → 管理后台静态文件
├── /api/        → Go 后端
└── /uploads/    → 素材文件
```

Docker 部署说明见 [`deploy/docker/xqyun/README.md`](../deploy/docker/xqyun/README.md)。

## 单文件分享版

```bash
NUXT_SHARE=1 npx nuxt generate
node scripts/build-share.mjs
```

生成 `dist-share/index.html`，资源内联、hash 路由，可直接双击打开。它只用于临时分享，不是正式部署产物。

## 说明

- 正式部署必须保持 SSR，不要改成纯 SPA。
- 中文和英文页面独立发布；另一语言未发布时返回 404 + `noindex`。
- 新闻 / 产品 / 案例内容库尚未实现，初版列表使用页面区块的「自定义卡片」。
- 区块协议与后台各自维护一份，新增字段时需同步 `admin/src/domain/blockSchema.ts`。
- 更多项目说明见根目录 [`README.md`](../README.md)。

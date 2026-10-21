# 小鹊云官网 · 前台

基于 **Nuxt 3 + TypeScript + @nuxtjs/i18n** 的双语官网，按《项目方案》的前台技术选型实现。官网不使用 Element Plus（Element Plus 只用于管理后台）。

## 开发

```bash
npm install
npm run dev
```

访问 http://localhost:3000/。

## 校验与构建

```bash
npm run typecheck   # vue-tsc / Nuxt 类型检查
npm run build       # SSR 生产构建，产物在 .output/
npm run preview     # 本地预览生产构建
```

生产部署使用 Node 运行时：

```bash
node .output/server/index.mjs
```

Nginx 将 `/` 反向代理到该 Node 服务；`/en` 走同一服务的前台路由。

## 代码分层

```
domain/        # 领域类型、区块协议（与页面无关）
data/          # 演示数据（页面、新闻、产品、案例、站点设置）
services/      # 数据访问边界：当前读本地 seed，接入 Go API 时只改这一层
composables/   # 业务逻辑：页面加载、SEO、语言切换、表单、Toast
layouts/       # 站点外壳
pages/         # 路由与页面组装
components/    # 纯 UI：区块、布局组件、基础组件
locales/       # 中英文 UI 词条
assets/css/    # 设计系统与基础样式
```

页面只负责路由、数据编排和 SEO；组件只负责展示；数据获取与业务规则集中在 `services/` 和 `composables/`。

## 路由与多语言

- 中文（默认）：`/`、`/products`、`/news/...`
- 英文：`/en`、`/en/products`、`/en/news/...`
- 每个页面输出 `<html lang>`、title/description、canonical 与 hreflang。
- 新闻详情的中文/英文 slug 独立，语言切换会定位到对应译文。

## 单文件分享版（临时预览）

```bash
NUXT_SHARE=1 npx nuxt generate
node scripts/build-share.mjs
```

生成 `dist-share/index.html`：Nuxt 客户端构建 + JS/CSS 全部内联，使用 hash 路由，可直接双击打开，不需要 Node 或 Nginx。`dist-share/xiaoque-website.html` 是同一文件的友好命名副本，`xiaoque-website-preview.zip` 是可直接发送的压缩包。

这是临时分享包，不用于正式部署；正式上线仍使用 `npm run build` 的 SSR 产物。

## 说明

- 当前数据为演示数据，没有接后端；`services/siteService.ts` 是替换为 Go render API 的唯一边界。
- 单文件分享版是客户端 SPA，不参与正式 SSR 部署。

# 小鹊云管理后台 · 站点编排（本地演示）

运营用来搭建官网页面、维护站点外壳的独立 SPA：选择区块、填写内容、实时预览、保存草稿、直接发布。

当前是**本地演示版**：页面、版本、顶部 / 底部导航和站点配置都保存在浏览器 `localStorage`（key：`xiaoque-admin-demo-v1`），不连接后端，也不会更新线上官网。

## 技术栈

Vue 3.5 + TypeScript + Vite 6 + Element Plus + Pinia + Vue Router（hash）+ TipTap + vue-draggable-plus。

## 运行

```bash
npm install
npm run dev        # http://localhost:5174/#/pages
npm run typecheck  # vue-tsc --noEmit
npm run build      # 产物 dist/（admin + preview 两个入口）
```

## 已实现的页面

- `/#/pages` 站点编排，内部用四个 Tab 区分：
  - **页面内容**：页面列表 → 三栏编辑器；筛选、新建（空白或复制首页）、复制、删除、预览、发布 / 下架、翻译配对。中英文独立发布，另一语言未发布时前台不展示对应语言。
  - **顶部导航**：菜单项、头部按钮、公告条、语言入口、滚动吸顶；中文 / EN、桌面 / 手机实时头部预览。
  - **底部导航**：页脚栏目 / 链接、订阅区、联系方式、品牌标语、版权 / ICP、底部链接；实时页脚预览。
  - **站点配置**：官网 Logo、品牌主色、页面底色、深色区块底色、搜索引擎收录；品牌预览。
- `/#/navigation` 旧入口，自动重定向到「站点编排 → 顶部导航」。
- `/#/pages/:id` 页面内容编辑器：
  - 左侧区块库与区块结构：添加、拖拽排序、复制、删除；
  - 中间由 `domain/blockSchema.ts` 自动生成属性表单；字段可按同级值条件显示（头部轮播的图片 / 视频字段）；页面设置面板包含路径、SEO、语言配对；
  - 富文本使用 TipTap：标题、列表、引用、代码块、链接、图片（工具栏上传 / 粘贴 / 拖入）；旧版自定义 HTML 区块加载时自动迁移；
  - 「产品列表 / 新闻列表 / 客户案例」可添加自定义卡片，图片、标签、标题、摘要、要点、按钮文案和跳转链接都能单独填写；卡片留空时回退为自动读取已发布内容；「显示条数」对两种来源都生效；
  - 「行动号召」可用主题底色或上传图片作为底图，并调图片遮罩强度；头部轮播视频模式可上传海报封面；
  - 「图文组合、首屏、头部轮播、视频、联系表单」的图片支持上传（点击或拖拽，自动压缩到 1600px 内的 WebP，也可切回内置图形）；
  - 「常见问题、行动号召、数据统计」等区块的侧栏文案、小标题、竖排小字也已开放编辑；
  - 右侧 iframe 实时预览，支持桌面 / 平板 / 手机；
  - 自动保存（停止编辑 3 秒）、撤销 / 重做、发布前校验、直接发布覆盖线上内容；
  - 中文 / EN 两个语言 Tab，草稿与发布状态互相独立。
- `preview.html` 草稿预览站：复用官网设计系统的区块渲染器，通过 `postMessage` 接收后台推送的页面数据。

## 保存策略

- 页面内容：停止编辑 3 秒自动保存，离开编辑器立即落盘。
- 顶部导航 / 底部导航 / 站点配置：修改后点击「保存设置」写入 `localStorage`，编辑过程中不自动保存。

## 目录

```
src/views/                 PagesView（四个 Tab：页面内容 / 顶部导航 / 底部导航 / 站点配置）、PageEditorView
src/components/navigation/ HeaderNavPanel、FooterNavPanel、FooterLinkRows
src/components/settings/   SiteSettingsPanel（Logo / 品牌色 / 底色 / 收录）
src/components/editor/     区块结构、区块库、页面设置、发布
src/components/form/       schema 表单、字段控件、TipTap 富文本、图片选择 / 上传
src/components/preview/    编辑器内 iframe 预览
src/domain/blockSchema.ts  区块协议与字段定义
src/services/              页面业务规则、站点设置兼容迁移、localStorage 持久化
src/stores/pages.ts        页面 / 站点设置状态（Pinia）
src/data/seed/             演示页面、内容与站点设置
src/preview/site/          预览站及官网区块组件
src/utils/image.ts         图片压缩（最长边 1600px WebP，SVG / GIF 保留原文件）
src/utils/brand.ts         品牌色 → CSS 变量
```

## 说明

- 前台官网在 `frontend/`（Nuxt 3 SSR）。后台与前台目前各自维护一份类型 / 区块 schema，后续可抽到 `packages/` 共享。
- 发布只在本地更新发布态，用于演示直接发布流程；接 Go API 时替换 `services/` 这一层。

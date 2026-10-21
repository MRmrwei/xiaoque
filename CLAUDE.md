# 小鹊项目 · Claude 工作指南

## 开始任何工作前

1. 先读 `docs/项目方案.md`。技术栈、架构、数据模型、SEO 与多语言要求都以它为准；实施状态以本文件为准。
2. 仓库包含三部分：`frontend/`（官网前台，Nuxt 3 SSR）、`admin/`（管理后台，登录已接 Go API，其余数据本地演示）和 `service/`（Go 后端，CloudWeGo Hertz）。`service/` 已完成 Hertz + MySQL/Redis 骨架、配置加载、统一响应、`/ping`、后台认证接口（登录 / 刷新 / 退出）、JWT 中间件与访问日志；页面 / 导航 / 站点 / 留言等业务接口仍未完成。除非用户明确要求，不要修改或运行它。
3. 后端接口与表结构分别以 `docs/Go后端接口设计.md`、`docs/MySQL数据库表设计.md` 为准；初版上线的范围、接口顺序与建表清单见 `docs/初版上线方案.md`；数据库初始化脚本在 `deploy/init.sql`。
4. 回复用中文，直接说明做了什么、怎么验证。预览 / 开发服务用完要关停，不要长期常驻。

## 当前状态（2026-10-21）

- `frontend/`：Nuxt 3 + TypeScript + Vite + `@nuxtjs/i18n`，SSR 官网，已完成演示版，数据为 mock；中英文、SEO、404、联系表单等按方案实现。
- `admin/`：Vue 3 + TypeScript + Vite + Element Plus 的管理后台演示版，已可运行；包含登录页（调用 Go API 登录 + 前端图形验证码）、基础设置分类页（「后台登录」，含验证码开关与登录页 Logo / 名称 / 展示文案）、管理员账号（侧边栏用户区进入，含资料 / 头像 / 修改密码）、站点编排（页面内容 / 顶部导航 / 底部导航 / 站点配置四个 Tab）、草稿预览、直接发布、中英翻译配对、图片上传，以及独立的「表单留言」线索管理。Access Token 与页面数据保存在浏览器 `localStorage`（页面数据与基础设置 key：`xiaoque-admin-demo-v1`，登录会话 key：`xiaoque-admin-auth-v1`，管理员资料 / 密码哈希 key：`xiaoque-admin-profile-v1` / `xiaoque-admin-password-v1`）；Refresh Token 由后端写入 HttpOnly Cookie，前端不可读取。**登录、静默刷新与退出登录已接 Go API，其余页面数据仍是本地演示，不会影响线上（frontend）内容**。
- `service/`（Go 后端）：CloudWeGo Hertz 骨架已完成，监听 `:8888`；已完成 YAML + env 配置（`service/config.yaml`；数据库用户名 / 密码、Redis 密码与 `JWT_SECRET` 只从 `service/.env` 读取）、日志初始化与轮转、MySQL / Redis 初始化与关闭、统一响应（HTTP 200 + `code`）、`GET /ping`、优雅退出。`app/models`（User，表名 `xqyun_users`）、`app/dto`、`app/services/auth`（bcrypt + HS256 JWT + 刷新令牌）、`app/services/adminauth`（登录用例）、`app/repository/user`（用户查询）、`app/repository/session`（Redis 登录会话）、`middleware/admin_auth.go`（Bearer JWT 认证）与 accesslog 访问日志已完成；`POST /api/v1/admin/auth/login` 已接通数据库并可登录，`POST /api/v1/admin/auth/refresh` 与 `POST /api/v1/admin/auth/logout` 已实现刷新令牌轮换与吊销，Refresh Token 通过 HttpOnly Cookie 下发。刷新令牌只存 Redis（`xqyun:refresh:{token_type}:{sha256}`，TTL 7 天），不落 MySQL；`deploy/migration/` 只建 `xqyun_users`，`xqyun` 库中已写入初始管理员 `admin` / `小阙`；当前用户及其余业务接口尚未实现。
- 2026-10-19 收尾：站点编排四个 Tab、导航 / 站点配置手动保存、TipTap 富文本（含图片）、站点品牌配置（Logo / 主色 / 页面底色 / 深色区块底色 / 收录开关）、CTA 图片底图、头部轮播条件字段与视频封面、产品 / 新闻 / 案例「显示条数」统一限制两种来源。
- 2026-10-20：管理后台新增登录页 + 图形验证码、基础设置分类页（目前只有「后台登录」，含验证码开关 / 登录页 Logo、名称、展示文案）、表单留言管理与详情抽屉；统一后台字号（正文 / 标签 ≥12px，辅助说明 ≥11px）与中文错误提示；Go 后端新增登录 service 实现、DTO、模型和前台 / 管理后台路由分组；数据库表名统一 `xqyun_` 前缀。
- 2026-10-21：后台登录页接 Go API（Access Token 本地保存、Refresh Token 走 HttpOnly Cookie）；新增统一请求层、JWT 按 admin / frontend 分组、管理后台 Bearer 中间件、刷新 Cookie、accesslog 访问日志、日志按环境输出与轮转；补充 JWT / 配置 / 中间件测试。同日完成：刷新令牌改存 Redis 并删除 `xqyun_refresh_tokens` 表；数据访问拆为 `repository/user` 与 `repository/session`；管理后台控制器移入 `app/handler/backend`；后台接入静默刷新与退出登录接口。
- 前后台目前各自维护一份区块类型与 schema，尚未抽到 `packages/` 共享；调整区块协议时两边都要检查。
- 设计系统：品牌「小鹊云 / Magpie Cloud」，暖纸色 `#F6F2EA`、墨黑 `#141518`、朱砂红 `#E8503A`，大字号编辑感排版；官网使用自定义组件。
- 旧 Vite 版源码备份：`_backup/frontend-vite-2026-10-17/`。
- 单文件分享版：`frontend/dist-share/index.html`（Nuxt 客户端构建 + 资源内联，hash 路由，可双击打开）。重建：`NUXT_SHARE=1 npx nuxt generate && node scripts/build-share.mjs`。`dist-share/xiaoque-website.html` 和 `xiaoque-website-preview.zip` 是发给别人的成品。它仅用于临时预览，不是正式部署产物。
- 仓库在 `G:\zhujiawei\xiaoque`，当前分支 `dev1019`；`service/`、`admin/`、`frontend/`、`docs/`、`deploy/` 等尚未全部纳入版本管理。除非用户明确要求，不要自行提交或推送。

## 管理后台（admin/）

入口：

- `/#/login` 登录页，填写用户名、密码与图形验证码（验证码可在基础设置关闭）；未登录访问其他后台路由会跳回登录页。
- `/#/pages` 站点编排，内部用四个 Tab 区分：
  - 页面内容：页面列表 → 三栏编辑器
  - 顶部导航：菜单项、头部按钮、公告条、语言入口、吸顶设置
  - 底部导航：页脚栏目、链接、订阅区、联系入口、版权与底部链接
  - 站点配置：官网 Logo、品牌主色、搜索引擎收录与站点信息概览
- `/#/navigation` 旧入口，自动重定向到「站点编排 → 顶部导航」。
- `/#/submissions` 表单留言，侧边栏独立一级菜单「线索 → 表单留言」。
- `/#/settings` 基础设置分类页，侧边栏独立一级菜单「设置 → 基础设置」。
- `/#/settings/login` 后台登录设置，从基础设置分类进入。
- `/#/settings/account` 管理员账号，从侧边栏用户区账号图标进入，修改显示名称、邮箱、手机号、头像与登录密码。

已实现：

- **登录**：`/#/login`，用户名 + 密码 + 图形验证码（可在基础设置中关闭）；调用 `POST /api/v1/admin/auth/login`，初始账号 `admin` / `123456`（本地开发提示，仅 DEV 显示）。成功后只在 `xiaoque-admin-auth-v1` 保存 Access Token 和用户信息；Refresh Token 由后端写入 HttpOnly Cookie，前端不可读取。Access Token 过期时，路由守卫与请求层都会调用 `POST /api/v1/admin/auth/refresh` 静默续期（请求层刷新成功后自动重试原请求一次，并发刷新共享同一次请求）；刷新失败才清除会话并回登录页。侧边栏退出登录会调用 `POST /api/v1/admin/auth/logout` 吊销服务端刷新令牌并清除 Cookie，即使接口失败也会清除本地状态。验证码目前仍是前端生成与校验，后端不校验。
- **基础设置**：独立一级菜单「设置 → 基础设置」；首页为分类列表，目前只有「后台登录」一个分类。进入 `/#/settings/login` 后可开关图形验证码，编辑登录页展示区的 Logo（上传 / 恢复内置图标）、Logo 名称、名称下方小字、眉标、主标题（支持换行）和说明文字，带登录页实时预览；点击「保存设置」写入 `xiaoque-admin-demo-v1`，支持恢复默认。
- **管理员账号**：`/#/settings/account`，从侧边栏用户区账号图标进入；显示登录账号 `admin`（不可修改），可编辑显示名称、邮箱、手机号并保存；支持上传自定义头像（经 `utils/image.ts` 压缩为 320px 以内的 WebP，未上传时保留名称首字文字头像，可随时移除恢复文字头像）；修改密码为本地演示（尚未接后端改密接口），新密码至少 8 个字符、不超过 72 字节；资料和密码哈希分别存 `xiaoque-admin-profile-v1` / `xiaoque-admin-password-v1`，保存后侧边栏名称、头像实时更新。接后端后改为调用账号资料 / 改密接口，密码由服务端 bcrypt 存储。
- **页面列表**：按语言 / 状态 / 关键词筛选；新建（空白或复制首页结构）；复制、删除、草稿预览、发布 / 下架；保留路径与重复路径校验；一键创建 / 前往翻译版本。
- **页面内容**：左侧区块结构（拖拽排序、复制、删除）、中间 schema 自动生成的属性表单、右侧 iframe 实时预览（桌面 / 平板 / 手机）；页面设置含路径、SEO、语言配对。
- **编辑能力**：停止编辑 3 秒自动保存、撤销 / 重做、发布前校验、直接发布覆盖线上内容；页面状态只有已发布 / 未发布，支持发布与下架；中英文草稿和发布状态互相独立。
- **翻译**：切换 EN / 中文时自动创建配对草稿，复制区块结构并清空可翻译文本，保留图片与链接。
- **顶部导航**：站点编排内的独立 Tab；菜单项增删、拖拽排序、显隐开关、中英文名称与跳转路径；站点名称、语言切换、滚动吸顶；公告条（开关、中英文案、高亮文案、链接文案与路径）；顶部按钮支持多个（增删、拖拽排序、显隐、中英文案、路径、主按钮 / 描边按钮 / 文字按钮三种样式，建议 1 主 + 最多 2 辅）；带中文 / EN、桌面 / 手机实时头部预览；修改后点击「保存设置」写入 `localStorage`，不再自动保存。旧的单按钮本地数据会自动迁移为一个主按钮。
- **底部导航**：站点编排内的独立 Tab；页脚栏目与链接的增删、拖拽排序、显隐、中英文文案与路径；订阅区、联系列（CTA）、品牌标语、底部大字标、版权 / ICP、状态提示、底部链接；带中文 / EN、桌面 / 手机实时页脚预览；修改后点击「保存设置」写入 `localStorage`，不再自动保存。导航设置同步到页面草稿预览的头部与页脚。
- **站点配置**：站点编排内的独立 Tab；官网 Logo 上传（自动压缩，空值回退内置品牌图标）、品牌主色（5 个预设 + 取色器，实时换算 `--accent` / `--accent-deep` / `--accent-soft`）、页面底色与深色区块底色（预设 + 取色器，换算 `--paper*` / `--ink*` / `--muted*` / `--line*`）、搜索引擎收录开关（前台据此输出 index / noindex）；带中文 / EN、桌面 / 手机品牌预览；修改后点击「保存设置」写入 `localStorage`。
- **区块内容**：14 种预置区块；富文本使用 TipTap 可视化编辑（标题、列表、引用、代码块、链接、图片、撤销重做），图片支持工具栏上传、粘贴与拖入，统一走 `utils/image.ts` 压缩为 WebP data URL；行动号召支持主题底色或上传图片作为底图，并可调图片遮罩强度；头部轮播按内容类型条件显示字段（图片模式只选图片，视频模式填视频地址并可上传海报封面，封面作为 `<video poster>`）；旧版「自定义 HTML」区块不再出现在区块库，加载时自动迁移为富文本区块。产品列表 / 新闻列表 / 客户案例支持「自定义卡片」—图片、标签、标题、摘要、要点、按钮文案和跳转链接均可编辑；自定义卡片留空时回退为自动读取已发布内容。「显示条数」对三种区块都同时限制自动读取和自定义卡片数量（取排序后的前 N 条）。客户案例的点击详情弹层对两种来源都生效；自定义卡片的链接作为弹层内的「查看完整案例」按钮。
- **图片上传**：图文组合、首屏、头部轮播、视频、联系表单等图片字段支持本地上传；自动等比压缩到最长边 1600px 的 WebP（SVG / GIF 保留原文件），以 data URL 随页面数据保存。上传控件在 `src/components/form/ImagePicker.vue`，处理逻辑在 `src/utils/image.ts`。
- **预览机制**：后台 `PreviewFrame` 通过 iframe 加载 `preview.html`，用 `postMessage` 推送页面、语言、内容集合与站点设置；`preview/site/` 是官网设计系统的副本渲染器。预览中的站内链接被禁用（避免 iframe 跳走），链接值仍会保存。
- **发布**：仅在本地更新发布态用于演示直接发布，不会写回 `frontend/`，也不会触发真实部署。
- **表单留言**：独立一级菜单「线索 → 表单留言」；统计全部 / 待处理 / 处理中 / 已完成，按关键词、状态、语言筛选；表格列表 + 详情抽屉；可标记处理中、已完成、垃圾或恢复待处理。演示留言在 `src/data/seed/submissions.ts`，状态与页面数据一起存 `localStorage`。

运行：

```bash
cd admin
npm run dev          # http://localhost:5174/#/pages
npm run typecheck    # vue-tsc --noEmit
npm run build        # 产物 dist/，包含 admin + preview 两个入口
```

- 登录接口需要 Go 后端（默认 `http://127.0.0.1:8888`）；Vite 开发 / 预览把 `/api` 代理到该地址，可用 `VITE_API_PROXY_TARGET` 覆盖。
- `src/services/api.ts` 默认自动带 `Authorization: Bearer`，统一解析 `{code,message,data}`；收到 `40100` / `40101` 时先静默刷新 Access Token 并重试原请求一次，刷新失败才清除登录会话并跳回登录页。支持 `query`、超时与 `get/post/put/patch/delete`。公开接口（登录、刷新、退出）传 `{ auth: false }`。
- 生产构建可用 `VITE_API_BASE_URL` 指定 API 前缀（默认 `/api/v1`）；前后端不同域时需要 Go 侧 CORS 或反向代理。

后台关键目录：

```
src/views/                 LoginView / BasicSettingsView（分类页）/ LoginSettingsView / AccountSettingsView / PagesView（含四个 Tab）/ PageEditorView / SubmissionsView
src/services/auth.ts       Go API 登录、会话与本地账号演示资料
src/services/api.ts        请求封装：Bearer、统一响应、40100 清理会话、超时与 query
src/data/seed/loginSettings.ts 登录页验证码与展示文案默认值
src/components/navigation/ HeaderNavPanel / FooterNavPanel / FooterLinkRows
src/components/editor/     区块结构、区块库、页面设置、发布
src/components/form/       schema 表单、字段控件、富文本、图片选择 / 上传
src/components/preview/    编辑器 iframe 预览
src/domain/blockSchema.ts  区块协议与字段定义（后台侧）
src/services/              页面业务规则、站点设置兼容迁移、localStorage 持久化
src/stores/pages.ts        页面、站点设置状态（Pinia）
src/data/seed/             演示页面、内容、站点设置（含顶部 / 底部导航）与表单留言
src/preview/               预览站入口、头部 / 页脚预览与官网区块组件
```

## Go 后端（service/）

- 框架 **CloudWeGo Hertz**（不是 Gin），Go 1.25，模块 `suki/xqyun`，监听 `:8888`。
- `config.yaml` 放非敏感配置；数据库用户名 / 密码、Redis 密码与 `JWT_SECRET` 只从 `.env` 读取（复制 `.env.example`），`.env` 不提交。`JWT_SECRET` 至少 32 字节，缺失时服务启动失败。
- MySQL 8 + GORM、Redis；初始化与关闭在 `app/core/application/`，启动装配在 `app/core/bootstrap/`。
- Redis 是登录会话的必需依赖：`redis.enabled=false` 时服务启动失败；生产环境建议开启 AOF，否则 Redis 重启会丢失全部刷新令牌（用户需重新登录）。
- 表模型统一放在 `service/app/models/`，数据库表名统一加 `xqyun_` 前缀（如 `xqyun_users`），通过模型的 `TableName()` 返回。`User` 包含 `email` / `mobile` / `avatar`。刷新令牌不是表模型，存 Redis。
- DTO 统一放在 `app/dto/`：`AdminLoginRequest` / `AdminLoginResponse`、`FrontendLoginRequest` / `FrontendLoginResponse`；service 接口不依赖 DTO。
- 认证接口在 `app/services/auth/`：
  - `auth.go`：`Authenticator`、`Credentials`、`TokenPair`、`LoginResult`；
  - `impl/admin.go`：`AdminLogin` 通过数据访问层查询 `xqyun_users`、校验 bcrypt、检查 `status`、签发 HS256 JWT Access Token（默认 15 分钟）、生成刷新令牌（SHA-256 写入 Redis，默认 TTL 7 天，明文通过 HttpOnly Cookie 下发）、更新 `last_login_at` / `last_login_ip`；
  - `password.go` / `token.go`：bcrypt 密码哈希与导出的 `SignAccessToken` / `NewRefreshToken` 工具；
  - 错误提示中文：`ErrInvalidCredentials` / `ErrAccountDisabled`。
- JWT claims 包含 `sub` / `username` / `iss` / `aud` / `token_type` / `iat` / `exp`；现阶段 `aud = admin-api`、`token_type = admin`，前台预留 `frontend-api` / `frontend`。Redis 会话的 key 与值都带 `token_type`，避免管理后台与前台用户 ID 冲突。
- 数据访问层按聚合拆分：`app/repository/user/` 的 `UserRepository` 负责用户查询（`FindByUsername` / `FindByID`）与最近登录信息更新（`UpdateLoginState`）；`app/repository/session/` 的 `Repository` 负责登录会话（`Save`、`FindByHash`、`Rotate`、`Revoke`），实现基于 Redis，key 为 `xqyun:refresh:{token_type}:{sha256}`，TTL 与 `refresh_ttl` 对齐。GORM 查询与 Redis 操作不进入 service / handler。
- `app/services/adminauth/service.go`：管理后台登录用例层，负责输入校验、认证错误映射与 `AdminLoginResponse` 组装；`app/handler/backend/login.go` 只做请求绑定、调用用例和输出统一响应。
- 依赖装配：`NewLoginController` 只保存 `Application`；`Login` 方法内用 `authimpl.NewAdminLogin` 创建认证实现，再用 `adminauth.NewLoginService` 包装成用例服务，避免在构造函数或路由注册阶段提前组装。
- 管理后台认证中间件在 `middleware/admin_auth.go`：校验 `Authorization: Bearer <token>`，验证 JWT 签名、`exp`、`iss`、`aud`、`token_type`，成功后把 claims 写入请求上下文；失败返回 HTTP 200 + `40100`。受保护路由注册在 `router/admin.go` 的 `authorized` 子组。
- JWT 配置按客户端分组：`config.yaml` 的 `jwt.clients.admin` / `jwt.clients.frontend`，各自包含 `audience` / `token_type` / `access_ttl` / `refresh_ttl`；`issuer` 与 `JWT_SECRET` 全局共享，具体客户端可用 `JWT_<CLIENT>_*` 环境变量覆盖。应用层通过 `application.JWT().Client(name)` 暴露给认证服务与中间件。
- Refresh Token Cookie 配置在 `jwt.cookie`：`name` / `path` / `domain` / `secure` / `same_site`，可用 `JWT_COOKIE_*` 覆盖；本地 HTTP 开发 `secure=false`，生产 HTTPS 必须设为 `true`。
- Cookie 读写工具在 `app/utils/cookie`：`SetRefresh` / `ClearRefresh` / `RefreshValue`；handler 只调用工具，不直接拼 Cookie 属性。
- 路由：`router.go` 的 `customizedRegister` 注册 accesslog 与 `/api/v1/frontend`、`/api/v1/admin` 两组；`router/admin.go` 注册公开的 `POST /auth/login` / `/auth/refresh` / `/auth/logout`，以及需要 Bearer Token 的 `authorized` 子组；`router/frontend.go` 预留公开接口；`router_gen.go` 目前只调用 `customizedRegister`（重新生成 hz 后需保留这些分组）。
- 刷新令牌接口：`app/services/auth/impl/admin_session.go` 实现 `auth.TokenRefresher` / `auth.Logouter`，校验哈希、过期、吊销与账号状态后轮换令牌（旧令牌立即失效，并发重放返回 `40100`）；`app/services/adminauth` 增加 `Refresh` / `Logout` 用例，`app/handler/backend/refresh.go`、`login_out.go` 只读写 Cookie 与组装响应。`POST /auth/refresh` 失败时仅对登录状态失效清 Cookie，Redis 临时故障保留 Cookie；`POST /auth/logout` 幂等，未携带令牌也返回成功并清 Cookie。路由中这三个接口都在 Bearer 中间件之外。
- **当前未完成**：当前用户接口（A4）尚未实现；页面、导航、站点配置、表单留言接口与其余业务表尚未完成。
- 统一响应 `{code, message, data}`：所有响应（包括失败）HTTP 200，成功 / 失败只按 `code` 判断。`0` 成功，`40000` 参数错误、`40100` 未登录 / 令牌无效、`40300` 无权限、`40400` 不存在、`50000` 服务端错误。
- 日志统一使用 `github.com/cloudwego/hertz/pkg/common/hlog`；handler 中优先用 `hlog.CtxErrorf(ctx, ...)` 携带请求上下文，不使用标准库 `log`。
- 日志初始化在 `app/core/logger`：`APP_ENV` 为 `development` / `dev` 等时输出控制台；为 `production` / `prod` / `release` 时写入 `log.dir/app.log`，使用 lumberjack 按大小轮转（默认 100MB、保留 30 份、30 天、压缩），可通过 `LOG_*` 环境变量覆盖。
- 开发环境终端日志按等级着色（Trace / Debug / Info / Notice / Warn / Error / Fatal）；仅在 stdout 是终端且 `log.color` / `LOG_COLOR` 开启时生效，`NO_COLOR` 可关闭；重定向输出和生产文件日志不带颜色。
- HTTP 访问日志使用 `hertz-contrib/logger/accesslog`，在 `router.go` 的 `customizedRegister` 注册；通过 `hlog.CtxInfof` 输出到与控制台 / 文件相同的日志。格式包含时间、状态、耗时、客户端 IP、方法、路径和 query；**不要加入 `${body}`、Authorization 或 Cookie，登录请求体包含密码**。
- 时间字段统一 Unix 时间戳（秒，`int64`），JSON 中 ID 用字符串。
- 接口按客户端分组：官网前台 `/api/v1/frontend/`，管理后台 `/api/v1/admin/`；详细契约见 `docs/Go后端接口设计.md`。
- `app/handler/` 是**控制器包（controller 层）**，专门处理 HTTP 请求：绑定 / 校验参数、调用 service 用例、组装统一响应；不写业务逻辑、不直接访问 GORM、不直接拼接 Cookie / Header 细节。`app/handler/backend/` 是管理后台接口的控制器包（登录、刷新令牌、退出登录等），后续管理后台业务控制器也放在这里。
- 新增接口：在 `app/handler/` 新建控制器（嵌入 `*controller.BaseController`），在 `router.go` 的 `customizedRegister` 注册，用 `response.Success` / `response.Fail` 输出；业务逻辑放 service 层，handler 只做参数解析与响应组装。
- 依赖：`golang.org/x/crypto`（bcrypt）、`github.com/redis/go-redis/v9`（刷新令牌会话）、`gopkg.in/natefinch/lumberjack.v2`（日志轮转）、`github.com/hertz-contrib/logger/accesslog`（访问日志）。

## 技术约定

- 官网必须保持 **SSR（Nuxt 3）**，不要改回纯 SPA；SEO 依赖服务端输出。
- 多语言：默认中文无前缀，英文 `/en`；`@nuxtjs/i18n` 的 `prefix_except_default`。
- **Element Plus 只能用于 `admin/`**，官网（`frontend/`）不要引入。
- 全量 TypeScript：`.ts` 与 `<script setup lang="ts">`，组件 props 要有类型；`frontend` 和 `admin` 提交前各自的 `npm run typecheck` 必须为 0 error。
- 后台演示数据只存 `localStorage`；不要把本地演示逻辑当成生产后端，也不要直接改 `frontend` 的 mock 数据来“模拟”后台发布。
- 上传图片必须走 `ImagePicker` / `utils/image.ts`，不要直接塞入未压缩的原图 data URL（会撑爆 localStorage）。
- 新增 / 修改区块字段时，检查后台 `admin/src/domain/blockSchema.ts` 与前台 `frontend` 对应类型 / 组件是否都需要同步；目前没有共享包。
- Schema 字段支持 `visibleWhen`（按同级字段值条件显示，见头部轮播的图片 / 视频字段）与 `uploadOnly`（只允许上传、不显示内置图形，见 CTA 底图 / 视频封面）；条件判断统一用 `isFieldVisible`，不要在组件里另写一套。
- 旧版 `html` 区块字段 `hidden: true`，加载时由 `migrateLegacyBlocks` 转为 `richtext`；不要重新把它暴露到区块库。
- **不要用过小字号**：后台 / 官网的正文与标签不小于 12px，辅助说明不小于 11px；避免 8–10px 的装饰性小字。统计卡片、表头、状态、详情标签等优先保证可读性。
- **错误提示优先用中文**：Go 的 `errors.New`、接口返回的 `message`、表单校验、后台 / 官网界面提示都用中文；只有日志上下文、第三方原始错误、技术排查信息可以保留英文。启动、配置、数据库 / Redis 等基础设施错误保持英文即可，不强制翻译。

## 前台业务与页面分层

```
domain/       领域类型、区块协议，不依赖页面
data/         演示数据（页面、新闻、产品、案例、站点设置）
services/     数据访问边界；将来接 Go API 只改这一层
composables/  业务逻辑：页面加载、SEO、语言切换、表单、Toast
layouts/      站点外壳
pages/        路由、页面组装、SEO 调用
components/   纯展示组件（区块、布局、基础 UI）
locales/      中英文 UI 词条
assets/css/   设计系统与基础样式
```

规则：

- 页面不直接读取 seed 数据；通过 `composables` → `services` 获取。
- 区块组件不直接 import `services` / `data`；数据源区块（产品/新闻/案例）通过 `PageRenderer` 传入 `collections`。
- 业务规则（发布态过滤、语言配对、表单提交、SEO）放在 services / composables，不散落在页面或 UI 组件里。
- 新增区块类型时，在 `domain/blocks.ts` 加类型，在 `components/blocks/` 加组件，并在 `BlockRenderer` 注册。

## SEO / 多语言

- `usePageSeo` 负责 title、description、canonical、hreflang（zh-CN / en / x-default）、`<html lang>`、robots。
- 普通页面 404：`SitePageView` / 新闻详情内联渲染 404 并设置 `setResponseStatus(..., 404)`，带 `noindex`；不要为了避免 404 改成客户端 200。
- 新闻详情中英 slug 独立；语言切换走 `useLocaleSwitch`，会查找 translation group 的对应 slug。
- 中英独立发布：另一语言未发布时，对应 URL 返回 404 + noindex；语言切换、hreflang、sitemap 只指向已发布译文。
- 内容只展示 `status: 'published'` 的记录。

## 常用命令

```bash
# 官网前台
cd frontend
npm run dev          # http://localhost:3000
npm run typecheck
npm run build        # SSR 生产构建，产物 .output/
node .output/server/index.mjs

# 管理后台
cd admin
npm run dev          # http://localhost:5174/#/pages
npm run typecheck
npm run build

# Go 后端（需要本地 MySQL / Redis，配置见 config.yaml）
cd service
go run .             # http://127.0.0.1:8888/ping
go build ./...
go vet ./...
```

## 改动后的验证清单

**通用**

1. 改动 `frontend/` 或 `admin/` 时，在该项目执行 `npm run typecheck` 与 `npm run build`，均须通过。
2. 改动 `service/` 时执行 `go build ./...` 与 `go vet ./...`，均须通过；`config.yaml` 不得写入密码。
3. 中英文各看一遍；检查桌面与移动端布局。
4. 验证完关停 dev / preview / Go 服务；不要留下长期运行端口。

**前台（frontend/）**

1. 检查 SSR HTML：页面标题、`<html lang>`、canonical、hreflang。
2. 语言切换（普通页 + 新闻详情）、联系表单成功态、404。
3. 预览 / 发布规则不被破坏。

**后台（admin/）**

1. 登录页可用：Go 后端已启动且 Vite 代理正常；错误密码 / 错误验证码提示、正确登录跳转、未登录访问其他路由跳回登录页、退出登录后会话清除。登录后 localStorage 不应出现 refreshToken，DevTools Cookie 中应存在 HttpOnly 的 `xqyun_admin_refresh`。退出登录应调用 `/admin/auth/logout` 并清除 Cookie；Access Token 过期后刷新页面或发请求应静默续期，不应直接跳登录页。
2. 基础设置可用：首页显示分类列表，能进入「后台登录」；管理员账号从侧边栏用户区账号图标进入；验证码开关、登录页 Logo / 名称 / 眉标 / 标题 / 说明文字保存后登录页同步生效；管理员资料保存后侧边栏名称更新；上传头像后账号页与侧边栏同步显示，移除后回退文字头像；本地改密只更新 localStorage，不影响 Go API 登录密码；恢复默认正常。
3. 页面列表 → 编辑器 → 预览 → 自动保存 → 直接发布链路可用。
4. 翻译创建后文本清空、图片与链接保留、中英配对正确。
5. 区块的图片上传、自定义卡片（产品 / 新闻 / 案例）正常渲染并持久化。
6. 顶部 / 底部导航与站点配置修改后，页面草稿预览的头部、页脚与品牌色同步更新；点击「保存设置」后刷新页面设置仍在。
7. 站点编排的四个 Tab（页面内容 / 顶部导航 / 底部导航 / 站点配置）切换正常，旧入口 `/#/navigation` 能重定向到顶部导航。
8. 站点配置修改 Logo / 主色 / 页面底色 / 深色区块底色后，官网头部、页脚、强调按钮与整体背景同步生效；关闭搜索引擎收录时前台输出 `noindex`。
9. 新字段或旧本地数据缺失字段时有默认值，不白屏。
10. 表单留言：列表、筛选、详情抽屉与状态标记可用；刷新后状态仍在，「重置演示数据」能恢复种子留言。

**Go 后端（service/）**

1. 执行 `go build ./...` 与 `go vet ./...`。
2. 登录相关改动先确认 `xqyun_users` 表已创建、Redis 已启动（`redis.enabled=true`）、`JWT_SECRET` 已配置；否则只能做编译级验证，不能声称登录已联调通过。
3. 新增路由后确认 `/api/v1/frontend` 与 `/api/v1/admin` 分组没有被 hz 重新生成覆盖。

## 不要做

- 不要把官网改成 SPA，不要移除 SSR。
- 不要在官网引入 Element Plus。
- 不要把业务逻辑写进页面模板或区块组件。
- 不要直接修改 mock 数据来绕开 services 的展示规则。
- 不要修改或启动 `service/`（除非用户明确要求）。
- 不要在业务表未建、`JWT_SECRET` 未配置、Redis 未启动、handler 未接 service 时声称登录接口已经可用。
- 不要把 `admin/` 的本地演示数据当成真实线上数据；后台发布不会影响 `frontend/`。
- 不要用未压缩原图或超出配额的内容写 `localStorage`。
- 不要在访问日志中记录 `${body}`、Authorization、Cookie 等敏感内容；登录接口的请求体包含密码。
- 不要删除 `_backup/` 和 `frontend/dist-share/`。

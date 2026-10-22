# 小鹊云 Docker 部署包

这个目录可以**整个上传**到服务器，不需要源码，也不需要 `docker build`。

## 上传

把本目录内容上传到服务器：

```
/xqyun/
```

最终结构：

```
/xqyun/
├── docker-compose.yml
├── .env.example
├── init.sql
├── migration/
├── data/
│   ├── nginx/conf/default.conf
│   ├── release/
│   │   ├── admin/                 # 管理后台静态文件
│   │   ├── frontend/.output/      # 官网 Nuxt SSR
│   │   └── backend/
│   │       ├── xqyun              # Linux Go 二进制
│   │       └── config.yaml
│   ├── uploads/                   # 素材（自动生成）
│   └── logs/                      # 后端日志（自动生成）
└── README.md
```

## 启动

```bash
cd /xqyun
chmod +x data/release/backend/xqyun
mkdir -p data/uploads data/logs
cp .env.example .env
# 编辑 .env：MySQL 密码、JWT_SECRET、UPLOAD_PUBLIC_BASE_URL
# 当前 Redis 不设置密码（REDIS_PASSWORD 留空）
docker compose up -d
```

不需要指定 `-f`，默认使用同目录的 `docker-compose.yml`。

## 访问

- 官网：`http://服务器IP/`
- 管理后台：`http://服务器IP/admin/`
- 健康检查：`http://服务器IP/ping`

初始管理员 `admin`，初始密码由部署时提供，登录后请立即修改密码。

## 说明

- 使用官方镜像：`nginx:1.27-alpine`、`node:22-alpine`、`alpine:3.22`、`mysql:8.0`、`redis:7-alpine`、`migrate/migrate:v4.18.3`，首次启动会自动拉取。
- 管理后台、官网、Go 服务都是挂载 `data/release` 下的构建产物，没有构建步骤。
- 数据库首次启动执行 `init.sql` 建库，随后 `migrate` 服务执行 `migration/` 下的迁移。
- 修改 Nginx 配置后：
  ```bash
  docker compose exec gateway nginx -t
  docker compose exec gateway nginx -s reload
  ```
- 更新版本时，替换 `data/release/admin`、`data/release/frontend/.output`、`data/release/backend/xqyun`，然后：
  ```bash
  docker compose restart backend frontend gateway
  ```
- 数据都在 `data/mysql`、`data/redis`、`data/uploads`、`data/logs`，删除容器不会丢数据。
- 开启 HTTPS 后把 `.env` 里的 `JWT_COOKIE_SECURE` 改为 `true`，并把 `UPLOAD_PUBLIC_BASE_URL` 改成 HTTPS 域名。
- 当前 Go 二进制为 **linux/amd64**；ARM 服务器需要重新编译。

import { proxyRequest } from 'h3'

/**
 * 浏览器侧 /api/backend/** 的运行时代理。
 *
 * 早期用 nuxt.config 的 routeRules.proxy 实现，目标地址会在构建时写死到产物里，
 * 容器内会指向 127.0.0.1:8888（前端容器没有这个端口）而返回 502。
 * 改为服务端路由后，目标取自 runtimeConfig.apiBase，可在运行期用 NUXT_API_BASE 覆盖，
 * 与 SSR 请求使用同一个后端地址。
 */
export default defineEventHandler((event) => {
  const { apiBase } = useRuntimeConfig()
  const base = String(apiBase).replace(/\/+$/, '')
  const url = getRequestURL(event)
  const path = url.pathname.replace(/^\/api\/backend\/?/, '')

  return proxyRequest(event, `${base}/${path}${url.search}`)
})

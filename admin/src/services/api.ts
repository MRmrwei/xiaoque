/**
 * 管理后台 HTTP 请求封装。
 *
 * - 默认请求 /api/v1，开发 / 预览环境由 Vite 代理到 Go 后端。
 * - 默认携带 Authorization: Bearer <accessToken>。
 * - 统一解析 { code, message, data }，code !== 0 时抛出 ApiError。
 * - 收到 40100 / 40101 时先用 HttpOnly Cookie 静默刷新一次 Access Token 并重试原请求；
 *   刷新失败才清除本地会话并跳转登录页。
 */

const configuredBase = import.meta.env.VITE_API_BASE_URL

export const API_BASE_URL = (configuredBase || '/api/v1').replace(/\/+$/, '')

const DEFAULT_TIMEOUT_MS = 15_000
const UNAUTHORIZED_CODES = new Set([40100, 40101])

export class ApiError extends Error {
  code: number

  constructor(code: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
  }
}

interface ApiEnvelope<T> {
  code: number
  message: string
  data?: T
}

export type QueryParams = Record<string, string | number | boolean | null | undefined>

export interface ApiRequestOptions extends Omit<RequestInit, 'body'> {
  /** 是否自动携带 Access Token，默认 true。登录等公开接口传 false。 */
  auth?: boolean
  /** 超时时间，毫秒，默认 15 秒。 */
  timeoutMs?: number
  /** 查询参数，自动拼接到 URL。 */
  query?: QueryParams
  /** 普通对象会自动 JSON 序列化；FormData / Blob / string 等原样发送。 */
  body?: unknown
}

interface ApiAuthProvider {
  getAccessToken: () => string | null
  /** 用 HttpOnly Cookie 换取新的 Access Token；刷新失败返回 null。 */
  refreshAccessToken: () => Promise<string | null>
  onUnauthorized: () => void
}

let authProvider: ApiAuthProvider | null = null

/** 由 auth 模块注入 token 读取与 401 处理逻辑，避免 api 与 auth 循环依赖。 */
export function configureApiAuth(provider: ApiAuthProvider | null): void {
  authProvider = provider
}

function buildPath(path: string, query?: QueryParams): string {
  if (!query) return path

  const params = new URLSearchParams()
  Object.entries(query).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      params.append(key, String(value))
    }
  })

  const search = params.toString()
  if (!search) return path
  return `${path}${path.includes('?') ? '&' : '?'}${search}`
}

function prepareBody(body: unknown): { value: BodyInit | null; isJson: boolean } {
  if (body === undefined || body === null) {
    return { value: null, isJson: false }
  }
  if (typeof body === 'string') {
    return { value: body, isJson: false }
  }
  if (typeof FormData !== 'undefined' && body instanceof FormData) {
    return { value: body, isJson: false }
  }
  if (typeof Blob !== 'undefined' && body instanceof Blob) {
    return { value: body, isJson: false }
  }
  if (typeof URLSearchParams !== 'undefined' && body instanceof URLSearchParams) {
    return { value: body, isJson: false }
  }
  if (body instanceof ArrayBuffer || ArrayBuffer.isView(body)) {
    return { value: body as BodyInit, isJson: false }
  }
  return { value: JSON.stringify(body), isJson: true }
}

export async function apiRequest<T>(
  path: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  return performRequest<T>(path, options, false)
}

async function performRequest<T>(
  path: string,
  options: ApiRequestOptions,
  retried: boolean
): Promise<T> {
  const {
    auth = true,
    timeoutMs = DEFAULT_TIMEOUT_MS,
    query,
    body: rawBody,
    headers: rawHeaders,
    signal: externalSignal,
    ...rest
  } = options

  const headers = new Headers(rawHeaders)
  headers.set('Accept', 'application/json')

  const prepared = prepareBody(rawBody)
  if (prepared.isJson && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  if (auth) {
    const token = authProvider?.getAccessToken()
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  }

  const controller = new AbortController()
  let timedOut = false
  const timer = window.setTimeout(() => {
    timedOut = true
    controller.abort()
  }, timeoutMs)

  if (externalSignal) {
    if (externalSignal.aborted) {
      controller.abort()
    } else {
      externalSignal.addEventListener('abort', () => controller.abort(), { once: true })
    }
  }

  let response: Response
  try {
    response = await fetch(`${API_BASE_URL}${buildPath(path, query)}`, {
      ...rest,
      // Refresh Token 在 HttpOnly Cookie 中，需要浏览器自动携带。
      credentials: rest.credentials ?? 'include',
      body: prepared.value,
      headers,
      signal: controller.signal
    })
  } catch (error) {
    if (timedOut) {
      throw new ApiError(0, '请求超时，请稍后重试')
    }
    if (error instanceof Error && error.name === 'AbortError') {
      throw new ApiError(0, '请求已取消')
    }
    throw new ApiError(0, '无法连接服务器，请确认 Go 后端已启动')
  } finally {
    window.clearTimeout(timer)
  }

  const raw = await response.text()
  let payload: ApiEnvelope<T> | null = null
  try {
    payload = raw ? (JSON.parse(raw) as ApiEnvelope<T>) : null
  } catch {
    payload = null
  }

  if (!payload || typeof payload.code !== 'number') {
    if (auth && response.status === 401) {
      if (await retryAfterRefresh(retried)) {
        return performRequest<T>(path, options, true)
      }
      authProvider?.onUnauthorized()
    }
    throw new ApiError(response.status, '服务器返回了无法解析的数据')
  }

  if (auth && (UNAUTHORIZED_CODES.has(payload.code) || response.status === 401)) {
    // Access Token 过期时先用刷新令牌换新 token，再重试一次原请求。
    // 刷新由 auth 模块单飞（并发请求共享同一次刷新），这里只负责重试。
    if (await retryAfterRefresh(retried)) {
      return performRequest<T>(path, options, true)
    }
    authProvider?.onUnauthorized()
  }

  if (payload.code !== 0) {
    throw new ApiError(payload.code, payload.message || '请求失败')
  }

  return payload.data as T
}

/** 调用注入的刷新方法；已经重试过或没有注入刷新能力时返回 false。 */
async function retryAfterRefresh(retried: boolean): Promise<boolean> {
  if (retried || !authProvider?.refreshAccessToken) return false
  try {
    const token = await authProvider.refreshAccessToken()
    return Boolean(token)
  } catch {
    return false
  }
}

export interface ApiDownloadResult {
  blob: Blob
  filename: string
}

/**
 * 下载后端文件（当前用于留言 CSV 导出）。
 * 与 JSON 请求共用 Bearer、Cookie 与静默刷新逻辑，返回文件 Blob 和文件名。
 */
export async function apiDownload(
  path: string,
  options: { query?: QueryParams; timeoutMs?: number } = {}
): Promise<ApiDownloadResult> {
  return performDownload(path, options, false)
}

async function performDownload(
  path: string,
  options: { query?: QueryParams; timeoutMs?: number },
  retried: boolean
): Promise<ApiDownloadResult> {
  const { query, timeoutMs = 30_000 } = options
  const headers = new Headers()
  headers.set('Accept', 'text/csv, application/json')

  const token = authProvider?.getAccessToken()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const controller = new AbortController()
  let timedOut = false
  const timer = window.setTimeout(() => {
    timedOut = true
    controller.abort()
  }, timeoutMs)

  let response: Response
  try {
    response = await fetch(`${API_BASE_URL}${buildPath(path, query)}`, {
      headers,
      credentials: 'include',
      signal: controller.signal
    })
  } catch (error) {
    if (timedOut) {
      throw new ApiError(0, '导出超时，请稍后重试')
    }
    if (error instanceof Error && error.name === 'AbortError') {
      throw new ApiError(0, '导出已取消')
    }
    throw new ApiError(0, '无法连接服务器，请确认 Go 后端已启动')
  } finally {
    window.clearTimeout(timer)
  }

  const contentType = response.headers.get('Content-Type') || ''
  if (!contentType.includes('text/csv')) {
    const payload = (await response.json().catch(() => null)) as ApiEnvelope<unknown> | null
    if (response.status === 401 || (payload && UNAUTHORIZED_CODES.has(payload.code))) {
      if (!retried && (await retryAfterRefresh(retried))) {
        return performDownload(path, options, true)
      }
      authProvider?.onUnauthorized()
    }
    if (payload && typeof payload.code === 'number') {
      throw new ApiError(payload.code, payload.message || '导出失败，请稍后重试')
    }
    throw new ApiError(response.status, '导出失败，请稍后重试')
  }

  return {
    blob: await response.blob(),
    filename: filenameFromDisposition(response.headers.get('Content-Disposition'))
  }
}

function filenameFromDisposition(header: string | null): string {
  if (!header) return 'submissions.csv'
  const match = /filename\*?=(?:UTF-8'')?"?([^";]+)"?/i.exec(header)
  if (!match) return 'submissions.csv'
  try {
    return decodeURIComponent(match[1])
  } catch {
    return match[1]
  }
}

type RequestOptions = Omit<ApiRequestOptions, 'method' | 'body'>

/** 常用方法快捷封装。 */
export const api = {
  get<T>(path: string, options: RequestOptions = {}): Promise<T> {
    return apiRequest<T>(path, { ...options, method: 'GET' })
  },
  post<T>(path: string, body?: unknown, options: RequestOptions = {}): Promise<T> {
    return apiRequest<T>(path, { ...options, method: 'POST', body })
  },
  put<T>(path: string, body?: unknown, options: RequestOptions = {}): Promise<T> {
    return apiRequest<T>(path, { ...options, method: 'PUT', body })
  },
  patch<T>(path: string, body?: unknown, options: RequestOptions = {}): Promise<T> {
    return apiRequest<T>(path, { ...options, method: 'PATCH', body })
  },
  delete<T>(path: string, options: RequestOptions = {}): Promise<T> {
    return apiRequest<T>(path, { ...options, method: 'DELETE' })
  }
}

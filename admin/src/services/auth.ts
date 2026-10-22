/**
 * 管理后台登录与当前管理员资料。
 *
 * 已接 Go API：
 * - POST /api/v1/admin/auth/login
 * - POST /api/v1/admin/auth/refresh（HttpOnly Cookie 静默续期）
 * - POST /api/v1/admin/auth/logout
 * - PUT  /api/v1/admin/user/profile
 * - PUT  /api/v1/admin/user/password
 *
 * Access Token 保存在 localStorage；Refresh Token 由后端写入 HttpOnly Cookie，前端不可读取。
 */
import { ref } from 'vue'
import type { AdminLoginResponse, AdminProfile, AdminUser } from '@/types/admin'
import { ApiError, api, configureApiAuth } from '@/services/api'
import { changeAdminPassword, updateProfile } from '@/services/userApi'

export interface AuthSession {
  username: string
  name: string
  loggedInAt: string
  accessToken: string
  /** Access Token 过期时间，Unix 秒。 */
  expiresAt: number
  user: AdminUser
}

export interface ChangePasswordResult {
  ok: boolean
  message: string
}

/** 账号页「恢复默认资料」使用的占位值；登录账号仍以 Go 后端为准。 */
export const demoCredentials = {
  username: 'admin',
  name: '小阙'
} as const

export const AUTH_STORAGE_KEY = 'xiaoque-admin-auth-v1'

const MIN_PASSWORD_CHARS = 8
const MAX_PASSWORD_BYTES = 72

const DEFAULT_PROFILE: AdminProfile = {
  username: demoCredentials.username,
  name: demoCredentials.name,
  email: '',
  mobile: '',
  avatar: ''
}

let memorySession: AuthSession | null = null
let memoryProfile: AdminProfile | null = null

function readStorage<T>(key: string): T | null {
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return null
    return JSON.parse(raw) as T
  } catch {
    return null
  }
}

function writeStorage(key: string, value: unknown): boolean {
  try {
    localStorage.setItem(key, JSON.stringify(value))
    return true
  } catch {
    return false
  }
}

function removeStorage(key: string): void {
  try {
    localStorage.removeItem(key)
  } catch {
    // ignore unavailable storage
  }
}

function readSession(): AuthSession | null {
  if (memorySession) return memorySession

  const parsed = readStorage<Partial<AuthSession>>(AUTH_STORAGE_KEY)
  if (
    !parsed ||
    typeof parsed.accessToken !== 'string' ||
    !parsed.accessToken ||
    !parsed.user ||
    typeof parsed.user.username !== 'string'
  ) {
    return null
  }

  const user = parsed.user
  const session: AuthSession = {
    username:
      typeof parsed.username === 'string' && parsed.username
        ? parsed.username
        : user.username,
    name: typeof parsed.name === 'string' && parsed.name ? parsed.name : user.name,
    loggedInAt:
      typeof parsed.loggedInAt === 'string' && parsed.loggedInAt
        ? parsed.loggedInAt
        : new Date().toISOString(),
    accessToken: parsed.accessToken,
    expiresAt: typeof parsed.expiresAt === 'number' ? parsed.expiresAt : 0,
    user
  }
  memorySession = session
  return session
}

/** 管理员资料以登录会话中的用户信息为准，不再单独存 localStorage。 */
function readProfile(): AdminProfile {
  if (memoryProfile) return memoryProfile

  const user = readSession()?.user
  const profile: AdminProfile = user
    ? {
        username: user.username,
        name: user.name || user.username,
        email: user.email,
        mobile: user.mobile,
        avatar: user.avatar
      }
    : { ...DEFAULT_PROFILE }

  memoryProfile = profile
  return profile
}

/** 当前登录会话；侧边栏等组件用它实时显示管理员名称。 */
export const sessionRef = ref<AuthSession | null>(readSession())

/** 当前管理员资料；侧边栏头像与账号页共用。 */
export const profileRef = ref<AdminProfile>(readProfile())

// 把 token 读取、静默刷新与 401 处理注入请求层，避免 api.ts 反向依赖 auth.ts。
configureApiAuth({
  getAccessToken: () => sessionRef.value?.accessToken ?? null,
  refreshAccessToken: async () => {
    const session = await refreshSession()
    return session?.accessToken ?? null
  },
  onUnauthorized: () => {
    clearSession()
    if (typeof window === 'undefined') return
    const current = window.location.hash.replace(/^#/, '') || '/pages'
    if (!current.startsWith('/login')) {
      window.location.hash = `#/login?redirect=${encodeURIComponent(current)}`
    }
  }
})

function passwordByteLength(value: string): number {
  return new TextEncoder().encode(value).length
}

/** 校验新密码长度；返回中文错误提示，通过时返回 null。 */
export function validateNewPassword(value: string): string | null {
  if (!value) return '请输入新密码'
  if (Array.from(value).length < MIN_PASSWORD_CHARS) {
    return `新密码至少 ${MIN_PASSWORD_CHARS} 个字符`
  }
  if (passwordByteLength(value) > MAX_PASSWORD_BYTES) {
    return `新密码不能超过 ${MAX_PASSWORD_BYTES} 字节（中文约 24 个字符）`
  }
  return null
}

/** 登录：调用 Go API，成功后保存 token 与用户信息。 */
export async function login(username: string, password: string): Promise<AuthSession> {
  const account = username.trim()
  if (!account || !password) {
    throw new Error('请输入用户名和密码')
  }

  const data = await api.post<AdminLoginResponse>(
    '/admin/auth/login',
    { username: account, password },
    { auth: false }
  )
  if (!data?.access_token || !data.user) {
    throw new Error('登录响应格式不正确')
  }

  const user = normalizeUser(data.user, account)
  const expiresIn = Number(data.expires_in)
  const session: AuthSession = {
    username: user.username,
    name: user.name,
    loggedInAt: new Date().toISOString(),
    accessToken: data.access_token,
    expiresAt: Math.floor(Date.now() / 1000) + (expiresIn > 0 ? expiresIn : 900),
    user
  }

  memorySession = session
  sessionRef.value = session
  writeStorage(AUTH_STORAGE_KEY, session)
  seedProfileFromUser(user)
  return session
}

/** 统一补齐登录 / 刷新响应里的用户字段，避免后端缺字段导致界面空白。 */
function normalizeUser(raw: AdminUser, fallbackUsername: string): AdminUser {
  return {
    ...raw,
    id: String(raw.id ?? ''),
    username: raw.username || fallbackUsername,
    name: raw.name || raw.username || fallbackUsername,
    email: raw.email || '',
    mobile: raw.mobile || '',
    avatar: raw.avatar || '',
    roles: Array.isArray(raw.roles) ? raw.roles : [],
    permissions: Array.isArray(raw.permissions) ? raw.permissions : []
  }
}

function seedProfileFromUser(user: AdminUser): void {
  const profile: AdminProfile = {
    username: user.username,
    name: user.name,
    email: user.email,
    mobile: user.mobile,
    avatar: user.avatar
  }
  memoryProfile = profile
  profileRef.value = { ...profile }
}

export function getSession(): AuthSession | null {
  return readSession()
}

/** 本地 Access Token 是否仍在有效期内；过期后由 refreshSession 尝试静默续期。 */
export function isAuthenticated(): boolean {
  const session = getSession()
  if (!session) return false
  return session.expiresAt > Math.floor(Date.now() / 1000)
}

let refreshPromise: Promise<AuthSession | null> | null = null

/**
 * 使用 HttpOnly Cookie 刷新 Access Token。
 * 并发调用共享同一次请求；刷新令牌失效时清除本地会话并返回 null。
 */
export function refreshSession(): Promise<AuthSession | null> {
  if (refreshPromise) return refreshPromise

  refreshPromise = performRefresh().finally(() => {
    refreshPromise = null
  })
  return refreshPromise
}

async function performRefresh(): Promise<AuthSession | null> {
  try {
    const data = await api.post<AdminLoginResponse>('/admin/auth/refresh', undefined, {
      auth: false,
      timeoutMs: 10_000
    })
    if (!data?.access_token || !data.user) return null

    const user = normalizeUser(data.user, data.user.username)

    const expiresIn = Number(data.expires_in)
    const session: AuthSession = {
      username: user.username,
      name: user.name,
      loggedInAt: sessionRef.value?.loggedInAt ?? new Date().toISOString(),
      accessToken: data.access_token,
      expiresAt: Math.floor(Date.now() / 1000) + (expiresIn > 0 ? expiresIn : 900),
      user
    }

    memorySession = session
    sessionRef.value = session
    writeStorage(AUTH_STORAGE_KEY, session)
    return session
  } catch (error) {
    // 只有服务端明确判定刷新令牌失效时才清除会话；网络 / 服务端故障保留会话，稍后可重试。
    if (error instanceof ApiError && (error.code === 40100 || error.code === 40101 || error.code === 401)) {
      clearSession()
    }
    return null
  }
}

/** 清除本地登录会话，不发起网络请求。 */
function clearSession(): void {
  memorySession = null
  sessionRef.value = null
  removeStorage(AUTH_STORAGE_KEY)
}

/**
 * 退出登录：先让服务端吊销 HttpOnly Cookie 中的刷新令牌，再清除本地会话。
 * 网络或服务端异常时仍然清除本地状态，避免用户卡在已登录界面。
 */
export async function logout(): Promise<void> {
  try {
    await api.post('/admin/auth/logout', undefined, { auth: false, timeoutMs: 5_000 })
  } catch {
    // 服务端不可用不影响本地退出；Cookie 中的令牌会自然过期。
  } finally {
    clearSession()
  }
}

export function getProfile(): AdminProfile {
  return { ...readProfile() }
}

/**
 * 保存显示名称、邮箱、手机号和头像到 Go 后端，并同步右上角用户区。
 * 登录账号固定不可修改。
 */
export async function saveProfile(
  input: Pick<AdminProfile, 'name' | 'email' | 'mobile'> & Partial<Pick<AdminProfile, 'avatar'>>
): Promise<AdminProfile> {
  const name = input.name.trim()
  if (!name) throw new Error('显示名称不能为空')

  const saved = await updateProfile({
    name,
    email: input.email.trim(),
    mobile: input.mobile.trim(),
    avatar: (input.avatar ?? '').trim()
  })

  const next: AdminProfile = {
    username: saved.username,
    name: saved.name,
    email: saved.email,
    mobile: saved.mobile,
    avatar: saved.avatar
  }
  memoryProfile = next
  profileRef.value = { ...next }

  const session = sessionRef.value
  if (session) {
    const updated: AuthSession = {
      ...session,
      name: next.name,
      user: {
        ...session.user,
        name: next.name,
        email: next.email,
        mobile: next.mobile,
        avatar: next.avatar
      }
    }
    memorySession = updated
    sessionRef.value = updated
    writeStorage(AUTH_STORAGE_KEY, updated)
  }

  return { ...next }
}

/**
 * 修改密码：调用 Go 后端，成功后服务端会吊销该账号全部刷新令牌并清除 Cookie。
 * 本地会话同步清除，调用方提示用户重新登录。
 */
export async function changePassword(
  currentPassword: string,
  newPassword: string
): Promise<ChangePasswordResult> {
  try {
    await changeAdminPassword(currentPassword, newPassword)
    clearSession()
    return { ok: true, message: '密码已修改，请使用新密码重新登录' }
  } catch (error) {
    return {
      ok: false,
      message: error instanceof Error ? error.message : '密码修改失败，请稍后重试'
    }
  }
}

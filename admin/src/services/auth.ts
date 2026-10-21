/**
 * 管理后台登录与本地演示账号资料。
 *
 * 登录 / 刷新令牌 / 退出登录已接 Go API：
 * - POST /api/v1/admin/auth/login
 * - POST /api/v1/admin/auth/refresh（HttpOnly Cookie 静默续期）
 * - POST /api/v1/admin/auth/logout
 *
 * Access Token 保存在 localStorage；Refresh Token 由后端写入 HttpOnly Cookie，前端不可读取。
 * 管理员资料与改密接口尚未接入，仍保存在 localStorage，仅用于演示。
 */
import { ref } from 'vue'
import type { AdminLoginResponse, AdminProfile, AdminUser } from '@/types/admin'
import { ApiError, api, configureApiAuth } from '@/services/api'

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

interface StoredPassword {
  algorithm: 'PBKDF2-SHA256'
  iterations: number
  salt: string
  hash: string
  updatedAt: string
}

/** 本地演示账号资料；登录接口已改用后端账号，这里仅供账号页默认值使用。 */
export const demoCredentials = {
  username: 'admin',
  password: '123456',
  name: '小阙'
} as const

export const AUTH_STORAGE_KEY = 'xiaoque-admin-auth-v1'
export const PROFILE_STORAGE_KEY = 'xiaoque-admin-profile-v1'
export const PASSWORD_STORAGE_KEY = 'xiaoque-admin-password-v1'

const PASSWORD_ITERATIONS = 100_000
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
let memoryPassword: StoredPassword | null = null

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

function readProfile(): AdminProfile {
  if (memoryProfile) return memoryProfile

  const parsed = readStorage<Partial<AdminProfile>>(PROFILE_STORAGE_KEY)
  const storedName = typeof parsed?.name === 'string' ? parsed.name.trim() : ''
  const profile: AdminProfile = {
    username: demoCredentials.username,
    name: storedName || DEFAULT_PROFILE.name,
    email: typeof parsed?.email === 'string' ? parsed.email : '',
    mobile: typeof parsed?.mobile === 'string' ? parsed.mobile : '',
    avatar: typeof parsed?.avatar === 'string' ? parsed.avatar : ''
  }
  memoryProfile = profile
  return profile
}

function readPassword(): StoredPassword | null {
  if (memoryPassword) return memoryPassword

  const parsed = readStorage<Partial<StoredPassword>>(PASSWORD_STORAGE_KEY)
  if (
    !parsed ||
    parsed.algorithm !== 'PBKDF2-SHA256' ||
    typeof parsed.salt !== 'string' ||
    typeof parsed.hash !== 'string' ||
    typeof parsed.iterations !== 'number'
  ) {
    return null
  }

  memoryPassword = {
    algorithm: 'PBKDF2-SHA256',
    iterations: parsed.iterations,
    salt: parsed.salt,
    hash: parsed.hash,
    updatedAt: typeof parsed.updatedAt === 'string' ? parsed.updatedAt : ''
  }
  return memoryPassword
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

function assertCrypto(): SubtleCrypto {
  const cryptoApi = globalThis.crypto
  if (!cryptoApi?.subtle || !cryptoApi.getRandomValues) {
    throw new Error('当前环境不支持安全密码存储，请通过 localhost 或 HTTPS 打开管理后台')
  }
  return cryptoApi.subtle
}

function toHex(bytes: Uint8Array): string {
  return Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('')
}

function fromHex(hex: string): Uint8Array<ArrayBuffer> {
  const bytes = new Uint8Array(new ArrayBuffer(hex.length / 2))
  for (let index = 0; index < bytes.length; index += 1) {
    bytes[index] = Number.parseInt(hex.slice(index * 2, index * 2 + 2), 16)
  }
  return bytes
}

function createSalt(): string {
  assertCrypto()
  const bytes = new Uint8Array(16)
  globalThis.crypto.getRandomValues(bytes)
  return toHex(bytes)
}

async function deriveHash(
  password: string,
  salt: string,
  iterations: number
): Promise<string> {
  const subtle = assertCrypto()
  const key = await subtle.importKey(
    'raw',
    new TextEncoder().encode(password),
    'PBKDF2',
    false,
    ['deriveBits']
  )
  const bits = await subtle.deriveBits(
    { name: 'PBKDF2', salt: fromHex(salt), iterations, hash: 'SHA-256' },
    key,
    256
  )
  return toHex(new Uint8Array(bits))
}

function safeEqual(left: string, right: string): boolean {
  if (left.length !== right.length) return false
  let difference = 0
  for (let index = 0; index < left.length; index += 1) {
    difference |= left.charCodeAt(index) ^ right.charCodeAt(index)
  }
  return difference === 0
}

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
  writeStorage(PROFILE_STORAGE_KEY, profile)
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
    // 账号资料目前仍是本地演示数据，刷新时保留用户已保存的显示名称 / 邮箱 / 手机 / 头像。
    const profile = readProfile()
    const mergedUser: AdminUser = {
      ...user,
      name: profile.name || user.name,
      email: profile.email,
      mobile: profile.mobile,
      avatar: profile.avatar
    }

    const expiresIn = Number(data.expires_in)
    const session: AuthSession = {
      username: mergedUser.username,
      name: mergedUser.name,
      loggedInAt: sessionRef.value?.loggedInAt ?? new Date().toISOString(),
      accessToken: data.access_token,
      expiresAt: Math.floor(Date.now() / 1000) + (expiresIn > 0 ? expiresIn : 900),
      user: mergedUser
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

/** 保存显示名称、邮箱、手机号和头像；登录账号固定不可修改。 */
export function saveProfile(
  input: Pick<AdminProfile, 'name' | 'email' | 'mobile'> & Partial<Pick<AdminProfile, 'avatar'>>
): AdminProfile {
  const current = readProfile()
  const next: AdminProfile = {
    username: current.username,
    name: input.name.trim(),
    email: input.email.trim(),
    mobile: input.mobile.trim(),
    avatar: (input.avatar ?? '').trim()
  }
  if (!next.name) throw new Error('显示名称不能为空')
  if (!writeStorage(PROFILE_STORAGE_KEY, next)) {
    throw new Error('浏览器本地存储不可用，资料未保存')
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

/** 本地演示改密；尚未接后端改密接口，不影响 xqyun_users。 */
export async function changePassword(
  currentPassword: string,
  newPassword: string
): Promise<ChangePasswordResult> {
  try {
    const stored = readPassword()
    if (stored) {
      const candidate = await deriveHash(currentPassword, stored.salt, stored.iterations)
      if (!safeEqual(candidate, stored.hash)) {
        return { ok: false, message: '当前密码不正确' }
      }
    } else if (currentPassword !== demoCredentials.password) {
      return { ok: false, message: '当前密码不正确' }
    }

    if (currentPassword === newPassword) {
      return { ok: false, message: '新密码不能与当前密码相同' }
    }

    const validation = validateNewPassword(newPassword)
    if (validation) return { ok: false, message: validation }

    const salt = createSalt()
    const hash = await deriveHash(newPassword, salt, PASSWORD_ITERATIONS)
    const record: StoredPassword = {
      algorithm: 'PBKDF2-SHA256',
      iterations: PASSWORD_ITERATIONS,
      salt,
      hash,
      updatedAt: new Date().toISOString()
    }

    if (!writeStorage(PASSWORD_STORAGE_KEY, record)) {
      return { ok: false, message: '浏览器本地存储不可用，密码未保存' }
    }

    memoryPassword = record
    return { ok: true, message: '密码已修改' }
  } catch (error) {
    return {
      ok: false,
      message: error instanceof Error ? error.message : '密码修改失败，请稍后重试'
    }
  }
}

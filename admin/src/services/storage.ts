import type { PersistedDemo } from '@/types/admin'

export const STORAGE_KEY = 'xiaoque-admin-demo-v1'

export function loadDemo(): PersistedDemo | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as PersistedDemo
    if (parsed?.schemaVersion !== 1 || !Array.isArray(parsed.pages)) {
      return null
    }
    return parsed
  } catch {
    return null
  }
}

export function saveDemo(demo: PersistedDemo): boolean {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(demo))
    return true
  } catch {
    return false
  }
}

export function clearDemo(): void {
  try {
    localStorage.removeItem(STORAGE_KEY)
  } catch {
    // ignore unavailable storage
  }
}

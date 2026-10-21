import type { PageRecord, SiteSettings } from '~/domain/types'
import type { PersistedDemo } from '@/types/admin'
import zhPages from './pages.zh'
import enPages from './pages.en'
import { contents } from './content'
import { settings } from './settings'
import { submissions } from './submissions'
import { loginSettings } from './loginSettings'

export const seedPages: PageRecord[] = [...zhPages, ...enPages]
export const seedContents = contents
export const seedSettings: SiteSettings = settings
export const seedSubmissions = submissions
export const seedLoginSettings = loginSettings

export function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

export function createSeedPersisted(): PersistedDemo {
  return {
    schemaVersion: 1,
    savedAt: new Date().toISOString(),
    pages: clone(seedPages),
    submissions: clone(seedSubmissions),
    loginSettings: clone(seedLoginSettings)
  }
}

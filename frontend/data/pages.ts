import type { PageRecord } from '~/domain/types'
import zhPages from './pages.zh'
import enPages from './pages.en'

export const pages: PageRecord[] = [...zhPages, ...enPages]

import type { PageBlock, PageRecord, SeoMeta } from '~/domain/types'
import type { Locale } from '~/domain/types'
import type { PageStatus } from '@/types/admin'
import { blockSchemas, clone, createBlock, uid } from '@/domain/blockSchema'
import type { BlockField } from '@/domain/blockSchema'

const RESERVED = new Set(['admin', 'api', 'preview', 'uploads', 'assets', 'en'])
const KEEP_TEXT_KEYS = new Set([
  'href',
  'url',
  'videoUrl',
  'primaryHref',
  'secondaryHref',
  'path',
  'key',
  'email',
  'phone',
  'slug',
  'value',
  'suffix',
  'metric',
  'metricSuffix',
  'icon',
  'variant',
  'theme',
  'layout',
  'tone',
  'visual',
  'cover',
  'mapStyle',
  'imageSide',
  'width',
  'image',
  'images',
  'limit',
  'interval',
  'autoplay',
  'showFilters',
  'marquee'
])

export function nowStamp(): string {
  const date = new Date()
  const pad = (value: number) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(
    date.getHours()
  )}:${pad(date.getMinutes())}`
}

export function formatClock(value: Date = new Date()): string {
  return new Intl.DateTimeFormat('zh-CN', { hour: '2-digit', minute: '2-digit' }).format(value)
}

export function normalizePath(path: string): string {
  let normalized = (path || '').trim()
  if (!normalized.startsWith('/')) normalized = `/${normalized}`
  if (normalized.length > 1) normalized = normalized.replace(/\/+$/, '')
  return normalized || '/'
}

export function isReservedPath(path: string): boolean {
  const firstSegment = normalizePath(path).split('/').filter(Boolean)[0]
  return firstSegment ? RESERVED.has(firstSegment) : false
}

export function isPathTaken(
  pages: PageRecord[],
  locale: Locale,
  path: string,
  excludeId?: string
): boolean {
  const normalized = normalizePath(path)
  return pages.some(
    (page) =>
      page.locale === locale &&
      page.id !== excludeId &&
      normalizePath(page.path) === normalized
  )
}

export function pageStatus(page: PageRecord): PageStatus {
  return page.status === 'published' ? 'published' : 'draft'
}

export function createBlankPage(name: string, path: string, locale: Locale): PageRecord {
  const pageId = uid('pg')
  const hero = createBlock('hero') as PageBlock
  return {
    id: pageId,
    groupId: uid('tg'),
    locale,
    path: normalizePath(path),
    name,
    status: 'draft',
    updatedAt: nowStamp(),
    seo: { title: name, description: '', keywords: [] },
    blocks: [{ ...hero, id: uid('blk') }],
    published: null
  }
}

export function duplicatePageRecord(page: PageRecord, pages: PageRecord[] = []): PageRecord {
  const basePath = normalizePath(page.path)
  const base = basePath === '/' ? '/home' : basePath
  let candidate = `${base}-copy`
  let suffix = 2
  while (
    pages.some(
      (item) => item.locale === page.locale && normalizePath(item.path) === candidate
    )
  ) {
    candidate = `${base}-copy-${suffix}`
    suffix += 1
  }
  return {
    ...clone(page),
    id: uid('pg'),
    groupId: uid('tg'),
    path: candidate,
    name: `${page.name} 副本`,
    status: 'draft',
    updatedAt: nowStamp(),
    published: null,
    blocks: clone(page.blocks).map((block) => ({ ...block, id: uid('blk') }) as PageBlock)
  }
}

function keepTextValue(fieldKey: string): boolean {
  return KEEP_TEXT_KEYS.has(fieldKey)
}

function clearByDefinition(definition: BlockField, value: unknown): unknown {
  if (definition.type === 'repeater') {
    if (!Array.isArray(value)) return value
    return value.map((item) => {
      if (!item || typeof item !== 'object') return item
      const source = { ...(item as Record<string, unknown>) }
      for (const itemField of definition.itemFields ?? []) {
        source[itemField.key] = clearByDefinition(itemField, source[itemField.key])
      }
      return source
    })
  }
  if (definition.type === 'group') {
    return value
  }
  if (
    (definition.type === 'text' ||
      definition.type === 'textarea' ||
      definition.type === 'richtext' ||
      definition.type === 'html') &&
    !keepTextValue(definition.key)
  ) {
    return ''
  }
  return value
}

export function clearTranslatableBlocks(blocks: PageBlock[]): PageBlock[] {
  return blocks.map((block) => {
    const schema = blockSchemas[block.type]
    const props = clone(block.props) as unknown as Record<string, unknown>
    if (!schema) return { ...clone(block), props } as unknown as PageBlock
    for (const field of schema.fields) {
      props[field.key] = clearByDefinition(field, props[field.key])
    }
    return { ...clone(block), id: uid('blk'), props } as unknown as PageBlock
  })
}

/** 旧版「自定义 HTML」区块并入富文本区块。 */
export function migrateLegacyBlocks(blocks: PageBlock[]): PageBlock[] {
  return blocks.map((block) => {
    if (block.type !== 'html') return block
    const legacy = block.props as { code?: unknown }
    return {
      id: block.id,
      type: 'rich-text',
      props: {
        content: typeof legacy.code === 'string' ? legacy.code : '',
        width: 'narrow'
      }
    } as PageBlock
  })
}

export function emptySeo(): SeoMeta {
  return { title: '', description: '', keywords: [] }
}

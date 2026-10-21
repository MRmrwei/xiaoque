import type { BlockInput, BlockPropsMap, BlockType, PageBlock } from '~/domain/blocks'
import type { Locale, PageRecord, SeoMeta } from '~/domain/types'

export interface MakePageInput {
  id: string
  group: string
  locale: Locale
  path: string
  name: string
  seo: SeoMeta
  blocks: BlockInput[]
  version?: number
  status?: 'published' | 'draft'
  updatedAt?: string
  publishedAt?: string
}

export function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

export const b = <K extends BlockType>(
  type: K,
  props: BlockPropsMap[K]
): BlockInput => ({ type, props }) as BlockInput

export function toBlocks(pageId: string, defs: BlockInput[]): PageBlock[] {
  return defs.map(
    (def, index) =>
      ({
        id: `${pageId}_b${index + 1}`,
        type: def.type,
        props: clone(def.props)
      }) as PageBlock
  )
}

export function makePage({
  id,
  group,
  locale,
  path,
  name,
  seo,
  blocks,
  version = 1,
  status = 'published',
  updatedAt = '2026-10-12T10:00:00+08:00',
  publishedAt = updatedAt
}: MakePageInput): PageRecord {
  const blockList = toBlocks(id, blocks)
  const page: PageRecord = {
    id,
    groupId: group,
    locale,
    path,
    name,
    status,
    updatedAt,
    seo: clone(seo),
    blocks: blockList,
    published: null
  }

  if (status === 'published') {
    page.published = {
      version,
      blocks: clone(blockList),
      seo: clone(seo),
      publishedAt,
      publishedBy: '林小鹊'
    }
  }

  return page
}

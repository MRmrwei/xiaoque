import { computed, type Ref } from 'vue'
import type { PageRecord } from '~/domain/types'
import { blockSchemas } from '@/domain/blockSchema'
import { isPathTaken, isReservedPath, normalizePath, pageStatus } from '@/services/pageService'

export interface ValidationCheck {
  key: string
  label: string
  ok: boolean
  warn?: boolean
}

function isEmpty(value: unknown): boolean {
  if (value === null || value === undefined) return true
  if (typeof value === 'string') return value.trim().length === 0
  if (Array.isArray(value)) return value.length === 0
  return false
}

function requiredBlockIssues(page: PageRecord): string[] {
  const issues: string[] = []
  for (const block of page.blocks) {
    const schema = blockSchemas[block.type]
    if (!schema) continue
    const props = block.props as unknown as Record<string, unknown>
    for (const field of schema.fields) {
      const value = props[field.key]
      if (field.type === 'repeater' && Array.isArray(value)) {
        value.forEach((item, index) => {
          const record = item as Record<string, unknown>
          for (const itemField of field.itemFields ?? []) {
            if (itemField.required && isEmpty(record?.[itemField.key])) {
              issues.push(`${schema.label} · 第 ${index + 1} 项「${itemField.label}」`)
            }
          }
        })
        continue
      }
      if (field.required && isEmpty(value)) {
        issues.push(`${schema.label} · ${field.label}`)
      }
    }
  }
  return issues
}

export function usePublishValidation(
  page: Readonly<Ref<PageRecord | null>>,
  pages: Readonly<Ref<PageRecord[]>>,
  pair: Readonly<Ref<PageRecord | undefined>>
) {
  const checks = computed<ValidationCheck[]>(() => {
    const current = page.value
    if (!current) return []

    const path = normalizePath(current.path)
    const requiredIssues = requiredBlockIssues(current)
    const pairNeedsPublish = pair.value && pageStatus(pair.value) === 'draft'

    return [
      {
        key: 'name',
        label: '页面名称已填写',
        ok: Boolean(current.name.trim())
      },
      {
        key: 'path',
        label: '访问路径有效且未被占用',
        ok:
          path.startsWith('/') &&
          !isReservedPath(path) &&
          !isPathTaken(pages.value, current.locale, path, current.id)
      },
      {
        key: 'blocks',
        label: '页面至少包含一个区块',
        ok: current.blocks.length > 0
      },
      {
        key: 'required',
        label: requiredIssues.length
          ? `区块必填项：${requiredIssues.slice(0, 2).join('；')}${requiredIssues.length > 2 ? '…' : ''}`
          : '区块必填项已完成',
        ok: requiredIssues.length === 0
      },
      {
        key: 'seo-title',
        label: '已配置 SEO 标题',
        ok: Boolean(current.seo.title.trim())
      },
      {
        key: 'seo-description',
        label: '建议补充 SEO 描述',
        ok: Boolean(current.seo.description.trim()),
        warn: true
      },
      {
        key: 'pair',
        label: pairNeedsPublish ? '另一种语言尚未发布' : '双语版本状态正常',
        ok: !pairNeedsPublish,
        warn: Boolean(pairNeedsPublish)
      }
    ]
  })

  const ready = computed(() => checks.value.every((check) => check.ok || check.warn))
  const errorCount = computed(
    () => checks.value.filter((check) => !check.ok && !check.warn).length
  )

  return { checks, ready, errorCount }
}

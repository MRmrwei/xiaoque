import type { Locale } from '~/domain/types'

function escapeHtml(value = ''): string {
  return String(value)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

export function titleHtml(value = ''): string {
  return escapeHtml(value)
    .replace(/&lt;em&gt;/gi, '<em>')
    .replace(/&lt;\/em&gt;/gi, '</em>')
    .replace(/&lt;br\s*\/?&gt;/gi, '<br />')
    .replace(/\n/g, '<br />')
}

export function stripTags(value = ''): string {
  return String(value).replace(/<[^>]+>/g, '')
}

export function formatDate(value: string, locale: Locale = 'zh'): string {
  if (!value) return ''
  const date = new Date(String(value).replace(' ', 'T'))
  if (Number.isNaN(date.getTime())) return String(value)
  return new Intl.DateTimeFormat(locale === 'en' ? 'en-GB' : 'zh-CN', {
    year: 'numeric',
    month: locale === 'en' ? 'short' : 'long',
    day: 'numeric'
  }).format(date)
}

export function formatDateTime(value: string): string {
  if (!value) return ''
  return String(value).replace('T', ' ').slice(0, 16)
}

export function copyText(text: string): Promise<void> {
  if (navigator.clipboard?.writeText) {
    return navigator.clipboard.writeText(text).catch(() => fallbackCopy(text))
  }
  fallbackCopy(text)
  return Promise.resolve()
}

function fallbackCopy(text: string): void {
  try {
    const area = document.createElement('textarea')
    area.value = text
    area.style.position = 'fixed'
    area.style.opacity = '0'
    document.body.appendChild(area)
    area.select()
    document.execCommand('copy')
    document.body.removeChild(area)
  } catch {
    // Clipboard access can be unavailable; the toast still confirms the intent.
  }
}

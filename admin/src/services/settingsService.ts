import type {
  FooterCta,
  FooterSettings,
  HeaderAnnouncement,
  HeaderCta,
  HeaderCtaVariant,
  HeaderSettings,
  Locale,
  SiteSettings
} from '~/domain/types'
import { clone } from '@/domain/blockSchema'
import { settings as seedSettings } from '@/data/seed/settings'

type LegacyFooterSettings = Partial<FooterSettings> & {
  cta?: Partial<FooterCta>
}

/** 兼容旧版单个顶部按钮（showCta / ctaLabel / ctaHref）。 */
type LegacyHeaderSettings = Partial<Omit<HeaderSettings, 'announcement'>> & {
  announcement?: Partial<HeaderAnnouncement>
  showCta?: boolean
  ctaLabel?: Record<Locale, string>
  ctaHref?: string
}

/** 旧版本地数据或 Go API 返回的部分设置，允许缺失新字段。 */
export type SiteSettingsInput = Omit<Partial<SiteSettings>, 'header' | 'footer'> & {
  header?: LegacyHeaderSettings
  footer?: LegacyFooterSettings
}

function normalizeHeaderCtas(value: unknown, defaults: HeaderCta[]): HeaderCta[] {
  if (!Array.isArray(value)) return clone(defaults)
  return value
    .filter((item): item is Partial<HeaderCta> => Boolean(item) && typeof item === 'object')
    .map((item, index) => {
      const variant: HeaderCtaVariant =
        item.variant === 'secondary' || item.variant === 'ghost' ? item.variant : 'primary'
      return {
        key: item.key || `header.cta.${index + 1}`,
        label: { zh: item.label?.zh ?? '', en: item.label?.en ?? '' },
        href: item.href ?? '',
        variant,
        visible: item.visible !== false
      }
    })
}

function resolveHeaderCtas(header: LegacyHeaderSettings, defaults: HeaderCta[]): HeaderCta[] {
  if (Array.isArray(header.ctas)) return normalizeHeaderCtas(header.ctas, defaults)
  // 旧数据：showCta + ctaLabel + ctaHref → 单个主按钮。
  if (header.showCta === false) return []
  if (header.showCta === true || header.ctaLabel || header.ctaHref) {
    return [
      {
        key: 'header.cta.demo',
        label: {
          zh: header.ctaLabel?.zh || defaults[0]?.label.zh || '预约演示',
          en: header.ctaLabel?.en || defaults[0]?.label.en || 'Book a demo'
        },
        href: header.ctaHref?.trim() || defaults[0]?.href || '/contact',
        variant: 'primary',
        visible: true
      }
    ]
  }
  return clone(defaults)
}

/**
 * Merge persisted demo settings with the current seed shape.
 * Older localStorage payloads may miss header/footer or newly added fields.
 */
export function normalizeSiteSettings(raw?: SiteSettingsInput | null): SiteSettings {
  const base = clone(seedSettings)
  if (!raw) return base

  const merged = { ...base, ...raw } as SiteSettings
  merged.nav = Array.isArray(raw.nav) ? raw.nav : base.nav
  merged.brand = { ...base.brand, ...(raw.brand ?? {}) }

  const header: LegacyHeaderSettings = raw.header ?? {}
  merged.header = {
    showLanguage: header.showLanguage ?? base.header.showLanguage,
    sticky: header.sticky ?? base.header.sticky,
    showAnnouncement: header.showAnnouncement ?? base.header.showAnnouncement,
    announcement: {
      ...base.header.announcement,
      ...(header.announcement ?? {}),
      text: { ...base.header.announcement.text, ...(header.announcement?.text ?? {}) },
      highlight: {
        ...base.header.announcement.highlight,
        ...(header.announcement?.highlight ?? {})
      },
      linkLabel: {
        ...base.header.announcement.linkLabel,
        ...(header.announcement?.linkLabel ?? {})
      }
    },
    ctas: resolveHeaderCtas(header, base.header.ctas)
  }

  const footer: LegacyFooterSettings = raw.footer ?? {}
  merged.footer = {
    ...base.footer,
    ...footer,
    cta: {
      ...base.footer.cta,
      ...(footer.cta ?? {}),
      label: { ...base.footer.cta.label, ...(footer.cta?.label ?? {}) }
    },
    columns: Array.isArray(footer.columns) ? footer.columns : base.footer.columns,
    bottomLinks: Array.isArray(footer.bottomLinks) ? footer.bottomLinks : base.footer.bottomLinks
  }

  merged.siteName = { ...base.siteName, ...(raw.siteName ?? {}) }
  merged.companyName = { ...base.companyName, ...(raw.companyName ?? {}) }
  merged.tagline = { ...base.tagline, ...(raw.tagline ?? {}) }
  merged.address = { ...base.address, ...(raw.address ?? {}) }
  merged.workHours = { ...base.workHours, ...(raw.workHours ?? {}) }
  merged.copyright = { ...base.copyright, ...(raw.copyright ?? {}) }
  merged.socials = Array.isArray(raw.socials) ? raw.socials : base.socials

  return merged
}

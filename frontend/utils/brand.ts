const DEFAULT_ACCENT = '#E8503A'
const DEFAULT_BACKGROUND = '#F6F2EA'
const DEFAULT_DARK = '#141518'

type Rgb = [number, number, number]

export interface BrandPalette {
  accent?: string
  background?: string
  darkBackground?: string
}

function parseHex(value: string): Rgb | null {
  let hex = value.trim().replace(/^#/, '')
  if (hex.length === 3) {
    hex = hex
      .split('')
      .map((char) => char + char)
      .join('')
  }
  if (!/^[0-9a-f]{6}$/i.test(hex)) return null
  return [
    Number.parseInt(hex.slice(0, 2), 16),
    Number.parseInt(hex.slice(2, 4), 16),
    Number.parseInt(hex.slice(4, 6), 16)
  ]
}

function toHex(rgb: Rgb): string {
  return `#${rgb.map((value) => value.toString(16).padStart(2, '0')).join('')}`
}

function toRgba(rgb: Rgb, alpha: number): string {
  return `rgba(${rgb[0]}, ${rgb[1]}, ${rgb[2]}, ${alpha})`
}

function mix(rgb: Rgb, target: Rgb, ratio: number): Rgb {
  return rgb.map((value, index) =>
    Math.round(value + (target[index] - value) * ratio)
  ) as Rgb
}

export function normalizeColor(value?: string | null, fallback = DEFAULT_ACCENT): string {
  const rgb = value ? parseHex(value) : null
  return rgb ? toHex(rgb).toUpperCase() : fallback
}

/** 把品牌配置换算成官网使用的一整套设计变量。 */
export function brandStyle(brand?: BrandPalette | null): Record<string, string> {
  const accent = normalizeColor(brand?.accent, DEFAULT_ACCENT)
  const paper = normalizeColor(brand?.background, DEFAULT_BACKGROUND)
  const ink = normalizeColor(brand?.darkBackground, DEFAULT_DARK)

  const accentRgb = parseHex(accent) ?? parseHex(DEFAULT_ACCENT)!
  const paperRgb = parseHex(paper) ?? parseHex(DEFAULT_BACKGROUND)!
  const inkRgb = parseHex(ink) ?? parseHex(DEFAULT_DARK)!

  return {
    '--accent': accent,
    '--accent-deep': toHex(mix(accentRgb, [0, 0, 0], 0.18)),
    '--accent-soft': toHex(mix(accentRgb, [255, 255, 255], 0.86)),
    '--paper': paper,
    '--paper-deep': toHex(mix(paperRgb, inkRgb, 0.06)),
    '--paper-light': toHex(mix(paperRgb, [255, 255, 255], 0.55)),
    '--ink': ink,
    '--ink-soft': toHex(mix(inkRgb, [255, 255, 255], 0.08)),
    '--muted': toHex(mix(inkRgb, paperRgb, 0.41)),
    '--muted-2': toHex(mix(inkRgb, paperRgb, 0.61)),
    '--line': toRgba(inkRgb, 0.12),
    '--line-soft': toRgba(inkRgb, 0.07),
    '--line-strong': toRgba(inkRgb, 0.22)
  }
}

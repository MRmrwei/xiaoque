/** 上传图片的统一处理：校验、等比压缩、转 data URL。 */

export const MAX_SOURCE_BYTES = 12 * 1024 * 1024
export const MAX_EDGE = 1600
const COMPRESS_QUALITY = 0.84

/** 内置图形是用 key 表示的，其余视为自定义图片地址。 */
export function isCustomImage(value: string | undefined | null): boolean {
  return /^(data:image\/|https?:\/\/|\/)/i.test((value ?? '').trim())
}

/** data URL 的实际字节数（base64 约为长度的 3/4）。 */
export function estimateDataUrlBytes(value: string): number {
  return Math.round((value.length * 3) / 4)
}

export function formatBytes(bytes: number): string {
  if (bytes >= 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(1)} MB`
  return `${Math.max(1, Math.round(bytes / 1024))} KB`
}

export async function fileToImageDataUrl(
  file: File,
  maxEdge: number = MAX_EDGE
): Promise<string> {
  if (!file.type.startsWith('image/')) {
    throw new Error('请选择图片文件')
  }
  if (file.size > MAX_SOURCE_BYTES) {
    throw new Error('图片不能超过 12MB')
  }

  // SVG 与 GIF 走原文件，保留矢量 / 动画。
  if (file.type === 'image/svg+xml' || file.type === 'image/gif') {
    if (file.size > 2 * 1024 * 1024) {
      throw new Error('SVG / GIF 请控制在 2MB 以内')
    }
    return readAsDataUrl(file)
  }

  const image = await loadImage(file)
  try {
    const edge = Math.max(1, Math.floor(maxEdge))
    const scale = Math.min(1, edge / Math.max(image.naturalWidth, image.naturalHeight))
    const width = Math.max(1, Math.round(image.naturalWidth * scale))
    const height = Math.max(1, Math.round(image.naturalHeight * scale))

    const canvas = document.createElement('canvas')
    canvas.width = width
    canvas.height = height
    const context = canvas.getContext('2d')
    if (!context) throw new Error('当前浏览器无法处理图片')
    context.drawImage(image, 0, 0, width, height)

    // WebP 体积最小；不支持时退回白底 JPEG。
    const webp = canvas.toDataURL('image/webp', COMPRESS_QUALITY)
    if (webp.startsWith('data:image/webp')) return webp

    const fallback = document.createElement('canvas')
    fallback.width = width
    fallback.height = height
    const fallbackContext = fallback.getContext('2d')
    if (!fallbackContext) throw new Error('当前浏览器无法处理图片')
    fallbackContext.fillStyle = '#ffffff'
    fallbackContext.fillRect(0, 0, width, height)
    fallbackContext.drawImage(image, 0, 0, width, height)
    return fallback.toDataURL('image/jpeg', COMPRESS_QUALITY)
  } finally {
    URL.revokeObjectURL(image.src)
  }
}

function readAsDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result ?? ''))
    reader.onerror = () => reject(new Error('图片读取失败'))
    reader.readAsDataURL(file)
  })
}

function loadImage(file: File): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const image = new Image()
    image.decoding = 'async'
    image.onload = () => resolve(image)
    image.onerror = () => {
      URL.revokeObjectURL(image.src)
      reject(new Error('图片读取失败，请换一张试试'))
    }
    image.src = URL.createObjectURL(file)
  })
}

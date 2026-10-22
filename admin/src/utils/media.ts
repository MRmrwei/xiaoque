import type { MediaKind } from '@/services/mediaApi'

export interface MediaRule {
  pattern: RegExp
  maxBytes: number
  accept: string
  label: string
  hint: string
}

/** 各素材类型的前端校验与上传提示；后端会再做一次校验。 */
export const MEDIA_RULES: Record<MediaKind, MediaRule> = {
  video: {
    pattern: /\.(mp4|webm|ogv|ogg|mov|m4v)$/i,
    maxBytes: 200 * 1024 * 1024,
    accept: '.mp4,.webm,.ogv,.ogg,.mov,.m4v,video/mp4,video/webm,video/ogg,video/quicktime',
    label: '视频',
    hint: '支持 MP4 / WebM / MOV，单个不超过 200MB'
  },
  image: {
    pattern: /\.(jpg|jpeg|png|webp|gif|avif)$/i,
    maxBytes: 10 * 1024 * 1024,
    accept:
      '.jpg,.jpeg,.png,.webp,.gif,.avif,image/jpeg,image/png,image/webp,image/gif,image/avif',
    label: '图片',
    hint: '支持 JPG / PNG / WebP / GIF / AVIF，单张不超过 10MB'
  }
}

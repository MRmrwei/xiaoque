/** 上传的图片、绝对地址或以 / 开头的站内资源都按真实图片渲染。 */
export function isCustomImage(value: string | undefined | null): boolean {
  return /^(data:image\/|https?:\/\/|\/)/i.test((value ?? '').trim())
}

/**
 * 登录页品牌区底图样式。
 * 未设置底图时返回空对象，由 CSS 的黑色背景兜底；设置后叠加一层深色遮罩，保证文案可读。
 */
export function loginBrandBackground(image: string | undefined): Record<string, string> {
  const value = image?.trim()
  if (!value) return {}

  return {
    backgroundImage: [
      'linear-gradient(rgba(10, 11, 13, 0.56), rgba(10, 11, 13, 0.74))',
      `url("${value}")`
    ].join(', '),
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    backgroundRepeat: 'no-repeat'
  }
}

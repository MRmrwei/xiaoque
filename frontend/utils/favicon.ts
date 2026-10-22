/** 未上传官网 Logo 时，浏览器标签页使用的内置图标。 */
export const DEFAULT_FAVICON =
  'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgdmlld0JveD0iMCAwIDQwIDQwIj4KICA8cmVjdCB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHJ4PSIxMiIgZmlsbD0iIzE0MTUxOCIvPgogIDxwYXRoIGQ9Ik0xMS4yIDI2LjhjLTEuOS0xLjItMi44LTMuMS0yLjQtNS40LjMtMS45IDEuNC0zLjYgMy00LjcgMS44LTEuMyA0LTEuOCA2LjItMS4zIDEuOS40IDMuNCAxLjYgNC4zIDMuMi44IDEuNCAxIDMuMS41IDQuOC0uNiAyLjEtMi4yIDMuNy00LjMgNC40LTIuNC44LTUgLjYtNy4zLS42WiIgZmlsbD0iI0Y2RjJFQSIvPgogIDxjaXJjbGUgY3g9IjI1LjciIGN5PSIxNS4zIiByPSI1LjEiIGZpbGw9IiNGNkYyRUEiLz4KICA8cGF0aCBkPSJNMzAuNSAxNC4yIDM2IDE2bC01LjIgMi4zLjYtMi4xLS45LTJaIiBmaWxsPSIjRTg1MDNBIi8+CiAgPHBhdGggZD0iTTEyLjMgMjUuNmMxLjQtNCA1LTYuNiA5LjMtNi42IDIuNiAwIDQuOCAxIDYuNCAyLjctMy4zLTEtNi44IDAtOS4xIDIuNS0uOSAxLTEuNSAyLjMtMS43IDMuNi0yLjItLjEtNC0uOS00LjktMi4yWiIgZmlsbD0iI0U4NTAzQSIvPgogIDxwYXRoIGQ9Ik0xMS42IDI1LjggNS4yIDMzLjJsMS41LTguNWMxLjQuOCAzLjEgMS4yIDQuOSAxLjFaIiBmaWxsPSIjRjZGMkVBIi8+CiAgPGNpcmNsZSBjeD0iMjcuMiIgY3k9IjE0LjMiIHI9IjEuMDUiIGZpbGw9IiMxNDE1MTgiLz4KPC9zdmc+Cg=='

/** 根据 data URL 前缀推断图标 MIME；远程地址返回 undefined，交给浏览器判断。 */
export function faviconType(href: string): string | undefined {
  if (href.startsWith('data:image/svg+xml')) return 'image/svg+xml'
  if (href.startsWith('data:image/png')) return 'image/png'
  if (href.startsWith('data:image/webp')) return 'image/webp'
  if (href.startsWith('data:image/jpeg')) return 'image/jpeg'
  if (href.startsWith('data:image/x-icon')) return 'image/x-icon'
  return undefined
}

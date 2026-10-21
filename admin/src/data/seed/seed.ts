import type { SiteSettings } from '~/domain/types'
import { contents } from './content'
import { settings } from './settings'
import { loginSettings } from './loginSettings'

// 页面、导航、站点配置与表单留言已改由 Go API 提供，
// 这里只保留仍需本地演示的内容库与登录页设置。
export const seedContents = contents
export const seedSettings: SiteSettings = settings
export const seedLoginSettings = loginSettings

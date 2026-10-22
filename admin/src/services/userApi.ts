/**
 * 当前管理员资料与密码接口封装。
 *
 * 路径相对于 API_BASE_URL（默认 /api/v1），请求层自动携带 Bearer Token 并处理静默刷新。
 * 用户 ID 由后端从 Access Token 解析，前端不传。
 */

import { api } from '@/services/api'

/** PUT /api/v1/admin/user/profile 的响应 data。 */
export interface ApiAdminProfile {
  id: string
  username: string
  name: string
  email: string
  mobile: string
  avatar: string
}

export interface ApiUpdateProfileRequest {
  name: string
  email: string
  mobile: string
  avatar: string
}

export function updateProfile(body: ApiUpdateProfileRequest): Promise<ApiAdminProfile> {
  return api.put<ApiAdminProfile>('/admin/user/profile', body)
}

/** 修改当前管理员密码；成功后下次登录使用新密码。 */
export function changeAdminPassword(
  currentPassword: string,
  newPassword: string
): Promise<void> {
  return api.put<void>('/admin/user/password', {
    current_password: currentPassword,
    new_password: newPassword
  })
}

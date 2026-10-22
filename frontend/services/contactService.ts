import type { ContactInput, Locale } from '~/domain/types'

interface ApiEnvelope<T> {
  code: number
  message: string
  data?: T
}

interface ApiSubmissionResult {
  id: string
}

/**
 * 表单提交边界：调用 Go 后端公开接口，由后端校验并写入后台「表单留言」。
 * 校验与 UI 状态仍由 useContactForm() 负责。
 */
export async function submitContact(
  input: ContactInput,
  locale: Locale
): Promise<void> {
  const config = useRuntimeConfig()
  // 单文件分享版没有后端，保持本地演示成功态。
  if (config.public.useMockData) {
    await new Promise((resolve) => setTimeout(resolve, 450))
    return
  }

  const base = import.meta.server ? config.apiBase : config.public.apiBase

  let body: ApiEnvelope<ApiSubmissionResult>
  try {
    body = await $fetch<ApiEnvelope<ApiSubmissionResult>>(
      `${base}/frontend/forms/contact/submissions`,
      {
        method: 'POST',
        body: {
          name: input.name,
          email: input.email,
          company: input.company,
          phone: input.phone,
          message: input.message,
          locale: locale === 'en' ? 'en' : 'zh-CN'
        }
      }
    )
  } catch {
    throw new Error('发送失败，请检查网络后重试')
  }

  if (!body || body.code !== 0) {
    throw new Error(body?.message || '发送失败，请稍后重试')
  }
}

import { reactive, ref } from 'vue'
import type { ContactInput } from '~/domain/types'
import { useToast } from './useToast'
import { useUiText } from './useUiText'

type ContactField = 'name' | 'email' | 'phone' | 'message'

/** Local-only form state for the preview canvas (no network, no backend). */
export function useContactForm() {
  const { t } = useUiText()
  const { notify } = useToast()

  const form = reactive<ContactInput>({
    name: '',
    email: '',
    company: '',
    phone: '',
    message: ''
  })
  const errors = reactive<Partial<Record<ContactField, string>>>({})
  const sending = ref(false)
  const sent = ref(false)
  /** 与前台保持同一份状态结构；预览内不会真正提交失败。 */
  const submitError = ref('')
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

  function validate(): boolean {
    errors.name = form.name.trim() ? undefined : t('form.required')
    // 联系电话必填，工作邮箱选填；填了邮箱才校验格式。
    errors.email = !form.email.trim()
      ? undefined
      : emailPattern.test(form.email)
        ? undefined
        : t('form.invalidEmail')
    errors.phone = form.phone.trim() ? undefined : t('form.required')
    errors.message = form.message.trim() ? undefined : t('form.required')
    return !errors.name && !errors.email && !errors.phone && !errors.message
  }

  async function submit(): Promise<boolean> {
    submitError.value = ''
    if (!validate()) return false
    sending.value = true
    await new Promise((resolve) => setTimeout(resolve, 450))
    sending.value = false
    sent.value = true
    notify(t('form.successTitle'))
    return true
  }

  function reset() {
    form.name = ''
    form.email = ''
    form.company = ''
    form.phone = ''
    form.message = ''
    errors.name = undefined
    errors.email = undefined
    errors.phone = undefined
    errors.message = undefined
    submitError.value = ''
    sent.value = false
  }

  return { form, errors, sending, sent, submitError, submit, reset }
}

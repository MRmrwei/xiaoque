import type { ContactInput } from '~/domain/types'
import { submitContact } from '~/services/contactService'

type ContactField = 'name' | 'email' | 'phone' | 'message'

export function useContactForm() {
  const { t, locale } = useUiText()
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
  /** 提交失败提示；显示在表单按钮上方，避免只出现在角落的 Toast 里。 */
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
    try {
      await submitContact({ ...form }, locale.value)
      sent.value = true
      notify(t('form.successTitle'))
      return true
    } catch (error) {
      submitError.value =
        error instanceof Error ? error.message : t('form.submitFailed')
      return false
    } finally {
      sending.value = false
    }
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

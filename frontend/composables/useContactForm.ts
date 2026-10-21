import type { ContactInput } from '~/domain/types'
import { submitContact } from '~/services/contactService'

type ContactField = 'name' | 'email' | 'message'

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

  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

  function validate(): boolean {
    errors.name = form.name.trim() ? undefined : t('form.required')
    errors.email = !form.email.trim()
      ? t('form.required')
      : emailPattern.test(form.email)
        ? undefined
        : t('form.invalidEmail')
    errors.message = form.message.trim() ? undefined : t('form.required')
    return !errors.name && !errors.email && !errors.message
  }

  async function submit(): Promise<boolean> {
    if (!validate()) return false
    sending.value = true
    await submitContact({ ...form })
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
    errors.message = undefined
    sent.value = false
  }

  return { form, errors, sending, sent, submit, reset }
}

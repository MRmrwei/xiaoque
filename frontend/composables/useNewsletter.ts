export function useNewsletter() {
  const { t } = useUiText()
  const { notify } = useToast()

  const email = ref('')
  const subscribed = ref(false)

  function subscribe() {
    if (!email.value) return
    subscribed.value = true
    email.value = ''
    notify(t('footer.subscribed'))
  }

  return { email, subscribed, subscribe }
}

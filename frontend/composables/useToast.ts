import type { ToastItem } from '~/domain/types'

export function useToast() {
  const toasts = useState<ToastItem[]>('site-toasts', () => [])

  function notify(text: string, type: ToastItem['type'] = 'success') {
    const id = Date.now() + Math.random()
    toasts.value.push({ id, text, type })
    setTimeout(() => {
      toasts.value = toasts.value.filter((toast) => toast.id !== id)
    }, 3200)
  }

  return { toasts, notify }
}

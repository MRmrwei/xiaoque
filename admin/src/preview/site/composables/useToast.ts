import { reactive } from 'vue'
import type { ToastItem } from '~/domain/types'

export const previewToasts = reactive<ToastItem[]>([])

export function useToast() {
  function notify(text: string, type: ToastItem['type'] = 'success') {
    const id = Date.now() + Math.random()
    previewToasts.push({ id, text, type })
    setTimeout(() => {
      const index = previewToasts.findIndex((toast) => toast.id === id)
      if (index >= 0) previewToasts.splice(index, 1)
    }, 3200)
  }

  return { toasts: previewToasts, notify }
}

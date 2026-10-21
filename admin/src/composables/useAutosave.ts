import { computed, onBeforeUnmount, ref } from 'vue'
import type { PageRecord } from '~/domain/types'
import type { EditorSaveState } from '@/types/admin'
import { formatClock } from '@/services/pageService'

export function useAutosave(
  getPage: () => PageRecord | null,
  save: (page: PageRecord) => Promise<void> | void,
  delay = 1500,
  onError?: (error: unknown) => void
) {
  const state = ref<EditorSaveState>('saved')
  const lastSaved = ref('')
  let timer: ReturnType<typeof setTimeout> | null = null
  let inFlight: Promise<boolean> | null = null

  function markDirty() {
    state.value = 'dirty'
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      void saveNow()
    }, delay)
  }

  async function saveNow(): Promise<boolean> {
    const page = getPage()
    if (!page) return true
    if (timer) {
      clearTimeout(timer)
      timer = null
    }

    // 同一时间只保留一个保存请求；等待前一次完成后再提交最新草稿。
    if (inFlight) await inFlight

    state.value = 'saving'
    const task = (async () => {
      try {
        await save(page)
        lastSaved.value = formatClock(new Date())
        if (state.value !== 'dirty') state.value = 'saved'
        return true
      } catch (error) {
        state.value = 'dirty'
        onError?.(error)
        return false
      }
    })()
    inFlight = task
    try {
      return await task
    } finally {
      if (inFlight === task) inFlight = null
    }
  }

  function flush(): Promise<boolean> {
    if (state.value === 'dirty' || state.value === 'saving') return saveNow()
    return Promise.resolve(true)
  }

  function reset() {
    if (timer) clearTimeout(timer)
    timer = null
    state.value = 'saved'
    lastSaved.value = ''
  }

  const label = computed(() => {
    if (state.value === 'dirty') return '有未保存修改'
    if (state.value === 'saving') return '保存中…'
    return lastSaved.value ? `已自动保存 ${lastSaved.value}` : '已保存'
  })

  onBeforeUnmount(() => {
    void flush()
  })

  return { state, lastSaved, label, markDirty, saveNow, flush, reset }
}

import { computed, onBeforeUnmount, ref } from 'vue'
import type { PageRecord } from '~/domain/types'
import type { EditorSaveState } from '@/types/admin'
import { formatClock } from '@/services/pageService'

export function useAutosave(
  getPage: () => PageRecord | null,
  save: (page: PageRecord) => void,
  delay = 1500
) {
  const state = ref<EditorSaveState>('saved')
  const lastSaved = ref('')
  let timer: ReturnType<typeof setTimeout> | null = null

  function markDirty() {
    state.value = 'dirty'
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      saveNow()
    }, delay)
  }

  function saveNow() {
    const page = getPage()
    if (!page) return
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
    state.value = 'saving'
    save(page)
    lastSaved.value = formatClock(new Date())
    state.value = 'saved'
  }

  function flush() {
    if (state.value === 'dirty') saveNow()
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

  onBeforeUnmount(flush)

  return { state, lastSaved, label, markDirty, saveNow, flush, reset }
}

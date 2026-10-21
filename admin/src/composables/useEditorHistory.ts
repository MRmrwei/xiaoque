import { computed, onBeforeUnmount, ref, type Ref } from 'vue'
import type { PageRecord, SeoMeta, PageBlock } from '~/domain/types'

interface EditorSnapshot {
  name: string
  path: string
  seo: SeoMeta
  blocks: PageBlock[]
}

export function useEditorHistory(page: Readonly<Ref<PageRecord | null>>) {
  const history = ref<string[]>([])
  const index = ref(-1)
  let groupTimer: ReturnType<typeof setTimeout> | null = null

  function snapshot(): string {
    if (!page.value) return ''
    return JSON.stringify({
      name: page.value.name,
      path: page.value.path,
      seo: page.value.seo,
      blocks: page.value.blocks
    })
  }

  function reset() {
    if (groupTimer) clearTimeout(groupTimer)
    groupTimer = null
    history.value = page.value ? [snapshot()] : []
    index.value = history.value.length - 1
  }

  function record(force = false) {
    if (!page.value) return
    const current = snapshot()
    if (!current || history.value[index.value] === current) return

    if (!force && groupTimer && index.value >= 0) {
      history.value[index.value] = current
    } else {
      history.value = history.value.slice(0, index.value + 1)
      history.value.push(current)
      if (history.value.length > 60) history.value.shift()
      index.value = history.value.length - 1
    }

    if (groupTimer) clearTimeout(groupTimer)
    groupTimer = setTimeout(() => {
      groupTimer = null
    }, 700)
  }

  function apply(snapshotValue: string) {
    if (!page.value) return
    const data = JSON.parse(snapshotValue) as EditorSnapshot
    page.value.name = data.name
    page.value.path = data.path
    page.value.seo = data.seo
    page.value.blocks = data.blocks
  }

  function undo(): boolean {
    if (index.value <= 0) return false
    index.value -= 1
    apply(history.value[index.value])
    return true
  }

  function redo(): boolean {
    if (index.value >= history.value.length - 1) return false
    index.value += 1
    apply(history.value[index.value])
    return true
  }

  const canUndo = computed(() => index.value > 0)
  const canRedo = computed(() => index.value < history.value.length - 1)

  onBeforeUnmount(() => {
    if (groupTimer) clearTimeout(groupTimer)
  })

  return { reset, record, undo, redo, canUndo, canRedo }
}

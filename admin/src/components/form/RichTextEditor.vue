<template>
  <div class="rt-editor">
    <div class="rt-toolbar" role="toolbar" aria-label="富文本工具栏">
      <div class="rt-group">
        <button
          type="button"
          class="rt-btn"
          title="撤销"
          :disabled="!state.canUndo"
          @click="editor?.chain().focus().undo().run()"
        >
          <UiIcon name="undo" :size="14" />
        </button>
        <button
          type="button"
          class="rt-btn"
          title="重做"
          :disabled="!state.canRedo"
          @click="editor?.chain().focus().redo().run()"
        >
          <UiIcon name="redo" :size="14" />
        </button>
      </div>

      <span class="rt-sep" />

      <div class="rt-group">
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="二级标题"
          :class="{ active: state.h2 }"
          @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()"
        >
          H2
        </button>
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="三级标题"
          :class="{ active: state.h3 }"
          @click="editor?.chain().focus().toggleHeading({ level: 3 }).run()"
        >
          H3
        </button>
      </div>

      <span class="rt-sep" />

      <div class="rt-group">
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="加粗"
          :class="{ active: state.bold }"
          @click="editor?.chain().focus().toggleBold().run()"
        >
          <strong>B</strong>
        </button>
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="斜体"
          :class="{ active: state.italic }"
          @click="editor?.chain().focus().toggleItalic().run()"
        >
          <em>I</em>
        </button>
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="下划线"
          :class="{ active: state.underline }"
          @click="editor?.chain().focus().toggleUnderline().run()"
        >
          <u>U</u>
        </button>
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="删除线"
          :class="{ active: state.strike }"
          @click="editor?.chain().focus().toggleStrike().run()"
        >
          <s>S</s>
        </button>
      </div>

      <span class="rt-sep" />

      <div class="rt-group">
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="无序列表"
          :class="{ active: state.bulletList }"
          @click="editor?.chain().focus().toggleBulletList().run()"
        >
          • 列表
        </button>
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="有序列表"
          :class="{ active: state.orderedList }"
          @click="editor?.chain().focus().toggleOrderedList().run()"
        >
          1. 列表
        </button>
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="引用"
          :class="{ active: state.blockquote }"
          @click="editor?.chain().focus().toggleBlockquote().run()"
        >
          ❝ 引用
        </button>
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="代码块"
          :class="{ active: state.codeBlock }"
          @click="editor?.chain().focus().toggleCodeBlock().run()"
        >
          &lt;/&gt;
        </button>
      </div>

      <span class="rt-sep" />

      <div class="rt-group">
        <button
          type="button"
          class="rt-btn"
          title="插入链接"
          :class="{ active: state.link }"
          @click="setLink"
        >
          <UiIcon name="link" :size="14" />
        </button>
        <button
          type="button"
          class="rt-btn"
          title="插入图片"
          @click="imageInput?.click()"
        >
          <UiIcon name="image" :size="14" />
        </button>
        <button
          type="button"
          class="rt-btn"
          title="分割线"
          @click="editor?.chain().focus().setHorizontalRule().run()"
        >
          —
        </button>
        <button
          type="button"
          class="rt-btn rt-btn--text"
          title="清除格式"
          @click="clearFormat"
        >
          Tx
        </button>
      </div>

      <input
        ref="imageInput"
        class="rt-file-input"
        type="file"
        accept="image/*"
        @change="onImagePick"
      />
    </div>

    <div class="rt-content">
      <EditorContent :editor="editor" />
      <span v-if="state.empty" class="rt-placeholder">在这里开始写作…</span>
    </div>

    <div class="rt-foot">
      <span>{{ state.characters }} 字</span>
      <span>支持标题、列表、引用、代码块与链接</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  estimateDataUrlBytes,
  fileToImageDataUrl,
  formatBytes
} from '@/utils/image'
import UiIcon from '~/components/UiIcon.vue'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const revision = ref(0)

function isSafeLink(url: string): boolean {
  const value = url.trim()
  return !value || /^(https?:\/\/|mailto:|tel:|\/)/i.test(value)
}

const editor = useEditor({
  content: props.modelValue || '',
  extensions: [
    StarterKit.configure({
      heading: { levels: [2, 3] },
      link: {
        openOnClick: false,
        autolink: true,
        defaultProtocol: 'https',
        isAllowedUri: (url) => isSafeLink(String(url ?? ''))
      }
    }),
    Image.configure({
      allowBase64: true,
      inline: false,
      HTMLAttributes: { class: 'rt-image' }
    })
  ],
  editorProps: {
    attributes: {
      class: 'rt-prose',
      'aria-label': '富文本内容'
    },
    handlePaste: (_view, event) => handleImageFiles(event.clipboardData?.files),
    handleDrop: (_view, event) =>
      handleImageFiles((event as DragEvent).dataTransfer?.files)
  },
  onUpdate: ({ editor: current }) => {
    emit('update:modelValue', current.getHTML())
  },
  onTransaction: () => {
    revision.value += 1
  }
})

const state = computed(() => {
  revision.value
  const instance = editor.value
  return {
    bold: instance?.isActive('bold') ?? false,
    italic: instance?.isActive('italic') ?? false,
    underline: instance?.isActive('underline') ?? false,
    strike: instance?.isActive('strike') ?? false,
    h2: instance?.isActive('heading', { level: 2 }) ?? false,
    h3: instance?.isActive('heading', { level: 3 }) ?? false,
    bulletList: instance?.isActive('bulletList') ?? false,
    orderedList: instance?.isActive('orderedList') ?? false,
    blockquote: instance?.isActive('blockquote') ?? false,
    codeBlock: instance?.isActive('codeBlock') ?? false,
    link: instance?.isActive('link') ?? false,
    canUndo: instance?.can().undo() ?? false,
    canRedo: instance?.can().redo() ?? false,
    empty: instance?.isEmpty ?? true,
    characters: instance?.getText().replace(/\s+/g, '').length ?? 0
  }
})

watch(
  () => props.modelValue,
  (value) => {
    const instance = editor.value
    if (!instance) return
    if (instance.getHTML() !== (value || '')) {
      instance.commands.setContent(value || '')
    }
  }
)

async function setLink() {
  const instance = editor.value
  if (!instance) return
  const previous = instance.getAttributes('link').href as string | undefined
  try {
    const { value } = await ElMessageBox.prompt(
      '支持 https://、mailto:、tel: 或以 / 开头的站内路径，留空可移除链接。',
      '插入链接',
      {
        inputValue: previous || 'https://',
        confirmButtonText: '应用',
        cancelButtonText: '取消',
        inputPlaceholder: 'https://example.com'
      }
    )
    const href = value.trim()
    if (!href) {
      instance.chain().focus().extendMarkRange('link').unsetLink().run()
      return
    }
    if (!isSafeLink(href)) {
      ElMessage.warning('请输入 http(s)、mailto、tel 或以 / 开头的站内路径')
      return
    }
    instance.chain().focus().extendMarkRange('link').setLink({ href }).run()
  } catch {
    // cancelled
  }
}

function clearFormat() {
  editor.value?.chain().focus().unsetAllMarks().clearNodes().run()
}

const imageInput = ref<HTMLInputElement | null>(null)

function onImagePick(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (file) void insertImageFile(file)
}

function handleImageFiles(files?: FileList | null): boolean {
  const file = Array.from(files ?? []).find((item) => item.type.startsWith('image/'))
  if (!file) return false
  void insertImageFile(file)
  return true
}

async function insertImageFile(file: File) {
  try {
    const dataUrl = await fileToImageDataUrl(file)
    editor.value
      ?.chain()
      .focus()
      .setImage({
        src: dataUrl,
        alt: file.name.replace(/\.[^.]+$/, ''),
        title: file.name
      })
      .run()
    ElMessage.success(`图片已插入（${formatBytes(estimateDataUrlBytes(dataUrl))}）`)
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '图片处理失败')
  }
}

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style scoped>
.rt-editor {
  border: 1px solid var(--adm-line);
  border-radius: 10px;
  background: #fff;
  overflow: hidden;
}

.rt-toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 2px;
  padding: 6px 8px;
  border-bottom: 1px solid var(--adm-line-soft);
  background: #faf8f4;
}

.rt-group {
  display: flex;
  align-items: center;
  gap: 2px;
}

.rt-sep {
  width: 1px;
  height: 16px;
  margin: 0 5px;
  background: var(--adm-line);
}

.rt-btn {
  min-width: 28px;
  height: 28px;
  display: inline-grid;
  place-items: center;
  padding: 0 7px;
  border: 0;
  border-radius: 5px;
  background: transparent;
  color: var(--adm-muted);
  font-size: 11px;
  line-height: 1;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.rt-btn:hover {
  background: rgba(20, 21, 24, 0.06);
  color: var(--adm-ink);
}

.rt-btn.active {
  background: var(--adm-ink);
  color: #f6f2ea;
}

.rt-btn:disabled {
  opacity: 0.3;
  cursor: default;
  background: transparent;
  color: var(--adm-muted);
}

.rt-btn--text {
  font-family: var(--adm-font);
  font-weight: 600;
  letter-spacing: -0.01em;
}

.rt-content {
  position: relative;
  min-height: 220px;
}

.rt-placeholder {
  position: absolute;
  left: 16px;
  top: 14px;
  font-size: 12.5px;
  color: var(--adm-muted-2);
  pointer-events: none;
}

.rt-content :deep(.rt-prose) {
  min-height: 220px;
  padding: 14px 16px;
  outline: none;
  color: var(--adm-ink);
  font-size: 13px;
  line-height: 1.75;
}

.rt-content :deep(.rt-prose > *:first-child) {
  margin-top: 0;
}

.rt-content :deep(.rt-prose > *:last-child) {
  margin-bottom: 0;
}

.rt-content :deep(.rt-prose p) {
  margin: 0 0 10px;
}

.rt-content :deep(.rt-prose h2) {
  margin: 18px 0 10px;
  font-size: 20px;
  line-height: 1.3;
  letter-spacing: -0.02em;
}

.rt-content :deep(.rt-prose h3) {
  margin: 16px 0 8px;
  font-size: 16px;
  line-height: 1.35;
  letter-spacing: -0.015em;
}

.rt-content :deep(.rt-prose ul),
.rt-content :deep(.rt-prose ol) {
  margin: 0 0 10px;
  padding-left: 22px;
}

.rt-content :deep(.rt-prose li) {
  margin: 3px 0;
}

.rt-content :deep(.rt-prose blockquote) {
  margin: 12px 0;
  padding: 4px 0 4px 14px;
  border-left: 2px solid var(--adm-accent);
  color: var(--adm-muted);
}

.rt-content :deep(.rt-prose pre) {
  margin: 12px 0;
  padding: 12px 14px;
  border-radius: 8px;
  background: #141518;
  color: #f6f2ea;
  font-family: var(--adm-mono);
  font-size: 11.5px;
  line-height: 1.7;
  overflow-x: auto;
}

.rt-content :deep(.rt-prose code) {
  padding: 1px 5px;
  border-radius: 4px;
  background: rgba(20, 21, 24, 0.06);
  font-family: var(--adm-mono);
  font-size: 0.92em;
}

.rt-content :deep(.rt-prose pre code) {
  padding: 0;
  background: transparent;
  color: inherit;
}

.rt-content :deep(.rt-prose a) {
  color: var(--adm-accent);
  text-decoration: underline;
  text-underline-offset: 2px;
}

.rt-content :deep(.rt-prose img) {
  display: block;
  max-width: 100%;
  height: auto;
  margin: 14px 0;
  border-radius: 8px;
}

.rt-content :deep(.rt-prose img.ProseMirror-selectednode) {
  outline: 2px solid var(--adm-accent);
  outline-offset: 2px;
}

.rt-file-input {
  display: none;
}

.rt-content :deep(.rt-prose hr) {
  margin: 18px 0;
  border: 0;
  border-top: 1px solid var(--adm-line);
}

.rt-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 6px 10px;
  border-top: 1px solid var(--adm-line-soft);
  background: #fdfcfa;
  font-size: 9px;
  color: var(--adm-muted-2);
}
</style>

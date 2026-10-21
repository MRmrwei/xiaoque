<template>
  <div class="footer-panel">
    <section class="panel-intro">
      <div>
        <p class="panel-kicker">FOOTER NAVIGATION</p>
        <h2 class="panel-title">底部导航</h2>
        <p class="panel-sub">
          编排官网页脚的栏目、链接、订阅区、联系入口与底部版权信息。中英文分开配置，修改后点击右上角「保存设置」，并同步到页面草稿预览。
        </p>
      </div>
      <div class="panel-actions">
        <span class="save-state" :class="saveState"><i />{{ saveLabel }}</span>
        <el-button size="small" @click="resetFooter">恢复默认</el-button>
        <el-button size="small" @click="addColumn">
          <el-icon><Plus /></el-icon>
          添加栏目
        </el-button>
        <el-button
          size="small"
          type="primary"
          :disabled="saveState === 'saved'"
          @click="saveFooter"
        >
          保存设置
        </el-button>
      </div>
    </section>

    <aside class="adm-card footer-preview-card">
      <div class="footer-preview-head">
        <div>
          <h3>实时预览</h3>
          <p>只显示页脚区域 · 中英文 / 桌面与手机</p>
        </div>
        <div class="footer-preview-actions">
          <div class="footer-segmented">
            <button :class="{ active: previewLocale === 'zh' }" @click="previewLocale = 'zh'">
              中文
            </button>
            <button :class="{ active: previewLocale === 'en' }" @click="previewLocale = 'en'">
              EN
            </button>
          </div>
          <div class="footer-segmented">
            <button :class="{ active: device === 'desktop' }" @click="device = 'desktop'">
              桌面
            </button>
            <button :class="{ active: device === 'mobile' }" @click="device = 'mobile'">
              手机
            </button>
          </div>
        </div>
      </div>

      <div class="footer-preview-stage">
        <div class="footer-preview-frame" :style="[brandVars, { maxWidth: previewWidth }]">
          <PreviewFooter :settings="settings" :locale="previewLocale" />
        </div>
      </div>

      <div class="footer-preview-note">
        <span class="mono">LIVE</span>
        <p>预览中的链接不可点击，路径会随设置一起保存；手机宽度下按容器宽度折叠栏目。</p>
      </div>
    </aside>

    <div class="footer-layout">
      <section class="adm-card footer-card">
        <div class="adm-card-head">
          <div>
            <h3>栏目编排</h3>
            <p>拖拽左侧手柄调整栏目顺序；栏目和链接都可以单独隐藏。</p>
          </div>
          <span class="mono">{{ footer.columns.length }} COLUMNS</span>
        </div>

        <div class="footer-columns-table">
          <div class="footer-columns-head">
            <span />
            <span>栏目名称（中文 / English）</span>
            <span>显示</span>
          </div>

          <VueDraggable
            :model-value="footer.columns"
            :animation="150"
            handle=".footer-column-drag"
            class="footer-columns"
            @update:model-value="reorderColumns"
          >
            <div
              v-for="(column, index) in footer.columns"
              :key="column.key"
              class="footer-column"
              :class="{ 'is-hidden': !column.visible }"
            >
              <div class="footer-column-head">
                <span class="footer-column-drag" title="拖动排序">
                  <UiIcon name="drag" :size="13" />
                </span>
                <div class="footer-column-titles">
                  <el-input
                    v-model="column.title.zh"
                    size="small"
                    class="footer-column-title footer-column-title--zh"
                    placeholder="栏目中文名称"
                  />
                  <el-input
                    v-model="column.title.en"
                    size="small"
                    class="footer-column-title footer-column-title--en"
                    placeholder="Column title"
                  />
                </div>
                <div class="footer-column-actions">
                  <el-switch
                    v-model="column.visible"
                    size="small"
                    :title="column.visible ? '显示中，点击隐藏' : '已隐藏，点击显示'"
                  />
                  <button
                    class="footer-column-remove"
                    title="删除栏目"
                    @click="removeColumn(index)"
                  >
                    <UiIcon name="trash" :size="13" />
                  </button>
                </div>
              </div>

              <FooterLinkRows v-model="column.links" />

              <button class="footer-add-link" @click="addLink(column)">
                <UiIcon name="plus" :size="12" />
                添加链接
              </button>
            </div>
          </VueDraggable>
        </div>

        <div v-if="!footer.columns.length" class="footer-empty">
          还没有页脚栏目，点右上角「添加栏目」开始。
        </div>
      </section>

      <div class="footer-side">
        <section class="adm-card footer-card">
          <div class="adm-card-head">
            <div>
              <h3>页脚信息</h3>
              <p>订阅区、联系入口与版权信息。</p>
            </div>
          </div>

          <div class="footer-settings">
            <div class="footer-setting-row">
              <div class="footer-setting-label">
                <strong>品牌标语</strong>
                <span>显示在页脚品牌名称下方</span>
              </div>
              <el-input
                v-model="settings.tagline.zh"
                size="small"
                type="textarea"
                :rows="2"
                placeholder="中文标语"
              />
              <el-input
                v-model="settings.tagline.en"
                size="small"
                type="textarea"
                :rows="2"
                placeholder="English tagline"
              />
            </div>

            <div class="footer-setting-row">
              <div class="footer-setting-label">
                <strong>联系方式</strong>
                <span>页脚联系列展示的邮箱、电话与地址</span>
              </div>
              <el-input v-model="settings.email" size="small" placeholder="hello@xiaoque.com" />
              <div class="footer-pair">
                <el-input v-model="settings.phone" size="small" placeholder="电话（中文）" />
                <el-input v-model="settings.phoneEn" size="small" placeholder="Phone (EN)" />
              </div>
              <div class="footer-pair">
                <el-input v-model="settings.address.zh" size="small" placeholder="中文地址" />
                <el-input v-model="settings.address.en" size="small" placeholder="Address (EN)" />
              </div>
            </div>

            <div class="footer-setting-row">
              <div class="footer-setting-label">
                <strong>订阅区</strong>
                <span>页脚左侧的邮件订阅入口</span>
              </div>
              <div class="footer-setting-inline">
                <el-switch v-model="footer.showNewsletter" size="small" />
                <span class="footer-setting-hint">关闭后整块隐藏</span>
              </div>
              <div class="footer-pair" :class="{ 'is-disabled': !footer.showNewsletter }">
                <el-input
                  v-model="footer.newsletterTitle.zh"
                  size="small"
                  placeholder="订阅标题"
                  :disabled="!footer.showNewsletter"
                />
                <el-input
                  v-model="footer.newsletterTitle.en"
                  size="small"
                  placeholder="Newsletter title"
                  :disabled="!footer.showNewsletter"
                />
              </div>
              <div class="footer-pair" :class="{ 'is-disabled': !footer.showNewsletter }">
                <el-input
                  v-model="footer.newsletterNote.zh"
                  size="small"
                  placeholder="订阅说明"
                  :disabled="!footer.showNewsletter"
                />
                <el-input
                  v-model="footer.newsletterNote.en"
                  size="small"
                  placeholder="Newsletter note"
                  :disabled="!footer.showNewsletter"
                />
              </div>
            </div>

            <div class="footer-setting-row">
              <div class="footer-setting-label">
                <strong>联系列</strong>
                <span>邮箱 / 电话 / 地址沿用站点设置</span>
              </div>
              <div class="footer-setting-inline">
                <el-switch v-model="footer.showContact" size="small" />
                <span class="footer-setting-hint">右侧联系栏目与 CTA</span>
              </div>
              <div class="footer-pair" :class="{ 'is-disabled': !footer.showContact }">
                <el-input
                  v-model="footer.cta.label.zh"
                  size="small"
                  placeholder="按钮中文文案"
                  :disabled="!footer.showContact"
                />
                <el-input
                  v-model="footer.cta.label.en"
                  size="small"
                  placeholder="CTA label"
                  :disabled="!footer.showContact"
                />
              </div>
              <el-input
                v-model="footer.cta.path"
                size="small"
                class="footer-path"
                placeholder="/contact"
                :disabled="!footer.showContact"
              />
            </div>

            <div class="footer-setting-row">
              <div class="footer-setting-label">
                <strong>底部大字标</strong>
                <span>页脚中部的大号品牌字</span>
              </div>
              <div class="footer-setting-inline">
                <el-switch v-model="footer.showWordmark" size="small" />
                <el-input
                  v-model="footer.wordmark"
                  size="small"
                  class="footer-wordmark-input"
                  placeholder="MAGPIE"
                  :disabled="!footer.showWordmark"
                />
              </div>
            </div>

            <div class="footer-setting-row">
              <div class="footer-setting-label">
                <strong>版权与备案</strong>
                <span>显示在页脚最底部</span>
              </div>
              <div class="footer-pair">
                <el-input v-model="settings.copyright.zh" size="small" placeholder="中文版权" />
                <el-input v-model="settings.copyright.en" size="small" placeholder="Copyright" />
              </div>
              <el-input v-model="settings.icp" size="small" placeholder="ICP 备案号" />
            </div>

            <div class="footer-setting-row">
              <div class="footer-setting-label">
                <strong>状态提示</strong>
                <span>页脚右下角的运行状态</span>
              </div>
              <div class="footer-setting-inline">
                <el-switch v-model="footer.showStatus" size="small" />
                <span class="footer-setting-hint">带绿点的状态文案</span>
              </div>
              <div class="footer-pair" :class="{ 'is-disabled': !footer.showStatus }">
                <el-input
                  v-model="footer.statusText.zh"
                  size="small"
                  placeholder="中文状态"
                  :disabled="!footer.showStatus"
                />
                <el-input
                  v-model="footer.statusText.en"
                  size="small"
                  placeholder="Status text"
                  :disabled="!footer.showStatus"
                />
              </div>
            </div>
          </div>
        </section>

        <section class="adm-card footer-card">
          <div class="adm-card-head">
            <div>
              <h3>底部链接</h3>
              <p>隐私政策、服务条款等，显示在版权旁边。</p>
            </div>
            <span class="mono">{{ footer.bottomLinks.length }} LINKS</span>
          </div>

          <div class="footer-bottom-editor">
            <FooterLinkRows v-model="footer.bottomLinks" />
            <button class="footer-add-link" @click="addBottomLink">
              <UiIcon name="plus" :size="12" />
              添加底部链接
            </button>
          </div>
        </section>

        <p v-if="issues.length" class="footer-warning">{{ issues.join('；') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { VueDraggable } from 'vue-draggable-plus'
import type { FooterColumn, FooterLink, Locale } from '~/domain/types'
import type { EditorSaveState } from '@/types/admin'
import { usePagesStore } from '@/stores/pages'
import PreviewFooter from '@/preview/PreviewFooter.vue'
import { brandStyle } from '@/utils/brand'
import FooterLinkRows from './FooterLinkRows.vue'
import UiIcon from '~/components/UiIcon.vue'

const pagesStore = usePagesStore()
const { settings } = storeToRefs(pagesStore)

const footer = computed(() => settings.value.footer)
const previewLocale = ref<Locale>('zh')
const device = ref<'desktop' | 'mobile'>('desktop')
const previewWidth = computed(() => (device.value === 'mobile' ? '420px' : '100%'))
const brandVars = computed(() => brandStyle(settings.value.brand))

const saveState = ref<EditorSaveState>('saved')
const saveLabel = computed(() => {
  if (saveState.value === 'dirty') return '有未保存修改'
  if (saveState.value === 'saving') return '正在保存…'
  return `已保存 ${pagesStore.lastSavedAt}`
})

watch(
  () => [
    settings.value.footer,
    settings.value.tagline,
    settings.value.copyright,
    settings.value.icp,
    settings.value.email,
    settings.value.phone,
    settings.value.phoneEn,
    settings.value.address
  ],
  () => {
    saveState.value = 'dirty'
  },
  { deep: true }
)

function saveFooter() {
  saveState.value = 'saving'
  pagesStore.saveSettings()
  saveState.value = 'saved'
  ElMessage.success('底部导航设置已保存')
}

const issues = computed(() => {
  const list: string[] = []
  const links = footer.value.columns.flatMap((column) => column.links)
  const empty = [...links, ...footer.value.bottomLinks].filter((link) => !link.path.trim()).length
  if (empty) list.push(`${empty} 条底部链接还没有填写跳转路径`)
  if (!footer.value.columns.some((column) => column.visible)) {
    list.push('当前没有可见的页脚栏目')
  }
  return list
})

function newLink(): FooterLink {
  return {
    key: `footer.custom_${Date.now().toString(36)}`,
    path: '/new-page',
    zh: '新链接',
    en: 'New link',
    visible: true
  }
}

function reorderColumns(value: FooterColumn[]) {
  settings.value.footer.columns = value
}

function addColumn() {
  settings.value.footer.columns.push({
    key: `footer.custom_col_${Date.now().toString(36)}`,
    title: { zh: '新栏目', en: 'New column' },
    visible: true,
    links: [newLink()]
  })
}

async function removeColumn(index: number) {
  const column = settings.value.footer.columns[index]
  if (!column) return
  try {
    await ElMessageBox.confirm(
      `确定删除栏目「${column.title.zh || column.title.en}」及其中 ${column.links.length} 条链接吗？`,
      '删除页脚栏目',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
    settings.value.footer.columns.splice(index, 1)
  } catch {
    // cancelled
  }
}

function addLink(column: FooterColumn) {
  column.links.push(newLink())
}

function addBottomLink() {
  settings.value.footer.bottomLinks.push(newLink())
}

async function resetFooter() {
  try {
    await ElMessageBox.confirm(
      '会把页脚栏目、链接和页脚信息恢复为初始数据，确定继续吗？',
      '恢复默认底部导航',
      { type: 'warning', confirmButtonText: '恢复默认', cancelButtonText: '取消' }
    )
    pagesStore.resetFooter()
    saveState.value = 'dirty'
    ElMessage.success('底部导航已恢复默认，点击「保存设置」生效')
  } catch {
    // cancelled
  }
}
</script>

<style scoped>
.panel-intro {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 18px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.panel-kicker {
  font-family: var(--adm-mono);
  font-size: 8px;
  letter-spacing: 0.16em;
  color: var(--adm-muted-2);
}

.panel-title {
  margin-top: 5px;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.025em;
}

.panel-sub {
  margin-top: 5px;
  font-size: 11.5px;
  color: var(--adm-muted);
  max-width: 76ch;
}

.panel-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.footer-panel .save-state {
  display: inline-flex;
}

.footer-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(330px, 0.75fr);
  gap: 16px;
  align-items: start;
}

.footer-side {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.footer-preview-card {
  margin-bottom: 16px;
  overflow: hidden;
}

.footer-preview-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 16px 18px 12px;
  flex-wrap: wrap;
}

.footer-preview-head h3 {
  font-size: 13px;
  font-weight: 600;
}

.footer-preview-head p {
  margin-top: 3px;
  font-size: 10px;
  color: var(--adm-muted-2);
}

.footer-preview-actions {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.footer-segmented {
  display: inline-flex;
  padding: 2px;
  gap: 2px;
  border: 1px solid var(--adm-line);
  border-radius: 7px;
  background: var(--adm-surface);
}

.footer-segmented button {
  border: 0;
  background: transparent;
  padding: 3px 9px;
  border-radius: 5px;
  font-size: 9.5px;
  color: var(--adm-muted);
  cursor: pointer;
}

.footer-segmented button.active {
  background: #fff;
  color: var(--adm-ink);
  box-shadow: 0 1px 2px rgba(20, 21, 24, 0.12);
}

.footer-preview-stage {
  padding: 16px;
  background:
    linear-gradient(45deg, rgba(20, 21, 24, 0.025) 25%, transparent 25%) 0 0 / 16px 16px,
    linear-gradient(-45deg, rgba(20, 21, 24, 0.025) 25%, transparent 25%) 0 8px / 16px 16px,
    linear-gradient(45deg, transparent 75%, rgba(20, 21, 24, 0.025) 75%) 8px -8px / 16px 16px,
    linear-gradient(-45deg, transparent 75%, rgba(20, 21, 24, 0.025) 75%) -8px 0 / 16px 16px,
    var(--adm-bg);
}

.footer-preview-frame {
  width: 100%;
  margin: 0 auto;
  border: 1px solid var(--adm-line);
  border-radius: 10px;
  overflow: hidden;
  box-shadow: var(--adm-shadow);
  background: #141518;

  /* 预览字体与预览站保持一致。 */
  --font-serif: var(--adm-serif);
  --font-mono: var(--adm-mono);
}

.footer-preview-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 18px 14px;
  color: var(--adm-muted-2);
}

.footer-preview-note .mono {
  color: var(--adm-accent);
  font-size: 8px;
  letter-spacing: 0.12em;
  padding-top: 1px;
}

.footer-preview-note p {
  font-size: 9.5px;
  line-height: 1.6;
}

.footer-card {
  padding-bottom: 12px;
}

.footer-card :deep(.adm-card-head) {
  padding: 16px 18px 12px;
}

.footer-card :deep(.adm-card-head p) {
  margin-top: 3px;
  font-size: 10.5px;
  color: var(--adm-muted);
}

.footer-card .mono {
  font-family: var(--adm-mono);
  font-size: 8.5px;
  color: var(--adm-muted-2);
}

.footer-columns-table {
  margin: 0 18px;
  border: 1px solid var(--adm-line);
  border-radius: 10px;
  overflow: hidden;
  background: #fff;
}

.footer-columns-head {
  display: grid;
  grid-template-columns: 30px minmax(0, 1fr) 82px;
  gap: 0;
  align-items: stretch;
  padding: 0;
  background: #f7f3ec;
  border-bottom: 1px solid var(--adm-line);
  font-size: 11.5px;
  font-weight: 600;
  color: var(--adm-ink-soft);
}

.footer-columns-head > span {
  display: flex;
  align-items: center;
  min-width: 0;
  padding: 9px 10px;
  border-left: 1px solid var(--adm-line-soft);
}

.footer-columns-head > span:first-child {
  border-left: 0;
}

.footer-columns-head > span:last-child {
  justify-content: flex-end;
}

.footer-columns {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 0;
}

.footer-column {
  padding: 10px 12px 8px;
  border: 0;
  border-bottom: 1px solid var(--adm-line);
  border-radius: 0;
  background: #fff;
}

.footer-column:last-child {
  border-bottom: 0;
}

.footer-column.is-hidden {
  opacity: 0.55;
}

.footer-column .footer-link-editor {
  margin-top: 8px;
}

.footer-column-head {
  display: grid;
  grid-template-columns: 18px minmax(0, 1fr) auto;
  gap: 7px;
  align-items: center;
  padding-bottom: 8px;
  margin-bottom: 2px;
  border-bottom: 1px solid var(--adm-line-soft);
}

.footer-column-drag {
  display: grid;
  place-items: center;
  height: 26px;
  color: var(--adm-muted-2);
  cursor: grab;
}

.footer-column-drag:active {
  cursor: grabbing;
}

.footer-column-titles {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px;
  min-width: 0;
}

.footer-column-title :deep(.el-input__wrapper) {
  padding: 2px 6px;
  background: transparent;
  box-shadow: none;
  border-radius: 5px;
  transition: background 0.15s ease, box-shadow 0.15s ease;
}

.footer-column-title :deep(.el-input__wrapper:hover) {
  background: rgba(20, 21, 24, 0.045);
}

.footer-column-title :deep(.el-input__wrapper.is-focus) {
  background: #fff;
  box-shadow: inset 0 0 0 1px var(--adm-line);
}

.footer-column-title--zh :deep(input) {
  color: var(--adm-ink);
  font-weight: 600;
}

.footer-column-title--en :deep(input) {
  color: var(--adm-muted);
}

.footer-column-actions {
  display: flex;
  align-items: center;
  gap: 2px;
}

.footer-column-remove {
  width: 26px;
  height: 26px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: var(--adm-muted-2);
  cursor: pointer;
  opacity: 0.35;
  transition: opacity 0.15s ease, background 0.15s ease, color 0.15s ease;
}

.footer-column:hover .footer-column-remove {
  opacity: 1;
}

.footer-column-remove:hover {
  background: var(--adm-red-soft);
  color: var(--adm-accent);
}

.footer-add-link {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  margin-top: 4px;
  padding: 6px 4px;
  border: 0;
  border-radius: 5px;
  background: transparent;
  font-size: 10px;
  color: var(--adm-muted);
  cursor: pointer;
}

.footer-add-link:hover {
  background: transparent;
  color: var(--adm-accent);
}

.footer-empty {
  margin: 10px 18px;
  padding: 18px;
  border: 1px dashed var(--adm-line);
  border-radius: 8px;
  text-align: center;
  font-size: 11px;
  color: var(--adm-muted-2);
}

.footer-settings {
  padding: 2px 18px 4px;
}

.footer-setting-row {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  padding-block: 12px;
  border-top: 1px solid var(--adm-line-soft);
}

.footer-setting-row:first-child {
  border-top: 0;
}

.footer-setting-label {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.footer-setting-label strong {
  font-size: 11.5px;
  font-weight: 600;
}

.footer-setting-label span {
  font-size: 9.5px;
  color: var(--adm-muted-2);
}

.footer-setting-inline {
  display: flex;
  align-items: center;
  gap: 9px;
  width: 100%;
}

.footer-setting-hint {
  font-size: 9.5px;
  color: var(--adm-muted-2);
}

.footer-pair {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 7px;
  width: 100%;
}

.footer-pair.is-disabled {
  opacity: 0.5;
}

.footer-path :deep(input) {
  font-family: var(--adm-mono);
  font-size: 10.5px;
}

.footer-wordmark-input {
  max-width: 180px;
}

.footer-wordmark-input :deep(input) {
  font-family: var(--adm-mono);
  font-size: 10.5px;
  letter-spacing: 0.08em;
}

.footer-bottom-editor {
  padding: 10px 18px 4px;
}

.footer-warning {
  margin: 0;
  padding: 8px 11px;
  border-radius: 7px;
  background: var(--adm-amber-soft);
  color: var(--adm-amber);
  font-size: 10px;
  line-height: 1.6;
}

@media (max-width: 1180px) {
  .footer-layout {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 760px) {
  .footer-columns-head {
    display: none;
  }
  .footer-columns-table {
    margin: 0 12px;
  }
  .footer-column {
    padding: 8px;
  }
  .footer-column-head {
    grid-template-columns: 18px minmax(0, 1fr);
    grid-template-areas:
      'drag titles'
      'drag actions';
  }
  .footer-column-drag {
    grid-area: drag;
    align-self: start;
  }
  .footer-column-titles {
    grid-area: titles;
    grid-template-columns: 1fr;
  }
  .footer-column-actions {
    grid-area: actions;
    justify-content: flex-end;
  }
  .footer-pair {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>

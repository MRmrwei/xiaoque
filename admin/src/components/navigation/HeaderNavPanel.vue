<template>
  <div class="nav-panel">
    <section class="panel-intro">
      <div>
        <p class="panel-kicker">HEADER NAVIGATION</p>
        <h2 class="panel-title">顶部导航</h2>
        <p class="panel-sub">
          管理官网顶部的菜单顺序、中英文名称和跳转路径，以及右侧按钮、语言切换等头部设置。修改后点击右上角「保存设置」，并同步到页面草稿预览。
        </p>
      </div>
      <div class="panel-actions">
        <span class="save-state" :class="saveState"><i />{{ saveLabel }}</span>
        <el-button size="small" @click="resetNav">恢复默认</el-button>
        <el-button size="small" @click="addItem">
          <el-icon><Plus /></el-icon>
          添加菜单项
        </el-button>
        <el-button
          size="small"
          type="primary"
          :disabled="saveState === 'saved'"
          @click="saveNav"
        >
          保存设置
        </el-button>
      </div>
    </section>

    <aside class="adm-card nav-preview-card">
      <div class="nav-preview-head">
        <div>
          <h3>实时预览</h3>
          <p>只显示顶部导航区域</p>
        </div>
        <div class="nav-preview-actions">
          <div class="nav-segmented">
            <button :class="{ active: previewLocale === 'zh' }" @click="previewLocale = 'zh'">
              中文
            </button>
            <button :class="{ active: previewLocale === 'en' }" @click="previewLocale = 'en'">
              EN
            </button>
          </div>
          <div class="nav-segmented">
            <button :class="{ active: device === 'desktop' }" @click="device = 'desktop'">
              桌面
            </button>
            <button :class="{ active: device === 'mobile' }" @click="device = 'mobile'">
              手机
            </button>
          </div>
        </div>
      </div>

      <div class="nav-preview-stage">
        <div class="nav-preview-frame" :style="[brandVars, { maxWidth: previewWidth }]">
          <PreviewChrome :settings="settings" :locale="previewLocale" />
          <div class="nav-preview-fill" aria-hidden="true">
            <span class="mono">PAGE CONTENT</span>
          </div>
        </div>
      </div>

      <div class="nav-preview-note">
        <span class="mono">LIVE</span>
        <p>预览中的站内跳转被禁用，链接会随设置一起保存。</p>
      </div>
    </aside>

    <div class="nav-layout">
      <div class="nav-main">
      <section class="adm-card nav-card">
        <div class="adm-card-head">
          <div>
            <h3>菜单项</h3>
            <p>拖拽左侧手柄调整顺序；关掉开关会从顶部隐藏。</p>
          </div>
          <span class="mono">{{ settings.nav.length }} ITEMS</span>
        </div>

        <div class="nav-table">
          <div class="nav-head-row nav-head-row--menu">
            <span />
            <span>中文名称</span>
            <span>English</span>
            <span>跳转路径</span>
            <span class="nav-head-actions">显示</span>
          </div>

          <VueDraggable
            :model-value="settings.nav"
            :animation="150"
            handle=".nav-drag"
            class="nav-rows"
            @update:model-value="onReorder"
          >
            <div
              v-for="(item, index) in settings.nav"
              :key="item.key"
              class="nav-row"
              :class="{ 'is-hidden': !item.visible }"
            >
              <span class="nav-drag" title="拖动排序">
                <UiIcon name="drag" :size="13" />
              </span>
              <el-input
                v-model="item.zh"
                size="small"
                class="nav-field nav-field--zh"
                placeholder="中文名称"
              />
              <el-input
                v-model="item.en"
                size="small"
                class="nav-field nav-field--en"
                placeholder="English label"
              />
              <div class="nav-row-path">
                <UiIcon name="link" :size="12" />
                <el-input v-model="item.path" size="small" placeholder="/path" />
              </div>
              <div class="nav-row-actions">
                <el-switch
                  v-model="item.visible"
                  size="small"
                  :title="item.visible ? '显示中，点击隐藏' : '已隐藏，点击显示'"
                />
                <button class="nav-remove" title="删除菜单项" @click="removeItem(index)">
                  <UiIcon name="trash" :size="13" />
                </button>
              </div>
            </div>
          </VueDraggable>
        </div>

        <div v-if="!settings.nav.length" class="nav-empty">
          还没有菜单项，点右上角「添加菜单项」开始。
        </div>
        <p v-if="navIssues.length" class="nav-warning">{{ navIssues.join('；') }}</p>
      </section>

      <section class="adm-card nav-card">
        <div class="adm-card-head">
          <div>
            <h3>顶部按钮</h3>
            <p>按顺序显示在右上角；建议 1 个主按钮 + 最多 2 个辅助按钮。</p>
          </div>
          <div class="nav-cta-head-actions">
            <span class="mono">{{ visibleCtaCount }} / {{ header.ctas.length }} VISIBLE</span>
            <el-button size="small" @click="addCta">
              <el-icon><Plus /></el-icon>
              添加按钮
            </el-button>
          </div>
        </div>

        <div class="nav-table">
          <div class="nav-head-row nav-head-row--cta">
            <span />
            <span>中文文案</span>
            <span>English</span>
            <span>跳转路径</span>
            <span>样式</span>
            <span class="nav-head-actions">显示</span>
          </div>

          <VueDraggable
            :model-value="header.ctas"
            :animation="150"
            handle=".nav-cta-drag"
            class="nav-cta-rows"
            @update:model-value="reorderCtas"
          >
            <div
              v-for="(cta, index) in header.ctas"
              :key="cta.key"
              class="nav-cta-row"
              :class="{ 'is-hidden': !cta.visible }"
            >
              <span class="nav-cta-drag" title="拖动排序">
                <UiIcon name="drag" :size="13" />
              </span>
              <el-input
                v-model="cta.label.zh"
                size="small"
                class="nav-field nav-field--zh"
                placeholder="中文文案"
              />
              <el-input
                v-model="cta.label.en"
                size="small"
                class="nav-field nav-field--en"
                placeholder="English label"
              />
              <div class="nav-row-path">
                <UiIcon name="link" :size="12" />
                <el-input v-model="cta.href" size="small" placeholder="/contact" />
              </div>
              <div class="nav-cta-variant">
                <span class="nav-cta-dot" :class="`is-${cta.variant}`" />
                <el-select v-model="cta.variant" size="small">
                  <el-option
                    v-for="option in ctaVariantOptions"
                    :key="option.value"
                    :label="option.label"
                    :value="option.value"
                  />
                </el-select>
              </div>
              <div class="nav-row-actions">
                <el-switch
                  v-model="cta.visible"
                  size="small"
                  :title="cta.visible ? '显示中，点击隐藏' : '已隐藏，点击显示'"
                />
                <button class="nav-cta-remove" title="删除按钮" @click="removeCta(index)">
                  <UiIcon name="trash" :size="13" />
                </button>
              </div>
            </div>
          </VueDraggable>
        </div>

        <div v-if="!header.ctas.length" class="nav-empty">
          还没有顶部按钮，点右上角「添加按钮」开始。
        </div>
        <p v-if="ctaIssues.length" class="nav-warning">{{ ctaIssues.join('；') }}</p>
      </section>
      </div>

      <section class="adm-card nav-card">
        <div class="adm-card-head">
          <div>
            <h3>头部设置</h3>
            <p>品牌名称、语言切换与滚动吸顶。</p>
          </div>
        </div>

        <div class="nav-settings">
          <div class="nav-setting-row">
            <div class="nav-setting-label">
              <strong>站点名称</strong>
              <span>显示在左上角品牌区</span>
            </div>
            <div class="nav-pair">
              <el-input v-model="settings.siteName.zh" size="small" placeholder="中文名称" />
              <el-input v-model="settings.siteName.en" size="small" placeholder="English name" />
            </div>
          </div>

          <div class="nav-setting-row">
            <div class="nav-setting-label">
              <strong>公告条</strong>
              <span>顶部导航上方的通知条</span>
            </div>
            <div class="nav-setting-inline">
              <el-switch v-model="header.showAnnouncement" size="small" />
              <span class="nav-setting-hint">关闭后整块隐藏</span>
            </div>
            <template v-if="header.showAnnouncement">
              <div class="nav-pair">
                <el-input
                  v-model="header.announcement.text.zh"
                  size="small"
                  placeholder="公告文案"
                />
                <el-input
                  v-model="header.announcement.text.en"
                  size="small"
                  placeholder="Announcement"
                />
              </div>
              <div class="nav-pair">
                <el-input
                  v-model="header.announcement.highlight.zh"
                  size="small"
                  placeholder="高亮文案"
                />
                <el-input
                  v-model="header.announcement.highlight.en"
                  size="small"
                  placeholder="Highlight"
                />
              </div>
              <div class="nav-pair">
                <el-input
                  v-model="header.announcement.linkLabel.zh"
                  size="small"
                  placeholder="链接文案"
                />
                <el-input
                  v-model="header.announcement.linkLabel.en"
                  size="small"
                  placeholder="Link label"
                />
              </div>
              <el-input
                v-model="header.announcement.linkHref"
                size="small"
                placeholder="/news"
              />
            </template>
          </div>

          <div class="nav-setting-row">
            <div class="nav-setting-label">
              <strong>语言切换</strong>
              <span>中 / EN 切换入口</span>
            </div>
            <el-switch v-model="header.showLanguage" size="small" />
          </div>

          <div class="nav-setting-row">
            <div class="nav-setting-label">
              <strong>滚动吸顶</strong>
              <span>页面向下滚动时固定在顶部</span>
            </div>
            <el-switch v-model="header.sticky" size="small" />
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { VueDraggable } from 'vue-draggable-plus'
import type { HeaderCta, HeaderCtaVariant, Locale, NavItem } from '~/domain/types'
import type { EditorSaveState } from '@/types/admin'
import { usePagesStore } from '@/stores/pages'
import PreviewChrome from '@/preview/PreviewChrome.vue'
import { brandStyle } from '@/utils/brand'
import UiIcon from '~/components/UiIcon.vue'

const pagesStore = usePagesStore()
const { settings } = storeToRefs(pagesStore)

const header = computed(() => settings.value.header)
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
  () => [settings.value.nav, settings.value.header, settings.value.siteName],
  () => {
    saveState.value = 'dirty'
  },
  { deep: true }
)

async function saveNav() {
  saveState.value = 'saving'
  try {
    await pagesStore.saveSettings()
    saveState.value = 'saved'
    ElMessage.success('顶部导航设置已保存')
  } catch (error) {
    saveState.value = 'dirty'
    ElMessage.error(error instanceof Error ? error.message : '保存失败，请稍后重试')
  }
}

const navIssues = computed(() => {
  const issues: string[] = []
  const empty = settings.value.nav.filter((item) => !item.path.trim()).length
  if (empty) issues.push(`${empty} 个菜单项还没有填写跳转路径`)

  const counts = new Map<string, number>()
  for (const item of settings.value.nav) {
    const path = item.path.trim()
    if (!path) continue
    counts.set(path, (counts.get(path) ?? 0) + 1)
  }
  const duplicates = [...counts.entries()].filter(([, count]) => count > 1).map(([path]) => path)
  if (duplicates.length) issues.push(`路径重复：${duplicates.join('、')}`)
  return issues
})

function onReorder(value: NavItem[]) {
  settings.value.nav = value
}

const ctaVariantOptions: Array<{ value: HeaderCtaVariant; label: string }> = [
  { value: 'primary', label: '主按钮' },
  { value: 'secondary', label: '描边按钮' },
  { value: 'ghost', label: '文字按钮' }
]

const visibleCtaCount = computed(() => header.value.ctas.filter((cta) => cta.visible).length)

const ctaIssues = computed(() => {
  const list: string[] = []
  const visible = header.value.ctas.filter((cta) => cta.visible)
  if (visible.length > 3) list.push('可见按钮超过 3 个，移动端可能拥挤')
  if (visible.filter((cta) => cta.variant === 'primary').length > 1) {
    list.push('建议只保留一个主按钮，避免视觉层级混乱')
  }
  const missingPath = visible.filter((cta) => !cta.href.trim()).length
  if (missingPath) list.push(`${missingPath} 个可见按钮还没有填写跳转路径`)
  const missingLabel = visible.filter(
    (cta) => !cta.label.zh.trim() || !cta.label.en.trim()
  ).length
  if (missingLabel) list.push(`${missingLabel} 个可见按钮缺少中文或英文文案`)
  return list
})

function newCta(): HeaderCta {
  return {
    key: `header.cta.custom_${Date.now().toString(36)}`,
    label: { zh: '新按钮', en: 'New button' },
    href: '/contact',
    variant: 'ghost',
    visible: true
  }
}

function reorderCtas(value: HeaderCta[]) {
  settings.value.header.ctas = value
}

function addCta() {
  settings.value.header.ctas.push(newCta())
}

async function removeCta(index: number) {
  const cta = settings.value.header.ctas[index]
  if (!cta) return
  try {
    await ElMessageBox.confirm(
      `确定删除按钮「${cta.label.zh || cta.label.en || cta.href}」吗？`,
      '删除顶部按钮',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
    settings.value.header.ctas.splice(index, 1)
  } catch {
    // cancelled
  }
}

function addItem() {
  const index = settings.value.nav.length + 1
  settings.value.nav.push({
    key: `nav.custom_${Date.now().toString(36)}`,
    path: `/new-page-${index}`,
    zh: '新菜单',
    en: 'New link',
    visible: true
  })
}

async function removeItem(index: number) {
  const item = settings.value.nav[index]
  if (!item) return
  try {
    await ElMessageBox.confirm(
      `确定删除菜单项「${item.zh || item.en || item.path}」吗？`,
      '删除菜单项',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
    settings.value.nav.splice(index, 1)
  } catch {
    // cancelled
  }
}

async function resetNav() {
  try {
    await ElMessageBox.confirm(
      '会把菜单项、站点名称和头部按钮恢复为初始数据，确定继续吗？',
      '恢复默认顶部导航',
      { type: 'warning', confirmButtonText: '恢复默认', cancelButtonText: '取消' }
    )
    pagesStore.resetNavigation()
    saveState.value = 'dirty'
    ElMessage.success('顶部导航已恢复默认，点击「保存设置」生效')
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

/* 编辑器在 <=1400px 会隐藏保存状态，这里保持显示。 */
.nav-panel .save-state {
  display: inline-flex;
}

.nav-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(330px, 0.75fr);
  gap: 16px;
  align-items: start;
}

.nav-preview-card {
  margin-bottom: 16px;
  overflow: hidden;
}

.nav-card {
  padding-bottom: 12px;
}

.nav-card :deep(.adm-card-head) {
  padding: 16px 18px 12px;
}

.nav-card :deep(.adm-card-head p) {
  margin-top: 3px;
  font-size: 10.5px;
  color: var(--adm-muted);
}

.nav-table {
  margin: 0 18px;
  border: 1px solid var(--adm-line);
  border-radius: 10px;
  overflow: hidden;
  background: #fff;
}

.nav-head-row {
  display: grid;
  align-items: stretch;
  gap: 0;
  background: #f7f3ec;
  border-bottom: 1px solid var(--adm-line);
}

.nav-head-row--menu {
  grid-template-columns: 30px minmax(0, 1.05fr) minmax(0, 1.05fr) minmax(0, 0.95fr) 82px;
}

.nav-head-row--cta {
  grid-template-columns: 30px minmax(0, 1.05fr) minmax(0, 1.05fr) minmax(0, 0.95fr) 116px 82px;
}

.nav-head-row > span {
  display: flex;
  align-items: center;
  min-width: 0;
  padding: 9px 10px;
  border-left: 1px solid var(--adm-line-soft);
  font-size: 11.5px;
  font-weight: 600;
  color: var(--adm-ink-soft);
}

.nav-head-row > span:first-child {
  border-left: 0;
}

.nav-head-actions {
  justify-content: flex-end;
}

.nav-rows,
.nav-cta-rows {
  display: flex;
  flex-direction: column;
  min-height: 24px;
  padding: 0;
}

.nav-row,
.nav-cta-row {
  display: grid;
  align-items: stretch;
  gap: 0;
  padding: 0;
  border-bottom: 1px solid var(--adm-line-soft);
  background: #fff;
  transition: background 0.15s ease;
}

.nav-row {
  grid-template-columns: 30px minmax(0, 1.05fr) minmax(0, 1.05fr) minmax(0, 0.95fr) 82px;
}

.nav-cta-row {
  grid-template-columns: 30px minmax(0, 1.05fr) minmax(0, 1.05fr) minmax(0, 0.95fr) 116px 82px;
}

.nav-row > *,
.nav-cta-row > * {
  display: flex;
  align-items: center;
  min-width: 0;
  padding: 6px 10px;
  border-left: 1px solid var(--adm-line-soft);
}

.nav-row > *:first-child,
.nav-cta-row > *:first-child {
  justify-content: center;
  padding: 6px 0;
  border-left: 0;
}

.nav-row:hover,
.nav-cta-row:hover {
  background: #faf8f4;
}

.nav-row:last-child,
.nav-cta-row:last-child {
  border-bottom: 0;
}

.nav-row.is-hidden,
.nav-cta-row.is-hidden {
  opacity: 0.45;
}

.nav-drag,
.nav-cta-drag {
  display: grid;
  place-items: center;
  height: 26px;
  color: var(--adm-muted-2);
  cursor: grab;
}

.nav-drag:active,
.nav-cta-drag:active {
  cursor: grabbing;
}

.nav-field :deep(.el-input__wrapper) {
  padding: 3px 6px;
  background: transparent;
  box-shadow: none;
  border-radius: 6px;
  transition: background 0.15s ease, box-shadow 0.15s ease;
}

.nav-field :deep(.el-input__wrapper:hover) {
  background: rgba(20, 21, 24, 0.045);
}

.nav-field :deep(.el-input__wrapper.is-focus) {
  background: #fff;
  box-shadow: inset 0 0 0 1px var(--adm-line);
}

.nav-field--zh :deep(input) {
  color: var(--adm-ink);
  font-weight: 600;
}

.nav-field--en :deep(input) {
  color: var(--adm-muted);
}

.nav-row-path {
  display: flex;
  align-items: center;
  gap: 5px;
  min-width: 0;
  color: var(--adm-muted-2);
}

.nav-row-path .ui-icon {
  opacity: 0.5;
}

.nav-row-path :deep(.el-input) {
  flex: 1;
  min-width: 0;
}

.nav-row-path :deep(.el-input__wrapper) {
  padding: 3px 4px;
  background: transparent;
  box-shadow: none;
}

.nav-row-path :deep(.el-input__wrapper:hover) {
  background: rgba(20, 21, 24, 0.045);
}

.nav-row-path :deep(.el-input__wrapper.is-focus) {
  background: #fff;
  box-shadow: inset 0 0 0 1px var(--adm-line);
}

.nav-row-path :deep(input) {
  font-family: var(--adm-mono);
  font-size: 10.5px;
  color: var(--adm-muted);
}

.nav-row-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  justify-content: flex-end;
}

.nav-remove,
.nav-cta-remove {
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

.nav-row:hover .nav-remove,
.nav-cta-row:hover .nav-cta-remove {
  opacity: 1;
}

.nav-remove:hover,
.nav-cta-remove:hover {
  background: var(--adm-red-soft);
  color: var(--adm-accent);
}

.nav-main {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.nav-cta-head-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.nav-cta-head-actions .mono {
  font-family: var(--adm-mono);
  font-size: 8.5px;
  color: var(--adm-muted-2);
}

.nav-cta-variant {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.nav-cta-variant :deep(.el-select) {
  width: 100%;
}

.nav-cta-variant :deep(.el-select__wrapper) {
  min-height: 24px;
  padding: 2px 8px;
  border-radius: 999px;
  background: rgba(20, 21, 24, 0.05);
  box-shadow: none;
  font-size: 10.5px;
}

.nav-cta-variant :deep(.el-select__wrapper:hover) {
  background: rgba(20, 21, 24, 0.09);
}

.nav-cta-variant :deep(.el-select__wrapper.is-focused) {
  background: #fff;
  box-shadow: inset 0 0 0 1px var(--adm-line);
}

.nav-cta-dot {
  width: 7px;
  height: 7px;
  flex: none;
  border-radius: 50%;
}

.nav-cta-dot.is-primary {
  background: var(--adm-accent);
}

.nav-cta-dot.is-secondary {
  border: 1.5px solid var(--adm-ink-soft);
  background: transparent;
}

.nav-cta-dot.is-ghost {
  background: var(--adm-muted-2);
}

.nav-warning {
  margin: 12px 18px 0;
  padding: 7px 10px;
  border-radius: 7px;
  background: var(--adm-amber-soft);
  color: var(--adm-amber);
  font-size: 10px;
  line-height: 1.6;
}

.nav-empty {
  margin: 8px 18px;
  padding: 18px;
  border: 1px dashed var(--adm-line);
  border-radius: 8px;
  text-align: center;
  font-size: 11px;
  color: var(--adm-muted-2);
}

.nav-settings {
  padding: 2px 18px 6px;
}

.nav-setting-row {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 9px;
  padding-block: 12px;
  border-top: 1px solid var(--adm-line-soft);
}

.nav-setting-row:first-child {
  border-top: 0;
}

.nav-setting-label {
  flex: none;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-setting-label strong {
  font-size: 11.5px;
  font-weight: 600;
}

.nav-setting-label span {
  font-size: 9.5px;
  color: var(--adm-muted-2);
}

.nav-setting-inline {
  display: flex;
  align-items: center;
  gap: 9px;
  width: 100%;
}

.nav-setting-hint {
  font-size: 9.5px;
  color: var(--adm-muted-2);
}

.nav-pair {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  width: 100%;
}

.nav-preview-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 16px 18px 12px;
  flex-wrap: wrap;
}

.nav-preview-head h3 {
  font-size: 13px;
  font-weight: 600;
}

.nav-preview-head p {
  margin-top: 3px;
  font-size: 10px;
  color: var(--adm-muted-2);
}

.nav-preview-actions {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.nav-segmented {
  display: inline-flex;
  padding: 2px;
  gap: 2px;
  border: 1px solid var(--adm-line);
  border-radius: 7px;
  background: var(--adm-surface);
}

.nav-segmented button {
  border: 0;
  background: transparent;
  padding: 3px 9px;
  border-radius: 5px;
  font-size: 9.5px;
  color: var(--adm-muted);
  cursor: pointer;
}

.nav-segmented button.active {
  background: #fff;
  color: var(--adm-ink);
  box-shadow: 0 1px 2px rgba(20, 21, 24, 0.12);
}

.nav-preview-stage {
  padding: 16px;
  background:
    linear-gradient(45deg, rgba(20, 21, 24, 0.025) 25%, transparent 25%) 0 0 / 16px 16px,
    linear-gradient(-45deg, rgba(20, 21, 24, 0.025) 25%, transparent 25%) 0 8px / 16px 16px,
    linear-gradient(45deg, transparent 75%, rgba(20, 21, 24, 0.025) 75%) 8px -8px / 16px 16px,
    linear-gradient(-45deg, transparent 75%, rgba(20, 21, 24, 0.025) 75%) -8px 0 / 16px 16px,
    var(--adm-bg);
}

.nav-preview-frame {
  width: 100%;
  margin: 0 auto;
  background: var(--paper);
  border: 1px solid var(--adm-line);
  border-radius: 10px;
  overflow: hidden;
  box-shadow: var(--adm-shadow);

  /* 预览站的设计变量与容器，限定在这块预览里，避免影响后台本身。 */
  --paper: #f6f2ea;
  --ink: #141518;
  --ink-soft: #292a2f;
  --muted: #6d6c66;
  --muted-2: #99978e;
  --line: rgba(20, 21, 24, 0.12);
  --line-soft: rgba(20, 21, 24, 0.07);
  --line-strong: rgba(20, 21, 24, 0.22);
  --accent: #e8503a;
  --font-sans: var(--adm-font);
  --font-mono: var(--adm-mono);
}

.nav-preview-frame :deep(.container) {
  width: 100%;
  max-width: 1240px;
  margin: 0 auto;
  padding-inline: clamp(14px, 3vw, 26px);
}

.nav-preview-frame :deep(.mono) {
  font-family: var(--adm-mono);
  letter-spacing: 0.02em;
}

.nav-preview-fill {
  height: 110px;
  display: grid;
  place-items: center;
  background:
    radial-gradient(ellipse 70% 60% at 50% 10%, rgba(232, 80, 58, 0.08), transparent 70%),
    var(--paper);
  border-top: 1px solid var(--line-soft);
  color: rgba(20, 21, 24, 0.18);
  font-size: 8px;
  letter-spacing: 0.16em;
}

.nav-preview-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 18px 14px;
  color: var(--adm-muted-2);
}

.nav-preview-note .mono {
  color: var(--adm-accent);
  font-size: 8px;
  letter-spacing: 0.12em;
  padding-top: 1px;
}

.nav-preview-note p {
  font-size: 9.5px;
  line-height: 1.6;
}

@media (max-width: 1180px) {
  .nav-layout {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 760px) {
  .nav-head-row {
    display: none;
  }
  .nav-table {
    margin: 0 12px;
  }
  .nav-row,
  .nav-cta-row {
    grid-template-columns: 22px minmax(0, 1fr) auto;
    grid-template-areas:
      'drag zh actions'
      'drag en en'
      'drag path path';
    align-items: start;
    gap: 6px;
    padding: 8px;
  }
  .nav-row > *,
  .nav-cta-row > * {
    padding: 0;
    border-left: 0;
  }
  .nav-cta-row {
    grid-template-areas:
      'drag zh actions'
      'drag en en'
      'drag path variant';
  }
  .nav-drag,
  .nav-cta-drag {
    grid-area: drag;
    align-self: start;
  }
  .nav-field--zh {
    grid-area: zh;
  }
  .nav-field--en {
    grid-area: en;
  }
  .nav-row-path {
    grid-area: path;
  }
  .nav-cta-variant {
    grid-area: variant;
    justify-self: end;
    min-width: 110px;
  }
  .nav-row-actions {
    grid-area: actions;
    justify-content: flex-end;
  }
  .nav-cta-head-actions {
    width: 100%;
    justify-content: space-between;
  }
  .nav-setting-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  .nav-pair {
    width: 100%;
  }
}
</style>

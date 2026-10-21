<template>
  <div class="site-panel">
    <section class="panel-intro">
      <div>
        <p class="panel-kicker">SITE SETTINGS</p>
        <h2 class="panel-title">站点配置</h2>
        <p class="panel-sub">
          配置官网的整体视觉与收录策略：品牌 Logo、主色、页面底色、深色区块底色和搜索引擎设置。修改后点击右上角「保存设置」，会同步到页面草稿预览的头部与页脚。
        </p>
      </div>
      <div class="panel-actions">
        <span class="save-state" :class="saveState"><i />{{ saveLabel }}</span>
        <el-button size="small" @click="resetSite">恢复默认</el-button>
        <el-button
          size="small"
          type="primary"
          :disabled="saveState === 'saved'"
          @click="saveSite"
        >
          保存设置
        </el-button>
      </div>
    </section>

    <div class="site-layout">
      <div class="site-main">
        <section class="adm-card site-card">
          <div class="adm-card-head">
            <div>
              <h3>品牌外观</h3>
              <p>头部、页脚、按钮与背景共用这套品牌视觉。</p>
            </div>
          </div>

          <div class="site-body">
            <div class="site-row">
              <div class="site-label">
                <strong>官网 Logo</strong>
                <span>显示在官网头部与页脚；为空时使用内置品牌图标</span>
              </div>
              <div class="logo-field">
                <span class="logo-preview">
                  <img v-if="settings.brand.logo" :src="settings.brand.logo" alt="Logo 预览" />
                  <BrandMark v-else :size="44" />
                </span>
                <div class="logo-actions">
                  <el-button size="small" :loading="logoLoading" @click="pickLogo">
                    {{ settings.brand.logo ? '更换 Logo' : '上传 Logo' }}
                  </el-button>
                  <el-button
                    v-if="settings.brand.logo"
                    size="small"
                    @click="settings.brand.logo = ''"
                  >
                    恢复内置图标
                  </el-button>
                  <span class="logo-hint">{{ logoHint }}</span>
                </div>
              </div>
              <input
                ref="logoInput"
                class="logo-input"
                type="file"
                accept="image/*"
                @change="onLogoPick"
              />
            </div>

            <div class="site-row">
              <div class="site-label">
                <strong>品牌主色</strong>
                <span>用于按钮、链接、标签与强调元素</span>
              </div>
              <div class="accent-field">
                <button
                  v-for="preset in accentPresets"
                  :key="preset.value"
                  type="button"
                  class="accent-swatch"
                  :class="{ active: isColor(normalizedAccent, preset.value) }"
                  :style="{ background: preset.value }"
                  :title="preset.label"
                  @click="settings.brand.accent = preset.value"
                >
                  <UiIcon
                    v-if="isColor(normalizedAccent, preset.value)"
                    name="check"
                    :size="12"
                  />
                </button>
                <el-color-picker
                  v-model="settings.brand.accent"
                  size="small"
                  color-format="hex"
                  :predefine="accentValues"
                />
                <span class="accent-value mono">{{ normalizedAccent }}</span>
              </div>
            </div>

            <div class="site-row">
              <div class="site-label">
                <strong>页面底色</strong>
                <span>页面、头部与浅色区块的整体背景</span>
              </div>
              <div class="accent-field">
                <button
                  v-for="preset in backgroundPresets"
                  :key="preset.value"
                  type="button"
                  class="accent-swatch accent-swatch--light"
                  :class="{ active: isColor(normalizedBackground, preset.value) }"
                  :style="{ background: preset.value }"
                  :title="preset.label"
                  @click="settings.brand.background = preset.value"
                >
                  <UiIcon
                    v-if="isColor(normalizedBackground, preset.value)"
                    name="check"
                    :size="12"
                  />
                </button>
                <el-color-picker
                  v-model="settings.brand.background"
                  size="small"
                  color-format="hex"
                />
                <span class="accent-value mono">{{ normalizedBackground }}</span>
              </div>
            </div>

            <div class="site-row">
              <div class="site-label">
                <strong>深色区块底色</strong>
                <span>页脚、深色横幅与反白区块使用这个底色</span>
              </div>
              <div class="accent-field">
                <button
                  v-for="preset in darkPresets"
                  :key="preset.value"
                  type="button"
                  class="accent-swatch"
                  :class="{ active: isColor(normalizedDark, preset.value) }"
                  :style="{ background: preset.value }"
                  :title="preset.label"
                  @click="settings.brand.darkBackground = preset.value"
                >
                  <UiIcon
                    v-if="isColor(normalizedDark, preset.value)"
                    name="check"
                    :size="12"
                  />
                </button>
                <el-color-picker
                  v-model="settings.brand.darkBackground"
                  size="small"
                  color-format="hex"
                />
                <span class="accent-value mono">{{ normalizedDark }}</span>
              </div>
            </div>

            <div class="site-row">
              <div class="site-label">
                <strong>搜索引擎收录</strong>
                <span>关闭时全站输出 noindex，适合预览或未上线站点</span>
              </div>
              <el-switch v-model="settings.allowIndex" size="small" />
            </div>
          </div>
        </section>

        <section class="adm-card site-card">
          <div class="adm-card-head">
            <div>
              <h3>站点信息</h3>
              <p>这些信息在顶部导航与底部导航中维护，避免重复配置。</p>
            </div>
          </div>
          <div class="site-body site-body--info">
            <div class="info-item">
              <span class="info-label">站点名称</span>
              <strong>{{ settings.siteName.zh }} / {{ settings.siteName.en }}</strong>
              <span class="info-hint">前往「顶部导航 → 头部设置」修改</span>
            </div>
            <div class="info-item">
              <span class="info-label">联系方式</span>
              <strong>{{ settings.email }} · {{ settings.phone }}</strong>
              <span class="info-hint">前往「底部导航 → 页脚信息」修改</span>
            </div>
          </div>
        </section>
      </div>

      <aside class="adm-card site-preview">
        <div class="site-preview-head">
          <div>
            <h3>品牌预览</h3>
            <p>Logo、主色与头部按钮的实时效果</p>
          </div>
          <div class="site-preview-actions">
            <div class="site-segmented">
              <button :class="{ active: previewLocale === 'zh' }" @click="previewLocale = 'zh'">
                中文
              </button>
              <button :class="{ active: previewLocale === 'en' }" @click="previewLocale = 'en'">
                EN
              </button>
            </div>
            <div class="site-segmented">
              <button :class="{ active: device === 'desktop' }" @click="device = 'desktop'">
                桌面
              </button>
              <button :class="{ active: device === 'mobile' }" @click="device = 'mobile'">
                手机
              </button>
            </div>
          </div>
        </div>

        <div class="site-preview-stage">
          <div class="site-preview-frame" :style="[brandVars, { maxWidth: previewWidth }]">
            <PreviewChrome :settings="settings" :locale="previewLocale" />
            <div class="site-preview-fill">
              <span class="mono sample-kicker">BRAND PREVIEW</span>
              <div class="sample-actions">
                <span class="sample-btn">预约演示 ↗</span>
                <span class="sample-link">了解更多 →</span>
              </div>
              <div class="sample-dark">
                <span class="mono">DARK SURFACE</span>
                <span>页脚与深色区块</span>
              </div>
            </div>
          </div>
        </div>

        <div class="site-preview-note">
          <span class="mono">LIVE</span>
          <p>修改 Logo、主色与背景色会立即反映在预览中；点击「保存设置」后写入本地演示数据。</p>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Locale } from '~/domain/types'
import type { EditorSaveState } from '@/types/admin'
import { usePagesStore } from '@/stores/pages'
import { brandStyle, normalizeColor } from '@/utils/brand'
import {
  estimateDataUrlBytes,
  fileToImageDataUrl,
  formatBytes
} from '@/utils/image'
import PreviewChrome from '@/preview/PreviewChrome.vue'
import BrandMark from '~/components/BrandMark.vue'
import UiIcon from '~/components/UiIcon.vue'

const pagesStore = usePagesStore()
const { settings } = storeToRefs(pagesStore)

const previewLocale = ref<Locale>('zh')
const device = ref<'desktop' | 'mobile'>('desktop')
const previewWidth = computed(() => (device.value === 'mobile' ? '420px' : '100%'))
const brandVars = computed(() => brandStyle(settings.value.brand))
const normalizedAccent = computed(() => normalizeColor(settings.value.brand.accent))
const normalizedBackground = computed(() =>
  normalizeColor(settings.value.brand.background, '#F6F2EA')
)
const normalizedDark = computed(() =>
  normalizeColor(settings.value.brand.darkBackground, '#141518')
)

const accentPresets = [
  { label: '朱砂红', value: '#E8503A' },
  { label: '松柏绿', value: '#2F6B55' },
  { label: '靛青', value: '#2E5AAC' },
  { label: '金棕', value: '#B07A2A' },
  { label: '紫罗兰', value: '#6B4FA0' }
]
const accentValues = accentPresets.map((preset) => preset.value)

const backgroundPresets = [
  { label: '暖纸', value: '#F6F2EA' },
  { label: '纯白', value: '#FFFFFF' },
  { label: '米沙', value: '#EFE7DA' },
  { label: '浅灰', value: '#F2F2EF' },
  { label: '雾蓝', value: '#EDF2F5' }
]
const darkPresets = [
  { label: '墨黑', value: '#141518' },
  { label: '深棕', value: '#241A16' },
  { label: '深绿', value: '#14251E' },
  { label: '深蓝', value: '#141D2E' }
]

const saveState = ref<EditorSaveState>('saved')
const saveLabel = computed(() => {
  if (saveState.value === 'dirty') return '有未保存修改'
  if (saveState.value === 'saving') return '正在保存…'
  return `已保存 ${pagesStore.lastSavedAt}`
})

watch(
  () => [settings.value.brand, settings.value.allowIndex],
  () => {
    saveState.value = 'dirty'
  },
  { deep: true }
)

function isColor(current: string, value: string): boolean {
  return current === normalizeColor(value)
}

const logoInput = ref<HTMLInputElement | null>(null)
const logoLoading = ref(false)
const logoHint = computed(() => {
  if (!settings.value.brand.logo) return 'PNG / SVG / WebP，会自动压缩到 1600px 以内'
  return `当前 ${formatBytes(estimateDataUrlBytes(settings.value.brand.logo))}`
})

function pickLogo() {
  logoInput.value?.click()
}

async function onLogoPick(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  logoLoading.value = true
  try {
    const dataUrl = await fileToImageDataUrl(file)
    settings.value.brand.logo = dataUrl
    ElMessage.success(`Logo 已更新（${formatBytes(estimateDataUrlBytes(dataUrl))}）`)
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '图片处理失败')
  } finally {
    logoLoading.value = false
  }
}

function saveSite() {
  saveState.value = 'saving'
  pagesStore.saveSettings()
  saveState.value = 'saved'
  ElMessage.success('站点配置已保存')
}

async function resetSite() {
  try {
    await ElMessageBox.confirm(
      '会把 Logo、品牌主色、页面底色、深色区块底色与搜索引擎收录恢复为初始数据，确定继续吗？',
      '恢复默认站点配置',
      { type: 'warning', confirmButtonText: '恢复默认', cancelButtonText: '取消' }
    )
    pagesStore.resetSite()
    saveState.value = 'dirty'
    ElMessage.success('站点配置已恢复默认，点击「保存设置」生效')
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

.site-panel .save-state {
  display: inline-flex;
}

.site-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(330px, 0.75fr);
  gap: 16px;
  align-items: start;
}

.site-main {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.site-card {
  padding-bottom: 14px;
}

.site-card :deep(.adm-card-head) {
  padding: 16px 18px 12px;
}

.site-card :deep(.adm-card-head p) {
  margin-top: 3px;
  font-size: 10.5px;
  color: var(--adm-muted);
}

.site-body {
  display: flex;
  flex-direction: column;
  padding: 0 18px;
}

.site-row {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
  padding-block: 14px;
  border-top: 1px solid var(--adm-line-soft);
}

.site-row:first-child {
  border-top: 0;
}

.site-label {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.site-label strong {
  font-size: 11.5px;
  font-weight: 600;
}

.site-label span {
  font-size: 9.5px;
  color: var(--adm-muted-2);
}

.logo-field {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
}

.logo-preview {
  width: 68px;
  height: 68px;
  flex: none;
  display: grid;
  place-items: center;
  border: 1px solid var(--adm-line);
  border-radius: 12px;
  background:
    linear-gradient(45deg, rgba(20, 21, 24, 0.03) 25%, transparent 25%) 0 0 / 12px 12px,
    linear-gradient(-45deg, rgba(20, 21, 24, 0.03) 25%, transparent 25%) 0 6px / 12px 12px,
    #fff;
  overflow: hidden;
}

.logo-preview img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: 10px;
}

.logo-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  min-width: 0;
}

.logo-hint {
  width: 100%;
  font-size: 9.5px;
  color: var(--adm-muted-2);
}

.logo-input {
  display: none;
}

.accent-field {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.accent-swatch {
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  border: 2px solid transparent;
  border-radius: 50%;
  color: #fff;
  cursor: pointer;
  box-shadow: inset 0 0 0 1px rgba(20, 21, 24, 0.12);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.accent-swatch:hover {
  transform: scale(1.08);
}

.accent-swatch.active {
  box-shadow: 0 0 0 2px #fff, 0 0 0 4px var(--adm-ink);
}

.accent-swatch--light {
  color: var(--adm-ink);
}

.accent-value {
  font-size: 9.5px;
  letter-spacing: 0.04em;
  color: var(--adm-muted);
}

.site-body--info {
  gap: 0;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 12px 0;
  border-top: 1px solid var(--adm-line-soft);
}

.info-item:first-child {
  border-top: 0;
  padding-top: 2px;
}

.info-label {
  font-family: var(--adm-mono);
  font-size: 8px;
  letter-spacing: 0.12em;
  color: var(--adm-muted-2);
}

.info-item strong {
  font-size: 12px;
  font-weight: 600;
}

.info-hint {
  font-size: 9.5px;
  color: var(--adm-muted-2);
}

.site-preview {
  overflow: hidden;
}

.site-preview-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 16px 18px 12px;
  flex-wrap: wrap;
}

.site-preview-head h3 {
  font-size: 13px;
  font-weight: 600;
}

.site-preview-head p {
  margin-top: 3px;
  font-size: 10px;
  color: var(--adm-muted-2);
}

.site-preview-actions {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.site-segmented {
  display: inline-flex;
  padding: 2px;
  gap: 2px;
  border: 1px solid var(--adm-line);
  border-radius: 7px;
  background: var(--adm-surface);
}

.site-segmented button {
  border: 0;
  background: transparent;
  padding: 3px 9px;
  border-radius: 5px;
  font-size: 9.5px;
  color: var(--adm-muted);
  cursor: pointer;
}

.site-segmented button.active {
  background: #fff;
  color: var(--adm-ink);
  box-shadow: 0 1px 2px rgba(20, 21, 24, 0.12);
}

.site-preview-stage {
  padding: 16px;
  background:
    linear-gradient(45deg, rgba(20, 21, 24, 0.025) 25%, transparent 25%) 0 0 / 16px 16px,
    linear-gradient(-45deg, rgba(20, 21, 24, 0.025) 25%, transparent 25%) 0 8px / 16px 16px,
    var(--adm-bg);
}

.site-preview-frame {
  width: 100%;
  margin: 0 auto;
  background: var(--paper);
  border: 1px solid var(--adm-line);
  border-radius: 10px;
  overflow: hidden;
  box-shadow: var(--adm-shadow);

  /* 预览站字体与容器，限定在预览区内部。 */
  --paper: #f6f2ea;
  --ink: #141518;
  --ink-soft: #292a2f;
  --muted: #6d6c66;
  --muted-2: #99978e;
  --line: rgba(20, 21, 24, 0.12);
  --line-soft: rgba(20, 21, 24, 0.07);
  --line-strong: rgba(20, 21, 24, 0.22);
  --font-sans: var(--adm-font);
  --font-mono: var(--adm-mono);
}

.site-preview-frame :deep(.container) {
  width: 100%;
  max-width: 1240px;
  margin: 0 auto;
  padding-inline: clamp(14px, 3vw, 26px);
}

.site-preview-frame :deep(.mono) {
  font-family: var(--adm-mono);
  letter-spacing: 0.02em;
}

.site-preview-fill {
  min-height: 170px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  gap: 16px;
  padding: clamp(22px, 4vw, 38px);
  background:
    radial-gradient(
      ellipse 70% 60% at 50% 10%,
      color-mix(in srgb, var(--accent) 8%, transparent),
      transparent 70%
    ),
    var(--paper);
  border-top: 1px solid var(--line-soft);
}

.sample-kicker {
  font-size: 8px;
  letter-spacing: 0.16em;
  color: var(--muted-2);
}

.sample-actions {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
}

.sample-btn {
  display: inline-flex;
  align-items: center;
  height: 36px;
  padding: 0 16px;
  border-radius: 999px;
  background: var(--accent);
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  box-shadow: 0 8px 20px color-mix(in srgb, var(--accent) 28%, transparent);
}

.sample-link {
  font-size: 11px;
  font-weight: 600;
  color: var(--accent);
}

.sample-dark {
  width: 100%;
  margin-top: 2px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 10px;
  background: var(--ink);
  color: var(--paper);
  font-size: 10px;
}

.sample-dark .mono {
  font-size: 8px;
  letter-spacing: 0.12em;
  opacity: 0.55;
}

.site-preview-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 18px 14px;
  color: var(--adm-muted-2);
}

.site-preview-note .mono {
  color: var(--adm-accent);
  font-size: 8px;
  letter-spacing: 0.12em;
  padding-top: 1px;
}

.site-preview-note p {
  font-size: 9.5px;
  line-height: 1.6;
}

@media (max-width: 1180px) {
  .site-layout {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 760px) {
  .logo-field {
    align-items: flex-start;
  }
  .site-preview-actions {
    width: 100%;
  }
}
</style>

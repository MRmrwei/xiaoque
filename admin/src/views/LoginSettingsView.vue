<template>
  <div class="adm-page">
    <router-link to="/settings" class="settings-back">← 基础设置</router-link>

    <header class="adm-page-head">
      <div>
        <p class="adm-kicker">BASIC SETTINGS · SIGN IN</p>
        <h1 class="adm-title">后台登录</h1>
        <p class="adm-sub">
          配置后台登录方式与登录页展示内容。修改后点击「保存设置」生效。
        </p>
      </div>
      <div class="adm-actions">
        <el-button @click="restoreDefaults">恢复默认</el-button>
        <el-button type="primary" @click="save">保存设置</el-button>
      </div>
    </header>

    <div class="settings-grid">
      <section class="adm-card settings-card">
        <div class="adm-card-head">
          <div>
            <h3>后台登录管理</h3>
            <p>控制登录页的安全校验与展示文案。</p>
          </div>
        </div>

        <div class="settings-body">
          <div class="setting-row">
            <div class="setting-text">
              <strong>图形验证码</strong>
              <span>开启后，登录时必须输入图片中的 4 位字符。</span>
            </div>
            <el-switch v-model="form.captchaEnabled" />
          </div>

          <div class="settings-divider" />

          <p class="settings-label">登录页展示区</p>

          <div class="background-setting">
            <span class="background-preview" :style="brandPreviewStyle">
              <span v-if="!form.backgroundImage" class="background-fallback">默认黑色底</span>
            </span>
            <div class="background-actions">
              <div class="background-buttons">
                <el-button size="small" :loading="backgroundLoading" @click="pickBackground">
                  {{ form.backgroundImage ? '更换底图' : '上传底图' }}
                </el-button>
                <el-button
                  v-if="form.backgroundImage"
                  size="small"
                  @click="form.backgroundImage = ''"
                >
                  使用黑色底
                </el-button>
              </div>
              <span class="background-hint">{{ backgroundHint }}</span>
            </div>
            <input
              ref="backgroundInput"
              class="background-input"
              type="file"
              accept="image/*"
              @change="onBackgroundPick"
            />
          </div>

          <label class="setting-field">
            <span>品牌名称</span>
            <el-input
              v-model="form.brandName"
              maxlength="20"
              show-word-limit
              placeholder="小鹊云"
            />
          </label>

          <label class="setting-field">
            <span>名称下方小字</span>
            <el-input
              v-model="form.brandSubName"
              maxlength="30"
              show-word-limit
              placeholder="SITE STUDIO"
            />
          </label>

          <label class="setting-field">
            <span>眉标</span>
            <el-input v-model="form.brandKicker" maxlength="40" show-word-limit />
          </label>

          <label class="setting-field">
            <span>主标题</span>
            <el-input
              v-model="form.brandTitle"
              type="textarea"
              :rows="2"
              maxlength="60"
              show-word-limit
              placeholder="支持换行"
            />
            <em>支持换行，展示在登录页左侧。</em>
          </label>

          <label class="setting-field">
            <span>说明文字</span>
            <el-input
              v-model="form.brandSubtitle"
              type="textarea"
              :rows="3"
              maxlength="160"
              show-word-limit
            />
          </label>

          <p class="settings-note">
            演示版验证码由前端生成与校验；接后端后改为服务端验证码，登录接口也需要增加验证码参数。
          </p>
        </div>
      </section>

      <aside class="settings-preview">
        <p class="preview-label">登录页预览</p>

        <div class="preview-brand" :style="brandPreviewStyle">
          <div class="preview-brand-head">
            <span>
              {{ form.brandName || '小鹊云' }}
              <em>{{ form.brandSubName || 'SITE STUDIO' }}</em>
            </span>
          </div>

          <div class="preview-copy">
            <p class="mono">{{ form.brandKicker || 'MAGPIE CLOUD' }}</p>
            <h2>{{ form.brandTitle || '把官网\n编排成作品。' }}</h2>
            <p>{{ form.brandSubtitle }}</p>
          </div>

          <div class="preview-foot">
            <span>LOCAL DEMO</span>
            <span>ADMIN v1</span>
          </div>
        </div>

        <div class="preview-form">
          <span class="preview-input">用户名</span>
          <span class="preview-input">密码</span>
          <span v-if="form.captchaEnabled" class="preview-captcha">
            <span class="preview-input">验证码</span>
            <span class="preview-code">AB7K</span>
          </span>
          <span class="preview-button">登录</span>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { usePagesStore } from '@/stores/pages'
import { seedLoginSettings } from '@/data/seed/seed'
import type { LoginSettings } from '@/types/admin'
import { loginBrandBackground } from '@/utils/loginBackground'
import {
  estimateDataUrlBytes,
  fileToImageDataUrl,
  formatBytes
} from '@/utils/image'

const pagesStore = usePagesStore()
const form = reactive<LoginSettings>({ ...pagesStore.loginSettings })

const backgroundInput = ref<HTMLInputElement | null>(null)
const backgroundLoading = ref(false)
const brandPreviewStyle = computed(() => loginBrandBackground(form.backgroundImage))
const backgroundHint = computed(() => {
  if (!form.backgroundImage) return '未设置时使用黑色底，建议上传横向图片'
  return `当前 ${formatBytes(estimateDataUrlBytes(form.backgroundImage))}`
})

function pickBackground(): void {
  backgroundInput.value?.click()
}

async function onBackgroundPick(event: Event): Promise<void> {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  backgroundLoading.value = true
  try {
    form.backgroundImage = await fileToImageDataUrl(file)
    ElMessage.success(
      `底图已更新（${formatBytes(estimateDataUrlBytes(form.backgroundImage))}）`
    )
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '图片处理失败')
  } finally {
    backgroundLoading.value = false
  }
}

function save(): void {
  pagesStore.saveLoginSettings({ ...form })
  ElMessage.success('基础设置已保存')
}

function restoreDefaults(): void {
  Object.assign(form, seedLoginSettings)
  ElMessage.info('已恢复默认值，点击「保存设置」生效')
}
</script>

<style scoped>
.settings-back {
  display: inline-flex;
  align-items: center;
  margin-bottom: 14px;
  font-size: 12px;
  color: var(--adm-muted);
  transition: color 0.18s ease;
}

.settings-back:hover {
  color: var(--adm-accent-deep);
}

.settings-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(300px, 0.65fr);
  gap: 18px;
  align-items: start;
}

.settings-card {
  overflow: visible;
}

.settings-card .adm-card-head > div p {
  margin-top: 3px;
  font-size: 11px;
  color: var(--adm-muted);
}

.settings-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.setting-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.setting-text strong {
  font-size: 13px;
  font-weight: 600;
}

.setting-text span {
  font-size: 11.5px;
  line-height: 1.6;
  color: var(--adm-muted);
}

.settings-divider {
  border-top: 1px solid var(--adm-line-soft);
}

.settings-label {
  font-family: var(--adm-mono);
  font-size: 8.5px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--adm-muted-2);
}

.background-setting {
  display: flex;
  align-items: center;
  gap: 14px;
}

.background-preview {
  width: 96px;
  height: 60px;
  flex: none;
  display: grid;
  place-items: center;
  border: 1px solid var(--adm-line);
  border-radius: 10px;
  background: #141518;
  background-size: cover;
  background-position: center;
  overflow: hidden;
}

.background-fallback {
  font-size: 11px;
  color: rgba(246, 242, 234, 0.55);
}

.background-actions {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.background-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.background-hint {
  font-size: 11px;
  color: var(--adm-muted-2);
}

.background-input {
  display: none;
}

.setting-field {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.setting-field > span {
  font-size: 11.5px;
  font-weight: 600;
  color: var(--adm-ink-soft);
}

.setting-field em {
  font-style: normal;
  font-size: 10.5px;
  color: var(--adm-muted-2);
}

.setting-field :deep(.el-input__wrapper),
.setting-field :deep(.el-textarea__inner) {
  background: #faf8f4;
  box-shadow: 0 0 0 1px var(--adm-line) inset;
  border-radius: 9px;
}

.setting-field :deep(.el-input__wrapper.is-focus),
.setting-field :deep(.el-textarea__inner:focus) {
  background: #fff;
  box-shadow: 0 0 0 1px var(--adm-accent) inset;
}

.settings-note {
  padding: 11px 13px;
  border: 1px solid var(--adm-line-soft);
  border-radius: 9px;
  background: #faf8f4;
  font-size: 10.5px;
  line-height: 1.65;
  color: var(--adm-muted);
}

.settings-preview {
  position: sticky;
  top: 20px;
}

.preview-label {
  margin-bottom: 9px;
  font-family: var(--adm-mono);
  font-size: 8.5px;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--adm-muted-2);
}

.preview-brand {
  display: flex;
  flex-direction: column;
  gap: 18px;
  min-height: 292px;
  padding: 22px;
  border-radius: 14px;
  background: #141518;
  color: #f6f2ea;
  overflow: hidden;
}

.preview-brand-head {
  display: flex;
  align-items: center;
  gap: 9px;
}

.preview-brand-head span {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 11px;
  font-weight: 600;
  line-height: 1;
}

.preview-brand-head em {
  font-family: var(--adm-mono);
  font-style: normal;
  font-size: 6px;
  letter-spacing: 0.14em;
  color: rgba(246, 242, 234, 0.34);
}

.preview-copy {
  margin-top: auto;
}

.preview-copy .mono {
  font-family: var(--adm-mono);
  font-size: 7px;
  letter-spacing: 0.14em;
  color: rgba(246, 242, 234, 0.38);
}

.preview-copy h2 {
  margin-top: 8px;
  font-size: 27px;
  font-weight: 600;
  letter-spacing: -0.05em;
  line-height: 1.1;
  white-space: pre-line;
  word-break: break-word;
}

.preview-copy p:last-child {
  margin-top: 10px;
  font-size: 10.5px;
  line-height: 1.65;
  color: rgba(246, 242, 234, 0.5);
  word-break: break-word;
}

.preview-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 12px;
  border-top: 1px solid rgba(246, 242, 234, 0.1);
  font-family: var(--adm-mono);
  font-size: 6.5px;
  letter-spacing: 0.14em;
  color: rgba(246, 242, 234, 0.32);
}

.preview-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 10px;
  padding: 14px;
  border: 1px solid var(--adm-line);
  border-radius: 12px;
  background: var(--adm-card);
}

.preview-input {
  display: flex;
  align-items: center;
  height: 34px;
  padding: 0 10px;
  border: 1px solid var(--adm-line);
  border-radius: 7px;
  background: #faf8f4;
  font-size: 10.5px;
  color: var(--adm-muted-2);
}

.preview-captcha {
  display: flex;
  gap: 7px;
}

.preview-captcha .preview-input {
  flex: 1;
  min-width: 0;
}

.preview-code {
  width: 82px;
  height: 34px;
  flex: none;
  display: grid;
  place-items: center;
  border: 1px solid var(--adm-line);
  border-radius: 7px;
  background: #f1ede5;
  font-family: var(--adm-mono);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  color: var(--adm-accent-deep);
}

.preview-button {
  display: grid;
  place-items: center;
  height: 36px;
  margin-top: 2px;
  border-radius: 7px;
  background: var(--adm-ink);
  color: #fff;
  font-size: 11.5px;
  font-weight: 600;
}

@media (max-width: 1080px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }

  .settings-preview {
    position: static;
    max-width: 420px;
  }
}

@media (max-width: 600px) {
  .settings-body {
    padding: 16px;
  }

  .setting-row {
    align-items: flex-start;
  }
}
</style>

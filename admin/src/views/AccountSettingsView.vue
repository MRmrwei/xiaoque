<template>
  <div class="adm-page">
    <router-link to="/pages" class="settings-back">← 返回工作台</router-link>

    <header class="adm-page-head">
      <div>
        <p class="adm-kicker">ACCOUNT &amp; SECURITY</p>
        <h1 class="adm-title">管理员账号</h1>
        <p class="adm-sub">
          维护当前管理员的显示资料与登录密码。演示版数据保存在当前浏览器，不会写入线上账号。
        </p>
      </div>
      <div class="adm-actions">
        <el-button @click="restoreProfile">恢复默认资料</el-button>
        <el-button type="primary" @click="submitProfile">保存资料</el-button>
      </div>
    </header>

    <div class="account-grid">
      <section class="adm-card account-card">
        <div class="adm-card-head">
          <div>
            <h3>管理员信息</h3>
            <p>登录账号不可修改，显示名称会同步到后台侧边栏。</p>
          </div>
        </div>

        <div class="account-body">
          <div class="account-identity">
            <span class="account-avatar" :class="{ 'has-image': profileForm.avatar }">
              <img v-if="profileForm.avatar" :src="profileForm.avatar" alt="管理员头像" />
              <template v-else>{{ avatarText }}</template>
            </span>
            <div class="account-identity-text">
              <strong>{{ profileForm.name.trim() || currentProfile.username }}</strong>
              <span>@{{ currentProfile.username }}</span>
            </div>
            <el-tag type="danger" size="small" effect="plain">超级管理员</el-tag>
          </div>

          <div class="account-avatar-actions">
            <el-button size="small" :loading="avatarLoading" @click="pickAvatar">
              {{ profileForm.avatar ? '更换头像' : '上传头像' }}
            </el-button>
            <el-button
              v-if="profileForm.avatar"
              size="small"
              @click="profileForm.avatar = ''"
            >
              使用文字头像
            </el-button>
            <span class="account-avatar-hint">{{ avatarHint }}</span>
            <input
              ref="avatarInput"
              class="account-avatar-input"
              type="file"
              accept="image/*"
              @change="onAvatarPick"
            />
          </div>

          <label class="account-field">
            <span>登录账号</span>
            <el-input :model-value="currentProfile.username" disabled />
            <em>登录账号创建后不可修改。</em>
          </label>

          <label class="account-field">
            <span>显示名称 <i class="req">*</i></span>
            <el-input
              v-model="profileForm.name"
              maxlength="100"
              show-word-limit
              placeholder="例如：小阙"
            />
          </label>

          <label class="account-field">
            <span>邮箱</span>
            <el-input
              v-model="profileForm.email"
              maxlength="255"
              placeholder="用于接收通知，可留空"
            />
          </label>

          <label class="account-field">
            <span>手机号</span>
            <el-input v-model="profileForm.mobile" maxlength="32" placeholder="可留空" />
          </label>

          <p class="account-note">
            登录已接 Go API；资料修改仍保存在 <span class="mono">localStorage</span>，尚未写回
            <span class="mono">xqyun_users</span>。
          </p>
        </div>
      </section>

      <section class="adm-card account-card">
        <div class="adm-card-head">
          <div>
            <h3>修改密码</h3>
            <p>当前为本地演示改密，不会修改 Go API 的登录密码。</p>
          </div>
        </div>

        <div class="account-body">
          <label class="account-field">
            <span>当前密码 <i class="req">*</i></span>
            <el-input
              v-model="passwordForm.current"
              type="password"
              show-password
              autocomplete="current-password"
              placeholder="请输入当前密码"
            />
          </label>

          <label class="account-field">
            <span>新密码 <i class="req">*</i></span>
            <el-input
              v-model="passwordForm.next"
              type="password"
              show-password
              autocomplete="new-password"
              placeholder="至少 8 个字符"
            />
            <em>至少 8 个字符；bcrypt 上限为 72 字节，中文约 24 个字符。</em>
          </label>

          <label class="account-field">
            <span>确认新密码 <i class="req">*</i></span>
            <el-input
              v-model="passwordForm.confirm"
              type="password"
              show-password
              autocomplete="new-password"
              placeholder="再次输入新密码"
              @keyup.enter="submitPassword"
            />
          </label>

          <div class="account-actions">
            <el-button type="primary" :loading="savingPassword" @click="submitPassword">
              修改密码
            </el-button>
          </div>

          <p class="account-note">
            登录已接 Go API；改密接口尚未接入，当前修改仅保存在本地演示数据，不会修改后端密码。
          </p>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { AdminProfile } from '@/types/admin'
import {
  changePassword,
  demoCredentials,
  getProfile,
  saveProfile,
  validateNewPassword
} from '@/services/auth'
import {
  estimateDataUrlBytes,
  fileToImageDataUrl,
  formatBytes
} from '@/utils/image'

const currentProfile = ref<AdminProfile>(getProfile())
const profileForm = reactive({
  name: currentProfile.value.name,
  email: currentProfile.value.email,
  mobile: currentProfile.value.mobile,
  avatar: currentProfile.value.avatar
})
const passwordForm = reactive({
  current: '',
  next: '',
  confirm: ''
})
const savingPassword = ref(false)
const avatarInput = ref<HTMLInputElement | null>(null)
const avatarLoading = ref(false)

const avatarText = computed(() => {
  const name = profileForm.name.trim() || currentProfile.value.username
  return name.charAt(0).toUpperCase() || '管'
})

const avatarHint = computed(() => {
  if (!profileForm.avatar) return '不上传时保留文字头像'
  if (!profileForm.avatar.startsWith('data:')) return '自定义头像'
  return `当前 ${formatBytes(estimateDataUrlBytes(profileForm.avatar))}，JPG / PNG / WebP 会压缩到 320px 以内`
})

function pickAvatar(): void {
  avatarInput.value?.click()
}

async function onAvatarPick(event: Event): Promise<void> {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  avatarLoading.value = true
  try {
    profileForm.avatar = await fileToImageDataUrl(file, 320)
    ElMessage.success('头像已更新，点击「保存资料」生效')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '头像处理失败')
  } finally {
    avatarLoading.value = false
  }
}

function submitProfile(): void {
  const name = profileForm.name.trim()
  if (!name) {
    ElMessage.error('请输入显示名称')
    return
  }

  const email = profileForm.email.trim()
  if (email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    ElMessage.error('邮箱格式不正确')
    return
  }

  const mobile = profileForm.mobile.trim()
  if (mobile && !/^[0-9+\-()\s]{6,32}$/.test(mobile)) {
    ElMessage.error('手机号格式不正确')
    return
  }

  try {
    const saved = saveProfile({ name, email, mobile, avatar: profileForm.avatar })
    currentProfile.value = saved
    Object.assign(profileForm, {
      name: saved.name,
      email: saved.email,
      mobile: saved.mobile,
      avatar: saved.avatar
    })
    ElMessage.success('管理员资料已保存')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '资料保存失败')
  }
}

async function restoreProfile(): Promise<void> {
  try {
    await ElMessageBox.confirm(
      '恢复默认资料只会重置显示名称、邮箱、手机号和头像，不会修改登录密码。',
      '恢复默认资料',
      { type: 'warning', confirmButtonText: '恢复默认', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  Object.assign(profileForm, {
    name: demoCredentials.name,
    email: '',
    mobile: '',
    avatar: ''
  })
  ElMessage.info('已恢复默认资料，点击「保存资料」生效')
}

async function submitPassword(): Promise<void> {
  if (savingPassword.value) return

  if (!passwordForm.current) {
    ElMessage.error('请输入当前密码')
    return
  }

  const validation = validateNewPassword(passwordForm.next)
  if (validation) {
    ElMessage.error(validation)
    return
  }

  if (passwordForm.next !== passwordForm.confirm) {
    ElMessage.error('两次输入的新密码不一致')
    return
  }

  if (passwordForm.next === passwordForm.current) {
    ElMessage.error('新密码不能与当前密码相同')
    return
  }

  savingPassword.value = true
  try {
    const result = await changePassword(passwordForm.current, passwordForm.next)
    if (!result.ok) {
      ElMessage.error(result.message)
      return
    }
    Object.assign(passwordForm, { current: '', next: '', confirm: '' })
    ElMessage.success('密码已修改，下次登录请使用新密码')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '密码修改失败')
  } finally {
    savingPassword.value = false
  }
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

.account-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  align-items: start;
  max-width: 1160px;
}

.account-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 18px;
}

.account-identity {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  border: 1px solid var(--adm-line-soft);
  border-radius: 11px;
  background: var(--adm-surface);
}

.account-avatar {
  width: 46px;
  height: 46px;
  flex: none;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: var(--adm-accent);
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  overflow: hidden;
}

.account-avatar.has-image {
  background: var(--adm-surface);
  border: 1px solid var(--adm-line);
}

.account-avatar img {
  display: block;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.account-avatar-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: -4px;
}

.account-avatar-hint {
  font-size: 11px;
  line-height: 1.5;
  color: var(--adm-muted-2);
}

.account-avatar-input {
  display: none;
}

.account-identity-text {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.account-identity-text strong {
  font-size: 14px;
  font-weight: 600;
}

.account-identity-text span {
  font-family: var(--adm-mono);
  font-size: 11.5px;
  color: var(--adm-muted);
}

.account-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.account-field > span {
  font-size: 12px;
  font-weight: 600;
  color: var(--adm-ink-soft);
}

.account-field em {
  font-size: 11px;
  font-style: normal;
  line-height: 1.6;
  color: var(--adm-muted-2);
}

.account-field .req {
  margin-left: 2px;
  font-style: normal;
  color: var(--adm-accent);
}

.account-field :deep(.el-input__wrapper) {
  box-shadow: 0 0 0 1px var(--adm-line) inset;
}

.account-field :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--adm-accent) inset;
}

.account-note {
  padding: 10px 12px;
  border-radius: 9px;
  background: var(--adm-surface);
  font-size: 11.5px;
  line-height: 1.65;
  color: var(--adm-muted);
}

.account-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 2px;
}

@media (max-width: 960px) {
  .account-grid {
    grid-template-columns: 1fr;
  }
}
</style>

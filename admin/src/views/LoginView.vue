<template>
  <div class="login-page">
    <div class="login-shell">
      <aside class="login-brand">
        <div class="brand-mark-row">
          <img
            v-if="loginSettings.logo"
            :src="loginSettings.logo"
            class="login-logo"
            alt="Logo"
          />
          <BrandMark v-else :size="42" inverse />
          <div class="brand-mark-text">
            <strong>{{ loginSettings.brandName || '小鹊云' }}</strong>
            <span>{{ loginSettings.brandSubName || 'SITE STUDIO' }}</span>
          </div>
        </div>

        <div class="brand-copy">
          <p class="mono">{{ loginSettings.brandKicker || 'MAGPIE CLOUD' }}</p>
          <h1 class="brand-title">{{ loginSettings.brandTitle || '把官网\n编排成作品。' }}</h1>
          <p>{{ loginSettings.brandSubtitle }}</p>
        </div>

        <div class="brand-foot">
          <span>LOCAL DEMO</span>
          <span>ADMIN v1</span>
        </div>
      </aside>

      <main class="login-panel">
        <header class="login-form-head">
          <p class="mono">ADMIN SIGN IN</p>
          <h2>登录管理后台</h2>
          <p>请输入用户名和密码。</p>
        </header>

        <form class="login-form" @submit.prevent="submit">
          <label class="field">
            <span>用户名</span>
            <el-input
              v-model="username"
              size="large"
              placeholder="请输入用户名"
              autocomplete="username"
              :disabled="loading"
            />
          </label>

          <label class="field">
            <span>密码</span>
            <el-input
              v-model="password"
              type="password"
              size="large"
              placeholder="请输入密码"
              autocomplete="current-password"
              show-password
              :disabled="loading"
              @keyup.enter="submit"
            />
          </label>

          <label v-if="loginSettings.captchaEnabled" class="field">
            <span>验证码</span>
            <div class="captcha-row">
              <el-input
                v-model="captchaInput"
                size="large"
                placeholder="请输入验证码"
                maxlength="4"
                autocomplete="off"
                :disabled="loading"
                @keyup.enter="submit"
              />
              <button
                type="button"
                class="captcha-image"
                title="点击换一张验证码"
                aria-label="点击换一张验证码"
                @click="drawCaptcha"
              >
                <canvas
                  ref="captchaCanvas"
                  width="120"
                  height="42"
                  role="img"
                  aria-label="图形验证码"
                />
              </button>
            </div>
          </label>

          <p v-if="error" class="login-error">{{ error }}</p>

          <el-button
            type="primary"
            size="large"
            native-type="submit"
            :loading="loading"
            class="login-submit"
          >
            登录
          </el-button>
        </form>

        <div class="login-hint">
          <span class="mono">ADMIN ACCOUNT</span>
          <p v-if="isDev">
            用户名 <code>admin</code>
            <span class="hint-dot">·</span>
            初始密码 <code>123456</code>
          </p>
          <p v-else>使用后台管理员账号登录</p>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BrandMark from '~/components/BrandMark.vue'
import { login } from '@/services/auth'
import { usePagesStore } from '@/stores/pages'

const router = useRouter()
const route = useRoute()
const pagesStore = usePagesStore()
const loginSettings = computed(() => pagesStore.loginSettings)
const isDev = import.meta.env.DEV

const username = ref('')
const password = ref('')
const captchaInput = ref('')
const captchaCode = ref('')
const captchaCanvas = ref<HTMLCanvasElement | null>(null)
const error = ref('')
const loading = ref(false)

const CAPTCHA_CHARS = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'

function randomInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min
}

/** 演示版图形验证码：前端生成并校验，接后端后改为服务端验证码。 */
function drawCaptcha(): void {
  const code = Array.from(
    { length: 4 },
    () => CAPTCHA_CHARS[randomInt(0, CAPTCHA_CHARS.length - 1)]
  ).join('')
  captchaCode.value = code

  const canvas = captchaCanvas.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const width = canvas.width
  const height = canvas.height
  ctx.clearRect(0, 0, width, height)

  const gradient = ctx.createLinearGradient(0, 0, width, height)
  gradient.addColorStop(0, '#fbf9f5')
  gradient.addColorStop(1, '#f1ede5')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, width, height)

  for (let i = 0; i < 26; i++) {
    ctx.fillStyle = `rgba(20, 21, 24, ${0.06 + Math.random() * 0.12})`
    ctx.beginPath()
    ctx.arc(randomInt(2, width - 2), randomInt(2, height - 2), randomInt(1, 2), 0, Math.PI * 2)
    ctx.fill()
  }

  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  code.split('').forEach((char, index) => {
    ctx.save()
    ctx.translate(18 + index * 27, height / 2 + randomInt(-2, 2))
    ctx.rotate((randomInt(-16, 16) * Math.PI) / 180)
    ctx.font = `700 ${randomInt(20, 24)}px "Manrope", "PingFang SC", sans-serif`
    ctx.fillStyle = index % 2 === 0 ? '#141518' : '#c93b26'
    ctx.fillText(char, 0, 0)
    ctx.restore()
  })

  for (let i = 0; i < 3; i++) {
    ctx.strokeStyle = `rgba(201, 59, 38, ${0.14 + Math.random() * 0.12})`
    ctx.lineWidth = 1
    ctx.beginPath()
    ctx.moveTo(randomInt(0, width), randomInt(0, height))
    ctx.lineTo(randomInt(0, width), randomInt(0, height))
    ctx.stroke()
  }
}

function resetCaptcha(): void {
  captchaInput.value = ''
  if (loginSettings.value.captchaEnabled) drawCaptcha()
}

function redirectTarget(): string {
  const raw = typeof route.query.redirect === 'string' ? route.query.redirect : ''
  return raw.startsWith('/') && !raw.startsWith('//') ? raw : '/pages'
}

async function submit(): Promise<void> {
  if (loading.value) return

  const account = username.value.trim()
  if (!account || !password.value) {
    error.value = '请输入用户名和密码'
    return
  }

  if (loginSettings.value.captchaEnabled) {
    const input = captchaInput.value.trim().toUpperCase()
    if (!input) {
      error.value = '请输入验证码'
      return
    }
    if (input !== captchaCode.value) {
      error.value = '验证码错误'
      resetCaptcha()
      return
    }
  }

  loading.value = true
  error.value = ''

  // 演示版模拟请求耗时；接后端后替换为真实登录接口。
  await new Promise((resolve) => setTimeout(resolve, 260))

  try {
    await login(account, password.value)
  } catch (loginError) {
    loading.value = false
    password.value = ''
    error.value =
      loginError instanceof Error ? loginError.message : '登录失败，请稍后重试'
    resetCaptcha()
    return
  }

  loading.value = false
  await router.replace(redirectTarget())
}

onMounted(() => {
  if (loginSettings.value.captchaEnabled) drawCaptcha()
})

watch(
  () => loginSettings.value.captchaEnabled,
  (enabled) => {
    if (enabled) {
      nextTick(drawCaptcha)
    } else {
      captchaInput.value = ''
    }
  }
)
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 32px 20px;
  background:
    radial-gradient(circle at 12% 10%, rgba(232, 80, 58, 0.1), transparent 34%),
    radial-gradient(circle at 88% 88%, rgba(20, 21, 24, 0.06), transparent 36%),
    var(--adm-bg);
}

.login-shell {
  width: min(100%, 960px);
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  background: var(--adm-card);
  border: 1px solid var(--adm-line);
  border-radius: 20px;
  box-shadow: 0 1px 2px rgba(20, 21, 24, 0.03), 0 30px 80px rgba(20, 21, 24, 0.12);
  overflow: hidden;
}

.login-brand {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 42px;
  padding: 42px;
  background: #141518;
  color: #f6f2ea;
  overflow: hidden;
}

.login-brand::before {
  content: '';
  position: absolute;
  right: -110px;
  bottom: -130px;
  width: 320px;
  height: 320px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(232, 80, 58, 0.32), transparent 68%);
}

.login-brand::after {
  content: '';
  position: absolute;
  left: -80px;
  top: -100px;
  width: 240px;
  height: 240px;
  border-radius: 50%;
  border: 1px solid rgba(246, 242, 234, 0.08);
}

.login-brand > * {
  position: relative;
  z-index: 1;
}

.brand-mark-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.login-logo {
  width: 42px;
  height: 42px;
  flex: none;
  object-fit: contain;
  padding: 5px;
  border-radius: 10px;
  background: #f6f2ea;
}

.brand-mark-text {
  display: flex;
  flex-direction: column;
  gap: 3px;
  line-height: 1;
}

.brand-mark-text strong {
  font-size: 14px;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.brand-mark-text span {
  font-family: var(--adm-mono);
  font-size: 7px;
  letter-spacing: 0.16em;
  color: rgba(246, 242, 234, 0.34);
}

.brand-copy .mono {
  font-family: var(--adm-mono);
  font-size: 8.5px;
  letter-spacing: 0.16em;
  color: rgba(246, 242, 234, 0.38);
}

.brand-copy h1 {
  margin-top: 14px;
  font-size: clamp(34px, 4vw, 54px);
  font-weight: 600;
  letter-spacing: -0.055em;
  line-height: 1.08;
  white-space: pre-line;
}

.brand-copy p:last-child {
  margin-top: 18px;
  max-width: 30ch;
  font-size: 12.5px;
  line-height: 1.75;
  color: rgba(246, 242, 234, 0.5);
}

.brand-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 18px;
  border-top: 1px solid rgba(246, 242, 234, 0.1);
  font-family: var(--adm-mono);
  font-size: 8px;
  letter-spacing: 0.14em;
  color: rgba(246, 242, 234, 0.34);
}

.login-panel {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 46px 44px;
}

.login-form-head .mono {
  font-family: var(--adm-mono);
  font-size: 8.5px;
  letter-spacing: 0.16em;
  color: var(--adm-muted-2);
}

.login-form-head h2 {
  margin-top: 9px;
  font-size: 26px;
  font-weight: 600;
  letter-spacing: -0.04em;
}

.login-form-head p:last-child {
  margin-top: 7px;
  font-size: 12.5px;
  color: var(--adm-muted);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 17px;
  margin-top: 28px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.field > span {
  font-size: 11.5px;
  font-weight: 600;
  color: var(--adm-ink-soft);
}

.login-form :deep(.el-input__wrapper) {
  padding: 2px 12px;
  border-radius: 9px;
  background: #faf8f4;
  box-shadow: 0 0 0 1px var(--adm-line) inset;
  transition: box-shadow 0.18s ease, background 0.18s ease;
}

.login-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--adm-line-strong, #cfc8bb) inset;
}

.login-form :deep(.el-input__wrapper.is-focus) {
  background: #fff;
  box-shadow: 0 0 0 1px var(--adm-accent) inset;
}

.login-form :deep(.el-input__inner) {
  height: 42px;
  font-size: 13px;
}

.captcha-row {
  display: flex;
  align-items: stretch;
  gap: 10px;
}

.captcha-row :deep(.el-input) {
  flex: 1;
  min-width: 0;
}

.captcha-image {
  width: 124px;
  height: 48px;
  flex: none;
  padding: 0;
  border: 1px solid var(--adm-line);
  border-radius: 9px;
  background: #faf8f4;
  overflow: hidden;
  cursor: pointer;
  transition: border-color 0.18s ease, box-shadow 0.18s ease;
}

.captcha-image:hover {
  border-color: var(--adm-accent);
  box-shadow: 0 0 0 2px var(--adm-accent-soft);
}

.captcha-image canvas {
  display: block;
  width: 100%;
  height: 100%;
}

.login-submit {
  width: 100%;
  height: 44px;
  margin-top: 5px;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.login-error {
  margin: -4px 0 0;
  font-size: 12px;
  color: var(--el-color-danger);
}

.login-hint {
  margin-top: 28px;
  padding-top: 18px;
  border-top: 1px solid var(--adm-line-soft);
}

.login-hint .mono {
  font-family: var(--adm-mono);
  font-size: 8px;
  letter-spacing: 0.14em;
  color: var(--adm-muted-2);
}

.login-hint p {
  margin-top: 7px;
  font-size: 11.5px;
  color: var(--adm-muted);
}

.login-hint code {
  padding: 2px 6px;
  border-radius: 5px;
  background: #f1ede5;
  font-family: var(--adm-mono);
  font-size: 10.5px;
  color: var(--adm-ink-soft);
}

.hint-dot {
  margin: 0 5px;
  color: var(--adm-muted-2);
}

@media (max-width: 760px) {
  .login-page {
    padding: 18px;
  }

  .login-shell {
    grid-template-columns: 1fr;
  }

  .login-brand {
    gap: 22px;
    padding: 26px;
  }

  .brand-copy h1 {
    font-size: 30px;
  }

  .brand-copy p:last-child,
  .brand-foot {
    display: none;
  }

  .login-panel {
    padding: 30px 26px 34px;
  }
}
</style>

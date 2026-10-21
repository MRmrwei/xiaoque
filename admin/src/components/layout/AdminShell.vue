<template>
  <div class="admin-shell">
    <aside class="admin-side">
      <router-link to="/pages" class="side-brand">
        <BrandMark :size="34" inverse />
        <span class="side-brand-text">
          <strong>小鹊云</strong>
          <span>SITE STUDIO</span>
        </span>
      </router-link>

      <p class="side-label">工作台</p>
      <nav class="side-nav">
        <router-link to="/pages" class="side-link" :class="{ active: isPagesActive }">
          <el-icon><Document /></el-icon>
          <span>站点编排</span>
          <span class="side-count">{{ pagesStore.pageCount }}</span>
        </router-link>
      </nav>
      <p class="side-note">页面内容 · 顶部导航 · 底部导航 · 站点配置</p>

      <p class="side-label">线索</p>
      <nav class="side-nav">
        <router-link
          to="/submissions"
          class="side-link"
          :class="{ active: isSubmissionsActive }"
        >
          <el-icon><Message /></el-icon>
          <span>表单留言</span>
          <span class="side-count">{{ pagesStore.newSubmissionCount }}</span>
        </router-link>
      </nav>
      <p class="side-note">联系表单提交 · 跟进状态</p>

      <p class="side-label">设置</p>
      <nav class="side-nav">
        <router-link
          to="/settings"
          class="side-link"
          :class="{ active: isSettingsActive }"
        >
          <el-icon><Setting /></el-icon>
          <span>基础设置</span>
        </router-link>
      </nav>
      <p class="side-note">后台登录 · 展示文案</p>

      <div class="side-spacer" />

      <div class="demo-note">
        <strong>本地演示</strong>
        <p>数据保存在当前浏览器，不会更新线上官网，也没有接后端。</p>
        <button class="reset-link" @click="resetDemo">重置演示数据</button>
      </div>

      <div class="side-user">
        <span class="avatar side-avatar">
          <img v-if="avatarUrl" :src="avatarUrl" alt="" />
          <template v-else>{{ avatarText }}</template>
        </span>
        <div class="side-user-text">
          <strong>{{ displayName }}</strong>
          <span>管理员</span>
        </div>
        <router-link to="/settings/account" class="side-account" title="账号与安全">
          <el-icon><User /></el-icon>
        </router-link>
        <button class="side-logout" title="退出登录" @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
        </button>
      </div>
    </aside>

    <div class="admin-main">
      <header class="admin-topbar">
        <div class="crumb">
          <span>工作台</span>
          <span>/</span>
          <strong>{{ route.meta.title || '站点编排' }}</strong>
          <el-tag v-if="pagesStore.storageError" type="danger" size="small" effect="plain">
            存储异常
          </el-tag>
        </div>
        <div class="top-actions">
          <span class="demo-pill"><i /> LOCAL DEMO</span>
          <a :href="websiteUrl" target="_blank" rel="noopener" class="top-link">查看官网 ↗</a>
          <span class="avatar side-avatar">
            <img v-if="avatarUrl" :src="avatarUrl" alt="" />
            <template v-else>{{ avatarText }}</template>
          </span>
        </div>
      </header>

      <main class="admin-content">
        <router-view />
      </main>
    </div>

    <el-alert
      v-if="pagesStore.storageError"
      class="storage-alert"
      :title="pagesStore.storageError"
      type="warning"
      :closable="false"
      show-icon
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Message, Setting, SwitchButton, User } from '@element-plus/icons-vue'
import BrandMark from '~/components/BrandMark.vue'
import { usePagesStore } from '@/stores/pages'
import { logout, profileRef, sessionRef } from '@/services/auth'

const route = useRoute()
const router = useRouter()
const pagesStore = usePagesStore()
const displayName = computed(
  () => profileRef.value.name.trim() || sessionRef.value?.name?.trim() || '管理员'
)
const avatarText = computed(() => displayName.value.charAt(0) || '管')
const avatarUrl = computed(() => profileRef.value.avatar)
const websiteUrl = import.meta.env.DEV ? 'http://localhost:3000' : '/'

const isPagesActive = computed(
  () => route.path === '/pages' || route.path.startsWith('/pages/')
)
const isSubmissionsActive = computed(
  () => route.path === '/submissions' || route.path.startsWith('/submissions/')
)
const isSettingsActive = computed(
  () => route.path === '/settings' || route.path.startsWith('/settings/')
)

async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定退出当前账号吗？', '退出登录', {
      type: 'warning',
      confirmButtonText: '退出登录',
      cancelButtonText: '取消'
    })
  } catch {
    return
  }

  await logout()
  ElMessage.success('已退出登录')
  await router.replace({ name: 'login' })
}

async function resetDemo() {
  try {
    await ElMessageBox.confirm(
      '这会清除你在演示中创建/修改的页面、版本、顶部导航、底部导航、站点配置和表单留言状态，恢复到初始种子数据。确定继续吗？',
      '重置演示数据',
      { type: 'warning', confirmButtonText: '确认重置', cancelButtonText: '取消' }
    )
    pagesStore.resetDemo()
    ElMessage.success('演示数据已重置')
  } catch {
    // cancelled
  }
}
</script>

<style scoped>
.storage-alert {
  position: fixed;
  right: 20px;
  bottom: 20px;
  width: 360px;
  z-index: 100;
}

.side-user-text {
  flex: 1;
  min-width: 0;
}

.side-avatar {
  overflow: hidden;
}

.side-avatar img {
  display: block;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.side-logout,
.side-account {
  width: 26px;
  height: 26px;
  flex: none;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 7px;
  background: transparent;
  color: rgba(246, 242, 234, 0.4);
  cursor: pointer;
  transition: background 0.18s ease, color 0.18s ease;
}

.side-logout:hover,
.side-account:hover {
  background: rgba(246, 242, 234, 0.08);
  color: #f6f2ea;
}
</style>

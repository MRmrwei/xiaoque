import { createRouter, createWebHashHistory } from 'vue-router'
import AdminShell from '@/components/layout/AdminShell.vue'
import PagesView from '@/views/PagesView.vue'
import PageEditorView from '@/views/PageEditorView.vue'
import SubmissionsView from '@/views/SubmissionsView.vue'
import BasicSettingsView from '@/views/BasicSettingsView.vue'
import LoginSettingsView from '@/views/LoginSettingsView.vue'
import AccountSettingsView from '@/views/AccountSettingsView.vue'
import LoginView from '@/views/LoginView.vue'
import NotFoundView from '@/views/NotFoundView.vue'
import { getSession, isAuthenticated, refreshSession } from '@/services/auth'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { title: '登录' }
    },
    {
      path: '/',
      component: AdminShell,
      children: [
        { path: '', redirect: '/pages' },
        { path: 'pages', name: 'pages', component: PagesView, meta: { title: '站点编排' } },
        {
          path: 'pages/:id',
          name: 'page-editor',
          component: PageEditorView,
          meta: { title: '站点编排' }
        },
        {
          path: 'submissions',
          name: 'submissions',
          component: SubmissionsView,
          meta: { title: '表单留言' }
        },
        {
          path: 'settings',
          name: 'basic-settings',
          component: BasicSettingsView,
          meta: { title: '基础设置' }
        },
        {
          path: 'settings/login',
          name: 'login-settings',
          component: LoginSettingsView,
          meta: { title: '后台登录' }
        },
        {
          path: 'settings/account',
          name: 'account-settings',
          component: AccountSettingsView,
          meta: { title: '管理员账号' }
        },
        // 旧入口：导航编排已并入站点编排，保留链接兼容。
        {
          path: 'navigation',
          redirect: { name: 'pages', query: { view: 'header' } }
        },
        { path: ':pathMatch(.*)*', name: 'not-found', component: NotFoundView }
      ]
    }
  ]
})

/**
 * 判断当前是否已登录。
 * Access Token 过期时尝试用 HttpOnly Cookie 静默刷新一次；
 * 没有本地会话时不发刷新请求，避免未登录用户打开登录页也要等待接口。
 */
async function ensureAuthenticated(): Promise<boolean> {
  if (isAuthenticated()) return true
  if (!getSession()) return false
  return (await refreshSession()) !== null
}

router.beforeEach(async (to) => {
  const authenticated = await ensureAuthenticated()

  if (to.name === 'login') {
    return authenticated ? { name: 'pages' } : true
  }

  if (!authenticated) {
    return {
      name: 'login',
      query: to.fullPath && to.fullPath !== '/' ? { redirect: to.fullPath } : undefined
    }
  }

  return true
})

export default router

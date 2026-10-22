<template>
  <div class="adm-page">
    <header class="adm-page-head">
      <div>
        <p class="adm-kicker">LEADS INBOX</p>
        <h1 class="adm-title">表单留言</h1>
        <p class="adm-sub">
          查看联系表单提交的线索，跟进后标记处理状态。
        </p>
      </div>
      <div class="adm-actions">
        <el-button :loading="exporting" @click="exportCsv">导出 CSV</el-button>
      </div>
    </header>

    <el-alert
      v-if="pagesStore.submissionsError"
      class="sub-error"
      :title="pagesStore.submissionsError"
      type="error"
      :closable="false"
      show-icon
    />

    <div class="sub-stats">
      <div v-for="item in stats" :key="item.label" class="sub-stat">
        <span class="sub-stat-label">{{ item.label }}</span>
        <strong class="sub-stat-value">{{ item.value }}</strong>
        <span class="sub-stat-note">{{ item.note }}</span>
      </div>
    </div>

    <section class="adm-card" v-loading="pagesStore.submissionsLoading">
      <div class="adm-toolbar">
        <el-input
          v-model="keyword"
          class="adm-search"
          placeholder="搜索姓名、邮箱、公司或留言…"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <el-select v-model="statusFilter" style="width: 132px">
          <el-option label="全部状态" value="all" />
          <el-option label="待处理" value="new" />
          <el-option label="处理中" value="processing" />
          <el-option label="已完成" value="done" />
          <el-option label="垃圾" value="spam" />
        </el-select>

        <el-select v-model="localeFilter" style="width: 120px">
          <el-option label="全部语言" value="all" />
          <el-option label="中文" value="zh" />
          <el-option label="English" value="en" />
        </el-select>

        <span class="toolbar-spacer" />
        <span class="mono list-count">
          {{ submissions.length }} / {{ pagesStore.submissionTotal }} MESSAGES
        </span>
      </div>

      <el-table
        v-if="submissions.length"
        :data="submissions"
        row-key="id"
        class="submissions-table"
        @row-click="openDetail"
      >
        <el-table-column label="联系人" min-width="170">
          <template #default="{ row }">
            <div class="contact-cell">
              <span class="sub-avatar">{{ initial(row.name) }}</span>
              <span class="contact-text">
                <strong>{{ row.name }}</strong>
                <span>{{ row.company || '未填写公司' }}</span>
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="联系方式" min-width="190">
          <template #default="{ row }">
            <div class="contact-lines">
              <span>{{ row.email }}</span>
              <span class="muted-2">{{ row.phone || '—' }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="留言" min-width="240">
          <template #default="{ row }">
            <p class="msg-preview">{{ row.message }}</p>
          </template>
        </el-table-column>

        <el-table-column label="来源" width="130">
          <template #default="{ row }">
            <div class="source-cell">
              <span class="lang-badge" :class="{ 'is-en': row.locale === 'en' }">
                {{ row.locale === 'en' ? 'EN' : '中' }}
              </span>
              <span class="source-page">{{ row.locale === 'en' ? 'English' : '中文' }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <span class="sub-status" :class="`sub-status--${statusMeta(row.status).tone}`">
              <i />
              {{ statusMeta(row.status).label }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="提交时间" width="120">
          <template #default="{ row }">
            <span class="muted-2">{{ formatWhen(row.createdAt) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="76" align="right">
          <template #default="{ row }">
            <el-button link type="primary" @click.stop="openDetail(row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty
        v-else
        :description="
          hasFilters ? '没有符合条件的留言，换个筛选条件试试' : '还没有收到留言'
        "
      />

      <div v-if="pagesStore.submissionTotal > pageSize" class="sub-pagination">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="pagesStore.submissionTotal"
          layout="prev, pager, next"
          background
          @current-change="changePage"
        />
      </div>
    </section>

    <el-drawer
      v-model="drawerOpen"
      class="submission-drawer"
      size="min(580px, 94vw)"
      title="留言详情"
      destroy-on-close
    >
      <div v-if="active" class="sub-detail">
        <div class="detail-hero">
          <span class="sub-avatar sub-avatar--lg">{{ initial(active.name) }}</span>
          <div class="detail-hero-text">
            <strong>{{ active.name }}</strong>
            <span>{{ active.company || '未填写公司' }}</span>
            <span v-if="active.email" class="detail-email">{{ active.email }}</span>
          </div>
          <span class="sub-status" :class="`sub-status--${statusMeta(active.status).tone}`">
            <i />
            {{ statusMeta(active.status).label }}
          </span>
        </div>

        <div class="detail-section">
          <p class="detail-label">留言内容</p>
          <p class="detail-message">{{ active.message }}</p>
        </div>

        <div class="detail-section">
          <p class="detail-label">联系方式</p>
          <dl class="detail-grid">
            <div>
              <dt>邮箱</dt>
              <dd>
                <a v-if="active.email" :href="`mailto:${active.email}`">{{ active.email }}</a>
                <span v-else>—</span>
              </dd>
            </div>
            <div>
              <dt>电话</dt>
              <dd>{{ active.phone || '—' }}</dd>
            </div>
            <div>
              <dt>公司</dt>
              <dd>{{ active.company || '—' }}</dd>
            </div>
          </dl>
        </div>

        <div class="detail-section">
          <p class="detail-label">来源</p>
          <dl class="detail-grid">
            <div>
              <dt>语言</dt>
              <dd>{{ active.locale === 'en' ? 'English' : '中文' }}</dd>
            </div>
            <div>
              <dt>提交时间</dt>
              <dd>{{ formatFull(active.createdAt) }}</dd>
            </div>
            <div>
              <dt>最近处理</dt>
              <dd>{{ formatFull(active.updatedAt) }}</dd>
            </div>
          </dl>
        </div>

        <div class="detail-actions">
          <p class="detail-label">标记状态</p>
          <div class="detail-action-row">
            <el-button
              :type="active.status === 'processing' ? 'primary' : 'default'"
              @click="setStatus('processing')"
            >
              处理中
            </el-button>
            <el-button
              :type="active.status === 'done' ? 'success' : 'default'"
              @click="setStatus('done')"
            >
              已完成
            </el-button>
            <el-button
              :type="active.status === 'spam' ? 'danger' : 'default'"
              :plain="active.status === 'spam'"
              @click="setStatus('spam')"
            >
              垃圾
            </el-button>
            <el-button v-if="active.status !== 'new'" @click="setStatus('new')">
              恢复待处理
            </el-button>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import type { ContactSubmissionRecord, SubmissionStatus } from '@/types/admin'
import { usePagesStore } from '@/stores/pages'
import { apiDownload, type QueryParams } from '@/services/api'

const pagesStore = usePagesStore()
const { submissions } = storeToRefs(pagesStore)

const keyword = ref('')
const statusFilter = ref<'all' | SubmissionStatus>('all')
const localeFilter = ref<'all' | 'zh' | 'en'>('all')
const drawerOpen = ref(false)
const active = ref<ContactSubmissionRecord | null>(null)
const page = ref(1)
const pageSize = ref(20)
const exporting = ref(false)
let keywordTimer: ReturnType<typeof setTimeout> | null = null

const statusMap: Record<SubmissionStatus, { label: string; tone: string }> = {
  new: { label: '待处理', tone: 'accent' },
  processing: { label: '处理中', tone: 'warning' },
  done: { label: '已完成', tone: 'success' },
  spam: { label: '垃圾', tone: 'neutral' }
}

const hasFilters = computed(
  () =>
    statusFilter.value !== 'all' ||
    localeFilter.value !== 'all' ||
    Boolean(keyword.value.trim())
)

const stats = computed(() => {
  const current = pagesStore.submissionStats
  return [
    {
      label: '全部留言',
      value: current.total,
      note: '累计提交'
    },
    {
      label: '待处理',
      value: current.new,
      note: '需要跟进'
    },
    {
      label: '处理中',
      value: current.processing,
      note: '正在联系'
    },
    {
      label: '已完成',
      value: current.done,
      note: '已关闭'
    }
  ]
})

function reload(): void {
  void pagesStore.loadSubmissions({
    status: statusFilter.value,
    locale: localeFilter.value,
    keyword: keyword.value,
    page: page.value,
    pageSize: pageSize.value
  })
}

function changePage(value: number): void {
  page.value = value
  reload()
}

watch([statusFilter, localeFilter], () => {
  page.value = 1
  reload()
})

watch(keyword, () => {
  if (keywordTimer) clearTimeout(keywordTimer)
  keywordTimer = setTimeout(() => {
    page.value = 1
    reload()
  }, 300)
})

onMounted(reload)

onBeforeUnmount(() => {
  if (keywordTimer) clearTimeout(keywordTimer)
})

function statusMeta(status: SubmissionStatus) {
  return statusMap[status]
}

function initial(name: string): string {
  return name.trim().charAt(0).toUpperCase() || '—'
}

function formatWhen(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

function formatFull(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

function openDetail(row: unknown): void {
  active.value = row as ContactSubmissionRecord
  drawerOpen.value = true
}

async function setStatus(status: SubmissionStatus): Promise<void> {
  const current = active.value
  if (!current) return
  try {
    await pagesStore.setSubmissionStatus(current.id, status)
    active.value = pagesStore.submissions.find((item) => item.id === current.id) ?? null
    if (!active.value) drawerOpen.value = false
    ElMessage.success(`已标记为「${statusMap[status].label}」`)
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '更新留言状态失败')
  }
}

async function exportCsv(): Promise<void> {
  if (exporting.value) return
  exporting.value = true
  try {
    const query: QueryParams = {}
    if (statusFilter.value !== 'all') query.status = statusFilter.value
    if (localeFilter.value !== 'all') {
      query.locale = localeFilter.value === 'en' ? 'en' : 'zh-CN'
    }
    if (keyword.value.trim()) query.keyword = keyword.value.trim()

    const { blob, filename } = await apiDownload('/admin/submissions/export', { query })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    // 等浏览器开始下载后再释放 Blob URL，避免部分浏览器取消下载。
    setTimeout(() => URL.revokeObjectURL(url), 1000)
    ElMessage.success('导出已开始')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '导出失败，请稍后重试')
  } finally {
    exporting.value = false
  }
}
</script>

<style scoped>
.sub-error {
  margin-bottom: 16px;
}

.sub-pagination {
  display: flex;
  justify-content: flex-end;
  padding: 16px 4px 4px;
}

.sub-stats {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 18px;
}

.sub-stat {
  background: var(--adm-card);
  border: 1px solid var(--adm-line);
  border-radius: var(--adm-radius);
  box-shadow: var(--adm-shadow);
  padding: 15px 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.sub-stat-label {
  font-size: 12.5px;
  font-weight: 600;
  letter-spacing: 0.01em;
  color: var(--adm-muted);
}

.sub-stat-value {
  font-size: 25px;
  font-weight: 600;
  letter-spacing: -0.03em;
  line-height: 1.1;
}

.sub-stat-note {
  font-size: 11.5px;
  color: var(--adm-muted);
}

.submissions-table :deep(.el-table__row) {
  cursor: pointer;
}

.contact-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.sub-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  flex: none;
  background: var(--adm-accent-soft);
  color: var(--adm-accent-deep);
  font-size: 13px;
  font-weight: 700;
}

.sub-avatar--lg {
  width: 46px;
  height: 46px;
  font-size: 17px;
}

.contact-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.contact-text strong {
  font-size: 13px;
  font-weight: 600;
}

.contact-text span {
  font-size: 11.5px;
  color: var(--adm-muted-2);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.contact-lines {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 12px;
  min-width: 0;
}

.contact-lines span:first-child {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.msg-preview {
  margin: 0;
  font-size: 12px;
  line-height: 1.6;
  color: var(--adm-muted);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.source-cell {
  display: flex;
  align-items: center;
  gap: 7px;
}

.source-page {
  font-size: 12px;
  color: var(--adm-muted);
}

.sub-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  font-size: 11.5px;
  font-weight: 600;
  white-space: nowrap;
}

.sub-status i {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: currentColor;
}

.sub-status--accent {
  background: var(--adm-accent-soft);
  color: var(--adm-accent-deep);
}

.sub-status--warning {
  background: var(--adm-amber-soft);
  color: var(--adm-amber);
}

.sub-status--success {
  background: var(--adm-green-soft);
  color: var(--adm-green);
}

.sub-status--neutral {
  background: #f0ece4;
  color: #8a867c;
}

/* Detail drawer */
:global(.submission-drawer .el-drawer__header) {
  margin-bottom: 0;
  padding: 20px 22px 15px;
  border-bottom: 1px solid var(--adm-line-soft);
}

:global(.submission-drawer .el-drawer__title) {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--adm-ink);
}

:global(.submission-drawer .el-drawer__body) {
  padding: 20px 22px 24px;
}

.sub-detail {
  display: flex;
  flex-direction: column;
  gap: 22px;
  padding-bottom: 12px;
}

.detail-hero {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  border: 1px solid var(--adm-line-soft);
  border-radius: 12px;
  background: #faf8f4;
}

.detail-hero-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}

.detail-hero-text strong {
  font-size: 17px;
  font-weight: 600;
}

.detail-hero-text span {
  font-size: 12px;
  color: var(--adm-muted);
}

.detail-hero-text .detail-email {
  font-family: var(--adm-mono);
  font-size: 11.5px;
  color: var(--adm-muted-2);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.detail-section {
  display: flex;
  flex-direction: column;
  gap: 9px;
}

.detail-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
  font-weight: 600;
  letter-spacing: 0.01em;
  color: var(--adm-ink-soft);
}

.detail-label::before {
  content: '';
  width: 3px;
  height: 12px;
  flex: none;
  border-radius: 2px;
  background: var(--adm-accent);
}

.detail-message {
  margin: 0;
  padding: 14px 15px;
  border-radius: 10px;
  background: #faf8f4;
  border: 1px solid var(--adm-line-soft);
  font-size: 13px;
  line-height: 1.8;
  color: var(--adm-ink-soft);
  white-space: pre-wrap;
}

.detail-grid {
  margin: 0;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 18px;
}

.detail-grid > div {
  min-width: 0;
}

.detail-grid dt {
  font-size: 11.5px;
  color: var(--adm-muted-2);
  margin-bottom: 3px;
}

.detail-grid dd {
  margin: 0;
  font-size: 13px;
  color: var(--adm-ink-soft);
  line-height: 1.55;
  overflow-wrap: anywhere;
}

.detail-grid dd a {
  color: var(--adm-accent-deep);
  text-decoration: underline;
  text-underline-offset: 2px;
}

.detail-actions {
  display: flex;
  flex-direction: column;
  gap: 9px;
  padding-top: 16px;
  border-top: 1px solid var(--adm-line-soft);
}

.detail-action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.detail-action-row :deep(.el-button + .el-button) {
  margin-left: 0;
}

@media (max-width: 980px) {
  .sub-stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 560px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>

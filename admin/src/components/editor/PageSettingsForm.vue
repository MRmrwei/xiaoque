<template>
  <div class="settings-form">
    <section class="settings-section">
      <h3>基本信息</h3>
      <p class="settings-hint">页面名称用于后台识别，路径决定访客看到的中文与英文 URL。</p>
      <div class="prop-field">
        <span class="prop-label">页面名称</span>
        <el-input
          :model-value="page.name"
          @update:model-value="emit('update-page', 'name', $event)"
        />
      </div>
      <div class="prop-field">
        <span class="prop-label">访问路径</span>
        <div class="path-input">
          <span class="mono">{{ locale === 'en' ? '/en' : '' }}</span>
          <el-input
            :model-value="page.path"
            @update:model-value="emit('update-page', 'path', $event)"
          />
        </div>
        <span class="prop-help">修改路径后，旧链接需要在站点设置中配置 301。</span>
      </div>
      <div class="meta-row">
        <span class="lang-badge" :class="{ 'is-en': locale === 'en' }">
          {{ locale === 'en' ? 'EN' : '中' }}
        </span>
        <span class="meta-text">{{ page.published ? `线上 v${page.published.version}` : '未发布' }}</span>
        <StatusTag :status="status" />
      </div>
    </section>

    <section class="settings-section">
      <h3>SEO 与分享</h3>
      <p class="settings-hint">每个语言独立配置。SEO 长度只提示，不阻止发布。</p>
      <div class="prop-field">
        <span class="prop-label">
          页面标题
          <span class="char-count" :class="{ over: seoTitleLength > 60 }">{{ seoTitleLength }} / 60</span>
        </span>
        <el-input
          :model-value="page.seo.title"
          @update:model-value="emit('update-seo', 'title', $event)"
        />
      </div>
      <div class="prop-field">
        <span class="prop-label">
          页面描述
          <span class="char-count" :class="{ over: seoDescriptionLength > 160 }">
            {{ seoDescriptionLength }} / 160
          </span>
        </span>
        <el-input
          type="textarea"
          :rows="3"
          :model-value="page.seo.description"
          @update:model-value="emit('update-seo', 'description', $event)"
        />
      </div>
      <div class="prop-field">
        <span class="prop-label">关键词</span>
        <el-input
          :model-value="(page.seo.keywords || []).join(', ')"
          placeholder="用逗号分隔"
          @update:model-value="emit('update-keywords', $event)"
        />
      </div>

      <div class="seo-preview">
        <span class="mono">SERP PREVIEW</span>
        <p class="seo-preview-title">{{ page.seo.title || '页面标题' }}</p>
        <p class="seo-preview-url">xiaoque.com{{ displayPath }}</p>
        <p class="seo-preview-desc">{{ page.seo.description || '页面描述会显示在这里。' }}</p>
      </div>
    </section>

    <section class="settings-section">
      <h3>语言与翻译</h3>
      <p class="settings-hint">中英文共享区块结构，但草稿与发布完全独立。</p>
      <div class="translation-card">
        <div class="translation-info">
          <span class="translation-flag">{{ locale === 'en' ? 'EN' : '中' }}</span>
          <div>
            <strong>{{ locale === 'en' ? 'English' : '中文' }}</strong>
            <span>{{ statusWord(status) }}</span>
          </div>
        </div>
        <div class="translation-info">
          <span class="translation-flag muted">{{ locale === 'en' ? '中' : 'EN' }}</span>
          <div>
            <strong>{{ locale === 'en' ? '中文' : 'English' }}</strong>
            <span>{{ pair ? statusWord(pairStatus) : '尚未创建' }}</span>
          </div>
        </div>
        <el-button size="small" @click="emit('switch-locale')">
          {{ pair ? '前往编辑' : '创建翻译版本' }}
        </el-button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PageRecord, Locale } from '~/domain/types'
import type { PageStatus } from '@/types/admin'
import StatusTag from '@/components/common/StatusTag.vue'

const props = defineProps<{
  page: PageRecord
  pair?: PageRecord
  pairStatus?: PageStatus | 'missing'
  status: PageStatus
  locale: Locale
}>()

const emit = defineEmits<{
  'update-page': [key: 'name' | 'path', value: string]
  'update-seo': [key: 'title' | 'description', value: string]
  'update-keywords': [value: string]
  'switch-locale': []
}>()

const seoTitleLength = computed(() => props.page.seo.title.length)
const seoDescriptionLength = computed(() => props.page.seo.description.length)
const displayPath = computed(() => {
  const prefix = props.locale === 'en' ? '/en' : ''
  return props.page.path === '/' ? prefix || '/' : `${prefix}${props.page.path}`
})

function statusWord(status?: PageStatus | 'missing') {
  if (!status || status === 'missing') return '尚未创建'
  if (status === 'published') return '已发布'
  return '未发布'
}
</script>

<style scoped>
.path-input {
  display: flex;
  align-items: center;
  border: 1px solid var(--adm-line);
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
}

.path-input > span {
  padding-left: 10px;
  font-size: 10px;
  color: var(--adm-muted-2);
  white-space: nowrap;
}

.path-input :deep(.el-input__wrapper) {
  box-shadow: none !important;
  background: transparent;
}

.meta-row .lang-badge {
  display: inline-grid;
  place-items: center;
}
</style>

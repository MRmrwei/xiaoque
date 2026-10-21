<template>
  <div v-if="item" class="article-view">
    <header class="article-head page-canvas">
      <div class="container">
        <NuxtLink :to="localePath('/news')" class="back-link">
          <UiIcon name="arrow-left" :size="15" /> {{ t('news.back') }}
        </NuxtLink>
        <div class="article-meta">
          <span class="tag tag--accent">{{ item.category }}</span>
          <span class="meta-text">{{ formatDate(item.publishedAt, locale) }}</span>
          <span class="meta-dot" />
          <span class="meta-text">{{ item.readMinutes }} {{ t('common.minRead') }}</span>
        </div>
        <h1>{{ item.title }}</h1>
        <p class="article-summary">{{ item.summary }}</p>
        <div class="article-byline">
          <span class="byline-avatar">鹊</span>
          <div>
            <strong>{{ t('news.editorial') }}</strong>
            <span>{{ t('news.published') }} {{ formatDate(item.publishedAt, locale) }}</span>
          </div>
          <button class="share-btn" @click="share">
            <UiIcon name="link" :size="15" /> {{ t('common.share') }}
          </button>
        </div>
      </div>
    </header>

    <div class="container">
      <div class="article-cover">
        <VisualArt :variant="item.cover" :tone="item.tone" />
        <span class="cover-index mono">JOURNAL / {{ String(item.publishedAt).slice(0, 4) }}</span>
      </div>
    </div>

    <div class="container article-body-grid">
      <aside class="article-aside">
        <div class="aside-share">
          <span class="mono">SHARE</span>
          <button @click="share" :aria-label="t('common.share')"><UiIcon name="link" :size="16" /></button>
        </div>
        <div class="aside-meta">
          <span class="mono">CATEGORY</span>
          <p>{{ item.category }}</p>
        </div>
      </aside>
      <div class="article-prose prose" v-html="item.body" />
    </div>

    <div class="container article-footer">
      <div class="article-author">
        <span class="byline-avatar byline-avatar--lg">鹊</span>
        <div>
          <strong>{{ t('news.editorial') }}</strong>
          <p>{{ locale === 'en' ? 'Notes from the team building Magpie Cloud.' : '来自小鹊云团队的产品与内容笔记。' }}</p>
        </div>
        <NuxtLink :to="localePath('/news')" class="article-more">
          {{ t('news.more') }} <UiIcon name="arrow-up-right" :size="15" />
        </NuxtLink>
      </div>

      <div v-if="related.length" class="related">
        <div class="related-head">
          <span class="eyebrow">{{ t('news.related') }}</span>
        </div>
        <div class="related-grid">
          <NuxtLink
            v-for="entry in related"
            :key="entry.id"
            :to="localePath(`/news/${entry.slug}`)"
            class="related-card"
          >
            <div class="related-cover">
              <VisualArt :variant="entry.cover" :tone="entry.tone" />
            </div>
            <div class="related-text">
              <span class="related-meta">{{ entry.category }} · {{ formatDate(entry.publishedAt, locale) }}</span>
              <h3>{{ entry.title }}</h3>
            </div>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
  <NotFoundContent v-else />
</template>

<script setup lang="ts">
import NotFoundContent from '~/components/NotFoundContent.vue'
import UiIcon from '~/components/UiIcon.vue'
import VisualArt from '~/components/VisualArt.vue'
import { useNewsArticle } from '~/composables/useNewsArticle'
import { usePageSeo } from '~/composables/usePageSeo'
import { useToast } from '~/composables/useToast'
import { useUiText } from '~/composables/useUiText'
import { copyText, formatDate } from '~/utils/format'

const route = useRoute()
const localePath = useLocalePath()
const { t, locale } = useUiText()
const { notify } = useToast()

const slug = String(route.params.slug ?? '')
const { article, translation, related } = await useNewsArticle(slug)
const item = article.value

if (item) {
  const translatedSlug = translation.value?.slug
  const zhSlug = locale.value === 'zh' ? item.slug : (translatedSlug ?? item.slug)
  const enSlug = locale.value === 'en' ? item.slug : (translatedSlug ?? item.slug)

  usePageSeo({
    title: `${item.title} · ${locale.value === 'en' ? 'Magpie Cloud' : '小鹊云'}`,
    description: item.summary,
    type: 'article',
    alternates: {
      zh: `/news/${zhSlug}`,
      en: `/news/${enSlug}`
    }
  })
} else {
  if (import.meta.server) {
    const event = useRequestEvent()
    if (event) setResponseStatus(event, 404)
  }
  usePageSeo({
    title: locale.value === 'en' ? 'Story not found · Magpie Cloud' : '文章未找到 · 小鹊云',
    description:
      locale.value === 'en'
        ? 'The story you are looking for may have moved, or never existed.'
        : '你要找的文章可能已经迁移，或者从未存在过。',
    path: '/',
    noindex: true
  })
}

function share() {
  copyText(window.location.href)
  notify(t('common.copied'))
}
</script>

<style scoped>
.article-view {
  padding-bottom: clamp(60px, 8vw, 110px);
}

.article-head {
  padding-top: clamp(34px, 5vw, 62px);
  padding-bottom: clamp(32px, 4.5vw, 54px);
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 12px;
  color: var(--muted);
  transition: color 0.2s ease;
}

.back-link:hover {
  color: var(--ink);
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: clamp(28px, 4vw, 46px);
}

.meta-text {
  font-size: 12px;
  color: var(--muted);
}

.meta-dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--muted-2);
}

.article-head h1 {
  margin-top: 20px;
  font-size: clamp(36px, 5.4vw, 72px);
  line-height: 1.06;
  letter-spacing: -0.045em;
  max-width: 20ch;
}

.article-summary {
  margin-top: 22px;
  max-width: 62ch;
  font-size: clamp(16px, 1.5vw, 19px);
  line-height: 1.75;
  color: var(--muted);
}

.article-byline {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 34px;
  padding-top: 24px;
  border-top: 1px solid var(--line);
}

.byline-avatar {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--ink);
  color: var(--paper);
  font-family: var(--font-serif);
  font-size: 17px;
  flex: none;
}

.byline-avatar--lg {
  width: 48px;
  height: 48px;
  font-size: 21px;
}

.article-byline div {
  display: flex;
  flex-direction: column;
  line-height: 1.5;
}

.article-byline strong {
  font-size: 12.5px;
}

.article-byline div span {
  font-size: 11.5px;
  color: var(--muted-2);
}

.share-btn {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  height: 38px;
  padding: 0 16px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: transparent;
  font-size: 12px;
  color: var(--muted);
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}

.share-btn:hover {
  background: var(--ink);
  border-color: var(--ink);
  color: var(--paper);
}

.article-cover {
  position: relative;
  border-radius: clamp(18px, 2.4vw, 30px);
  overflow: hidden;
  aspect-ratio: 21 / 9;
  border: 1px solid var(--line-soft);
}

.cover-index {
  position: absolute;
  left: 18px;
  bottom: 16px;
  font-size: 9px;
  letter-spacing: 0.16em;
  color: rgba(246, 242, 234, 0.72);
  background: rgba(20, 21, 24, 0.35);
  backdrop-filter: blur(6px);
  padding: 6px 10px;
  border-radius: 999px;
}

.article-body-grid {
  display: grid;
  grid-template-columns: 160px minmax(0, 740px);
  justify-content: center;
  gap: clamp(24px, 5vw, 72px);
  padding-top: clamp(40px, 6vw, 70px);
}

.article-aside {
  position: sticky;
  top: 110px;
  align-self: start;
  display: flex;
  flex-direction: column;
  gap: 28px;
  padding-top: 6px;
}

.aside-share,
.aside-meta {
  display: flex;
  flex-direction: column;
  gap: 9px;
}

.article-aside .mono {
  font-size: 9px;
  letter-spacing: 0.16em;
  color: var(--muted-2);
}

.aside-share button {
  width: 40px;
  height: 40px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 50%;
  background: transparent;
  color: var(--muted);
  transition: background 0.2s ease, color 0.2s ease;
}

.aside-share button:hover {
  background: var(--accent);
  border-color: var(--accent);
  color: #fff;
}

.aside-meta p {
  font-size: 12.5px;
  color: var(--ink-soft);
}

.article-footer {
  max-width: 980px;
  padding-top: clamp(46px, 6vw, 72px);
}

.article-author {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 0;
  border-top: 1px solid var(--line);
  border-bottom: 1px solid var(--line);
}

.article-author div {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.article-author strong {
  font-size: 13px;
}

.article-author p {
  font-size: 12px;
  color: var(--muted);
}

.article-more {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--muted);
  transition: color 0.2s ease;
}

.article-more:hover {
  color: var(--accent);
}

.related {
  margin-top: clamp(40px, 5vw, 64px);
}

.related-head {
  margin-bottom: 20px;
}

.related-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.related-card {
  display: flex;
  gap: 18px;
  padding: 14px;
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  background: rgba(251, 249, 245, 0.55);
  transition: border-color 0.25s ease, transform 0.35s var(--ease), box-shadow 0.25s ease;
}

.related-card:hover {
  border-color: var(--line-strong);
  transform: translateY(-3px);
  box-shadow: var(--shadow-sm);
}

.related-cover {
  width: 118px;
  aspect-ratio: 1 / 1;
  border-radius: 12px;
  overflow: hidden;
  flex: none;
}

.related-text {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
}

.related-meta {
  font-size: 10px;
  letter-spacing: 0.05em;
  color: var(--muted-2);
}

.related-text h3 {
  font-size: 15px;
  line-height: 1.4;
  font-weight: 600;
  letter-spacing: -0.02em;
}

@media (max-width: 900px) {
  .article-body-grid {
    grid-template-columns: 1fr;
    gap: 28px;
    max-width: 740px;
  }
  .article-aside {
    position: static;
    flex-direction: row;
    gap: 30px;
    padding-top: 0;
    border-top: 1px solid var(--line);
    padding-top: 18px;
  }
  .aside-meta {
    margin-left: auto;
    text-align: right;
  }
  .aside-share {
    flex-direction: row;
    align-items: center;
  }
}

@media (max-width: 720px) {
  .article-cover {
    aspect-ratio: 4 / 3;
  }
  .article-head h1 {
    font-size: clamp(32px, 9vw, 46px);
  }
  .article-byline {
    flex-wrap: wrap;
  }
  .share-btn {
    margin-left: 0;
  }
  .related-grid {
    grid-template-columns: 1fr;
  }
  .article-author {
    flex-wrap: wrap;
  }
}
</style>

<template>
  <section
    class="hero-carousel"
    :class="{ 'hero-carousel--en': locale === 'en' }"
    role="region"
    aria-roledescription="carousel"
    :aria-label="locale === 'en' ? 'Homepage highlights' : '首页精选'"
    tabindex="0"
    @mouseenter="paused = true"
    @mouseleave="paused = false"
    @keydown.left.prevent="prev"
    @keydown.right.prevent="next"
  >
    <div class="carousel-grid-bg" aria-hidden="true" />
    <div class="carousel-glow" aria-hidden="true" />

    <div class="container carousel-inner">
      <div class="carousel-copy">
        <div class="carousel-meta">
          <span class="carousel-count mono">
            {{ pad(active + 1) }} <i>/</i> {{ pad(slides.length) }}
          </span>
          <span class="eyebrow">{{ slide.eyebrow }}</span>
        </div>

        <Transition name="carousel-copy" mode="out-in">
          <div :key="active" class="carousel-copy-inner" aria-live="polite">
            <h1 v-html="titleHtml(slide.title)" />
            <p v-if="slide.subtitle" class="carousel-sub">{{ slide.subtitle }}</p>
            <div v-if="slide.primaryLabel || slide.secondaryLabel" class="btn-row carousel-actions">
              <BlockLink
                v-if="slide.primaryLabel"
                :href="slide.primaryHref"
                :preview="preview"
                class="btn btn--primary btn--lg"
              >
                {{ slide.primaryLabel }} <span class="arr">↗</span>
              </BlockLink>
              <BlockLink
                v-if="slide.secondaryLabel"
                :href="slide.secondaryHref"
                :preview="preview"
                class="btn btn--ghost btn--lg"
              >
                {{ slide.secondaryLabel }}
              </BlockLink>
            </div>
          </div>
        </Transition>

        <div class="carousel-controls">
          <div class="carousel-arrows">
            <button
              type="button"
              :aria-label="locale === 'en' ? 'Previous slide' : '上一张'"
              @click="prev"
            >
              <UiIcon name="arrow-left" :size="17" />
            </button>
            <button
              type="button"
              :aria-label="locale === 'en' ? 'Next slide' : '下一张'"
              @click="next"
            >
              <UiIcon name="arrow" :size="17" />
            </button>
          </div>

          <div class="carousel-dots" role="tablist">
            <button
              v-for="(item, index) in slides"
              :key="item.id"
              type="button"
              role="tab"
              :aria-selected="index === active"
              :aria-label="
                locale === 'en' ? `Go to slide ${index + 1}` : `切换到第 ${index + 1} 张`
              "
              :class="{ active: index === active }"
              @click="go(index)"
            >
              <span />
            </button>
          </div>

          <span v-if="isVideoSlide" class="video-autoplay-note mono">
            {{
              videoEnded
                ? locale === 'en'
                  ? 'VIDEO ENDED · RESUMING'
                  : '视频结束 · 继续轮播'
                : locale === 'en'
                  ? 'VIDEO · AUTOPLAY PAUSED'
                  : '视频播放中 · 暂停轮播'
            }}
          </span>
          <button
            v-else
            type="button"
            class="autoplay-btn"
            :title="autoplayOn ? (locale === 'en' ? 'Pause autoplay' : '暂停自动轮播') : locale === 'en' ? 'Start autoplay' : '开始自动轮播'"
            @click="toggleAutoplay"
          >
            <UiIcon :name="autoplayOn ? 'pause' : 'play'" :size="12" />
            <span>{{ autoplayOn ? (locale === 'en' ? 'Auto' : '自动') : locale === 'en' ? 'Play' : '播放' }}</span>
          </button>
        </div>

        <p v-if="data.note" class="carousel-note">{{ data.note }}</p>
      </div>

      <div
        class="carousel-media"
        @pointerdown="onPointerDown"
        @pointerup="onPointerUp"
        @pointercancel="onPointerUp"
      >
        <div class="media-frame" :class="{ 'is-video': isVideoSlide }">
          <div class="media-track" :style="{ transform: `translateX(-${active * 100}%)` }">
            <div
              v-for="(item, index) in slides"
              :key="item.id"
              class="media-slide"
              :aria-hidden="index !== active"
            >
              <template v-if="item.mediaType === 'video' && item.videoUrl && index === active && !videoError">
                <video
                  :ref="setVideoRef"
                  :src="item.videoUrl"
                  :poster="videoPoster(item)"
                  autoplay
                  muted
                  playsinline
                  @play="onVideoPlay"
                  @pause="onVideoPause"
                  @ended="onVideoEnded"
                  @timeupdate="onVideoTimeUpdate"
                  @loadedmetadata="onVideoMetadata"
                  @error="videoError = true"
                />
                <span class="media-scrim" aria-hidden="true" />
              </template>
              <HeroBuilderMock
                v-else-if="item.mediaType === 'image' && item.media === 'builder'"
                :locale="locale"
              />
              <BlockImage
                v-else
                :src="imageSource(item)"
                fallback="orbit"
              />

              <span class="media-type mono">
                {{ item.mediaType === 'video' ? 'VIDEO' : 'IMAGE' }}
              </span>
              <span v-if="item.mediaType === 'video'" class="video-badge">
                <UiIcon name="play" :size="14" />
                {{ item.videoUrl ? (locale === 'en' ? 'Video' : '视频') : locale === 'en' ? 'Add video URL' : '待配置视频' }}
              </span>
              <span
                v-if="item.mediaType === 'video' && (!item.videoUrl || videoError)"
                class="video-cover-play"
                aria-hidden="true"
              >
                <UiIcon name="play" :size="28" />
              </span>
            </div>
          </div>

          <div v-if="isVideoSlide" class="video-controls" @pointerdown.stop @pointerup.stop @click.stop>
            <button
              type="button"
              class="vc-play"
              :aria-label="videoPlaying ? (locale === 'en' ? 'Pause video' : '暂停视频') : locale === 'en' ? 'Play video' : '播放视频'"
              @click="toggleVideoPlay"
            >
              <UiIcon :name="videoPlaying ? 'pause' : 'play'" :size="14" />
            </button>
            <span class="vc-time mono">
              {{ formatTime(currentTime) }} / {{ formatTime(duration) }}
            </span>
            <input
              class="vc-range"
              type="range"
              min="0"
              :max="duration || 0"
              step="0.1"
              :value="currentTime"
              :disabled="!duration"
              :aria-label="locale === 'en' ? 'Video position' : '视频播放进度'"
              @input="seekVideo"
            />
            <span class="vc-note mono">
              {{
                videoEnded
                  ? locale === 'en'
                    ? 'ENDED · RESUMING'
                    : '已播完 · 继续轮播'
                  : locale === 'en'
                    ? 'PLAYING · NO AUTOPLAY'
                    : '播放中 · 不轮播'
              }}
            </span>
          </div>

          <div class="media-progress" aria-hidden="true">
            <span
              :key="`progress-${active}`"
              :style="{
                animationDuration: `${safeInterval}ms`,
                animationPlayState: autoplayOn && !paused && !isVideoSlide ? 'running' : 'paused'
              }"
            />
          </div>

          <div class="media-footer">
            <span class="mono">{{ slide.mediaType === 'video' ? 'VIDEO SLIDE' : 'MEDIA SLIDE' }}</span>
            <span class="mono">{{ pad(active + 1) }} / {{ pad(slides.length) }}</span>
          </div>
        </div>
        <span class="carousel-side-note mono" aria-hidden="true">HEADER / 01</span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { ComponentPublicInstance } from 'vue'
import BlockImage from '~/components/BlockImage.vue'
import BlockLink from '~/components/BlockLink.vue'
import HeroBuilderMock from '~/components/blocks/HeroBuilderMock.vue'
import UiIcon from '~/components/UiIcon.vue'
import type { HeroCarouselProps, HeroCarouselSlide } from '~/domain/blocks'
import type { Locale } from '~/domain/types'
import { titleHtml } from '~/utils/format'
import { isCustomImage } from '@/utils/image'

const props = withDefaults(
  defineProps<{
    propsData: HeroCarouselProps
    locale: Locale
    preview?: boolean
  }>(),
  { preview: false }
)

const fallbackSlide: HeroCarouselSlide = {
  id: 'fallback',
  mediaType: 'image',
  media: 'orbit',
  videoUrl: '',
  eyebrow: '小鹊云',
  title: '让每一次表达\n都成为<em>增长</em>。',
  subtitle: '在这里添加你的第一张轮播内容。',
  primaryLabel: '预约演示',
  primaryHref: '/contact',
  secondaryLabel: '',
  secondaryHref: ''
}

const data = computed(() => props.propsData)
const slides = computed<HeroCarouselSlide[]>(() =>
  data.value.slides?.length ? data.value.slides : [fallbackSlide]
)
const active = ref(0)
const paused = ref(false)
const autoplayOn = ref(data.value.autoplay !== false)
const videoPlaying = ref(false)
const videoError = ref(false)
const reducedMotion = ref(false)
const videoEl = ref<HTMLVideoElement | null>(null)
const currentTime = ref(0)
const duration = ref(0)
const videoEnded = ref(false)

const slide = computed(() => slides.value[Math.min(active.value, slides.value.length - 1)] ?? fallbackSlide)
const safeInterval = computed(() => Math.max(3000, Number(data.value.interval) || 7000))
const isVideoSlide = computed(
  () => slide.value.mediaType === 'video' && Boolean(slide.value.videoUrl) && !videoError.value
)
const holdAutoplay = computed(() => isVideoSlide.value && !videoEnded.value)

let timer: ReturnType<typeof setInterval> | undefined
let resumeTimer: ReturnType<typeof setTimeout> | undefined

function videoPoster(item: HeroCarouselSlide): string | undefined {
  if (isCustomImage(item.videoPoster)) return item.videoPoster
  // 兼容旧数据：视频封面曾经保存在 media 字段里。
  return item.mediaType === 'video' && isCustomImage(item.media) ? item.media : undefined
}

function imageSource(item: HeroCarouselSlide): string {
  const source = item.mediaType === 'video' ? item.videoPoster || item.media : item.media
  return source === 'builder' ? '' : source
}

function pad(value: number): string {
  return String(value).padStart(2, '0')
}

function stopTimer() {
  if (timer) clearInterval(timer)
  timer = undefined
}

function clearResumeTimer() {
  if (resumeTimer) clearTimeout(resumeTimer)
  resumeTimer = undefined
}

function restartTimer() {
  stopTimer()
  if (
    typeof window === 'undefined' ||
    !autoplayOn.value ||
    paused.value ||
    reducedMotion.value ||
    holdAutoplay.value ||
    slides.value.length < 2
  ) {
    return
  }
  timer = setInterval(() => {
    go((active.value + 1) % slides.value.length, false)
  }, safeInterval.value)
}

function go(index: number, restart = true) {
  const length = slides.value.length
  if (!length) return
  active.value = ((index % length) + length) % length
  videoPlaying.value = false
  videoError.value = false
  videoEnded.value = false
  currentTime.value = 0
  duration.value = 0
  clearResumeTimer()
  if (restart) restartTimer()
}

function next() {
  go(active.value + 1)
}

function prev() {
  go(active.value - 1)
}

function toggleAutoplay() {
  autoplayOn.value = !autoplayOn.value
  restartTimer()
}

function setVideoRef(el: Element | ComponentPublicInstance | null) {
  videoEl.value = el instanceof HTMLVideoElement ? el : null
  if (!videoEl.value) {
    currentTime.value = 0
    duration.value = 0
  }
}

function onVideoPlay() {
  videoPlaying.value = true
  videoEnded.value = false
  clearResumeTimer()
  restartTimer()
}

function onVideoPause() {
  videoPlaying.value = false
  restartTimer()
}

function onVideoEnded() {
  videoPlaying.value = false
  videoEnded.value = true
  clearResumeTimer()
  // Give the viewer a beat before moving on, then resume the normal carousel.
  resumeTimer = setTimeout(() => {
    if (isVideoSlide.value) next()
  }, 1500)
}

function onVideoTimeUpdate() {
  if (videoEl.value) currentTime.value = videoEl.value.currentTime
}

function onVideoMetadata() {
  if (videoEl.value) duration.value = videoEl.value.duration || 0
}

function toggleVideoPlay() {
  const video = videoEl.value
  if (!video) return
  if (video.ended || videoEnded.value) {
    video.currentTime = 0
    videoEnded.value = false
    clearResumeTimer()
    void video.play()
    return
  }
  if (video.paused) void video.play()
  else video.pause()
}

function seekVideo(event: Event) {
  const value = Number((event.target as HTMLInputElement).value)
  if (!Number.isFinite(value)) return
  videoEnded.value = false
  clearResumeTimer()
  if (videoEl.value) videoEl.value.currentTime = value
  currentTime.value = value
  restartTimer()
}

function formatTime(value: number): string {
  if (!Number.isFinite(value) || value < 0) return '00:00'
  const total = Math.floor(value)
  const minutes = Math.floor(total / 60)
  const seconds = total % 60
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
}

let pointerStartX = 0

function onPointerDown(event: PointerEvent) {
  pointerStartX = event.clientX
}

function onPointerUp(event: PointerEvent) {
  const delta = event.clientX - pointerStartX
  if (Math.abs(delta) < 48) return
  if (delta < 0) next()
  else prev()
}

watch(
  () => [autoplayOn.value, paused.value, slides.value.length, safeInterval.value],
  () => restartTimer()
)

watch(active, () => {
  videoPlaying.value = false
})

watch(holdAutoplay, () => restartTimer())

onMounted(() => {
  reducedMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  restartTimer()
})

onBeforeUnmount(() => {
  stopTimer()
  clearResumeTimer()
})
</script>

<style scoped>
.hero-carousel {
  position: relative;
  overflow: hidden;
  padding-block: clamp(42px, 6vw, 88px) clamp(50px, 7vw, 96px);
  background: var(--paper);
  border-bottom: 1px solid var(--line-soft);
  outline: none;
}

.hero-carousel:focus-visible {
  box-shadow: inset 0 0 0 2px var(--accent);
}

.carousel-grid-bg {
  position: absolute;
  inset: 0;
  background-image: linear-gradient(to right, rgba(20, 21, 24, 0.035) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(20, 21, 24, 0.035) 1px, transparent 1px);
  background-size: 88px 88px;
  mask-image: radial-gradient(ellipse 80% 75% at 55% 40%, #000 10%, transparent 76%);
  -webkit-mask-image: radial-gradient(ellipse 80% 75% at 55% 40%, #000 10%, transparent 76%);
  pointer-events: none;
}

.carousel-glow {
  position: absolute;
  width: 700px;
  height: 700px;
  right: -160px;
  top: -260px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(232, 80, 58, 0.11), transparent 66%);
  pointer-events: none;
}

.carousel-inner {
  position: relative;
  z-index: 2;
  display: grid;
  grid-template-columns: minmax(0, 1.02fr) minmax(0, 0.98fr);
  align-items: center;
  gap: clamp(34px, 5vw, 72px);
}

.carousel-copy {
  min-width: 0;
}

.carousel-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.carousel-count {
  font-size: 9px;
  letter-spacing: 0.12em;
  color: var(--muted-2);
  display: inline-flex;
  align-items: baseline;
  gap: 3px;
}

.carousel-count i {
  font-style: normal;
  color: var(--line-strong);
}

.carousel-copy-inner {
  min-height: 330px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.carousel-copy-inner h1 {
  font-size: clamp(40px, 4.8vw, 68px);
  line-height: 1.03;
  letter-spacing: -0.05em;
  font-weight: 600;
  max-width: 12ch;
}

.hero-carousel--en .carousel-copy-inner h1 {
  font-size: clamp(34px, 4.1vw, 56px);
  letter-spacing: -0.045em;
}

.carousel-copy-inner h1 :deep(em) {
  font-family: var(--font-serif);
  font-style: normal;
  font-weight: 500;
  color: var(--accent);
  letter-spacing: -0.015em;
}

.carousel-sub {
  margin-top: 22px;
  max-width: 46ch;
  color: var(--muted);
  font-size: clamp(14.5px, 1.3vw, 17px);
  line-height: 1.75;
}

.carousel-actions {
  margin-top: 30px;
}

.carousel-controls {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 32px;
}

.carousel-arrows {
  display: inline-flex;
  gap: 6px;
}

.carousel-arrows button {
  width: 40px;
  height: 40px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line-strong);
  border-radius: 50%;
  background: transparent;
  color: var(--ink);
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease,
    transform 0.3s var(--ease);
}

.carousel-arrows button:hover {
  background: var(--ink);
  border-color: var(--ink);
  color: var(--paper);
  transform: translateY(-2px);
}

.carousel-dots {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-left: 4px;
}

.carousel-dots button {
  padding: 8px 2px;
  border: 0;
  background: transparent;
  display: grid;
  place-items: center;
}

.carousel-dots button span {
  display: block;
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: var(--line-strong);
  transition: width 0.35s var(--ease), background 0.25s ease;
}

.carousel-dots button.active span {
  width: 24px;
  background: var(--accent);
}

.carousel-dots button:hover span {
  background: var(--ink);
}

.autoplay-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  margin-left: 6px;
  border: 0;
  background: transparent;
  padding: 6px 0;
  font-family: var(--font-mono);
  font-size: 8.5px;
  letter-spacing: 0.1em;
  color: var(--muted-2);
  white-space: nowrap;
  transition: color 0.2s ease;
}

.autoplay-btn:hover {
  color: var(--ink);
}

.video-autoplay-note {
  min-width: 0;
  margin-left: 6px;
  font-size: 8px;
  letter-spacing: 0.08em;
  color: var(--muted-2);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.carousel-note {
  margin-top: 22px;
  font-size: 10.5px;
  letter-spacing: 0.01em;
  color: var(--muted-2);
}

.carousel-media {
  position: relative;
  min-width: 0;
  touch-action: pan-y;
}

.media-frame {
  position: relative;
  aspect-ratio: 1.08 / 1;
  border-radius: clamp(18px, 2.2vw, 28px);
  overflow: hidden;
  background: var(--paper-deep);
  box-shadow: var(--shadow-lg);
  border: 1px solid rgba(20, 21, 24, 0.08);
}

.media-track {
  display: flex;
  height: 100%;
  transition: transform 0.8s var(--ease);
  will-change: transform;
}

.media-slide {
  position: relative;
  flex: 0 0 100%;
  height: 100%;
  overflow: hidden;
  background: var(--paper-deep);
}

.media-slide > :deep(.visual-art) {
  width: 100%;
  height: 100%;
}

.media-slide video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.media-slide :deep(.builder-mock) {
  height: 100%;
  display: flex;
  align-items: center;
}

.media-slide :deep(.builder-mock .browser) {
  width: 100%;
}

.media-scrim {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(20, 21, 24, 0.48), transparent 55%);
  pointer-events: none;
}

.media-type {
  position: absolute;
  top: 16px;
  left: 16px;
  font-size: 8px;
  letter-spacing: 0.14em;
  color: rgba(246, 242, 234, 0.78);
  background: rgba(20, 21, 24, 0.35);
  backdrop-filter: blur(8px);
  padding: 5px 9px;
  border-radius: 999px;
}

.video-badge {
  position: absolute;
  right: 16px;
  top: 16px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 9px;
  font-weight: 600;
  color: var(--paper);
  background: var(--accent);
  padding: 5px 9px;
  border-radius: 999px;
}

.video-cover-play {
  position: absolute;
  left: 50%;
  top: 50%;
  z-index: 2;
  width: clamp(58px, 7vw, 78px);
  height: clamp(58px, 7vw, 78px);
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: rgba(246, 242, 234, 0.88);
  color: var(--ink);
  backdrop-filter: blur(10px);
  transform: translate(-50%, -50%);
  box-shadow: 0 12px 34px rgba(20, 21, 24, 0.2);
  pointer-events: none;
}

.media-frame.is-video .media-footer,
.media-frame.is-video .media-progress {
  display: none;
}

.video-controls {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 4;
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 10px 14px 12px;
  background: linear-gradient(to top, rgba(20, 21, 24, 0.82), rgba(20, 21, 24, 0.42) 60%, transparent);
  color: var(--paper);
}

.vc-play {
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(246, 242, 234, 0.35);
  border-radius: 50%;
  background: rgba(246, 242, 234, 0.12);
  color: var(--paper);
  flex: none;
  transition: background 0.2s ease, color 0.2s ease;
}

.vc-play:hover {
  background: var(--accent);
  border-color: var(--accent);
  color: #fff;
}

.vc-time {
  font-size: 8px;
  letter-spacing: 0.04em;
  color: rgba(246, 242, 234, 0.76);
  white-space: nowrap;
  flex: none;
}

.vc-range {
  flex: 1;
  min-width: 40px;
  height: 16px;
  margin: 0;
  background: transparent;
  cursor: pointer;
  -webkit-appearance: none;
  appearance: none;
}

.vc-range:disabled {
  opacity: 0.5;
  cursor: default;
}

.vc-range::-webkit-slider-runnable-track {
  height: 3px;
  border-radius: 999px;
  background: rgba(246, 242, 234, 0.28);
}

.vc-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 11px;
  height: 11px;
  margin-top: -4px;
  border-radius: 50%;
  background: var(--paper);
  border: 2px solid var(--accent);
}

.vc-range::-moz-range-track {
  height: 3px;
  border-radius: 999px;
  background: rgba(246, 242, 234, 0.28);
}

.vc-range::-moz-range-thumb {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: var(--paper);
  border: 2px solid var(--accent);
}

.vc-note {
  font-size: 7px;
  letter-spacing: 0.08em;
  color: rgba(246, 242, 234, 0.58);
  white-space: nowrap;
  flex: none;
}

.media-progress {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 2px;
  background: rgba(246, 242, 234, 0.2);
  z-index: 3;
}

.media-progress span {
  display: block;
  height: 100%;
  width: 100%;
  background: var(--accent);
  transform-origin: left center;
  animation-name: carousel-progress;
  animation-timing-function: linear;
  animation-fill-mode: forwards;
}

@keyframes carousel-progress {
  from {
    transform: scaleX(0);
  }
  to {
    transform: scaleX(1);
  }
}

.media-footer {
  position: absolute;
  left: 14px;
  right: 14px;
  bottom: 12px;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding: 5px 10px;
  border-radius: 999px;
  background: rgba(20, 21, 24, 0.34);
  backdrop-filter: blur(8px);
  color: rgba(246, 242, 234, 0.78);
  pointer-events: none;
}

.media-footer .mono {
  font-size: 7.5px;
  letter-spacing: 0.14em;
}

.carousel-side-note {
  position: absolute;
  right: 12px;
  bottom: 10px;
  writing-mode: vertical-rl;
  font-size: 7.5px;
  letter-spacing: 0.22em;
  color: var(--muted-2);
  opacity: 0.65;
}

.carousel-copy-enter-active,
.carousel-copy-leave-active {
  transition: opacity 0.28s ease, transform 0.42s var(--ease);
}

.carousel-copy-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

.carousel-copy-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

@container canvas (max-width: 980px) {
  .carousel-inner {
    grid-template-columns: 1fr;
    gap: 38px;
  }
  .carousel-copy-inner {
    min-height: 0;
  }
  .carousel-copy-inner h1 {
    max-width: 100%;
  }
  .media-frame {
    aspect-ratio: 16 / 10;
    max-width: 720px;
  }
  .carousel-side-note {
    display: none;
  }
}

@container canvas (max-width: 620px) {
  .hero-carousel {
    padding-block: 34px 52px;
  }
  .carousel-meta {
    margin-bottom: 14px;
  }
  .carousel-copy-inner h1 {
    font-size: clamp(34px, 10.5vw, 48px);
  }
  .carousel-sub {
    margin-top: 16px;
    font-size: 14.5px;
  }
  .carousel-actions {
    margin-top: 22px;
  }
  .carousel-actions .btn {
    height: 47px;
    padding: 0 19px;
    font-size: 13px;
  }
  .carousel-controls {
    margin-top: 24px;
    gap: 12px;
  }
  .autoplay-btn {
    margin-left: 0;
  }
  .media-frame {
    aspect-ratio: 4 / 3.2;
    border-radius: 18px;
  }
  .carousel-note {
    margin-top: 16px;
    max-width: 34ch;
  }
  .media-type,
  .video-badge {
    top: 10px;
  }
  .media-type {
    left: 10px;
  }
  .video-badge {
    right: 10px;
  }
  .video-controls {
    gap: 7px;
    padding: 8px 10px 10px;
  }
  .vc-time {
    font-size: 7.5px;
  }
  .vc-note {
    display: none;
  }
}
</style>

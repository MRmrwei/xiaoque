<template>
  <section class="contact-block section">
    <div class="container contact-grid">
      <div class="contact-intro">
        <p v-if="kicker" class="eyebrow">{{ kicker }}</p>
        <h2 v-html="title" />
        <p v-if="data.subtitle" class="contact-sub">{{ data.subtitle }}</p>

        <div class="contact-art">
          <BlockImage :src="data.mapStyle" tone="ink" fallback="orbit" />
          <span v-if="mapNote" class="contact-art-label mono">{{ mapNote }}</span>
        </div>

        <div class="contact-info">
          <div class="info-item">
            <span class="info-icon"><UiIcon name="mail" :size="16" /></span>
            <div>
              <span class="info-label mono">{{ t('contact.email') }}</span>
              <a :href="`mailto:${data.email}`">{{ data.email }}</a>
            </div>
          </div>
          <div class="info-item">
            <span class="info-icon"><UiIcon name="phone" :size="16" /></span>
            <div>
              <span class="info-label mono">{{ t('contact.phone') }}</span>
              <a :href="`tel:${data.phone}`">{{ data.phone }}</a>
            </div>
          </div>
          <div class="info-item">
            <span class="info-icon"><UiIcon name="pin" :size="16" /></span>
            <div>
              <span class="info-label mono">{{ t('contact.address') }}</span>
              <p>{{ data.address }}</p>
            </div>
          </div>
          <div class="info-item">
            <span class="info-icon"><UiIcon name="clock" :size="16" /></span>
            <div>
              <span class="info-label mono">{{ t('contact.hours') }}</span>
              <p>{{ data.hours }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="form-card">
        <div class="form-card-head">
          <span v-if="formKicker" class="mono">{{ formKicker }}</span>
          <span v-if="replyNote" class="form-live"><i />{{ replyNote }}</span>
        </div>

        <Transition name="form-fade" mode="out-in">
          <form v-if="!sent" key="form" class="c-form" novalidate @submit.prevent="onSubmit">
            <div class="form-row">
              <div class="c-field" :class="{ invalid: errors.name }">
                <label for="cf-name">{{ t('form.name') }} <em>*</em></label>
                <input
                  id="cf-name"
                  v-model="form.name"
                  type="text"
                  :placeholder="t('form.namePlaceholder')"
                  @input="errors.name = ''"
                />
                <span v-if="errors.name" class="error-text">{{ errors.name }}</span>
              </div>
              <div class="c-field">
                <label for="cf-company">{{ t('form.company') }}</label>
                <input
                  id="cf-company"
                  v-model="form.company"
                  type="text"
                  :placeholder="t('form.companyPlaceholder')"
                />
              </div>
            </div>

            <div class="form-row">
              <div class="c-field" :class="{ invalid: errors.email }">
                <label for="cf-email">{{ t('form.email') }} <em>*</em></label>
                <input
                  id="cf-email"
                  v-model="form.email"
                  type="email"
                  :placeholder="t('form.emailPlaceholder')"
                  @input="errors.email = ''"
                />
                <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
              </div>
              <div class="c-field">
                <label for="cf-phone">{{ t('form.phone') }}</label>
                <input
                  id="cf-phone"
                  v-model="form.phone"
                  type="text"
                  :placeholder="t('form.phonePlaceholder')"
                />
              </div>
            </div>

            <div class="c-field" :class="{ invalid: errors.message }">
              <label for="cf-message">{{ t('form.message') }} <em>*</em></label>
              <textarea
                id="cf-message"
                v-model="form.message"
                rows="5"
                :placeholder="t('form.messagePlaceholder')"
                @input="errors.message = ''"
              />
              <span v-if="errors.message" class="error-text">{{ errors.message }}</span>
            </div>

            <button type="submit" class="btn btn--primary form-submit" :disabled="sending">
              <span v-if="!sending">{{ data.submitLabel || t('form.submit') }} <span class="arr">↗</span></span>
              <span v-else>{{ t('form.submitting') }}</span>
            </button>

            <p class="form-privacy">{{ t('form.privacy') }}</p>
          </form>

          <div v-else key="success" class="form-success">
            <span class="success-icon"><UiIcon name="check" :size="26" /></span>
            <h3>{{ data.successTitle || t('form.successTitle') }}</h3>
            <p>{{ data.successText || t('form.successText') }}</p>
            <button class="btn btn--ghost btn--sm" @click="reset">{{ t('form.another') }}</button>
            <span class="success-note mono">MESSAGE · RECEIVED</span>
          </div>
        </Transition>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BlockImage from '~/components/BlockImage.vue'
import UiIcon from '~/components/UiIcon.vue'
import { useContactForm } from '~/composables/useContactForm'
import { useUiText } from '~/composables/useUiText'
import type { ContactFormProps } from '~/domain/blocks'
import type { Locale } from '~/domain/types'
import { titleHtml } from '~/utils/format'

const props = withDefaults(
  defineProps<{
    propsData: ContactFormProps
    locale: Locale
    preview?: boolean
  }>(),
  { preview: false }
)

const { t, locale } = useUiText()
const data = computed(() => props.propsData)
const title = computed(() => titleHtml(data.value.title || ''))
const kicker = computed(() =>
  data.value.kicker === undefined ? 'CONTACT / 01' : data.value.kicker
)
const mapNote = computed(() =>
  data.value.mapNote === undefined ? 'SHANGHAI · 31.2304° N, 121.4737° E' : data.value.mapNote
)
const formKicker = computed(() =>
  data.value.formKicker === undefined ? 'PROJECT INQUIRY' : data.value.formKicker
)
const replyNote = computed(() =>
  data.value.replyNote === undefined
    ? locale.value === 'en'
      ? 'Replies within 1 day'
      : '一个工作日内回复'
    : data.value.replyNote
)
const { form, errors, sending, sent, submit, reset } = useContactForm()

async function onSubmit() {
  if (props.preview) return
  await submit()
}
</script>

<style scoped>
.contact-block {
  background: var(--paper);
}

.contact-grid {
  display: grid;
  grid-template-columns: minmax(0, 0.92fr) minmax(0, 1.08fr);
  gap: clamp(36px, 6vw, 90px);
  align-items: start;
}

.contact-intro h2 {
  margin-top: 18px;
  font-size: clamp(30px, 3.9vw, 52px);
  max-width: 13ch;
}

.contact-sub {
  margin-top: 16px;
  max-width: 40ch;
  color: var(--muted);
  font-size: 14.5px;
  line-height: 1.7;
}

.contact-art {
  position: relative;
  margin-top: 30px;
  aspect-ratio: 16 / 8.6;
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid var(--line-soft);
}

.contact-art-label {
  position: absolute;
  left: 14px;
  bottom: 12px;
  font-size: 7.5px;
  letter-spacing: 0.13em;
  color: rgba(246, 242, 234, 0.6);
}

.contact-info {
  margin-top: 28px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px 24px;
  border-top: 1px solid var(--line);
  padding-top: 24px;
}

.info-item {
  display: flex;
  gap: 11px;
  align-items: flex-start;
}

.info-icon {
  width: 32px;
  height: 32px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 50%;
  color: var(--accent);
  flex: none;
}

.info-item > div {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.info-label {
  font-size: 8px;
  letter-spacing: 0.13em;
  color: var(--muted-2);
}

.info-item a,
.info-item p {
  font-size: 13px;
  color: var(--ink-soft);
  line-height: 1.45;
  word-break: break-word;
}

.info-item a:hover {
  color: var(--accent);
}

.form-card {
  background: var(--paper-light);
  border: 1px solid var(--line);
  border-radius: clamp(18px, 2.4vw, 26px);
  padding: clamp(24px, 3.4vw, 38px);
  box-shadow: var(--shadow-sm);
  position: relative;
  overflow: hidden;
}

.form-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--accent), #f1b7a8 55%, transparent);
}

.form-card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  padding-bottom: 18px;
  margin-bottom: 22px;
  border-bottom: 1px solid var(--line-soft);
}

.form-card-head .mono {
  font-size: 8.5px;
  letter-spacing: 0.16em;
  color: var(--muted-2);
}

.form-live {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  color: var(--muted);
}

.form-live i {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #6aa374;
}

.c-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.c-field {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.c-field label {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.02em;
  color: var(--muted);
}

.c-field label em {
  font-style: normal;
  color: var(--accent);
  margin-left: 2px;
}

.c-field input,
.c-field textarea {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: #fff;
  padding: 11px 13px;
  font-size: 13px;
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.c-field textarea {
  resize: vertical;
  min-height: 110px;
}

.c-field input::placeholder,
.c-field textarea::placeholder {
  color: var(--muted-2);
}

.c-field input:focus,
.c-field textarea:focus {
  border-color: var(--ink);
  box-shadow: 0 0 0 3px rgba(20, 21, 24, 0.05);
}

.c-field.invalid input,
.c-field.invalid textarea {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(232, 80, 58, 0.08);
}

.error-text {
  font-size: 10.5px;
  color: var(--accent-deep);
}

.form-submit {
  width: 100%;
  margin-top: 4px;
  height: 50px;
}

.form-submit:disabled {
  opacity: 0.6;
  cursor: wait;
}

.form-privacy {
  font-size: 10.5px;
  line-height: 1.6;
  color: var(--muted-2);
}

.form-success {
  min-height: 380px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 20px;
}

.success-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: var(--accent-soft);
  color: var(--accent-deep);
  margin-bottom: 22px;
}

.form-success h3 {
  font-size: 26px;
  letter-spacing: -0.03em;
}

.form-success p {
  margin-top: 10px;
  max-width: 36ch;
  font-size: 13.5px;
  line-height: 1.7;
  color: var(--muted);
}

.form-success .btn {
  margin-top: 24px;
}

.success-note {
  margin-top: 30px;
  font-size: 8px;
  letter-spacing: 0.16em;
  color: var(--muted-2);
}

.form-fade-enter-active,
.form-fade-leave-active {
  transition: opacity 0.25s ease, transform 0.3s var(--ease);
}

.form-fade-enter-from,
.form-fade-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

@container canvas (max-width: 960px) {
  .contact-grid {
    grid-template-columns: 1fr;
    gap: 40px;
  }
  .contact-intro h2 {
    max-width: 100%;
  }
}

@container canvas (max-width: 540px) {
  .form-row {
    grid-template-columns: 1fr;
  }
  .contact-info {
    grid-template-columns: 1fr;
  }
  .form-card {
    padding: 22px 18px;
  }
}
</style>

<template>
  <div class="site-toast-host">
    <TransitionGroup name="site-toast">
      <div v-for="toast in toasts" :key="toast.id" class="site-toast">
        <span class="site-toast-icon">
          <UiIcon :name="toast.type === 'error' ? 'close' : 'check'" :size="12" />
        </span>
        <span>{{ toast.text }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import UiIcon from '~/components/UiIcon.vue'
import { useToast } from '~/composables/useToast'

const { toasts } = useToast()
</script>

<style scoped>
.site-toast-host {
  position: fixed;
  right: 22px;
  bottom: 22px;
  z-index: 9000;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
  pointer-events: none;
}

.site-toast {
  display: flex;
  align-items: center;
  gap: 9px;
  max-width: min(360px, calc(100vw - 40px));
  padding: 10px 15px 10px 11px;
  border-radius: 999px;
  background: var(--ink);
  color: var(--paper);
  font-size: 12px;
  font-weight: 500;
  box-shadow: 0 14px 38px rgba(20, 21, 24, 0.28);
}

.site-toast-icon {
  width: 20px;
  height: 20px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  background: rgba(246, 242, 234, 0.12);
  color: #a9d0a9;
  flex: none;
}

.site-toast-enter-active,
.site-toast-leave-active {
  transition: opacity 0.3s ease, transform 0.35s var(--ease);
}

.site-toast-enter-from,
.site-toast-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.97);
}

@media (max-width: 600px) {
  .site-toast-host {
    right: 12px;
    bottom: 12px;
  }
}
</style>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="modal-overlay" @click.self="$emit('close')">
        <div class="modal-card" role="dialog" aria-modal="true">
          <button class="modal-close" aria-label="Close" @click="$emit('close')">
            <UiIcon name="close" :size="18" />
          </button>
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, watch } from 'vue'
import UiIcon from './UiIcon.vue'

const props = withDefaults(
  defineProps<{
    open?: boolean
  }>(),
  { open: false }
)

const emit = defineEmits<{ close: [] }>()

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.open) emit('close')
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})

watch(
  () => props.open,
  (open) => {
    document.body.style.overflow = open ? 'hidden' : ''
  }
)
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 5000;
  background: rgba(20, 21, 24, 0.55);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.modal-card {
  position: relative;
  width: min(920px, 100%);
  max-height: min(86vh, 760px);
  overflow-y: auto;
  background: var(--paper-light);
  border-radius: clamp(18px, 2.4vw, 28px);
  box-shadow: 0 40px 100px rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(246, 242, 234, 0.2);
}

.modal-close {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 3;
  width: 40px;
  height: 40px;
  display: grid;
  place-items: center;
  border: 1px solid var(--line);
  border-radius: 50%;
  background: rgba(251, 249, 245, 0.86);
  backdrop-filter: blur(8px);
  color: var(--ink);
  transition: background 0.2s ease, color 0.2s ease, transform 0.3s var(--ease);
}

.modal-close:hover {
  background: var(--ink);
  color: var(--paper);
  transform: rotate(90deg);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-active .modal-card,
.modal-leave-active .modal-card {
  transition: transform 0.45s var(--ease), opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-card,
.modal-leave-to .modal-card {
  transform: translateY(24px) scale(0.98);
  opacity: 0;
}
</style>

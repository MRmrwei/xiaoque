<template>
  <svg
    class="ui-icon"
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    :stroke-width="strokeWidth"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
    v-html="markup"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    name?: string
    size?: number | string
    strokeWidth?: number | string
  }>(),
  { name: 'spark', size: 20, strokeWidth: 1.7 }
)

const icons: Record<string, string> = {
  spark:
    '<path d="M12 2.8l2.05 5.7 5.7 2.05-5.7 2.05L12 18.3l-2.05-5.7-5.7-2.05 5.7-2.05L12 2.8Z"/><path d="M19 15.6l.75 2.1 2.1.75-2.1.75L19 21.3l-.75-2.1-2.1-.75 2.1-.75L19 15.6Z"/>',
  blocks:
    '<rect x="3" y="3" width="7.5" height="7.5" rx="1.6"/><rect x="13.5" y="3" width="7.5" height="7.5" rx="1.6"/><rect x="3" y="13.5" width="7.5" height="7.5" rx="1.6"/><rect x="13.5" y="13.5" width="7.5" height="7.5" rx="1.6"/>',
  globe:
    '<circle cx="12" cy="12" r="9"/><path d="M3 12h18"/><path d="M12 3c2.5 2.4 4 5.6 4 9s-1.5 6.6-4 9c-2.5-2.4-4-5.6-4-9s1.5-6.6 4-9Z"/>',
  search: '<circle cx="11" cy="11" r="7"/><path d="M20.5 20.5 16.2 16.2"/>',
  history:
    '<path d="M3 12a9 9 0 1 0 2.64-6.36L3 8"/><path d="M3 3v5h5"/><path d="M12 7v5l3.2 1.9"/>',
  image:
    '<rect x="3" y="4" width="18" height="16" rx="2.2"/><circle cx="8.5" cy="9.5" r="1.5"/><path d="M21 15.5 16.5 11 5 21"/>',
  shield:
    '<path d="M12 3l7 3v6c0 4.4-3 7.6-7 9-4-1.4-7-4.6-7-9V6l7-3Z"/><path d="M9 12l2 2 4-4"/>',
  zap: '<path d="M13 2 4 14h6l-1 8 9-12h-6l1-8Z"/>',
  brush:
    '<path d="M14.5 4.5l5 5L9 20H4v-5L14.5 4.5Z"/><path d="M12 7l5 5"/>',
  layers:
    '<path d="M12 3l9 5-9 5-9-5 9-5Z"/><path d="M3 13l9 5 9-5"/><path d="M3 17l9 5 9-5"/>',
  book:
    '<path d="M4 4.5A2.5 2.5 0 0 1 6.5 2H20v20H6.5A2.5 2.5 0 0 0 4 19.5v-15Z"/><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/>',
  users:
    '<circle cx="9" cy="8" r="3.4"/><path d="M2.6 20c0-3.5 2.9-6 6.4-6s6.4 2.5 6.4 6"/><path d="M16.8 5.2a3.4 3.4 0 0 1 0 6.4"/><path d="M18.2 14.6c2 .7 3.3 2.4 3.3 4.4"/>',
  type: '<path d="M4 6V4h16v2"/><path d="M12 4v16"/><path d="M8.5 20h7"/>',
  quote:
    '<path d="M9 7H5v5h4v5H4"/><path d="M20 7h-4v5h4v5h-5"/>',
  chart: '<path d="M3 20h18"/><path d="M6 20v-6"/><path d="M11 20V8"/><path d="M16 20V4"/><path d="M21 20v-9"/>',
  play: '<circle cx="12" cy="12" r="9"/><path d="M10 8.5l6 3.5-6 3.5v-7Z"/>',
  form:
    '<rect x="3" y="4" width="18" height="16" rx="2.2"/><path d="M7 9h10"/><path d="M7 13h6"/>',
  star:
    '<path d="M12 3l2.7 5.5 6.1.9-4.4 4.3 1 6.1-5.4-2.9-5.4 2.9 1-6.1L3.2 9.4l6.1-.9L12 3Z"/>',
  check: '<path d="M4 12.5l5 5L20 6.5"/>',
  arrow: '<path d="M5 12h14"/><path d="M13 6l6 6-6 6"/>',
  'arrow-up-right': '<path d="M7 17 17 7"/><path d="M8 7h9v9"/>',
  plus: '<path d="M12 5v14"/><path d="M5 12h14"/>',
  minus: '<path d="M5 12h14"/>',
  drag: '<circle cx="9" cy="6" r="1.1" fill="currentColor" stroke="none"/><circle cx="15" cy="6" r="1.1" fill="currentColor" stroke="none"/><circle cx="9" cy="12" r="1.1" fill="currentColor" stroke="none"/><circle cx="15" cy="12" r="1.1" fill="currentColor" stroke="none"/><circle cx="9" cy="18" r="1.1" fill="currentColor" stroke="none"/><circle cx="15" cy="18" r="1.1" fill="currentColor" stroke="none"/>',
  trash:
    '<path d="M4 7h16"/><path d="M9 7V4h6v3"/><path d="M6 7l1 13h10l1-13"/><path d="M10 11v5"/><path d="M14 11v5"/>',
  copy:
    '<rect x="8" y="8" width="12" height="12" rx="2"/><path d="M16 8V6a2 2 0 0 0-2-2H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h2"/>',
  eye: '<path d="M2.5 12S6 6 12 6s9.5 6 9.5 6-3.5 6-9.5 6-9.5-6-9.5-6Z"/><circle cx="12" cy="12" r="3"/>',
  clock: '<circle cx="12" cy="12" r="9"/><path d="M12 7v5l3.4 2"/>',
  layout:
    '<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18"/><path d="M9 9v12"/>',
  edit:
    '<path d="M12 20h9"/><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5Z"/>',
  more: '<circle cx="5" cy="12" r="1.25" fill="currentColor" stroke="none"/><circle cx="12" cy="12" r="1.25" fill="currentColor" stroke="none"/><circle cx="19" cy="12" r="1.25" fill="currentColor" stroke="none"/>',
  close: '<path d="M6 6l12 12"/><path d="M18 6 6 18"/>',
  chevron: '<path d="M6 9l6 6 6-6"/>',
  'chevron-right': '<path d="M9 6l6 6-6 6"/>',
  monitor:
    '<rect x="2" y="4" width="20" height="13" rx="2"/><path d="M8 21h8"/><path d="M12 17v4"/>',
  tablet:
    '<rect x="5" y="2" width="14" height="20" rx="2"/><path d="M11 18h2"/>',
  phone:
    '<rect x="7" y="2" width="10" height="20" rx="2.2"/><path d="M11 18h2"/>',
  settings:
    '<path d="M4 6h16"/><path d="M4 12h16"/><path d="M4 18h16"/><circle cx="9" cy="6" r="2" fill="var(--paper, #f6f2ea)"/><circle cx="15" cy="12" r="2" fill="var(--paper, #f6f2ea)"/><circle cx="7" cy="18" r="2" fill="var(--paper, #f6f2ea)"/>',
  home:
    '<path d="M3 10.5 12 3l9 7.5"/><path d="M5 9.5V21h14V9.5"/><path d="M9 21v-6h6v6"/>',
  pages:
    '<path d="M6 2h8l4 4v16H6V2Z"/><path d="M14 2v5h5"/><path d="M9 12h6"/><path d="M9 16h6"/>',
  menu: '<path d="M4 6h16"/><path d="M4 12h16"/><path d="M4 18h16"/>',
  logout:
    '<path d="M9 21H4V3h5"/><path d="M16 17l5-5-5-5"/><path d="M21 12H9"/>',
  undo: '<path d="M9 14l-5-5 5-5"/><path d="M4 9h9a7 7 0 0 1 0 14h-2"/>',
  redo: '<path d="M15 14l5-5-5-5"/><path d="M20 9h-9a7 7 0 0 0 0 14h2"/>',
  publish:
    '<path d="M12 16V4"/><path d="M7 9l5-5 5 5"/><path d="M4 16v4h16v-4"/>',
  link:
    '<path d="M10 13a5 5 0 0 0 7.1 0l2-2A5 5 0 0 0 12 4l-1.1 1.1"/><path d="M14 11a5 5 0 0 0-7.1 0l-2 2A5 5 0 0 0 12 20l1.1-1.1"/>',
  filter: '<path d="M3 5h18l-7 8v6l-4 2v-8L3 5Z"/>',
  download: '<path d="M12 3v12"/><path d="M7 10l5 5 5-5"/><path d="M4 19h16"/>',
  upload: '<path d="M12 17V4"/><path d="M7 9l5-5 5 5"/><path d="M4 20h16"/>',
  grid:
    '<rect x="3" y="3" width="7" height="7" rx="1.4"/><rect x="14" y="3" width="7" height="7" rx="1.4"/><rect x="3" y="14" width="7" height="7" rx="1.4"/><rect x="14" y="14" width="7" height="7" rx="1.4"/>',
  mail: '<rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 7.5 12 13l9-5.5"/>',
  pin:
    '<path d="M12 21s7-6.2 7-11a7 7 0 1 0-14 0c0 4.8 7 11 7 11Z"/><circle cx="12" cy="10" r="2.5"/>',
  external:
    '<path d="M14 4h6v6"/><path d="M20 4 10 14"/><path d="M18 14v6H4V6h6"/>',
  refresh: '<path d="M20 12a8 8 0 1 1-2.34-5.66"/><path d="M20 4v5h-5"/>',
  globe2:
    '<circle cx="12" cy="12" r="9"/><path d="M3.5 9h17"/><path d="M3.5 15h17"/><path d="M12 3c2 2.4 3 5.5 3 9s-1 6.6-3 9c-2-2.4-3-5.5-3-9s1-6.6 3-9Z"/>',
  code: '<path d="M8 8l-5 4 5 4"/><path d="M16 8l5 4-5 4"/><path d="M13.5 4l-3 16"/>',
  sparkles:
    '<path d="M12 3l1.8 5.2L19 10l-5.2 1.8L12 17l-1.8-5.2L5 10l5.2-1.8L12 3Z"/><path d="M19 16l.7 2 2 .7-2 .7-.7 2-.7-2-2-.7 2-.7.7-2Z"/>',
  panel:
    '<rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18"/>',
  send: '<path d="M22 2 11 13"/><path d="M22 2 15 22l-4-9-9-4 20-7Z"/>',
  'arrow-left': '<path d="M19 12H5"/><path d="M11 6l-6 6 6 6"/>',
  filter2: '<path d="M4 6h16"/><path d="M7 12h10"/><path d="M10 18h4"/>',
  bell:
    '<path d="M18 8a6 6 0 0 0-12 0c0 6.2-2.5 7.5-2.5 7.5h17S18 14.2 18 8Z"/><path d="M13.7 20a2 2 0 0 1-3.4 0"/>',
  save:
    '<path d="M5 3h11l3 3v15H5V3Z"/><path d="M8 3v6h8V3"/><path d="M8 21v-8h8v8"/>',
  help:
    '<circle cx="12" cy="12" r="9"/><path d="M9.2 9a3 3 0 0 1 5.8 1c0 2-3 2.3-3 4.2"/><circle cx="12" cy="17.4" r=".8" fill="currentColor" stroke="none"/>',
  lock:
    '<rect x="4.5" y="10" width="15" height="11" rx="2"/><path d="M8 10V7a4 4 0 0 1 8 0v3"/>',
  pause:
    '<rect x="6.5" y="5" width="4" height="14" rx="1" fill="currentColor" stroke="none"/><rect x="13.5" y="5" width="4" height="14" rx="1" fill="currentColor" stroke="none"/>'
}

const markup = computed(() => icons[props.name] || icons.spark)
</script>

<style scoped>
.ui-icon {
  flex: none;
  display: block;
}
</style>

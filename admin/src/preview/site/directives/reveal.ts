import type { Directive } from 'vue'

/** Preview canvas: content is visible immediately; no scroll animations. */
const reveal: Directive<HTMLElement, { delay?: number } | undefined> = {
  mounted(el) {
    el.classList.add('is-visible')
  }
}

export default reveal

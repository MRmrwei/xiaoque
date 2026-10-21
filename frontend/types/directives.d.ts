import type { Directive } from 'vue'

declare module 'vue' {
  interface GlobalDirectives {
    vReveal: Directive<HTMLElement, { delay?: number } | undefined>
  }
}

export {}

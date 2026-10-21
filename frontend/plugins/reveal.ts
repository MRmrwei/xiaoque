import type { Directive } from 'vue'

interface RevealOptions {
  delay?: number
}

type RevealElement = HTMLElement & {
  _revealObserver?: IntersectionObserver
}

const reveal: Directive<HTMLElement, RevealOptions | undefined> = {
  getSSRProps() {
    return {}
  },
  mounted(el, binding) {
    const options = binding.value ?? {}

    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      el.classList.add('is-visible')
      return
    }

    // Avoid hiding content that is already on screen when hydration finishes.
    const rect = el.getBoundingClientRect()
    if (rect.top < window.innerHeight * 0.95) {
      el.classList.add('is-visible')
      return
    }

    el.classList.add('reveal')
    if (options.delay) el.style.transitionDelay = `${options.delay}ms`

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry?.isIntersecting) {
          el.classList.add('is-visible')
          observer.disconnect()
        }
      },
      { threshold: 0.08, rootMargin: '0px 0px -40px 0px' }
    )
    observer.observe(el)
    ;(el as RevealElement)._revealObserver = observer
  },
  unmounted(el) {
    ;(el as RevealElement)._revealObserver?.disconnect()
  }
}

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.directive('reveal', reveal)
})

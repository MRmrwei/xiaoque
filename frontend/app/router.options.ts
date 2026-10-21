import type { RouterConfig } from '@nuxt/schema'
import { createMemoryHistory, createWebHashHistory, createWebHistory } from 'vue-router'

/**
 * Normal deployments use clean URLs. A static share build opened from disk
 * (file://) uses hash routing so navigation works without a web server.
 */
export default <RouterConfig>{
  history: (base) => {
    // Nuxt creates the router during SSR as well; web history touches window.
    if (typeof window === 'undefined') {
      return createMemoryHistory(base)
    }
    if (window.location.protocol === 'file:') {
      // Without an explicit base, Vue Router treats the file path as the base
      // and the initial route becomes /G:/.../index.html.
      return createWebHashHistory('/')
    }
    return createWebHistory(base)
  }
}

/**
 * When the single-file share build is opened directly from disk, the browser
 * URL has no hash and Nuxt/Vue Router can interpret the file pathname as the
 * initial route. Redirect that case to the site root.
 *
 * No effect when the site is served over http/https.
 */
export default defineNuxtPlugin(async () => {
  if (typeof window === 'undefined' || window.location.protocol !== 'file:') return

  const router = useRouter()
  await router.isReady()

  const path = router.currentRoute.value.path
  if (path.includes('index.html') || /^\/[A-Za-z]:\//.test(path)) {
    await router.replace('/')
  }
})

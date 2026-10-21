import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

const resolvePath = (path: string) => fileURLToPath(new URL(path, import.meta.url))

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), 'VITE_')
  const proxyTarget = env.VITE_API_PROXY_TARGET || 'http://127.0.0.1:8888'
  const proxy = {
    '/api': {
      target: proxyTarget,
      changeOrigin: true
    }
  }

  return {
    // Hash routing + relative assets keep the demo portable (static host or file://).
    base: './',
    plugins: [vue()],
    resolve: {
      alias: [
        { find: '@', replacement: resolvePath('./src') },
        { find: '~', replacement: resolvePath('./src/preview/site') }
      ],
      dedupe: ['vue', 'vue-router']
    },
    server: {
      port: 5174,
      host: '0.0.0.0',
      proxy
    },
    preview: {
      port: 4174,
      host: '0.0.0.0',
      proxy
    },
    build: {
      rollupOptions: {
        input: {
          admin: resolvePath('./index.html'),
          preview: resolvePath('./preview.html')
        }
      }
    }
  }
})

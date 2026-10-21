import fs from 'node:fs'
import path from 'node:path'

const root = path.resolve('share-output/public')
const htmlPath = path.join(root, 'index.html')
const outDir = path.resolve('dist-share')
const outFile = path.join(outDir, 'index.html')

if (!fs.existsSync(htmlPath)) {
  console.error('share-output/public/index.html not found. Run: NUXT_SHARE=1 npx nuxt generate')
  process.exit(1)
}

let html = fs.readFileSync(htmlPath, 'utf8')

function resolveAsset(url) {
  const clean = url.replace(/^\.?\//, '').split('?')[0]
  return path.join(root, clean)
}

// Inline local stylesheets; keep external Google Fonts links untouched.
html = html.replace(/<link\b[^>]*rel="stylesheet"[^>]*>/g, (tag) => {
  const href = /href="([^"]+)"/.exec(tag)?.[1]
  if (!href || /^https?:\/\//.test(href)) return tag
  const file = resolveAsset(href)
  if (!fs.existsSync(file)) {
    console.warn('Missing stylesheet:', href)
    return tag
  }
  const css = fs.readFileSync(file, 'utf8')
  return `<style>${css.replace(/<\/style/gi, '<\\/style')}</style>`
})

// Modulepreload hints are pointless once the entry is inline.
html = html.replace(/<link\b[^>]*rel="modulepreload"[^>]*>\s*/g, '')

// Inline the client entry.
html = html.replace(
  /<script\b([^>]*)\bsrc="([^"]+)"([^>]*)><\/script>/g,
  (tag, before, src, after) => {
    if (/^https?:\/\//.test(src)) return tag
    const file = resolveAsset(src)
    if (!fs.existsSync(file)) {
      console.warn('Missing script:', src)
      return tag
    }
    const code = fs.readFileSync(file, 'utf8').replace(/<\/script/gi, '<\\/script')
    return `<script${before}${after}>${code}</script>`
  }
)

// Double-clicking a local file has no hash. Wait for the SPA to mount, then
// navigate to the root. Setting the hash after mount triggers Vue Router's
// hashchange listener, which is reliable on file://.
const autoRoute = `<script>(function(){if(location.protocol!=='file:')return;var n=0,t=setInterval(function(){var root=document.querySelector('#__nuxt');if(root&&root.children.length){clearInterval(t);if(!location.hash||location.hash.indexOf('index.html')!==-1||document.querySelector('.not-found')){location.hash='/'}return}if(++n>60){clearInterval(t)}},50)})()</script>`
html = html.replace('</body>', `${autoRoute}</body>`)

fs.mkdirSync(outDir, { recursive: true })
fs.writeFileSync(outFile, html, 'utf8')
console.log(`Wrote ${outFile} (${(fs.statSync(outFile).size / 1024).toFixed(0)} KB)`)

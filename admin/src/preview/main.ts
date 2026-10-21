import { createApp } from 'vue'
import PreviewApp from './PreviewApp.vue'
import reveal from './site/directives/reveal'
import './site/assets/base.css'
import './preview.css'

const app = createApp(PreviewApp)
app.directive('reveal', reveal)
app.mount('#preview-app')

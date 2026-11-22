import { createApp } from 'vue'
import App from './App.vue'
import './assets/styles.css'

if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('./sw.js').catch(() => {})
  })
}

createApp(App).mount('#app')

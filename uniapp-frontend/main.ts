import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import uview from 'uview-ui'

export function createApp() {
  const app = createSSRApp(App)
  const pinia = createPinia()

  app.use(pinia)
  app.use(uview)

  return {
    app,
    pinia
  }
}

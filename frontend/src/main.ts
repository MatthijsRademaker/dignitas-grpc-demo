import { createApp } from 'vue'
import App from './App.vue'
import vuetify0 from './plugins/vuetify0'
import '@fontsource/fira-code/latin-400.css'
import '@fontsource/fira-code/latin-600.css'
import './styles/fonts.css'
import './styles/main.css'

const app = createApp(App)

vuetify0(app)

app.mount('#app')

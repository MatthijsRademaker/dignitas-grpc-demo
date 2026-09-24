import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    // The browser only ever talks to the BFF, over plain HTTP.
    proxy: { '/api': process.env.BFF_URL ?? 'http://localhost:5080' },
  },
})

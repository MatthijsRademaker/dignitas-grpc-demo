import { createThemePlugin } from '@vuetify/v0'
import type { App } from 'vue'

export default function vuetify0 (app: App) {
  app.use(
    createThemePlugin({
      default: 'light',
      themes: {
        light: {
          dark: false,
          colors: {
            'primary': '#733FCB',
            'on-primary': '#ffffff',
            'surface': '#ffffff',
            'on-surface': '#1f1b2e',
            'background': '#f5f3fa',
            'muted': '#6b6680',
            'border': '#e4dff0',
            'success': '#15803d',
            'warning': '#b45309',
            'error': '#b91c1c',
          },
        },
      },
    }),
  )
}

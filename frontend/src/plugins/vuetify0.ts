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
            // Gold is for artwork, borders and glows only, never for text.
            'gold-50': '#fdf8e6',
            'gold-100': '#faedc0',
            'gold-200': '#f6d77a',
            'gold-300': '#eec455',
            'gold-400': '#e0ad32',
            'gold-500': '#c9971c',
            'gold-600': '#a67a0f',
            'gold-700': '#7a5a00',
            'gold-800': '#5c4300',
            'gold-900': '#3d2c00',
            // The one gold that may be text: 6.4:1 on surface, 5.8:1 on background.
            'bronze': '#7a5a00',
          },
        },
      },
    }),
  )
}

import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'

import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    tailwindcss(),
    AutoImport({
      imports: [
        'vue',
        'vue-router',
        'pinia',
        '@vueuse/core',
        {
          '@/utils': ['formatDate', 'handleApiError'],
        },
        {
          '@yuelioi/toast': ['toast'], // ✅ 指定 toast 函数
        },
      ],
      dts: 'src/auto-imports.d.ts', // 自动生成类型声明
      dirs: ['./src/stores', './src/composables', './src/models/', './src/api/'],

      eslintrc: {
        enabled: true, // 生成 ESLint 配置，避免报 no-undef
        filepath: './.eslintrc-auto-import.json',
        globalsPropValue: true,
      },
    }),
    Components({
      dirs: ['src/components'], // 自动导入的组件目录
      extensions: ['vue'],
      deep: true, // 支持子目录
      dts: 'src/components.d.ts', // 生成类型声明文件
      resolvers: [],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:9000',
        changeOrigin: true,
      },
    },
  },
})

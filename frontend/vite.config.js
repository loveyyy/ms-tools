/*
 * @Description: 
 */
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import VueSetupExtend from 'vite-plugin-vue-setup-extend'
import { visualizer } from 'rollup-plugin-visualizer'
import viteCompression from 'vite-plugin-compression'//Rollup 的虚拟模块
import { viteExternalsPlugin } from 'vite-plugin-externals';
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'


import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'

import path from 'path'
export default defineConfig(({ mode }) => {
  const { VITE_look_visualizer = false, VITE_APP_ENV } = loadEnv(mode, process.cwd());
  return {
    plugins: [
      vueJsx(),
      vue(),
      AutoImport({
        imports: [
          {
            'naive-ui': ['useDialog', 'useMessage', 'useNotification', 'useLoadingBar', 'createDiscreteApi']
          }
        ],
        // vue3 组件 js 语句中自动引入组件
        resolvers: [NaiveUiResolver()]
      }),
      VueSetupExtend(),
      viteExternalsPlugin({
        '@wangeditor/editor': 'wangEditor',
        'echarts': 'echarts'
      }),
      // vue3 组件中自动引入组件
      Components({
        resolvers: [NaiveUiResolver()]
      })
    ],
    resolve: {
      alias: {
        // @ts-ignore
        '@': path.resolve(__dirname, 'src')
      }
    },
    css: {
      preprocessorOptions: {
        less: {
          javascriptEnabled: true
        }
      }
    },
    server: {
      host: '0.0.0.0',
      port: 4417,
      open: false,
      https: false,
      proxy: {}
    },
    build: {
      minify: 'terser',
      terserOptions: {
        compress: {
          /* 清除console */
          drop_console: VITE_APP_ENV === 'production',
          /* 清除debugger */
          drop_debugger: VITE_APP_ENV === 'production'
        }
      },
      rollupOptions: {
        plugins: [
          // 判断是否使用
          VITE_look_visualizer && visualizer({
            open: true, // 是否在打包完成后自动打开生成的报告文件（默认为 `false`）
            filename: 'dist/stats.html', // 指定报告文件的输出路径（默认为 `stats.html`）
            gzipSize: true, // 是否显示 Gzipped 文件的大小（默认为 `false`）
            brotliSize: true // 是否显示 Brotli 压缩后的文件大小（默认为 `false`）
          }),
          // gzip压缩 生产环境生成 .gz 文件
          viteCompression({
            verbose: true,
            disable: false,
            threshold: 10240,
            algorithm: 'gzip',
            ext: '.gz'
          })
        ],
        output: {
          chunkFileNames: 'static/js/[name]-[hash].js',
          entryFileNames: 'static/js/[name]-[hash].js',
          assetFileNames: 'static/[ext]/[name]-[hash].[ext]',
          manualChunks(id) {
            //静态资源分拆打包
            if (id.includes('node_modules')) {
              return id.toString().split('node_modules/')[1].split('/')[0].toString()
            }
          }
        }
      },
      sourcemap: mode === 'production'
    }
  }
})
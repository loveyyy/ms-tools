/*
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-03-18 20:55:56
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-07 17:30:42
 * @Description: 
 */
import {createApp} from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import { router } from "./router/index";
import { createDiscreteApi } from 'naive-ui'


const app = createApp(App)

// 创建 Naive UI 的离散 API
const { message, dialog, notification, loadingBar } = createDiscreteApi([
  'message',
  'dialog', 
  'notification',
  'loadingBar'
])

// 将实例挂载到 window 对象，以便在非组件上下文中使用
window.$naive = {
  $message: message,
  $dialog: dialog,
  $notification: notification,
  $loadingBar: loadingBar
}

app.use(naive)
app.use(router)

app.mount('#app')
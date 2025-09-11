<!--
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-09-11 00:00:00
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-11 16:18:02
 * @Description: 项目首页，Markdown 介绍与使用说明
-->
<template>
  <div class="home-container">
    <n-card size="large" embedded>
      <div class="markdown-body" v-html="html"></div>
    </n-card>
  </div>

</template>

<script setup>
import { computed } from 'vue'
import { renderMarkdown } from '@/utils/md'
import '@/assets/css/github-markdown-light.css'
import 'highlight.js/styles/github.css'
import "katex/dist/katex.min.css";

const md = `
# MS-TOOLS

> 一个用于学习与技术研究的桌面工具（Wails + Vue3 + Naive UI）。
##
### 技术栈
  - 前端：Vue3、Vite、Naive UI、Vue-Router
  - Markdown 渲染：marked + highlight.js + KaTeX
  - 桌面容器：Wails（Go + Web 前端）

##
### 功能概览
  - 小红书：检索并查看笔记详情、用户信息；
  - 抖音：检索并查看视频详情。

##
### 使用说明
#### 运行/编译项目
- wails dev/build。具体请参考:https://wails.io/docs/gettingstarted/running
#### 关于小红书不登陆的局限性
- 不登陆无法进行搜索.
- 搜索后只能出现一页数据，且无法获取到笔记ID，意味着无法查看详情.
#### 关于小红书web_seesion获取方式：
- 使用浏览器打开小红书APP，登录后。
- 在开发者工具中查看Cookie,找到web_session字段。
- 将web_session字段复制到页面上，点击“设置登录Session”按钮即可。
##
### 声明
- 本项目仅供学习与技术研究使用，禁止用于任何商业用途；
- 使用者需遵守目标平台的服务条款、用户协议与相关法律；
- 未经授权请勿抓取、存储、传播他人受保护的数据或内容；
- 开发者不对使用者的任何不当使用承担责任，若涉及第三方权利，请立即停止并删除相关数据。

##
### TODO
- [ ] 集成gochorme 拉起页面进行登录或者使用手机号进行登录并且持久化至本地
- [ ] 集成轻量聊天与 AI 助手（合规 API，本地优先）
- [ ] 增加一键清理本地日志入口
- [ ] 小红书与抖音检索历史本地化存储
`

const html = computed(() => renderMarkdown(md))
</script>

<style scoped>
.home-container {
  padding: 16px;
}
</style>

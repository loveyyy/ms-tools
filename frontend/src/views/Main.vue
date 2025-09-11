<!--
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-03-20 21:14:07
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-11 15:46:07
 * @Description: 
-->
<template>
  <div class="main-container">
    <n-layout style="height: 100%">
      <n-layout-header style="height: 64px" bordered>
        <n-flex style="height: 100%" align="center">
          <div class="head-title">MS-TOOLS</div>
        </n-flex>
      </n-layout-header>
      <n-layout has-sider style="height: calc(100% - 64px)">
        <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="160" :collapsed="collapsed"
          show-trigger @collapse="collapsed = true" @expand="collapsed = false">
          <n-menu style="height: 100%" :collapsed="collapsed" default-value="home" :collapsed-width="64"
            :collapsed-icon-size="22" :options="menuOptions" @update:value="handleMenuChange" />
        </n-layout-sider>
        <n-layout-content style="position: relative">
          <router-view v-slot="{ Component, route }">
            <transition name="fade-transform" mode="out-in">
              <keep-alive :include="keepAliveRoutes">
                <component :is="Component" :key="route.path" />
              </keep-alive>
            </transition>
          </router-view>
        </n-layout-content>
      </n-layout>
    </n-layout>
  </div>

</template>


<script setup>
import { h, ref, shallowRef } from 'vue'
import { Analytics, Book, ChatboxEllipses, Settings, Videocam } from '@vicons/ionicons5'
import { useRouter } from 'vue-router';

const router = useRouter();

const keepAliveRoutes = ['home', 'xhs', 'dy']

const menuOptions = shallowRef([
  {
    label: '首页',
    key: 'home',
    icon: renderMenuIcon(ChatboxEllipses),
  },
  {
    label: '音视频',
    key: "voice",
    icon: renderMenuIcon(Analytics),
    children: [
      {
        label: '小红书',
        key: "xhs",
        icon: renderMenuIcon(Book),
      },
      {
        label: '抖音',
        key: "dy",
        icon: renderMenuIcon(Videocam),
      },
    ]
  },

  // {
  //   label: '设置',
  //   key: "setting",
  //   icon: renderMenuIcon(Settings),
  // }
])

const collapsed = ref(false)

function handleMenuChange (key) {
  router.push(key)
}


function renderMenuIcon (name) {
  return () => h(NIcon, null, { default: () => h(name) });
}
</script>

<style scoped>
.main-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.head-title {
  font-size: 20px;
  font-weight: bold;
  text-align: center;
  width: 160px;
}
</style>
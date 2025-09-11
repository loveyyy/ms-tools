<!--
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-03-18 20:55:56
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-09 22:08:33
 * @Description: 
-->
<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme-overrides="themeOverrides" style="height: 100vh;">
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-loading-bar-provider>
            <router-view />
          </n-loading-bar-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script>
import { dateZhCN, zhCN, NConfigProvider } from 'naive-ui'
import { defineComponent, ref } from "vue";

export default defineComponent({
  setup () {
    return {
      zhCN,
      dateZhCN,
      locale: ref(null),
      dateLocale: ref(null),
      themeOverrides: {
        common: {
          primaryColor: '#409eff',
          primaryColorHover: '#66b1ff',
          primaryColorPressed: '#85ceff',
          primaryColorSuppl: '#5cadff',
        },
      }
    }
  }
});
</script>


<style></style>

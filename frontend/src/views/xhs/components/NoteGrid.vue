<!--
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-09-10 14:54:51
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-10 22:03:10
 * @Description: 
-->
<template>
  <n-row :gutter="[12, 12]" style="width: 100%;">
    <n-col :span="8" v-for="(item, index) in list" :key="index" @click.stop="() => $emit('selectNote', item)"
      style="cursor: pointer;">
      <div style="position: relative;">
        <n-image width="100%" height="360px" style="border-radius: 10px"
          :src="proxyXhs(item.note_card.cover.url_default || item.note_card.cover.url)" object-fit="contain" preview-disabled />
        <n-icon color="white" size="20px" style="position: absolute;right: 10px;top: 10px;"
          v-if="(item.note_card.type) === 'video'">
          <PlayCircleOutline />
        </n-icon>
      </div>
      <n-ellipsis class="note_title" style="width: 100%;text-align: start">
        {{ item.note_card.display_title }}
      </n-ellipsis>
      <div style="display: flex;gap: 5px;align-items: center">
        <n-image width="20px" height="20px" style="border-radius: 10px;width: 20px;"
          :src="proxyXhs(item.note_card.user.avatar)" preview-disabled @click.stop="() => emitOpenUser(item)" />
        <n-text @click.stop="() => emitOpenUser(item)" style="cursor: pointer;">{{ item.note_card.user.nick_name
        }}</n-text>
        <div style="margin-left: auto;font-size: 13px">{{ item.note_card.interact_info.liked_count }}</div>
      </div>
    </n-col>
  </n-row>
</template>

<script setup>
import { PlayCircleOutline } from "@vicons/ionicons5"

const emits = defineEmits(['selectNote', 'openUser'])
const props = defineProps({
  list: {
    type: Array,
    default: () => []
  }
})

function emitOpenUser (item) {
  emits('openUser', item)
}

function proxyXhs (url) {
  if (!url) return url
  try {
    return `http://127.0.0.1:17892/proxy/xhs?u=${encodeURIComponent(url)}`
  } catch (e) {
    return url
  }
}
</script>

<style scoped>
.note_title {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  width: 100%;
  text-align: start;
}
</style>

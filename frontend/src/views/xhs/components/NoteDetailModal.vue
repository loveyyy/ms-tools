<template>
  <n-modal v-model:show="modelShow" :mask-closable="true" preset="card" style="width:80vw; height :80vh"
    :title="detail?.note_card?.title || '内容详情'">
    <div v-if="detail" class="card-content">
      <div style="height: 100%; display: flex;max-width: 50%;">
        <div v-if="detail.note_card.type === 'video'"
          style="width: 100%;height: 100%;flex: 1; background: #000; border-radius: 8px; overflow: hidden;">
          <video v-if="detail.note_card.type === 'video'"
            style="width: 100%; height: 100%; object-fit: fit; background: #000;" controls playsinline autoplay
            :src="proxyXhs(getVideoUrl(detail))" />
        </div>
        <div v-else>
          <n-carousel autoplay :show-arrow="detail.note_card.image_list && detail.note_card.image_list.length > 1">
            <n-carousel-item v-for="(item, idx) in detail.note_card.image_list" :key="idx">
              <n-image height="100%" width="100%" fit="cover" style="border-radius: 10px;width: 100%;height: 100%;"
                :src="proxyXhs(item.url_default)" />
            </n-carousel-item>
            <template #arrow="{ prev, next }"
              v-if="detail.note_card.image_list && detail.note_card.image_list.length > 1">
              <div class="custom-arrow-left">
                <n-icon @click="prev" size="30px" color="grey">
                  <ArrowBackCircle />
                </n-icon>
              </div>
              <div class="custom-arrow-right">
                <n-icon @click="next" size="30px" color="grey">
                  <ArrowForwardCircle />
                </n-icon>
              </div>
            </template>
          </n-carousel>
        </div>
      </div>
      <div style="display: flex; flex-direction: column;align-items: center;height: 100%;flex: 1;">
        <n-flex vertical style="height: 100%;width: 100%;">
          <div style="display: flex; align-items: center; gap: 8px;">
            <n-image width="40px" height="40px" style="border-radius: 50%;width: 40px;"
              :src="detail.note_card.user.avatar" preview-disabled />
            <div style="display: flex;flex-direction: column;">
              <n-text style="text-align: start;">{{ detail.note_card.user.nickname }}</n-text>
              <n-text depth="3" style="font-size: 12px;">{{ formatTime(detail.note_card.time) }}</n-text>
            </div>
          </div>
          <n-scrollbar style="height: 100%;">
            <div style="padding-right: 20px;font-weight: 400;font-size: 16px;"
              v-html="filterDesc(detail.note_card.desc)">
            </div>
            <n-flex warap style="margin-top: 20px;">
              <n-tag v-for="(item, index) in (detail.note_card.tag_list || [])" round size="small" type="primary"
                :key="index" :disabled="false">#{{ item.name }}</n-tag>
            </n-flex>
          </n-scrollbar>
          <n-space style="margin-top: auto;">
            <n-tag round>赞 {{ formatNum(detail.note_card.interact_info.liked_count) }}</n-tag>
            <n-tag round>收藏 {{ formatNum(detail.note_card.interact_info.collected_count) }}</n-tag>
            <n-tag round>评论 {{ formatNum(detail.note_card.interact_info.comment_count) }}</n-tag>
          </n-space>
        </n-flex>
      </div>
    </div>
    <div v-else style="padding: 20px; text-align: center; color: #888;">暂无详情</div>
  </n-modal>
</template>

<script setup>
import { ArrowBackCircle, ArrowForwardCircle } from '@vicons/ionicons5'
import { formatTime } from '@/utils'

const modelShow = defineModel('show', { type: Boolean, default: false })
const props = defineProps({
  detail: {
    type: Object,
    default: null
  }
})

function filterDesc (desc) {
  desc = (desc || '').replace(/#[^#]+#/g, '').replace(/#/g, '')
  return desc
}

function formatNum (n) {
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n || 0)
}

function getVideoUrl (item) {
  if (!item || !item.note_card || !item.note_card.video) return ""
  let streamItem = item.note_card.video.media.stream
  for (let key in streamItem) {
    let value = streamItem[key]
    if (value && value.length > 0) {
      if (value[0].master_url) return value[0].master_url
    }
  }
  return ''
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
.custom-arrow-left {
  display: flex;
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  left: 10px;
  z-index: 10;
}

.custom-arrow-right {
  display: flex;
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: 10px;
  z-index: 10;
}

.card-content {
  position: absolute;
  top: 67px;
  left: 0;
  height: calc(100% - 87px);
  display: flex;
  gap: 20px;
  padding: 0 20px;
}
</style>

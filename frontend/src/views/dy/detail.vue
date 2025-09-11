<template>
  <div class="main-container">
    <div style="display: flex; flex-direction: column; height: 100%; gap: 12px;">
      <n-flex align="center" style="gap: 10px;">
        <n-icon size="20px" @click="handleBack">
          <ArrowBack />
        </n-icon>

        <n-input style="max-width: 320px;border-radius: 8px;text-align: start" v-model:value="keyword" clearable
          placeholder="搜索用户或关键词">
          <template #suffix>
            <n-button text @click="refresh">搜索</n-button>
          </template>
        </n-input>
      </n-flex>

      <n-tabs v-model:value="searchType" type="segment" @update:value="refresh">
        <n-tab v-for="(item, idx) in typeList" :key="idx" :tab="item.name" :name="item.value" />
      </n-tabs>

      <n-back-top :bottom="40" :visibility-height="200" :show="showBackTop" @click="scrollToTop" />

      <n-scrollbar ref="scrollbarRef" style="height: 100%; position: relative;" @scroll="onScroll"
        v-if="videoList && videoList.length > 0">
        <n-spin size="small" stroke="#409eff" v-if="isRefresh" />
        <n-row :gutter="[12, 12]" style="width: 100%;">
          <n-col :span="6" v-for="(item, index) in videoList" :key="item.aweme_id || index" style="cursor: pointer;"
            @click="openPlayer(item)">
            <div style="position: relative;">
              <n-image width="100%" height="250px" style="border-radius: 10px" :src="proxy(item.cover)" preview-disabled
                fit="cover" />
              <n-icon color="white" size="20px" style="position: absolute;right: 10px;top: 10px;">
                <PlayCircleOutline />
              </n-icon>
            </div>
            <n-ellipsis class="note_title" style="width: 100%;text-align: start" line-clamp="2">
              {{ item.desc || '无标题' }}
            </n-ellipsis>
            <div style="display: flex;gap: 5px;align-items: center">
              <n-text>@{{ item.author || '抖音用户' }}</n-text>
              <div style="margin-left: auto;font-size: 13px">{{ formatNum(item.digg_count) }}</div>
            </div>
          </n-col>
        </n-row>
        <n-flex wrap align="center">
          <div v-if="loading"
            style="font-weight: bold;margin-top: 10px;width: 100%;display: flex;align-items: center;justify-content: center;gap: 10px;">
            <n-spin size="small" stroke="#409eff" />
            <span> 加载中...</span>
          </div>
          <div v-if="noMore" style="font-weight: bold;margin-top: 10px;width: 100%;">
            没有更多了 🤪
          </div>
        </n-flex>
      </n-scrollbar>
      <div v-else style="height: 100%;align-items: center;display: flex;justify-content: center;">
        <n-empty description="暂无数据">
          <template #extra>
            <n-button size="small" @click="refresh">刷新</n-button>
          </template>
        </n-empty>
      </div>
    </div>

    <n-modal v-model:show="showPlayer" :mask-closable="true" preset="card" :style="{ width: '960px' }"
      :title="currentVideo?.desc || '视频播放'">
      <div style="display: flex; gap: 16px; align-items: flex-start;">
        <div style="flex: 2; background: #000; border-radius: 8px; overflow: hidden;">
          <video v-if="currentVideo?.play_url" :src="proxy(currentVideo.play_url)"
            style="width: 100%; height: 520px; object-fit: contain; background: #000;" controls autoplay loop
            playsinline muted></video>
          <div v-else style="height: 520px; display: flex; align-items: center; justify-content: center; color: #999;">
            暂无可播放地址
          </div>
        </div>
        <div style="flex: 1; display: flex; flex-direction: column; gap: 10px;">
          <div style="font-size: 16px; font-weight: 600;">{{ currentVideo?.desc || '无标题' }}</div>
          <div style="display: flex; align-items: center; gap: 8px;">
            <n-image height="28px" style="border-radius: 50%;width: 28px;" :src="currentVideo?.avatar"
              preview-disabled />
            <n-flex>
              <n-text>{{ currentVideo?.author || '抖音用户' }}</n-text>
              <n-text depth="3">发布于：{{ formatTime(currentVideo?.create_time) }}</n-text>
            </n-flex>
          </div>
          <n-text depth="3">时长：{{ formatDuration(currentVideo?.duration) }}</n-text>
          <n-text depth="3">点赞：{{ formatNum(currentVideo?.digg_count) }}</n-text>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Search } from '../../../wailsjs/go/dy/Dy'
import { ArrowBack, PlayCircleOutline } from '@vicons/ionicons5'
import { formatTime } from '@/utils'

const route = useRoute()

const keyword = ref('')
const searchType = ref('video')
const typeList = ref([
  { name: '综合', value: 'normal' },
  { name: '视频', value: 'video' },
  { name: '用户', value: 'user' },
])

const searchParams = ref({
  search_channel: 'aweme_video_web',
  enable_history: '1',
  keyword: '',
  search_source: 'normal_search',
  query_correct_type: '1',
  is_filter_search: '0',
  from_group_id: '',
  disable_rs: '1',
  offset: "0",
  count: "10",
  need_filter_settings: "1",
  list_type: 'single',
  search_id: '111111'
})

const videoList = ref([])
const loading = ref(false)
const isRefresh = ref(false)
const noMore = ref(false)

const scrollbarRef = ref(null)
const showBackTop = ref(false)

const showPlayer = ref(false)
const currentVideo = ref(null)

onMounted(() => {
  const query = route.query
  if (!query.keyword) {
    history.back()
    return
  }
  keyword.value = String(query.keyword || '')
  searchParams.value.keyword = query.keyword
  refresh()
})

function handleBack () {
  history.back()
}

function refresh () {
  isRefresh.value = true
  noMore.value = false
  videoList.value = []
  searchParams.value.offset = '0'
  scrollToTop()
  getList(false)
}

function onScroll (e) {
  const target = e.target
  const distanceToBottom = target.scrollHeight - target.scrollTop - target.clientHeight
  showBackTop.value = target.scrollTop > 200
  if (!loading.value && !noMore.value && distanceToBottom <= 250) {
    handleLoad()
  }
}

function handleLoad () {
  if (loading.value || noMore.value) return
  loading.value = true
  console.log(searchParams.value)
  getList(true)
}

function scrollToTop () {
  const target = scrollbarRef.value
  if (target) target.scrollTo({ top: 0, behavior: 'smooth' })
}

function getList (isMore = false) {
  loading.value = true
  Search(searchType.value, searchParams.value).then(res => {
    console.log(res)
    searchParams.value.offset = "10"
    let list = res.data.map(a => ({
      aweme_id: a.aweme_info.aweme_id,
      desc: a.aweme_info.desc,
      cover: a.aweme_info.video?.cover?.url_list?.[0] || a.aweme_info.video?.origin_cover?.url_list?.[0],
      author: a.aweme_info.author?.nickname,
      avatar: a.aweme_info.author?.avatar_thumb?.url_list?.[0],
      digg_count: a.aweme_info.statistics?.digg_count || 0,
      play_url: a.aweme_info.video?.play_addr?.url_list?.[0],
      duration: a.aweme_info.video?.duration,
      create_time: a.aweme_info.create_time,
    }))
    applyList(list, isMore)
  }).finally(() => {
    loading.value = false
    isRefresh.value = false
  })
}

function applyList (list, isMore) {
  if (isMore) {
    if (!list || list.length === 0) {
      noMore.value = true
      return
    }
    videoList.value = videoList.value.concat(list)
  } else {
    videoList.value = list
  }
}

function openPlayer (item) {
  currentVideo.value = item
  showPlayer.value = true
}

function formatDuration (msOrSec) {
  if (!msOrSec && msOrSec !== 0) return '-'
  const totalSec = msOrSec > 10000 ? Math.floor(msOrSec / 1000) : Math.floor(msOrSec)
  const m = Math.floor(totalSec / 60)
  const s = totalSec % 60
  const pad = (n) => (n < 10 ? '0' + n : n)
  return `${pad(m)}:${pad(s)}`
}

function formatNum (n) {
  if (n >= 10000) return (n / 10000).toFixed(1) + '万'
  return String(n || 0)
}

function proxy (url) {
  if (!url) return url
  try {
    return `http://127.0.0.1:17892/proxy/dy?u=${encodeURIComponent(url)}`
  } catch (e) {
    return url
  }
}
</script>

<style scoped>
.main-container {
  height: calc(100% - 20px);
  width: calc(100% - 20px);
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  display: flex;
  padding: 10px;
  flex-direction: column;
}

.note_title {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  width: 100%;
  text-align: start;
}
</style>

<template>
    <div class="main-container">
        <div style="display: flex;flex-direction: column;height: 100%;gap: 10px;">
            <n-flex align="center" style="gap: 10px;">
                <n-dropdown key-field="word" label-field="word" trigger="hover" :options="hotWords"
                    @select="handleSelect">
                    <n-input style="max-width: 320px;border-radius: 8px;text-align: start" v-model:value="keyword"
                        clearable placeholder="搜索用户或关键词">
                        <template #suffix>
                            <n-button text @click="doSearch">搜索</n-button>
                        </template>
                    </n-input>
                </n-dropdown>
            </n-flex>
            <n-back-top :bottom="40" :visibility-height="200" :show="showBackTop" @click="scrollToTop" />
            <n-scrollbar ref="scrollbarRef" style="height: 100%; position: relative;" @scroll="onScroll"
                v-if="videoList && videoList.length > 0">
                <n-row :gutter="[12, 12]" style="width: 100%;">
                    <n-col :span="6" v-for="(item, index) in videoList" :key="item.aweme_id || index"
                        style="cursor: pointer;" @click="openPlayer(item)">
                        <n-image width="100%" height="250px;width: 100%;" style="border-radius: 10px"
                            :src="proxy(item.cover)" preview-disabled fit="contain" />
                        <n-ellipsis class="note_title" style="width: 100%;text-align: start" line-clamp="2">
                            {{ item.desc || '无标题' }}
                        </n-ellipsis>
                        <div style="display: flex;gap: 5px;align-items: center">
                            <n-text>@{{ item.author }}</n-text>
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
                    <div v-if="noMore" style="font-weight: bold;margin-top: 10px;width: 100%;">没有更多了 🤪</div>
                </n-flex>
            </n-scrollbar>
            <div v-else style="height: 100%;align-items: center;display: flex;justify-content: center;">
                <n-empty description="暂无数据">
                    <template #extra>
                        <n-button size="small" @click="doSearch">刷新</n-button>
                    </template>
                </n-empty>
            </div>
        </div>
        <n-modal v-model:show="showPlayer" :mask-closable="true" preset="card" :style="{ width: '960px' }"
            :title="currentVideo?.desc || '视频播放'">
            <div style="display: flex; gap: 16px; align-items: flex-start;">
                <div style="flex: 2; background: #000; border-radius: 8px; overflow: hidden;">
                    <video v-if="currentVideo?.play_url" :src="proxy(currentVideo.play_url)"
                        style="width: 100%; height: 520px; object-fit: contain; background: #000;" controls autoplay
                        loop playsinline muted></video>
                    <div v-else
                        style="height: 520px; display: flex; align-items: center; justify-content: center; color: #999;">
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
import { HomeList, HotList } from '../../../wailsjs/go/dy/Dy'
import { formatTime } from '@/utils'
import { useRouter } from 'vue-router'

const router = useRouter()

const keyword = ref('')
const hotWords = ref([])

const homgeListparams = ref({
    count: "20",
    refresh_index: "1",
    refer_id: "",
    refer_type: "10",
    pull_type: "1"
})

const videoList = ref([])
const loading = ref(false)
const isLoadMore = ref(false)
const noMore = ref(false)

const scrollbarRef = ref(null)
const showBackTop = ref(false)
const showPlayer = ref(false)
const currentVideo = ref(null)

function onScroll (e) {
    const target = e.target
    const distanceToBottom = target.scrollHeight - target.scrollTop - target.clientHeight
    showBackTop.value = target.scrollTop > 200
    if (!loading.value && !noMore.value && distanceToBottom <= 250) {
        loadMore()
    }
}

function scrollToTop () {
    const target = scrollbarRef.value
    if (target) target.scrollTo({ top: 0, behavior: 'smooth' })
}

async function fetchHot () {
    const res = await HotList()
    if (res && res.data && res.data.word_list) {
        //取前10个
        hotWords.value = res.data.word_list.slice(0, 10)
    }
}


function getHomeList () {
    loading.value = true
    console.log(homgeListparams.value)
    HomeList(homgeListparams.value).then(res => {
        console.log(res)
        let list = res.aweme_list.map(a => ({
            aweme_id: a.aweme_id,
            desc: a.desc,
            cover: a.video?.cover?.url_list?.[0] || a.video?.origin_cover?.url_list?.[0],
            author: a.author?.nickname,
            avatar: a.author?.avatar_thumb?.url_list?.[0],
            digg_count: a.statistics?.digg_count || 0,
            play_url: a.video?.play_addr?.url_list?.[0],
            duration: a.video?.duration,
            create_time: a.create_time,
        }))
        if (isLoadMore.value) {
            if (videoList.length === 0) {
                noMore.value = true
            }
            videoList.value = videoList.value.concat(list)
        } else {
            videoList.value = list
        }
    }).finally(() => {
        loading.value = false
        isLoadMore.value = false
    })
}


function handleSelect (e) {
    keyword.value = e
    doSearch()
}


function doSearch () {
    router.push({
        name: 'dyDetail',
        query: {
            "keyword": keyword.value,
        }
    })
}

function loadMore () {
    noMore.value = false
    isLoadMore.value = true
    homgeListparams.value.refresh_index = String(Number(homgeListparams.value.refresh_index) + 1)
    getHomeList()
}

onMounted(() => {
    fetchHot()
    getHomeList()
})

function openPlayer (item) {
    currentVideo.value = item
    showPlayer.value = true
}

function formatDuration (msOrSec) {
    if (!msOrSec && msOrSec !== 0) return '-'
    // 抖音 duration 多为毫秒
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
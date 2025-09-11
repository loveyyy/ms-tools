<template>
  <div class="main-container">
    <div style="display: flex; flex-direction: column; height: 100%; gap: 12px;">
      <n-flex align="center" style="gap: 10px;">
        <n-icon size="20px" @click="handleBack">
          <ArrowBack />
        </n-icon>

        <n-input style="max-width: 300px;border-radius: 8px;text-align: start" v-model:value="searchParams.keyword"
          clearable placeholder="请输入内容">
          <template #suffix>
            <n-button text @click="refresh">搜索</n-button>
          </template>
        </n-input>
      </n-flex>

      <n-tabs v-model:value="searchParams.note_type" type="segment" @update:value="refresh">
        <n-tab :tab="item.name" :name="item.value" v-for="(item, index) in typeList" :key="index">
        </n-tab>
      </n-tabs>
      <n-back-top :bottom="40" :visibility-height="200" :show="showBackTop" @click="scrollToTop">
      </n-back-top>
      <n-scrollbar ref="scrollbarRef" style="height: 100%; position: relative;" @scroll="onScroll"
        v-if="searchList && searchList.length > 0">
        <n-spin size="small" stroke="#409eff" v-if="isRefresh" />
        <NoteGrid :list="searchList" @selectNote="openDetail" @openUser="openUser" />
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
    </div>

    <NoteDetailModal v-model:show="showDetail" :detail="detail" />
  </div>

</template>

<script setup>
import { onMounted, ref } from 'vue'
import NoteGrid from './components/NoteGrid.vue'
import NoteDetailModal from './components/NoteDetailModal.vue'
import { useRoute } from 'vue-router'
import { useRouter } from 'vue-router'
import { GetData } from "../../../wailsjs/go/xhs/Xhs";
import { ArrowBack } from '@vicons/ionicons5'
const route = useRoute()
const router = useRouter()

const searchList = ref(null)

const detail = ref(null)
const showDetail = ref(false)

const isRefresh = ref(false)
const loading = ref(false)
const noMore = ref(false)

const scrollbarRef = ref(null)
const showBackTop = ref(false)

const searchParams = ref({
  page: 1,
  page_size: 20,
  keyword: '',
  sort: "general",
  note_type: 0,
  ext_flags: [],
  geo: "",
  image_formats: ["jpg", "webp", "avif"],
})

const typeList = ref([
  { name: '全部', value: '0' },
  { name: '视频', value: '1' },
  { name: '图文', value: '2' }
])

onMounted(() => {
  const query = route.query
  if (!query.keyword) {
    history.back()
    return
  }
  searchParams.value.keyword = query.keyword
  getSearchInfo(false);
})


function refresh () {
  searchParams.value.page = 1
  isRefresh.value = true;
  scrollToTop()
  getSearchInfo(false)
}


function getSearchInfo (isMore = false) {
  GetData(5, searchParams.value).then(res => {
    if (res && res.items) {
      console.log(res)
      let resultList = res.items.filter((item) => item.model_type == "note")
      if (isMore) {
        searchList.value = [...searchList.value, ...resultList]
      } else {
        searchList.value = resultList
      }
    } else {
      noMore.value = true
    }
  }).finally(() => {
    loading.value = false
    isRefresh.value = false
  })
}


function openDetail (item) {
  showDetail.value = true
  detail.value = null
  const params = {
    source_note_id: item.id,
    image_formats: ["jpg", "webp", "avif"],
    extra: { need_body_topic: "1" },
    xsec_source: "pc_feed",
    xsec_token: item.xsec_token
  }
  GetData(4, params).then((res => {
    console.log(res)
    if (res && res.items) {
      detail.value = res.items[0]
    }
  }))
}

function openUser (item) {
  const userId = item.note_card.user.user_id
  const xsec_token = item.note_card.user.xsec_token
  if (!userId || !xsec_token) {
    window.$naive.message.error('用户信息获取失败')
    return
  }
  router.push({
    path: '/xhs/user',
    query: {
      userId,
      xsec_token
    }
  })
}

function handleBack () {
  history.back()
}

function onScroll (e) {
  const target = e.target
  const distanceToBottom = target.scrollHeight - target.scrollTop - target.clientHeight
  if (target.scrollTop <= 200) {
    showBackTop.value = false
  } else {
    showBackTop.value = true
  }
  if (!loading.value && !noMore.value && distanceToBottom <= 200) {
    handleLoad()
  }
}

function handleLoad () {
  if (loading.value || noMore.value) {
    return
  }
  console.log('开始加载1')
  loading.value = true
  searchParams.value.page = searchParams.value.page + 1
  getSearchInfo(true)
}


function scrollToTop () {
  const target = scrollbarRef.value
  if (target) {
    target.scrollTo({ top: 0, behavior: 'smooth' })
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

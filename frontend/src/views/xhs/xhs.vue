<template>
  <div class="main-container">
    <div style="display: flex;flex-direction: column;height: 100%;gap: 10px;">
      <n-flex align="center" style="gap: 10px;">
        <n-input style="max-width: 300px;border-radius: 8px;text-align: start" v-model:value="searchContent" clearable
          placeholder="请输入内容" :disabled="!loginSession">
          <template #suffix>
            <n-button text @click="startSearch" :disabled="!loginSession">搜索</n-button>
          </template>
        </n-input>
        <n-button type="primary" style="border-radius: 8px;" size="small"
          @click="showSetSessionModal = true">设置登录会话开始搜索</n-button>
      </n-flex>
      <n-tabs v-model:value="params.category" type="segment" @update:value="categoryChange">
        <n-tab :tab="item.name" :name="item.id" v-for="(item, index) in categoryList" :key="index">
        </n-tab>
      </n-tabs>
      <n-back-top :bottom="40" :visibility-height="200" :show="showBackTop" @click="scrollToTop">
      </n-back-top>
      <n-scrollbar ref="scrollbarRef" style="height: 100%; position: relative;" @scroll="onScroll"
        v-if="xhsList && xhsList.length > 0">
        <n-spin size="small" stroke="#409eff" v-if="isRefresh" />
        <NoteGrid :list="xhsList" @selectNote="openDetail" @openUser="openUser" />
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
            <n-button size="small" @click="refresh">
              刷新
            </n-button>
          </template>
        </n-empty>
      </div>
    </div>
    <n-modal v-model:show="showSetSessionModal" preset="card" title="设置登录会话" style="width: 500px;" size="small">
      <n-flex style="gap: 15px;padding: 20px;">
        <n-text style="font-weight: 400;font-size: 14px;text-align: start;">登录会话</n-text>
        <n-input style="border-radius: 6px" v-model:value="loginSession" placeholder="请输入登录会话" clearable>
        </n-input>
        <n-button block type="primary" round :disabled="!loginSession" @click="confirmSetLoginSession">
          保存
        </n-button>
      </n-flex>
    </n-modal>
    <NoteDetailModal v-model:show="showDetail" :detail="detail" />
  </div>
</template>


<script setup>
import { onMounted, ref } from 'vue'
import NoteGrid from './components/NoteGrid.vue'
import NoteDetailModal from './components/NoteDetailModal.vue'
import { GetData, IsLoginSession, SetLoginSession } from "../../../wailsjs/go/xhs/Xhs";
import { useRouter } from 'vue-router';


const router = useRouter();

const searchContent = ref('')

const isRefresh = ref(false)
const loading = ref(false)
const noMore = ref(false)

const showBackTop = ref(false)
const scrollbarRef = ref(null)

const categoryList = ref([])

const params = ref({
  cursor_score: '',
  num: 30,
  refresh_type: 1,
  note_index: 0,
  unread_begin_note_id: '',
  unread_end_note_id: '',
  unread_note_count: 0,
  category: 'homefeed_recommend',
  search_key: '',
  need_num: 15,
  image_formats: ['jpg', 'webp', 'avif'],
  need_filter_image: false,
})

const xhsList = ref(null)
const showSetSessionModal = ref(false)
const loginSession = ref('')

const showDetail = ref(false)
const detail = ref(null)


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

function scrollToTop () {
  const target = scrollbarRef.value
  if (target) {
    target.scrollTo({ top: 0, behavior: 'smooth' })
  }
}


function startSearch () {
  router.push({
    name: 'xhsDetail',
    query: {
      "keyword": searchContent.value,
    }
  })
}


function getIsLoginSession () {
  IsLoginSession().then((res => {
    console.log(res)
    loginSession.value = res
  }))
}


function confirmSetLoginSession () {
  SetLoginSession(loginSession.value).then(() => {
    showSetSessionModal.value = false
    getIsLoginSession()
  })
}


function handleLoad () {
  console.log('开始加载')
  if (loading.value || noMore.value) {
    return
  }
  console.log('开始加载1')
  loading.value = true
  getList(true)
}


function getList (isMore = false) {
  params.value.refresh_type = isMore ? 3 : 1
  GetData(2, params.value).then((res => {
    if (res) {
      if (res.items) {
        params.value.cursor_score = res.cursor_score
        if (isMore) {
          xhsList.value = [...xhsList.value, ...res.items]
        } else {
          xhsList.value = res.items
        }
        params.value.note_index = xhsList.value.length + 1
      } else {
        noMore.value = true
      }
    }
  }))
    .finally(() => {
      loading.value = false
      isRefresh.value = false
    })
}

function refresh () {
  params.value.cursor_score = "";
  params.value.refresh_type = 1
  isRefresh.value = true
  getList()
}


function categoryChange (val) {
  params.value.category = val
  params.value.cursor_score = "";
  params.value.refresh_type = 1
  isRefresh.value = true
  scrollToTop()
  getList()
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
    if (res && res.items) {
      detail.value = res.items[0]
    }
  }))
}

function openUser (item) {
  const user_id = item.note_card.user.user_id
  const xsec_token = item.note_card.user.xsec_token
  if (!user_id || !xsec_token) {
    window.$naive.$message.error('用户信息获取失败')
    return
  }
  router.push({
    path: 'xhsUser',
    query: {
      user_id,
      xsec_token
    }
  })
}

// 详情展示已抽到 NoteDetailModal

onMounted(() => {
  GetData(1, {}).then((res => {
    let categorys = [
      { "id": "homefeed_recommend", "name": "推荐" }
    ]
    if (res && res.categories) {
      categorys = categorys.concat(res.categories)
    }
    categoryList.value = categorys
  }))
  isRefresh.value = true
  getList();
  getIsLoginSession()
})

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
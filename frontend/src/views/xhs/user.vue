<template>
  <div class="main-container">
    <div style="display: flex; flex-direction: column; height: 100%; gap: 12px;">
      <n-flex align="center" style="gap: 10px;">
        <n-button text @click="handleBack">
          <n-icon size="18">
            <ArrowBack />
          </n-icon>
          返回
        </n-button>
      </n-flex>

      <n-card size="small">
        <n-flex align="center" style="gap: 20px;">
          <div>
            <n-image height="120px" style="border-radius: 50%;width: 120px;" object-fit="cover"
              :src="userInfo.avatar_url" />
          </div>
          <n-flex vertical align="start" style="flex: 1;">
            <n-text style="font-weight: bold;font-size: 16px;">{{ userInfo.nick_name }}</n-text>
            <n-text depth="3" style="font-size: 12px;">{{ userInfo.id }} | {{ userInfo.ip_address }}</n-text>
            <n-ellipsis style="max-width: 600px; text-align: start;">{{ userInfo.desc }}</n-ellipsis>
            <n-space>
              <n-tag size="small" v-for="(t, i) in userInfo.tags" :key="i" round>{{ t }}</n-tag>
            </n-space>
            <n-space style="margin-left: auto;">
              <n-tag round size="small">关注 {{ userInfo.follow_count }}</n-tag>
              <n-tag round size="small">粉丝 {{ userInfo.fans_count }}</n-tag>
              <n-tag round size="small">获赞与收藏 {{ userInfo.like_count }}</n-tag>
            </n-space>
          </n-flex>
        </n-flex>
      </n-card>
      <n-flex style="padding: 0 20px;">
        <n-button style="margin-left: auto;" size="small" type="primary" @click="handleExport">导出当前</n-button>
      </n-flex>
      <n-back-top :bottom="40" :visibility-height="200" :show="showBackTop" @click="scrollToTop" />
      <n-scrollbar ref="scrollbarRef" style="height: 100%; position: relative;" @scroll="onScroll">
        <n-spin size="small" stroke="#409eff" v-if="isRefresh" />
        <NoteGrid :list="noteList" @selectNote="openDetail" />
        <n-flex wrap align="center">
          <div v-if="loading"
            style="font-weight: bold;margin-top: 10px;width: 100%;display: flex;align-items: center;justify-content: center;gap: 10px;">
            <n-spin size="small" stroke="#409eff" />
            <span> 加载中...</span>
          </div>
          <div v-if="noMore" style="font-weight: bold;margin-top: 10px;width: 100%;">没有更多了 🤪</div>
        </n-flex>
      </n-scrollbar>

      <NoteDetailModal :show="showDetail" :detail="detail" />
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NoteGrid from './components/NoteGrid.vue'
import NoteDetailModal from './components/NoteDetailModal.vue'
import { GetData, IsLoginSession } from "../../../wailsjs/go/xhs/Xhs"
import { ArrowBack } from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()

const userInfo = ref({
  avatar_url: '',
  nick_name: '',
  id: '',
  ip_address: '',
  desc: '',
  tags: [],
  follow_count: '',
  fans_count: '',
  like_count: ''
})

const isSetLoginSession = ref(false)

const noteList = ref([])

const showDetail = ref(false)
const detail = ref(null)

const isRefresh = ref(false)
const loading = ref(false)
const noMore = ref(false)

const scrollbarRef = ref(null)
const showBackTop = ref(false)

const params = ref({
  user_id: '',
  channel_type: "explore_feed",
  parent_page_channel_type: 'web_profile_board',
  xsec_token: '',
  xsec_source: "pc_feed",
  num: '',
  cursor: '',
})

onMounted(() => {
  const q = route.query
  if (!q.user_id && !q.xsec_token) {
    history.back()
    return
  }

  IsLoginSession().then(res => {
    isSetLoginSession.value = res
  })

  params.value.user_id = q.user_id
  params.value.xsec_token = q.xsec_token
  loadUser()
  isRefresh.value = true
})

function handleBack () {
  history.back()
}

async function handleExport () {
  if (!isSetLoginSession.value) {
    window.$navite.$message.info('需要设置登录session才能获取数据')
    return
  }

  if (!noteList.value || noteList.value.length === 0) {
    window.$naive.$message.warning('暂无可导出数据')
    return
  }

  const rows = noteList.value.map((item) => {
    const card = item.note_card
    const user =card.user
    const type = card.type || ''
    let imageUrls = ''
    if (Array.isArray(card.image_list) && card.image_list.length > 0) {
      imageUrls = card.image_list.map(i => i.url_default || i.url).filter(Boolean).join(', ')
    } else if (card.cover) {
      imageUrls = card.cover.url_default || card.cover.url || ''
    }
    const videoUrl = getVideoUrl(item)
    const desc = filterDesc(card.desc || '')
    const liked = formatNum(card.interact_info.liked_count)
    const collected = formatNum(card.interact_info.collected_count)
    return {
      '标题': card.display_title || '',
      '描述': desc,
      '类型': type,
      '图片地址': imageUrls,
      '视频地址': videoUrl,
      '点赞数量': liked,
      '收藏数量': collected,
      '作者昵称': user.nick_name ||'',
      '作者ID': user.user_id || '',
      '作者头像': user.avatar || ''
    }
  })

  const fileName = `${userInfo.value.nick_name || 'xhs'}_${new Date().toISOString().slice(0, 10)}.xlsx`

  try {
    const XLSX = await import('xlsx')
    const ws = XLSX.utils.json_to_sheet(rows)
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, '笔记')
    XLSX.writeFile(wb, fileName)
    return
  } catch (e) {
    // fallback CSV
    try {
      const headers = Object.keys(rows[0])
      const csv = [headers.join(',')].concat(
        rows.map(r => headers.map(h => {
          const val = (r[h] ?? '').toString().replace(/"/g, '""')
          return `"${val}"`
        }).join(','))
      ).join('\n')
      const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
      const link = document.createElement('a')
      const url = URL.createObjectURL(blob)
      link.href = url
      link.download = fileName.replace(/\.xlsx$/, '.csv')
      link.click()
      URL.revokeObjectURL(url)
    } catch (_) { }
  }
}

function loadUser () {
  GetData(6, params.value).then(res => {
    console.log(res)
    if (res) {
      userInfo.value = res
      var newList = []
      res.feed_list.forEach(element => {
        newList.push(objectToUnderline(element))
      });
      console.log(newList)
      noteList.value = newList
      if (res.hasMore) {
        params.value.cursor = res.cursor
        params.value.num = String(res.num)
      }
      noMore.value = !res.hasMore
    }
  }).finally(() => {
    loading.value = false
    isRefresh.value = false
  })
}


//将数组对象中的key由驼峰转化为下划线 如果是对象或者数组 递归处理
function objectToUnderline (input) {
  if (Array.isArray(input)) {
    return input.map(item => objectToUnderline(item))
  }
  if (input && typeof input === 'object') {
    const result = {}
    Object.keys(input).forEach((key) => {
      const snakeKey = key
        .replace(/([A-Z])/g, '_$1')
        .replace(/-/g, '_')
        .toLowerCase()
      result[snakeKey] = objectToUnderline(input[key])
    })
    return result
  }
  return input
}

function getList () {
  if (!isSetLoginSession.value) {
    window.$navite.$message.info('需要设置登录session才能获取数据')
    return
  }
  if (noMore.value) {
    return
  }
  GetData(7, params.value).then(res => {
    console.log(res)
    if (res && res.notes) {
      noMore.value = !res.has_more
      console.log(noMore.value)
      params.value.cursor = res.cursor
      var newList = res.notes.map((item) => {
        return {
          id: item.note_id,
          note_card: item,
          ssr_rendered: item.ssr_rendered,
          xsec_token: item.xsec_token
        }
      })
      console.log(newList)
      noteList.value = [...noteList.value, ...newList]
    }
  }).finally(() => {
    loading.value = false
    isRefresh.value = false
  })
}

function onScroll (e) {
  const target = e.target
  const distanceToBottom = target.scrollHeight - target.scrollTop - target.clientHeight
  showBackTop.value = target.scrollTop > 200
  if (!loading.value && !noMore.value && distanceToBottom <= 200) {
    handleLoad()
  }
}

function handleLoad () {
  if (loading.value || noMore.value) return
  loading.value = true
  getList(true)
}

function scrollToTop () {
  const target = scrollbarRef.value
  if (target) target.scrollTo({ top: 0, behavior: 'smooth' })
}

function openDetail (item) {
  if (!isSetLoginSession.value) {
    window.$navite.$message.info('需要设置登录session才能获取数据')
    return
  }
  showDetail.value = true
  detail.value = null
  const p = {
    source_note_id: item.id,
    image_formats: ["jpg", "webp", "avif"],
    extra: { need_body_topic: "1" },
    xsec_source: "pc_feed",
    xsec_token: item.xsec_token
  }
  GetData(4, p).then((res => {
    if (res && res.items) {
      detail.value = res.items[0]
    }
  }))
}

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

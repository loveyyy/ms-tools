<!--
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-03-22 21:30:22
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-11 16:59:09
 * @Description:
-->
<template>
  <div class="ai-chat-container">
    <n-layout position="absolute">
      <n-layout position="absolute" style="bottom: 64px;background-color: white; padding: 10px; border-radius: 7px"
      >
        <BubbleList
            :roles="roles"
            :items="messages"/>
                <n-flex vertical>
                  <template v-for="(item, index) in messages" :key="index">
                    <Bubble
                        v-if="item.role === 'user'"
                        placement="end"
                        :content="item.content">
                      <template #footer>
                        <n-flex >
                          <n-icon>
                            <RefreshSharp/>
                          </n-icon>
                          <n-icon>
                            <CopyOutline/>
                          </n-icon>
                        </n-flex>
                      </template>
                    </Bubble>
                    <Bubble
                        v-else
                        placement="start"
                        :content="item.content"
                        :loading="isLoading"
                        :messageRender="renderMarkdown"
                        :typing="{ step: 2, interval: 50 }"
                        variant="filled">
                      <template #footer>
                        <n-flex >
                          <n-icon>
                            <CopyOutline/>
                          </n-icon>
                        </n-flex>
                      </template>
                    </Bubble>
                  </template>
                </n-flex>
      </n-layout>
      <n-layout-footer position="absolute"
                       style="background-color: transparent">
        <Sender v-model:value="userInput" placeholder="输入消息......" class="version_select"
                @submit="sendMessage">
          <!--          <template #actions>-->
          <!--            <n-button-->
          <!--                class="send-btn"-->
          <!--                style="background-color: #21c16bff"-->
          <!--                @click="sendMessage()"-->
          <!--                :loading="false"-->
          <!--            >-->
          <!--              <template #icon>-->
          <!--                <n-icon size="15" color="#ffffff">-->
          <!--                  <SendOutline/>-->
          <!--                </n-icon>-->
          <!--              </template>-->
          <!--            </n-button>-->
          <!--          </template>-->
          <!--          <template #prefix v-if="propsData.canUpload">-->
          <!--            <n-upload @on-before-upload="handleChange()">-->
          <!--              <n-button class="send-btn" type="primary" text @click="triggerOpen">-->
          <!--                <template #icon>-->
          <!--                  <n-icon size="25" color="#000">-->
          <!--                    <AttachOutline/>-->
          <!--                  </n-icon>-->
          <!--                </template>-->
          <!--              </n-button>-->
          <!--            </n-upload>-->
          <!--          </template>-->
        </Sender>
      </n-layout-footer>
    </n-layout>
  </div>
</template>

<script setup>
import {Bubble, Sender} from "ant-design-x-vue";
import {h, onMounted, ref} from "vue";
import OpenAI from 'openai'
import markdownit from 'markdown-it';
import {NIcon, NText} from 'naive-ui';
import {CopyOutline, RefreshSharp} from '@vicons/ionicons5'

const md = markdownit({html: true, breaks: true});

const openai = ref('')

const abortController = ref('')

let propsData = defineProps({
  version: {
    type: String,
    default: "pro"
  },
  canUpload: {
    type: Boolean,
    default: true
  }
});

const userInput = ref('');

const messages = ref([])

const isLoading = ref(false)

onMounted(() => {
  initOpenAI()
});

function initOpenAI() {
  const config = {
    apiKey: '',
    baseURL: '',
    dangerouslyAllowBrowser: true
  }
  openai.value = new OpenAI(config)
}

async function callAIAPI(message, useWebSearch = false) {
  if (!openai.value) {
    throw new Error('请先配置 OpenAI API Key')
  }

  abortController.value = new AbortController()

  let sendMessage = []
  //加载上下文 //通过设置限制上下文条数
  var contentMsg = messages.value.slice(messages.value.length - 10, messages.value.length)
  contentMsg.forEach((msg) => {
    if (msg) {
      sendMessage.push({
        "role": msg.role,
        "content": msg.content,
      })
    }
  })

  // 添加一个标志来检查是否被中断
  let isAborted = false
  abortController.value.signal.addEventListener('abort', () => {
    isAborted = true
  })

  try {
    if (useWebSearch) {

    }

    const response = await openai.value.chat.completions.create({
      model: "deepseek-r1",
      messages: sendMessage,
      stream: true,
      temperature: 1,
      top_p: 1,
      frequency_penalty: 0,
      presence_penalty: 0,
      n: 1,
      signal: abortController.value.signal,
    })

    // 创建消息占位
    const responseMessage = {
      role: 'assistant',
      content: '',
      isStreaming: true
    }
    messages.value.push(responseMessage)

    const messageIndex = messages.value.length - 1

    var totalContent = ''
    try {
      // 流式读取数据
      for await (const chunk of response) {
        // 检查是否被中断
        if (isAborted) {
          messages.value[messageIndex].content += '\n[已停止生成]'
          break
        }
        if (isLoading.value) {
          isLoading.value = false
        }
        const content = chunk.choices[0]?.delta?.content || ""
        if (content) {
          totalContent += content
          messages.value[messageIndex].content = totalContent
        }
      }
    } catch (error) {
      if (error.name === 'AbortError') {
        messages.value[messageIndex].content += '\n[已停止生成]'
        return
      }
      throw error
    }

    // 流式传输完成
    delete messages.value[messageIndex].isStreaming
    return messages.value[messageIndex].content
  } catch (error) {
    if (error.name === 'AbortError') {
      console.log('请求已被中断')
      return
    }
    if (error.response?.status === 401) {
      throw new Error('API Key 无效，请检查设置')
    }
    throw error
  } finally {
    abortController.value = null
  }
}


function handleChange(info) {
}

async function sendMessage() {
  if (!userInput.value.trim()) return

  if (!openai.value) {
    window.$message.error('请先在设置中配置 OpenAI API Key')
    return
  }

  let messageContent = userInput.value.trim()
  try {
    // // 处理附件内容
    // if (this.selectedFile) {
    //   const fileContent = await this.readFileContent(this.selectedFile)
    //   let fileMessage = `[附件: ${this.selectedFile.name}]\n`
    //
    //   if (this.selectedFile.type.includes('image')) {
    //     // 如果是图片，添加图片描述
    //     fileMessage += `[图片内容: ${fileContent}]\n`
    //   } else {
    //     // 如果是文本文件，添加文件内容
    //     fileMessage += `文件内容:\n${fileContent}\n`
    //   }
    //
    //   messageContent = fileMessage + messageContent
    // }

    // 处理粘贴内容
    // if (this.pastedContent) {
    //   messageContent = this.pastedContent.content + '\n' + messageContent
    // }

    // 添加用户消息
    messages.value.push({
      role: 'user',
      content: messageContent
    })

    userInput.value = ''
    isLoading.value = true

    // 调用 AI API
    await callAIAPI(messageContent)
  } catch (error) {
    console.error('发送消息失败:', error)
    let errorMessage = '发送消息失败'
    if (error.response) {
      switch (error.response.status) {
        case 401:
          errorMessage = 'API Key 无效，请检查设置'
          break
        case 429:
          errorMessage = '请求过于频繁，请稍后再试'
          break
        case 500:
          errorMessage = 'OpenAI 服务器错误，请稍后重试'
          break
        default:
          errorMessage = `请求失败: ${error.response.status}`
      }
    } else if (error.message) {
      errorMessage = error.message
    }
    window.$message.error(errorMessage)
    messages.value.push({
      role: 'assistant',
      content: errorMessage
    })
  } finally {
    isLoading.value = false
  }
}

function renderMarkdown(content) {
  return h(NText, null, {
    default: () => h('div', {innerHTML: md.render(content)}),
  })
}

function copyMessage(index) {
  console.log(index)
}

</script>

<style scoped>
.ai-chat-container {
  display: flex;
  flex-direction: column;
  position: relative;
  height: calc(100vh - 84px);
  background-color: white;
}

.version_select {
  background: white;
  border-radius: 15px 15px 15px 15px;
  font-weight: normal;
  font-size: 13px;
  color: #c0c3cf;
}

:deep(.ant-input) {
  font-weight: normal;
  font-size: 15px;
  color: #000000;
  letter-spacing: 1.2px;
}

:deep(.n-layout) {
  --n-color: transparent !important;
}

:deep(.n-layout-footer) {
  --n-color: transparent !important;
}

.send-btn {
  width: 25px;
  height: 25px;
  border-radius: 7px;
}
</style>

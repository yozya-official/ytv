<template>
  <!-- 加载中状态 -->
  <div v-if="isLoading" class="w-full max-w-7xl mx-auto py-2 px-2 sm:px-4">
    <div
      class="relative bg-gradient-to-br from-black/80 to-black/60 backdrop-blur-sm rounded-xl overflow-hidden max-h-full w-full aspect-video"
    >
      <div class="absolute inset-0 flex flex-col gap-6 items-center justify-center">
        <!-- 加载动画 -->
        <div class="relative">
          <div class="w-20 h-20 border-4 border-primary/30 rounded-full"></div>
          <div
            class="absolute top-0 left-0 w-20 h-20 border-4 border-primary border-t-transparent rounded-full animate-spin"
          ></div>
        </div>

        <!-- 加载文字 -->
        <div class="text-center space-y-2">
          <p class="text-xl font-semibold text-white">正在加载视频</p>
          <p class="text-sm text-white/60">请稍候...</p>
        </div>

        <!-- 加载进度指示器 -->
        <div class="w-64 h-1 bg-white/10 rounded-full overflow-hidden">
          <div
            class="h-full bg-gradient-to-r from-primary to-primary/60 rounded-full animate-pulse"
            style="width: 60%"
          ></div>
        </div>
      </div>

      <!-- 装饰性背景元素 -->
      <div class="absolute top-10 left-10 w-32 h-32 bg-primary/5 rounded-full blur-3xl"></div>
      <div class="absolute bottom-10 right-10 w-40 h-40 bg-accent/5 rounded-full blur-3xl"></div>
    </div>
  </div>

  <!-- 加载失败状态 -->
  <div v-else-if="failed" class="w-full max-w-7xl mx-auto py-2 px-2 sm:px-4">
    <div
      class="relative bg-gradient-to-br from-destructive/10 via-black/80 to-black/60 backdrop-blur-sm rounded-xl overflow-hidden border border-destructive/20 h-[40rem] max-h-full w-full aspect-video"
    >
      <div class="absolute inset-0 flex flex-col gap-6 items-center justify-center p-8 z-10">
        <!-- 错误图标 -->
        <div class="relative">
          <div class="w-24 h-24 rounded-full bg-destructive/10 flex items-center justify-center">
            <svg
              class="w-12 h-12 text-destructive"
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
            >
              <g
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
              >
                <circle cx="12" cy="12" r="10" />
                <path d="M12 8v4m0 4h.01" />
              </g>
            </svg>
          </div>
          <!-- 脉冲动画 -->
          <div class="absolute inset-0 rounded-full bg-destructive/20 animate-ping"></div>
        </div>

        <!-- 错误信息 -->
        <div class="text-center space-y-3 max-w-md">
          <h3 class="text-2xl font-bold text-white">视频加载失败</h3>
          <p class="text-white/70 leading-relaxed">
            无法加载视频内容，请检查链接是否正确或稍后重试
          </p>
          <p class="text-sm text-white/50">如问题持续存在，请联系站长获取帮助</p>
        </div>

        <!-- 操作按钮 -->
        <div class="flex items-center gap-3 mt-4">
          <button @click="handleRetry" class="btn btn-outline gap-2 group">
            <svg
              class="size-4"
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
            >
              <g
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
              >
                <path d="M3 12a9 9 0 0 1 9-9a9.75 9.75 0 0 1 6.74 2.74L21 8" />
                <path d="M21 3v5h-5m5 4a9 9 0 0 1-9 9a9.75 9.75 0 0 1-6.74-2.74L3 16" />
                <path d="M8 16H3v5" />
              </g>
            </svg>
            重新加载
          </button>

          <button @click="router.push({ name: 'home' })" class="btn btn-primary gap-2">
            <svg
              class="size-4"
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
            >
              <g
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
              >
                <path d="M15 21v-8a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v8" />
                <path
                  d="M3 10a2 2 0 0 1 .709-1.528l7-6a2 2 0 0 1 2.582 0l7 6A2 2 0 0 1 21 10v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"
                />
              </g>
            </svg>
            返回主页
          </button>
        </div>

        <!-- 帮助提示 -->
        <div class="mt-8 p-4 bg-white/5 rounded-lg border border-white/10 max-w-md">
          <p class="text-sm text-white/60 text-center">
            💡 提示：尝试切换不同的视频源或检查网络连接
          </p>
        </div>
      </div>

      <!-- 装饰性背景元素 -->
      <div class="absolute top-0 left-0 w-full h-full opacity-5">
        <div class="absolute top-20 left-20 w-40 h-40 bg-destructive rounded-full blur-3xl"></div>
        <div
          class="absolute bottom-20 right-20 w-48 h-48 bg-destructive rounded-full blur-3xl"
        ></div>
      </div>
    </div>
  </div>

  <div v-else-if="selectedVod" class="w-full max-w-7xl mx-auto py-2 px-2 sm:px-4">
    <div class="relative bg-black rounded-lg overflow-hidden">
      <VideoPlayer
        :player_url="playerUrl"
        :episode_index="episodeIndex"
        :episode_title="''"
        :source_key="selectedVod?.source_key || ''"
        :vod_id="vodId"
        :vod_name="selectedVod.vod_name"
        @play-end="goToNextEpisode"
      ></VideoPlayer>
    </div>
    <div class="w-full max-w-7xl mx-auto pb-2 my-4">
      <div class="p-4 flex flex-col gap-4 card bg-card">
        <div class="flex gap-3">
          <h3 class="font-bold">{{ selectedVod.vod_name }}</h3>
          <button class="btn rounded-xl border-primary btn-outline btn-xs" @click="handleSearch">
            切换资源
          </button>
        </div>

        <div class="flex justify-between">
          <button
            class="btn btn-secondary"
            @click="goToPreviousEpisode"
            :disabled="isDisabledPrevious"
          >
            上一集
          </button>

          <div v-if="selectedVod">
            第{{ episodeIndex + 1 }} / {{ selectedVod.episodes.length }} 集
          </div>
          <div v-else>加载中...</div>

          <button class="btn btn-secondary" @click="goToNextEpisode" :disabled="isDisabledNext">
            下一集
          </button>
        </div>
        <div class="border"></div>
        <EpisodeList
          v-if="selectedVod"
          v-model:select-vod="selectedVod"
          :episode-index="episodeIndex"
        ></EpisodeList>
      </div>
    </div>
  </div>

  <dialog class="dialog" ref="switchRef">
    <div class="dialog-body bg-background w-[64rem]">
      <!-- 关闭按钮 -->
      <button @click="switchRef?.close()" class="btn-close">
        <svg class="size-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>

      <SwitchResult
        @send-close="close"
        :search-results="searchResults"
        :sourceKey="sourceKey"
        :vodId="vodId"
      ></SwitchResult>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import router from '@/router'
import type { VodItem } from '@/models'
import { toast } from '@yuelioi/toast'
const switchRef = ref<HTMLDialogElement>()

const selectedVod = ref<VodItem>()

const route = useRoute()

const vodId = computed(() => Number(route.query.vodId))
const episodeIndex = computed(() => Number(route.query.episodeIndex ?? 0))
const sourceKey = computed(() => route.query.sourceKey as string)

const searchResults = ref<VodItem[]>([])

const playerUrl = ref('')

const isLoading = ref(false)
const failed = ref(false)

const handleRetry = () => {
  location.reload()
}

const handleSearch = async () => {
  try {
    toast.info('切换资源中...请稍后', { duration: 2000 })
    const resp = await videoApi.searchAll(selectedVod.value?.vod_name || '')
    switchRef.value?.showModal()
    searchResults.value = resp.data.data.list
  } catch (err: unknown) {
    handleApiError(err)
  }
}

// 监听 id 以及 资源id, 用于判断是否更换视频
watch(
  [vodId, sourceKey, episodeIndex],
  async ([newId, newSourceKey, newEpisodeIndex]) => {
    if (!newId || !newSourceKey) {
      failed.value = true
    }

    try {
      isLoading.value = true
      const result = await videoApi.searchById(newSourceKey, newId, newEpisodeIndex)
      if (result.data.data) {
        selectedVod.value = result.data.data

        playerUrl.value = selectedVod.value.episodes[episodeIndex.value].url
      } else {
        selectedVod.value = undefined
      }
    } catch (err) {
      console.error('查询视频失败', err)
      toast.error('查询视频失败')
      selectedVod.value = undefined
      failed.value = true
    } finally {
      isLoading.value = false
    }
  },
  { immediate: true },
)

watch(episodeIndex, () => {
  const episode = selectedVod.value?.episodes.find((e) => e.episode_index === episodeIndex.value)
  if (!episode?.url) return
  playerUrl.value = episode?.url
})

const close = () => {
  switchRef.value?.close()
}

// 判断上一集按钮是否禁用
const isDisabledPrevious = computed(() => {
  return episodeIndex.value <= 0
})

// 判断下一集按钮是否禁用
const isDisabledNext = computed(() => {
  if (!selectedVod.value) return
  return episodeIndex.value >= selectedVod.value.episodes.length - 1
})

// 切换到上一集
const goToPreviousEpisode = () => {
  if (isDisabledPrevious.value || !selectedVod.value) return

  router.push({
    name: route.name,
    query: {
      ...route.query,
      episodeIndex: episodeIndex.value - 1,
    },
  })
}

// 切换到下一集
const goToNextEpisode = () => {
  if (isDisabledNext.value || !selectedVod.value) return

  router.push({
    name: route.name,
    query: {
      ...route.query,
      episodeIndex: episodeIndex.value + 1,
    },
  })
}
</script>

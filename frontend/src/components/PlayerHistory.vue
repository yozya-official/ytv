<template>
  <teleport to="body">
    <div class="p-6">
      <!-- 对话框 -->
      <dialog :open="isOpen" class="dialog">
        <div class="dialog-body max-w-2xl w-full">
          <!-- 头部 -->
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center gap-3">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="w-6 h-6 text-primary"
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
                  <path d="M12 6v6l4 2" />
                  <circle cx="12" cy="12" r="10" />
                </g>
              </svg>

              <h2 class="text-2xl font-bold">播放历史</h2>
              <span v-if="playHistory.length > 0" class="badge badge-muted">
                共 {{ playHistory.length }} 条
              </span>
            </div>
            <div class="flex items-center gap-2">
              <button
                v-if="playHistory.length > 0"
                @click="clearAllHistory"
                class="btn btn-ghost btn-sm text-destructive hover:bg-destructive/10"
              >
                清空全部
              </button>
              <button @click="closeDialog" class="btn btn-ghost btn-icon btn-sm">
                <svg
                  class="w-5 h-5"
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                >
                  <path
                    fill="none"
                    stroke="currentColor"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M18 6L6 18M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          </div>

          <div class="pr-2">
            <!-- 搜索框 -->
            <div class="search-input mb-3 search-input-primary max-w-2xl w-full shadow-md relative">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="size-5 text-muted-foreground flex-shrink-0"
                viewBox="0 0 24 24"
              >
                <g
                  fill="none"
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                >
                  <circle cx="11" cy="11" r="8" />
                  <path d="m21 21l-4.34-4.34" />
                </g>
              </svg>

              <input
                type="text"
                v-model.lazy="search"
                @focus="isOpen = true"
                class="flex-1 bg-transparent border-none outline-none text-sm placeholder:text-muted-foreground"
                placeholder="搜索名称, 剧集..."
              />
              <svg
                @click="search = ''"
                xmlns="http://www.w3.org/2000/svg"
                class="size-5"
                width="24"
                height="24"
                viewBox="0 0 24 24"
              >
                <path
                  fill="none"
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M18 6L6 18M6 6l12 12"
                />
              </svg>
              <button class="btn btn-primary btn-xs">
                <span>搜索</span>
              </button>
            </div>

            <!-- 历史记录列表 -->
            <div class="space-y-3 max-h-96 overflow-y-auto">
              <div v-if="playHistory.length === 0" class="text-center py-12">
                <svg
                  class="w-16 h-16 mx-auto text-muted-foreground/30 mb-4"
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
                    <path d="M12 6v6l4 2" />
                    <circle cx="12" cy="12" r="10" />
                  </g>
                </svg>
                <p class="text-muted-foreground text-lg">暂无播放记录</p>
              </div>

              <div
                v-for="item in filterHistory"
                :key="`${item.vod_id}-${item.episode_index}-${item.sourceKey}`"
                class="card group hover:shadow-lg transition-all"
              >
                <div class="p-4">
                  <div class="flex items-start gap-4">
                    <!-- 播放图标 -->
                    <div
                      @click="goToPlayer(item)"
                      class="flex-shrink-0 w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors"
                    >
                      <svg
                        class="w-6 h-6 text-primary"
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                      >
                        <path
                          fill="none"
                          stroke="currentColor"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M5 5a2 2 0 0 1 3.008-1.728l11.997 6.998a2 2 0 0 1 .003 3.458l-12 7A2 2 0 0 1 5 19z"
                        />
                      </svg>
                    </div>

                    <!-- 内容区域 -->
                    <div class="flex-1 min-w-0">
                      <!-- 标题和集数 -->
                      <div class="flex items-start justify-between gap-2 mb-2">
                        <div class="flex-1 min-w-0">
                          <h3 class="font-semibold text-foreground truncate">
                            {{ item.name }}
                          </h3>
                          <div class="flex items-center gap-2 mt-1">
                            <span class="badge badge-outline badge-sm">
                              第 {{ item.episode_index + 1 }} 集
                            </span>
                            <span class="text-xs text-muted-foreground">
                              {{ formatLastPlayTime(item.lastPlayTime) }}
                            </span>
                          </div>
                        </div>

                        <!-- 删除按钮 -->
                        <button
                          @click="deleteHistory(item.vod_id, item.episode_index, item.sourceKey)"
                          class="btn btn-ghost btn-icon btn-xs text-muted-foreground hover:text-destructive opacity-0 group-hover:opacity-100 transition-opacity"
                          title="删除记录"
                        >
                          <svg
                            class="w-4 h-4"
                            xmlns="http://www.w3.org/2000/svg"
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                          >
                            <path
                              fill="none"
                              stroke="currentColor"
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M10 11v6m4-6v6m5-11v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                            />
                          </svg>
                        </button>
                      </div>

                      <!-- 进度条 -->
                      <div class="space-y-1">
                        <div
                          class="flex items-center justify-between text-xs text-muted-foreground"
                        >
                          <span
                            >已观看
                            {{
                              Math.round(getProgressPercent(item.progress, item.duration))
                            }}%</span
                          >
                          <span
                            >{{ formatTime(item.progress) }} / {{ formatTime(item.duration) }}</span
                          >
                        </div>
                        <div class="h-1.5 bg-muted rounded-full overflow-hidden">
                          <div
                            class="h-full bg-primary rounded-full transition-all"
                            :style="{
                              width: `${getProgressPercent(item.progress, item.duration)}%`,
                            }"
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </dialog></div
  ></teleport>
</template>

<script setup lang="ts">
import type { PlayHistory } from '@/models/history'
import router from '@/router'

// 模拟数据
const { playHistory } = storeToRefs(useHistoryStore())

const search = ref('')

const isOpen = defineModel<boolean>('showHistory')

const goToPlayer = (item: PlayHistory) => {
  router.push({
    name: 'player',
    query: {
      vodId: item.vod_id,
      sourceKey: item.sourceKey,
      episodeIndex: item.episode_index,
    },
  })
  isOpen.value = false
}

// 格式化时间
const formatTime = (seconds: number): string => {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)

  if (h > 0) {
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  }
  return `${m}:${s.toString().padStart(2, '0')}`
}

const filterHistory = computed(() => {
  if (!search.value) return playHistory.value
  return playHistory.value.filter((h) => {
    return h.episode_title.includes(search.value) || h.name.includes(search.value)
  })
})

// 格式化播放时间
const formatLastPlayTime = (timestamp: number): string => {
  const now = Date.now()
  const diff = now - timestamp
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else {
    return `${days}天前`
  }
}

// 计算进度百分比
const getProgressPercent = (progress: number, duration: number): number => {
  return Math.min((progress / duration) * 100, 100)
}

// 删除单条记录
const deleteHistory = (vod_id: number, episode_index: number, sourceKey: string) => {
  playHistory.value = playHistory.value.filter(
    (h) => !(h.vod_id === vod_id && h.episode_index === episode_index && h.sourceKey === sourceKey),
  )
}

// 清空所有记录
const clearAllHistory = () => {
  playHistory.value = []
}

// 关闭对话框
const closeDialog = () => {
  isOpen.value = false
}
</script>

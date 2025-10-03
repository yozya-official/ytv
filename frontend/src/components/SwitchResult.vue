<template>
  <div class="flex flex-col overflow-hidden">
    <!-- 搜索头部区域 -->
    <div class="border-b border-border">
      <div class="max-w-7xl mx-auto px-6 py-6 space-y-4">
        <h2 class="text-center">更换资源</h2>

        <!-- 搜索框 -->
        <div class="relative group">
          <div
            class="flex items-center gap-3 bg-card border border-border rounded-xl px-4 py-2 shadow-sm hover:shadow-md transition-all duration-300 focus-within:border-primary focus-within:ring-2 focus-within:ring-primary/20"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="size-5 text-muted-foreground"
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
              type="search"
              v-model.lazy="search"
              class="flex-1 bg-transparent border-none outline-none text-sm placeholder:text-muted-foreground"
              placeholder="搜索名称、年代、来源..."
            />

            <div class="flex items-center gap-3">
              <span
                class="hidden sm:block text-sm text-muted-foreground font-medium px-3 py-1 bg-muted rounded-full"
              >
                {{ filteredVods.length }} 个结果
              </span>
              <button
                class="px-4 py-2 bg-primary text-primary-foreground rounded-lg font-medium hover:bg-primary/90 transition-colors"
              >
                搜索
              </button>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <DataPagination
          :total="filteredVods.length"
          v-model:current-page="page"
          v-model:page-size="pageSize"
        />
      </div>
    </div>
    <!-- 视频网格 -->
    <div class="max-w-7xl mx-auto px-4 py-6 flex-1 overflow-y-scroll">
      <div
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4"
      >
        <div
          v-for="vod in pagedVods"
          :key="vod.vod_id"
          @click="selectVod(vod)"
          :class="{ 'pointer-events-none': isCurrentVod(vod) }"
          class="group cursor-pointer"
        >
          <!-- 卡片容器 -->
          <div
            class="bg-card rounded-xl overflow-hidden border border-border hover:border-primary/50 shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col h-full"
          >
            <!-- 图片区域 -->
            <div class="relative aspect-[2/3] overflow-hidden bg-muted">
              <img
                class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                :src="vod.vod_pic_thumb || vod.vod_pic"
                :alt="vod.vod_name"
                loading="lazy"
              />

              <!-- 当前播放浮层 -->
              <div
                v-if="isCurrentVod(vod)"
                class="absolute inset-0 via-primary/85 to-primary/95 backdrop-blur-sm flex items-center justify-center"
              >
                <div class="flex flex-col items-center gap-3 animate-pulse-slow">
                  <!-- 播放图标 -->
                  <div class="relative">
                    <div class="absolute inset-0 bg-white/30 rounded-full blur-xl"></div>
                    <div
                      class="relative bg-white/20 backdrop-blur-md rounded-full p-4 border-2 border-white/50 shadow-2xl"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="size-8 text-white"
                        viewBox="0 0 24 24"
                        fill="currentColor"
                      >
                        <path d="M8 5v14l11-7z" />
                      </svg>
                    </div>
                  </div>
                  <!-- 文字提示 -->
                  <div
                    class="bg-white/20 backdrop-blur-md px-4 py-2 rounded-full border border-white/30 shadow-lg"
                  >
                    <span class="text-white font-bold text-sm tracking-wide">正在播放</span>
                  </div>
                </div>
              </div>

              <!-- 年份标签 -->
              <div class="absolute top-2 right-2">
                <span
                  class="px-2.5 py-1 rounded-lg text-xs font-bold bg-primary/90 text-primary-foreground backdrop-blur-sm shadow-lg"
                >
                  {{ vod.vod_year }}
                </span>
              </div>

              <!-- 悬浮操作按钮 -->
              <div
                class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300"
              >
                <div class="bg-background/90 backdrop-blur-sm rounded-full p-3 shadow-lg">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="size-6 text-primary"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                  >
                    <path d="M8 5v14l11-7z" />
                  </svg>
                </div>
              </div>
            </div>

            <!-- 信息区域 -->
            <div class="p-3 flex flex-col gap-2 flex-1">
              <!-- 标题 -->

              <h3
                class="text-sm font-semibold text-foreground overflow-hidden line-clamp-2 group-hover:text-primary transition-colors min-h-[2.5rem]"
              >
                {{ vod.vod_name }}
              </h3>

              <!-- 来源标签 -->
              <div class="flex items-center gap-1.5 mt-auto">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="size-3.5 text-muted-foreground flex-shrink-0"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"
                  />
                </svg>
                <span class="text-xs text-muted-foreground truncate font-medium">
                  {{ vod.source_name }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div
        v-if="pagedVods.length === 0"
        class="flex flex-col items-center justify-center py-20 text-center"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="size-20 text-muted-foreground/40 mb-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
        <h3 class="text-lg font-semibold text-foreground mb-2">未找到相关内容</h3>
        <p class="text-sm text-muted-foreground">尝试使用不同的关键词搜索</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { VodItem } from '@/models'
import router from '@/router'

const search = ref('')

const page = ref(1)
const pageSize = ref(24)

const selectVod = async (vod: VodItem) => {
  await router.push({
    name: 'player',
    query: {
      vodId: vod.vod_id,
      sourceKey: vod.source_key,
      episodeIndex: 0,
    },
  })
  location.reload()
}

const props = defineProps<{
  containerClass?: string
  sourceKey: String
  vodId: Number
  searchResults: VodItem[]
}>()

const isCurrentVod = (vod: VodItem) => {
  return vod.vod_id === props.vodId && vod.source_key === props.sourceKey
}

// 把当前资源放最开头
const filteredVods = computed(() => {
  const currentVod = props.searchResults.find((vod) => isCurrentVod(vod))
  const excludeCurrentResult = props.searchResults.filter(
    (vod) => !(vod.vod_id === props.vodId && vod.source_key === props.sourceKey),
  )

  if (!search.value) {
    return currentVod ? [currentVod, ...excludeCurrentResult] : excludeCurrentResult
  }

  page.value = 1

  const filtered = excludeCurrentResult.filter((vod) => {
    return (
      vod.vod_year === search.value ||
      vod.source_name.includes(search.value) ||
      vod.vod_name.includes(search.value)
    )
  })

  return currentVod ? [currentVod, ...filtered] : filtered
})

const pagedVods = computed(() => {
  const start = (page.value - 1) * pageSize.value
  const end = start + pageSize.value

  return filteredVods.value.slice(start, end)
})
</script>

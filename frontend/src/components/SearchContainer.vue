<template>
  <div v-if="loading" class="space-y-8 px-6 py-8">
    <div class="skeleton h-6 w-32"></div>
    <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-6">
      <div v-for="j in 24" :key="j" class="space-y-3">
        <div class="skeleton aspect-[2/3] w-full rounded-lg"></div>
        <div class="skeleton h-4 w-full"></div>
        <div class="skeleton h-3 w-2/3"></div>
      </div>
    </div>
  </div>

  <div v-else class="px-6 py-8">
    <div v-if="vods.length">
      <div class="space-y-8">
        <!-- 搜索结果提示 -->
        <div
          class="flex items-center gap-4 px-5 py-3.5 rounded-xl bg-gradient-to-r from-primary/5 via-primary/3 to-transparent text-sm max-w-sm border border-primary/10"
        >
          <div
            class="flex-shrink-0 w-11 h-11 rounded-xl bg-primary/10 flex items-center justify-center ring-4 ring-primary/5"
          >
            <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
              />
            </svg>
          </div>
          <div class="flex flex-col gap-0.5">
            <span class="text-xs text-muted-foreground font-medium uppercase tracking-wide"
              >搜索结果</span
            >
            <span class="text-foreground font-bold text-lg"
              >{{ vods.length }}
              <span class="text-sm font-normal text-muted-foreground">个视频</span></span
            >
          </div>
        </div>

        <!-- 视频网格 -->
        <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-6">
          <div
            class="bg-card rounded-xl shadow-sm hover:shadow-xl hover:-translate-y-1.5 transition-all duration-300 overflow-hidden group cursor-pointer border border-border hover:border-primary/40 hover:ring-2 hover:ring-primary/20"
            v-for="vod in pagedVideos"
            :key="vod.vod_id"
          >
            <SearchCard @click="selectVod(vod)" :vod="vod" />
          </div>
        </div>

        <!-- 分页 -->
        <div class="pt-8 pb-4">
          <DataPagination
            class="mx-auto"
            :total="vods.length"
            v-model:current-page="page"
            v-model:page-size="pageSize"
          ></DataPagination>
        </div>
      </div>

      <!-- 详情对话框 - 响应式优化 -->
      <dialog class="dialog backdrop:bg-black/60 backdrop:backdrop-blur-sm" ref="vodInfoRef">
        <div
          class="dialog-body max-w-4xl w-full max-h-[95vh] md:max-h-[92vh] bg-background rounded-t-3xl md:rounded-3xl shadow-2xl overflow-hidden"
        >
          <!-- 顶部栏 - 固定 -->
          <div
            class="sticky top-0 z-50 bg-background/95 backdrop-blur-md border-b border-border px-4 md:px-6 py-4 flex items-center justify-between"
          >
            <div class="flex items-center gap-2 md:gap-3">
              <div
                class="w-8 h-8 md:w-10 md:h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0"
              >
                <svg
                  class="w-4 h-4 md:w-5 md:h-5 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z"
                  />
                </svg>
              </div>
              <div class="hidden md:block">
                <h3 class="text-sm font-medium text-muted-foreground">视频详情</h3>
                <p class="text-xs text-muted-foreground/70">Video Details</p>
              </div>
              <h3 class="md:hidden text-sm font-medium text-foreground">详情</h3>
            </div>

            <button
              @click="vodInfoRef?.close()"
              class="w-9 h-9 md:w-10 md:h-10 rounded-full bg-muted/50 hover:bg-muted flex items-center justify-center transition-colors group flex-shrink-0"
            >
              <svg
                class="w-5 h-5 text-muted-foreground group-hover:text-foreground transition-colors"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

          <!-- 内容区域 - 可滚动 -->
          <div
            class="overflow-y-auto max-h-[calc(95vh-65px)] md:max-h-[calc(92vh-73px)] px-4 md:px-6 lg:px-8 py-5 md:py-6 lg:py-8"
          >
            <div v-if="selectedVod" class="space-y-6 md:space-y-8 lg:space-y-10">
              <!-- 标题和标签区域 -->
              <div class="space-y-3 md:space-y-4">
                <h1
                  class="text-2xl md:text-3xl lg:text-4xl font-bold text-foreground leading-tight tracking-tight"
                >
                  {{ selectedVod.vod_name }}
                </h1>
                <div class="flex flex-wrap items-center gap-2">
                  <span
                    v-if="selectedVod.source_name"
                    class="px-3 py-1 text-xs md:text-sm font-medium text-primary bg-primary/10 rounded-full border border-primary/20"
                  >
                    {{ selectedVod.source_name }}
                  </span>
                  <span
                    v-if="selectedVod.type_name"
                    class="px-3 py-1 text-xs md:text-sm font-medium text-foreground/80 bg-muted/80 rounded-full"
                  >
                    {{ selectedVod.type_name }}
                  </span>
                  <span
                    v-if="selectedVod.vod_year"
                    class="px-3 py-1 text-xs md:text-sm font-medium text-foreground/80 bg-muted/80 rounded-full"
                  >
                    {{ selectedVod.vod_year }}
                  </span>
                </div>
              </div>

              <!-- 基本信息卡片 -->
              <div
                class="bg-gradient-to-br from-muted/40 to-muted/20 rounded-xl md:rounded-2xl p-4 md:p-6 lg:p-8 border border-border/50"
              >
                <div class="grid grid-cols-1 gap-4 md:gap-6 lg:gap-8">
                  <div class="space-y-3 md:space-y-4">
                    <div class="flex items-start gap-3">
                      <div
                        class="w-9 h-9 md:w-10 md:h-10 rounded-lg md:rounded-xl bg-background flex items-center justify-center flex-shrink-0"
                      >
                        <svg
                          class="w-4 h-4 md:w-5 md:h-5 text-primary"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                          />
                        </svg>
                      </div>
                      <div class="flex-1 space-y-0.5 md:space-y-1 min-w-0">
                        <p
                          class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
                        >
                          分类
                        </p>
                        <p class="text-sm md:text-base text-foreground font-medium truncate">
                          {{ selectedVod.type_name || '暂无' }}
                        </p>
                      </div>
                    </div>

                    <div class="flex items-start gap-3">
                      <div
                        class="w-9 h-9 md:w-10 md:h-10 rounded-lg md:rounded-xl bg-background flex items-center justify-center flex-shrink-0"
                      >
                        <svg
                          class="w-4 h-4 md:w-5 md:h-5 text-primary"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                          />
                        </svg>
                      </div>
                      <div class="flex-1 space-y-0.5 md:space-y-1 min-w-0">
                        <p
                          class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
                        >
                          年份
                        </p>
                        <p class="text-sm md:text-base text-foreground font-medium">
                          {{ selectedVod.vod_year || '暂无' }}
                        </p>
                      </div>
                    </div>

                    <div class="flex items-start gap-3">
                      <div
                        class="w-9 h-9 md:w-10 md:h-10 rounded-lg md:rounded-xl bg-background flex items-center justify-center flex-shrink-0"
                      >
                        <svg
                          class="w-4 h-4 md:w-5 md:h-5 text-primary"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                          />
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                          />
                        </svg>
                      </div>
                      <div class="flex-1 space-y-0.5 md:space-y-1 min-w-0">
                        <p
                          class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
                        >
                          地区
                        </p>
                        <p class="text-sm md:text-base text-foreground font-medium">
                          {{ selectedVod.vod_area || '暂无' }}
                        </p>
                      </div>
                    </div>

                    <div class="flex items-start gap-3">
                      <div
                        class="w-9 h-9 md:w-10 md:h-10 rounded-lg md:rounded-xl bg-background flex items-center justify-center flex-shrink-0"
                      >
                        <svg
                          class="w-4 h-4 md:w-5 md:h-5 text-primary"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
                          />
                        </svg>
                      </div>
                      <div class="flex-1 space-y-0.5 md:space-y-1 min-w-0">
                        <p
                          class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
                        >
                          导演
                        </p>
                        <p class="text-sm md:text-base text-foreground font-medium truncate">
                          {{ selectedVod.vod_director || selectedVod.vod_remarks || '暂无' }}
                        </p>
                      </div>
                    </div>

                    <div class="flex items-start gap-3">
                      <div
                        class="w-9 h-9 md:w-10 md:h-10 rounded-lg md:rounded-xl bg-background flex items-center justify-center flex-shrink-0"
                      >
                        <svg
                          class="w-4 h-4 md:w-5 md:h-5 text-primary"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                          />
                        </svg>
                      </div>
                      <div class="flex-1 space-y-0.5 md:space-y-1 min-w-0">
                        <p
                          class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
                        >
                          主演
                        </p>
                        <p class="text-sm md:text-base text-foreground font-medium leading-relaxed">
                          {{ selectedVod.vod_actor || '暂无' }}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 简介卡片 -->
              <div class="space-y-3 md:space-y-4">
                <div class="flex items-center gap-2 md:gap-3">
                  <div class="w-1 h-5 md:h-6 bg-primary rounded-full"></div>
                  <h3 class="text-base md:text-lg font-semibold text-foreground">剧情简介</h3>
                </div>
                <div
                  class="bg-gradient-to-br from-muted/30 to-muted/10 rounded-xl md:rounded-2xl p-4 md:p-6 lg:p-8 border border-border/50"
                >
                  <p
                    v-if="selectedVod.vod_content"
                    class="text-foreground/90 leading-relaxed md:leading-loose text-sm md:text-base whitespace-pre-wrap"
                  >
                    {{ stripHtml(selectedVod.vod_content) }}
                  </p>
                  <p
                    v-else-if="selectedVod.vod_blurb"
                    class="text-foreground/90 leading-relaxed md:leading-loose text-sm md:text-base whitespace-pre-wrap"
                  >
                    {{ selectedVod.vod_blurb }}
                  </p>
                  <p
                    v-else
                    class="text-muted-foreground italic text-center py-6 md:py-8 text-sm md:text-base"
                  >
                    暂无简介信息
                  </p>
                </div>
              </div>

              <!-- 剧集列表 -->
              <div class="space-y-3 md:space-y-4 pb-4 md:pb-6 lg:pb-8">
                <div class="flex items-center gap-2 md:gap-3">
                  <div class="w-1 h-5 md:h-6 bg-primary rounded-full"></div>
                  <h3 class="text-base md:text-lg font-semibold text-foreground">播放列表</h3>
                </div>
                <EpisodeList
                  v-model:select-vod="selectedVod"
                  :key="selectedVod.vod_id + selectedVod.source_name"
                ></EpisodeList>
              </div>
            </div>
          </div>
        </div>
      </dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ vods: VodItem[]; loading: boolean }>()

const vodInfoRef = ref<HTMLDialogElement>()

const selectedVod = ref<VodItem>()

const page = ref(1)
const pageSize = ref(16)

const stripHtml = (html: string) => {
  const div = document.createElement('div')
  div.innerHTML = html
  return div.innerText.trim() || div.textContent?.trim() || ''
}

const selectVod = (vod: VodItem) => {
  selectedVod.value = vod
  vodInfoRef.value?.show()
}

const pagedVideos = computed(() => {
  const start = (page.value - 1) * pageSize.value
  const end = start + pageSize.value

  return props.vods.slice(start, end)
})
</script>

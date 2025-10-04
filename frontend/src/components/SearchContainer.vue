<template>
  <div v-if="loading" class="space-y-8">
    <div class="skeleton h-6 w-32"></div>
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
      <div v-for="j in 24" :key="j" class="space-y-2">
        <div class="skeleton aspect-[2/3] w-full"></div>
        <div class="skeleton h-4 w-full"></div>
        <div class="skeleton h-3 w-2/3"></div>
      </div>
    </div>
  </div>

  <div v-else>
    <div v-if="vods.length">
      <!-- 搜索结果列表 -->
      <div class="space-y-6">
        <!-- 视频网格 -->
        <div class="flex items-center gap-3 px-4 py-2 rounded-lg bg-muted text-sm mb-4">
          <svg
            class="w-4 h-4 text-muted-foreground"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
            />
          </svg>
          <span class="text-muted-foreground">
            找到 <span class="font-semibold text-foreground">{{ vods.length }}</span> 个结果
          </span>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          <div
            class="bg-card rounded-lg shadow-md hover:shadow-lg transition-all duration-300 overflow-hidden group cursor-pointer border border-border"
            v-for="vod in pagedVideos"
            :key="vod.vod_id"
          >
            <SearchCard @click="selectVod(vod)" :vod="vod" />
          </div>
        </div>
        <DataPagination
          class="mx-auto"
          :total="vods.length"
          v-model:current-page="page"
          v-model:page-size="pageSize"
        ></DataPagination>
      </div>

      <dialog class="dialog" ref="vodInfoRef">
        <div class="dialog-body max-w-4xl w-full max-h-[90vh] overflow-y-auto">
          <!-- 关闭按钮 -->
          <button @click="vodInfoRef?.close()" class="btn-close">
            <svg class="size-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>

          <div v-if="selectedVod" class="space-y-6">
            <!-- 标题 -->
            <div class="space-y-2">
              <h2 class="text-2xl font-bold text-foreground">
                {{ selectedVod.vod_name }}
                <span v-if="selectedVod.vod_remarks" class="text-sm font-normal text-primary ml-2">
                  ({{ selectedVod.source_name }})
                </span>
              </h2>
            </div>

            <!-- 基本信息表格 -->
            <div class="grid grid-cols-2 gap-x-8 gap-y-3 text-sm">
              <div class="flex">
                <span class="text-muted-foreground w-16 flex-shrink-0">类型:</span>
                <span class="text-foreground">{{ selectedVod.type_name }}</span>
              </div>

              <div class="flex">
                <span class="text-muted-foreground w-16 flex-shrink-0">年份:</span>
                <span class="text-foreground">{{ selectedVod.vod_year }}</span>
              </div>

              <div class="flex">
                <span class="text-muted-foreground w-16 flex-shrink-0">地区:</span>
                <span class="text-foreground">{{ selectedVod.vod_area }}</span>
              </div>

              <div class="flex">
                <span class="text-muted-foreground w-16 flex-shrink-0">导演:</span>
                <span class="text-foreground">{{ selectedVod.vod_remarks }}</span>
              </div>

              <div class="flex col-span-2">
                <span class="text-muted-foreground w-16 flex-shrink-0">主演:</span>
                <span class="text-foreground">{{ selectedVod.vod_actor || '暂无' }}</span>
              </div>

              <div class="flex col-span-2">
                <span class="text-muted-foreground w-16 flex-shrink-0">备注:</span>
                <span class="text-primary">{{ selectedVod.type_name }}</span>
              </div>
            </div>

            <div class="border"></div>

            <!-- 简介 -->
            <div class="space-y-2">
              <div class="text-sm text-muted-foreground">简介:</div>

              <div v-if="selectedVod.vod_content">{{ stripHtml(selectedVod.vod_content) }}</div>

              <div v-else-if="selectedVod.vod_blurb">{{ selectedVod.vod_blurb }}</div>

              <div v-else>{{ '暂无简介' }}</div>
            </div>

            <EpisodeList
              v-model:select-vod="selectedVod"
              :key="selectedVod.vod_id + selectedVod.source_name"
            ></EpisodeList>
          </div>
        </div>
      </dialog>
    </div>
    <!-- 空状态 -->
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ vods: VodItem[]; loading: boolean }>()

const vodInfoRef = ref<HTMLDialogElement>()

const selectedVod = ref<VodItem>()

const page = ref(1)
const pageSize = ref(24)

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

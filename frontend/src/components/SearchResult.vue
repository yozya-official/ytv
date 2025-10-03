<template>
  <!-- 搜索结果列表 -->
  <div class="space-y-6">
    <!-- 视频网格 -->
    <div>{{ vods.length }} 个结果</div>

    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
      <div
        class="bg-card rounded-lg shadow-md hover:shadow-lg transition-all duration-300 overflow-hidden group cursor-pointer border border-border"
        v-for="vod in pagedVideos"
        :key="vod.vod_id"
      >
        <div @click="selectVod(vod)"><SearchCard :vod="vod" /></div>
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
      <button
        @click="vodInfoRef?.close()"
        class="absolute top-4 right-4 text-muted-foreground hover:text-foreground transition-colors z-50"
      >
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
</template>

<script setup lang="ts">
import type { VodItem } from '@/models'

const vodInfoRef = ref<HTMLDialogElement>()

const selectedVod = ref<VodItem>()

const props = defineProps<{ vods: VodItem[]; containerClass?: string }>()

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

onMounted(() => {
  selectedVod.value = undefined
})
</script>

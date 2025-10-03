<template>
  <!-- 剧集列表 -->
  <div class="space-y-3">
    <div class="flex items-center gap-3 text-sm">
      <div class="badge badge-accent">{{ selectVod.source_name }}</div>

      <!-- 分隔线 -->
      <div class="h-4 w-px border"></div>
      <!-- 视频图标 -->
      <div class="flex items-center gap-1.5 text-muted-foreground">
        <svg class="size-4" fill="currentColor" viewBox="0 0 20 20">
          <path
            d="M2 6a2 2 0 012-2h6a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z"
          />
        </svg>
        <span class="font-medium">共 {{ selectedVod.episodes.length }} 集</span>
      </div>

      <!-- 分隔线 -->
      <div class="h-4 w-px border"></div>
      <!-- 排序按钮 -->
      <button
        class="btn btn-ghost flex items-center gap-1.5 px-3 py-1.5 rounded-lg"
        @click="episodeReverse = !episodeReverse"
      >
        <!-- 排序图标 -->
        <svg
          class="size-4 transition-transform duration-200"
          :class="{ 'rotate-180': episodeReverse }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
          />
        </svg>
        <span class="font-medium">{{ episodeReverse ? '倒序' : '正序' }}</span>
      </button>
    </div>

    <!-- 剧集按钮网格 -->
    <div class="grid grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-2">
      <button
        v-for="episode in pagesEpisodes"
        :key="episode.url"
        class="btn btn-outline text-sm py-2 hover:btn-primary transition-colors"
        :class="{ 'btn-primary': episodeIndex == episode.episode_index }"
        @click="playEpisode(episode)"
      >
        {{ episode.episode_title }}
      </button>
    </div>

    <DataPagination
      :total="selectedVod.episodes.length"
      v-model:current-page="episodePage"
      v-model:page-size="episodePageSize"
    ></DataPagination>
  </div>
</template>

<script setup lang="ts">
import type { Episode, VodItem } from '@/models'
import router from '@/router'

const playEpisode = (episode: Episode) => {
  if (!selectedVod.value) {
    console.warn('没有选中的视频')
    return
  }

  router.push({
    name: 'player',
    query: {
      vodId: selectedVod.value.vod_id,
      sourceKey: selectedVod.value.source_key,
      episodeIndex: episode.episode_index,
    },
  })
}

const selectedVod = defineModel<VodItem>('select-vod', { required: true })

const props = defineProps<{ episodeIndex?: number }>()

const episodeReverse = ref(false)

const episodePage = ref(1)
const episodePageSize = ref(48)

const pagesEpisodes = computed(() => {
  if (!selectedVod.value?.episodes) return []

  const episodes = [...selectedVod.value.episodes]

  if (episodeReverse.value) {
    episodes.reverse()
  }

  // 分页切片
  const start = (episodePage.value - 1) * episodePageSize.value
  const end = episodePage.value * episodePageSize.value
  return episodes.slice(start, end)
})

onMounted(() => {
  episodePage.value = props.episodeIndex ? Math.ceil(props.episodeIndex / episodePageSize.value) : 1
})
</script>

<template>
  <!-- 热搜区域 -->
  <div class="mt-16 max-w-6xl mx-auto">
    <div class="p-8 card-deco card-deco-primary card my-8 shadow-xl relative overflow-hidden">
      <div class=""></div>

      <div class="flex items-center justify-between">
        <div class="flex gap-1.5">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="size-8 text-primary/70"
            viewBox="0 0 24 24"
            fill="currentColor"
          >
            <path
              d="M18 3v2h-2V3H8v2H6V3H4v18h2v-2h2v2h8v-2h2v2h2V3h-2zM8 17H6v-2h2v2zm0-4H6v-2h2v2zm0-4H6V7h2v2zm10 8h-2v-2h2v2zm0-4h-2v-2h2v2zm0-4h-2V7h2v2z"
            />
          </svg>
          <h2 class="text-2xl font-bold">豆瓣热门</h2>
        </div>

        <button class="btn btn-destructive" @click="randPage">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="size-4"
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
              <path d="M21 12a9 9 0 0 0-9-9a9.75 9.75 0 0 0-6.74 2.74L3 8" />
              <path d="M3 3v5h5m-5 4a9 9 0 0 0 9 9a9.75 9.75 0 0 0 6.74-2.74L21 16" />
              <path d="M16 16h5v5" />
            </g>
          </svg>
          <span class="pl-1.5">换一批</span>
        </button>
      </div>

      <!-- 类型切换 -->
      <div class="flex justify-between">
        <div role="tablist" class="tabs tabs-box w-max my-4">
          <button
            class="btn btn-sm"
            v-for="type in contentTypes"
            :key="type.value"
            @click="selectedType = type.value"
            :class="['btn btn-sm', selectedType === type.value ? 'btn-primary' : 'btn-ghost']"
          >
            {{ type.label }}
          </button>
        </div>
      </div>
      <!-- 标签列表 -->
      <div class="flex flex-wrap gap-2">
        <div v-for="tag in tags" :key="tag" class="relative group">
          <button
            @click="selectedTag = tag"
            :class="[
              'btn btn-sm ',
              selectedTag === tag ? 'btn-primary' : 'bg-muted btn-outline border',
            ]"
          >
            {{ tag }}
          </button>

          <!-- 删除按钮 -->
          <button
            v-if="!defaultTags.includes(tag)"
            @click="removeTag(tag)"
            class="absolute -top-2 -right-2 btn btn-circle btn-xs btn-error opacity-60 group-hover:opacity-100 transition-opacity"
            title="删除标签"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="size-3"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        </div>
        <!-- 添加标签按钮 -->
        <button @click="addTagDialog?.show" class="btn btn-sm border btn-outline text-primary">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="size-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M12 5v14M5 12h14" />
          </svg>
          <span class="text-sm">添加标签</span>
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="space-y-8">
      <div class="skeleton h-6 w-32"></div>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-8">
        <div v-for="i in pageSize" :key="i" class="space-y-2">
          <div class="skeleton aspect-[2/3] w-full"></div>
          <div class="skeleton h-4 w-full"></div>
          <div class="skeleton h-3 w-2/3"></div>
        </div>
      </div>
    </div>

    <!-- 热搜列表 -->
    <div
      v-else-if="hotMovies.length > 0"
      class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-8"
    >
      <div
        v-for="(movie, index) in hotMovies"
        :key="index"
        @click="sendTitle(movie.title)"
        class="group relative card-hover card transition-all duration-300 cursor-pointer overflow-hidden"
      >
        <!-- 排名角标 -->
        <div class="absolute top-3 right-3 z-20 drop-shadow-2xl">
          <div class="relative">
            <div
              class="relative badge badge-sm bg-gradient-to-br from-background/90 to-background/75 backdrop-blur-md border border-amber-500/30 font-bold shadow-lg"
            >
              <div class="flex items-center gap-1.5">
                <!-- 评分星星 -->
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="size-4 text-amber-500"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
                  />
                </svg>
                <span
                  class="text-sm font-bold bg-gradient-to-r from-amber-600 to-amber-400 bg-clip-text text-transparent"
                >
                  {{ movie.rate || 'N/A' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 海报图片 -->
        <figure class="relative aspect-[2/3] overflow-hidden">
          <img
            :src="movie.cover_base64 || movie.cover"
            :alt="movie.title"
            referrerpolicy="no-referrer"
            class="group-hover:scale-110 h-full transition-transform duration-300"
          />

          <!-- 悬浮操作按钮 -->
          <div class="image-mask flex items-center justify-center">
            <div class="bg-background/90 backdrop-blur-sm rounded-full p-3 shadow-lg">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="size-12 text-primary"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path d="M8 5v14l11-7z" />
              </svg>
            </div>
          </div>
        </figure>

        <!-- 影片信息 -->

        <h3
          class="font-bold p-2 text-center text-sm line-clamp-2 min-h-[2.5rem]"
          :title="movie.title"
        >
          {{ movie.title }}
        </h3>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="text-center py-12">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="size-16 mx-auto mb-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z"
        />
      </svg>
      <p class="text-lg">暂无数据</p>
    </div>
  </div>

  <!-- 添加标签弹窗 -->
  <dialog class="dialog" ref="addTagDialog" @click.self="addTagDialog?.close()">
    <div class="dialog-body p-5 gap-2.5">
      <div class="flex items-center gap-2 mb-4">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="size-6 text-primary"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"
          />
          <line x1="7" y1="7" x2="7.01" y2="7" />
        </svg>
        <h3 class="font-bold text-left text-lg">添加新标签</h3>
      </div>

      <input
        v-model="newTag"
        type="text"
        placeholder="输入标签名称..."
        class="input input-primary w-full"
        @keyup.enter="addTag"
      />

      <div class="flex gap-2 p-2 justify-between">
        <button @click="addTagDialog?.close()" class="btn btn-ghost">取消</button>
        <button @click="addTag" class="btn btn-primary" :disabled="!newTag.trim()">添加</button>
      </div>
    </div>
  </dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useStorage } from '@vueuse/core'
import { videoApi } from '@/api'

const loading = ref(false)
const addTagDialog = ref<HTMLDialogElement>()
const newTag = ref('')

// 内容类型
const contentTypes = [
  { label: '电影', value: 'movie' },
  { label: '电视剧', value: 'tv' },
]

const selectedType = ref('movie')
const page = ref(1)
const pageSize = ref(12)

// 默认标签（不可删除）
const defaultTags = ['热门', '美剧', '英剧', '日剧', '韩剧', '港剧', '综艺', '纪录片']

const randPage = () => {
  page.value = Math.floor(Math.random() * 100) + 1
}

// 从 localStorage 读取自定义标签
const customTags = useStorage<string[]>('custom-tags', [])

// 合并默认标签和自定义标签
const tags = computed(() => [...defaultTags, ...customTags.value])

const selectedTag = ref('热门')

// 热搜数据
const hotMovies = ref<HotMovie[]>([])

// 添加标签
const addTag = () => {
  const trimmedTag = newTag.value.trim()
  if (trimmedTag && !tags.value.includes(trimmedTag)) {
    customTags.value.push(trimmedTag)
    newTag.value = ''
    addTagDialog.value?.close()
  }
}

// 删除标签
const removeTag = (tag: string) => {
  if (defaultTags.includes(tag)) return

  customTags.value = customTags.value.filter((t) => t !== tag)

  // 如果删除的是当前选中的标签，切换到"热门"
  if (selectedTag.value === tag) {
    selectedTag.value = '热门'
  }
}

const emit = defineEmits(['send-title'])

const sendTitle = (title: string) => {
  emit('send-title', title)
}

// 监听hot参数 页码,类别,tag
watch(
  [page, pageSize, selectedType, selectedTag],
  async () => {
    try {
      loading.value = true

      // 先请求主要信息
      const res = await videoApi.getHotVideos(
        page.value,
        pageSize.value,
        selectedType.value,
        selectedTag.value,
      )
      hotMovies.value = res.data.data.list

      console.log(hotMovies.value)
    } catch (err) {
      console.error('获取热门影片失败:', err)
      hotMovies.value = []
    } finally {
      loading.value = false
    }
  },
  { immediate: true, deep: true },
)
</script>

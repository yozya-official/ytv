<template>
  <PageHeader title="YTV" subtitle="自由观影,畅享精彩">
    <template #icon>
      <svg class="size-8" fill="white" viewBox="0 0 20 20">
        <path
          d="M2 6a2 2 0 012-2h6a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z"
        ></path>
      </svg>
    </template>
  </PageHeader>

  <div class="container mx-auto px-4 pt-12 pb-8">
    <!-- 搜索框 -->
    <div class="flex items-center justify-center mb-12">
      <div class="search-input search-input-primary max-w-2xl w-full shadow-md relative">
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
          @keydown.enter="handleSearch"
          @focus="showHistory = true"
          class="flex-1 bg-transparent border-none outline-none text-sm placeholder:text-muted-foreground"
          placeholder="搜索影视剧、综艺、动漫..."
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
        <button
          class="btn btn-primary btn-sm"
          @click="handleSearch"
          :disabled="!search.trim() || loading"
        >
          <span v-if="!loading">搜索</span>
          <span v-else class="flex items-center gap-2">
            <span class="spinner spinner-sm"></span>
            搜索中
          </span>
        </button>

        <!-- 搜索历史下拉框 -->
        <div
          v-if="showHistory && baseStore.recentSearchList.length > 0 && !search"
          class="absolute top-full left-0 right-0 mt-2 bg-background border border-border rounded-lg shadow-lg z-50"
        >
          <div class="flex items-center justify-between px-4 py-3 border-b border-border">
            <span class="text-sm font-medium text-foreground">最近搜索</span>
            <button
              @click.stop="handleClearHistory"
              class="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              清空
            </button>
          </div>
          <div class="py-2">
            <div
              v-for="keyword in baseStore.recentSearchList"
              :key="keyword"
              class="flex items-center justify-between px-4 py-2 hover:bg-muted/50 cursor-pointer group"
              @mousedown.prevent="handleHistoryClick(keyword)"
            >
              <div class="flex items-center gap-3 flex-1" @click="handleHistoryClick(keyword)">
                <svg
                  class="size-4 text-muted-foreground flex-shrink-0"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <span class="text-sm text-foreground">{{ keyword }}</span>
              </div>
              <button
                @click.stop="handleRemoveHistory(keyword)"
                class="opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <svg
                  class="size-4 text-muted-foreground hover:text-foreground"
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
          </div>
        </div>
      </div>
    </div>

    <!-- 热门视频 -->
    <HotVideos v-if="!search && !searchResults.length" @send-title="handleTitle" />
    <div v-else-if="failed" class="flex flex-col items-center justify-center py-20">
      <div class="card max-w-md w-full text-center">
        <div class="card-content p-6">
          <svg
            class="size-16 mx-auto mb-4 text-muted-foreground"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <h3 class="card-title text-lg mb-2">未找到相关内容</h3>
          <p class="card-description">试试搜索其他关键词吧</p>
          <button class="btn btn-primary mt-4" @click="reload">回到主页</button>
        </div>
      </div>
    </div>

    <SearchContainer v-else :loading="loading" :results="searchResults"></SearchContainer>
  </div>
</template>

<script setup lang="ts">
const baseStore = useBaseStore()

const searchResults = ref<VodItem[]>([])

const search = ref('')
const showHistory = ref(false)

const loading = ref(false)
const failed = ref(false)

// 执行搜索
const handleSearch = async () => {
  const trimmedInput = search.value.trim()

  // 如果输入为空,不执行搜索
  if (!trimmedInput) {
    return
  }

  searchResults.value = []

  search.value = trimmedInput
  showHistory.value = false

  // 添加到搜索历史
  baseStore.addSearch(trimmedInput)

  try {
    loading.value = true
    const resp = await videoApi.searchAll(search.value)
    searchResults.value = resp.data.data.list
    if (searchResults.value.length == 0) {
      failed.value = true
    } else {
      failed.value = false
    }
  } catch (err: unknown) {
    handleApiError(err)
    failed.value = true
  } finally {
    loading.value = false
  }
}

const reload = () => {
  search.value = ''
  location.reload()
}

// 处理热门视频点击
const handleTitle = async (title: string) => {
  search.value = title
  await handleSearch()
}

// 点击历史记录
const handleHistoryClick = (keyword: string) => {
  search.value = keyword
  showHistory.value = false
  handleSearch()
}

// 删除单条历史记录
const handleRemoveHistory = (keyword: string) => {
  baseStore.removeSearch(keyword)
}

// 清空所有历史记录
const handleClearHistory = () => {
  baseStore.clearSearch()
  showHistory.value = false
}
</script>

<template>
  <PageHeader title="YTV" subtitle="自由观影,畅享精彩" class="cursor-pointer" @click="reload()">
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
      <SearchInput
        v-model:history="baseStore.recentSearchList"
        :placeholder="'搜索影视剧、综艺、动漫...'"
        @search="handleSearch"
        :loading="loading"
        :enable-history="true"
        class="search-input search-input-primary max-w-2xl w-full"
      ></SearchInput>
    </div>

    <!-- 热门视频 -->
    <HotVideos v-if="!searchResults.length" @search-movie="handleSearchHot" />
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

    <SearchContainer v-else :loading="loading" :vods="searchResults"></SearchContainer>
  </div>
</template>

<script setup lang="ts">
const baseStore = useBaseStore()

const searchResults = ref<VodItem[]>([])

const showHistory = ref(false)

const loading = ref(false)
const failed = ref(false)

// 执行搜索
const handleSearch = async (keyword: string) => {
  // 如果输入为空,不执行搜索
  if (!keyword) {
    toast.info('搜索内容为空')
    return
  }

  searchResults.value = []

  showHistory.value = false

  // 添加到搜索历史
  baseStore.addSearch(keyword)

  try {
    loading.value = true
    const resp = await videoApi.searchAll(keyword)
    searchResults.value = resp.data.data.list
    if (searchResults.value.length == 0) {
      failed.value = true
      setTimeout(() => {
        toast.warning('没有找到相关资源...')
      }, 2500)
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
  location.reload()
}

// 处理热门视频点击
const handleSearchHot = async (title: string) => {
  toast.info('搜索中...', { duration: 2000 })
  await handleSearch(title)
}
</script>

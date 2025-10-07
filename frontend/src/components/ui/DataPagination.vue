<template>
  <div v-if="total > pageSize" class="w-max">
    <div
      class="flex [&_.btn]:btn-sm sm:[&_.btn]:btn items-center justify-center sm:space-x-1 bg-card border border-border rounded-xl p-1.5 sm:p-2 shadow-sm w-full"
    >
      <!-- 上一页 -->
      <button
        @click="currentPage = Math.max(1, currentPage - 1)"
        :disabled="currentPage === 1"
        class="btn rounded-lg text-muted-foreground hover:text-foreground hover:bg-accent disabled:opacity-30 disabled:cursor-not-allowed transition-all duration-200 flex-shrink-0"
        :class="{ 'hover:bg-transparent': currentPage === 1 }"
        aria-label="上一页"
      >
        <svg class="size-3.5 sm:size-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M15 19l-7-7 7-7"
          />
        </svg>
      </button>

      <!-- 页码按钮 -->
      <div
        class="flex items-center space-x-0.5 sm:space-x-1 min-w-0 overflow-x-auto scrollbar-hide"
      >
        <template v-for="page in visiblePages" :key="page">
          <button
            v-if="typeof page === 'number'"
            @click="currentPage = page"
            class="btn rounded-lg text-xs sm:text-sm font-medium transition-all duration-200 flex-shrink-0 min-w-[28px] sm:min-w-[32px]"
            :class="
              currentPage === page
                ? 'bg-primary text-primary-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent'
            "
          >
            {{ page }}
          </button>
          <span
            v-else
            class="btn rounded-lg text-xs sm:text-sm font-medium transition-all duration-200 flex-shrink-0 min-w-[28px] sm:min-w-[32px] cursor-default text-muted-foreground"
          >
            {{ page }}
          </span>
        </template>
      </div>

      <!-- 下一页 -->
      <button
        @click="currentPage = Math.min(totalPages, currentPage + 1)"
        :disabled="currentPage === totalPages"
        class="btn rounded-lg btn-ghost flex-shrink-0 text-muted-foreground hover:text-foreground hover:bg-accent disabled:opacity-30 disabled:cursor-not-allowed transition-all duration-200"
        :class="{ 'hover:bg-transparent': currentPage === totalPages }"
        aria-label="下一页"
      >
        <svg class="size-3.5 sm:size-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const currentPage = defineModel<number>('currentPage', { default: 1 })
const pageSize = defineModel<number>('pageSize', { default: 12 })
const total = defineModel<number>('total', { required: true })

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

// 响应式判断是否为移动端
const isMobile = computed(() => {
  if (typeof window === 'undefined') return false
  return window.innerWidth < 640
})

const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const totalPageCount = totalPages.value
  const current = currentPage.value

  // 移动端显示更少的页码
  const maxVisiblePages = isMobile.value ? 5 : 8

  // 总页数较少，直接显示所有页
  if (totalPageCount <= maxVisiblePages) {
    for (let i = 1; i <= totalPageCount; i++) pages.push(i)
    return pages
  }

  // 移动端逻辑：只显示 首页 + 当前页附近 + 尾页
  if (isMobile.value) {
    pages.push(1) // 首页

    if (current > 3) {
      pages.push('...')
    }

    // 当前页附近只显示 1 个
    const start = Math.max(2, current - 1)
    const end = Math.min(totalPageCount - 1, current + 1)

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    if (current < totalPageCount - 2) {
      pages.push('...')
    }

    pages.push(totalPageCount) // 尾页
  } else {
    // 桌面端逻辑：显示更多页码
    pages.push(1)

    if (current > 3) {
      pages.push('...')
    }

    const start = Math.max(2, current - 2)
    const end = Math.min(totalPageCount - 1, current + 2)

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    if (current < totalPageCount - 3) {
      pages.push('...')
    }

    pages.push(totalPageCount)
  }

  return pages
})

const resetPage = () => {
  currentPage.value = 1
}

const setPage = (page: number) => {
  currentPage.value = page
}

defineExpose({
  resetPage,
  setPage,
})
</script>

<style scoped>
/* 隐藏滚动条但保持滚动功能 */
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
</style>

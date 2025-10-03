<template>
  <div v-if="total > pageSize" class="w-max">
    <div class="flex items-center space-x-2 bg-card border border-border rounded-xl p-2 shadow-sm">
      <!-- 上一页 -->
      <button
        @click="currentPage = Math.max(1, currentPage - 1)"
        :disabled="currentPage === 1"
        class="px-3 py-2 rounded-lg text-muted-foreground hover:text-foreground hover:bg-accent disabled:opacity-30 disabled:cursor-not-allowed transition-all duration-200"
        :class="{ 'hover:bg-transparent': currentPage === 1 }"
        aria-label="上一页"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M15 19l-7-7 7-7"
          />
        </svg>
      </button>

      <!-- 页码按钮 -->
      <div class="flex items-center space-x-1">
        <template v-for="page in visiblePages" :key="page">
          <button
            v-if="typeof page === 'number'"
            @click="currentPage = page"
            class="min-w-[2.5rem] px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200"
            :class="
              currentPage === page
                ? 'bg-primary text-primary-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent'
            "
          >
            {{ page }}
          </button>
          <button
            v-else
            class="min-w-[2.5rem] px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200"
          >
            {{ page }}
          </button>
        </template>
      </div>

      <!-- 下一页 -->
      <button
        @click="currentPage = Math.min(totalPages, currentPage + 1)"
        :disabled="currentPage === totalPages"
        class="px-3 py-2 rounded-lg text-muted-foreground hover:text-foreground hover:bg-accent disabled:opacity-30 disabled:cursor-not-allowed transition-all duration-200"
        :class="{ 'hover:bg-transparent': currentPage === totalPages }"
        aria-label="下一页"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

const visiblePages = computed(() => {
  const pages: (number | string)[] = [] // 用 string 表示省略号
  const totalPageCount = totalPages.value
  const current = currentPage.value

  // 总页数 <= 6，直接显示所有页
  if (totalPageCount <= 8) {
    for (let i = 1; i <= totalPageCount; i++) pages.push(i)
  } else {
    pages.push(1) // 首页

    if (current > 3) {
      pages.push('...') // 前面省略
    }

    // 中间页，保证当前页左右各 1 个
    const start = Math.max(2, current - 2)
    const end = Math.min(totalPageCount - 1, current + 2)

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    if (current < totalPageCount - 3) {
      pages.push('...') // 后面省略
    }

    pages.push(totalPageCount) // 尾页
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

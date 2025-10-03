<template>
  <div class="relative flex items-center gap-2 px-4 py-2 rounded-lg">
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
      v-model.lazy="searchText"
      @keydown.enter="handleSearch"
      @focus="showHistory = true"
      @blur="handleBlur"
      @compositionstart="isComposing = true"
      @compositionend="isComposing = false"
      class="flex-1 bg-transparent border-none outline-none text-sm placeholder:text-muted-foreground"
      :placeholder="placeholder"
    />

    <svg
      v-if="searchText || !isComposing"
      @click="handleClear"
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

    <button class="btn btn-primary btn-sm" @click="handleSearch" :disabled="loading || isComposing">
      <span v-if="!loading">搜索</span>
      <span v-else class="flex items-center gap-2">
        <span class="spinner spinner-sm"></span>
        搜索中
      </span>
    </button>

    <!-- 搜索历史 -->
    <div
      v-if="enableHistory && showHistory && historyList.length > 0 && !searchText"
      class="absolute top-full left-0 right-0 mt-2 bg-background border border-border rounded-lg shadow-lg z-50"
    >
      <div class="flex items-center justify-between px-4 py-3 border-b border-border">
        <span class="text-sm font-medium">最近搜索</span>
        <button
          @click.stop="handleClearHistory"
          class="text-xs text-muted-foreground hover:text-foreground"
        >
          清空
        </button>
      </div>
      <div class="py-2">
        <div
          v-for="keyword in historyList"
          :key="keyword"
          class="flex items-center justify-between px-4 py-2 hover:bg-muted/50 cursor-pointer group"
          @mousedown.prevent="handleHistoryClick(keyword)"
        >
          <div class="flex items-center gap-3 flex-1">
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
            <span class="text-sm">{{ keyword }}</span>
          </div>
          <button
            @click.stop="handleRemoveHistory(keyword)"
            class="opacity-0 group-hover:opacity-100"
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
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps({
  placeholder: {
    type: String,
    default: '搜索...',
  },

  loading: {
    type: Boolean,
    default: false,
  },
  enableHistory: {
    type: Boolean,
    default: false,
  },
})

const historyList = defineModel<string[]>('history', { default: [] })
const emit = defineEmits(['search', 'clear'])

const searchText = ref('')
const showHistory = ref(false)
const isComposing = ref(false)

const handleClear = () => {
  searchText.value = ''
  emit('clear')
}

const handleSearch = () => {
  if (isComposing.value || !searchText.value.trim()) return
  emit('search', searchText.value.trim())
}

const handleBlur = () => {
  emit('search', searchText.value.trim())
  showHistory.value = false
}

const handleHistoryClick = (keyword: string) => {
  searchText.value = keyword
  showHistory.value = false
  emit('search', keyword)
}

// 删除单条历史记录
const handleRemoveHistory = (keyword: string) => {
  const index = historyList.value.indexOf(keyword)
  if (index !== -1) historyList.value.splice(index, 1)
}

// 清空所有历史记录
const handleClearHistory = () => {
  historyList.value = []
  showHistory.value = false
}
</script>

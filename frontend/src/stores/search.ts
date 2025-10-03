import { useStorage } from '@vueuse/core'
import { defineStore } from 'pinia'

export const useBaseStore = defineStore('base', () => {
  const recentSearchList = useStorage<string[]>('recently-search', [])

  /** 增加搜索记录，最新的放前面，最多保留 5 个 */
  function addSearch(keyword: string) {
    if (!keyword) return

    const index = recentSearchList.value.indexOf(keyword)
    if (index !== -1) recentSearchList.value.splice(index, 1)

    recentSearchList.value.unshift(keyword)

    // 保留最多 5 个
    if (recentSearchList.value.length > 5) {
      recentSearchList.value = recentSearchList.value.slice(0, 5)
    }
  }

  return {
    recentSearchList,
    addSearch,
  }
})

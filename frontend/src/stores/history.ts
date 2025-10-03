import { useStorage } from '@vueuse/core'

// 播放记录管理
export const useHistoryStore = defineStore('history', () => {
  const playHistory = useStorage<PlayHistory[]>('video_play_history', [])

  // 保存播放记录
  const savePlayHistory = (history: PlayHistory) => {
    try {
      // 查找是否已存在该视频的记录
      const existingIndex = playHistory.value.findIndex(
        (h) =>
          h.vod_id === history.vod_id &&
          h.episode_index === history.episode_index &&
          h.sourceKey === history.sourceKey,
      )

      if (existingIndex !== -1) {
        // 更新现有记录
        playHistory.value[existingIndex] = history
      } else {
        // 添加新记录
        playHistory.value.unshift(history)
      }

      // 只保留最近100条记录
      if (playHistory.value.length > 100) {
        playHistory.value.splice(100)
      }

      console.log('播放记录已保存:', history)
    } catch (error) {
      console.error('保存播放记录失败:', error)
    }
  }

  return {
    playHistory,
    savePlayHistory,
  }
})

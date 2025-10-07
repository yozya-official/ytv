<template>
  <video
    class="w-full aspect-video plyr-video"
    ref="videoRef"
    playsinline
    preload="metadata"
  ></video>
</template>

<script setup lang="ts">
import Hls from 'hls.js'
import Plyr from 'plyr'
import 'plyr/dist/plyr.css'

const store = useHistoryStore()
const { savePlayHistory } = store

const { playHistory } = storeToRefs(store)

const videoRef = ref<HTMLVideoElement | null>(null)
let hls: Hls | null = null
const plyr = ref<Plyr | null>(null)

const props = defineProps({
  player_url: { type: String, required: true },
  vod_id: { type: Number, required: true }, // 视频ID
  episode_index: { type: Number, required: true }, // 剧集序号
  episode_title: { type: String, required: true }, // 剧集标题
  vod_name: { type: String, required: true }, // 标题
  source_key: { type: String, required: true }, // 播放源标识
})

const emits = defineEmits(['play-end'])

let progressSaveTimer: number | null = null

// 获取当前视频的播放记录
const getCurrentPlayHistory = (): PlayHistory | null => {
  return (
    playHistory.value.find(
      (h) =>
        h.vod_id === props.vod_id &&
        h.episode_index === props.episode_index &&
        h.sourceKey === props.source_key,
    ) || null
  )
}

// 更新播放进度
const updateProgress = () => {
  const video = videoRef.value
  if (!video) return

  const currentTime = video.currentTime
  const duration = video.duration

  if (!duration || isNaN(duration)) return
  if (currentTime > 5 && currentTime < duration - 10 && currentTime % 5 < 0.25) {
    const history: PlayHistory = {
      vod_id: props.vod_id,
      episode_index: props.episode_index,
      name: props.vod_name,
      episode_title: props.episode_title,
      sourceKey: props.source_key,
      lastPlayTime: Date.now(),
      progress: Math.floor(currentTime),
      duration: Math.floor(duration),
    }
    savePlayHistory(history)
  }
}

// 恢复播放进度
const restoreProgress = () => {
  const video = videoRef.value
  if (!video) return

  const history = getCurrentPlayHistory()

  if (history && history.progress > 5) {
    console.log(`恢复播放进度: ${history.progress}秒`)

    // 等待视频可以播放后再跳转
    const seekToSaved = () => {
      video.currentTime = history.progress
      video.removeEventListener('canplay', seekToSaved)
    }

    if (video.readyState >= 2) {
      video.currentTime = history.progress
    } else {
      video.addEventListener('canplay', seekToSaved, { once: true })
    }
  }
}

// 视频播放结束处理
const handleVideoEnded = () => {
  emits('play-end')
}

// 清理 HLS 实例
const cleanupHlsOnly = () => {
  const video = videoRef.value

  if (hls) {
    console.log('销毁 HLS 实例')
    hls.destroy()
    hls = null
  }

  if (video) {
    // 移除所有事件监听器
    video.removeEventListener('ended', handleVideoEnded)
    video.removeEventListener('timeupdate', updateProgress)

    video.src = ''
    video.load()
  }

  // 清除进度保存定时器
  if (progressSaveTimer) {
    clearTimeout(progressSaveTimer)
    progressSaveTimer = null
  }
}

// 初始化播放器
const initPlayer = async (url: string) => {
  const video = videoRef.value
  if (!video) return

  console.log('开始加载视频:', url)

  cleanupHlsOnly()
  await nextTick()

  // 首次初始化 Plyr
  if (!plyr.value) {
    console.log('初始化 Plyr 实例')
    plyr.value = new Plyr(video, {
      autoplay: true,
      controls: [
        'play-large',
        'play',
        'progress',
        'current-time',
        'duration',
        'mute',
        'volume',
        'settings',
        'pip',
        'fullscreen',
      ],
      settings: ['quality', 'speed'],
    })
  }

  // 添加事件监听器
  video.addEventListener('ended', handleVideoEnded)
  video.addEventListener('timeupdate', updateProgress) // 监听播放进度

  // 元数据加载完成后恢复进度
  const handleLoadedMetadata = () => {
    console.log('视频元数据加载完成')
    plyr.value?.restart()
    restoreProgress() // 恢复播放进度
    video.removeEventListener('loadedmetadata', handleLoadedMetadata)
  }
  video.addEventListener('loadedmetadata', handleLoadedMetadata)

  // 加载 HLS 源
  if (video.canPlayType('application/vnd.apple.mpegurl')) {
    console.log('使用原生 HLS 播放')
    video.src = url
  } else if (Hls.isSupported()) {
    console.log('使用 hls.js 播放')

    hls = new Hls({ debug: false, enableWorker: true, maxBufferHole: 0.5 })
    hls.attachMedia(video)
    hls.loadSource(url)

    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      console.log('HLS manifest 解析完成')
      plyr.value?.play()
    })

    hls.on(Hls.Events.ERROR, (event, data) => {
      if (data.fatal) {
        console.error('HLS 致命错误:', data)
        switch (data.type) {
          case Hls.ErrorTypes.NETWORK_ERROR:
          case Hls.ErrorTypes.MEDIA_ERROR:
            hls?.recoverMediaError()
            break
          default:
            cleanupHlsOnly()
            break
        }
      }
    })
  } else {
    console.error('当前浏览器不支持 HLS 播放')
  }
}

// 监听 URL 变化
watch(
  () => props.player_url,
  () => {
    initPlayer(props.player_url)
  },
)

onMounted(() => {
  initPlayer(props.player_url)
})

// 组件卸载清理
onUnmounted(() => {
  console.log('组件卸载，清理资源')
  cleanupHlsOnly()

  if (plyr.value) {
    plyr.value.destroy()
    plyr.value = null
  }
})
</script>

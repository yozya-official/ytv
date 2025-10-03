export interface PlayHistory {
  vod_id: number
  episode_index: number
  name: string
  episode_title: string
  sourceKey: string
  lastPlayTime: number
  progress: number // 播放进度（秒）
  duration: number // 视频总时长（秒）
}

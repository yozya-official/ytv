/**
 * VodItem 爬取结果一个视频条目的所有属性。
 */
export interface VodItem {
  source_name: string
  source_key: string
  episodes: Episode[]
  // 基本信息
  vod_id: number // 视频ID
  type_id: number // 类型ID
  type_id_1: number // 一级类型ID
  type_name: string // 类型名称
  vod_name: string // 视频名称
  vod_sub: string // 副标题
  vod_en: string // 视频英文名/拼音
  vod_year: string // 年份
  vod_area: string // 地区
  vod_lang: string // 语言
  vod_remarks: string // 备注/集数状态（如：正片, HD中字, 更新第11集）
  vod_status: number // 状态（1表示正常）
  vod_isend: number // 是否完结（1表示完结, 0表示未完结）
  vod_total: number // 总集数
  vod_serial: string // 当前连载集数

  // 描述和分类
  vod_class: string // 分类（如：悬疑, 剧情）
  vod_tag: string // 标签
  vod_blurb: string // 简介摘要
  vod_content: string // 详细内容/剧情简介（可能包含HTML标签）

  // 人员信息
  vod_actor: string // 演员
  vod_director: string // 导演
  vod_writer: string // 编剧
  vod_author: string // 原作者

  // 评分和热度
  vod_score: string // 评分
  vod_douban_id: number // 豆瓣ID
  vod_douban_score: string // 豆瓣评分
  vod_hits: number // 总点击量
  vod_hits_day: number // 日点击量
  vod_hits_week: number // 周点击量
  vod_hits_month: number // 月点击量
  vod_up: number // 点赞数
  vod_down: number // 点踩数

  // 时间信息
  vod_pubdate: string // 发布日期
  vod_time: string // 数据更新时间
  vod_duration: string // 时长/单集时长

  // 媒体资源
  vod_pic: string // 封面图片URL
  vod_pic_thumb: string // 缩略图URL
  vod_play_from: string // 播放源标识
  vod_play_url: string // 播放地址（通常是"名称$URL#名称$URL"格式）

  // 其他/冗余信息
  vod_letter: string // 首字母
  vod_color: string // 颜色标记
  vod_version: string // 版本（如：高清版）
  vod_state: string // 状态（如：正片）
}

export interface Episode {
  episode_index: number
  episode_title: string
  url: string
}

// 统一信息
export interface APIResponse<T, K> {
  code: number
  msg: string
  data: T
  extra: K
}

interface VodList {
  list: VodItem[]
}

// 热搜参数额外信息
export interface SearchHotExtra {
  page_limit: '16'
  page_start: '0'
  sort: 'recommend'
  tag: '热门'
  type: 'movie'
}

// 搜索关键词参数额外信息
export interface SearchKeywordExtra {
  failed_count: number
  keyword: string
  page: number
  success_count: number
  total_sources: number
}

// 搜索单个参数额外信息
export interface SearchDetailExtra {
  source_key: string
  vod_id: number
}

// 搜索关键词结果
export type SearchAllResult = APIResponse<VodList, SearchKeywordExtra>

// 搜索单个结果
export type SearchDetailResult = APIResponse<VodItem, SearchDetailExtra>

export type HotVideosResult = APIResponse<HotMovies, SearchHotExtra>

//  资源api
export interface VideoSource {
  api: string
  name: string
  detail?: string
  adult?: boolean
}

// 热搜
export interface HotMovies {
  list: HotMovie[]
}

export interface HotMovie {
  episodes_info: string
  rate: string
  cover_x: number
  title: string
  url: string
  playable: boolean
  cover: string
  id: string
  cover_y: number
  is_new: boolean

  cover_base64: string // 需要请求 hot-covers
}

/* 观影记录 */
export interface ViewingRecord {
  title: string
  directVideoUrl: string
  url: string
  episodeIndex: number
  sourceName: string
  vod_id: string
  sourceCode: string
  showIdentifier: string
  timestamp: number
  playbackPosition: number
  duration: number
  episodes: string[]
}

import axios from 'axios'

// 创建 axios 实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000, // 并行请求多个源，增加超时时间
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    console.log('Request:', config.method?.toUpperCase(), config.url)
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    console.error('API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  },
)

// Video API 接口
export const videoApi = {
  /**
   * 搜索所有视频源
   * @param keyword 搜索关键词
   * @param page 页码，默认 1
   * @param includeAdult 是否包含成人内容，默认 false
   */
  searchAll: (keyword: string, page: number = 1, includeAdult: boolean = false) => {
    return api.get<SearchAllResult>('/search', {
      params: {
        wd: keyword,
        pg: page,
        adult: includeAdult,
      },
    })
  },

  /**
   * todo 搜索单个视频源
   * @param source 视频源 key (如: ffzy, zy360)
   * @param keyword 搜索关键词
   * @param page 页码，默认 1
   */
  searchSingle: (source: string, keyword: string, page: number = 1) => {
    return api.get<SearchAllResult>(`/search/${source}`, {
      params: {
        wd: keyword,
        pg: page,
      },
    })
  },

  /**
   * 获取热门影片
   * @param page 页码，默认 1
   * @param options.type 类型 movie=电影 tv=电视剧，默认 movie
   * @param options.tag 标签，默认 热门
   */
  getHotVideos: (page_start: number = 1, page_limit: number = 12, type = 'movie', tag = '热门') => {
    return api.get<HotVideosResult>('/hot', {
      params: {
        page_start,
        page_limit,
        type,
        tag,
      },
    })
  },

  /**
   * 根据id搜索vod
   */
  searchById: (sourceKey: string, vodId: number) => {
    return api.get<SearchDetailResult>('/vod', {
      params: {
        sourceKey,
        vodId,
      },
    })
  },

  // /**
  //  * 获取所有可用的视频源列表
  //  * @param includeAdult 是否包含成人内容，默认 false
  //  */
  // getSources: (includeAdult: boolean = false) => {
  //   return api.get<SourcesResponse>('/sources', {
  //     params: {
  //       adult: includeAdult,
  //     },
  //   })
  // },
}

export default api

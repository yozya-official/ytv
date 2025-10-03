package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"tv/cache"
	"tv/conf"
	"tv/models"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

var videoAPI = NewVideoAPI()

type VideoAPI struct {
	client *resty.Client
}

func NewVideoAPI() *VideoAPI {
	c := resty.New().
		SetTimeout(3*time.Second).
		SetHeader("User-Agent", "Mozilla/5.0 (compatible; VideoAPI/1.0)")
	return &VideoAPI{client: c}
}

// 原始API响应（内部使用）
type videoAPIResponse struct {
	Code      int              `json:"code"`
	Msg       string           `json:"msg"`
	Page      json.Number      `json:"page"`
	Pagecount int              `json:"pagecount"`
	Limit     json.Number      `json:"limit"`
	Total     json.Number      `json:"total"`
	List      []models.VodItem `json:"list"`
}

// 单个源的抓取结果（内部使用）
type sourceResult struct {
	SourceKey  string
	SourceName string
	Items      []models.VodItem
	Error      error
	Duration   int64
}

// 从单个源获取数据
func (api *VideoAPI) fetchFromSource(sourceKey string, source models.VideoSource, params map[string]string) sourceResult {
	start := time.Now()
	result := sourceResult{SourceKey: sourceKey, SourceName: source.Name}

	log.Debug().
		Str("source", sourceKey).
		Str("api", source.API).
		Fields(params).
		Msg("开始请求视频源")

	resp, err := api.client.R().SetQueryParams(params).Get(source.API)
	if err != nil {
		result.Error = fmt.Errorf("请求失败: %v", err)
		result.Duration = time.Since(start).Milliseconds()

		log.Error().
			Str("source", sourceKey).
			Err(err).
			Int64("duration_ms", result.Duration).
			Msg("请求视频源失败")

		return result
	}

	var apiResp videoAPIResponse
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		result.Error = fmt.Errorf("解析 JSON 失败: %v", err)
		result.Duration = time.Since(start).Milliseconds()

		log.Error().
			Str("source", sourceKey).
			Err(err).
			Int64("duration_ms", result.Duration).
			Msg("解析视频源响应失败")

		return result
	}

	result.Items = make([]models.VodItem, len(apiResp.List))
	for i, item := range apiResp.List {
		item.SourceKey = sourceKey
		item.SourceName = source.Name
		item.Episodes = parseVodPlayURL(item.VodPlayURL)
		result.Items[i] = item
	}
	result.Duration = time.Since(start).Milliseconds()

	log.Debug().
		Str("source", sourceKey).
		Int("items", len(result.Items)).
		Int64("duration_ms", result.Duration).
		Msg("视频源请求完成")

	return result
}

// 并行抓取多个源
func (api *VideoAPI) fetchParallel(sources map[string]models.VideoSource, fetcher func(string, models.VideoSource) sourceResult) []sourceResult {
	var wg sync.WaitGroup
	resultChan := make(chan sourceResult, len(sources))

	for key, source := range sources {
		wg.Add(1)
		go func(k string, s models.VideoSource) {
			defer wg.Done()
			resultChan <- fetcher(k, s)
		}(key, source)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []sourceResult
	for r := range resultChan {
		results = append(results, r)
	}
	return results
}

// 搜索关键词
func (api *VideoAPI) SearchByKeyword(keyword, page string, includeAdult bool) (any, any, error) {
	start := time.Now()
	sources := conf.Cfg.GetActiveVideoSources()

	log.Info().
		Str("keyword", keyword).
		Str("page", page).
		Bool("adult", includeAdult).
		Int("sources", len(sources)).
		Msg("开始关键词搜索")

	results := api.fetchParallel(sources, func(key string, s models.VideoSource) sourceResult {
		params := map[string]string{"ac": "videolist", "wd": keyword, "pg": page}
		return api.fetchFromSource(key, s, params)
	})

	all := make([]models.VodItem, 0)
	successCount := 0
	failedCount := 0
	for _, r := range results {
		if r.Error == nil {
			successCount++
			all = append(all, r.Items...)
		} else {
			failedCount++
		}
	}

	duration := time.Since(start).Milliseconds()
	log.Info().
		Str("keyword", keyword).
		Str("page", page).
		Int("success", successCount).
		Int("failed", failedCount).
		Int("total_items", len(all)).
		Int64("duration_ms", duration).
		Msg("关键词搜索完成")

	data := gin.H{"list": all, "total": len(all)}
	extra := gin.H{
		"keyword":       keyword,
		"page":          page,
		"success_count": successCount,
		"failed_count":  failedCount,
		"total_sources": len(sources),
	}
	return data, extra, nil
}

// 根据ID搜索
func (api *VideoAPI) SearchByID(sourceKey string, vodID int) (any, any, error) {
	start := time.Now()
	source, ok := conf.Cfg.GetVideoSource(sourceKey)
	if !ok {
		log.Warn().
			Str("source_key", sourceKey).
			Int("vod_id", vodID).
			Msg("视频源不存在")
		return nil, gin.H{"source_key": sourceKey, "vod_id": vodID}, fmt.Errorf("视频源不存在")
	}

	params := map[string]string{"ac": "videolist", "ids": strconv.Itoa(vodID)}
	result := api.fetchFromSource(sourceKey, source, params)
	duration := time.Since(start).Milliseconds()

	if result.Error != nil {
		log.Error().
			Str("source_key", sourceKey).
			Int("vod_id", vodID).
			Err(result.Error).
			Int64("duration_ms", duration).
			Msg("ID 搜索失败")
		return nil, gin.H{"source_key": sourceKey, "vod_id": vodID}, result.Error
	}

	if len(result.Items) == 0 {
		log.Warn().
			Str("source_key", sourceKey).
			Int("vod_id", vodID).
			Int64("duration_ms", duration).
			Msg("ID 搜索无结果")
		return nil, nil, fmt.Errorf("没有查到相关信息")
	}

	log.Info().
		Str("source_key", sourceKey).
		Int("vod_id", vodID).
		Int("items", len(result.Items)).
		Int64("duration_ms", duration).
		Msg("ID 搜索成功")

	data := result.Items[0]
	extra := gin.H{"source_key": sourceKey, "vod_id": vodID}
	return data, extra, nil
}

// ============ Handler ============

// ============ Handler ============

func SearchVideoAPI(c *gin.Context) {
	// --- 【补充点 1：函数入口日志】 ---
	log.Info().Msg("处理视频关键词搜索请求")

	keyword := c.Query("wd")
	if keyword == "" {
		// --- 【补充点 2：参数校验失败日志】 ---
		log.Warn().Msg("请求参数 'wd' (关键词) 不能为空")
		Error(c, 400, "搜索关键词不能为空", nil)
		return
	}
	page := c.DefaultQuery("pg", "1")
	includeAdult := c.DefaultQuery("adult", "false") == "true"

	log.Debug().
		Str("keyword", keyword).
		Str("page", page).
		Bool("adult", includeAdult).
		Msg("请求参数解析成功")

	// 获取缓存
	cacheKey := cache.SearchParams{
		Keyword:      keyword,
		Page:         page,
		IncludeAdult: includeAdult,
	}
	cacher := cache.GetCacher()
	res, ok := cacher.GetKeyword(cacheKey)
	if ok {
		log.Info().
			Str("keyword", keyword).
			Str("page", page).
			Msg("关键词搜索请求命中缓存")
		Success(c, res.Data, res.Extra)
		return
	}
	log.Debug().
		Str("keyword", keyword).
		Str("page", page).
		Msg("关键词搜索请求未命中缓存，将调用后端服务")

	data, extra, err := videoAPI.SearchByKeyword(keyword, page, includeAdult)
	if err != nil {
		log.Error().
			Str("keyword", keyword).
			Str("page", page).
			Err(err).
			Msg("调用 SearchByKeyword 失败")
		Error(c, 500, err.Error(), extra)
		return
	}

	// 缓存结果的日志，可以放在 cacher.SetKeyword 之后或之前
	cacher.SetKeyword(cacheKey, models.APIResponse{
		Code:    0,
		Message: "",
		Data:    data,
		Extra:   extra,
	})
	log.Debug().
		Str("keyword", keyword).
		Str("page", page).
		Msg("关键词搜索结果已存入缓存")

	Success(c, data, extra)
	log.Info().
		Str("keyword", keyword).
		Str("page", page).
		Msg("关键词搜索请求处理完成并成功返回")
}

func SearchVideoById(c *gin.Context) {
	log.Info().Msg("处理视频 ID 搜索请求")

	sourceKey := c.Query("sourceKey")
	if sourceKey == "" {
		log.Warn().Msg("请求参数 'sourceKey' 不能为空")
		Error(c, 400, "sourceKey不能为空", nil)
		return
	}
	vodIDStr := c.Query("vodId")
	if vodIDStr == "" {
		log.Warn().
			Str("source_key", sourceKey).
			Msg("请求参数 'vodId' 不能为空")
		Error(c, 400, "vodId不能为空", nil)
		return
	}
	vodID, err := strconv.Atoi(vodIDStr)
	if err != nil {
		log.Warn().
			Str("source_key", sourceKey).
			Str("vod_id_str", vodIDStr).
			Err(err).
			Msg("请求参数 'vodId' 格式错误")
		Error(c, 400, "vodId 格式错误", vodIDStr)
		return
	}

	log.Debug().
		Str("source_key", sourceKey).
		Int("vod_id", vodID).
		Msg("请求参数解析成功")

	// 获取缓存
	cacheKey := cache.IDParams{
		SourceKey: sourceKey,
		VodID:     vodID,
	}
	cacher := cache.GetCacher()
	res, ok := cacher.GetByID(cacheKey)
	if ok {
		log.Info().
			Str("source_key", sourceKey).
			Int("vod_id", vodID).
			Msg("ID 搜索请求命中缓存")
		Success(c, res.Data, res.Extra)
		return
	}
	log.Debug().
		Str("source_key", sourceKey).
		Int("vod_id", vodID).
		Msg("ID 搜索请求未命中缓存，将调用后端服务")

	data, extra, err := videoAPI.SearchByID(sourceKey, vodID)
	if err != nil {
		log.Error().
			Str("source_key", sourceKey).
			Int("vod_id", vodID).
			Err(err).
			Msg("调用 SearchByID 失败")
		Error(c, 500, err.Error(), extra)
		return
	}

	// 缓存结果的日志
	cacher.SetByID(cacheKey, models.APIResponse{
		Code:    0,
		Message: "",
		Data:    data,
		Extra:   extra,
	})
	log.Debug().
		Str("source_key", sourceKey).
		Int("vod_id", vodID).
		Msg("ID 搜索结果已存入缓存")

	Success(c, data, extra)
	// --- 【补充点 8：请求成功日志】 ---
	log.Info().
		Str("source_key", sourceKey).
		Int("vod_id", vodID).
		Msg("ID 搜索请求处理完成并成功返回")
}

// ====== 工具函数 =======
// 解析播放URL
func parseVodPlayURL(playURL string) []models.Episode {
	var episodes []models.Episode
	items := strings.Split(playURL, "#")
	for i, item := range items {
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, "$", 2)
		if len(parts) == 2 {
			episodes = append(episodes, models.Episode{EpisodeIndex: i, EpisodeTitle: parts[0], URL: parts[1]})
		} else {
			episodes = append(episodes, models.Episode{EpisodeIndex: i, URL: parts[0]})
		}
	}
	return episodes
}

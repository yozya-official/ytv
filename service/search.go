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
	Limit     string           `json:"limit"`
	Total     int              `json:"total"`
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

	resp, err := api.client.R().
		SetQueryParams(params).
		Get(source.API)
	if err != nil {
		result.Error = fmt.Errorf("请求失败: %v", err)
		result.Duration = time.Since(start).Milliseconds()
		return result
	}

	var apiResp videoAPIResponse
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		result.Error = fmt.Errorf("解析 JSON 失败: %v", err)
		result.Duration = time.Since(start).Milliseconds()
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
	sources := conf.GetActiveVideoSources()
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
	source, ok := conf.GetVideoSource(sourceKey)
	if !ok {
		return nil, gin.H{"source_key": sourceKey, "vod_id": vodID}, fmt.Errorf("视频源不存在")
	}
	params := map[string]string{"ac": "videolist", "ids": strconv.Itoa(vodID)}
	result := api.fetchFromSource(sourceKey, source, params)
	if result.Error != nil {
		return nil, gin.H{"source_key": sourceKey, "vod_id": vodID}, result.Error
	}

	if len(result.Items) == 0 {
		return nil, nil, fmt.Errorf("没有查到相关信息")
	}
	data := result.Items[0]
	extra := gin.H{"source_key": sourceKey, "vod_id": vodID}
	return data, extra, nil
}

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

// ============ Handler ============

func SearchVideoAPI(c *gin.Context) {
	keyword := c.Query("wd")
	if keyword == "" {
		Error(c, 400, "搜索关键词不能为空", nil)
		return
	}
	page := c.DefaultQuery("pg", "1")
	includeAdult := c.DefaultQuery("adult", "false") == "true"

	// 获取缓存
	cacheKey := cache.SearchParams{
		Keyword:      keyword,
		Page:         page,
		IncludeAdult: includeAdult,
	}
	cacher := cache.GetCacher()
	res, ok := cacher.GetKeyword(cacheKey)
	if ok {
		Success(c, res.Data, res.Extra)
		return
	}

	data, extra, err := videoAPI.SearchByKeyword(keyword, page, includeAdult)
	if err != nil {
		Error(c, 500, err.Error(), extra)
		return
	}

	cacher.SetKeyword(cacheKey, models.APIResponse{
		Code:    0,
		Message: "",
		Data:    data,
		Extra:   extra,
	})

	Success(c, data, extra)
}

func SearchVideoById(c *gin.Context) {
	sourceKey := c.Query("sourceKey")
	if sourceKey == "" {
		Error(c, 400, "sourceKey不能为空", nil)
		return
	}
	vodIDStr := c.Query("vodId")
	if vodIDStr == "" {
		Error(c, 400, "vodId不能为空", nil)
		return
	}
	vodID, err := strconv.Atoi(vodIDStr)
	if err != nil {
		Error(c, 400, "vodId 格式错误", vodIDStr)
		return
	}

	// 获取缓存
	cacheKey := cache.IDParams{
		SourceKey: sourceKey,
		VodID:     vodID,
	}
	cacher := cache.GetCacher()
	res, ok := cacher.GetByID(cacheKey)
	if ok {
		Success(c, res.Data, res.Extra)
		return
	}
	data, extra, err := videoAPI.SearchByID(sourceKey, vodID)
	if err != nil {
		Error(c, 500, err.Error(), extra)
		return
	}

	cacher.SetByID(cacheKey, models.APIResponse{
		Code:    0,
		Message: "",
		Data:    data,
		Extra:   extra,
	})

	Success(c, data, extra)
}

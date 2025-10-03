package service

import (
	"encoding/json"
	"time"
	"tv/cache"
	"tv/models"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)

var doubanClient = NewClient(
	"https://movie.douban.com",
	5*time.Second,
	map[string]string{
		"User-Agent": "Mozilla/5.0",
	},
)

// 电影/电视剧结构体
type Movie struct {
	EpisodesInfo string `json:"episodes_info"`
	Rate         string `json:"rate"`
	CoverX       int    `json:"cover_x"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Playable     bool   `json:"playable"`
	Cover        string `json:"cover"`
	ID           string `json:"id"`
	CoverY       int    `json:"cover_y"`
	IsNew        bool   `json:"is_new"`

	CoverBase64 string `json:"cover_base64,omitempty"` // 自定义字段 用于储存图片数据
}

// 豆瓣 API 原始响应
type DoubanResponse struct {
	Subjects []Movie `json:"subjects"`
}

type DoubanRespReturn struct {
	List  []Movie `json:"list"`
	Total int     `json:"total"`
}

func HotMovies(c *gin.Context) {
	// 获取查询参数
	params := map[string]string{
		"type":       c.DefaultQuery("type", "movie"),
		"tag":        c.DefaultQuery("tag", "热门"),
		"sort":       c.DefaultQuery("sort", "recommend"),
		"page_limit": c.DefaultQuery("page_limit", "16"),
		"page_start": c.DefaultQuery("page_start", "0"),
	}
	log.Debug().Str("path", "/hots").Interface("params", params).Msg("开始处理 HotMovies 请求")

	// 构建缓存 key
	cacheKey := cache.HotParams{
		Type:      params["type"],
		Tag:       params["tag"],
		Sort:      params["sort"],
		PageLimit: params["page_limit"],
		PageStart: params["page_start"],
	}

	// 读取缓存
	if cachedData, ok := cache.GetCacher().GetHot(cacheKey); ok {
		log.Debug().Interface("cacheKey", cacheKey).Msg("命中缓存，直接返回数据")
		Success(c, cachedData.Data, cachedData.Extra)
		return
	}
	log.Debug().Interface("cacheKey", cacheKey).Msg("缓存未命中，准备请求豆瓣 API")

	// 请求豆瓣 API
	resp, err := doubanClient.resty.R().
		SetQueryParams(params).
		Get("/j/search_subjects")
	if err != nil {
		log.Error().Err(err).Str("url", "/j/search_subjects").Interface("params", params).Msg("请求豆瓣热搜失败")
		Error(c, 500, "请求豆瓣热搜失败", nil)
		return
	}
	log.Debug().Int("status_code", resp.StatusCode()).Msg("豆瓣 API 请求完成")

	// 解析响应
	var doubanResp DoubanResponse
	if err := json.Unmarshal(resp.Body(), &doubanResp); err != nil {
		log.Error().Err(err).Str("body", string(resp.Body())).Msg("解析 JSON 失败")
		Error(c, 500, "解析 JSON 失败", string(resp.Body()))
		return
	}
	log.Debug().Int("subjects_count", len(doubanResp.Subjects)).Msg("豆瓣热搜解析成功")

	// 构建响应数据
	data := DoubanRespReturn{
		Total: len(doubanResp.Subjects),
		List:  doubanResp.Subjects,
	}

	// 缓存数据
	cache.GetCacher().SetHot(cacheKey, models.APIResponse{
		Code:    0,
		Message: "",
		Data:    data,
		Extra:   params,
	})
	log.Debug().Msg("数据已缓存成功")

	// 返回响应
	Success(c, data, params)
	log.Debug().Msg("HotMovies 请求处理完成")
}

package cache

import (
	"fmt"
	"sync"
	"time"
	"tv/conf"
	"tv/models"

	"github.com/rs/zerolog/log"
)

var (
	instance *SearchCache
	once     sync.Once
)

func GetCacher() *SearchCache {
	once.Do(func() {
		instance = &SearchCache{
			items: make(map[string]cacheItem),
		}
		log.Info().Msg("搜索缓存已就绪")
	})
	return instance
}

type cacheItem struct {
	data      models.APIResponse
	expiresAt time.Time
}

// 定义缓存类型
type CacheType string

const (
	CacheTypeSearch CacheType = "search"
	CacheTypeID     CacheType = "id"
	CacheTypeHot    CacheType = "hot"
)

// 缓存过期时间配置
var cacheTTL = map[CacheType]time.Duration{
	CacheTypeSearch: conf.Cfg.Cache.Search,
	CacheTypeID:     conf.Cfg.Cache.ID,
	CacheTypeHot:    conf.Cfg.Cache.Hot,
}

type SearchCache struct {
	sync.RWMutex
	items map[string]cacheItem
}

func NewSearchCache() *SearchCache {
	return &SearchCache{
		items: make(map[string]cacheItem),
	}
}

// ============ 通用缓存操作 ============

// 获取缓存
func (c *SearchCache) get(key string) (models.APIResponse, bool) {
	c.RLock()
	defer c.RUnlock()

	if item, ok := c.items[key]; ok && time.Now().Before(item.expiresAt) {
		log.Info().Str("key", key).Msg("[Cache Hit]")
		return item.data, true
	}
	return models.APIResponse{}, false
}

// 设置缓存
func (c *SearchCache) set(key string, data models.APIResponse, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	c.items[key] = cacheItem{
		data:      data,
		expiresAt: time.Now().Add(ttl),
	}
	log.Info().Str("key", key).Dur("ttl", ttl).Msg("[Cache Set]")
}

// ============ Key 生成器 ============

// key生成器
func makeKey(cacheType CacheType, parts ...any) string {
	key := string(cacheType)
	for _, part := range parts {
		key += fmt.Sprintf("|%v", part)
	}
	return key
}

// ============ 热搜缓存 ============

type HotParams struct {
	Type      string
	Tag       string
	Sort      string
	PageLimit string
	PageStart string
}

func (c *SearchCache) GetHot(params HotParams) (models.APIResponse, bool) {
	key := makeKey(CacheTypeHot, params.Type, params.Tag, params.Sort, params.PageLimit, params.PageStart)
	return c.get(key)
}

func (c *SearchCache) SetHot(params HotParams, data models.APIResponse) {
	key := makeKey(CacheTypeHot, params.Type, params.Tag, params.Sort, params.PageLimit, params.PageStart)
	c.set(key, data, cacheTTL[CacheTypeHot])
}

// ============ 关键词搜索缓存 ============

type SearchParams struct {
	Keyword      string
	Page         string
	IncludeAdult bool
}

func (c *SearchCache) GetKeyword(params SearchParams) (models.APIResponse, bool) {
	key := makeKey(CacheTypeSearch, params.Keyword, params.Page, params.IncludeAdult)
	return c.get(key)
}

func (c *SearchCache) SetKeyword(params SearchParams, data models.APIResponse) {
	key := makeKey(CacheTypeSearch, params.Keyword, params.Page, params.IncludeAdult)
	c.set(key, data, cacheTTL[CacheTypeSearch])
}

// ============ ID搜索缓存 ============

type IDParams struct {
	SourceKey string
	VodID     int
}

func (c *SearchCache) GetByID(params IDParams) (models.APIResponse, bool) {
	key := makeKey(CacheTypeID, params.SourceKey, params.VodID)
	return c.get(key)
}

func (c *SearchCache) SetByID(params IDParams, data models.APIResponse) {
	key := makeKey(CacheTypeID, params.SourceKey, params.VodID)
	c.set(key, data, cacheTTL[CacheTypeID])
}

// ============ 清理过期缓存 ============

// 清理所有过期的缓存项
func (c *SearchCache) CleanExpired() int {
	c.Lock()
	defer c.Unlock()

	now := time.Now()
	count := 0
	for key, item := range c.items {
		if now.After(item.expiresAt) {
			delete(c.items, key)
			count++
		}
	}

	if count > 0 {
		log.Info().Int("count", count).Msg("Cleaned expired cache items")
	}
	return count
}

// 清空指定类型的缓存
func (c *SearchCache) Clear(cacheType CacheType) int {
	c.Lock()
	defer c.Unlock()

	prefix := string(cacheType) + "|"
	count := 0
	for key := range c.items {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.items, key)
			count++
		}
	}

	log.Info().Str("type", string(cacheType)).Int("count", count).Msg("Cleared cache")
	return count
}

// 返回缓存项数量
func (c *SearchCache) Size() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.items)
}

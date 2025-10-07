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
			ttl: map[CacheType]time.Duration{
				CacheTypeSearch: conf.Cfg.Cache.Search,
				CacheTypeID:     conf.Cfg.Cache.ID,
				CacheTypeHot:    conf.Cfg.Cache.Hot,
			},
		}

		log.Info().
			Dur("search", instance.ttl[CacheTypeSearch]).
			Dur("id", instance.ttl[CacheTypeID]).
			Dur("hot", instance.ttl[CacheTypeHot]).
			Msg("搜索缓存已就绪")

		// 启动定期清理协程
		go instance.startCleanup()
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

type SearchCache struct {
	sync.RWMutex
	items map[string]cacheItem
	ttl   map[CacheType]time.Duration
}

// ============ 通用缓存操作 ============

// 获取缓存
func (c *SearchCache) get(key string) (models.APIResponse, bool) {
	c.RLock()
	defer c.RUnlock()

	item, exists := c.items[key]

	if !exists {
		log.Debug().
			Str("key", key).
			Int("total_items", len(c.items)).
			Msg("缓存 Key 不存在")
		return models.APIResponse{}, false
	}

	now := time.Now()
	if now.After(item.expiresAt) {
		log.Debug().
			Str("key", key).
			Time("expired_at", item.expiresAt).
			Time("now", now).
			Dur("expired_for", now.Sub(item.expiresAt)).
			Msg("缓存已过期")
		return models.APIResponse{}, false
	}

	log.Debug().
		Str("key", key).
		Time("expires_at", item.expiresAt).
		Dur("remaining", item.expiresAt.Sub(now)).
		Msg("缓存命中")

	return item.data, true
}

// 设置缓存
func (c *SearchCache) set(key string, data models.APIResponse, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	expiresAt := time.Now().Add(ttl)
	c.items[key] = cacheItem{
		data:      data,
		expiresAt: expiresAt,
	}

	log.Debug().
		Str("key", key).
		Dur("ttl", ttl).
		Time("expires_at", expiresAt).
		Int("total_items", len(c.items)).
		Msg("缓存已设置")
}

// ============ Key 生成器 ============

// key生成器
func makeKey(cacheType CacheType, parts ...any) string {
	key := string(cacheType)
	for _, part := range parts {
		key += fmt.Sprintf("|%v", part)
	}

	log.Debug().
		Str("cache_type", string(cacheType)).
		Str("generated_key", key).
		Interface("parts", parts).
		Msg("生成缓存 Key")

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
	c.set(key, data, c.ttl[CacheTypeHot])
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
	c.set(key, data, c.ttl[CacheTypeSearch])
}

// ============ ID搜索缓存 ============

type IDParams struct {
	SourceKey string
	VodID     int
	Index     int
}

func (c *SearchCache) GetByID(params IDParams) (models.APIResponse, bool) {
	key := makeKey(CacheTypeID, params.SourceKey, params.VodID, params.Index)
	return c.get(key)
}

func (c *SearchCache) SetByID(params IDParams, data models.APIResponse) {
	key := makeKey(CacheTypeID, params.SourceKey, params.VodID, params.Index)
	c.set(key, data, c.ttl[CacheTypeID])
}

// ============ 清理过期缓存 ============

// startCleanup 启动定期清理协程
func (c *SearchCache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Info().Msg("缓存清理协程已启动")

	for range ticker.C {
		cleaned := c.CleanExpired()
		if cleaned > 0 {
			log.Info().
				Int("cleaned", cleaned).
				Int("remaining", c.Size()).
				Msg("定期清理过期缓存完成")
		}
	}
}

// CleanExpired 清理所有过期的缓存项
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
		log.Debug().
			Int("count", count).
			Int("remaining", len(c.items)).
			Msg("清理过期缓存项")
	}
	return count
}

// Clear 清空指定类型的缓存
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

	log.Info().
		Str("type", string(cacheType)).
		Int("count", count).
		Msg("清空指定类型缓存")
	return count
}

// ClearAll 清空所有缓存
func (c *SearchCache) ClearAll() {
	c.Lock()
	defer c.Unlock()

	count := len(c.items)
	c.items = make(map[string]cacheItem)

	log.Info().
		Int("count", count).
		Msg("清空所有缓存")
}

// Size 返回缓存项数量
func (c *SearchCache) Size() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.items)
}

// Stats 获取缓存统计信息
func (c *SearchCache) Stats() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()

	now := time.Now()
	expired := 0
	byType := make(map[CacheType]int)

	for key, item := range c.items {
		if now.After(item.expiresAt) {
			expired++
		}

		// 统计每种类型的数量
		if len(key) > 0 {
			parts := []byte(key)
			for i, b := range parts {
				if b == '|' {
					cacheType := CacheType(key[:i])
					byType[cacheType]++
					break
				}
			}
		}
	}

	return map[string]interface{}{
		"total_items":   len(c.items),
		"expired_items": expired,
		"valid_items":   len(c.items) - expired,
		"by_type":       byType,
		"ttl_config": map[string]interface{}{
			"search": c.ttl[CacheTypeSearch].String(),
			"id":     c.ttl[CacheTypeID].String(),
			"hot":    c.ttl[CacheTypeHot].String(),
		},
	}
}

// DumpKeys 导出所有缓存 Key（用于调试）
func (c *SearchCache) DumpKeys() []string {
	c.RLock()
	defer c.RUnlock()

	keys := make([]string, 0, len(c.items))
	now := time.Now()

	for key, item := range c.items {
		status := "valid"
		if now.After(item.expiresAt) {
			status = "expired"
		}
		keys = append(keys, fmt.Sprintf("%s [%s, expires: %s]", key, status, item.expiresAt.Format("15:04:05")))
	}

	return keys
}

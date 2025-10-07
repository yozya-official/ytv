package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"tv/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

// ==================== 公开方法 ====================

// 通过关键词搜索 Omo 视频
func (api *VideoAPI) SearchOmo(keyword string) sourceResult {
	start := time.Now()
	result := sourceResult{
		SourceKey:  "omo",
		SourceName: "Omo",
	}

	log.Debug().
		Str("source", "omo").
		Str("keyword", keyword).
		Msg("开始搜索 Omo")

	items, err := api.scrapeOmoSearch(keyword)
	result.Duration = time.Since(start).Milliseconds()

	if err != nil {
		result.Error = fmt.Errorf("搜索失败: %v", err)
		log.Error().
			Str("source", "omo").
			Err(err).
			Int64("duration_ms", result.Duration).
			Msg("Omo 搜索失败")
		return result
	}

	result.Items = items
	log.Debug().
		Str("source", "omo").
		Int("items", len(result.Items)).
		Int64("duration_ms", result.Duration).
		Msg("Omo 搜索完成")

	return result
}

// GetOmoDetail 通过 ID 和集数索引获取 Omo 视频详情（包含播放地址）
// vodID: 视频 ID
// index: 集数索引（从 0 开始）
func (api *VideoAPI) GetOmoDetail(vodID int, index int) sourceResult {
	start := time.Now()
	result := sourceResult{
		SourceKey:  "omo",
		SourceName: "Omo",
	}

	log.Debug().
		Str("source", "omo").
		Int("vod_id", vodID).
		Int("index", index).
		Msg("开始获取 Omo 详情")

	item, err := api.scrapeOmoPlayPage(vodID, index)
	result.Duration = time.Since(start).Milliseconds()

	if err != nil {
		result.Error = fmt.Errorf("获取详情失败: %v", err)
		log.Error().
			Str("source", "omo").
			Int("vod_id", vodID).
			Err(err).
			Int64("duration_ms", result.Duration).
			Msg("Omo 详情获取失败")
		return result
	}

	result.Items = []models.VodItem{item}
	log.Debug().
		Str("source", "omo").
		Int("vod_id", vodID).
		Int("episodes", len(item.Episodes)).
		Int64("duration_ms", result.Duration).
		Msg("Omo 详情获取完成")

	return result
}

// ==================== 私有爬虫方法 ====================

// scrapeOmoSearch 爬取搜索结果页
func (api *VideoAPI) scrapeOmoSearch(keyword string) ([]models.VodItem, error) {
	searchCollector := colly.NewCollector(
		colly.AllowedDomains("www.omofun.link", "omofun.link"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	var mu sync.Mutex
	vodList := make([]models.VodItem, 0)
	vodMap := make(map[int]*models.VodItem)
	var scrapeErr error

	// 1. 爬取搜索页 - 获取视频基本信息
	searchCollector.OnHTML("div.module-card-item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a.module-card-item-poster", "href")
		idStr := strings.TrimSuffix(strings.TrimPrefix(link, "/vod/detail/id/"), ".html")
		vodID, err := strconv.Atoi(idStr)
		if err != nil {
			log.Warn().Str("id_str", idStr).Msg("无法解析 VOD ID")
			return
		}

		vod := models.VodItem{
			SourceKey:  "omo",
			SourceName: "Omo",
			VodID:      vodID,
			TypeName:   strings.TrimSpace(e.ChildText("div.module-card-item-class")),
			VodName:    strings.TrimSpace(e.ChildText("div.module-card-item-title a strong")),
			VodSerial:  strings.TrimSpace(e.ChildText("div.module-item-note")),
			VodPic:     e.ChildAttr("div.module-item-pic img", "data-original"),
			Episodes:   make([]models.Episode, 0),
		}

		mu.Lock()
		vodList = append(vodList, vod)
		vodMap[vodID] = &vodList[len(vodList)-1]
		mu.Unlock()
	})

	// 错误处理
	searchCollector.OnError(func(r *colly.Response, err error) {
		scrapeErr = err
		log.Error().Err(err).Str("url", r.Request.URL.String()).Msg("Omo 搜索页错误")
	})

	// 并发限制
	searchCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*omofun.link*",
		Parallelism: 2,
		Delay:       200 * time.Millisecond,
	})

	// 超时设置
	searchCollector.SetRequestTimeout(3 * time.Second)

	// 访问搜索页
	searchURL := "https://www.omofun.link/vod/search/page/1/wd/" + keyword + ".html"
	if err := searchCollector.Visit(searchURL); err != nil {
		return nil, fmt.Errorf("访问搜索页失败: %v", err)
	}

	searchCollector.Wait()

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	// 为每个视频获取剧集列表（不包含播放地址）
	log.Debug().Int("vod_count", len(vodList)).Msg("开始获取剧集列表")

	var wg sync.WaitGroup
	for i := range vodList {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			episodes, err := api.scrapeEpisodeList(vodList[idx].VodID)
			if err != nil {
				log.Warn().
					Int("vod_id", vodList[idx].VodID).
					Err(err).
					Msg("获取剧集列表失败")
				return
			}

			mu.Lock()
			vodList[idx].Episodes = episodes
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	return vodList, nil
}

// scrapeEpisodeList 获取指定视频的剧集列表（不包含播放地址）
func (api *VideoAPI) scrapeEpisodeList(vodID int) ([]models.Episode, error) {
	detailURL := fmt.Sprintf("https://www.omofun.link/vod/detail/id/%d.html", vodID)

	resp, err := api.client.R().Get(detailURL)
	if err != nil {
		return nil, fmt.Errorf("访问详情页失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP 状态码 %d", resp.StatusCode())
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %v", err)
	}

	episodes := make([]models.Episode, 0)

	// 提取第二个 module-list 的剧集
	doc.Find(".module-list:nth-of-type(2) a.module-play-list-link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		title := strings.TrimSpace(s.Find("span").Text())

		// 提取 nid 参数作为索引（网站上的 nid 从 1 开始）
		episodeIndex := 0
		if nidParts := strings.Split(href, "/nid/"); len(nidParts) == 2 {
			idxStr := strings.TrimSuffix(nidParts[1], ".html")
			if nid, err := strconv.Atoi(idxStr); err == nil {
				episodeIndex = nid - 1 // 转换为从 0 开始的索引
			}
		}

		episode := models.Episode{
			EpisodeIndex: episodeIndex,
			EpisodeTitle: title,
			URL:          "", // 搜索时不获取播放地址
		}

		episodes = append(episodes, episode)
	})

	return episodes, nil
}

// getPlayerUrl 从播放页面提取真实播放地址
func (api *VideoAPI) getPlayerUrl(url string) (string, error) {
	resp, err := api.client.R().Get(url)
	if err != nil {
		return "", fmt.Errorf("访问播放页失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("HTTP 状态码 %d", resp.StatusCode())
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return "", fmt.Errorf("解析 HTML 失败: %v", err)
	}

	var playURL string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptText := s.Text()

		if strings.Contains(scriptText, "player_aaaa") {
			re := regexp.MustCompile(`"url"\s*:\s*"([^"]+)"`)
			if m := re.FindStringSubmatch(scriptText); len(m) > 1 {
				raw := m[1]
				playURL = strings.ReplaceAll(raw, `\/`, `/`)
				log.Debug().Str("play_url", playURL).Msg("成功提取播放地址")
			}
		}
	})

	if playURL == "" {
		return "", fmt.Errorf("找不到播放地址")
	}

	return playURL, nil
}

// scrapeOmoPlayPage 直接爬取播放页面获取视频信息和播放地址
func (api *VideoAPI) scrapeOmoPlayPage(vodID int, index int) (models.VodItem, error) {
	// 构建播放页 URL（网站的 nid 从 1 开始，所以要 +1）
	playPageURL := fmt.Sprintf("https://www.omofun.link/vod/play/id/%d/sid/8/nid/%d.html", vodID, index+1)

	log.Debug().
		Str("play_page_url", playPageURL).
		Msg("开始访问播放页")

	resp, err := api.client.R().Get(playPageURL)
	if err != nil {
		return models.VodItem{}, fmt.Errorf("访问播放页失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return models.VodItem{}, fmt.Errorf("HTTP 状态码 %d", resp.StatusCode())
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return models.VodItem{}, fmt.Errorf("解析 HTML 失败: %v", err)
	}

	vod := models.VodItem{
		SourceKey:  "omo",
		SourceName: "Omo",
		VodID:      vodID,
		Episodes:   make([]models.Episode, 0),
	}

	// 1. 从 script 标签中提取视频名称
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptText := s.Text()
		if strings.Contains(scriptText, "var vod_name=") {
			re := regexp.MustCompile(`var vod_name='([^']+)'`)
			if matches := re.FindStringSubmatch(scriptText); len(matches) > 1 {
				vod.VodName = matches[1]
				log.Debug().Str("vod_name", vod.VodName).Msg("提取到视频名称")
			}
		}
	})

	// 2. 提取剧集（第二个 module-list）
	log.Debug().Msg("开始提取剧集列表")

	doc.Find(".player-list .module-list:nth-of-type(2) a.module-play-list-link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		title := strings.TrimSpace(s.Find("span").Text())

		// 提取 nid 参数作为索引
		episodeIndex := 0
		if nidParts := strings.Split(href, "/nid/"); len(nidParts) == 2 {
			idxStr := strings.TrimSuffix(nidParts[1], ".html")
			if nid, err := strconv.Atoi(idxStr); err == nil {
				episodeIndex = nid - 1 // 转换为从 0 开始
			}
		}

		episode := models.Episode{
			EpisodeIndex: episodeIndex,
			EpisodeTitle: title,
			URL:          "https://www.omofun.link" + href,
		}

		vod.Episodes = append(vod.Episodes, episode)
	})

	log.Debug().
		Int("total_episodes", len(vod.Episodes)).
		Msg("剧集列表提取完成")

	// 验证视频信息
	if vod.VodName == "" {
		return vod, fmt.Errorf("未找到视频信息")
	}

	if len(vod.Episodes) == 0 {
		return vod, fmt.Errorf("未找到剧集列表")
	}

	// 验证请求的集数是否存在
	if index < 0 || index >= len(vod.Episodes) {
		return vod, fmt.Errorf("集数索引 %d 超出范围 (0-%d)", index, len(vod.Episodes)-1)
	}

	// 3. 获取指定集数的播放地址
	playURL, err := api.getPlayerUrl(vod.Episodes[index].URL)
	if err != nil {
		return vod, fmt.Errorf("获取播放地址失败: %v", err)
	}

	vod.Episodes[index].URL = playURL
	log.Debug().
		Int("episode_index", index).
		Str("play_url", playURL).
		Msg("成功获取播放地址")

	return vod, nil
}

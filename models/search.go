package models

type Episode struct {
	EpisodeIndex int    `json:"episode_index"`
	EpisodeTitle string `json:"episode_title"`
	URL          string `json:"url"`
}

type VodItem struct {
	SourceKey  string    `json:"source_key"`
	SourceName string    `json:"source_name"`
	Episodes   []Episode `json:"episodes"`

	// 基本信息
	VodID      int    `json:"vod_id"`
	TypeID     int    `json:"type_id"`
	TypeID1    int    `json:"type_id_1"`
	TypeName   string `json:"type_name"`
	VodName    string `json:"vod_name"`
	VodSub     string `json:"vod_sub"`
	VodEn      string `json:"vod_en"`
	VodYear    string `json:"vod_year"`
	VodArea    string `json:"vod_area"`
	VodLang    string `json:"vod_lang"`
	VodRemarks string `json:"vod_remarks"`
	VodStatus  int    `json:"vod_status"`
	VodIsEnd   int    `json:"vod_isend"`
	VodTotal   int    `json:"vod_total"`
	VodSerial  string `json:"vod_serial"`

	// 描述和分类
	VodClass   string `json:"vod_class"`
	VodTag     string `json:"vod_tag"`
	VodBlurb   string `json:"vod_blurb"`
	VodContent string `json:"vod_content"`

	// 人员信息
	VodActor    string `json:"vod_actor"`
	VodDirector string `json:"vod_director"`
	VodWriter   string `json:"vod_writer"`
	VodAuthor   string `json:"vod_author"`

	// 评分和热度
	VodScore       string `json:"vod_score"`
	VodDoubanID    int    `json:"vod_douban_id"`
	VodDoubanScore string `json:"vod_douban_score"`
	VodHits        int    `json:"vod_hits"`
	VodHitsDay     int    `json:"vod_hits_day"`
	VodHitsWeek    int    `json:"vod_hits_week"`
	VodHitsMonth   int    `json:"vod_hits_month"`
	VodUp          int    `json:"vod_up"`
	VodDown        int    `json:"vod_down"`

	// 时间信息
	VodPubdate  string `json:"vod_pubdate"`
	VodTime     string `json:"vod_time"`
	VodDuration string `json:"vod_duration"`

	// 媒体资源
	VodPic      string `json:"vod_pic"`
	VodPicThumb string `json:"vod_pic_thumb"`
	VodPlayFrom string `json:"vod_play_from"`
	VodPlayURL  string `json:"vod_play_url"`

	// 其他信息
	VodLetter  string `json:"vod_letter"`
	VodColor   string `json:"vod_color"`
	VodVersion string `json:"vod_version"`
	VodState   string `json:"vod_state"`
}

// 统一的API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Extra   interface{} `json:"extra,omitempty"`
}

// 搜索数据结构
type SearchData struct {
	List  []VodItem `json:"list"`
	Total int       `json:"total"`
}

// 搜索额外信息
type SearchExtra struct {
	Keyword      string `json:"keyword"`
	Page         string `json:"page"`
	SuccessCount int    `json:"success_count"`
	FailedCount  int    `json:"failed_count"`
	TotalSources int    `json:"total_sources"`
}

// ID查询数据结构
type DetailData struct {
	List  []VodItem `json:"list"`
	Total int       `json:"total"`
}

// ID查询额外信息
type DetailExtra struct {
	SourceKey string `json:"source_key"`
	VodID     int    `json:"vod_id"`
}

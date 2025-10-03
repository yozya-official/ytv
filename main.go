package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"tv/conf"
	"tv/middleware"
	"tv/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {

	// 配置日志
	logger := initLog()
	log.Logger = logger

	// 初始化配置
	if err := conf.InitConfig("data/config.yaml"); err != nil {
		log.Fatal().Msg("加载配置失败")
		panic(err)
	}

	log.Debug().Msg("加载配置")

	// 初始化Gin
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	r := gin.New()
	r.Use(middleware.GinLogger(logger), middleware.GinRecovery(logger))

	mode := strings.ToLower(os.Getenv("APP_MODE"))

	// 基于部署方式动态修改参数
	addr := ":" + conf.Cfg.App.Port
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		// 全局使用
		addr = "0.0.0.0" + addr
		// CORS 中间件
		r.Use(middleware.Cors())
	}

	// API 路由组
	api := r.Group("/api/" + conf.Cfg.App.APIVersion)
	api.Use((middleware.CacheMiddleware(conf.Cfg.Cache.Header)))
	{
		api.GET("/search", service.SearchVideoAPI)
		api.GET("/hot", service.HotMovies)
		api.GET("/vod", service.SearchVideoById)
	}

	// 静态资源路由
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.svg", "./frontend/dist/favicon.svg")
	r.StaticFile("/logo.png", "./frontend/dist/logo.png")

	// 处理 Vue SPA 的所有其他路由
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API 路由未找到"})
			return
		}
		c.File("./frontend/dist/index.html")
	})

	// 启动服务器
	logger.Info().Str("地址", addr).Msg("服务已启动")
	if err := r.Run(addr); err != nil {
		logger.Fatal().Err(err).Msg("服务加载失败")
	}
}

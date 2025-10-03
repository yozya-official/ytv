package main

import (
	"io"
	"net/http"
	"time"
	"tv/conf"
	"tv/middleware"
	"tv/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {

	// 加载api
	err := conf.LoadVideoSources("data/api.yaml")
	if err != nil {
		panic(err)
	}

	// 配置 zerolog
	logger := initLog()

	log.Logger = logger

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	// 设置 Gin 模式为 release（去除 debug 信息）
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(middleware.GinLogger(logger), middleware.GinRecovery(logger))

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API 路由组
	api := r.Group("/api/v1")
	api.Use(middleware.CacheMiddleware(1 * time.Hour))
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
			c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}
		c.File("./frontend/dist/index.html")
	})

	// 启动服务器
	logger.Info().Str("address", "0.0.0.0:9000").Msg("Starting server")
	if err := r.Run("0.0.0.0:9000"); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

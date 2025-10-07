package main

import (
	"io"
	"os"
	"strings"
	"tv/conf"
	"tv/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/Yuelioi/gkit/logx/zero"
	"github.com/Yuelioi/gkit/web/gin/middleware/cachecontrol"
	"github.com/Yuelioi/gkit/web/gin/middleware/log/gzero"
	"github.com/Yuelioi/gkit/web/gin/middleware/ratelimit"
	"github.com/Yuelioi/gkit/web/gin/server"
)

func main() {

	// 配置日志
	logger := zero.Default()
	log.Logger = logger

	// 初始化配置
	if err := conf.InitConfig("data/config.yaml"); err != nil {
		log.Fatal().Msg("加载配置失败")
		panic(err)
	}

	// 禁用gin log
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	mw := []gin.HandlerFunc{
		gzero.Default(logger),
		gzero.GinRecovery(logger),
		// cachecontrol.Default(),
		ratelimit.Default(),
	}

	server.Start(server.ServerConfig{
		Addr:        ":9000",
		Logger:      logger,
		Mode:        strings.ToLower(os.Getenv("APP_MODE")),
		EnableCORS:  true,
		APIPrefix:   "/api/v1",
		Middlewares: mw,
		SPAPath:     "./frontend/dist",
	}, func(api *gin.RouterGroup) {

		api.Use((cachecontrol.Default()), ratelimit.Default())
		{
			api.GET("/search", service.SearchVideoAPI)
			api.GET("/hot", service.HotMovies)
			api.GET("/vod", service.SearchVideoById)
		}
	})

}

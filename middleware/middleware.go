package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func GinRecovery(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error().
					Interface("error", err).
					Str("path", c.Request.URL.Path).
					Msg("Panic recovered")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// 使用gozerolog
func GinLogger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 继续执行后续 handler
		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// 获取 gin 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if errorMessage == "" && statusCode >= 400 {
			// 如果没有显式错误信息，但 status >= 400，可以写默认错误
			errorMessage = fmt.Sprintf("HTTP %d", statusCode)
		}

		// 拼接 query
		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		// 根据 statusCode 设置日志等级
		var logEvent *zerolog.Event
		switch {
		case statusCode >= 500:
			logEvent = logger.Error()
		case statusCode >= 400:
			logEvent = logger.Warn()
		case statusCode >= 300:
			logEvent = logger.Info()
		case statusCode >= 200:
			logEvent = logger.Debug()
		default:
			logEvent = logger.Info()
		}

		logEvent.
			Str("client_ip", clientIP).
			Str("full_path", fullPath).
			Str("method", method).
			Str("path", path).
			Str("query", raw).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Str("error", errorMessage).
			Msg("HTTP request")
	}
}

// 缓存get请求头
func CacheMiddleware(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age="+
			fmt.Sprintf("%d", int(maxAge.Seconds())))
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

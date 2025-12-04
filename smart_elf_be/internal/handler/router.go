package handler

import (
    "smart_elf_standalone/internal/auth"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

// SetupRouter 设置路由
func SetupRouter(h *Handler, feishuAuth *auth.FeishuAuth, projectWebHost string) *gin.Engine {
	// 创建Gin引擎
	router := gin.Default()

	// 添加中间件
	router.Use(loggerMiddleware())
	router.Use(corsMiddleware())

    // 健康检查
    router.GET("/health", h.HealthCheck)
    proxyHandler := NewProxyHandler(projectWebHost, feishuAuth)

	// API路由组
	api := router.Group("/api/v1")
	{
		// 飞书事件回调
		api.POST("/lark/event", h.HandleLarkEvent)

		// 配置管理
		config := api.Group("/config")
		{
			config.POST("/update", h.UpdateConfig)
			config.GET("/query", h.QueryConfig)
			config.POST("/signature", h.GetSignature)
		}
	}
	router.Any("/proxy/*path", proxyHandler.ProxyRequest)

	log.Info().Msg("路由设置完成")
	return router
}

// loggerMiddleware 日志中间件
func loggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Float64("latency", float64(param.Latency.Microseconds())/1000.0).
			Str("client_ip", param.ClientIP).
			Msg("HTTP Request")
		return ""
	})
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, x-user-key,locale")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

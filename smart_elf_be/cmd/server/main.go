package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"smart_elf_standalone/internal"
	"smart_elf_standalone/internal/auth"
	"smart_elf_standalone/internal/handler"
	"smart_elf_standalone/internal/service"
	"smart_elf_standalone/pkg/config"
	"smart_elf_standalone/pkg/database"
	"syscall"
	"time"
)

func main() {
	// 初始化日志
	initLogger() // 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("错误: 加载配置失败: %v\n", err)
    }
    // 设置全局配置单例，便于其他模块直接使用
    config.SetGlobalConfig(cfg)

	// 初始化数据库
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("错误: 初始化数据库失败: %v\n", err)
	}
	defer func() {
		if err := database.CloseDB(db); err != nil {
			log.Printf("错误: 关闭数据库连接失败: %v\n", err)
		}
	}()

	// 数据库迁移
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("错误: 数据库迁移失败: %v\n", err)
	}
	// 初始化飞书认证
	feishuAuth := auth.NewFeishuAuth(cfg.Feishu.ProjectAPIHost, cfg.Feishu.PluginID, cfg.Feishu.PluginSecret)

	// 初始化服务
	configService := service.NewConfigService(db)
	eventService := service.NewEventService(db, configService, cfg.Feishu)

	// 初始化SmartElf核心组件
	smartElf := internal.NewSmartElf(configService, eventService)

	// 初始化处理器
	h := handler.NewHandler(smartElf)

	// 设置路由
	router := handler.SetupRouter(h, feishuAuth, cfg.Feishu.ProjectWebHost)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Printf("信息: HTTP服务器启动: addr=%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("错误: HTTP服务器启动失败: %v\n", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("信息: 正在关闭服务器...")

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("错误: 服务器关闭失败: %v\n", err)
	} else {
		log.Println("信息: 服务器已关闭")
	}
}

// initLogger 初始化日志
func initLogger() {
	// 简化的日志初始化
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 从环境变量读取日志级别（简化版本）
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		log.Printf("信息: 日志级别: %s\n", level)
	}
}

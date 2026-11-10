package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"blog-api/config"
	"blog-api/repositories"
	"blog-api/router"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置 Gin 运行模式
	gin.SetMode(cfg.ServerMode)

	// 连接数据库（自动迁移 + 幂等种子数据）
	db, err := repositories.InitDB(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 注册路由
	engine := router.Setup(cfg, db)

	// 启动服务
	addr := cfg.ServerHost + ":" + cfg.ServerPort
	log.Printf("服务监听地址: http://%s", addr)
	log.Printf("数据库文件路径: %s", cfg.DBPath)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}

package main

import (
	"github.com/gin-gonic/gin"
	"monopoly/web2"
)

func main() {
	router := gin.Default()

	// 创建游戏管理器
	gameManager := web2.NewGameManager()

	// 创建Handler
	handler := web2.NewHandler(gameManager)

	// 设置路由
	handler.SetupRoutes(router)

	// 启动服务器
	router.Run(":8080")
}

package main

import (
	"github.com/gin-gonic/gin"
)

type ClashConfig map[string]interface{}

func main() {
	r := gin.Default()

	// 创建 TemplateManager 实例
	tm := NewTemplateManager()

	// 设置模板管理路由
	SetupTemplateRoutes(r, tm)

	r.GET("/api/process", processTemplate)

	// 提供静态文件服务
	r.Static("/static", "./static")

	// 提供主页
	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})

	// 原有的处理逻辑
	// r.GET("/process", processTemplate)

	r.Run(":8080")
}

func SetupTemplateProcessRoutes(r *gin.Engine, tm *TemplateProcessor) {

}

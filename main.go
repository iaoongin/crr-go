package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type ClashConfig map[string]interface{}

func main() {
	r := gin.Default()

    // 应用全局中间件
    // r.Use(AuthMiddleware())

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

// 权限校验中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization") // 获取请求头中的 Token
		if token == "" {
			token = c.Query("token")
		}
        
        // 假设我们有一个函数 validateToken 用来校验 Token 是否有效
        if token != "123456" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort() // 阻止后续的处理
            return
        }
        
        c.Next() // 继续处理请求
    }
}
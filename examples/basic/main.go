package main

import (
	"net/http"
	"time"

	startergin "github.com/CrazyLionCat/go-spring-starter-gin"
	"github.com/go-spring/spring-core/gs"
)

// userHandler 演示把业务处理器作为 go-spring Bean 管理。
type userHandler struct{}

// newUserHandler 是 userHandler 的构造函数，供 gs.Provide 注册。
func newUserHandler() *userHandler {
	return &userHandler{}
}

// Profile 返回用户详情，示例中直接使用路径参数作为用户 ID。
func (h *userHandler) Profile(c *startergin.Context) {
	c.JSON(http.StatusOK, startergin.H{
		"id":   c.Param("id"),
		"name": "demo-user",
	})
}

// Create 演示 JSON 请求体绑定和参数校验失败时的错误响应。
func (h *userHandler) Create(c *startergin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, startergin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, startergin.H{
		"id":   "1001",
		"name": body.Name,
	})
}

// routes 聚合注册路由所需的 Bean。
//
// go-spring 会自动注入 starter 创建的 Engine 和业务处理器。
type routes struct {
	Engine      *startergin.Engine `autowire:""`
	UserHandler *userHandler       `autowire:""`
}

func init() {
	// 注册业务处理器和路由聚合对象；Init 会在依赖注入完成后调用 Register。
	gs.Provide(newUserHandler)
	gs.Provide(&routes{}).Init((*routes).Register)
}

// Register 演示直接使用 Gin API 注册全局中间件、兜底路由和分组路由。
func (r *routes) Register() {
	r.Engine.Use(requestID())
	r.Engine.NoRoute(func(c *startergin.Context) {
		c.JSON(http.StatusNotFound, startergin.H{"error": "route not found"})
	})

	r.Engine.GET("/hello", func(c *startergin.Context) {
		c.String(http.StatusOK, "hello from go-spring-starter-gin")
	})

	group := r.Engine.Group("/api/users", requireHeader("X-Demo-Token", "demo-token"))
	group.GET("/:id", r.UserHandler.Profile)
	group.POST("", r.UserHandler.Create)
	group.DELETE("/:id", requireHeader("X-Demo-Role", "admin"), func(c *startergin.Context) {
		c.Status(http.StatusNoContent)
	})
}

// requestID 为每个响应写入 X-Request-ID，优先复用客户端传入的值。
func requestID() startergin.HandlerFunc {
	return func(c *startergin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = "demo-" + time.Now().UTC().Format("20060102150405")
		}
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// requireHeader 返回一个简单的 Header 校验中间件，用于演示路由级认证。
func requireHeader(name string, value string) startergin.HandlerFunc {
	return func(c *startergin.Context) {
		if c.GetHeader(name) != value {
			c.AbortWithStatusJSON(http.StatusUnauthorized, startergin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func main() {
	// 启动 go-spring 应用，starter 会在容器初始化阶段创建 Gin Engine。
	gs.Run()
}

package main

import (
	"net/http"
	"time"

	startergin "github.com/JavaLionLi/go-spring-starter-gin"
	"github.com/go-spring/spring-core/gs"
)

type userHandler struct{}

func newUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) Profile(c *startergin.Context) {
	c.JSON(http.StatusOK, startergin.H{
		"id":   c.Param("id"),
		"name": "demo-user",
	})
}

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

type routes struct {
	Engine      *startergin.Engine `autowire:""`
	UserHandler *userHandler       `autowire:""`
}

func init() {
	gs.Provide(newUserHandler)
	gs.Provide(&routes{}).Init((*routes).Register)
}

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
	gs.Run()
}

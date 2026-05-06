package main

import (
	"net/http"
	"time"

	startergin "github.com/JavaLionLi/go-spring-starter-gin"
	"github.com/JavaLionLi/go-spring-starter-gin/routekit"
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

func init() {
	gs.Provide(newUserHandler)
	gs.Provide(requestIDMiddleware)
	gs.Provide(routeHandlers)
	gs.Provide(engineConfigurer)
	gs.Provide(publicRoutes)
	gs.Provide(userRoutes)
}

func requestIDMiddleware() startergin.Middleware {
	return startergin.NewMiddleware(100, func(c *startergin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = "demo-" + time.Now().UTC().Format("20060102150405")
		}
		c.Header("X-Request-ID", requestID)
		c.Next()
	})
}

func routeHandlers() routekit.KitItem {
	return routekit.NewKitItem(100, func(kit *routekit.Kit) {
		kit.SetHandler("login", requireHeader("X-Demo-Token", "demo-token"))
		kit.SetHandler("admin", requireHeader("X-Demo-Role", "admin"))
	})
}

func engineConfigurer() startergin.EngineConfigurer {
	return startergin.NewEngineConfigurer(100, func(engine *startergin.Engine) {
		engine.NoRoute(func(c *startergin.Context) {
			c.JSON(http.StatusNotFound, startergin.H{"error": "route not found"})
		})
	})
}

func publicRoutes() routekit.Registrar {
	return routekit.NewRegistrar(100, func(engine *routekit.Engine, kit routekit.Kit) {
		engine.GET("/hello", func(c *routekit.Context) {
			c.String(http.StatusOK, "hello from go-spring-starter-gin")
		})
	})
}

func userRoutes(handler *userHandler) routekit.Registrar {
	return routekit.NewRegistrar(200, func(engine *routekit.Engine, kit routekit.Kit) {
		group := engine.Group("/api/users", kit.Handler("login"))
		group.GET("/:id", handler.Profile)
		group.POST("", handler.Create)
		group.DELETE("/:id", kit.Handler("admin"), func(c *routekit.Context) {
			c.Status(http.StatusNoContent)
		})
	})
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

package main

import (
	"net/http"

	_ "github.com/JavaLionLi/go-spring-starter-gin"
	"github.com/JavaLionLi/go-spring-starter-gin/routekit"
	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-core/gs"
)

func init() {
	gs.Provide(func() routekit.Registrar {
		return routekit.NewRegistrar(100, func(engine *gin.Engine, kit routekit.Kit) {
			engine.GET("/hello", func(c *gin.Context) {
				c.String(http.StatusOK, "hello from go-spring-starter-gin")
			})
		})
	})
}

func main() {
	gs.Run()
}

# go-spring-starter-gin

[中文](README_zh-CN.md) | English

`go-spring-starter-gin` is a Gin starter for [Go-Spring](https://github.com/go-spring/spring-core).

It follows the official Go-Spring starter style:

- use a blank import to enable auto-configuration
- register beans during `init()`
- auto-create `*gin.Engine`
- auto-create `*gs.HttpServeMux` for Go-Spring's HTTP server
- let business modules contribute routes through ordered `routekit.Registrar` beans

## Install

```bash
go get github.com/JavaLionLi/go-spring-starter-gin
```

## Quick Start

Import the starter for side effects, then start Go-Spring normally.

```go
package main

import (
	"net/http"

	"github.com/JavaLionLi/go-spring-starter-gin/routekit"
	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-core/gs"
	_ "github.com/JavaLionLi/go-spring-starter-gin"
)

func init() {
	gs.Provide(func() routekit.Registrar {
		return routekit.NewRegistrar(100, func(engine *gin.Engine, kit routekit.Kit) {
			engine.GET("/hello", func(c *gin.Context) {
				c.String(http.StatusOK, "hello")
			})
		})
	})
}

func main() {
	gs.Run()
}
```

Start the application and visit:

```bash
curl http://localhost:9090/hello
curl http://localhost:9090/healthz
curl http://localhost:9090/ping
```

The default HTTP port comes from Go-Spring's built-in HTTP server config. Set `spring.http.server.addr` to change it.

## Configuration

The starter is enabled by default. You can disable it with:

```yaml
spring:
  gin:
    enabled: false
```

Complete example:

```yaml
spring:
  http:
    server:
      enabled: true
      addr: ":8080"
      headerTimeout: 5s
      readTimeout: 10s
      writeTimeout: 0s
      idleTimeout: 60s
  gin:
    enabled: true
    mode: release
    logger: true
    recovery: true
    trusted-proxies: []
    health:
      enabled: true
      healthz: /healthz
      ping: /ping
    cors:
      enabled: true
      allow-origins: ["*"]
      allow-methods: [GET, POST, PUT, PATCH, DELETE, OPTIONS]
      allow-headers: [Origin, Content-Type, Accept, Authorization]
      expose-headers: [Content-Length]
      allow-credentials: false
      max-age: 12h
```

## Register Routes

Each business module can provide one `routekit.Registrar`. Registrars are sorted by `Order()`.

```go
package user

import (
	"github.com/JavaLionLi/go-spring-starter-gin/routekit"
	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-core/gs"
)

func init() {
	gs.Provide(NewHandler)
	gs.Provide(NewRoutes)
}

func NewRoutes(handler *Handler) routekit.Registrar {
	return routekit.NewRegistrar(500, func(engine *gin.Engine, kit routekit.Kit) {
		group := engine.Group("/system/user", kit.Handler("login"))
		group.GET("", handler.List)
		group.POST("", handler.Create)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	})
}
```

If you keep a central module aggregator, import modules once:

```go
package all

import (
	_ "your-app/internal/modules/auth"
	_ "your-app/internal/modules/system/user"
)
```

Then import the aggregator in `main`:

```go
import (
	_ "github.com/JavaLionLi/go-spring-starter-gin"
	_ "your-app/internal/modules/all"
)
```

## Named Route Handlers

`routekit.Kit` is a small extension point for route-level middleware such as login, permission, role, encryption, or rate limiting.

Register named handlers:

```go
package auth

import (
	"github.com/JavaLionLi/go-spring-starter-gin/routekit"
	"github.com/go-spring/spring-core/gs"
)

func init() {
	gs.Provide(func() routekit.KitItem {
		return routekit.NewKitItem(100, func(kit *routekit.Kit) {
			kit.SetHandler("login", RequireLogin())
			kit.SetHandler("admin", RequireRole("admin"))
		})
	})
}
```

Use them in routes:

```go
func NewRoutes(handler *Handler) routekit.Registrar {
	return routekit.NewRegistrar(500, func(engine *gin.Engine, kit routekit.Kit) {
		group := engine.Group("/system/user", kit.Handler("login"))
		group.DELETE("/:id", kit.Handler("admin"), handler.Delete)
	})
}
```

Missing names return a no-op middleware, so modules can call `kit.Handler("login")` without nil checks.

## Global Middlewares

Use `startergin.Middleware` when a middleware should be mounted globally with `engine.Use(...)`.

```go
package httpx

import (
	startergin "github.com/JavaLionLi/go-spring-starter-gin"
	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-core/gs"
)

func init() {
	gs.Provide(func() startergin.Middleware {
		return startergin.NewMiddleware(100, RequestID())
	})
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Request-ID", "example")
		c.Next()
	}
}
```

Middlewares are sorted by order before registration.

## Engine Customization

Use `startergin.EngineConfigurer` for direct access to `*gin.Engine`, for example `NoRoute`, static files, trusted platform settings, or custom route groups.

```go
package httpx

import (
	"net/http"

	startergin "github.com/JavaLionLi/go-spring-starter-gin"
	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-core/gs"
)

func init() {
	gs.Provide(func() startergin.EngineConfigurer {
		return startergin.NewEngineConfigurer(100, func(engine *gin.Engine) {
			engine.NoRoute(func(c *gin.Context) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			})
		})
	})
}
```

Configurers run after global middlewares and health routes, before business registrars.

## Override Auto Configuration

If your application provides its own `*gin.Engine`, this starter will not create another one.

```go
func init() {
	gs.Provide(func() *gin.Engine {
		engine := gin.New()
		engine.Use(gin.Recovery())
		return engine
	})
}
```

If your application provides its own `*gs.HttpServeMux`, this starter will not create another one.

```go
func init() {
	gs.Provide(func(engine *gin.Engine) *gs.HttpServeMux {
		return &gs.HttpServeMux{Handler: engine}
	})
}
```

## Suggested Project Layout

```text
your-app/
  cmd/server/main.go
  configs/application.yml
  internal/modules/all/all.go
  internal/modules/auth/
  internal/modules/system/user/
```

`cmd/server/main.go`:

```go
package main

import (
	"github.com/go-spring/spring-core/gs"
	_ "github.com/JavaLionLi/go-spring-starter-gin"
	_ "your-app/internal/modules/all"
)

func main() {
	gs.Run()
}
```

## Demo

This repository includes a runnable demo in `examples/basic`.

```bash
cd examples/basic
go run .
```

Then visit:

```bash
curl http://localhost:8080/hello
curl http://localhost:8080/healthz
curl http://localhost:8080/ping
```

The demo uses `examples/basic/config/application.yml` to set `spring.http.server.addr` to `:8080`.

## Tests

Run all tests:

```bash
go test ./...
```

## Release

After pushing to GitHub, create a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Consumers can then install a fixed version:

```bash
go get github.com/JavaLionLi/go-spring-starter-gin@v0.1.0
```

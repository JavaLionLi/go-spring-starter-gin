# go-spring-starter-gin

[中文](README_zh-CN.md) | English

`go-spring-starter-gin` is a Gin starter for [Go-Spring](https://github.com/go-spring/spring-core).

It follows the official Go-Spring starter style:

- use a blank import to enable auto-configuration
- register beans during `init()`
- auto-create `*startergin.Engine`
- auto-create `*gs.HttpServeMux` for Go-Spring's HTTP server
- let business modules inject `*startergin.Engine` and register routes explicitly

## Install

```bash
go get github.com/CrazyLionCat/go-spring-starter-gin
```

## Quick Start

Import the starter for side effects, then start Go-Spring normally.

```go
package main

import (
	"net/http"

	startergin "github.com/CrazyLionCat/go-spring-starter-gin"
	"github.com/go-spring/spring-core/gs"
	_ "github.com/CrazyLionCat/go-spring-starter-gin"
)

type routes struct {
	Engine *startergin.Engine `autowire:""`
}

func init() {
	gs.Provide(&routes{}).Init((*routes).Register)
}

func (r *routes) Register() {
	r.Engine.GET("/hello", func(c *startergin.Context) {
		c.String(http.StatusOK, "hello")
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

The recommended style is to inject `*startergin.Engine` into a business bean and register routes during bean initialization. This keeps routes, middlewares, and dependencies close to the business code.

```go
package user

import (
	startergin "github.com/CrazyLionCat/go-spring-starter-gin"
	"github.com/go-spring/spring-core/gs"
)

type Routes struct {
	Engine  *startergin.Engine `autowire:""`
	Handler *Handler           `autowire:""`
}

func init() {
	gs.Provide(NewHandler)
	gs.Provide(&Routes{}).Init((*Routes).Register)
}

func (r *Routes) Register() {
	group := r.Engine.Group("/system/user", RequireLogin())
	group.GET("", r.Handler.List)
	group.POST("", r.Handler.Create)
	group.PUT("/:id", r.Handler.Update)
	group.DELETE("/:id", r.Handler.Delete)
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
	_ "github.com/CrazyLionCat/go-spring-starter-gin"
	_ "your-app/internal/modules/all"
)
```

## Named Route Handlers

For route-level middleware, use the regular Gin style and pass handlers to `Group` or individual routes.

```go
func (r *Routes) Register() {
	group := r.Engine.Group("/system/user", RequireLogin())
	group.DELETE("/:id", RequireRole("admin"), r.Handler.Delete)
}
```

## Global Middlewares

For global middleware, call `engine.Use(...)` from the same explicit initialization entry point.

```go
package httpx

import (
	startergin "github.com/CrazyLionCat/go-spring-starter-gin"
	"github.com/go-spring/spring-core/gs"
)

type Routes struct {
	Engine *startergin.Engine `autowire:""`
}

func init() {
	gs.Provide(&Routes{}).Init((*Routes).Register)
}

func (r *Routes) Register() {
	r.Engine.Use(RequestID())
}

func RequestID() startergin.HandlerFunc {
	return func(c *startergin.Context) {
		c.Header("X-Request-ID", "example")
		c.Next()
	}
}
```

If you need a custom CORS middleware outside the YAML config, use the starter wrapper instead of importing `gin-contrib/cors`:

```go
engine.Use(startergin.CORS(startergin.CORSHandlerConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
}))
```

## Engine Customization

For direct access to `*startergin.Engine`, use the same explicit initialization entry point. This is suitable for `NoRoute`, static files, trusted platform settings, or custom route groups.

```go
func (r *Routes) Register() {
	r.Engine.NoRoute(func(c *startergin.Context) {
		c.JSON(http.StatusNotFound, startergin.H{"error": "not found"})
	})
}
```

The starter's built-in setup runs logger/recovery, CORS, and health routes before business bean initialization.

## Advanced Auto Collection

For larger modular applications, you can still use the retained auto-collection extension points. The starter collects `startergin.Middleware`, `startergin.EngineConfigurer`, `routekit.KitItem`, and `routekit.Registrar`, then sorts them by `Order()`.

```go
package user

import (
	"github.com/CrazyLionCat/go-spring-starter-gin/routekit"
	"github.com/go-spring/spring-core/gs"
)

func init() {
	gs.Provide(NewHandler)
	gs.Provide(NewRoutes)
}

func NewRoutes(handler *Handler) routekit.Registrar {
	return routekit.NewRegistrar(500, func(engine *routekit.Engine, kit routekit.Kit) {
		group := engine.Group("/system/user", kit.Handler("login"))
		group.GET("", handler.List)
		group.POST("", handler.Create)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	})
}
```

`routekit.Kit` can hold named route middleware:

```go
func init() {
	gs.Provide(func() routekit.KitItem {
		return routekit.NewKitItem(100, func(kit *routekit.Kit) {
			kit.SetHandler("login", RequireLogin())
			kit.SetHandler("admin", RequireRole("admin"))
		})
	})
}
```

## Override Auto Configuration

If your application provides its own `*startergin.Engine`, this starter will not create another one.

```go
func init() {
	gs.Provide(func() *startergin.Engine {
		engine := startergin.New()
		engine.Use(startergin.Recovery())
		return engine
	})
}
```

If your application provides its own `*gs.HttpServeMux`, this starter will not create another one.

```go
func init() {
	gs.Provide(func(engine *startergin.Engine) *gs.HttpServeMux {
		return &gs.HttpServeMux{Handler: engine}
	})
}
```

## Suggested Project Layout

```text
your-app/
  cmd/server/main.go
  configs/application.yml
  internal/modules/auth/
  internal/modules/system/user/
```

`cmd/server/main.go`:

```go
package main

import (
	"github.com/go-spring/spring-core/gs"
	_ "github.com/CrazyLionCat/go-spring-starter-gin"
	_ "your-app/internal/modules/auth"
	_ "your-app/internal/modules/system/user"
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
curl -H "X-Demo-Token: demo-token" http://localhost:8080/api/users/42
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "X-Demo-Token: demo-token" \
  -d '{"name":"alice"}'
curl -X DELETE http://localhost:8080/api/users/42 \
  -H "X-Demo-Token: demo-token" \
  -H "X-Demo-Role: admin"
```

The demo uses `examples/basic/config/application.yml` to set `spring.http.server.addr` to `:8080` and show Gin mode, health, CORS, global middleware, route registration, and `NoRoute` customization.

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
go get github.com/CrazyLionCat/go-spring-starter-gin@v0.1.0
```

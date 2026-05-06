# go-spring-starter-gin

中文 | [English](README.md)

`go-spring-starter-gin` 是一个面向 [Go-Spring](https://github.com/go-spring/spring-core) 的 Gin starter。

它按 Go-Spring 官方 starter 的方式工作：

- 通过空白导入启用自动装配
- 在 `init()` 阶段注册 Bean
- 自动创建 `*gin.Engine`
- 自动创建 `*gs.HttpServeMux`，交给 Go-Spring 内置 HTTP Server 启动
- 业务模块通过有序的 `routekit.Registrar` 注册路由

## 安装

```bash
go get github.com/JavaLionLi/go-spring-starter-gin
```

## 快速开始

在入口文件中空白导入 starter，然后正常调用 `gs.Run()`。

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

启动后访问：

```bash
curl http://localhost:9090/hello
curl http://localhost:9090/healthz
curl http://localhost:9090/ping
```

默认端口来自 Go-Spring 内置 HTTP Server 配置。需要修改端口时，配置 `spring.http.server.addr`。

## 配置

starter 默认启用。如需关闭：

```yaml
spring:
  gin:
    enabled: false
```

完整配置示例：

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

## 注册业务路由

每个业务模块可以提供一个 `routekit.Registrar`。多个 Registrar 会按 `Order()` 从小到大注册。

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

如果项目里有统一的模块聚合包，可以只在一个地方导入所有模块：

```go
package all

import (
	_ "your-app/internal/modules/auth"
	_ "your-app/internal/modules/system/user"
)
```

入口文件中导入 starter 和聚合包：

```go
import (
	_ "github.com/JavaLionLi/go-spring-starter-gin"
	_ "your-app/internal/modules/all"
)
```

## 路由级中间件

`routekit.Kit` 用来放路由级通用中间件，比如登录校验、权限校验、角色校验、接口加密、限流等。

先注册命名中间件：

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

在路由里使用：

```go
func NewRoutes(handler *Handler) routekit.Registrar {
	return routekit.NewRegistrar(500, func(engine *gin.Engine, kit routekit.Kit) {
		group := engine.Group("/system/user", kit.Handler("login"))
		group.DELETE("/:id", kit.Handler("admin"), handler.Delete)
	})
}
```

如果指定名称不存在，`kit.Handler("xxx")` 会返回空操作中间件，因此业务模块不需要额外判空。

## 全局中间件

需要通过 `engine.Use(...)` 挂载的全局中间件，可以注册为 `startergin.Middleware`。

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

多个全局中间件会按 order 排序后注册。

## 自定义 Gin Engine

需要直接操作 `*gin.Engine` 时，可以注册 `startergin.EngineConfigurer`。适合配置 `NoRoute`、静态资源、特殊路由组等。

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

Configurer 的执行顺序在全局中间件和健康检查之后、业务路由注册之前。

## 覆盖自动装配

如果业务项目自己提供了 `*gin.Engine`，starter 不会重复创建。

```go
func init() {
	gs.Provide(func() *gin.Engine {
		engine := gin.New()
		engine.Use(gin.Recovery())
		return engine
	})
}
```

如果业务项目自己提供了 `*gs.HttpServeMux`，starter 也不会重复创建。

```go
func init() {
	gs.Provide(func(engine *gin.Engine) *gs.HttpServeMux {
		return &gs.HttpServeMux{Handler: engine}
	})
}
```

## 推荐项目结构

```text
your-app/
  cmd/server/main.go
  configs/application.yml
  internal/modules/all/all.go
  internal/modules/auth/
  internal/modules/system/user/
```

`cmd/server/main.go`：

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

## 使用 Demo

仓库内提供了一个可运行示例：`examples/basic`。

```bash
cd examples/basic
go run .
```

启动后访问：

```bash
curl http://localhost:8080/hello
curl http://localhost:8080/healthz
curl http://localhost:8080/ping
```

Demo 使用 `examples/basic/config/application.yml` 将 `spring.http.server.addr` 配置为 `:8080`。

## 测试

运行全部测试：

```bash
go test ./...
```

## 发布版本

推送到 GitHub 后，打一个 tag：

```bash
git tag v0.1.0
git push origin v0.1.0
```

其他项目即可固定版本引用：

```bash
go get github.com/JavaLionLi/go-spring-starter-gin@v0.1.0
```

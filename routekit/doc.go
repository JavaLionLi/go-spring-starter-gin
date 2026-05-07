// Package routekit 提供轻量的 Gin 路由注册辅助层。
//
// 它把路由注册器和可复用的命名处理器拆开，便于在 go-spring 容器中
// 按模块装配路由，同时保留 Gin 原生的路由声明方式。
package routekit

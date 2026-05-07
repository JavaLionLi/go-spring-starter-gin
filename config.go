package startergin

import "time"

// Config 控制 starter 自动创建的 Gin Engine。
//
// 字段通过 "spring.gin" 配置前缀绑定，例如 spring.gin.mode。
type Config struct {
	Mode           string       `value:"${mode:=release}"`     // Gin 运行模式，默认 release。
	TrustedProxies []string     `value:"${trusted-proxies:=}"` // 可信代理地址或 CIDR，留空时不调用 SetTrustedProxies。
	Logger         bool         `value:"${logger:=true}"`      // 是否挂载 Gin 默认访问日志中间件。
	Recovery       bool         `value:"${recovery:=true}"`    // 是否挂载 Gin 默认 panic 恢复中间件。
	Health         HealthConfig `value:"${health}"`            // 健康检查路由配置。
	CORS           CORSConfig   `value:"${cors}"`              // 跨域中间件配置。
}

// HealthConfig 控制 starter 自动注册的健康检查端点。
type HealthConfig struct {
	Enabled bool   `value:"${enabled:=true}"`     // 是否注册健康检查路由。
	Healthz string `value:"${healthz:=/healthz}"` // 返回 JSON 状态的健康检查路径，留空则不注册。
	Ping    string `value:"${ping:=/ping}"`       // 返回 pong 文本的探活路径，留空则不注册。
}

// CORSConfig 映射到 gin-contrib/cors.Config，用于自动挂载 CORS 中间件。
type CORSConfig struct {
	Enabled          bool          `value:"${enabled:=true}"`                                           // 是否启用 CORS 中间件。
	AllowOrigins     []string      `value:"${allow-origins:=*}"`                                        // 允许的来源列表，空白项会被忽略。
	AllowMethods     []string      `value:"${allow-methods:=GET,POST,PUT,PATCH,DELETE,OPTIONS}"`        // 允许的 HTTP 方法。
	AllowHeaders     []string      `value:"${allow-headers:=Origin,Content-Type,Accept,Authorization}"` // 允许的请求头。
	ExposeHeaders    []string      `value:"${expose-headers:=Content-Length}"`                          // 允许浏览器读取的响应头。
	AllowCredentials bool          `value:"${allow-credentials:=false}"`                                // 是否允许携带 Cookie 或认证信息。
	MaxAge           time.Duration `value:"${max-age:=12h}"`                                            // 预检请求结果缓存时间。
}

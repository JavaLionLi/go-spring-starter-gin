package startergin

import (
	"net/http"
	"sort"
	"strings"

	"github.com/JavaLionLi/go-spring-starter-gin/routekit"
	"github.com/gin-contrib/cors"
)

// EngineDeps 汇总创建 Engine 时可选注入的扩展点。
//
// go-spring 会自动装配实现了这些接口的 Bean；所有字段都标记为可选，
// 因此业务侧只注册自己需要的 Middleware、Configurer 或路由注册器即可。
type EngineDeps struct {
	Middlewares []Middleware         `autowire:"?"` // 按 Order 升序挂载到 Engine 的全局中间件。
	Configurers []EngineConfigurer   `autowire:"?"` // 按 Order 升序执行的 Engine 自定义配置。
	KitItems    []routekit.KitItem   `autowire:"?"` // 用于组装 routekit.Kit 的命名处理器。
	Registrars  []routekit.Registrar `autowire:"?"` // 按 Order 升序执行的路由注册器。
}

// NewEngine 根据配置和依赖创建一个可直接作为 HTTP Handler 使用的 Gin Engine。
//
// 创建顺序是固定的：设置 Gin 模式、配置可信代理、挂载内置中间件、
// 注册健康检查、挂载业务中间件、执行自定义配置，最后注册 routekit 路由。
func NewEngine(cfg Config, deps *EngineDeps) (*Engine, error) {
	if deps == nil {
		deps = &EngineDeps{}
	}
	if cfg.Mode != "" {
		SetMode(cfg.Mode)
	}

	engine := New()
	trustedProxies := cleanStrings(cfg.TrustedProxies)
	if len(trustedProxies) > 0 {
		// Gin 会校验代理地址格式；这里把错误返回给容器启动流程。
		if err := engine.SetTrustedProxies(trustedProxies); err != nil {
			return nil, err
		}
	}
	if cfg.Logger {
		engine.Use(Logger())
	}
	if cfg.Recovery {
		engine.Use(Recovery())
	}
	if cfg.CORS.Enabled {
		engine.Use(cors.New(cors.Config{
			// 配置中心传入空字符串时视为未配置，回退到 starter 的安全默认值。
			AllowOrigins:     defaultStrings(cleanStrings(cfg.CORS.AllowOrigins), []string{"*"}),
			AllowMethods:     defaultStrings(cleanStrings(cfg.CORS.AllowMethods), []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
			AllowHeaders:     defaultStrings(cleanStrings(cfg.CORS.AllowHeaders), []string{"Origin", "Content-Type", "Accept", "Authorization"}),
			ExposeHeaders:    defaultStrings(cleanStrings(cfg.CORS.ExposeHeaders), []string{"Content-Length"}),
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           cfg.CORS.MaxAge,
		}))
	}

	registerHealth(engine, cfg.Health)
	registerMiddlewares(engine, deps.Middlewares)
	runConfigurers(engine, deps.Configurers)
	routekit.RegisterAll(engine, routekit.NewKit(deps.KitItems), deps.Registrars)

	return engine, nil
}

// registerHealth 根据配置注册内置健康检查路由。
func registerHealth(engine *Engine, cfg HealthConfig) {
	if !cfg.Enabled {
		return
	}
	if cfg.Healthz != "" {
		engine.GET(cfg.Healthz, func(c *Context) {
			c.JSON(http.StatusOK, H{"status": "ok"})
		})
	}
	if cfg.Ping != "" {
		engine.GET(cfg.Ping, func(c *Context) {
			c.String(http.StatusOK, "pong")
		})
	}
}

// registerMiddlewares 按 Order 升序挂载业务提供的全局中间件。
//
// nil Middleware 或 nil Handler 会被忽略，便于条件装配场景直接返回空实现。
func registerMiddlewares(engine *Engine, middlewares []Middleware) {
	ordered := append([]Middleware(nil), middlewares...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i] == nil {
			return false
		}
		if ordered[j] == nil {
			return true
		}
		return ordered[i].Order() < ordered[j].Order()
	})
	for _, middleware := range ordered {
		if middleware == nil {
			continue
		}
		if handler := middleware.Handler(); handler != nil {
			engine.Use(handler)
		}
	}
}

// runConfigurers 按 Order 升序执行 Engine 配置器。
//
// 配置器适合注册 NoRoute、静态文件、模板、路由组或其他 Gin 原生设置。
func runConfigurers(engine *Engine, configurers []EngineConfigurer) {
	ordered := append([]EngineConfigurer(nil), configurers...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i] == nil {
			return false
		}
		if ordered[j] == nil {
			return true
		}
		return ordered[i].Order() < ordered[j].Order()
	})
	for _, configurer := range ordered {
		if configurer != nil {
			configurer.Configure(engine)
		}
	}
}

// defaultStrings 在清洗后的配置为空时返回默认值。
func defaultStrings(values []string, defaults []string) []string {
	if len(values) == 0 {
		return defaults
	}
	return values
}

// cleanStrings 去掉字符串切片中的前后空白和空字符串。
func cleanStrings(values []string) []string {
	cleaned := values[:0]
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			cleaned = append(cleaned, value)
		}
	}
	return cleaned
}

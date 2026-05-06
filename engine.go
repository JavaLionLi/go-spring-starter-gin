package startergin

import (
	"net/http"
	"sort"
	"strings"

	"github.com/JavaLionLi/go-spring-starter-gin/routekit"
	"github.com/gin-contrib/cors"
)

type EngineDeps struct {
	Middlewares []Middleware         `autowire:"?"`
	Configurers []EngineConfigurer   `autowire:"?"`
	KitItems    []routekit.KitItem   `autowire:"?"`
	Registrars  []routekit.Registrar `autowire:"?"`
}

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

func defaultStrings(values []string, defaults []string) []string {
	if len(values) == 0 {
		return defaults
	}
	return values
}

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

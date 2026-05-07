package startergin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CrazyLionCat/go-spring-starter-gin/routekit"
)

// TestNewEngineRegistersHealthMiddlewareConfigurerAndRoutes 覆盖 Engine 创建主流程：
// 健康检查、业务中间件排序、Configurer 执行和 routekit 路由注册。
func TestNewEngineRegistersHealthMiddlewareConfigurerAndRoutes(t *testing.T) {
	cfg := Config{
		Mode:     TestMode,
		Logger:   false,
		Recovery: false,
		Health: HealthConfig{
			Enabled: true,
			Healthz: "/healthz",
			Ping:    "/ping",
		},
		CORS: CORSConfig{Enabled: false},
	}
	deps := &EngineDeps{
		Middlewares: []Middleware{
			NewMiddleware(20, func(c *Context) {
				c.Header("X-Middleware-2", c.GetHeader("X-Middleware-1"))
				c.Next()
			}),
			NewMiddleware(10, func(c *Context) {
				c.Request.Header.Set("X-Middleware-1", "ok")
				c.Next()
			}),
		},
		Configurers: []EngineConfigurer{
			NewEngineConfigurer(10, func(engine *Engine) {
				engine.NoRoute(func(c *Context) {
					c.String(http.StatusNotFound, "missing")
				})
			}),
		},
		KitItems: []routekit.KitItem{
			routekit.NewKitItem(10, func(kit *routekit.Kit) {
				kit.SetHandler("mark", func(c *Context) {
					c.Header("X-Kit", "marked")
					c.Next()
				})
			}),
		},
		Registrars: []routekit.Registrar{
			routekit.NewRegistrar(10, func(engine *routekit.Engine, kit routekit.Kit) {
				engine.GET("/hello", kit.Handler("mark"), func(c *routekit.Context) {
					c.String(http.StatusOK, "hello")
				})
			}),
		},
	}

	engine, err := NewEngine(cfg, deps)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	assertResponse(t, engine, http.MethodGet, "/healthz", http.StatusOK, "")
	assertResponse(t, engine, http.MethodGet, "/ping", http.StatusOK, "pong")

	recorder := assertResponse(t, engine, http.MethodGet, "/hello", http.StatusOK, "hello")
	if got := recorder.Header().Get("X-Middleware-2"); got != "ok" {
		t.Fatalf("middleware order not applied, got %q", got)
	}
	if got := recorder.Header().Get("X-Kit"); got != "marked" {
		t.Fatalf("kit handler not applied, got %q", got)
	}

	assertResponse(t, engine, http.MethodGet, "/not-found", http.StatusNotFound, "missing")
}

// TestNewEngineWithNilDepsRegistersDefaultHealth 验证 deps 为 nil 时仍能创建 Engine，
// 并按配置注册默认健康检查路由。
func TestNewEngineWithNilDepsRegistersDefaultHealth(t *testing.T) {
	cfg := Config{
		Mode:     TestMode,
		Logger:   false,
		Recovery: false,
		Health: HealthConfig{
			Enabled: true,
			Healthz: "/ready",
			Ping:    "/live",
		},
		CORS: CORSConfig{Enabled: false},
	}

	engine, err := NewEngine(cfg, nil)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	assertResponse(t, engine, http.MethodGet, "/ready", http.StatusOK, `{"status":"ok"}`)
	assertResponse(t, engine, http.MethodGet, "/live", http.StatusOK, "pong")
}

// TestNewEngineSkipsDisabledHealthAndBlankHealthPaths 验证健康检查关闭或路径为空时不会注册路由。
func TestNewEngineSkipsDisabledHealthAndBlankHealthPaths(t *testing.T) {
	tests := []struct {
		name   string
		health HealthConfig
		paths  []string
	}{
		{
			name: "disabled",
			health: HealthConfig{
				Enabled: false,
				Healthz: "/healthz",
				Ping:    "/ping",
			},
			paths: []string{"/healthz", "/ping"},
		},
		{
			name: "blank paths",
			health: HealthConfig{
				Enabled: true,
				Healthz: "",
				Ping:    "",
			},
			paths: []string{"/healthz", "/ping"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine(Config{
				Mode:     TestMode,
				Logger:   false,
				Recovery: false,
				Health:   tt.health,
				CORS:     CORSConfig{Enabled: false},
			}, nil)
			if err != nil {
				t.Fatalf("new engine: %v", err)
			}

			for _, path := range tt.paths {
				assertResponse(t, engine, http.MethodGet, path, http.StatusNotFound, "404 page not found")
			}
		})
	}
}

// TestNewEngineConfiguresCORS 验证 spring.gin.cors 配置会转换成可用的 CORS 中间件。
func TestNewEngineConfiguresCORS(t *testing.T) {
	cfg := Config{
		Mode:     TestMode,
		Logger:   false,
		Recovery: false,
		Health:   HealthConfig{Enabled: false},
		CORS: CORSConfig{
			Enabled:          true,
			AllowOrigins:     []string{" https://app.example.com ", ""},
			AllowMethods:     []string{"GET", "POST", ""},
			AllowHeaders:     []string{"Authorization", "X-Request-ID"},
			ExposeHeaders:    []string{"X-Trace-ID"},
			AllowCredentials: true,
			MaxAge:           time.Hour,
		},
	}

	engine, err := NewEngine(cfg, &EngineDeps{
		Registrars: []routekit.Registrar{
			routekit.NewRegistrar(10, func(engine *routekit.Engine, kit routekit.Kit) {
				engine.POST("/api/items", func(c *routekit.Context) {
					c.Header("X-Trace-ID", "trace-1")
					c.String(http.StatusCreated, "created")
				})
			}),
		},
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	preflight := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/api/items", nil)
	request.Header.Set("Origin", "https://app.example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	request.Header.Set("Access-Control-Request-Headers", "Authorization")
	engine.ServeHTTP(preflight, request)

	if preflight.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d, want %d; body=%s", preflight.Code, http.StatusNoContent, preflight.Body.String())
	}
	if got := preflight.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("allow origin = %q", got)
	}
	if got := preflight.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("allow credentials = %q", got)
	}
	assertHeaderContains(t, preflight, "Access-Control-Allow-Methods", http.MethodPost)
	assertHeaderContains(t, preflight, "Access-Control-Allow-Headers", "Authorization")

	actual := httptest.NewRecorder()
	actualRequest := httptest.NewRequest(http.MethodPost, "/api/items", nil)
	actualRequest.Header.Set("Origin", "https://app.example.com")
	engine.ServeHTTP(actual, actualRequest)

	if actual.Code != http.StatusCreated {
		t.Fatalf("actual status = %d, want %d; body=%s", actual.Code, http.StatusCreated, actual.Body.String())
	}
	if got := actual.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("actual allow origin = %q", got)
	}
	assertHeaderContains(t, actual, "Access-Control-Expose-Headers", "X-Trace-ID")
}

// TestNewEngineUsesDefaultCORSListsWhenConfiguredListsAreBlank 验证 CORS 列表配置只有空白项时
// 会回退到 starter 默认值。
func TestNewEngineUsesDefaultCORSListsWhenConfiguredListsAreBlank(t *testing.T) {
	engine, err := NewEngine(Config{
		Mode:     TestMode,
		Logger:   false,
		Recovery: false,
		Health:   HealthConfig{Enabled: false},
		CORS: CORSConfig{
			Enabled:      true,
			AllowOrigins: []string{"", "   "},
			AllowMethods: []string{"", "   "},
			AllowHeaders: []string{"", "   "},
		},
	}, nil)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/anything", nil)
	request.Header.Set("Origin", "https://any.example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d, want %d; body=%s", recorder.Code, http.StatusNoContent, recorder.Body.String())
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("allow origin = %q, want *", got)
	}
}

// TestNewEngineIgnoresBlankTrustedProxies 验证空白可信代理不会触发 Gin 的格式校验错误。
func TestNewEngineIgnoresBlankTrustedProxies(t *testing.T) {
	cfg := Config{
		Mode:           TestMode,
		Logger:         false,
		Recovery:       false,
		TrustedProxies: []string{"", "   "},
		Health:         HealthConfig{Enabled: false},
		CORS:           CORSConfig{Enabled: false},
	}

	if _, err := NewEngine(cfg, nil); err != nil {
		t.Fatalf("new engine should ignore blank trusted proxies: %v", err)
	}
}

// TestNewEngineReturnsTrustedProxyErrors 验证非法可信代理配置会阻止 Engine 创建成功。
func TestNewEngineReturnsTrustedProxyErrors(t *testing.T) {
	cfg := Config{
		Mode:           TestMode,
		Logger:         false,
		Recovery:       false,
		TrustedProxies: []string{"not-a-cidr"},
		Health:         HealthConfig{Enabled: false},
		CORS:           CORSConfig{Enabled: false},
	}

	if _, err := NewEngine(cfg, nil); err == nil {
		t.Fatal("expected trusted proxy error")
	}
}

// TestNewEngineIgnoresNilDependencies 验证各类可选依赖中的 nil 值都会被安全跳过。
func TestNewEngineIgnoresNilDependencies(t *testing.T) {
	engine, err := NewEngine(Config{
		Mode:     TestMode,
		Logger:   false,
		Recovery: false,
		Health:   HealthConfig{Enabled: false},
		CORS:     CORSConfig{Enabled: false},
	}, &EngineDeps{
		Middlewares: []Middleware{
			nil,
			NewMiddleware(10, nil),
			NewMiddleware(20, func(c *Context) {
				c.Header("X-Middleware", "applied")
				c.Next()
			}),
		},
		Configurers: []EngineConfigurer{
			nil,
			NewEngineConfigurer(10, nil),
			NewEngineConfigurer(20, func(engine *Engine) {
				engine.GET("/configured", func(c *Context) {
					c.String(http.StatusOK, "configured")
				})
			}),
		},
		KitItems: []routekit.KitItem{nil},
		Registrars: []routekit.Registrar{
			nil,
			routekit.NewRegistrar(10, nil),
		},
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	recorder := assertResponse(t, engine, http.MethodGet, "/configured", http.StatusOK, "configured")
	if got := recorder.Header().Get("X-Middleware"); got != "applied" {
		t.Fatalf("middleware header = %q", got)
	}
}

// assertResponse 发送请求并断言状态码和完整响应体。
func assertResponse(t *testing.T, engine *Engine, method string, path string, status int, contains string) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != status {
		t.Fatalf("%s %s status = %d, want %d; body=%s", method, path, recorder.Code, status, recorder.Body.String())
	}
	if contains != "" && recorder.Body.String() != contains {
		t.Fatalf("%s %s body = %q, want %q", method, path, recorder.Body.String(), contains)
	}
	return recorder
}

// assertHeaderContains 以大小写不敏感方式断言响应头包含指定值。
func assertHeaderContains(t *testing.T, recorder *httptest.ResponseRecorder, name string, value string) {
	t.Helper()

	got := strings.ToLower(recorder.Header().Get(name))
	want := strings.ToLower(value)
	if !strings.Contains(got, want) {
		t.Fatalf("%s = %q, want to contain %q", name, got, value)
	}
}

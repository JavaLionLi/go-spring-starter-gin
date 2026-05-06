package startergin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JavaLionLi/go-spring-starter-gin/routekit"
	"github.com/gin-gonic/gin"
)

func TestNewEngineRegistersHealthMiddlewareConfigurerAndRoutes(t *testing.T) {
	cfg := Config{
		Mode:     gin.TestMode,
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
			NewMiddleware(20, func(c *gin.Context) {
				c.Header("X-Middleware-2", c.GetHeader("X-Middleware-1"))
				c.Next()
			}),
			NewMiddleware(10, func(c *gin.Context) {
				c.Request.Header.Set("X-Middleware-1", "ok")
				c.Next()
			}),
		},
		Configurers: []EngineConfigurer{
			NewEngineConfigurer(10, func(engine *gin.Engine) {
				engine.NoRoute(func(c *gin.Context) {
					c.String(http.StatusNotFound, "missing")
				})
			}),
		},
		KitItems: []routekit.KitItem{
			routekit.NewKitItem(10, func(kit *routekit.Kit) {
				kit.SetHandler("mark", func(c *gin.Context) {
					c.Header("X-Kit", "marked")
					c.Next()
				})
			}),
		},
		Registrars: []routekit.Registrar{
			routekit.NewRegistrar(10, func(engine *gin.Engine, kit routekit.Kit) {
				engine.GET("/hello", kit.Handler("mark"), func(c *gin.Context) {
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

func TestNewEngineIgnoresBlankTrustedProxies(t *testing.T) {
	cfg := Config{
		Mode:           gin.TestMode,
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

func assertResponse(t *testing.T, engine *gin.Engine, method string, path string, status int, contains string) *httptest.ResponseRecorder {
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

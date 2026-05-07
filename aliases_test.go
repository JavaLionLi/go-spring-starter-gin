package startergin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestReExportedGinHelpersBuildRoutesWithoutGinImport 验证业务代码只导入 startergin
// 时，仍能使用常见 Gin 类型、认证中间件和标准库 Handler 包装能力。
func TestReExportedGinHelpersBuildRoutesWithoutGinImport(t *testing.T) {
	SetMode(TestMode)

	engine := New()
	engine.GET("/json", func(c *Context) {
		c.JSON(http.StatusOK, H{"ok": true})
	})
	engine.GET("/wrapped", WrapF(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("wrapped"))
	}))
	engine.GET("/secure", BasicAuth(Accounts{"demo": "secret"}), func(c *Context) {
		c.String(http.StatusOK, c.GetString(AuthUserKey))
	})

	json := aliasRequest(engine, http.MethodGet, "/json", "")
	if json.Code != http.StatusOK || json.Body.String() != `{"ok":true}` {
		t.Fatalf("/json status = %d body = %q", json.Code, json.Body.String())
	}

	wrapped := aliasRequest(engine, http.MethodGet, "/wrapped", "")
	if wrapped.Code != http.StatusAccepted || wrapped.Body.String() != "wrapped" {
		t.Fatalf("/wrapped status = %d body = %q", wrapped.Code, wrapped.Body.String())
	}

	unauthorized := aliasRequest(engine, http.MethodGet, "/secure", "")
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("/secure unauthorized status = %d", unauthorized.Code)
	}

	authorized := aliasRequest(engine, http.MethodGet, "/secure", "Basic ZGVtbzpzZWNyZXQ=")
	if authorized.Code != http.StatusOK || authorized.Body.String() != "demo" {
		t.Fatalf("/secure authorized status = %d body = %q", authorized.Code, authorized.Body.String())
	}
}

// TestCORSHelperBuildsCorsMiddlewareWithoutCorsImport 验证 CORS 辅助函数可直接
// 基于 startergin 暴露的配置类型创建跨域中间件。
func TestCORSHelperBuildsCorsMiddlewareWithoutCorsImport(t *testing.T) {
	SetMode(TestMode)

	engine := New()
	engine.Use(CORS(CORSHandlerConfig{
		AllowOrigins: []string{"https://app.example.com"},
		AllowMethods: []string{http.MethodPost},
		AllowHeaders: []string{"Authorization"},
	}))
	engine.POST("/items", func(c *Context) {
		c.Status(http.StatusCreated)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/items", nil)
	request.Header.Set("Origin", "https://app.example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	request.Header.Set("Access-Control-Request-Headers", "Authorization")
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d body = %q", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("allow origin = %q", got)
	}
}

// TestCreateTestContextAliases 验证测试辅助函数的类型别名能正常返回 Context 和 Engine。
func TestCreateTestContextAliases(t *testing.T) {
	recorder := httptest.NewRecorder()
	context, engine := CreateTestContext(recorder)
	if context == nil {
		t.Fatal("expected test context")
	}
	if engine == nil {
		t.Fatal("expected test engine")
	}
}

// aliasRequest 封装测试请求创建，减少别名透传测试中的重复代码。
func aliasRequest(engine *Engine, method string, path string, authorization string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, nil)
	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}
	engine.ServeHTTP(recorder, request)
	return recorder
}

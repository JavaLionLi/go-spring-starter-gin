package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	startergin "github.com/JavaLionLi/go-spring-starter-gin"
)

func TestDemoRoutes(t *testing.T) {
	startergin.SetMode(startergin.TestMode)

	engine, err := startergin.NewEngine(startergin.Config{
		Mode:     startergin.TestMode,
		Logger:   false,
		Recovery: false,
		Health:   startergin.HealthConfig{Enabled: false},
		CORS:     startergin.CORSConfig{Enabled: false},
	}, nil)
	if err != nil {
		t.Fatalf("new demo engine: %v", err)
	}
	(&routes{
		Engine:      engine,
		UserHandler: newUserHandler(),
	}).Register()

	public := demoRequest(engine, http.MethodGet, "/hello", nil)
	if public.Code != http.StatusOK || public.Body.String() != "hello from go-spring-starter-gin" {
		t.Fatalf("/hello status = %d body = %q", public.Code, public.Body.String())
	}
	if got := public.Header().Get("X-Request-ID"); got == "" {
		t.Fatal("expected request id header")
	}

	anonymous := demoRequest(engine, http.MethodGet, "/api/users/42", nil)
	if anonymous.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous user status = %d", anonymous.Code)
	}

	authenticated := demoRequest(engine, http.MethodGet, "/api/users/42", map[string]string{
		"X-Demo-Token": "demo-token",
	})
	if authenticated.Code != http.StatusOK {
		t.Fatalf("authenticated user status = %d body = %q", authenticated.Code, authenticated.Body.String())
	}
	if !strings.Contains(authenticated.Body.String(), `"id":"42"`) {
		t.Fatalf("unexpected profile body: %q", authenticated.Body.String())
	}

	created := demoJSONRequest(engine, http.MethodPost, "/api/users", `{"name":"alice"}`, map[string]string{
		"X-Demo-Token": "demo-token",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d body = %q", created.Code, created.Body.String())
	}
	if !strings.Contains(created.Body.String(), `"name":"alice"`) {
		t.Fatalf("unexpected create body: %q", created.Body.String())
	}

	invalid := demoJSONRequest(engine, http.MethodPost, "/api/users", `{}`, map[string]string{
		"X-Demo-Token": "demo-token",
	})
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid create status = %d body = %q", invalid.Code, invalid.Body.String())
	}

	withoutRole := demoRequest(engine, http.MethodDelete, "/api/users/42", map[string]string{
		"X-Demo-Token": "demo-token",
	})
	if withoutRole.Code != http.StatusUnauthorized {
		t.Fatalf("delete without role status = %d", withoutRole.Code)
	}

	deleted := demoRequest(engine, http.MethodDelete, "/api/users/42", map[string]string{
		"X-Demo-Token": "demo-token",
		"X-Demo-Role":  "admin",
	})
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d body = %q", deleted.Code, deleted.Body.String())
	}

	missing := demoRequest(engine, http.MethodGet, "/missing", nil)
	if missing.Code != http.StatusNotFound || !strings.Contains(missing.Body.String(), "route not found") {
		t.Fatalf("missing status = %d body = %q", missing.Code, missing.Body.String())
	}
}

func demoRequest(engine *startergin.Engine, method string, path string, headers map[string]string) *httptest.ResponseRecorder {
	return demoJSONRequest(engine, method, path, "", headers)
}

func demoJSONRequest(engine *startergin.Engine, method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	engine.ServeHTTP(recorder, request)
	return recorder
}

package startergin

import "testing"

// TestNewHTTPServeMux 验证 Gin Engine 会被原样包装到 go-spring HttpServeMux 中。
func TestNewHTTPServeMux(t *testing.T) {
	engine := New()
	mux, err := NewHTTPServeMux(engine)
	if err != nil {
		t.Fatalf("new serve mux: %v", err)
	}
	if mux.Handler != engine {
		t.Fatalf("unexpected handler: %#v", mux.Handler)
	}
}

// TestNewHTTPServeMuxRejectsNilEngine 验证缺失 Engine 时返回明确错误，避免启动空服务。
func TestNewHTTPServeMuxRejectsNilEngine(t *testing.T) {
	if _, err := NewHTTPServeMux(nil); err == nil {
		t.Fatal("expected error")
	}
}

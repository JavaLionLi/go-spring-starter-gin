package routekit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRegisterAllSortsRegistrarsByOrder 验证路由注册器按 Order 稳定排序，
// 并且 nil 或空函数注册器不会影响其他注册器执行。
func TestRegisterAllSortsRegistrarsByOrder(t *testing.T) {
	SetMode(TestMode)

	engine := New()
	calls := make([]int, 0, 4)

	RegisterAll(engine, Kit{}, []Registrar{
		NewRegistrar(20, func(*Engine, Kit) {
			calls = append(calls, 20)
		}),
		nil,
		NewRegistrar(10, nil),
		NewRegistrar(10, func(*Engine, Kit) {
			calls = append(calls, 10)
		}),
		NewRegistrar(10, func(*Engine, Kit) {
			calls = append(calls, 11)
		}),
	})

	if len(calls) != 3 || calls[0] != 10 || calls[1] != 11 || calls[2] != 20 {
		t.Fatalf("unexpected register order: %#v", calls)
	}
}

// TestKitHandlerUsesNamedHandler 验证 Kit 会返回最后一次写入的同名处理器。
func TestKitHandlerUsesNamedHandler(t *testing.T) {
	SetMode(TestMode)

	kit := NewKit([]KitItem{
		NewKitItem(20, func(kit *Kit) {
			kit.SetHandler("auth", func(c *Context) {
				c.Header("X-Auth", "late")
				c.Next()
			})
		}),
		NewKitItem(10, func(kit *Kit) {
			kit.SetHandler("auth", func(c *Context) {
				c.Header("X-Auth", "early")
				c.Next()
			})
		}),
	})

	engine := New()
	engine.GET("/secure", kit.Handler("auth"), func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/secure", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	if got := recorder.Header().Get("X-Auth"); got != "late" {
		t.Fatalf("unexpected handler result: %q", got)
	}
}

// TestKitIgnoresNilItemsBlankNamesAndNilHandlers 验证 Kit 对 nil、空名称和 nil Handler 的容错。
func TestKitIgnoresNilItemsBlankNamesAndNilHandlers(t *testing.T) {
	kit := NewKit([]KitItem{
		nil,
		NewKitItem(10, nil),
		NewKitItem(20, func(kit *Kit) {
			kit.SetHandler("", func(c *Context) {
				c.AbortWithStatus(http.StatusInternalServerError)
			})
			kit.SetHandler("nil", nil)
			kit.SetHandler("ok", func(c *Context) {
				c.Header("X-Kit", "ok")
				c.Next()
			})
		}),
	})

	engine := New()
	engine.GET("/ok", kit.Handler("ok"), func(c *Context) {
		c.String(http.StatusOK, "ok")
	})
	engine.GET("/blank", kit.Handler(""), func(c *Context) {
		c.String(http.StatusOK, "blank")
	})
	engine.GET("/nil", kit.Handler("nil"), func(c *Context) {
		c.String(http.StatusOK, "nil")
	})

	ok := performRequest(engine, http.MethodGet, "/ok")
	if ok.Code != http.StatusOK {
		t.Fatalf("/ok status = %d", ok.Code)
	}
	if got := ok.Header().Get("X-Kit"); got != "ok" {
		t.Fatalf("kit header = %q", got)
	}

	blank := performRequest(engine, http.MethodGet, "/blank")
	if blank.Code != http.StatusOK || blank.Body.String() != "blank" {
		t.Fatalf("/blank status = %d body = %q", blank.Code, blank.Body.String())
	}

	nilHandler := performRequest(engine, http.MethodGet, "/nil")
	if nilHandler.Code != http.StatusOK || nilHandler.Body.String() != "nil" {
		t.Fatalf("/nil status = %d body = %q", nilHandler.Code, nilHandler.Body.String())
	}
}

// TestSetHandlerInitializesNilKitMap 验证零值 Kit 也可以直接注册命名处理器。
func TestSetHandlerInitializesNilKitMap(t *testing.T) {
	var kit Kit
	kit.SetHandler("trace", func(c *Context) {
		c.Header("X-Trace", "set")
		c.Next()
	})

	engine := New()
	engine.GET("/trace", kit.Handler("trace"), func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	recorder := performRequest(engine, http.MethodGet, "/trace")
	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	if got := recorder.Header().Get("X-Trace"); got != "set" {
		t.Fatalf("unexpected trace header: %q", got)
	}
}

// TestMissingKitHandlerIsNoOp 验证缺失的命名处理器会退化为 NoOp，不阻断后续处理。
func TestMissingKitHandlerIsNoOp(t *testing.T) {
	SetMode(TestMode)

	engine := New()
	engine.GET("/open", Kit{}.Handler("missing"), func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/open", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
}

// TestNoOpContinuesToNextHandler 验证 NoOp 只继续执行后续 Handler。
func TestNoOpContinuesToNextHandler(t *testing.T) {
	SetMode(TestMode)

	engine := New()
	engine.GET("/next", NoOp(), func(c *Context) {
		c.String(http.StatusAccepted, "next")
	})

	recorder := performRequest(engine, http.MethodGet, "/next")
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	if recorder.Body.String() != "next" {
		t.Fatalf("unexpected body: %q", recorder.Body.String())
	}
}

// performRequest 封装 routekit 测试中的 HTTP 请求发送。
func performRequest(engine *Engine, method string, path string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, nil)
	engine.ServeHTTP(recorder, request)
	return recorder
}

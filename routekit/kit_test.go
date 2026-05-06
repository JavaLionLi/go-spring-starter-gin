package routekit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterAllSortsRegistrarsByOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	calls := make([]int, 0, 4)

	RegisterAll(engine, Kit{}, []Registrar{
		NewRegistrar(20, func(*gin.Engine, Kit) {
			calls = append(calls, 20)
		}),
		nil,
		NewRegistrar(10, nil),
		NewRegistrar(10, func(*gin.Engine, Kit) {
			calls = append(calls, 10)
		}),
		NewRegistrar(10, func(*gin.Engine, Kit) {
			calls = append(calls, 11)
		}),
	})

	if len(calls) != 3 || calls[0] != 10 || calls[1] != 11 || calls[2] != 20 {
		t.Fatalf("unexpected register order: %#v", calls)
	}
}

func TestKitHandlerUsesNamedHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	kit := NewKit([]KitItem{
		NewKitItem(20, func(kit *Kit) {
			kit.SetHandler("auth", func(c *gin.Context) {
				c.Header("X-Auth", "late")
				c.Next()
			})
		}),
		NewKitItem(10, func(kit *Kit) {
			kit.SetHandler("auth", func(c *gin.Context) {
				c.Header("X-Auth", "early")
				c.Next()
			})
		}),
	})

	engine := gin.New()
	engine.GET("/secure", kit.Handler("auth"), func(c *gin.Context) {
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

func TestKitIgnoresNilItemsBlankNamesAndNilHandlers(t *testing.T) {
	kit := NewKit([]KitItem{
		nil,
		NewKitItem(10, nil),
		NewKitItem(20, func(kit *Kit) {
			kit.SetHandler("", func(c *gin.Context) {
				c.AbortWithStatus(http.StatusInternalServerError)
			})
			kit.SetHandler("nil", nil)
			kit.SetHandler("ok", func(c *gin.Context) {
				c.Header("X-Kit", "ok")
				c.Next()
			})
		}),
	})

	engine := gin.New()
	engine.GET("/ok", kit.Handler("ok"), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	engine.GET("/blank", kit.Handler(""), func(c *gin.Context) {
		c.String(http.StatusOK, "blank")
	})
	engine.GET("/nil", kit.Handler("nil"), func(c *gin.Context) {
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

func TestSetHandlerInitializesNilKitMap(t *testing.T) {
	var kit Kit
	kit.SetHandler("trace", func(c *gin.Context) {
		c.Header("X-Trace", "set")
		c.Next()
	})

	engine := gin.New()
	engine.GET("/trace", kit.Handler("trace"), func(c *gin.Context) {
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

func TestMissingKitHandlerIsNoOp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.GET("/open", Kit{}.Handler("missing"), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/open", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
}

func TestNoOpContinuesToNextHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.GET("/next", NoOp(), func(c *gin.Context) {
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

func performRequest(engine *gin.Engine, method string, path string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, nil)
	engine.ServeHTTP(recorder, request)
	return recorder
}

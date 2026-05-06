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
	calls := make([]int, 0, 2)

	RegisterAll(engine, Kit{}, []Registrar{
		NewRegistrar(20, func(*gin.Engine, Kit) {
			calls = append(calls, 20)
		}),
		NewRegistrar(10, func(*gin.Engine, Kit) {
			calls = append(calls, 10)
		}),
	})

	if len(calls) != 2 || calls[0] != 10 || calls[1] != 20 {
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

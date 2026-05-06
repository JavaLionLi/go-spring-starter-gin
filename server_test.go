package startergin

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewHTTPServeMux(t *testing.T) {
	engine := gin.New()
	mux, err := NewHTTPServeMux(engine)
	if err != nil {
		t.Fatalf("new serve mux: %v", err)
	}
	if mux.Handler != engine {
		t.Fatalf("unexpected handler: %#v", mux.Handler)
	}
}

func TestNewHTTPServeMuxRejectsNilEngine(t *testing.T) {
	if _, err := NewHTTPServeMux(nil); err == nil {
		t.Fatal("expected error")
	}
}

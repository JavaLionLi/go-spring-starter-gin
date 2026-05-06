package startergin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-core/gs"
)

func NewHTTPServeMux(engine *gin.Engine) (*gs.HttpServeMux, error) {
	if engine == nil {
		return nil, fmt.Errorf("gin engine is not initialized")
	}
	return &gs.HttpServeMux{Handler: engine}, nil
}

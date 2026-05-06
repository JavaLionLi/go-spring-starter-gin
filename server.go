package startergin

import (
	"fmt"

	"github.com/go-spring/spring-core/gs"
)

func NewHTTPServeMux(engine *Engine) (*gs.HttpServeMux, error) {
	if engine == nil {
		return nil, fmt.Errorf("gin engine is not initialized")
	}
	return &gs.HttpServeMux{Handler: engine}, nil
}

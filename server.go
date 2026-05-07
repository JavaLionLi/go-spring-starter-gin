package startergin

import (
	"fmt"

	"github.com/go-spring/spring-core/gs"
)

// NewHTTPServeMux 将 Gin Engine 包装成 go-spring 使用的 HttpServeMux。
//
// go-spring 的 HTTP server 组件会消费 HttpServeMux，因此这里是 Gin 与
// go-spring Web 启动流程之间的适配层。
func NewHTTPServeMux(engine *Engine) (*gs.HttpServeMux, error) {
	if engine == nil {
		return nil, fmt.Errorf("gin engine is not initialized")
	}
	return &gs.HttpServeMux{Handler: engine}, nil
}

package startergin

import "github.com/gin-gonic/gin"

type Middleware interface {
	Order() int
	Handler() gin.HandlerFunc
}

type middleware struct {
	order   int
	handler gin.HandlerFunc
}

func NewMiddleware(order int, handler gin.HandlerFunc) Middleware {
	return middleware{order: order, handler: handler}
}

func (m middleware) Order() int {
	return m.order
}

func (m middleware) Handler() gin.HandlerFunc {
	return m.handler
}

type EngineConfigurer interface {
	Order() int
	Configure(*gin.Engine)
}

type engineConfigurer struct {
	order int
	fn    func(*gin.Engine)
}

func NewEngineConfigurer(order int, fn func(*gin.Engine)) EngineConfigurer {
	return engineConfigurer{order: order, fn: fn}
}

func (c engineConfigurer) Order() int {
	return c.order
}

func (c engineConfigurer) Configure(engine *gin.Engine) {
	if c.fn != nil {
		c.fn(engine)
	}
}

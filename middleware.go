package startergin

type Middleware interface {
	Order() int
	Handler() HandlerFunc
}

type middleware struct {
	order   int
	handler HandlerFunc
}

func NewMiddleware(order int, handler HandlerFunc) Middleware {
	return middleware{order: order, handler: handler}
}

func (m middleware) Order() int {
	return m.order
}

func (m middleware) Handler() HandlerFunc {
	return m.handler
}

type EngineConfigurer interface {
	Order() int
	Configure(*Engine)
}

type engineConfigurer struct {
	order int
	fn    func(*Engine)
}

func NewEngineConfigurer(order int, fn func(*Engine)) EngineConfigurer {
	return engineConfigurer{order: order, fn: fn}
}

func (c engineConfigurer) Order() int {
	return c.order
}

func (c engineConfigurer) Configure(engine *Engine) {
	if c.fn != nil {
		c.fn(engine)
	}
}

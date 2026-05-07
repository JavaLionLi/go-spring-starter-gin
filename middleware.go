package startergin

// Middleware 表示可被 starter 自动挂载到 Gin Engine 的全局中间件。
//
// 多个 Middleware 会按 Order 升序执行；相同 Order 保持 go-spring 注入顺序。
type Middleware interface {
	// Order 返回中间件挂载顺序，数值越小越早执行。
	Order() int
	// Handler 返回实际的 Gin 处理函数；返回 nil 时会被忽略。
	Handler() HandlerFunc
}

type middleware struct {
	order   int
	handler HandlerFunc
}

// NewMiddleware 用函数式方式创建一个 Middleware Bean。
//
// 业务代码通常将返回值注册到 go-spring 容器中，由 starter 自动收集并挂载。
func NewMiddleware(order int, handler HandlerFunc) Middleware {
	return middleware{order: order, handler: handler}
}

// Order 返回该中间件的排序值。
func (m middleware) Order() int {
	return m.order
}

// Handler 返回该中间件持有的 Gin HandlerFunc。
func (m middleware) Handler() HandlerFunc {
	return m.handler
}

// EngineConfigurer 表示在 Engine 创建完成后执行的自定义配置扩展点。
//
// 它适合做路由分组、模板、静态资源、NoRoute 等需要直接访问 Engine 的设置。
type EngineConfigurer interface {
	// Order 返回配置器执行顺序，数值越小越早执行。
	Order() int
	// Configure 对 Engine 应用自定义配置。
	Configure(*Engine)
}

type engineConfigurer struct {
	order int
	fn    func(*Engine)
}

// NewEngineConfigurer 用函数式方式创建一个 EngineConfigurer Bean。
func NewEngineConfigurer(order int, fn func(*Engine)) EngineConfigurer {
	return engineConfigurer{order: order, fn: fn}
}

// Order 返回该配置器的排序值。
func (c engineConfigurer) Order() int {
	return c.order
}

// Configure 在函数不为空时执行实际配置逻辑。
func (c engineConfigurer) Configure(engine *Engine) {
	if c.fn != nil {
		c.fn(engine)
	}
}

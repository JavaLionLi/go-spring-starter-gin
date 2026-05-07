package routekit

import "sort"

// Registrar 表示一组路由注册逻辑。
//
// starter 会收集所有 Registrar 并按 Order 升序执行，适合把不同业务模块的
// 路由声明拆成独立 Bean。
type Registrar interface {
	// Order 返回路由注册顺序，数值越小越早执行。
	Order() int
	// RegisterRoutes 在指定 Engine 上注册路由，并可从 Kit 取命名处理器。
	RegisterRoutes(*Engine, Kit)
}

type registrar struct {
	order int
	fn    func(*Engine, Kit)
}

// NewRegistrar 用函数式方式创建一个 Registrar。
func NewRegistrar(order int, fn func(*Engine, Kit)) Registrar {
	return registrar{order: order, fn: fn}
}

// Order 返回该注册器的排序值。
func (r registrar) Order() int {
	return r.order
}

// RegisterRoutes 在函数不为空时执行实际路由注册逻辑。
func (r registrar) RegisterRoutes(engine *Engine, kit Kit) {
	if r.fn != nil {
		r.fn(engine, kit)
	}
}

// RegisterAll 按 Order 升序执行所有路由注册器。
//
// nil Registrar 会被忽略，便于条件装配或测试场景复用同一组依赖。
func RegisterAll(engine *Engine, kit Kit, registrars []Registrar) {
	ordered := append([]Registrar(nil), registrars...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i] == nil {
			return false
		}
		if ordered[j] == nil {
			return true
		}
		return ordered[i].Order() < ordered[j].Order()
	})
	for _, registrar := range ordered {
		if registrar != nil {
			registrar.RegisterRoutes(engine, kit)
		}
	}
}

// Kit 保存按名称索引的 Gin HandlerFunc。
//
// 它用于在路由注册器之间共享中间件或处理器，例如 auth、trace、tenant 等。
type Kit struct {
	handlers map[string]HandlerFunc
}

// NewKit 根据 KitItem 列表构建 Kit。
//
// KitItem 会按 Order 升序应用；相同名称的 Handler 后写入者覆盖先写入者。
func NewKit(items []KitItem) Kit {
	kit := Kit{handlers: map[string]HandlerFunc{}}
	ordered := append([]KitItem(nil), items...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i] == nil {
			return false
		}
		if ordered[j] == nil {
			return true
		}
		return ordered[i].Order() < ordered[j].Order()
	})
	for _, item := range ordered {
		if item != nil {
			item.Apply(&kit)
		}
	}
	return kit
}

// Handler 返回指定名称的处理器。
//
// 未找到名称或 Kit 尚未初始化时返回 NoOp，这样路由声明可以保持链式写法，
// 不需要在每个调用点单独判断 nil。
func (k Kit) Handler(name string) HandlerFunc {
	if k.handlers == nil {
		return NoOp()
	}
	if handler := k.handlers[name]; handler != nil {
		return handler
	}
	return NoOp()
}

// SetHandler 注册一个命名处理器。
//
// 空名称或 nil Handler 会被忽略，避免错误配置污染 Kit。
func (k *Kit) SetHandler(name string, handler HandlerFunc) {
	if k.handlers == nil {
		k.handlers = map[string]HandlerFunc{}
	}
	if name != "" && handler != nil {
		k.handlers[name] = handler
	}
}

// NoOp 返回一个只调用 Next 的空处理器。
//
// 它作为缺失命名处理器的安全 fallback，不会中断后续路由处理。
func NoOp() HandlerFunc {
	return func(c *Context) {
		c.Next()
	}
}

// KitItem 表示对 Kit 的一次配置。
//
// 业务侧可通过 KitItem 注册命名中间件，再在 Registrar 中按名称引用。
type KitItem interface {
	// Order 返回 KitItem 应用顺序，数值越小越早执行。
	Order() int
	// Apply 对 Kit 写入或调整命名处理器。
	Apply(*Kit)
}

type kitItem struct {
	order int
	fn    func(*Kit)
}

// NewKitItem 用函数式方式创建一个 KitItem。
func NewKitItem(order int, fn func(*Kit)) KitItem {
	return kitItem{order: order, fn: fn}
}

// Order 返回该 KitItem 的排序值。
func (i kitItem) Order() int {
	return i.order
}

// Apply 在函数不为空时执行实际 Kit 配置逻辑。
func (i kitItem) Apply(kit *Kit) {
	if i.fn != nil {
		i.fn(kit)
	}
}

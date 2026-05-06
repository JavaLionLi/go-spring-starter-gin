package routekit

import "sort"

type Registrar interface {
	Order() int
	RegisterRoutes(*Engine, Kit)
}

type registrar struct {
	order int
	fn    func(*Engine, Kit)
}

func NewRegistrar(order int, fn func(*Engine, Kit)) Registrar {
	return registrar{order: order, fn: fn}
}

func (r registrar) Order() int {
	return r.order
}

func (r registrar) RegisterRoutes(engine *Engine, kit Kit) {
	if r.fn != nil {
		r.fn(engine, kit)
	}
}

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

type Kit struct {
	handlers map[string]HandlerFunc
}

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

func (k Kit) Handler(name string) HandlerFunc {
	if k.handlers == nil {
		return NoOp()
	}
	if handler := k.handlers[name]; handler != nil {
		return handler
	}
	return NoOp()
}

func (k *Kit) SetHandler(name string, handler HandlerFunc) {
	if k.handlers == nil {
		k.handlers = map[string]HandlerFunc{}
	}
	if name != "" && handler != nil {
		k.handlers[name] = handler
	}
}

func NoOp() HandlerFunc {
	return func(c *Context) {
		c.Next()
	}
}

type KitItem interface {
	Order() int
	Apply(*Kit)
}

type kitItem struct {
	order int
	fn    func(*Kit)
}

func NewKitItem(order int, fn func(*Kit)) KitItem {
	return kitItem{order: order, fn: fn}
}

func (i kitItem) Order() int {
	return i.order
}

func (i kitItem) Apply(kit *Kit) {
	if i.fn != nil {
		i.fn(kit)
	}
}

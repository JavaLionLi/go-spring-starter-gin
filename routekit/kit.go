package routekit

import (
	"sort"

	"github.com/gin-gonic/gin"
)

type Registrar interface {
	Order() int
	RegisterRoutes(*gin.Engine, Kit)
}

type registrar struct {
	order int
	fn    func(*gin.Engine, Kit)
}

func NewRegistrar(order int, fn func(*gin.Engine, Kit)) Registrar {
	return registrar{order: order, fn: fn}
}

func (r registrar) Order() int {
	return r.order
}

func (r registrar) RegisterRoutes(engine *gin.Engine, kit Kit) {
	if r.fn != nil {
		r.fn(engine, kit)
	}
}

func RegisterAll(engine *gin.Engine, kit Kit, registrars []Registrar) {
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
	handlers map[string]gin.HandlerFunc
}

func NewKit(items []KitItem) Kit {
	kit := Kit{handlers: map[string]gin.HandlerFunc{}}
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

func (k Kit) Handler(name string) gin.HandlerFunc {
	if k.handlers == nil {
		return NoOp()
	}
	if handler := k.handlers[name]; handler != nil {
		return handler
	}
	return NoOp()
}

func (k *Kit) SetHandler(name string, handler gin.HandlerFunc) {
	if k.handlers == nil {
		k.handlers = map[string]gin.HandlerFunc{}
	}
	if name != "" && handler != nil {
		k.handlers[name] = handler
	}
}

func NoOp() gin.HandlerFunc {
	return func(c *gin.Context) {
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

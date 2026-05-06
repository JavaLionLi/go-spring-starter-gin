package startergin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-spring/spring-core/gs"
	"github.com/go-spring/stdlib/flatten"
)

func init() {
	gs.Module(
		gs.OnProperty("spring.gin.enabled").HavingValue("true").MatchIfMissing(),
		func(r gs.BeanProvider, _ flatten.Storage) error {
			r.Provide(&EngineDeps{})
			r.Provide(
				NewEngine,
				gs.IndexArg(0, gs.TagArg("${spring.gin}")),
			).Condition(gs.OnMissingBean[*gin.Engine]())
			r.Provide(NewHTTPServeMux).Condition(gs.OnMissingBean[*gs.HttpServeMux]())
			return nil
		},
	)
}

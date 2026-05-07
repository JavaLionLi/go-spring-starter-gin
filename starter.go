package startergin

import (
	"github.com/go-spring/spring-core/gs"
	"github.com/go-spring/stdlib/flatten"
)

func init() {
	// 通过 go-spring 模块机制按配置开关注册 Gin starter。
	// spring.gin.enabled 未配置时默认启用，显式设置为 false 时跳过整个模块。
	gs.Module(
		gs.OnProperty("spring.gin.enabled").HavingValue("true").MatchIfMissing(),
		func(r gs.BeanProvider, _ flatten.Storage) error {
			// EngineDeps 作为聚合 Bean 承接可选扩展点，NewEngine 会读取其中的切片。
			r.Provide(&EngineDeps{})
			r.Provide(
				NewEngine,
				gs.IndexArg(0, gs.TagArg("${spring.gin}")),
			).Condition(gs.OnMissingBean[*Engine]())
			// 若业务侧没有提供 HttpServeMux，则使用 Gin Engine 作为 HTTP Handler。
			r.Provide(NewHTTPServeMux).Condition(gs.OnMissingBean[*gs.HttpServeMux]())
			return nil
		},
	)
}

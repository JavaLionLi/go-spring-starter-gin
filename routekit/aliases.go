package routekit

import "github.com/gin-gonic/gin"

// 这些类型别名让 routekit 的注册器可以直接表达 Gin 路由签名，
// 调用方无需在 routekit 扩展代码中额外导入 gin 包。
type Accounts = gin.Accounts
type ContextKeyType = gin.ContextKeyType
type Engine = gin.Engine
type Context = gin.Context
type Error = gin.Error
type ErrorType = gin.ErrorType
type HandlersChain = gin.HandlersChain
type HandlerFunc = gin.HandlerFunc
type H = gin.H
type IRouter = gin.IRouter
type IRoutes = gin.IRoutes
type LogFormatter = gin.LogFormatter
type LogFormatterParams = gin.LogFormatterParams
type LoggerConfig = gin.LoggerConfig
type Negotiate = gin.Negotiate
type OnlyFilesFS = gin.OnlyFilesFS
type OptionFunc = gin.OptionFunc
type Param = gin.Param
type Params = gin.Params
type RecoveryFunc = gin.RecoveryFunc
type ResponseWriter = gin.ResponseWriter
type RouteInfo = gin.RouteInfo
type RouterGroup = gin.RouterGroup
type RoutesInfo = gin.RoutesInfo
type Skipper = gin.Skipper

const (
	// 常用 MIME 类型常量，保持与 gin 包定义一致。
	MIMEJSON              = gin.MIMEJSON
	MIMEHTML              = gin.MIMEHTML
	MIMEXML               = gin.MIMEXML
	MIMEXML2              = gin.MIMEXML2
	MIMEPlain             = gin.MIMEPlain
	MIMEPOSTForm          = gin.MIMEPOSTForm
	MIMEMultipartPOSTForm = gin.MIMEMultipartPOSTForm
	MIMEYAML              = gin.MIMEYAML
	MIMEYAML2             = gin.MIMEYAML2
	MIMETOML              = gin.MIMETOML
	MIMEPROTOBUF          = gin.MIMEPROTOBUF
	MIMEBSON              = gin.MIMEBSON

	// Gin 在可信平台代理场景下使用的平台标识。
	PlatformGoogleAppEngine = gin.PlatformGoogleAppEngine
	PlatformCloudflare      = gin.PlatformCloudflare
	PlatformFlyIO           = gin.PlatformFlyIO

	// Gin 运行模式常量。
	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
	TestMode    = gin.TestMode

	// Gin 在上下文、绑定和运行时模式中使用的标准 key。
	AuthProxyUserKey = gin.AuthProxyUserKey
	AuthUserKey      = gin.AuthUserKey
	BindKey          = gin.BindKey
	BodyBytesKey     = gin.BodyBytesKey
	ContextKey       = gin.ContextKey
	EnvGinMode       = gin.EnvGinMode
	Version          = gin.Version

	ContextRequestKey = gin.ContextRequestKey

	// Gin 错误分类常量，用于按类型处理错误。
	ErrorTypeBind    = gin.ErrorTypeBind
	ErrorTypeRender  = gin.ErrorTypeRender
	ErrorTypePrivate = gin.ErrorTypePrivate
	ErrorTypePublic  = gin.ErrorTypePublic
	ErrorTypeAny     = gin.ErrorTypeAny
)

// New 创建一个未默认挂载中间件的 Gin Engine。
func New(opts ...OptionFunc) *Engine {
	return gin.New(opts...)
}

// SetMode 设置 Gin 运行模式，例如 debug、release 或 test。
func SetMode(value string) {
	gin.SetMode(value)
}
